package dynamo

import (
	"strings"

	"git.devops.com/go/odm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func init() {
	dialect := &dynamoDialect{}
	odm.RegisterDialect("dynamo", dialect)
	odm.RegisterDialect("dynamodb", dialect)
}

type dynamoDialect struct {
}

func (d *dynamoDialect) Open(connectString string) (odm.DialectDB, error) {
	cfg, err := ParseConnectString(connectString)
	if err != nil {
		return nil, err
	}
	return OpenDB(cfg)
}

func (d *dynamoDialect) GetName() string {
	return "dynamodb"
}

func ParseConnectString(connectString string) (*aws.Config, error) {
	parts := strings.Split(connectString, ";")
	cfg := &aws.Config{}
	var accessKey, secretKey, token string
	for _, part := range parts {
		part = strings.Trim(part, " ")
		if part != "" {
			kv := strings.SplitN(part, "=", 2)
			v := kv[1]
			switch strings.ToLower(kv[0]) {
			case "accesskey":
				accessKey = v
				break
			case "secretkey":
				secretKey = v
			case "region":
				cfg.Region = aws.String(v)
			case "token":
				token = v
				break
			case "endpoint":
				cfg.Endpoint = aws.String(v)
				break
			default:
			}
		}
	}
	if accessKey != "" && secretKey != "" {
		cfg.Credentials = credentials.NewStaticCredentials(accessKey, secretKey, token)
	}
	return cfg, nil
}

// Open dynamodb client
func openConnection(cfg *aws.Config) (*dynamodb.DynamoDB, error) {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	return svc, nil
}

func OpenDB(cfg *aws.Config) (*DB, error) {
	conn, err := openConnection(cfg)
	if err != nil {
		return nil, err
	}
	db := &DB{
		conn:                conn,
		enableTableCreation: *cfg.Region == "localhost",
		tableMap:            make(map[string]*Table),
		tableMetaMap:        make(map[string]*odm.TableMeta),
	}
	return db, nil
}

type DB struct {
	conn *dynamodb.DynamoDB
	// if this is true, then auto create table if not exists.
	enableTableCreation bool
	// cache for Describe table
	// TODO: what if table changed while running?
	tableMetaMap map[string]*odm.TableMeta
	// cache for Table
	tableMap map[string]*Table
}

func (db *DB) createTableIfNotExists(meta *odm.TableMeta) error {
	_meta, err := db.GetTableMeta(meta.TableName)
	if err != nil {
		if aerr, ok := err.(awserr.Error); !ok || aerr.Code() != dynamodb.ErrCodeResourceNotFoundException {
			return err
		}
	}
	if _meta != nil {
		// table exists
		return nil
	}
	return db.createTable(meta)
}

func (db *DB) createTable(tableMeta *odm.TableMeta) error {
	conn := db.GetConn()
	keySchema := []*dynamodb.KeySchemaElement{
		&dynamodb.KeySchemaElement{
			AttributeName: aws.String(tableMeta.PartitionKey),
			KeyType:       aws.String("HASH"),
		},
	}
	if tableMeta.SortingKey != "" {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{
			AttributeName: aws.String(tableMeta.SortingKey),
			KeyType:       aws.String("RANGE"),
		})
	}
	out, err := conn.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(tableMeta.TableName),
		KeySchema: keySchema,
	})
	if out != nil {
		db.tableMetaMap[tableMeta.TableName] = db.updateTableDescription(out.TableDescription)
	}
	return err
}

func (db *DB) updateTableDescription(tableDesc *dynamodb.TableDescription) *odm.TableMeta {
	meta := &odm.TableMeta{
		TableName: *tableDesc.TableName,
	}
	for _, key := range tableDesc.KeySchema {
		if *key.KeyType == "HASH" {
			meta.PartitionKey = *key.AttributeName
		} else if *key.KeyType == "RANGE" {
			meta.SortingKey = *key.AttributeName
		}
	}
	db.tableMetaMap[*tableDesc.TableName] = meta
	// result.Table.LocalSecondaryIndexes
	// result.Table.GlobalSecondaryIndexes
	return meta
}

func (db *DB) GetTableMeta(tableName string) (*odm.TableMeta, error) {
	meta := db.tableMetaMap[tableName]
	if meta != nil {
		return meta, nil
	}
	conn := db.GetConn()
	result, err := conn.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if result != nil {
		meta = db.updateTableDescription(result.Table)
		db.tableMetaMap[tableName] = meta
	}
	return meta, err
}

func (db *DB) GetConn() *dynamodb.DynamoDB {
	return db.conn
}

func (db *DB) GetTable(name string) odm.Table {
	table := new(Table)
	table.db = db
	table.TableName = name
	return table
}

func (db *DB) GetDialectTable(meta *odm.TableMeta) odm.Table {
	return &Table{
		db:        db,
		TableMeta: *meta,
	}
}

func (db *DB) Close() {
	// Nothing to do.
}

func (db *DB) BatchGetItem(options []odm.BatchGet, unprocessedItems *[]odm.BatchGet, results ...interface{}) error {
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{},
	}
	var resultsMap map[string]interface{}
	for i, opt := range options {
		if opt.TableName == "" {
			panic("BatchGetItem TableName is required, but got nil.")
		}
		if input.RequestItems[opt.TableName] != nil {
			panic("BatchGetItem options TableName <" + opt.TableName + "> duplicated. ")
		}
		if opt.Keys == nil || len(opt.Keys) == 0 {
			continue
		}
		optIn := &dynamodb.KeysAndAttributes{}
		if opt.Consistent {
			optIn.ConsistentRead = aws.Bool(opt.Consistent)
		}
		if opt.Select != "" {
			optIn.ProjectionExpression = aws.String(opt.Select)
		}
		if opt.NameParams != nil && len(opt.NameParams) > 0 {
			convertAttributeNames(opt.NameParams, optIn.ExpressionAttributeNames)
		}
		for _, key := range opt.Keys {
			keyMap, err := dynamodbattribute.MarshalMap(key)
			if err != nil {
				return err
			}
			optIn.Keys = append(optIn.Keys, keyMap)
		}
		input.RequestItems[opt.TableName] = optIn
		resultsMap[opt.TableName] = results[i]
	}
	out, err := db.GetConn().BatchGetItem(input)
	// Handle output
	for tableName, items := range out.Responses {
		err = dynamodbattribute.UnmarshalListOfMaps(items, resultsMap[tableName])
		if err != nil {
			return err
		}
	}
	// Handle UnprocessedKyes
	// 测试时以 input.RequestItem 代替
	for tableName, requestItem := range out.UnprocessedKeys {
		rawItem := odm.BatchGet{
			TableName:  tableName,
			NameParams: map[string]string{},
			Keys:       []odm.Key{},
		}
		if requestItem.ConsistentRead != nil {
			rawItem.Consistent = *requestItem.ConsistentRead
		}
		if requestItem.ExpressionAttributeNames != nil {
			revertAttributeNames(rawItem.NameParams, requestItem.ExpressionAttributeNames)
		}
		if requestItem.ProjectionExpression != nil {
			rawItem.Select = *requestItem.ProjectionExpression
		}
		if requestItem.Keys != nil {
			for _, keyMap := range requestItem.Keys {
				key := make(odm.Key)
				err = dynamodbattribute.UnmarshalMap(keyMap, key)
				if err != nil {
					return err
				}
				rawItem.Keys = append(rawItem.Keys, key)
			}
		}
		*unprocessedItems = append(*unprocessedItems, rawItem)
	}
	return err
}

func (db *DB) BatchWriteItem(options []odm.BatchWrite, unprocessedItems *[]odm.BatchWrite) error {
	panic("not implemented") // TODO: Implement
}

func (db *DB) TransactGetItems(gets []odm.TransGet, results ...odm.Model) error {
	panic("not implemented") // TODO: Implement
}

func (db *DB) TransactWriteItems(writes []odm.TransWrite) error {
	panic("not implemented") // TODO: Implement
}

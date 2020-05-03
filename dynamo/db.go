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

var dbName = "dynamodb"

func init() {
	dialect := &dynamoDialect{}
	odm.RegisterDialect("dynamo", dialect)
	odm.RegisterDialect(dbName, dialect)
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
	return dbName
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
		enableTableDeletion: *cfg.Region == "localhost",
		tableMap:            make(map[string]*Table),
		tableMetaMap:        make(map[string]*odm.TableMeta),
	}
	return db, nil
}

type DB struct {
	conn *dynamodb.DynamoDB
	// if this is true, then auto create table if not exists.
	enableTableCreation bool
	enableTableDeletion bool
	// cache for Describe table
	// TODO: what if table changed while running?
	tableMetaMap map[string]*odm.TableMeta
	// cache for Table
	tableMap map[string]*Table
}

// DropTable only allowed on localhost
func (db *DB) DropTable(tableName string) error {
	if !db.enableTableDeletion {
		panic("DropTable is not allowed")
	}
	conn := db.GetConn()
	_, err := conn.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	return err
}

func (db *DB) CreateTableIfNotExists(meta *odm.TableMeta) error {
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
	return db.CreateTable(meta)
}

func (db *DB) getFieldName(f *odm.FieldDefine) string {
	return f.GetDBFieldName(dbName)
}

func (db *DB) CreateTable(tableMeta *odm.TableMeta) error {
	conn := db.GetConn()
	// Key definition.
	keySchema := []*dynamodb.KeySchemaElement{
		&dynamodb.KeySchemaElement{
			AttributeName: aws.String(db.getFieldName(tableMeta.PK)),
			KeyType:       aws.String("HASH"),
		},
	}
	// AttributeDefinitions
	attrs := []*dynamodb.AttributeDefinition{
		&dynamodb.AttributeDefinition{
			AttributeName: aws.String(db.getFieldName(tableMeta.PK)),
			AttributeType: aws.String(tableMeta.PK.Type),
		},
	}
	if tableMeta.SK != nil {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{
			AttributeName: aws.String(db.getFieldName(tableMeta.SK)),
			KeyType:       aws.String("RANGE"),
		})
		attrs = append(attrs, &dynamodb.AttributeDefinition{
			AttributeName: aws.String(db.getFieldName(tableMeta.SK)),
			AttributeType: aws.String(tableMeta.SK.Type),
		})
	}
	// for _, f := range tableMeta.Fields {
	// 	attrs = append(attrs, &dynamodb.AttributeDefinition{
	// 		AttributeName: aws.String(db.getFieldName(f)),
	// 		AttributeType: aws.String(f.Type),
	// 	})
	// }
	// TODO: GSI
	// TODO: LSI
	out, err := conn.CreateTable(&dynamodb.CreateTableInput{
		TableName:   aws.String(tableMeta.TableName),
		KeySchema:   keySchema,
		BillingMode: aws.String("PAY_PER_REQUEST"), // PAY_PER_REQUEST, PROVISIONED
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1000),
			WriteCapacityUnits: aws.Int64(1000),
		},
		// GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{},
		// LocalSecondaryIndexes:  []*dynamodb.LocalSecondaryIndex{},
		AttributeDefinitions: attrs,
	})
	if err == nil && out != nil && out.TableDescription != nil {
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
			meta.PK = &odm.FieldDefine{
				SchemaFieldName: map[string]string{
					dbName: *key.AttributeName,
				},
			}
		} else if *key.KeyType == "RANGE" {
			meta.SK = &odm.FieldDefine{
				SchemaFieldName: map[string]string{
					dbName: *key.AttributeName,
				},
			}
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
	if err == nil && result != nil && result.Table != nil {
		meta = db.updateTableDescription(result.Table)
		db.tableMetaMap[tableName] = meta
	}
	return meta, err
}

func (db *DB) GetConn() *dynamodb.DynamoDB {
	return db.conn
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

func (db *DB) BatchGetItem(options []*odm.BatchGet, unprocessedItems *[]*odm.BatchGet, results ...interface{}) error {
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
			Keys:       []odm.Map{},
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
				key := make(odm.Map)
				err = dynamodbattribute.UnmarshalMap(keyMap, key)
				if err != nil {
					return err
				}
				rawItem.Keys = append(rawItem.Keys, key)
			}
		}
		*unprocessedItems = append(*unprocessedItems, &rawItem)
	}
	return err
}

func (db *DB) BatchWriteItem(options []*odm.BatchWrite, unprocessedItems *[]*odm.BatchWrite) error {
	panic("not implemented") // TODO: Implement
}

func (db *DB) TransactGetItems(gets []*odm.TransactGet, results ...odm.Model) error {
	panic("not implemented") // TODO: Implement
}

func (db *DB) TransactWriteItems(writes []*odm.TransactWrite) error {
	items := []*dynamodb.TransactWriteItem{}
	for _, write := range writes {
		item := &dynamodb.TransactWriteItem{}
		if write.ConditionCheck != nil {

		}
		if write.Delete != nil {

		}
		if write.Update != nil {

		}
		if write.Put != nil {

		}
		items = append(items, item)
	}
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	}
	// input.ClientRequestToken = aws.String("")
	_, err := db.GetConn().TransactWriteItems(input)
	return err
}

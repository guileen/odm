package dynamo

import (
	"strings"

	"git.devops.com/go/odm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func init() {
	dialect := &DynamodbDialectDB{}
	odm.RegisterDialect("dynamo", dialect)
	odm.RegisterDialect("dynamodb", dialect)
}

type DynamodbDialectDB struct {
}

func (d *DynamodbDialectDB) Open(connectString string) (odm.DialectDB, error) {
	cfg, err := ParseConnectString(connectString)
	if err != nil {
		return nil, err
	}
	return OpenDB(cfg)
}

func (d *DynamodbDialectDB) GetName() string {
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
	db := new(DB)
	db.conn = conn
	return db, nil
}

type DB struct {
	conn *dynamodb.DynamoDB
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

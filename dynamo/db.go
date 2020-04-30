package dynamo

import (
	"git.devops.com/go/odm/meta"
	"git.devops.com/go/odm/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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

func (db *DB) Table(model types.Model) types.Table {
	modelMeta := meta.GetModelMeta(model)
	table := new(Table)
	table.conn = db.conn
	table.TableMeta = *modelMeta
	return table
}

func (db *DB) GetTable(name string) types.Table {
	table := new(Table)
	table.conn = db.conn
	table.TableName = name
	return table
}

func (db *DB) Close() {
	// Nothing to do.
}

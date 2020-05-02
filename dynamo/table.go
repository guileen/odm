package dynamo

import (
	"errors"
	"fmt"

	"git.devops.com/go/odm"
	"git.devops.com/go/odm/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Table of DynamoDB implemetation.
type Table struct {
	odm.TableMeta
	db *DB
}

// GetDB of current table
func (t *Table) GetDB() odm.DialectDB {
	return t.db
}

// GetConn the Connection
func (t *Table) GetConn() (*dynamodb.DynamoDB, error) {
	fmt.Println("GetConn", t.TableMeta)
	if t.PartitionKey == "" {
		// TableMeta not initialized. 使用数据库来初始化
		meta, err := t.db.GetTableMeta(t.TableName)
		if err != nil {
			return nil, err
		}
		t.TableMeta = *meta
	} else {
		fmt.Println("createTableIfNotExists", t.TableMeta)
		err := t.db.createTableIfNotExists(&t.TableMeta)
		if err != nil {
			return nil, err
		}
	}
	return t.db.GetConn(), nil
}

func convertAttributeNames(params map[string]string, targetMap map[string]*string) {
	for k, v := range params {
		targetMap[k] = aws.String(v)
	}
}

func revertAttributeNames(params map[string]string, attrNames map[string]*string) {
	for k, v := range attrNames {
		params[k] = *v
	}
}

// PutItem put a item, will replace entire item.
func (t *Table) PutItem(item odm.Model, cond *odm.WriteOption) error {
	conn, err := t.GetConn()
	if err != nil {
		return err
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}
	// Create item in table Movies
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(t.TableName),
	}
	if cond != nil {
		if cond.ValueParams != nil {
			input.ExpressionAttributeValues, err = dynamodbattribute.MarshalMap(cond.ValueParams)
			if err != nil {
				return err
			}
		}
		if cond.Filter != "" {
			input.ConditionExpression = aws.String(cond.Filter)
		}
		if cond.NameParams != nil {
			input.ExpressionAttributeNames = make(map[string]*string)
			convertAttributeNames(cond.NameParams, input.ExpressionAttributeNames)
		}
	}
	_, err = conn.PutItem(input)
	return err
}

// UpdateItem attributes. item will fill base on ReturnValues.
func (t *Table) UpdateItem(key odm.Key, updateExpression string, cond *odm.WriteOption, result odm.Model) error {
	conn, err := t.GetConn()
	if err != nil {
		return err
	}
	keyMap, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return err
	}
	input := &dynamodb.UpdateItemInput{
		TableName:        aws.String(t.TableName),
		Key:              keyMap,
		UpdateExpression: aws.String(updateExpression),
	}
	if cond != nil {
		if cond.ValueParams != nil {
			input.ExpressionAttributeValues, err = dynamodbattribute.MarshalMap(cond.ValueParams)
			if err != nil {
				return err
			}
		}
		if cond.Filter != "" {
			input.ConditionExpression = aws.String(cond.Filter)
		}
		if cond.NameParams != nil {
			input.ExpressionAttributeNames = make(map[string]*string)
			convertAttributeNames(cond.NameParams, input.ExpressionAttributeNames)
		}
		if result != nil {
			input.ReturnValues = aws.String("UPDATED_NEW")
		}
	}
	out, err := conn.UpdateItem(input)
	if result != nil && err == nil {
		_ = dynamodbattribute.UnmarshalMap(out.Attributes, result)
	}
	return err
}

// GetItem get an item
func (t *Table) GetItem(key odm.Key, opt *odm.GetOption, item odm.Model) error {
	conn, err := t.GetConn()
	if err != nil {
		return err
	}
	keyMap, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return err
	}
	input := &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: aws.String(t.TableName),
	}
	if opt != nil {
		if opt.Consistent {
			input.ConsistentRead = aws.Bool(opt.Consistent)
		}
		if opt.Select != "" {
			input.ProjectionExpression = aws.String(opt.Select)
		}
		if opt.NameParams != nil {
			convertAttributeNames(opt.NameParams, input.ExpressionAttributeNames)
		}
	}
	result, err := conn.GetItem(input)
	if err != nil {
		return err
	}
	if item != nil && result != nil && result.Item != nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, item)
	}
	return err
}

// DeleteItem returns deleted item if item provide
func (t *Table) DeleteItem(key odm.Key, cond *odm.WriteOption, result odm.Model) error {
	conn, err := t.GetConn()
	if err != nil {
		return err
	}
	keyMap, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return err
	}
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(t.TableName),
		Key:       keyMap,
	}
	if cond != nil {
		if cond.ValueParams != nil {
			input.ExpressionAttributeValues, err = dynamodbattribute.MarshalMap(cond.ValueParams)
			if err != nil {
				return err
			}
		}
		if cond.Filter != "" {
			input.ConditionExpression = aws.String(cond.Filter)
		}
		if cond.NameParams != nil {
			input.ExpressionAttributeNames = make(map[string]*string)
			convertAttributeNames(cond.NameParams, input.ExpressionAttributeNames)
		}
		if result != nil {
			input.ReturnValues = aws.String("ALL_OLD")
		}
	}
	out, err := conn.DeleteItem(input)
	if result != nil && err == nil {
		_ = dynamodbattribute.UnmarshalMap(out.Attributes, result)
	}
	return err
}

func (t *Table) Scan(query *odm.QueryOption, offsetKey odm.Key, items interface{}) error {
	panic("Not implement scan")
}

// Query and fill in items, StartKey will be replaced after query
func (t *Table) Query(query *odm.QueryOption, offsetKey odm.Key, items interface{}) error {
	if query == nil {
		return errors.New("QueryOptions is required for Table.Query, ")
	}
	if query.KeyFilter == "" {
		return t.Scan(query, offsetKey, items)
	}
	conn, err := t.GetConn()
	if err != nil {
		return err
	}
	input := &dynamodb.QueryInput{
		TableName:              aws.String(t.TableName),
		KeyConditionExpression: aws.String(query.KeyFilter),
	}
	if offsetKey != nil && len(offsetKey) > 0 {
		input.ExclusiveStartKey, err = dynamodbattribute.MarshalMap(offsetKey)
		if err != nil {
			return err
		}
	}
	if query.ValueParams != nil {
		input.ExpressionAttributeValues, err = dynamodbattribute.MarshalMap(query.ValueParams)
		if err != nil {
			return err
		}
	}
	if query.NameParams != nil {
		input.ExpressionAttributeNames = make(map[string]*string)
		convertAttributeNames(query.NameParams, input.ExpressionAttributeNames)
	}
	if query.Filter != "" {
		input.FilterExpression = aws.String(query.Filter)
	}
	if query.Select != "" {
		input.ProjectionExpression = aws.String(query.Select)
	}
	if query.Consistent {
		input.ConsistentRead = aws.Bool(query.Consistent)
	}
	if query.Limit != 0 {
		input.Limit = aws.Int64(query.Limit)
	}
	if query.Desc {
		input.ScanIndexForward = aws.Bool(false)
	}
	if query.IndexName != "" {
		input.IndexName = aws.String(query.IndexName)
	}
	out, err := conn.Query(input)
	if err != nil {
		return fmt.Errorf("Fail to execute Query on %s. %w", t.TableName, err)
	}
	if out == nil {
		util.ClearSlice(items)
	} else {
		err = dynamodbattribute.UnmarshalListOfMaps(out.Items, items)
		if offsetKey != nil && err == nil {
			err = dynamodbattribute.UnmarshalMap(out.LastEvaluatedKey, &offsetKey)
		}
	}
	return err
}

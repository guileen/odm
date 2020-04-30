package dynamo

import (
	"errors"

	"git.devops.com/go/odm/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Table of DynamoDB implemetation.
type Table struct {
	types.TableMeta
	db *DB
}

// GetDB of current table
func (t *Table) GetDB() types.DB {
	return t.db
}

// GetConn the Connection
func (t *Table) GetConn() *dynamodb.DynamoDB {
	return t.db.GetConn()
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
func (t *Table) PutItem(item types.Model, cond *types.Condition) error {
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
		if cond.ConditionExpression != "" {
			input.ConditionExpression = aws.String(cond.ConditionExpression)
		}
		if cond.NameParams != nil {
			input.ExpressionAttributeNames = make(map[string]*string)
			convertAttributeNames(cond.NameParams, input.ExpressionAttributeNames)
		}
	}
	_, err = t.GetConn().PutItem(input)
	return err
}

// UpdateItem attributes. item will fill base on ReturnValues.
func (t *Table) UpdateItem(key types.Key, updateExpression string, cond *types.Condition, result types.Model) error {
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
		if cond.ConditionExpression != "" {
			input.ConditionExpression = aws.String(cond.ConditionExpression)
		}
		if cond.NameParams != nil {
			input.ExpressionAttributeNames = make(map[string]*string)
			convertAttributeNames(cond.NameParams, input.ExpressionAttributeNames)
		}
		if result != nil {
			input.ReturnValues = aws.String("UPDATED_NEW")
		}
	}
	out, err := t.GetConn().UpdateItem(input)
	if result != nil && err == nil {
		_ = dynamodbattribute.UnmarshalMap(out.Attributes, result)
	}
	return err
}

// GetItem get an item
func (t *Table) GetItem(key types.Key, opt *types.GetOption, item types.Model) error {
	keyMap, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return err
	}
	input := &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: aws.String(t.TableName),
	}
	if opt != nil {
		if opt.ConsistentRead {
			input.ConsistentRead = aws.Bool(opt.ConsistentRead)
		}
		if opt.ProjectionExpression != "" {
			input.ProjectionExpression = aws.String(opt.ProjectionExpression)
		}
		if opt.NameParams != nil {
			convertAttributeNames(opt.NameParams, input.ExpressionAttributeNames)
		}
	}
	result, err := t.GetConn().GetItem(input)
	if err != nil {
		return err
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	return err
}

// DeleteItem returns deleted item if item provide
func (t *Table) DeleteItem(key types.Key, cond *types.Condition, result types.Model) error {
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
		if cond.ConditionExpression != "" {
			input.ConditionExpression = aws.String(cond.ConditionExpression)
		}
		if cond.NameParams != nil {
			input.ExpressionAttributeNames = make(map[string]*string)
			convertAttributeNames(cond.NameParams, input.ExpressionAttributeNames)
		}
		if result != nil {
			input.ReturnValues = aws.String("ALL_OLD")
		}
	}
	out, err := t.GetConn().DeleteItem(input)
	if result != nil && err == nil {
		_ = dynamodbattribute.UnmarshalMap(out.Attributes, result)
	}
	return err
}

func (t *Table) Scan(query *types.QueryOption, offsetKey types.Key, items interface{}) error {
	panic("Not implement scan")
}

// Query and fill in items, StartKey will be replaced after query
func (t *Table) Query(query *types.QueryOption, offsetKey types.Key, items interface{}) error {
	if query == nil {
		return errors.New("QueryOptions is required for Table.Query, ")
	}
	if query.KeyConditionExpression == "" {
		return t.Scan(query, offsetKey, items)
	}
	input := &dynamodb.QueryInput{
		TableName:              aws.String(t.TableName),
		KeyConditionExpression: aws.String(query.KeyConditionExpression),
	}
	var err error
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
	if query.FilterExpression != "" {
		input.FilterExpression = aws.String(query.FilterExpression)
	}
	if query.ProjectionExpression != "" {
		input.ProjectionExpression = aws.String(query.ProjectionExpression)
	}
	if query.ConsistentRead {
		input.ConsistentRead = aws.Bool(query.ConsistentRead)
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
	out, err := t.GetConn().Query(input)
	if out != nil && err == nil {
		err = dynamodbattribute.UnmarshalListOfMaps(out.Items, items)
		if offsetKey != nil && err == nil {
			err = dynamodbattribute.UnmarshalMap(out.LastEvaluatedKey, &offsetKey)
		}
	}
	return err
}

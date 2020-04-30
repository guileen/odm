
## DyanamoDB input options explain

// Determines the read consistency model: If set to true, then the operation
// uses strongly consistent reads; otherwise, the operation uses eventually
// consistent reads.
ConsistentRead bool

// One or more substitution tokens for attribute names in an expression. The
// following are some use cases for using ExpressionAttributeNames:
//
//    * To access an attribute whose name conflicts with a DynamoDB reserved
//    word.
//
//    * To create a placeholder for repeating occurrences of an attribute name
//    in an expression.
//
//    * To prevent special characters in an attribute name from being misinterpreted
//    in an expression.
//
// Use the # character in an expression to dereference an attribute name. For
// example, consider the following attribute name:
//
//    * Percentile
//
// The name of this attribute conflicts with a reserved word, so it cannot be
// used directly in an expression. (For the complete list of reserved words,
// see Reserved Words (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ReservedWords.html)
// in the Amazon DynamoDB Developer Guide). To work around this, you could specify
// the following for ExpressionAttributeNames:
//
//    * {"#P":"Percentile"}
//
// You could then use this substitution in an expression, as in this example:
//
//    * #P = :val
//
// Tokens that begin with the : character are expression attribute values, which
// are placeholders for the actual value at runtime.
//
// For more information on expression attribute names, see Specifying Item Attributes
// (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.AccessingItemAttributes.html)
// in the Amazon DynamoDB Developer Guide.
ExpressionAttributeNames map[string]*string

ProjectionExpression     *string `type:"string"`
// A map of attribute names to AttributeValue objects, representing the primary
// key of the item to retrieve.
//
// For the primary key, you must provide all of the attributes. For example,
// with a simple primary key, you only need to provide a value for the partition
// key. For a composite primary key, you must provide values for both the partition
// key and the sort key.
//
// Key is a required field
Key map[string]interface{} `type:"map" required:"true"`

// A condition that must be satisfied in order for a conditional update to succeed.
//
// An expression can contain any of the following:
//
//    * Functions: attribute_exists | attribute_not_exists | attribute_type
//    | contains | begins_with | size These function names are case-sensitive.
//
//    * Comparison operators: = | <> | < | > | <= | >= | BETWEEN | IN
//
//    * Logical operators: AND | OR | NOT
//
// For more information about condition expressions, see Specifying Conditions
// (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.SpecifyingConditions.html)
// in the Amazon DynamoDB Developer Guide.
ConditionExpression       *string                    `type:"string"`


// The primary key of the first item that this operation will evaluate. Use
// the value that was returned for LastEvaluatedKey in the previous operation.
//
// The data type for ExclusiveStartKey must be String, Number, or Binary. No
// set data types are allowed.
ExclusiveStartKey        map[string]*AttributeValue `type:"map"`

// One or more values that can be substituted in an expression.
//
// Use the : (colon) character in an expression to dereference an attribute
// value. For example, suppose that you wanted to check whether the value of
// the ProductStatus attribute was one of the following:
//
// Available | Backordered | Discontinued
//
// You would first need to specify ExpressionAttributeValues as follows:
//
// { ":avail":{"S":"Available"}, ":back":{"S":"Backordered"}, ":disc":{"S":"Discontinued"}
// }
//
// You could then use these values in an expression, such as this:
//
// ProductStatus IN (:avail, :back, :disc)
//
// For more information on expression attribute values, see Specifying Conditions
// (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.SpecifyingConditions.html)
// in the Amazon DynamoDB Developer Guide.
ExpressionAttributeValues map[string]*AttributeValue `type:"map"`

// A string that contains conditions that DynamoDB applies after the Query operation,
// but before the data is returned to you. Items that do not satisfy the FilterExpression
// criteria are not returned.
//
// A FilterExpression does not allow key attributes. You cannot define a filter
// expression based on a partition key or a sort key.
//
// A FilterExpression is applied after the items have already been read; the
// process of filtering does not consume any additional read capacity units.
//
// For more information, see Filter Expressions (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/QueryAndScan.html#FilteringResults)
// in the Amazon DynamoDB Developer Guide.
FilterExpression *string `type:"string"`

// The name of an index to query. This index can be any local secondary index
// or global secondary index on the table. Note that if you use the IndexName
// parameter, you must also provide TableName.
IndexName *string `min:"3" type:"string"`

// The condition that specifies the key values for items to be retrieved by
// the Query action.
//
// The condition must perform an equality test on a single partition key value.
//
// The condition can optionally perform one of several comparison tests on a
// single sort key value. This allows Query to retrieve one item with a given
// partition key value and sort key value, or several items that have the same
// partition key value but different sort key values.
//
// The partition key equality test is required, and must be specified in the
// following format:
//
// partitionKeyName = :partitionkeyval
//
// If you also want to provide a condition for the sort key, it must be combined
// using AND with the condition for the sort key. Following is an example, using
// the = comparison operator for the sort key:
//
// partitionKeyName = :partitionkeyval AND sortKeyName = :sortkeyval
//
// Valid comparisons for the sort key condition are as follows:
//
//    * sortKeyName = :sortkeyval - true if the sort key value is equal to :sortkeyval.
//
//    * sortKeyName < :sortkeyval - true if the sort key value is less than
//    :sortkeyval.
//
//    * sortKeyName <= :sortkeyval - true if the sort key value is less than
//    or equal to :sortkeyval.
//
//    * sortKeyName > :sortkeyval - true if the sort key value is greater than
//    :sortkeyval.
//
//    * sortKeyName >= :sortkeyval - true if the sort key value is greater than
//    or equal to :sortkeyval.
//
//    * sortKeyName BETWEEN :sortkeyval1 AND :sortkeyval2 - true if the sort
//    key value is greater than or equal to :sortkeyval1, and less than or equal
//    to :sortkeyval2.
//
//    * begins_with ( sortKeyName, :sortkeyval ) - true if the sort key value
//    begins with a particular operand. (You cannot use this function with a
//    sort key that is of type Number.) Note that the function name begins_with
//    is case-sensitive.
//
// Use the ExpressionAttributeValues parameter to replace tokens such as :partitionval
// and :sortval with actual values at runtime.
//
// You can optionally use the ExpressionAttributeNames parameter to replace
// the names of the partition key and sort key with placeholder tokens. This
// option might be necessary if an attribute name conflicts with a DynamoDB
// reserved word. For example, the following KeyConditionExpression parameter
// causes an error because Size is a reserved word:
//
//    * Size = :myval
//
// To work around this, define a placeholder (such a #S) to represent the attribute
// name Size. KeyConditionExpression then is as follows:
//
//    * #S = :myval
//
// For a list of reserved words, see Reserved Words (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ReservedWords.html)
// in the Amazon DynamoDB Developer Guide.
//
// For more information on ExpressionAttributeNames and ExpressionAttributeValues,
// see Using Placeholders for Attribute Names and Values (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ExpressionPlaceholders.html)
// in the Amazon DynamoDB Developer Guide.
KeyConditionExpression *string `type:"string"`

// Specifies the order for index traversal: If true (default), the traversal
// is performed in ascending order; if false, the traversal is performed in
// descending order.
//
// Items with the same partition key value are stored in sorted order by sort
// key. If the sort key data type is Number, the results are stored in numeric
// order. For type String, the results are stored in order of UTF-8 bytes. For
// type Binary, DynamoDB treats each byte of the binary data as unsigned.
//
// If ScanIndexForward is true, DynamoDB returns the results in the order in
// which they are stored (by sort key value). This is the default behavior.
// If ScanIndexForward is false, DynamoDB reads the results in reverse order
// by sort key value, and then returns the results to the client.
ScanIndexForward *bool `type:"boolean"`
// The attributes to be returned in the result. You can retrieve all item attributes,
// specific item attributes, the count of matching items, or in the case of
// an index, some or all of the attributes projected into the index.
//
//    * ALL_ATTRIBUTES - Returns all of the item attributes from the specified
//    table or index. If you query a local secondary index, then for each matching
//    item in the index, DynamoDB fetches the entire item from the parent table.
//    If the index is configured to project all item attributes, then all of
//    the data can be obtained from the local secondary index, and no fetching
//    is required.
//
//    * ALL_PROJECTED_ATTRIBUTES - Allowed only when querying an index. Retrieves
//    all attributes that have been projected into the index. If the index is
//    configured to project all attributes, this return value is equivalent
//    to specifying ALL_ATTRIBUTES.
//
//    * COUNT - Returns the number of matching items, rather than the matching
//    items themselves.
//
//    * SPECIFIC_ATTRIBUTES - Returns only the attributes listed in AttributesToGet.
//    This return value is equivalent to specifying AttributesToGet without
//    specifying any value for Select. If you query or scan a local secondary
//    index and request only attributes that are projected into that index,
//    the operation will read only the index and not the table. If any of the
//    requested attributes are not projected into the local secondary index,
//    DynamoDB fetches each of these attributes from the parent table. This
//    extra fetching incurs additional throughput cost and latency. If you query
//    or scan a global secondary index, you can only request attributes that
//    are projected into the index. Global secondary index queries cannot fetch
//    attributes from the parent table.
//
// If neither Select nor AttributesToGet are specified, DynamoDB defaults to
// ALL_ATTRIBUTES when accessing a table, and ALL_PROJECTED_ATTRIBUTES when
// accessing an index. You cannot use both Select and AttributesToGet together
// in a single request, unless the value for Select is SPECIFIC_ATTRIBUTES.
// (This usage is equivalent to specifying AttributesToGet without any value
// for Select.)
//
// If you use the ProjectionExpression parameter, then the value for Select
// can only be SPECIFIC_ATTRIBUTES. Any other value for Select will return an
// error.
Select *string `type:"string" enum:"Select"`

// For a parallel Scan request, Segment identifies an individual segment to
// be scanned by an application worker.
//
// Segment IDs are zero-based, so the first segment is always 0. For example,
// if you want to use four application threads to scan a table or an index,
// then the first thread specifies a Segment value of 0, the second thread specifies
// 1, and so on.
//
// The value of LastEvaluatedKey returned from a parallel Scan request must
// be used as ExclusiveStartKey with the same segment ID in a subsequent Scan
// operation.
//
// The value for Segment must be greater than or equal to 0, and less than the
// value provided for TotalSegments.
//
// If you provide Segment, you must also provide TotalSegments.
Segment *int64  `type:"integer"`

// For a parallel Scan request, TotalSegments represents the total number of
// segments into which the Scan operation will be divided. The value of TotalSegments
// corresponds to the number of application workers that will perform the parallel
// scan. For example, if you want to use four application threads to scan a
// table or an index, specify a TotalSegments value of 4.
//
// The value for TotalSegments must be greater than or equal to 1, and less
// than or equal to 1000000. If you specify a TotalSegments value of 1, the
// Scan operation will be sequential rather than parallel.
//
// If you specify TotalSegments, you must also specify Segment.
TotalSegments *int64 `min:"1" type:"integer"`
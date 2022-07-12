# README

This library is a set of Go based AWS SDK V2 helper tools for DynamoDB.

If there is a better way of doing this, please send me a message and I'll adopt that. I did take a look making the best of my Google powers and didn't find anything.

__Installation__

```bash
# At the time of release, that is 0.1.0
go get github.com/davedotdev/ddbtools@latest
```

__SetEquals__

The SetEquals function makes it easier to create an expression update string and deals with adding AttributeValues to the update map.
It means your update queries become tider! An example of using it is below.


```go
import (
	"github.com/davedotdev/ddbtools
)

func (da *Data) UpdateSetting(firstName, lastName, userGUID string) error {

	client := dynamodb.NewFromConfig(da.DDBConfig, func(o *dynamodb.Options) {})

	// The entry needs to exist first and foremost!
	CondExpr := "attribute_exists(#SK)"

  // Boilerplate
	PKKey := "setting"
	SKKey := userGUID
	upd := dynamodb.UpdateItemInput{}
	upd.ConditionExpression = &CondExpr
	upd.TableName = &da.TableName
	upd.Key = map[string]types.AttributeValue{}
	upd.ExpressionAttributeNames = make(map[string]string)
	upd.ExpressionAttributeValues = make(map[string]types.AttributeValue)	
	upd.Key["PK"] = &types.AttributeValueMemberS{Value: PKKey}
	upd.Key["SK"] = &types.AttributeValueMemberS{Value: SKKey}
	upd.ExpressionAttributeNames["#SK"] = "SK"

  // The magic is here: handling each attribute
  // pass in the UpdateItemInput instance (upd) and the string value of exprStr
  // If the fourth input is true, then it creates a new string for exprStr
	exprStr := ddbtools.SetEquals(firstName, "firstName", "", true, &upd)
	exprStr = ddbtools.SetEquals(lastName, "lastName", exprStr, false, &upd)
	// This SetEquals func isn't too disruptive to normal patterns
  // and doesn't do any real 'magic', so it's easy to debug

	upd.UpdateExpression = &exprStr

  // Do the transaction
	_, err := client.UpdateItem(context.TODO(), &upd)

	if err != nil {
		return err
	}

	return nil
}
```

__IncrementValue__

Works in a very similar way, adding `1` to the Number type. Ensure the attribute is a number and not a string, else DynamoDB will error out.


```go
import (
	"github.com/davedotdev/ddbtools
)

func (da *Data) UpdateSetting(firstName, lastName, userGUID string) error {

	client := dynamodb.NewFromConfig(da.DDBConfig, func(o *dynamodb.Options) {})

	items := []types.TransactWriteItem{}

	// The entry needs to exist first and foremost!
	CondExpr := "attribute_exists(#SK)"

    // Boilerplate
	PKKey := "setting"
	SKKey := userGUID
	upd := types.TransactWriteItem{}
	upd.Update = &types.Update{}
	upd.Upate.ConditionExpression = &CondExpr
	upd.Upate.TableName = &da.TableName
	upd.Upate.Key = map[string]types.AttributeValue{}
	upd.Upate.ExpressionAttributeNames = make(map[string]string)
	upd.Upate.ExpressionAttributeValues = make(map[string]types.AttributeValue)	
	upd.Upate.Key["PK"] = &types.AttributeValueMemberS{Value: PKKey}
	upd.Upate.Key["SK"] = &types.AttributeValueMemberS{Value: SKKey}
	upd.Upate.ExpressionAttributeNames["#SK"] = "SK"

	// The magic is here: handling each attribute
	// pass in the Update instance (upd) and the string value of exprStr
	// If the fourth input is true, then it creates a new string for exprStr
	// This is the equiv of 'ADD #thingCount :1'
	exprStr := ddbtools.TransactIncrementValue("thingCount", "", true, &upd)
	upd.UpdateExpression = &exprStr	

	// Add upd to the transaction slice
	items = append(items, upd)

  	// Do the transaction
	_, err = client.TransactWriteItems(context.TODO(), &dynamodb.TransactWriteItemsInput{
	    TransactItems: items,
	})

	if err != nil {
		return err
	}

	return nil
}
```

There is also another version that does non-transaction based increment updates called `IncremenetValue`. The core difference is, it accepts a `dynamodb.UpdateItemInput` as the last argument.

## Support & Issues

Raise an issue on this repository.

## PRs

Welcome.
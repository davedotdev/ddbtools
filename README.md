# README

This library is a set of Go based AWS SDK V2 helper tools.

Currently it's one tool, but I suspect more will be added over time.

__SetEquals__

The SetEquals function makes it easier to create an expression update string and deals with adding AttributeValues to the update map.
It means your update queries become tider! An example of using it is below.

import (
	"github.com/davedotdev/ddbtools
)

```go
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

  // The magic is here: handling each property
  // pass in the UpdateItemInput instance (upd) and the string value of exprStr
  // If the fourth input is true, then it creates a new string for exprStr
	exprStr := ddbtools.SetEquals(firstName, "firstName", "", true, &upd)
	exprStr = ddbtools.SetEquals(lastName, "lastName", exprStr, false, &upd)
	upd.UpdateExpression = &exprStr

  // Do the transaction
	_, err := client.UpdateItem(context.TODO(), &upd)

	if err != nil {
		return err
	}

	return nil
}
```

## Support & Issues

Raise an issue on this repository.

## PRs

Welcome.
package ddbtools

// Common DynamoDB code

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// SetEquals accepts some input values, an inputName (the DDB attribute name), an expression string, whether it's the first entry and the DDB UpdateItemInput isntance
// The reason it needs 'first' is because it will create the expression string with SET. It could be auto detected, but this way means the SetEquals function
// has hybrid use cases and is easy for brownfield insertion.
func SetEquals(inputValue interface{}, inputName string, expr string, first bool, update *dynamodb.UpdateItemInput) string {

	ExprNameTemplate := "#%v"
	// equiv of "#BusinessName"
	ExprName := fmt.Sprintf(ExprNameTemplate, inputName)
	ExprNames := map[string]string{ExprName: inputName}

	ExprTemplate := "#%v = :%v"
	ExprStringAddition := fmt.Sprintf(ExprTemplate, inputName, inputName)

	exprStr := ""
	if first {
		exprStr = "SET "
		exprStr += ExprStringAddition
	} else {
		exprStr = expr
		exprStr += ", " + ExprStringAddition
	}

	for k, v := range ExprNames {
		update.ExpressionAttributeNames[k] = v
	}

	valName := ":" + inputName

	switch v := inputValue.(type) {
	case int, float32, float64:
		update.ExpressionAttributeValues[valName] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%v", v)}
	case string:
		update.ExpressionAttributeValues[valName] = &types.AttributeValueMemberS{Value: v}
	case bool:
		update.ExpressionAttributeValues[valName] = &types.AttributeValueMemberBOOL{Value: v}
	}

	return exprStr
}

// Incremement Value accepts an inputName (the DDB attribute name), an expression string, whether
// it's the first entry and the DDB UpdateItemInput isntance
// The reason it needs 'first' is because it will create the expression string with ADD pre-appended.
// It could be auto detected, but this way means it has hybrid use cases and is easy for brownfield insertion.
// This function has no intelligence to know if you're incrementing a bool attribute.
func TransactIncrementValue(inputName string, expr string, first bool, update *dynamodb.UpdateItemInput) string {

	ExprNameTemplate := "#%v"
	// equiv of "#BusinessName"
	ExprName := fmt.Sprintf(ExprNameTemplate, inputName)
	ExprNames := map[string]string{ExprName: inputName}

	ExprStringAddition := ExprName + " :inc" // "#thing :inc"

	exprStr := ""
	if first {
		exprStr = "ADD "
		exprStr += ExprStringAddition
	} else {
		exprStr = expr
		exprStr += ", " + ExprStringAddition
	}

	for k, v := range ExprNames {
		update.ExpressionAttributeNames[k] = v
	}

	update.ExpressionAttributeValues[":inc"] = &types.AttributeValueMemberN{Value: "1"}

	return exprStr
}

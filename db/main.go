package db

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var Instance *dynamodb.DynamoDB

func Init() {
	Instance = dynamodb.New(session.New())
}

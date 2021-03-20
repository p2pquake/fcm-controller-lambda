package tokens

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/p2pquake/fcm-controller-lambda/db"
)

const Prefix = "/v1/tokens"

var notFound = events.APIGatewayV2HTTPResponse{
	StatusCode: 404,
}
var pattern = regexp.MustCompile(Prefix + "/(quake|foreign|tsunami|userquake|eew)/([^/]*)$")
var tableMap = map[string]string{
	"quake":     "P2PQuakeMobilePushQuake",
	"foreign":   "P2PQuakeMobilePushForeign",
	"tsunami":   "P2PQuakeMobilePushTsunami",
	"userquake": "P2PQuakeMobilePushUserquake",
	"eew":       "P2PQuakeMobilePushEEW",
}

type createParams struct {
	UUID string `json:"uuid"`
	Min  int    `json:"min"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	matches := pattern.FindStringSubmatch(request.RequestContext.HTTP.Path)

	if matches == nil {
		return notFound, nil
	}

	switch request.RequestContext.HTTP.Method {
	case "PUT":
		return createOrUpdateToken(ctx, matches[1], matches[2], request.RequestContext.TimeEpoch, request.Body)
	case "DELETE":
		return deleteToken(ctx, matches[1], matches[2])
	}

	return notFound, nil
}

func createOrUpdateToken(ctx context.Context, table string, token string, timeEpoch int64, body string) (events.APIGatewayV2HTTPResponse, error) {
	params := createParams{}
	json.Unmarshal([]byte(body), &params)

	item := dynamodb.UpdateItemInput{
		TableName: aws.String(tableMap[table]),
		Key: map[string]*dynamodb.AttributeValue{
			"Token": {
				S: aws.String(token),
			},
		},
		UpdateExpression: aws.String("SET #createdAt = if_not_exists(#createdAt, :createdAt), #updatedAt = :updatedAt"),
		ExpressionAttributeNames: map[string]*string{
			"#createdAt": aws.String("CreatedAt"),
			"#updatedAt": aws.String("UpdatedAt"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":createdAt": {
				N: aws.String(strconv.FormatInt(timeEpoch, 10)),
			},
			":updatedAt": {
				N: aws.String(strconv.FormatInt(timeEpoch, 10)),
			},
		},
	}

	_, err := db.Instance.UpdateItem(&item)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Printf("PutItem(%s, %s) aws error occurred: %s", table, token, aerr.Error())
			switch aerr.Code() {
			case dynamodb.ErrCodeRequestLimitExceeded:
				fallthrough
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 503,
				}, nil
			default:
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 500,
				}, nil
			}
		} else {
			log.Printf("PutItem(%s, %s) error occurred: %s", table, token, err.Error())
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 500,
			}, nil
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
	}, nil
}

func deleteToken(ctx context.Context, table string, token string) (events.APIGatewayV2HTTPResponse, error) {
	item := dynamodb.DeleteItemInput{
		TableName: aws.String(tableMap[table]),
		Key: map[string]*dynamodb.AttributeValue{
			"Token": {
				S: aws.String(token),
			},
		},
	}

	_, err := db.Instance.DeleteItem(&item)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Printf("DeleteItem(%s, %s) aws error occurred: %s", table, token, aerr.Error())
			switch aerr.Code() {
			case dynamodb.ErrCodeRequestLimitExceeded:
				fallthrough
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 503,
				}, nil
			default:
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 500,
				}, nil
			}
		} else {
			log.Printf("DeleteItem(%s, %s) error occurred: %s", table, token, err.Error())
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 500,
			}, nil
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
	}, nil
}

package main

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/p2pquake/fcm-controller-lambda/db"
	"github.com/p2pquake/fcm-controller-lambda/notifications"
	"github.com/p2pquake/fcm-controller-lambda/tokens"
)

func init() {
	db.Init()
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if strings.HasPrefix(request.RequestContext.HTTP.Path, tokens.Prefix) {
		return tokens.HandleRequest(ctx, request)
	}

	if strings.HasPrefix(request.RequestContext.HTTP.Path, notifications.Prefix) {
		return notifications.HandleRequest(ctx, request)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 404,
	}, nil
}

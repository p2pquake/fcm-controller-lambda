package notifications

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

const Prefix = "/v1/notifications"

func HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       "Hello, notifications!",
	}, nil
}

package messages

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/fatih/structs"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func LambdaHandler(request Request) (Response, error) {
	structs.DefaultTagName = "json" // https://github.com/fatih/structs/issues/25
	body := []byte(request.Body)
	log.Printf("Receiving request body %s", body)
	resp := Response{
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: 200,
	}
	challengeResponse, err := SlackHandler(body)
	if err != nil {
		resp.StatusCode = 500
	} else {
		resp.StatusCode = 200
		resp.Body = challengeResponse
		resp.Headers["Content-Type"] = "text"
	}

	return resp, nil
}

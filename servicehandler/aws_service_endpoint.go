package servicehandler

import (
	"context"
	"go-micro/logger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type AWSServiceEndpoint struct {
	handler interface{}
}

// ServiceFunction is the function type of microservice funtion implementation
type ServiceFunction func(se ServiceEvent, logger logger.Logger) string

// awsLambdaStart is the trigger for lambda execution that can me mocked in testing
var awsLambdaStart = func(handler interface{}) {
	lambda.Start(handler)
}

// NewServiceEndpoint will create the aws service enpoint instance
func NewServiceEndpoint(es EventSpec, sf ServiceFunction, lgr logger.Logger,
	retHeaders map[string]string) *AWSServiceEndpoint {
	defaultRetHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	genericServiceEnpoint := func(ctx context.Context, event events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse) {
		// Append Return Headers
		for k, v := range retHeaders {
			defaultRetHeaders[k] = v
		}

		// Initialize Service Handler
		lgr.LogTxt(logger.INFO, "Initializing AWS Service Handler..")
		svh := AWSServiceHandler{
			Event:  event,
			Logger: lgr,
		}

		// Handle Http Exceptions
		defer func() {
			err := recover()
			if err != nil {
				response = svh.HandleExceptions(
					err,
					defaultRetHeaders,
				).(events.APIGatewayProxyResponse)
			}
			lgr.DisplayLogsBackward()
		}()

		se := svh.NewServiceEvent(es)

		// Execute the service function
		lgr.LogTxt(logger.INFO, "Executing Service Function..")
		responseBody := sf(se, lgr)

		// Generate New HTTP Response
		lgr.LogTxt(logger.INFO, "Building Response..")
		response = svh.NewHTTPResponse(ServiceResponse{
			StatusCode:    200,
			ReturnBody:    responseBody,
			ReturnHeaders: defaultRetHeaders,
		}).(events.APIGatewayProxyResponse)

		return response
	}

	return &AWSServiceEndpoint{
		handler: genericServiceEnpoint,
	}

}

// Execute will trigger the execution of aws lambda
func (ae AWSServiceEndpoint) Execute() {
	awsLambdaStart(ae.handler)
}

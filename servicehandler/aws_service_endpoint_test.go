package servicehandler

import (
	"context"
	"go-micro/logger"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// Testing constants
const (
	TEST_AWS_REQUEST_RESOURCE      = "mockResource"
	TEST_AWS_REQUEST_PATH          = "mockPath"
	TEST_AWS_RESPONSE_OK           = "{\"message\": \"OK\"}"
	TEST_AWS_BAD_REQUEST           = "Error in Query Parameter, MISSING ATTRIBUTE ERROR. missing attribute 'testQparam'"
	TEST_AWS_INTERNAL_SERVER_ERROR = "\"message\": \"Internal Server Error\""
	TEST_EXTRA_HEADER_KEY          = "extra-header"
	TEST_EXTRA_HEADER_VALUE        = "test-value"
	TEST_SUCCESS_CONTENT_TYPE      = "application/json"
	TEST_ERROR_CONTENT_TYPE        = "text/plain"
)

// serviceEndpointTests for table testing of service endpoint
var serviceEndpointTests = []struct {
	testName     string
	eventSpec    EventSpec
	function     ServiceFunction
	requestEvent events.APIGatewayProxyRequest
	headers      map[string]string
	want         map[string]string
}{
	{
		"test service endpoint valid request",
		EventSpec{},
		func(se ServiceEvent, logger logger.Logger) string { return "{\"message\": \"OK\"}" },
		events.APIGatewayProxyRequest{},
		map[string]string{
			TEST_EXTRA_HEADER_KEY: TEST_EXTRA_HEADER_VALUE,
		},
		map[string]string{
			"testAWSResponse":       TEST_AWS_RESPONSE_OK,
			"testHeaderContentType": TEST_SUCCESS_CONTENT_TYPE,
			TEST_EXTRA_HEADER_KEY:   TEST_EXTRA_HEADER_VALUE,
		},
	},
	{
		"test service endpoint bad request",
		EventSpec{
			RequiredQueryParams: ReqEventSpec{
				ReqEventAttributes: map[string]interface{}{
					"testQparam": NewReqEvenAttrib("string", true, 4, 50),
				},
			},
		},
		func(se ServiceEvent, logger logger.Logger) string { return "" },
		events.APIGatewayProxyRequest{},
		map[string]string{
			TEST_EXTRA_HEADER_KEY: TEST_EXTRA_HEADER_VALUE,
		},
		map[string]string{
			"testAWSResponse":       TEST_AWS_BAD_REQUEST,
			"testHeaderContentType": TEST_ERROR_CONTENT_TYPE,
			TEST_EXTRA_HEADER_KEY:   TEST_EXTRA_HEADER_VALUE,
		},
	},
}

func TestServiceValidRequest(t *testing.T) {
	for _, tt := range serviceEndpointTests {
		t.Run(tt.testName, func(t *testing.T) {
			// mock aws lambda start for testing expected values
			awsLambdaStart = func(handler interface{}) {
				var ctx context.Context
				response := handler.(func(context.Context, events.APIGatewayProxyRequest) events.APIGatewayProxyResponse)(
					ctx,
					tt.requestEvent,
				)
				if response.Body != tt.want["testAWSResponse"] {
					t.Error("service enpoint response not matched")
				}
				if response.Headers["Content-Type"] != tt.want["testHeaderContentType"] {
					t.Error("invalid value for response header Content-Type")
				}
				if response.Headers[TEST_EXTRA_HEADER_KEY] != tt.want[TEST_EXTRA_HEADER_KEY] {
					t.Errorf("invalid value for response header extra-heder")
				}
			}

			lgr := logger.NewLogger()
			testServiceEndpoint := NewServiceEndpoint(tt.eventSpec, tt.function, lgr, tt.headers)
			testServiceEndpoint.Execute()
		})
	}
}

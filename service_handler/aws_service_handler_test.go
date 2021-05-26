package service_handler

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var newServiceTests = []struct {
	testName string
	awsEvent events.APIGatewayProxyRequest
	isValid  bool
}{
	{
		"valid request",
		newAWSMockEvent(
			map[string]string{
				"fields": "firstname",
			},
			map[string]string{
				"department": "IT",
			},
			`
				{
					"firstname" : "juan",
					"middlename" : "ponce",
					"lastname" : "dela cruz"
				}
			`,
		),
		true,
	},
	{
		"invalid request body",
		newAWSMockEvent(
			map[string]string{
				"fields": "firstname",
			},
			map[string]string{
				"department": "IT",
			},
			`
				{
					"firstname" : ""
				}
			`,
		),
		false,
	},
	{
		"invalid query params",
		newAWSMockEvent(
			map[string]string{},
			map[string]string{
				"department": "IT",
			},
			`
				{
					"firstname": "juan",
					"middlename": "ponce",
					"lastname": "dela cruz"
				}
			`,
		),
		false,
	},
	{
		"invalid path params",
		newAWSMockEvent(
			map[string]string{
				"fields": "firstname",
			},
			map[string]string{},
			`
				{
					"firstname": "juan",
					"middlename": "ponce",
					"lastname": "dela cruz"
				}
			`,
		),
		false,
	},
}

func newAWSMockEvent(qParam map[string]string, pParam map[string]string,
	body string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		Resource:              "mockResource",
		Path:                  "mockPath",
		QueryStringParameters: qParam,
		PathParameters:        pParam,
		RequestContext: events.APIGatewayProxyRequestContext{
			ResourcePath: "mockResourcePath",
			Identity:     events.APIGatewayRequestIdentity{},
		},
		Body: body,
	}
}

func TestNewService(t *testing.T) {
	serviceSpec := ServiceSpec{
		RequiredRequestBody: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"firstname":  NewReqEvenAttrib("string", true, 2, 50),
				"lastname":   NewReqEvenAttrib("string", true, 2, 50),
				"middlename": NewReqEvenAttrib("string", true, 2, 50),
			}},
		RequiredQueryParams: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"fields": NewReqEvenAttrib("string", true, 4, 50),
			},
		},
		RequiredPathParams: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"department": NewReqEvenAttrib("string", true, 2, 50),
			},
		},
	}

	for _, tt := range newServiceTests {
		t.Run(tt.testName, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil && !tt.isValid {
					t.Error("Invalid request not caught")
				}
			}()
			var service_handler ServiceHandler = AWSServiceHandler{
				event: tt.awsEvent,
			}
			serviceEvent := service_handler.NewService(serviceSpec)
			if reflect.TypeOf(serviceEvent).String() != "service_handler.ServiceEvent" {
				t.Error("Invalid ReqEventAttrib")
			}
		})
	}

}

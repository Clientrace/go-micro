package servicehandler

import (
	"fmt"
	"go-micro/logger"
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
			var serviceHandler ServiceHandler = AWSServiceHandler{
				Event: tt.awsEvent,
			}
			serviceEvent := serviceHandler.NewService(serviceSpec)
			if reflect.TypeOf(serviceEvent).String() != "servicehandler.ServiceEvent" {
				t.Error("Invalid ReqEventAttrib")
			}
		})
	}

}

func TestAWSNewResponse(t *testing.T) {
	logger := logger.NewLogger()
	var serviceHandler = AWSServiceHandler{
		Event:  events.APIGatewayProxyRequest{},
		Logger: logger,
	}

	want := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            `{"message" : "OK"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	got := serviceHandler.NewHTTPResponse(ServiceResponse{
		StatusCode: 200,
		ReturnBody: `{"message" : "OK"}`,
		ReturnHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	})

	if want.StatusCode != got.(events.APIGatewayProxyResponse).StatusCode {
		t.Errorf("Invalid AWS Service Response Status Code")
	}

	if want.Body != got.(events.APIGatewayProxyResponse).Body {
		t.Errorf("Invalid AWS Service Response Body")
	}

	logger.DisplayLogsBackward()

}

func testBadRequest() (response events.APIGatewayProxyResponse) {
	returnHeaders := map[string]string{
		"Content-Type": "application/json",
	}
	logger := logger.NewLogger()
	sh := AWSServiceHandler{
		Event: newAWSMockEvent(
			map[string]string{},
			map[string]string{},
			`{"testParam": "testParamValue"}`,
		),
		Logger: logger,
	}
	requestSpec := ServiceSpec{
		RequiredRequestBody: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"username": map[string]interface{}{
					"firstName":  NewReqEvenAttrib("string", true, 4, 15),
					"lastName":   NewReqEvenAttrib("string", true, 4, 255),
					"middleName": NewReqEvenAttrib("string", true, 4, 255),
				},
			},
		},
		RequiredQueryParams: ReqEventSpec{},
		RequiredPathParams:  ReqEventSpec{},
	}
	defer func() {
		response = sh.HandleExceptions(
			recover(),
			returnHeaders,
		).(events.APIGatewayProxyResponse)
		logger.DisplayLogsBackward()
	}()

	service := sh.NewService(requestSpec)
	fmt.Println(service.PathParams)
	fmt.Println(service.QueryParams)
	fmt.Println(service.RequestBody)

	return response
}

func testInternalServerError() (response events.APIGatewayProxyResponse) {
	returnHeaders := map[string]string{
		"Content-Type": "application/json",
	}
	logger := logger.NewLogger()
	sh := AWSServiceHandler{
		Event: newAWSMockEvent(
			map[string]string{},
			map[string]string{},
			`{
				"username": {
					"firstName": "clarence",
					"lastName" : "penaflor",
					"middleName": "par"
				}
			}`,
		),
		Logger: logger,
	}
	requestSpec := ServiceSpec{
		RequiredRequestBody: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"username": map[string]interface{}{
					"firstName":  NewReqEvenAttrib("string", true, 4, 15),
					"lastName":   NewReqEvenAttrib("string", true, 4, 255),
					"middleName": NewReqEvenAttrib("string", true, 3, 255),
				},
			},
		},
		RequiredQueryParams: ReqEventSpec{},
		RequiredPathParams:  ReqEventSpec{},
	}
	defer func() {
		response = sh.HandleExceptions(
			recover(),
			returnHeaders,
		).(events.APIGatewayProxyResponse)
		logger.DisplayLogsBackward()
	}()

	service := sh.NewService(requestSpec)
	fmt.Println(service.PathParams)

	varA := 0
	varB := 1

	// Intentional Error for testin
	fmt.Println(varB / varA)
	return response

}

func TestBadRequestException(t *testing.T) {
	gotResponse := testBadRequest()
	wantResponse := events.APIGatewayProxyResponse{
		StatusCode:      400,
		IsBase64Encoded: false,
		Body:            "Error in Request Body, MISSING ATTRIBUTE ERROR. missin attribute 'middleName'",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	if gotResponse.StatusCode != wantResponse.StatusCode {
		t.Errorf("Invalid response status code for bad request testing")
	}
}

func TestInternalServerErrorException(t *testing.T) {
	gotResponse := testInternalServerError()
	wantResponse := events.APIGatewayProxyResponse{
		StatusCode:      500,
		IsBase64Encoded: false,
		Body:            "Internal server error",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	fmt.Println(gotResponse)

	if gotResponse.StatusCode != wantResponse.StatusCode {
		t.Errorf("Invalid repsonse status code for internalserver error testing")
	}

}

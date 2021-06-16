package servicehandler

import (
	"fmt"
	"go-micro/logger"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var returnHeaders = map[string]string{
	"Content-Type": "application/json",
}
var handleExceptionTests = []struct {
	testName string
	input    interface{}
	want     events.APIGatewayProxyResponse
}{
	{
		"internal server error test (string)",
		"error test",
		events.APIGatewayProxyResponse{
			StatusCode:      500,
			IsBase64Encoded: false,
			Body:            "Internal Server Error",
		},
	},
	{
		"internal server error test (struct)",
		map[string]string{
			"testError": "testError",
		},
		events.APIGatewayProxyResponse{
			StatusCode:      500,
			IsBase64Encoded: false,
			Body:            "Internal Server Error",
		},
	},
}

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

func TestNewServiceEvent(t *testing.T) {
	eventSpec := EventSpec{
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

	logger := logger.NewLogger()
	for _, tt := range newServiceTests {
		t.Run(tt.testName, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil && !tt.isValid {
					t.Error("Invalid request not caught")
				}
			}()
			var serviceHandler ServiceHandler = AWSServiceHandler{
				Event:  tt.awsEvent,
				Logger: logger,
			}
			serviceEvent := serviceHandler.NewServiceEvent(eventSpec)
			if reflect.TypeOf(serviceEvent).String() != "servicehandler.ServiceEvent" {
				t.Error("Invalid ReqEventAttrib")
			}
		})
	}

}

func TestAWSNewResponse(t *testing.T) {
	logger := logger.NewLogger()
	returnBody := "{\"Body\": \"OK\"}"

	var serviceHandler = AWSServiceHandler{
		Event:  events.APIGatewayProxyRequest{},
		Logger: logger,
	}

	want := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            returnBody,
		Headers:         returnHeaders,
	}

	got := serviceHandler.NewHTTPResponse(ServiceResponse{
		StatusCode:    200,
		ReturnBody:    returnBody,
		ReturnHeaders: returnHeaders,
	})

	if want.StatusCode != got.(events.APIGatewayProxyResponse).StatusCode {
		t.Errorf("Invalid AWS Service Response Status Code")
	}

	if want.Body != got.(events.APIGatewayProxyResponse).Body {
		t.Errorf("Invalid AWS Service Response Body")
	}

	logger.DisplayLogsBackward()

}

func testValidRequest() (response events.APIGatewayProxyResponse) {
	logger := logger.NewLogger()
	sh := AWSServiceHandler{
		Event: newAWSMockEvent(
			map[string]string{},
			map[string]string{},
			`{
				"username" : {
					"firstName" : "juan",
					"lastName" : "delacruz",
					"middleName" : "ponce"
				}
			}`,
		),
		Logger: logger,
	}
	requestSpec := EventSpec{
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
		if ex := sh.HandleExceptions(recover(), returnHeaders); ex != nil {
			response = ex.(events.APIGatewayProxyResponse)
		}
		logger.DisplayLogsBackward()
	}()

	service := sh.NewServiceEvent(requestSpec)
	fmt.Println(service.PathParams)
	fmt.Println(service.QueryParams)
	fmt.Println(service.RequestBody)

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            `{"message" : "OK"}`,
		Headers:         returnHeaders,
	}
}

func testBadRequest() (response events.APIGatewayProxyResponse) {
	logger := logger.NewLogger()
	sh := AWSServiceHandler{
		Event: newAWSMockEvent(
			map[string]string{},
			map[string]string{},
			`{}`,
		),
		Logger: logger,
	}
	requestSpec := EventSpec{
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
		if ex := sh.HandleExceptions(recover(), returnHeaders); ex != nil {
			response = ex.(events.APIGatewayProxyResponse)
		}
		logger.DisplayLogsBackward()
	}()

	service := sh.NewServiceEvent(requestSpec)
	fmt.Println(service.PathParams)
	fmt.Println(service.QueryParams)
	fmt.Println(service.RequestBody)

	return response
}

func testInternalServerError() (response events.APIGatewayProxyResponse) {
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
	requestSpec := EventSpec{
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
		if ex := sh.HandleExceptions(recover(), returnHeaders); ex != nil {
			response = ex.(events.APIGatewayProxyResponse)
		}
		logger.DisplayLogsBackward()
	}()

	service := sh.NewServiceEvent(requestSpec)
	fmt.Println(service.PathParams)

	varA := 0
	varB := 1

	// Intentional Error for testin
	fmt.Println(varB / varA)
	return response

}

func TestValidRequest(t *testing.T) {
	gotResponse := testValidRequest()
	wantResponse := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "",
		Headers:         returnHeaders,
	}
	if gotResponse.StatusCode != wantResponse.StatusCode {
		t.Errorf("Invalid response status code for valid http request")
	}
}

func TestBadRequestException(t *testing.T) {
	gotResponse := testBadRequest()
	wantResponse := events.APIGatewayProxyResponse{
		StatusCode:      400,
		IsBase64Encoded: false,
		Body:            "Error in Request Body, MISSING ATTRIBUTE ERROR. missin attribute 'middleName'",
		Headers:         returnHeaders,
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
		Headers:         returnHeaders,
	}

	fmt.Println(gotResponse)

	if gotResponse.StatusCode != wantResponse.StatusCode {
		t.Errorf("Invalid repsonse status code for internalserver error testing")
	}

}

func TestAWSHTTPExceptions(t *testing.T) {
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

	for _, tt := range handleExceptionTests {
		t.Run(tt.testName, func(t *testing.T) {
			defer func() {
				got := sh.HandleExceptions(recover(), returnHeaders).(events.APIGatewayProxyResponse)
				if tt.want.StatusCode != got.StatusCode {
					t.Errorf("Invalid aws http exception. status code does not match")
				}
			}()
			panic(tt.input)
		})
	}

}

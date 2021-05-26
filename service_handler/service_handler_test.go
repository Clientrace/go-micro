package service_handler

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var AGE_INVALID_FLOAT32_VALUE float32 = 0.11
var AGE_INVALID_FLOAT64_VALUE float64 = 0.11111111111111111
var AGE_INVALID_INT_VALUE int = 10000000000

// All the subtest for recursive attribute checking
var attribCheckTests = []struct {
	testName   string
	attributes map[string]interface{}
	want       int
}{
	{
		"attribute OK",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  "firstname",
				"lastName":   "latname",
				"middleName": "middleName",
			},
			"email":      "test@gmail.com",
			"age":        5,
			"isEmployed": true,
		},
		ATTRIBUTE_OK,
	},
	{
		"missing key check",
		map[string]interface{}{
			"email": "test@gmail.com",
		},
		MISSING_ATTRIBUTE_ERROR,
	},
	{
		"invalid parent object attribute key type",
		map[string]interface{}{
			"username":   "test",
			"email":      "test@email.com",
			"age":        2,
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_TYPE_ERROR,
	},
	{
		"missing child keys from parent attribute test",
		map[string]interface{}{
			"username":   map[string]interface{}{},
			"email":      "test@email.com",
			"age":        2,
			"isEmployed": false,
		},
		MISSING_ATTRIBUTE_ERROR,
	},
	{
		"attribute type checking [invalid string]",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  0,
				"lastName":   0,
				"middleName": 1,
			},
			"email":      1,
			"age":        2,
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_TYPE_ERROR,
	},
	{
		"attribute type checking [invalid number]",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  "testFirstName",
				"lastName":   "testLastName",
				"middleName": "testMiddleName",
			},
			"email":      "test email",
			"age":        "testInvalidValue",
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_TYPE_ERROR,
	},
	{
		"attribute type checking [invalid boolean]",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  "testFirstName",
				"lastName":   "testLastName",
				"middleName": "testMiddleName",
			},
			"email":      "test email",
			"age":        5,
			"isEmployed": "testValue",
		},
		INVALID_ATTRIBUTE_TYPE_ERROR,
	},
	{
		"missing child key check",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstname": "test",
			},
			"email":      "test@email.com",
			"age":        2,
			"isEmployed": false,
		},
		MISSING_ATTRIBUTE_ERROR,
	},
	{
		"invalid child attrib type check",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  0,
				"lastName":   "testLastname",
				"middleName": "testMiddlename",
			},
			"email":      "test@email.com",
			"age":        2,
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_TYPE_ERROR,
	},
	{
		"invalid length, string length too short",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  "1",
				"lastName":   "testLastname",
				"middleName": "testMiddlename",
			},
			"email":      "test@email.com",
			"age":        2,
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_LENGTH_ERROR,
	},
	{
		"invalid length, string length too long",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  "abcdefghijklmnopq",
				"lastName":   "testLastname",
				"middleName": "testMiddlename",
			},
			"email":      "test@email.com",
			"age":        2,
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_LENGTH_ERROR,
	},
	{
		"invalid length, string length too short",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  "",
				"lastName":   "testLastname",
				"middleName": "testMiddlename",
			},
			"email":      "test@email.com",
			"age":        2,
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_LENGTH_ERROR,
	},
	{
		"invalid length, number length too short",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  "testFirstname",
				"lastName":   "testLastname",
				"middleName": "testMiddlename",
			},
			"email":      "test@email.com",
			"age":        AGE_INVALID_INT_VALUE,
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_LENGTH_ERROR,
	},
	{
		"invalid length, number length too short",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  "testFirstname",
				"lastName":   "testLastname",
				"middleName": "testMiddlename",
			},
			"email":      "test@email.com",
			"age":        AGE_INVALID_FLOAT64_VALUE,
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_LENGTH_ERROR,
	},
	{
		"invalid length, number length too short",
		map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  "testFirstname",
				"lastName":   "testLastname",
				"middleName": "testMiddlename",
			},
			"email":      "test@email.com",
			"age":        AGE_INVALID_FLOAT32_VALUE,
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_LENGTH_ERROR,
	},
}

var eventAttributeTests = []struct {
	testName       string
	reqEventAttrib map[string]interface{}
	isValid        bool
}{
	{
		"valid event attributes [string]",
		map[string]interface{}{
			"DataType":   "string",
			"IsRequired": true,
			"MinLength":  10,
			"MaxLength":  10,
		},
		true,
	},
	{
		"valid event attributes [number]",
		map[string]interface{}{
			"DataType":   "number",
			"IsRequired": false,
			"MinLength":  10,
			"MaxLength":  10,
		},
		true,
	},
	{
		"valid event attributes [boolean]",
		map[string]interface{}{
			"DataType":   "boolean",
			"IsRequired": false,
			"MinLength":  10,
			"MaxLength":  10,
		},
		true,
	}, {
		"invalid event attribute",
		map[string]interface{}{
			"DataType":   "test",
			"IsRequired": false,
			"MinLength":  10,
			"MaxLength":  10,
		},
		false,
	},
}

var serviceHandler ServiceHandler = AWSServiceHandler{
	event: events.APIGatewayProxyRequest{
		Resource: "mockResource",
		Path:     "mockPath",
		QueryStringParameters: map[string]string{
			"fields": "firstname",
		},
		PathParameters: map[string]string{
			"department": "IT",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			ResourcePath: "mockResourcepath",
			Identity:     events.APIGatewayRequestIdentity{},
		},
		Body: `
				{
					"firstname" : "juan",
					"middlename" : "ponce",
					"lastname" : "dela cruz"
				}
			`,
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

func TestRecursiveAtribCheck(t *testing.T) {
	var serviceHandler ServiceHandler = AWSServiceHandler{}
	requestSpec := ReqEventSpec{
		ReqEventAttributes: map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  serviceHandler.NewReqEvenAttrib("string", true, 4, 15),
				"lastName":   serviceHandler.NewReqEvenAttrib("string", true, 4, 255),
				"middleName": serviceHandler.NewReqEvenAttrib("string", true, 4, 255),
			},
			"email":      serviceHandler.NewReqEvenAttrib("string", true, 4, 250),
			"age":        serviceHandler.NewReqEvenAttrib("number", true, 1, 1000),
			"isEmployed": serviceHandler.NewReqEvenAttrib("boolean", true, 0, 0),
		},
	}
	for _, tt := range attribCheckTests {
		t.Run(tt.testName, func(t *testing.T) {
			got, _ := recursiveAttributeCheck("testEndpoint", requestSpec, tt.attributes, 0)
			if got != tt.want {
				t.Errorf("recursive attribute check got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEventAttrib(t *testing.T) {
	var serviceHandler ServiceHandler = AWSServiceHandler{}

	for _, tt := range eventAttributeTests {
		t.Run(tt.testName, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil && !tt.isValid {
					t.Error("Invalid attribute not caught")
				}
			}()
			reqAttribInstance := serviceHandler.NewReqEvenAttrib(
				tt.reqEventAttrib["DataType"].(string),
				tt.reqEventAttrib["IsRequired"].(bool),
				tt.reqEventAttrib["MinLength"].(int),
				tt.reqEventAttrib["MaxLength"].(int),
			)
			if reflect.TypeOf(reqAttribInstance).String() != "service_handler.ReqEventAttrib" {
				t.Error("Invalid ReqEventAttrib")
			}
		})
	}
}

func TestNewService(t *testing.T) {
	serviceSpec := ServiceSpec{
		RequiredRequestBody: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"firstname":  serviceHandler.NewReqEvenAttrib("string", true, 2, 50),
				"lastname":   serviceHandler.NewReqEvenAttrib("string", true, 2, 50),
				"middlename": serviceHandler.NewReqEvenAttrib("string", true, 2, 50),
			}},
		RequiredQueryParams: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"fields": serviceHandler.NewReqEvenAttrib("string", true, 4, 50),
			},
		},
		RequiredPathParams: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"department": serviceHandler.NewReqEvenAttrib("string", true, 2, 50),
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

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
	var serviceHandler ServiceHandler = AWSServiceHandler{
		event: events.APIGatewayProxyRequest{
			Resource: "mockResource",
			Path:     "mockPath",
			QueryStringParameters: map[string]string{
				"first":             "queryStringValue1",
				"queryStringParam2": "queryStringValue2",
				"queryStringParam3": "queryStringValue3",
			},
			PathParameters: map[string]string{
				"pathParam1": "pathParamValue1",
				"pathParam2": "pathParamValue2",
				"pathParam3": "pathParamValue3",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: "mockResourcepath",
				Identity:     events.APIGatewayRequestIdentity{},
			},
			Body: `
				{
					"bodyParam1" : "bodyParamValue1",
					"bodyParam2" : "bodyParamValue2",
					"bodyParam3" : "bodyParamValue3"
				}
			`,
		},
	}
	serviceSpec := ServiceSpec{
		RequiredRequestBody: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"firstname":  serviceHandler.NewReqEvenAttrib("string", true, 10, 10),
				"lastname":   serviceHandler.NewReqEvenAttrib("string", true, 10, 10),
				"middlename": serviceHandler.NewReqEvenAttrib("string", true, 10, 10),
			}},
		RequiredQueryParams: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"filterBy": serviceHandler.NewReqEvenAttrib("string", true, 10, 10),
			},
		},
		RequiredPathParams: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"category": serviceHandler.NewReqEvenAttrib("string", true, 10, 10),
			},
		},
	}
	var testService = serviceHandler.NewService(serviceSpec)
	wantIdentity := Identity{
		Email:    "test@gmail.com",
		Username: "testusername",
		Role:     "testuserrole",
	}

	if testService.Identity != wantIdentity {
		t.Error("Invalid identity")
	}
}

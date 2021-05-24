package service_handler

import (
	"testing"
)

// All the subtest for recursive attribute checking
var attribCheckTests = []struct {
	testName   string
	attributes map[string]interface{}
	want       string
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
			"username":   "test_username",
			"email":      "test email",
			"age":        "testInvalidValue",
			"isEmployed": false,
		},
		INVALID_ATTRIBUTE_TYPE_ERROR,
	},
	{
		"attribute type checking [invalid boolean]",
		map[string]interface{}{
			"username":   "test_username",
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
}

func TestRecursiveAtribCheck(t *testing.T) {
	// Test Specifications for request attribute
	requestSpec := ReqEventSpec{
		ReqEventAttributes: map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  NewReqEvenAttrib("string", true, 4, 15),
				"lastName":   NewReqEvenAttrib("string", true, 4, 255),
				"middleName": NewReqEvenAttrib("string", true, 4, 255),
			},
			"email":      NewReqEvenAttrib("string", true, 4, 250),
			"age":        NewReqEvenAttrib("number", true, 1, 1000),
			"isEmployed": NewReqEvenAttrib("boolean", true, 0, 0),
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

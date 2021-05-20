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
		"missing key check",
		map[string]interface{}{
			"email": "test@gmail.com",
		},
		MISSING_ATTRIBUTE_ERROR,
	},
	{
		"attribute type checking [invalid string]",
		map[string]interface{}{
			"username": 0,
			"email":    1,
			"age":      2,
		},
		INVALID_ATTRIBUTE_TYPE_ERROR,
	},
	{
		"attribute type checking [invalid number]",
		map[string]interface{}{
			"username": "test_username",
			"email":    "test email",
			"age":      "testInvalidValue",
		},
		INVALID_ATTRIBUTE_TYPE_ERROR,
	},
	{
		"missing child key check",
		map[string]interface{}{
			"username": map[string]interface{}{},
			"email":    "test@email.com",
			"age":      2,
		},
		MISSING_ATTRIBUTE_ERROR,
	},
	{
		"invalid child attrib type check",
		map[string]interface{}{
			"username": map[string]interface{}{"firstName": 0},
			"email":    "test@email.com",
			"age":      2,
		},
		INVALID_ATTRIBUTE_TYPE_ERROR,
	},
}

func TestRecursiveAtribCheck(t *testing.T) {
	// Test Specifications for request attribute
	requestSpec := ReqEventSpec{
		ReqEventAttributes: map[string]interface{}{
			"username": map[string]interface{}{
				"firstName":  NewReqEvenAttrib("string", true, 4, 4),
				"lastName":   NewReqEvenAttrib("string", true, 4, 255),
				"middleName": NewReqEvenAttrib("string", true, 4, 255),
			},
			"email":      NewReqEvenAttrib("string", true, 4, 4),
			"age":        NewReqEvenAttrib("string", true, 1, 1000),
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

package service_handler

import (
	"fmt"
	"testing"
)

func TestRecursiveAtribCheck(t *testing.T) {
	requestSpec := ReqEventSpec{
		ReqEventAttributes: map[string]interface{}{
			"username": map[string]interface{}{
				"firstName": ReqEventAttrib{
					DataType:   "string",
					IsRequired: true,
					MinLength:  4,
					MaxLength:  4,
				},
				"lastName": ReqEventAttrib{
					DataType:   "string",
					IsRequired: true,
					MinLength:  4,
					MaxLength:  4,
				},
				"middleName": ReqEventAttrib{
					DataType:   "string",
					IsRequired: false,
					MinLength:  4,
					MaxLength:  255,
				},
			},
			"email": ReqEventAttrib{
				DataType:   "string",
				IsRequired: true,
				MinLength:  4,
				MaxLength:  255,
			},
			"age": ReqEventAttrib{
				DataType:   "number",
				IsRequired: true,
				MinLength:  0,
				MaxLength:  0,
			},
		},
	}

	fmt.Println("Testing recursive attribute check if there's a missing key..")
	got, _ := recursiveAttributeCheck(
		"testEndpoint",
		requestSpec, map[string]interface{}{
			"email": "test@gmail.com",
			"age":   0,
		},
	)

	want := MISSING_ATTRIBUTE_ERROR
	if got != want {
		t.Errorf("Failed Missing attrib check. Got %v, expecting %v", got, want)
	} else {
		fmt.Println("=> OK")
	}

	fmt.Println("Testing recursive attribute type checking..")
	got, _ = recursiveAttributeCheck(
		"testEndpoint",
		requestSpec, map[string]interface{}{
			"username": 0,
			"email":    1,
			"age":      2,
		},
	)

	want = INVALID_ATTRIBUTE_TYPE_ERROR
	if got != want {
		t.Errorf("Failed Attribute Type Check. Got %v, expecting %v", got, want)
	} else {
		fmt.Println("=> OK")
	}

}

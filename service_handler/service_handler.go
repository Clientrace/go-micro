package service_handler

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
)

/* Attribute Parse Codes */
const (
	ATTRIBUTE_OK                   = iota
	MISSING_ATTRIBUTE_ERROR        = iota
	INVALID_ATTRIBUTE_TYPE_ERROR   = iota
	INVALID_ATTRIBUTE_LENGTH_ERROR = iota
)

/* Param Type for Parse Code */
const (
	REQ_BODY     = iota
	QUERY_PARAMS = iota
	PATH_PARAMS  = iota
)

/* Required Event Specification Attribute */
type ReqEventAttrib struct {
	DataType   string
	IsRequired bool
	MinLength  int
	MaxLength  int
}

/* Service Event specification */
type ServiceSpec struct {
	RequiredRequestBody ReqEventSpec
	RequiredQueryParams ReqEventSpec
	RequiredPathParams  ReqEventSpec
}

/* Required Event specification */
type ReqEventSpec struct {
	ReqEventAttributes map[string]interface{}
}

/* Request on Service Event */
type ServiceEvent struct {
	Identity    events.APIGatewayRequestIdentity
	RequestBody map[string]interface{}
	QueryParams map[string]interface{}
	PathParams  map[string]interface{}
}

type ServiceResponse struct {
	StatusCode    int
	ReturnBody    string
	ReturnHeaders map[string]string
}

type ServiceHandler interface {
	NewService(ServiceSpec) ServiceEvent
	NewHTTPResponse(ServiceResponse) interface{}
}

var parameterMap = map[int]string{
	QUERY_PARAMS: "Query Parameter",
	PATH_PARAMS:  "Path Parameter",
	REQ_BODY:     "Request Body",
}

var errMsgMap = map[int]string{
	INVALID_ATTRIBUTE_LENGTH_ERROR: "INVALID ATTRIBUTE LENGTH",
	INVALID_ATTRIBUTE_TYPE_ERROR:   "INVALID ATTRIBUTE TYPE",
}

/*
	Cause a panic in service handler for the
	http_exceptions to catch orrecover from.
*/
func causePanic(paramType int, parseCode int, errorMsg string) {
	panic(HTTPException{
		StatusCode:   BAD_REQUEST,
		ErrorMessage: fmt.Sprintf("Error in %v, %v. %v", parameterMap[paramType], errMsgMap[parseCode], errorMsg),
	})
}

/*
	Create new Required Event Attributes
*/
func NewReqEvenAttrib(dataType string, isRequired bool, minLength int, maxLength int) ReqEventAttrib {
	validDataTypes := []string{"string", "number", "boolean"}
	invalidDataType := true
	for _, v := range validDataTypes {
		if v == dataType {
			invalidDataType = false
			break
		}
	}
	if invalidDataType {
		panic("invalid attribute type, attribute type can only be of the ff [string ,number, boolean]")
	}
	return ReqEventAttrib{
		DataType:   dataType,
		IsRequired: isRequired,
		MinLength:  minLength,
		MaxLength:  maxLength,
	}

}

// Check request attributes deep. See if passed attribute match the specs
func recursiveAttributeCheck(endpoint string, reqEventSpec ReqEventSpec, attributes map[string]interface{}, depth int) (int, string) {
	rqa := reqEventSpec.ReqEventAttributes

	// Iterate through required attributes
	for k := range rqa {
		if _, ok := attributes[k]; ok {
			if reflect.TypeOf(rqa[k]).Kind().String() == "map" {
				if reflect.TypeOf(attributes[k]).String() == "map[string]interface {}" {
					// Create ReqEventSpec for child attribute
					reqEventChild := ReqEventSpec{
						ReqEventAttributes: rqa[k].(map[string]interface{}),
					}
					ret, errMsg := recursiveAttributeCheck(endpoint, reqEventChild, attributes[k].(map[string]interface{}), depth+1)
					if ret != ATTRIBUTE_OK {
						return ret, errMsg
					}
				} else {
					return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Sprintf(
						"invalid type of attribute '%v'. expected object got %v",
						k,
						reflect.TypeOf(attributes[k]),
					)
				}

			} else {

				// Required VS Found Attribute Type
				reqAttributeType := rqa[k].(ReqEventAttrib).DataType
				reqAttributeMaxLength := rqa[k].(ReqEventAttrib).MaxLength
				reqAttributeMinLength := rqa[k].(ReqEventAttrib).MinLength
				foundAttribType := reflect.TypeOf(attributes[k]).String()

				if reqAttributeType == "string" && foundAttribType == "string" {
					if !(len(attributes[k].(string)) >= reqAttributeMinLength && len(attributes[k].(string)) <=
						reqAttributeMaxLength) {
						return INVALID_ATTRIBUTE_LENGTH_ERROR, fmt.Sprintf(
							"invalid length of attribute %v. min length: %d, max length: %d",
							k,
							reqAttributeMinLength,
							reqAttributeMaxLength,
						)
					}
				} else if reqAttributeType == "string" && foundAttribType != "string" {
					return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Sprintf(
						"invalid type of attribute %v. expected string, got %v",
						k,
						reflect.TypeOf(attributes[k]).String(),
					)
				}
				if reqAttributeType == "number" && (foundAttribType == "int" || foundAttribType == "float32" ||
					foundAttribType == "float64") {
					if foundAttribType == "int" {
						if !(attributes[k].(int) >= reqAttributeMinLength && attributes[k].(int) <=
							reqAttributeMaxLength) {
							return INVALID_ATTRIBUTE_LENGTH_ERROR, fmt.Sprintf(
								"invalid length of attribute %v. min length: %d, max length: %d",
								k,
								reqAttributeMinLength,
								reqAttributeMaxLength,
							)
						}
					}
					if foundAttribType == "float64" {
						if !(attributes[k].(float64) >= float64(reqAttributeMinLength) && attributes[k].(float64) <=
							float64(reqAttributeMaxLength)) {
							return INVALID_ATTRIBUTE_LENGTH_ERROR, fmt.Sprintf(
								"invalid length of attribute %v. min length: %d, max length: %d",
								k,
								reqAttributeMinLength,
								reqAttributeMaxLength,
							)
						}
					}
					if foundAttribType == "float32" {
						if !(attributes[k].(float32) >= float32(reqAttributeMinLength) && attributes[k].(float32) <=
							float32(reqAttributeMaxLength)) {
							return INVALID_ATTRIBUTE_LENGTH_ERROR, fmt.Sprintf(
								"invalid length of attribute %v. min length: %d, max length: %d",
								k,
								reqAttributeMinLength,
								reqAttributeMaxLength,
							)
						}
					}
				} else if reqAttributeType == "number" && !(foundAttribType == "int" || foundAttribType == "float32" ||
					foundAttribType == "float64") {
					return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Sprintf(
						"invalid type of attribute %v. expected number(int or float) got %v",
						k,
						attributes[k],
					)
				}

				if reqAttributeType == "boolean" && foundAttribType != "bool" {
					return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Sprintf(
						"invalid type of attribute %v. expected boolean",
						k,
					)
				}
			}
		} else {
			if reflect.TypeOf(rqa[k]).String() == "map[string]interface {}" {
				return MISSING_ATTRIBUTE_ERROR, fmt.Sprintf("missing attribute '%v'", k)
			} else {
				if rqa[k].(ReqEventAttrib).IsRequired {
					return MISSING_ATTRIBUTE_ERROR, fmt.Sprintf("missing attribute '%v'", k)
				}
			}
		}
	}

	return ATTRIBUTE_OK, "OK"
}

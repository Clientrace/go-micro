package service_handler

import (
	"encoding/json"
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

/* Requester Identity */
type Identity struct {
	Email    string
	Username string
	Role     string
}

/* Request on Service Event */
type ServiceEvent struct {
	Identity    Identity
	RequestBody map[string]interface{}
	QueryParams map[string]interface{}
	PathParams  map[string]interface{}
}

type ServiceHandler interface {
	NewService(ServiceSpec) ServiceEvent
	NewReqEvenAttrib(string, bool, int, int) ReqEventAttrib
}

type AWSServiceHandler struct {
	event events.APIGatewayProxyRequest
}

/*
	Cause a panic in service handler for the
	http_exceptions to catch orrecover from.
*/
func causePanic(paramType int, parseCode int, errorMsg string) {
	panic(map[string]interface{}{
		"paramType": paramType,
		"errorMsg":  errorMsg,
		"parseCode": parseCode,
	})
}

/*
	@Exported Function
	Parse AWS Event to get the identity and requests objects
	event AWS  HTTP Event
	requestFmt Required event body format
*/
func (ah AWSServiceHandler) NewService(ss ServiceSpec) ServiceEvent {
	requestEndpoint := ah.event.RequestContext.ResourcePath

	var requestBody map[string]interface{}
	queryParamsMapBuffer := ah.event.QueryStringParameters
	pathParamsMapBuffer := ah.event.PathParameters

	// Convert JSON String body to map
	json.Unmarshal([]byte(ah.event.Body), &requestBody)

	parseCode, errMsg := recursiveAttributeCheck(requestEndpoint, ss.RequiredRequestBody, requestBody, 0)
	if parseCode != ATTRIBUTE_OK {
		causePanic(REQ_BODY, parseCode, errMsg)
	}

	// Covert queryParamsBuffer of map[string]string type to map[string]interface{}
	queryParams := make(map[string]interface{}, len(queryParamsMapBuffer))
	for k, v := range queryParamsMapBuffer {
		queryParams[k] = v
	}
	parseCode, errMsg = recursiveAttributeCheck(requestEndpoint, ss.RequiredQueryParams, queryParams, 0)
	if parseCode != ATTRIBUTE_OK {
		causePanic(QUERY_PARAMS, parseCode, errMsg)
	}

	// Covert pathParams of map[string]string type to map[string]interface{}
	pathParams := make(map[string]interface{}, len(pathParamsMapBuffer))
	for k, v := range queryParamsMapBuffer {
		pathParams[k] = v
	}
	parseCode, errMsg = recursiveAttributeCheck(requestEndpoint, ss.RequiredPathParams, pathParams, 0)
	if parseCode != ATTRIBUTE_OK {
		causePanic(PATH_PARAMS, parseCode, errMsg)
	}

	//TODO: implement proper identity parser
	return ServiceEvent{
		PathParams:  pathParams,
		RequestBody: requestBody,
		QueryParams: queryParams,
		Identity: Identity{
			Email:    "test@email.com",
			Username: "testusername",
			Role:     "testRole",
		},
	}
}

// Create new Required Event Attributes
func (ah AWSServiceHandler) NewReqEvenAttrib(dataType string, isRequired bool, minLength int, maxLength int) ReqEventAttrib {
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

/*
	@Internal Function
	Check request attributes deep. See if passed attribute match the specs
*/
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
					if !(len(attributes[k].(string)) >= reqAttributeMinLength && len(attributes[k].(string)) <
						reqAttributeMaxLength) {
						return INVALID_ATTRIBUTE_LENGTH_ERROR, fmt.Sprintf(
							"invalid length of attribute %v. expected",
							k,
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
						if !(attributes[k].(int) >= reqAttributeMinLength && attributes[k].(int) <
							reqAttributeMaxLength) {
							return INVALID_ATTRIBUTE_LENGTH_ERROR, fmt.Sprintf(
								"invalid length of attribute %v. expected",
								k,
							)
						}
					}
					if foundAttribType == "float64" {
						if !(attributes[k].(float64) >= float64(reqAttributeMinLength) && attributes[k].(float64) <
							float64(reqAttributeMaxLength)) {
							return INVALID_ATTRIBUTE_LENGTH_ERROR, fmt.Sprintf(
								"invalid length of attribute %v. expected",
								k,
							)
						}
					}
					if foundAttribType == "float32" {
						if !(attributes[k].(float32) >= float32(reqAttributeMinLength) && attributes[k].(float32) <
							float32(reqAttributeMaxLength)) {
							return INVALID_ATTRIBUTE_LENGTH_ERROR, fmt.Sprintf(
								"invalid length of attribute %v. expected",
								k,
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

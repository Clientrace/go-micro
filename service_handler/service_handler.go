package service_handler

import (
	"fmt"
	"reflect"
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
	PATH_PARAM   = iota
)

/* Required Event Specification Attribute */
type ReqEventAttrib struct {
	DataType   string
	IsRequired bool
	MinLength  int
	MaxLength  int
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
	Identity
	RequestBody map[string]interface{}
	QueryParams map[string]interface{}
	PathParams  map[string]interface{}
}

type Service struct {
}

type ServiceHandler interface {
	NewService(ServiceEvent)
}

/*
	@Exported Function
	Parse AWS Event to get the identity and requests objects
	event AWS  HTTP Event
	requestFmt Required event body format
*/
// func (sh ServiceEvent) NewService(event events.APIGatewayProxyRequest,
// 	bodyFmt ReqEventSpec, queryPrmsFmt ReqEventSpec, pathPrmsFmt ReqEventSpec) (ServiceEvent, error) {
// 	requestEndpoint := event.RequestContext.ResourcePath

// 	var requestBody map[string]interface{}
// 	queryParams := event.QueryStringParameters
// 	pathParams := event.PathParameters

// 	// Convert JSON String body to map
// 	json.Unmarshal([]byte(event.Body), &requestBody)

// 	parseCode, err_message := recursiveAttributeCheck(requestEndpoint, bodyFmt, requestBody, 0)
// 	if parseCode != ATTRIBUTE_OK {
// 		panic(map[string]interface{}{
// 			"paramType": "REQ_BODY",
// 			"errorMsg":  err_message,
// 			"parseCode": parseCode,
// 		})
// 	}

// }

// Create new Required Event Attributes
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

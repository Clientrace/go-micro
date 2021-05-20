package service_handler

import (
	"fmt"
	"reflect"
)

/* Error Codes */
const ATTRIBUTE_OK string = "ATTRIBUTE_OK"
const MISSING_ATTRIBUTE_ERROR string = "MISSING_ATTRIBUTE"
const INVALID_ATTRIBUTE_TYPE_ERROR string = "INVALID_ATTRIBUTE_TYPE"

/* Type Map Reference */
var ATTRIB_TYPE_MAP = map[string]interface{}{
	"string":  []string{"string"},
	"boolean": []string{"bool"},
	"number":  []string{"int", "float64"},
}

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
	ParseEvent(ServiceEvent)
}

/*
	@Exported Function
	Parse AWS Event to get the identity and requests objects
	event AWS  HTTP Event
	requestFmt Required event body format
*/
// func (sh ServiceEvent) ParseEvent(event events.APIGatewayProxyRequest, requestFmt map[string]interface{}) ServiceEvent {
// HTTP Request URL Endpiont
// requestEndpoint := event.RequestContext.ResourcePath
// }

/*
	@Internal Function
	Check request attributes deep. See if passed attribute match the specs
*/
func recursiveAttributeCheck(endpoint string, reqEventSpec ReqEventSpec, attributes map[string]interface{}, depth int) (string, error) {
	fmt.Println("\tObject Depth:", depth)
	rqa := reqEventSpec.ReqEventAttributes

	// Iterate through required attributes
	for k := range rqa {
		if val, ok := attributes[k]; ok {
			fmt.Println("\tChecking Required Attribute ", k, val)
			if reflect.TypeOf(rqa[k]).Kind().String() == "map" {
				if reflect.TypeOf(attributes[k]).String() == "map[string]interface {}" {
					fmt.Println("\t> Parsing Parent Attribute")
					// Create ReqEventSpec for child attribute
					reqEventChild := ReqEventSpec{
						ReqEventAttributes: rqa[k].(map[string]interface{}),
					}
					return recursiveAttributeCheck(endpoint, reqEventChild, attributes[k].(map[string]interface{}), depth+1)
				} else {
					return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Errorf(
						"invalid type of attribute '%v'. expected object got %v",
						k,
						reflect.TypeOf(attributes[k]),
					)
				}

			} else {
				fmt.Println("\t\t> Parsing Child Attribute", k)

				// Required VS Found Attribute Type
				reqAttributeType := rqa[k].(ReqEventAttrib).DataType
				reqAttributeMaxLength := rqa[k].(ReqEventAttrib).MaxLength
				reqAttributeMinLength := rqa[k].(ReqEventAttrib).MinLength
				foundAttribType := reflect.TypeOf(attributes[k]).String()

				if reqAttributeType == "string" && foundAttribType == "string" {
					if !(len(attributes[k].(string)) >= reqAttributeMinLength && len(attributes[k].(string)) <
						reqAttributeMaxLength) {
						return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Errorf(
							"invalid length of attribute %v. expected",
							k,
						)
					}
				} else {
					return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Errorf(
						"invalid type of attribute %v. expected string, got %v",
						k,
						reflect.TypeOf(attributes[k]).String(),
					)
				}
				if reqAttributeType == "number" && (foundAttribType == "int" || foundAttribType == "float32" ||
					foundAttribType == "float64") {
					if !(len(attributes[k].(string)) >= reqAttributeMinLength && len(attributes[k].(string)) <
						reqAttributeMaxLength) {
						return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Errorf(
							"invalid length of attribute %v. expected",
							k,
						)
					}
				} else {
					return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Errorf(
						"invalid type of attribute %v. expected number(int or float) got %v",
						k,
						attributes[k],
					)
				}
			}
		} else {
			if reflect.TypeOf(rqa[k]).String() == "map[string]interface {}" {
				return MISSING_ATTRIBUTE_ERROR, fmt.Errorf("missing attribute '%v'", k)
			} else {
				if rqa[k].(ReqEventAttrib).IsRequired {
					return MISSING_ATTRIBUTE_ERROR, fmt.Errorf("missing attribute '%v'", k)
				}
			}
		}
	}

	return ATTRIBUTE_OK, nil
}

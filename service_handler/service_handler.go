package service_handler

import (
	"fmt"
	"reflect"
)

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

/* @Internal Function */
func recursiveAttributeCheck(endpoint string, reqEventSpec ReqEventSpec, attributes map[string]interface{}) (string, error) {
	rqa := reqEventSpec.ReqEventAttributes

	// Iterate through required attributes
	for k := range rqa {
		if val, ok := attributes[k]; ok {
			fmt.Println("\t", val, "==>", reflect.Map)
			if reflect.TypeOf(rqa[k]).Kind() == reflect.Map {
				fmt.Println("\t\t> Parsing Parent Attribute")
				if reflect.TypeOf(attributes[k]) == reflect.TypeOf(ReqEventAttrib{}) {
					fmt.Printf("Child")
				} else {
					return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Errorf(
						"invalid type of attribute '%v'. expected object got %v",
						k,
						reflect.TypeOf(attributes[k]),
					)
				}

			} else {
				fmt.Println("\t\t> Parsing Parent Attribute")
			}
		} else {
			return MISSING_ATTRIBUTE_ERROR, fmt.Errorf("missing attribute '%v'", k)
		}
	}
	// for k, v := range attributes {
	// 	fmt.Println(k, v)
	// }
	return ATTRIBUTE_OK, nil
}

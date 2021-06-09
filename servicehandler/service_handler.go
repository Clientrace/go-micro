package servicehandler

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
type EventSpec struct {
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
	NewServiceEvent(EventSpec) ServiceEvent
	NewHTTPResponse(ServiceResponse) interface{}
	HandleExceptions(interface{}, map[string]string) interface{}
}

var parameterMap = map[int]string{
	QUERY_PARAMS: "Query Parameter",
	PATH_PARAMS:  "Path Parameter",
	REQ_BODY:     "Request Body",
}

var errMsgMap = map[int]string{
	MISSING_ATTRIBUTE_ERROR:        "MISSING ATTRIBUTE ERROR",
	INVALID_ATTRIBUTE_LENGTH_ERROR: "INVALID ATTRIBUTE LENGTH",
	INVALID_ATTRIBUTE_TYPE_ERROR:   "INVALID ATTRIBUTE TYPE",
}

// causePanic will raise an http exception via that'll cause a panic and should be recovered
func causePanic(paramType int, parseCode int, errorMsg string) {
	RaiseHTTPException(
		BAD_REQUEST,
		fmt.Sprintf("Error in %v, %v. %v", parameterMap[paramType], errMsgMap[parseCode], errorMsg),
	)
}

// NewReqEventAttrib will create a new ReqEventAttrib object
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

// isInRange will check if in is between the given range
func isInRange(in interface{}, min int, max int) bool {
	vType := reflect.TypeOf(in).String()
	switch vType {
	case "string":
		vLen := len(in.(string))
		if !(vLen >= min && vLen <= max) {
			return false
		}
	case "int":
		vLen := in.(int)
		if !(vLen >= min && vLen <= max) {
			return false
		}
	case "float32":
		vLen := in.(float32)
		if !(vLen >= float32(min) && vLen <= float32(max)) {
			return false
		}
	case "float64":
		vLen := in.(float64)
		if !(vLen >= float64(min) && vLen <= float64(max)) {
			return false
		}
	}
	return true
}

// attribCheck is a helper function of recursiveAttributeCheck that checks the attrib type
// and required min/max length.
func attribCheck(attribName string, rqa ReqEventAttrib, attribute interface{}) (int, string) {
	reqType := rqa.DataType
	raMaxLen := rqa.MaxLength
	raMinLen := rqa.MinLength
	gotType := reflect.TypeOf(attribute).String()
	resCode := ATTRIBUTE_OK
	retMsg := "OK"

	invalidTypeMsg := fmt.Sprintf(
		"invalid type of attribute %v. expected %v, got %v", attribName, reqType, gotType,
	)
	invalidLengthMsg := fmt.Sprintf(
		"invalid length of attribute %v. min length: %d, max length: %d", attribName, raMinLen, raMaxLen,
	)

	typeValidatorMap := map[string]interface{}{
		"string":  []string{"string"},
		"number":  []string{"int", "float32", "float64"},
		"boolean": []string{"bool"},
	}

	typeValidator := func(got string, want string) bool {
		for _, v := range typeValidatorMap[want].([]string) {
			if v == got {
				return true
			}
		}
		return false
	}

	if !typeValidator(gotType, reqType) {
		resCode = INVALID_ATTRIBUTE_TYPE_ERROR
		retMsg = invalidTypeMsg
	} else if reqType == "string" || reqType == "number" {
		if !isInRange(attribute, raMinLen, raMaxLen) {
			resCode = INVALID_ATTRIBUTE_LENGTH_ERROR
			retMsg = invalidLengthMsg
		}
	}
	return resCode, retMsg
}

// recursiveAttributeCheck will check request attributes deep; see if passed attribute match the specs
func recursiveAttributeCheck(endpoint string, reqEventSpec ReqEventSpec, attributes map[string]interface{}, depth int) (int, string) {
	rqa := reqEventSpec.ReqEventAttributes
	retCode := ATTRIBUTE_OK
	retMsg := "OK"
	mapType := "map[string]interface {}"

	// Iterate through required attributes
	for k := range rqa {
		// Missing Attribute
		rqaType := reflect.TypeOf(rqa[k]).String()
		if _, ok := attributes[k]; !ok && (rqaType == mapType || (rqa[k].(ReqEventAttrib).IsRequired && rqaType != mapType)) {
			retCode = MISSING_ATTRIBUTE_ERROR
			retMsg = fmt.Sprintf("missing attribute '%v'", k)
			return retCode, retMsg
		}

		// Recurse through the parent attribute
		if reflect.TypeOf(rqa[k]).Kind().String() == "map" && reflect.TypeOf(attributes[k]).String() == mapType {
			retCode, retMsg = recursiveAttributeCheck(
				endpoint,
				ReqEventSpec{
					ReqEventAttributes: rqa[k].(map[string]interface{}),
				},
				attributes[k].(map[string]interface{}),
				depth+1,
			)
			if retCode != ATTRIBUTE_OK {
				return retCode, retMsg
			}
		} else if reflect.TypeOf(rqa[k]).Kind().String() == "map" && reflect.TypeOf(attributes[k]).String() != mapType {
			return INVALID_ATTRIBUTE_TYPE_ERROR, fmt.Sprintf(
				"invalid type of attribute '%v'. expected object got %v",
				k,
				reflect.TypeOf(attributes[k]),
			)
		} else if retCode, retMsg := attribCheck(k, rqa[k].(ReqEventAttrib), attributes[k]); retCode != ATTRIBUTE_OK {
			return retCode, retMsg
		}
	}
	return retCode, retMsg
}

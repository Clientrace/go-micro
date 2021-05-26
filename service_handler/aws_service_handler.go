package service_handler

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

/*
	AWS Implementation of Service Handler
*/
type AWSServiceHandler struct {
	event events.APIGatewayProxyRequest
}

/*
	Parse AWS Event to get the identity and requests objects
	and returns ServiceEvent object.
*/
func (ah AWSServiceHandler) NewService(ss ServiceSpec) ServiceEvent {
	requestEndpoint := ah.event.RequestContext.ResourcePath

	identity := ah.event.RequestContext.Identity

	var requestBody map[string]interface{}
	queryParamsMapBuffer := ah.event.QueryStringParameters
	pathParamsMapBuffer := ah.event.PathParameters

	// Convert JSON String body to map
	json.Unmarshal([]byte(ah.event.Body), &requestBody)

	parseCode, errMsg := recursiveAttributeCheck(requestEndpoint, ss.RequiredRequestBody, requestBody, 0)
	if parseCode != ATTRIBUTE_OK {
		log.Println("Invalid Request Body", errMsg)
		causePanic(REQ_BODY, parseCode, errMsg)
	}

	// Covert queryParamsBuffer of map[string]string type to map[string]interface{}
	queryParams := make(map[string]interface{}, len(queryParamsMapBuffer))
	for k, v := range queryParamsMapBuffer {
		queryParams[k] = v
	}
	parseCode, errMsg = recursiveAttributeCheck(requestEndpoint, ss.RequiredQueryParams, queryParams, 0)
	if parseCode != ATTRIBUTE_OK {
		log.Println("Invalid Query Params", errMsg)
		causePanic(QUERY_PARAMS, parseCode, errMsg)
	}

	// Covert pathParams of map[string]string type to map[string]interface{}
	pathParams := make(map[string]interface{}, len(pathParamsMapBuffer))
	for k, v := range pathParamsMapBuffer {
		pathParams[k] = v
	}
	parseCode, errMsg = recursiveAttributeCheck(requestEndpoint, ss.RequiredPathParams, pathParams, 0)
	if parseCode != ATTRIBUTE_OK {
		log.Println("Invalid Path Params", errMsg)
		causePanic(PATH_PARAMS, parseCode, errMsg)
	}

	return ServiceEvent{
		PathParams:  pathParams,
		RequestBody: requestBody,
		QueryParams: queryParams,
		Identity:    identity,
	}
}
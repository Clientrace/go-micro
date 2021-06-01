package servicehandler

import (
	"encoding/json"
	"go-micro/logger"
	"log"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
)

// AWSServiceHandler - AWS Implement of Service Handler
type AWSServiceHandler struct {
	Event  events.APIGatewayProxyRequest
	Logger logger.Logger
}

// AWSServiceHandler.NewService - Creates new AWSServiceHandler instance
func (ah AWSServiceHandler) NewService(ss ServiceSpec) ServiceEvent {
	ah.Logger.Log(logger.INFO, "Creating new service", "aws_service_handler", map[string]interface{}{})

	requestEndpoint := ah.Event.RequestContext.ResourcePath

	identity := ah.Event.RequestContext.Identity

	var requestBody map[string]interface{}
	queryParamsMapBuffer := ah.Event.QueryStringParameters
	pathParamsMapBuffer := ah.Event.PathParameters

	// Convert JSON String body to map
	json.Unmarshal([]byte(ah.Event.Body), &requestBody)

	ah.Logger.Log(logger.INFO, "Parsing Request Body", "aws_service_handler", requestBody)
	parseCode, errMsg := recursiveAttributeCheck(requestEndpoint, ss.RequiredRequestBody, requestBody, 0)
	if parseCode != ATTRIBUTE_OK {
		ah.Logger.Log(logger.ERROR, "Invalid Request Body, "+errMsg, "aws_service_handler", nil)
		causePanic(REQ_BODY, parseCode, errMsg)
	}

	// Covert queryParamsBuffer of map[string]string type to map[string]interface{}
	ah.Logger.Log(logger.INFO, "Parsing Query Params", "aws_service_handler", requestBody)
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

func (ah AWSServiceHandler) NewHTTPResponse(sr ServiceResponse) interface{} {
	ah.Logger.Log(logger.INFO, "Creating New HTTP Response", "aws_service_handler", nil)
	return events.APIGatewayProxyResponse{
		StatusCode:      sr.StatusCode,
		IsBase64Encoded: false,
		Body:            sr.ReturnBody,
		Headers:         sr.ReturnHeaders,
	}
}

func (ah AWSServiceHandler) HandleExceptions(recoverPayload interface{}, returnHeaders map[string]string) interface{} {
	if recoverPayload != nil {
		if reflect.TypeOf(recoverPayload).String() != "servicehandler.HTTPException" {
			switch reflect.TypeOf(recoverPayload).String() {
			case "string":
				ah.Logger.Log(logger.FATAL, "Internal Server Error. "+recoverPayload.(string), "aws_service_handler", nil)
			case "runtime.errorString":
				errorString := recoverPayload.(error).Error()
				ah.Logger.Log(logger.FATAL, "Internal Server Error. "+errorString, "aws_service_handler", nil)
			case "error":
				errorString := recoverPayload.(error).Error()
				ah.Logger.Log(logger.FATAL, "Internal Server Error. "+errorString, "aws_service_handler", nil)
			}
			return ah.NewHTTPResponse(ServiceResponse{
				StatusCode:    int(INTERNAL_SERVER_ERROR),
				ReturnBody:    "Internal Server Error",
				ReturnHeaders: returnHeaders,
			}).(events.APIGatewayProxyResponse)
		}
		ah.Logger.Log(logger.FATAL, recoverPayload.(HTTPException).ErrorMessage, "aws_service_handler", nil)
		return ah.NewHTTPResponse(ServiceResponse{
			StatusCode:    recoverPayload.(HTTPException).StatusCode,
			ReturnBody:    recoverPayload.(HTTPException).ErrorMessage,
			ReturnHeaders: returnHeaders,
		}).(events.APIGatewayProxyResponse)
	}
	return nil
}

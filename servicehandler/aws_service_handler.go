package servicehandler

import (
	"encoding/json"
	"go-micro/logger"
	"reflect"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

// AWSServiceHandler is the aws implementation of ServiceHandler
type AWSServiceHandler struct {
	Event  events.APIGatewayProxyRequest
	Logger logger.Logger
}

// NewService will crete new AWSServiceHandler instance
func (ah AWSServiceHandler) NewServiceEvent(es EventSpec) ServiceEvent {
	ah.Logger.LogTxt(logger.INFO, "Creating new service")

	requestEndpoint := ah.Event.RequestContext.ResourcePath

	identity := ah.Event.RequestContext.Identity

	var requestBody map[string]interface{}
	queryParamsMapBuffer := ah.Event.QueryStringParameters
	pathParamsMapBuffer := ah.Event.PathParameters

	// Convert JSON String body to map
	json.Unmarshal([]byte(ah.Event.Body), &requestBody)

	ah.Logger.LogObj(logger.INFO, "Parsing Request Body", requestBody, "", false)
	parseCode, errMsg := recursiveAttributeCheck(requestEndpoint, es.RequiredRequestBody, requestBody, 0)
	if parseCode != ATTRIBUTE_OK {
		ah.Logger.LogTxt(logger.ERROR, "Invalid Request Body, "+errMsg)
		causePanic(REQ_BODY, parseCode, errMsg)
	}

	// Covert queryParamsBuffer of map[string]string type to map[string]interface{}
	queryParams := make(map[string]interface{}, len(queryParamsMapBuffer))
	for k, v := range queryParamsMapBuffer {
		queryParams[k] = v
	}
	ah.Logger.LogObj(logger.INFO, "Parsing Query Params", queryParams, "", false)
	parseCode, errMsg = recursiveAttributeCheck(requestEndpoint, es.RequiredQueryParams, queryParams, 0)
	if parseCode != ATTRIBUTE_OK {
		ah.Logger.LogTxt(logger.ERROR, "Invalid Query Params, "+errMsg)
		causePanic(QUERY_PARAMS, parseCode, errMsg)
	}

	// Covert pathParams of map[string]string type to map[string]interface{}
	pathParams := make(map[string]interface{}, len(pathParamsMapBuffer))
	for k, v := range pathParamsMapBuffer {
		pathParams[k] = v
	}
	ah.Logger.LogObj(logger.INFO, "Parsing Path Params", pathParams, "", false)
	parseCode, errMsg = recursiveAttributeCheck(requestEndpoint, es.RequiredPathParams, pathParams, 0)
	if parseCode != ATTRIBUTE_OK {
		ah.Logger.LogTxt(logger.ERROR, "Invalid Path params, "+errMsg)
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
	ah.Logger.LogTxt(
		logger.INFO,
		"Creating new HTTP Response. Status Code <"+strconv.Itoa(sr.StatusCode)+">. Return Body:\n "+sr.ReturnBody,
	)
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
				ah.Logger.LogTxt(
					logger.FATAL,
					"Internal Server Error. "+recoverPayload.(string),
				)
			case "runtime.errorString":
				errorString := recoverPayload.(error).Error()
				ah.Logger.LogTxt(
					logger.FATAL,
					"Internal Server Error. "+errorString,
				)
			}
			return ah.NewHTTPResponse(ServiceResponse{
				StatusCode:    int(INTERNAL_SERVER_ERROR),
				ReturnBody:    "Internal Server Error",
				ReturnHeaders: returnHeaders,
			}).(events.APIGatewayProxyResponse)
		}
		ah.Logger.LogTxt(
			logger.ERROR,
			recoverPayload.(HTTPException).ErrorMessage,
		)
		return ah.NewHTTPResponse(ServiceResponse{
			StatusCode:    recoverPayload.(HTTPException).StatusCode,
			ReturnBody:    recoverPayload.(HTTPException).ErrorMessage,
			ReturnHeaders: returnHeaders,
		}).(events.APIGatewayProxyResponse)
	}
	return nil
}

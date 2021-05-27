package service_handler

import "reflect"

const (
	BAD_REQUEST           = 400
	RESOURCE_CONFLICT     = 409
	INTERNAL_SERVER_ERROR = 500
)

type HTTPException struct {
	StatusCode   int
	ErrorMessage string
}

func HandleExceptions(sh ServiceHandler, sr map[string]string) interface{} {
	if err := recover(); err != nil {
		if reflect.TypeOf(err).String() != "service_handler.ServiceResponse" {
			return sh.NewHTTPResponse(ServiceResponse{
				StatusCode:    INTERNAL_SERVER_ERROR,
				ReturnBody:    "Internal Server Error",
				ReturnHeaders: sr,
			})
		}

		return sh.NewHTTPResponse(ServiceResponse{
			StatusCode:    err.(HTTPException).StatusCode,
			ReturnBody:    err.(HTTPException).ErrorMessage,
			ReturnHeaders: sr,
		})
	}
	return nil
}

package service_handler

type StatusCode int

const (
	BAD_REQUEST           StatusCode = 400
	RESOURCE_CONFLICT                = 409
	INTERNAL_SERVER_ERROR            = 500
)

type HTTPException struct {
	StatusCode   int
	ErrorMessage string
}

/* Check if status code is valie */
func (sc StatusCode) isValid() bool {
	switch sc {
	case BAD_REQUEST, RESOURCE_CONFLICT, INTERNAL_SERVER_ERROR:
		return true
	}
	return false
}

func RaiseHTTPException(sc StatusCode, errMsg string) {
	if !sc.isValid() {
		panic("Invalid Status Code")
	}
	panic(HTTPException{
		StatusCode:   int(sc),
		ErrorMessage: errMsg,
	})
}

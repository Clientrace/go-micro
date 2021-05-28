package service_handler

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func badRequestTest() (response events.APIGatewayProxyResponse) {
	returnHeaders := map[string]string{
		"Content-Type": "application/json",
	}
	sh := AWSServiceHandler{
		event: newAWSMockEvent(
			map[string]string{},
			map[string]string{},
			`
				{
				}
			`,
		),
	}

	sh.HandleExceptions(&response, returnHeaders)

	requestSpec := ServiceSpec{
		RequiredRequestBody: ReqEventSpec{
			ReqEventAttributes: map[string]interface{}{
				"username": map[string]interface{}{
					"firstName":  NewReqEvenAttrib("string", true, 4, 15),
					"lastName":   NewReqEvenAttrib("string", true, 4, 255),
					"middleName": NewReqEvenAttrib("string", true, 4, 255),
				},
			},
		},
		RequiredQueryParams: ReqEventSpec{},
		RequiredPathParams:  ReqEventSpec{},
	}

	service := sh.NewService(requestSpec)

	fmt.Println(service.PathParams)
	fmt.Println(service.QueryParams)
	fmt.Println(service.RequestBody)
	return response
}

func TestHandleExceptions(t *testing.T) {
	got := badRequestTest()
	fmt.Println("GOT", got)
	want := events.APIGatewayProxyResponse{
		StatusCode:      400,
		IsBase64Encoded: false,
		Body:            "Error in Request Body, MISSING ATTRIBUTE ERROR. missin attribute 'middleName'",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	if want.StatusCode != got.StatusCode || want.Headers["Content-Type"] != got.Headers["Content-Type"] ||
		want.Body != got.Body {
		fmt.Println(got)
		t.Errorf("Invalid Bad Request Return Output")
	}
}

func cPanic() {
	panic("test")
}

func testFunction() (s string) {
	ret := ""

	defer func() {
		if err := recover(); err != nil {
			ret = "RECOVERED"
		}
	}()
	cPanic()
	return ret
}

func TestDefer(t *testing.T) {
	fmt.Println("Return value", testFunction())
}

# Go Micro
### Serverless Go lang Microservice Template
This project is a module abstracting the [serverless](http://serverless.com/) app in Go to speed up the API creation of a microservice project. It provide structures for creating event specification that validates incomming requests for the format needed. It currently supports AWS for the infra.

## Creating a project
If you're new with serverless, click [here](https://www.serverless.com/blog/framework-example-golang-lambda-support)
```
go get https://github.com/Clientrace/go-micro.git
```

## Sample Usage
Say we want to create a service for creating a user in AWS Dynamodb. Your handler should look like this:

lambda_handlers/create_user/main.go
```
package create_user
import (
	"context"
	"github.com/Clientrace/go-micro/logger"
	"github.com/Clientrace/go-micro/servicehandler"
)

```
### **Service Function**
Creating the service implementation function. This is where the business logic is implemented. The function exepects a context, service event, and a loggger as its arguments. The servicehandler.ServiceEvent contains the properties extracted from the client side's API request (se.PathParams, se.QueryParams, se.RequestBody, se.Identity).
```
// createUserHandler is service implementation of create user
func createUserHandler(ctx context.Context, se servicehandler.ServiceEvent, lgr logger.Logger) string {
	// Initialize dyanmodb Repo
	lgr.LogTxt(logger.INFO, "Initializing dynamodb repo")
	r := repo.NewDynamodbRepo(ctx, lgr, se.Options)

	// Create new user
	lgr.LogTxt(logger.INFO, "Creating new service")
	s := createNewService(r, lgr)
	s.CreateUser(se)

	return `"message": "OK"`
}
```
### **Service Endpoint**
Creating the endpoing along with the service specification. This is where we define what the service is expecting the service event would look like. You can specify here the required request body, required path params or the required query params and the service_handler will counter check it versus all the aws event the service will receive.
```
// NewCreateUserEndpoint will generate the service endpoint for create user
func NewCreateUserEndpoint(options interface{}) *servicehandler.AWSServiceEndpoint {
	return servicehandler.NewServiceEndpoint(
		servicehandler.EventSpec{
			RequiredRequestBody: servicehandler.ReqEventSpec{
				ReqEventAttributes: map[string]interface{}{
					"firstName":     servicehandler.NewReqEvenAttrib("string", true, 4, 75),
					"lastName":      servicehandler.NewReqEvenAttrib("string", true, 4, 75),
					"emailAddress":  servicehandler.NewReqEvenAttrib("string", true, 8, 250),
				},
			},
		},
		createUserHandler,
		logger.NewLogger(),
		map[string]string{},
		options,
	)
}

```


### Running the Unit Tests
```
go test ./...
```

### Test Coverage
```
go test -coverprofile coverage.html ./...
go tool cover -html=coverage.html
```

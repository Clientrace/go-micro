package servicehandler

import (
	"go-micro/logger"

	"github.com/aws/aws-lambda-go/events"
)

type Endpoint interface {
	Execute(se ServiceEvent, logger logger.Logger) (response events.APIGatewayProxyResponse)
}

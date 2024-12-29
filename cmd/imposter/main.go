package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gatehill/imposter-go/internal/adapter/awslambda"
	"github.com/gatehill/imposter-go/internal/adapter/httpserver"
)

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(awslambda.HandleLambdaRequest)
	} else {
		httpserver.StartServer()
	}
}

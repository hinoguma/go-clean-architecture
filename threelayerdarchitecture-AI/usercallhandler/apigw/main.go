package main

import (
	"app/threelayerdarchitecture/setup"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	container := setup.NewContainer()
	router := NewRouter(container)

	lambda.Start(router.Route)
}

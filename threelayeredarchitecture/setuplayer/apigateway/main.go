package apigateway

import (
	"app/threelayeredarchitecture/appinfralayer"
	"app/threelayeredarchitecture/config"
	"app/threelayeredarchitecture/crosscuttinglayer"
	"app/threelayeredarchitecture/crosscuttinglayer/infra"
	"app/threelayeredarchitecture/setuplayer/registry"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// logger, timer, id generator
	crosscuttinglayer.SetGlobalTimeGenerator(
		crosscuttinglayer.NewStdTimeGenerator(),
	)
	crosscuttinglayer.SetGlobalTempIDGenerator(
		infra.NewUUIDV4Generator(),
	)

	// db setup
	dbConf := appinfralayer.DBConfig{
		Host:     config.GetDBHost(),
		Port:     config.GetDBPort(),
		User:     config.GetDBUser(),
		Password: config.GetDBPassword(),
		DBName:   config.GetDBName(),
		SSLMode:  config.GetSSLMode(),
	}
	err := appinfralayer.InitDB(dbConf)
	if err != nil {
		panic(err)
	}
	defer appinfralayer.CloseDB()

	// registry
	registry.Init(appinfralayer.GetPostgreDB())

	// cold start done
	// start lambda
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	switch event.RawPath {
	case "/transfer":
		if event.RequestContext.HTTP.Method == "POST" {
			return transferHandler(ctx, event)
		}
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 404,
		Body:       "Not Found",
	}, nil
}

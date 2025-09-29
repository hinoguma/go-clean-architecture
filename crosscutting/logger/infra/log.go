package infra

import (
	"app/crosscutting/logger/infrainterface"
	"app/crosscutting/utils"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type cloudWatchLogger struct {
	logGroupName  string
	logStreamName string
	client        cloudwatchlogs.Client
}

func (logger cloudWatchLogger) Info(ctx context.Context, req utils.LogRequest) {
	logjson := req.LogJson()
	logjson.WithLevel(utils.LogLevelInfo).WithRequestIDByCtx(ctx)
	putEvt := logger.logJsonToPutLogEventsInput(logjson)
	_, err := logger.client.PutLogEvents(ctx, &putEvt)
	if err != nil {
		return
	}
}

func (logger cloudWatchLogger) Error(ctx context.Context, req utils.LogRequest) {
	logjson := req.LogJson()
	logjson.WithLevel(utils.LogLevelError).WithRequestIDByCtx(ctx)
	putEvt := logger.logJsonToPutLogEventsInput(logjson)
	_, err := logger.client.PutLogEvents(ctx, &putEvt)
	if err != nil {
		return
	}
}

func (logger cloudWatchLogger) Debug(ctx context.Context, req utils.LogRequest) {
	logjson := req.LogJson()
	logjson.WithLevel(utils.LogLevelDebug).WithRequestIDByCtx(ctx)
	putEvt := logger.logJsonToPutLogEventsInput(logjson)
	_, err := logger.client.PutLogEvents(ctx, &putEvt)
	if err != nil {
		return
	}
}

func NewCloudWatchLogger(ctx context.Context) (infrainterface.Logger, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := cloudwatchlogs.NewFromConfig(cfg)
	return cloudWatchLogger{client: *client}, nil
}

func (logger cloudWatchLogger) logJsonToPutLogEventsInput(lg utils.LogJson) cloudwatchlogs.PutLogEventsInput {
	return cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logger.logGroupName),
		LogStreamName: aws.String(logger.logStreamName),
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(lg.JsonString()),
				Timestamp: aws.Int64(utils.GetNowUnixTime().Int64())},
		},
	}
}

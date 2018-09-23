package main

import (
	"context"
	"encoding/json"
	"hackathon-iguatemi/iguatemi-server/config"
	"hackathon-iguatemi/iguatemi-server/database"
	"runtime"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/graphql-go/graphql"
)

type RequestBody struct {
	Query          string                 `json:"query"`
	VariableValues map[string]interface{} `json:"variables"`
	OperationName  string                 `json:"operationName"`
}

func executeQuery(ctx context.Context, request RequestBody, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		VariableValues: request.VariableValues,
		RequestString:  request.Query,
		OperationName:  request.OperationName,
		Context:        ctx,
	})
	return result
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestBody := RequestBody{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	c := context.WithValue(ctx, "id", "1")
	graphQLResult := executeQuery(c, requestBody, config.Schema)
	responseJSON, err := json.Marshal(graphQLResult)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	return events.APIGatewayProxyResponse{Body: string(responseJSON[:]), StatusCode: 200}, nil
}

func main() {
	database.LoadDatabase()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	lambda.Start(Handler)
}

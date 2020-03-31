package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Tai-Yaki/lambda-practice-go/handlers/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

type User struct {
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	CreatedTime time.Time `json:"created_time"`
	UpdateTime  time.Time `json:"updated_time"`
}

type Response struct {
	User string `json:"user"`
}

type request struct {
	UserID string `json:"userID"`
}

var DynamoDB db.DB

func init() {
	DynamoDB = db.New()
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// parsedRequest, err := parseUrlRequest(request)
	// parsedRequest := request.PathParameters
	// if err != nil {
	// 	return response(
	// 		http.StatusBadRequest,
	// 		errorResponseBody(err.Error()),
	// 	), nil
	// }

	user, err := DynamoDB.GetItem(request.PathParameters["userID"])
	if err != nil {
		return response(
			http.StatusNoContent,
			errorResponseBody(err.Error()),
		), nil
	}

	body, err := responseBody(user)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody(err.Error()),
		), nil
	}

	return response(http.StatusOK, body), nil
}

func response(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func parseUrlRequest(req events.APIGatewayProxyRequest) (*request, error) {
	if req.HTTPMethod != http.MethodGet {
		return nil, fmt.Errorf("use GET request")
	}

	var r request
	err := json.Unmarshal([]byte(req.Body), &r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse request")
	}

	return &r, nil
}

func responseBody(user db.User) (string, error) {
	response, err := json.Marshal(user)
	if err != nil {
		return "", nil
	}

	return string(response), nil
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}

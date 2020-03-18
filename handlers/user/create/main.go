package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Tai-Yaki/lambda-practice-go/handlers/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Email       string    `json:"Email"`
	Password    string    `json:"password"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

type Response struct {
	User string `json:"user"`
}

var DynamoDB db.DB

func init() {
	DynamoDB = db.New()
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, err := parseRequest(request)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(request.PathParameters["password"]), bcrypt.DefaultCost)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	user := db.User{
		UserID:      xid.New().String(),
		Name:        request.PathParameters["name"],
		Email:       request.PathParameters["email"],
		Password:    string(password_hash),
		CreatedTime: time.Now().UTC(),
		UpdatedTime: time.Now().UTC(),
	}

	_, err = DynamoDB.PutItem(user)
	if err != nil {
		return response(
			http.StatusInternalServerError,
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

func parseRequest(req events.APIGatewayProxyRequest) (string, error) {
	if req.HTTPMethod != http.MethodPost {
		return "", fmt.Errorf("use POST request")
	}

	if _, ok := req.PathParameters["Name"]; !ok {
		return "", fmt.Errorf("parameter error")
	}

	if _, ok := req.PathParameters["Email"]; !ok {
		return "", fmt.Errorf("parameter error")
	}

	if _, ok := req.PathParameters["Password"]; !ok {
		return "", fmt.Errorf("parameter error")
	}

	return "", nil
}

func response(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}

func responseBody(user db.User) (string, error) {
	response, err := json.Marshal(user)
	if err != nil {
		return "", nil
	}

	return string(response), nil
}

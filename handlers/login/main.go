package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Tai-Yaki/lambda-practice-go/handlers/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type request struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

type responseData struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

var DynamoDB db.DB

func init() {
	DynamoDB = db.New()
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	parsedRequest, err := parseRequest(request)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}
	user, err := DynamoDB.GetItem(parsedRequest.UserID)
	if err != nil {
		return response(
			http.StatusNoContent,
			errorResponseBody(err.Error()),
		), nil
	}

	err = passwordCompare(user.Password, parsedRequest.Password)
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

func parseRequest(req events.APIGatewayProxyRequest) (*request, error) {
	if req.HTTPMethod != http.MethodPost {
		return nil, fmt.Errorf("use POST request")
	}

	var r request
	err := json.Unmarshal([]byte(req.Body), &r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse request")
	}

	return &r, nil
}

func passwordCompare(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.Wrapf(err, "failed to match password")
		}
		return err
	}

	return nil
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
	var responseUser responseData
	responseUser.UserID = user.UserID
	responseUser.Email = user.Email
	responseUser.Name = user.Name
	response, err := json.Marshal(responseUser)
	if err != nil {
		return "", nil
	}

	return string(response), nil
}

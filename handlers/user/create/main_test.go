package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Tai-Yaki/lambda-practice-go/handler/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

const exitError = 1

func TestHandler(t *testing.T) {
	tests := []struct {
		name, email, password, method string
		status                        int
	}{
		{"test", "test@example.com", "aaaa1234", http.MethodPost, http.StatusOK},
		{"test", "test@example.com", "", http.MethodPost, http.StatusBadRequest},
		{"test", "test@example.com", "aaaa1234", http.MethodGet, http.StatusBadRequest},
	}

	for _, te := range tests {
		res, _ := handler(events.APIGatewayProxyRequest{
			HTTPMethod: te.method,
			Body: "{\"name\": \"" + te.name +
				"\", \"email\": \"" + te.email +
				", \"password\": \"" + te.password + "\"}",
		})

		if res.StatusCode != te.status {
			t.Errorf("ExitStatus=%d, want %d", res.StatusCode, te.status)
		}
	}
}

func TestMain(m *testing.M) {
	if err := prepare(); err != nil {
		fmt.Println(err)
		os.Exit(exitError)
	}
	exitCode := m.Run()
	if err := cleanUp(); err != nil {
		fmt.Println(err)
		os.Exit(exitError)
	}
	os.Exit(exitCode)
}

func prepare() error {
	DynamoDB = db.TestNew()

	ok, err := DynamoDB.LinkTableExists()
	if err != nil {
		return errors.Wrap(err, "failed to check table existence")
	}
	if ok {
		if err := DynamoDB.DeleteLinkTable(); err != nil {
			return errors.Wrap(err, "failed to delete link table")
		}
	}

	if err := DynamoDB.CreateLinkTable(); err != nil {
		return errors.Wrap(err, "failed tp create link table")
	}

	return nil
}

func cleanUp() error {
	ok, err := DynamoDB.LinkTableExists()
	if err != nil {
		return errors.Wrap(err, "failed to check table existence")
	}
	if ok {
		if err := DynamoDB.DeleteLinkTable(); err != nil {
			return errors.Wrap(err, "failed to delete link table")
		}
	}

	DynamoDB = db.DB{}

	return nil
}

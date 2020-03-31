package db

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

var (
	LinkTableName = os.Getenv("LINK_TABLE")
	Region        = os.Getenv("REGION")
)

type DB struct {
	Instance *dynamo.DB
}

type User struct {
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

func New() DB {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String(Region)})

	return DB{Instance: db}
}

func (db DB) GetItem(UserID string) (User, error) {
	var result User
	table := db.Instance.Table(LinkTableName)
	err := table.Get("UserID", UserID).One(&result)
	if err != nil {
		return result, errors.Wrapf(err, "failed to get item")
	}
	fmt.Println("LinkTableName: " + LinkTableName)
	return result, nil
}

func (db DB) PutItem(user User) (User, error) {
	table := db.Instance.Table(LinkTableName)
	err := table.Put(user).Run()
	if err != nil {
		return user, errors.Wrapf(err, "failed to put item")
	}

	return user, nil
}

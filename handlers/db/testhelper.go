package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

var (
	Endpoint = "http://localhost:8000"
)

type CreateTable struct {
	db            *DB
	tableName     string
	attribs       []*dynamodb.AttributeDefinition
	schema        []*dynamodb.KeySchemaElement
	globalIndices map[string]dynamodb.GlobalSecondaryIndex
	localIndices  map[string]dynamodb.LocalSecondaryIndex
	readUnits     int64
	writeUnits    int64
	ondemand      bool
	tags          []*dynamodb.Tag
	err           error
}

func TestNew() DB {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String(Region), Endpoint: aws.String(Endpoint)})

	return DB{Instance: db}
}

func (d DB) CreateLinkTable() error {
	if err := d.Instance.CreateTable(LinkTableName, User{}).Run(); err != nil {
		return err
	}

	return nil
}

func (d DB) DeleteLinkTable() error {
	if err := d.Instance.Table(LinkTableName).DeleteTable().Run(); err != nil {
		return err
	}

	return nil
}

func (d DB) LinkTableExists() (bool, error) {
	output, err := d.Instance.ListTables().All()
	if err != nil {
		return false, err
	}

	if contains(output, LinkTableName) {
		return true, nil
	}

	return false, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == "" {
			continue
		}

		if a == e {
			return true
		}
	}

	return false
}

package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var tableName = "TouristEmails"

type TouristEmailItem struct {
	TourDate      string   `json:"TourDate"`
	TouristEmails []string `json:"TourstEmails"`
}

func main() {
	//initialize an aws session
	AwsSess, err := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	//Create a DynamoDb client
	svc := dynamodb.New(AwsSess)
	// define the dynamo db table schema
	params := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("TourDate"),
				KeyType:       aws.String("HASH"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefintion{
			{
				AttributeName: aws.String("TourDate"),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}
	//Create a dynamo db table
	_, err = svc.CreateTable(params)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

}

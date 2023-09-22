package main

import (
	"fmt"
	"log"

	"github.com/Tour-Feedback/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	//initialize an aws session
	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	//Create a DynamoDb client
	svc := dynamodb.New(sess)

	// define the dynamo db table schema
	params := &dynamodb.CreateTableInput{
		TableName: aws.String(dynamodb.TableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("TourDate"),
				KeyType:       aws.String("HASH"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
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

	fmt.Printf("Table %s created succesfully. \n", dynamodb.TableName)

	//inserting dummy data to dynamodb
	sampleData := []dynamodb.TouristEmailItem{
		{
			TourDate: "2023-09-01",
			TouristEmails: []string{
				"danielavelarde4@gmail.com",
				"dani@dani.com",
				"worldtraveler@gmail.com",
				"tourst@test.com",
			},
		},
		{
			TourDate: "2023-09-02",
			TouristEmails: []string{
				"dani@dani.com",
				"newintech@tech.com",
				"danielavelarde4@gmail.com",
				"serge@serge.com",
			},
		},
	}

	for _, item := range sampleData {
		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			log.Fatalf("error Marshaling dynamodb item %v", err)
		}

		_, err = svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(dynamodb.TableName),
			Item:      av,
		})
		if err != nil {
			log.Fatalf("error inserting item into DynamoDb %v", err)
		}
	}
	fmt.Println("dummy data added into Dynamodb")
}
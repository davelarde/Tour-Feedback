package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	dynamoDBClient *dynamodb.DynamoDB
	tableName      = "TouristEmails"
)

type TouristEmailItem struct {
	TourDate      string   `json:"TourDate"`
	TouristEmails []string `json:"TouristEmails"`
}

func init() {
	//initialize an aws session
	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	//Create a DynamoDb client
	dynamoDBClient = dynamodb.New(sess)
}

// this function will tell lambda if the tour was completed or not
func CompleteTourHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tourDate := request.QueryStringParameters["tourDate"]
	toursCompleted := request.QueryStringParameters["completed"] // this is the parameter to indicate if the tour was completed.

	if tourDate == "" {
		log.Print("tour date parameter is missing")
		return events.APIGatewayProxyResponse{
			statusCode: 400,
			Body:       "Tour date parameter is missing",
		}, nil
	}
	//update tour status in dynamoDb
	err := updateTourStatus(tourDate, tourCompleted)
	if err != nil {
		log.Print("Error updating tour status")
		return events.APIGatewayProxyResponse{
			statusCode: 500,
			Body:       fmt.Sprintf("Error updating tour status %v", err),
		}, nil
	}
	//if the tour is completed, optionally record the email address in dynamo db
	if tourCompleted == "true" {
		emailAddresses := []string{"danielavelarde4@gmail.com", "dani@dani.com", "worldtraveler@gmail.com"}
		err := updateEmailAddresses(tourDate, emailAddresses)
		if err != nil {
			log.Print("error updating email addresses")
			return events, APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("error updating email addresses %v", err),
			}, nil
		}

	}
	log.Print("Tour Status updated successfully")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Tour status updated successfully",
	}, nil
}

func updateTourStatus(TourDate string, tourCompleted string) error {
	_, err := dynamoDBClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"TourDate": {
				S: aws.String(tourDate),
			},
		},
		UpdateExpression: aws.String("Set TourStatus = :status")
		ExpressionAttributeValues : map[string]*dynamodb.AttributeValue{
			":status":{
				S: aws.String(tourCompleted),
			},
		},
	})
	if err != nil{
		log.Print(err)
		return err
	}
	return nil
}



func main() {

	// define the dynamo db table schema
	params := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
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

	fmt.Printf("Table %s created succesfully. \n", tableName)

	//inserting dummy data to dynamodb
	sampleData := []TouristEmailItem{
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
			TableName: aws.String(tableName),
			Item:      av,
		})
		if err != nil {
			log.Fatalf("error inserting item into DynamoDb %v", err)
		}
	}
	fmt.Println("dummy data added into Dynamodb")

}

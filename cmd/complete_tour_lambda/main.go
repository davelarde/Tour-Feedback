package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var (
	dynamoDBClient *dynamodb.DynamoDB
	tableName      = "TouristEmails"
)

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
		UpdateExpression: aws.String("SET TourStatus = :status"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(tourCompleted),
			},
		},
	})
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func updateEmailAddresses(tourDate string, emailAddresses []string) error {
	_, err := dynamoDBClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"TourDate": {
				S: aws.String(tourDate),
			},
		},
		UpdateExpression: aws.String("SET TouristEmails = list_append(TouristEmails, :emails)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":emails": {
				SS: aws.StringSlice(emailAddresses),
			},
		},
	})
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func main() {
	lambda.Start(CompleteTourHandler)

}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/davelarde/Tour-Feedback/dynamodbfol"
)

func main() {
	lambda.Start(CompleteTourHandler)

}

// this function will tell lambda if the tour was completed or not
func CompleteTourHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tourDate := request.QueryStringParameters["tourDate"]
	tourCompleted := request.QueryStringParameters["completed"] // this is the parameter to indicate if the tour was completed.

	if tourDate == "" {
		log.Print("tour date parameter is missing")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Tour date parameter is missing",
		}, nil
	}
	//update tour status in dynamoDb
	err := updateTourStatus(tourDate, tourCompleted)
	if err != nil {
		log.Print("Error updating tour status")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error updating tour status %v", err),
		}, nil
	}
	//if the tour is completed, optionally record the email address in dynamo db
	if tourCompleted == "true" {
		emailAddresses := []string{"danielavelarde44@gmail.com", "dani@dani.com", "worldtraveler@gmail.com"}
		err := dynamodbfol.UpdateEmailAddresses(tourDate, emailAddresses)
		if err != nil {
			log.Print("error updating email addresses")
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("error updating email addresses %v", err),
			}, nil
		}
		//send survey emails to the recorded emails
		err = sendSurveyEmail(emailAddresses)
		if err != nil {
			log.Print("Error sending survey emails")
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Error sending survey emails %v", err),
			}, nil

		}
	}
	log.Print("Tour Status updated successfully")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Tour status updated successfully",
	}, nil
}

func updateTourStatus(tourDate string, tourCompleted string) error {
	_, err := dynamodbfol.DynamoDBClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(dynamodbfol.TableName),
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

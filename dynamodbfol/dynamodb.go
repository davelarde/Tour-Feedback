package dynamodbfol

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var TableName = "TouristEmails"

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
	DynamoDBClient = dynamodb.New(sess)
}

var DynamoDBClient *dynamodb.DynamoDB

func UpdateEmailAddresses(tourDate string, emailAddresses []string) error {
	_, err := DynamoDBClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(TableName),
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

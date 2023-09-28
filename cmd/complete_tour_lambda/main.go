package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	// "github.com/aws/aws-sdk-go/service/lambda"
)

var dummyData = map[string][]string{
	"2023-09-01": {"danielavelarde4@gmail.com", "dan@dani.com", "serge@worldtraveller.com"},
	"2023-09-02": {"danielavelarde4@gmail.com", "tourist2@dani.com", "tourist3@worldtraveller.com"},
}

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
	//check if tour date exists in the dummy data
	touristEmails, ok := dummyData[tourDate]
	if !ok {
		log.Print("Tour date not found")
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Tour date not found",
		}, nil
	}

	log.Printf("Tour status for %s updated succesfully", tourDate)
	//if the tour is completed, optionally record the email address in dynamo db
	if tourCompleted == "true" {
		err := SendSurveyEmail(touristEmails)
		if err != nil {
			log.Printf("error sending survey emails")
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("error sending survey email %v", err),
			}, nil
		}
	}

	log.Print("Tour status updated succesfully")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Tour status updated succesfully",
	}, nil

}

func SendSurveyEmail(emailAddresses []string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatalf("error creating session")
	}

	svc := ses.New(sess)

	subject := "Tour Survey"
	body := "hello please take our tour survey"

	for _, emailAddress := range emailAddresses {
		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				ToAddresses: []*string{aws.String(emailAddress)},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Text: &ses.Content{
						Data: aws.String(body),
					},
				},
				Subject: &ses.Content{
					Data: aws.String(subject),
				},
			},
			Source: aws.String("danielavelarde4@gmail.com"),
		}
		// send email
		_, err := svc.SendEmail(input)
		if err != nil {
			log.Printf("Error sending email to %s: %v", emailAddress, err)
			return err
		}
	}
	return nil
}

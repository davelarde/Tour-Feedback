package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

	fmt.Println(AwsSess)
}

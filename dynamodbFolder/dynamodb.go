package dynamodb

var TableName = "TouristEmails"

type TouristEmailItem struct {
	TourDate      string   `json:"TourDate"`
	TouristEmails []string `json:"TouristEmails"`
}

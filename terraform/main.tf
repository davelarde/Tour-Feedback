provider "aws" {
    region = "us-east-1"
  
}
//create an Iam role for the lambda function
//create lambda func
resource "aws_lambda_function" "complete_tour_lambda" {
    function_name = "completeTourLambda"
    runtime = "go1.x"
    handler = "CompleteTourHandler"
    role = aws_iam_role.lambda_role.arn
    filename = "../cmd/complete_tour_lambda/main.zip"
    source_code_hash = filebase64sha256("../cmd/complete_tour_lambda/main.zip")
  
}
resource "aws_iam_role" "lambda_role" {
  name = "lambda-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement =[
        {
        Action = "sts:AssumeRole",
        Effect = "Allow"
        Principal = {
            Service = "lambda.amazonaws.com"
        }
        }
    ]
  })

}

resource "aws_iam_policy" "lambda_policy" {
  name = "lambda-policy"
  description = "IAM policy for lambda function"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement =[
        {
        Action = ["lambda:InvokeFunction"],
        Effect = "Allow"
        Resource = aws_lambda_function.complete_tour_lambda.arn
       
        },
    ]
  })

}
// attach an aws managed policy
resource "aws_iam_policy_attachment" "lambda_attachment" {
    name = "lambda-policy-attachment"
    policy_arn = aws_iam_policy.lambda_policy.arn
    roles = [aws_iam_role.lambda_role.name]
  
}


resource "aws_dynamodb_table" "tourist_emails" {
  name = "TouristEmails"
  billing_mode = "PROVISIONED"
  read_capacity = 1
  write_capacity = 1
 
  attribute {
    name = "TourDate"
    type = "S" 
  }
  hash_key = "TourDate"
  
}

output "table_name"{
    value = aws_dynamodb_table.tourist_emails.name
}

output "lambda_function_arn" {
  value = aws_lambda_function.complete_tour_lambda.arn
}
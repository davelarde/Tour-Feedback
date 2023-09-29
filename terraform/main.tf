provider "aws" {
    region = "us-east-1"
  
}
//create an Iam role for the lambda function
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


//define IAM policy for SES permissions
resource "aws_iam_policy" "ses_policy" {
  name = "ses-policy"
  description = "IAM policty for SES to send raw emails"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
        {
            Action = [
            "ses: SendRawEmail", 
            "ses:SendEmail"],
            Effect = "Allow",
            Resource = "*",
          
        }
    ]
  })
}
// attach Iam policy for SES to lambda execution role 
resource "aws_iam_role_policy_attachment" "ses_attachment" {
    policy_arn = aws_iam_policy.ses_policy.arn
    role = aws_iam_role.lambda_role.name
  
}
//define IAM policy for lambda function permission 
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

//create lambda func
resource "aws_lambda_function" "complete_tour_lambda" {
    function_name = "completeTourLambda"
    runtime = "go1.x"
    handler = "main"
    role = aws_iam_role.lambda_role.arn
    filename = "../cmd/main.zip"
    source_code_hash = filebase64sha256("../cmd/main.zip")
  
}
// attach an aws managed policy
resource "aws_iam_policy_attachment" "lambda_attachment" {
    name = "lambda-policy-attachment"
    policy_arn = aws_iam_policy.lambda_policy.arn
    roles = [aws_iam_role.lambda_role.name]
  
}



//create api gateway
resource "aws_api_gateway_rest_api" "complete_tour_api" {
    name = "complete-tour-api"
    description = "Complete tour api"
  
}

resource "aws_api_gateway_resource" "complete_tour" {
    rest_api_id = aws_api_gateway_rest_api.complete_tour_api.id
    parent_id = aws_api_gateway_rest_api.complete_tour_api.root_resource_id
    path_part = "completeTour"
  
}
resource "aws_api_gateway_method" "complete_tour" {
    rest_api_id = aws_api_gateway_rest_api.complete_tour_api.id
    resource_id = aws_api_gateway_resource.complete_tour.id
    http_method = "POST"
    authorization = "NONE"
  
}

resource "aws_api_gateway_integration" "complete_tour" {
  rest_api_id = aws_api_gateway_rest_api.complete_tour_api.id
  resource_id = aws_api_gateway_resource.complete_tour.id
  http_method = aws_api_gateway_method.complete_tour.http_method
  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = aws_lambda_function.complete_tour_lambda.invoke_arn
}

resource "aws_lambda_permission" "complete_tour" {
  statement_id = "AllowCompleteTourInvoke"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.complete_tour_lambda.function_name
  principal = "apigateway.amazonaws.com"
}

//deploy the api on stage
resource "aws_api_gateway_deployment" "complete_tour" {
  depends_on = [ aws_api_gateway_integration.complete_tour ]
  rest_api_id = aws_api_gateway_rest_api.complete_tour_api.id
  stage_name = "prod"
}



output "complete_tour_api_endpoint_url" {
  value = aws_api_gateway_deployment.complete_tour.invoke_url
}



output "lambda_function_arn" {
  value = aws_lambda_function.complete_tour_lambda.arn
}
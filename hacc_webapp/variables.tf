variable "hacc_app_username" {
  description = "Name of the SSM parameter to store HACC parameters"
  type        = string
}

variable "hacc_app_temp_password" {
  description = "Temporary password for the Cognito user pool"
  type        = string
}

variable "hacc_app_email" {
  description = "Email address for the Cognito user pool user"
  type        = string
}

variable "hacc_param_path" {
  description = "Path in SSM Parameter Store to store HACC parameters"
  type        = string
}

variable "aws_region" {
  description = "AWS region to deploy to"
  type        = string
}

variable "s3_bucket_name" {
  description = "Name of the S3 bucket to create for hosting the website"
  type        = string
}

variable "cloudfront_cache_policy_id" {
  description = "ID of the CloudFront cache policy to use for the distribution"
  type        = string
}

variable "cognito_first_auth_factors" {
  description = "List of allowed first authentication factors for Cognito user pool"
  type        = list(string)
}

variable "lambda_function_name" {
  description = "Name of the Lambda function to create for handling form submissions"
  type        = string
}

variable "lambda_handler" {
  description = "Handler for the Lambda function (e.g. 'index.handler')"
  type        = string
}

variable "lambda_runtime" {
  description = "Runtime for the Lambda function (e.g. 'python3.8')"
  type        = string
}
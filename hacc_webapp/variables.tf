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

variable "hacc_param_kms_key_id" {
  description = "KMS key ID to encrypt parameters in SSM Parameter Store (optional, default aws/ssm)"
  type        = string
  default     = "aws/ssm"
}

variable "aws_region" {
  description = "AWS region to deploy to"
  type        = string
}

variable "r53_zone_name" {
  description = "Name of existing Route53 public domain under which a zone will be created to manage DNS records for the website"
  type        = string
}

variable "r53_subdomain_name" {
  description = "Subdomain name for the website (e.g. 'hacc-app.example.com')"
  type        = string
}

variable "s3_bucket_name" {
  description = "Name of the S3 bucket to create for hosting the website"
  type        = string
}

variable "cloudfront_url_obfuscation_string" {
  description = "Path prefix to add to CloudFront distribution for obfuscation (e.g. 'asfdasjhlkjdasfklj')"
  type        = string
  default     = ""
}

variable "cloudfront_cache_policy_id" {
  description = "ID of the CloudFront cache policy to use for the distribution"
  type        = string
}

variable "cloudfront_origin_policy_id" {
  description = "ID of the CloudFront origin request policy to use for the distribution"
  type        = string
}

variable "cognito_first_auth_factors" {
  description = "List of allowed first authentication factors for Cognito user pool"
  type        = list(string)
}

variable "cognito_access_token_validity" {
  description = "Validity period of the Cognito access token in minutes"
  type        = number
  default     = 10
}

variable "cognito_id_token_validity" {
  description = "Validity period of the Cognito ID token in minutes"
  type        = number
  default     = 10
}

variable "cognito_refresh_token_validity" {
  description = "Validity period of the Cognito refresh token in minutes (min 60)"
  type        = number
  default     = 60
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
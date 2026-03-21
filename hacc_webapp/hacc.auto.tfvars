hacc_app_username      = "" ## Fill in
hacc_app_email         = "" ## Fill in
hacc_app_temp_password = "Temp_p@ss_W3lcome" ## will be reset on first login - needs lowercase, uppercase, number, special, 12 char min
hacc_param_path        = "hacc-lambda"

aws_region                 = "us-east-1"
s3_bucket_name             = "hacc-website-bucket"
cloudfront_cache_policy_id = "" ## Fill in No-cache managed policy in AWS account (within CloudFront service -> policies)
cognito_first_auth_factors = ["PASSWORD"]
lambda_function_name       = "hacc-lambda"
lambda_handler             = "hacc_lambda.lambda_handler"
lambda_runtime             = "python3.14"

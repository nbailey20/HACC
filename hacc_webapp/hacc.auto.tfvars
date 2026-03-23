hacc_app_username      = ""                  ## Fill in
hacc_app_email         = ""                  ## Fill in
hacc_app_temp_password = "Temp_p@ss_W3lcome" ## will be reset on first login - needs lowercase, uppercase, number, special, 12 char min
hacc_param_path        = "hacc-lambda"

aws_region                        = "us-east-1"
r53_zone_name                     = ""       ## Fill in existing DNS zone where subdomain should be created
r53_subdomain_name                = ""       ## Fill in
cloudfront_url_obfuscation_string = "kljgkl" ## Fill in random hard-to-guess string
s3_bucket_name                    = "hacc-website-bucket"
cloudfront_cache_policy_id        = "" ## Fill in No-cache managed policy in AWS account (within CloudFront service -> policies)
cloudfront_origin_policy_id       = "" ## Fill in Managed AllViewerExceptHostHeader policy in AWS account
cognito_first_auth_factors        = ["PASSWORD"]
lambda_function_name              = "hacc-lambda"
lambda_handler                    = "hacc_lambda.lambda_handler"
lambda_runtime                    = "python3.14"

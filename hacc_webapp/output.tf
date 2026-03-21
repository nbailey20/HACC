output "website_url" {
  description = "URL of the CloudFront distribution for the website"
  value       = aws_cloudfront_distribution.website_distribution.domain_name
}

output "username" {
  description = "Username for the Cognito user pool"
  value       = var.hacc_app_username
}

output "temporary_password" {
  description = "Initial password for the Cognito user pool (must be changed on first login)"
  value       = var.hacc_app_temp_password
}

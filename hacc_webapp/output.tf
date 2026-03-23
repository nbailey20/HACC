output "website_url" {
  description = "Subdomain pointing to the CloudFront distribution for the website"
  value       = "${var.r53_subdomain_name}/${var.cloudfront_url_obfuscation_string}/"
}

output "username" {
  description = "Username for the Cognito user pool"
  value       = var.hacc_app_username
}

output "temporary_password" {
  description = "Initial password for the Cognito user pool (must be changed on first login)"
  value       = var.hacc_app_temp_password
}

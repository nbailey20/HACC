// Deploys a basic (but secure) website using:
//   S3/CloudFront with static frontend HTML/JS for HACC search functionality
//   API Gateway/Lambda backend with minimal IAM to invoke HACC
//   Route53 / ACM to put everything under same subdomain and enable HTTPS
//   Cognito authentication with MFA, validation performed in Lambda

data "aws_caller_identity" "current" {}

// S3 resources / frontend hosting for website
resource "aws_s3_bucket" "website_bucket" {
  bucket = var.s3_bucket_name
}

resource "aws_s3_bucket_public_access_block" "website_bucket_public_access_block" {
  bucket = aws_s3_bucket.website_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_object" "index_html" {
  bucket = aws_s3_bucket.website_bucket.id
  key    = "index.html"
  content = templatefile("${path.module}/src/index.html", {
    CLIENT_ID        = aws_cognito_user_pool_client.user_pool_client.id
    COGNITO_DOMAIN   = "${aws_cognito_user_pool_domain.domain.domain}.auth.${var.aws_region}.amazoncognito.com"
    FRONTEND_DOMAIN  = var.r53_subdomain_name
    OBFUSCATE_STRING = var.cloudfront_url_obfuscation_string
  })
  content_type = "text/html"
}

## Use OAC to securely allow CloudFront to access S3 bucket without making it public
resource "aws_s3_bucket_policy" "website_bucket_policy" {
  bucket = aws_s3_bucket.website_bucket.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          "Service" : "cloudfront.amazonaws.com"
        }
        Action   = "s3:GetObject"
        Resource = "${aws_s3_bucket.website_bucket.arn}/*"
        Condition = {
          StringEquals = {
            "aws:SourceArn" = aws_cloudfront_distribution.website_distribution.arn
          }
        }
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.website_bucket_public_access_block]
}


// Route53 and ACM resources for custom domain and HTTPS
data "aws_route53_zone" "existing" {
  name = var.r53_zone_name
}

resource "aws_acm_certificate" "website_cert" {
  domain_name       = var.r53_subdomain_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

## Adds CNAME record to R53 zone to prove domain ownership for ACM certificate validation
resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.website_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      type   = dvo.resource_record_type
      record = dvo.resource_record_value
    }
  }

  zone_id = data.aws_route53_zone.existing.zone_id
  name    = each.value.name
  type    = each.value.type
  records = [each.value.record]
  ttl     = 60
}

resource "aws_route53_record" "website_alias" {
  zone_id = data.aws_route53_zone.existing.zone_id
  name    = var.r53_subdomain_name
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.website_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.website_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}


// CloudFront distribution for the website
resource "aws_cloudfront_distribution" "website_distribution" {
  origin {
    domain_name              = aws_s3_bucket.website_bucket.bucket_regional_domain_name
    origin_id                = "HACC-S3-frontend-origin"
    origin_access_control_id = aws_cloudfront_origin_access_control.oac.id
    s3_origin_config {
      origin_access_identity = "" # REQUIRED when using OAC
    }
  }

  origin {
    domain_name = replace(aws_apigatewayv2_api.hacc_api.api_endpoint, "https://", "")
    origin_id   = "HACC-APIGW-backend-origin"
    origin_path = ""
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "HACC website distribution"
  default_root_object = "index.html"

  default_cache_behavior {
    target_origin_id       = "HACC-S3-frontend-origin"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    cache_policy_id        = var.cloudfront_cache_policy_id

    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.path_guard.arn
    }
  }

  ordered_cache_behavior {
    path_pattern             = "${var.cloudfront_url_obfuscation_string}/api/*"
    target_origin_id         = "HACC-APIGW-backend-origin"
    allowed_methods          = ["GET", "HEAD", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"]
    cached_methods           = ["GET", "HEAD", "OPTIONS"]
    cache_policy_id          = var.cloudfront_cache_policy_id
    viewer_protocol_policy   = "redirect-to-https"
    origin_request_policy_id = var.cloudfront_origin_policy_id

    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.strip_path_prefix.arn
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate.website_cert.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
      locations        = []
    }
  }

  aliases = [var.r53_subdomain_name]
}

resource "aws_cloudfront_origin_access_control" "oac" {
  name                              = "hacc-s3-oac"
  description                       = "Access control for S3"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_function" "path_guard" {
  name    = "hacc-app-path-guard"
  runtime = "cloudfront-js-1.0"

  code = <<EOF
function handler(event) {
  var request = event.request;
  var uri = request.uri;

  var prefix = '${var.cloudfront_url_obfuscation_string}';

  // Block everything outside prefix
  if (!uri.startsWith('/'+prefix+'/')) {
    return {
      statusCode: 404,
      statusDescription: 'Not Found'
    };
  }

  // Strip prefix before sending to origin
  var newUri = uri.substring(prefix.length+1);

  // Default to index.html
  if (newUri === '' || newUri === '/') {
    newUri = '/index.html';
  }

  request.uri = newUri;

  return request;
}
EOF
}

resource "aws_cloudfront_function" "strip_path_prefix" {
  name    = "hacc-app-strip-path-prefix"
  runtime = "cloudfront-js-1.0"
  code    = <<EOF
function handler(event) {
  var request = event.request;
  if (request.uri.startsWith('/${var.cloudfront_url_obfuscation_string}/api/')) {
    request.uri = request.uri.substring('/${var.cloudfront_url_obfuscation_string}'.length + 4);
  }
  return request;
}
EOF
}


// Cognito user pool, user, MFA-enabled login domain for authentication
// Immutable email address attribute used for account recovery and verification
resource "aws_cognito_user_pool" "user_pool" {
  name = "hacc-website-user-pool"

  admin_create_user_config {
    allow_admin_create_user_only = true
  }

  email_verification_subject = "Verify your email for HACC App"
  email_verification_message = "Your verification code is {####}"

  sign_in_policy {
    allowed_first_auth_factors = var.cognito_first_auth_factors
  }

  mfa_configuration = "ON"
  software_token_mfa_configuration {
    enabled = true
  }

  dynamic "password_policy" {
    for_each = contains(var.cognito_first_auth_factors, "PASSWORD") ? [1] : []
    content {
      minimum_length    = 12
      require_uppercase = true
      require_lowercase = true
      require_numbers   = true
      require_symbols   = true
    }
  }

  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }
}

resource "aws_cognito_user_pool_domain" "domain" {
  domain       = "hacc-auth-user-pool-domain" # must be globally unique
  user_pool_id = aws_cognito_user_pool.user_pool.id
}

resource "aws_cognito_user_pool_client" "user_pool_client" {
  name                                 = "hacc-website-user-pool-client"
  user_pool_id                         = aws_cognito_user_pool.user_pool.id
  callback_urls                        = ["https://${var.r53_subdomain_name}/${var.cloudfront_url_obfuscation_string}/api/callback"]
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  explicit_auth_flows                  = ["ALLOW_USER_AUTH"]
  allowed_oauth_scopes                 = ["openid", "email"]
  supported_identity_providers         = ["COGNITO"]
  generate_secret                      = false

  access_token_validity  = var.cognito_access_token_validity
  id_token_validity      = var.cognito_id_token_validity
  refresh_token_validity = var.cognito_refresh_token_validity
  token_validity_units {
    access_token  = "minutes"
    id_token      = "minutes"
    refresh_token = "minutes"
  }
}

resource "aws_cognito_user" "app_user" {
  user_pool_id       = aws_cognito_user_pool.user_pool.id
  username           = var.hacc_app_username
  temporary_password = var.hacc_app_temp_password
  attributes = {
    email          = var.hacc_app_email
    email_verified = "true"
  }
}


// IAM role and policy for Lambda to execute HACC
resource "aws_iam_role" "lambda_execution_role" {
  name = "hacc-lambda-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_policy" "lambda_execution_policy" {
  name        = "hacc-lambda-execution-policy"
  description = "Policy for Lambda function to call HACC executable"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        "Effect" : "Allow",
        "Action" : "logs:CreateLogGroup",
        "Resource" : "arn:aws:logs:${var.aws_region}:${data.aws_caller_identity.current.account_id}:*"
      },
      {
        "Effect" : "Allow",
        "Action" : [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        "Resource" : [
          "arn:aws:logs:${var.aws_region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${var.lambda_function_name}:*"
        ]
      },
      {
        "Sid" : "HaccSsmPerms",
        "Effect" : "Allow",
        "Action" : [
          "ssm:GetParametersByPath",
          "ssm:GetParameter"
        ],
        "Resource" : [
          "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/${var.hacc_param_path}/*",
          "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/${var.hacc_param_path}/"
        ]
      },
      {
        "Sid" : "HaccKmsPerms",
        "Effect" : "Allow",
        "Action" : [
          "kms:Encrypt",
          "kms:Decrypt"
        ],
        "Resource" : [
          "arn:aws:kms:${var.aws_region}:${data.aws_caller_identity.current.account_id}:key/${var.hacc_param_kms_key_id}",
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_execution_role.name
  policy_arn = aws_iam_policy.lambda_execution_policy.arn
}

resource "aws_lambda_function" "hacc_api_lambda" {
  function_name    = var.lambda_function_name
  runtime          = var.lambda_runtime
  handler          = var.lambda_handler
  timeout          = 30
  role             = aws_iam_role.lambda_execution_role.arn
  filename         = "src/build/hacc_lambda.zip"
  source_code_hash = filebase64sha256("src/build/hacc_lambda.zip")
  environment {
    variables = {
      HACC_CONFIG          = "./config.yaml"
      COGNITO_REDIRECT_URI = "https://${var.r53_subdomain_name}/${var.cloudfront_url_obfuscation_string}/api/callback"
      COGNITO_CLIENT_ID    = aws_cognito_user_pool_client.user_pool_client.id
      COGNITO_USER_POOL_ID = aws_cognito_user_pool.user_pool.id
      COGNITO_TOKEN_URL    = "https://${aws_cognito_user_pool_domain.domain.domain}.auth.${var.aws_region}.amazoncognito.com/oauth2/token"
      FRONTEND_URL         = "https://${var.r53_subdomain_name}/${var.cloudfront_url_obfuscation_string}/"
    }
  }
}

resource "aws_lambda_permission" "apigw_invoke" {
  for_each = toset(["search", "callback", "auth-check"])

  statement_id  = "AllowAPIGatewayInvoke-${each.key}"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.hacc_api_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.hacc_api.execution_arn}/*/*/${each.key}"
}


// API Gateway to expose the Lambda function as an API endpoint
resource "aws_apigatewayv2_api" "hacc_api" {
  name          = "hacc-website-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id             = aws_apigatewayv2_api.hacc_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.hacc_api_lambda.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "authcheck_route" {
  api_id    = aws_apigatewayv2_api.hacc_api.id
  route_key = "GET /auth-check"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "callback_route" {
  api_id    = aws_apigatewayv2_api.hacc_api.id
  route_key = "GET /callback"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "search_route" {
  api_id    = aws_apigatewayv2_api.hacc_api.id
  route_key = "POST /search"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.hacc_api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_apigatewayv2_authorizer" "cognito_authorizer" {
  api_id          = aws_apigatewayv2_api.hacc_api.id
  name            = "CognitoAuthorizer"
  authorizer_type = "JWT"

  identity_sources = ["$request.header.Authorization"]

  jwt_configuration {
    audience = [aws_cognito_user_pool_client.user_pool_client.id]
    issuer   = "https://cognito-idp.${var.aws_region}.amazonaws.com/${aws_cognito_user_pool.user_pool.id}"
  }
}

## Enhancements to consider:
## add obfuscation path to app under subdomain (e.g. https://hacc-app.nick-bailey.com/asfdasjhlkjdasfklj)

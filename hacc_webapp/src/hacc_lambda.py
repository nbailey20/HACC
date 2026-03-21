import os
import subprocess
import json
import urllib.request
import urllib.parse
import urllib.error
from jwt import PyJWKClient, decode


# ========== CONFIG ==========
COGNITO_CLIENT_ID    = os.environ.get("COGNITO_CLIENT_ID")
COGNITO_TOKEN_URL    = os.environ.get("COGNITO_TOKEN_URL")
COGNITO_REDIRECT_URI = os.environ.get("COGNITO_REDIRECT_URI")
COGNITO_USER_POOL_ID = os.environ.get("COGNITO_USER_POOL_ID")
FRONTEND_URL         = os.environ.get("FRONTEND_URL")


# ========= COOKIE HELPERS ==========
def parse_cookies(cookie_header):
    """Return dict of cookie_name -> value"""
    cookies = {}
    if not cookie_header:
        return cookies
    for part in cookie_header.split(";"):
        if "=" in part:
            name, value = part.strip().split("=", 1)
            cookies[name] = value
    return cookies

def set_cookie_header(name, value, max_age=600):
    return f"{name}={value}; HttpOnly; Secure; Path=/; Max-Age={max_age}; SameSite=None"


# ========= TOKEN HELPERS ==========
def exchange_code_for_tokens(code):
    """Call Cognito /oauth2/token"""
    data = {
        "grant_type": "authorization_code",
        "client_id": COGNITO_CLIENT_ID,
        "code": code,
        "redirect_uri": COGNITO_REDIRECT_URI,
    }

    encoded_data = urllib.parse.urlencode(data).encode('utf-8')
    req = urllib.request.Request(COGNITO_TOKEN_URL, data=encoded_data, headers={"Content-Type": "application/x-www-form-urlencoded"})
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            if resp.status != 200:
                raise Exception(f"Token exchange failed: HTTP {resp.status}")
            return json.loads(resp.read().decode('utf-8'))
    except urllib.error.HTTPError as e:
        body = e.read().decode()
        raise Exception(f"Token exchange failed: HTTP {e.code}: {e.reason} {body}")


def validate_token(token):
    """Validate JWT token signature and claims using PyJWT"""
    if not token:
        print("No token provided for validation")
        return False
    try:
        # Extract region from user pool ID (format: region_poolId)
        region = COGNITO_USER_POOL_ID.split('_')[0] if COGNITO_USER_POOL_ID else 'us-east-1'
        issuer = f"https://cognito-idp.{region}.amazonaws.com/{COGNITO_USER_POOL_ID}"
        jwks_url = f"{issuer}/.well-known/jwks.json"
        
        # PyJWT handles JWKS fetching and key selection
        jwks_client = PyJWKClient(jwks_url)
        signing_key = jwks_client.get_signing_key_from_jwt(token)
        
        # Decode and validate all claims in one call
        decode(
            token,
            signing_key.key,
            algorithms=["RS256"],
            audience=COGNITO_CLIENT_ID,
            issuer=issuer
        )
        return True
    except Exception as e:
        print(f"Token validation error: {e}")
        return False


# =========== ROUTE HELPERS =============
def serve_callback_path(params):
    # Parse query params
    code = params.get("code")
    if not code:
        print("Missing code parameter in callback")
        return {
            "statusCode": 400,
            "body": json.dumps({"error": "Missing code parameter"})
        }

    # Exchange code for tokens
    tokens = exchange_code_for_tokens(code)
    id_token = tokens.get("id_token")
    if not id_token:
        print("No id_token returned from token exchange")
        return {
            "statusCode": 500,
            "body": json.dumps({"error": "No id_token returned"})
        }

    # Set HttpOnly cookie and redirect back to frontend
    print("Token exchange successful, setting cookie and redirecting")
    return {
        "statusCode": 302,
        "headers": {
            "Set-Cookie": set_cookie_header("session", id_token),
            "Location": FRONTEND_URL
        },
        "body": ""
    }


def serve_authcheck_path(headers):
    # Check if session cookie exists and is valid
    cookie_header = headers.get("Cookie", "")
    cookies = parse_cookies(cookie_header)
    id_token = cookies.get("session")
    valid_token = validate_token(id_token)
    if valid_token:
        print("User authenticated")
        return {
            "statusCode": 200,
            "body": json.dumps({"authenticated": True})
        }
    print("User not authenticated")
    return {
        "statusCode": 401,
        "body": json.dumps({"authenticated": False, "error": "Unauthorized"})
    }


def serve_search_path(headers, body):
    # Read token from cookie
    cookie_header = headers.get("Cookie", "")
    cookies = parse_cookies(cookie_header)
    id_token = cookies.get("session")
    valid_token = validate_token(id_token)
    if not valid_token:
        print("Unauthorized search attempt")
        return {
            "statusCode": 401,
            "body": json.dumps({"error": "Unauthorized"})
        }

    # Parse request body
    try:
        params = json.loads(body)
    except json.JSONDecodeError:
        return {"statusCode": 400, "body": json.dumps({"error": "Invalid JSON"})}
    query = params.get("query", "")
    if not query:
        return {"statusCode": 200, "body": json.dumps({})}

    # Execute credential client
    print(f"Executing search for query: {query}")
    result = subprocess.run(
        [
            "./hacc-linux-amd64",
            "search",
            query,
            "--json-output"
        ],
        capture_output=True,
        text=True,
        check=True,
        timeout=10
    )
    return {
        "statusCode": 200,
        "headers": {"Content-Type": "application/json"},
        "body": result.stdout
    }


# ========== LAMBDA HANDLER ===========
def lambda_handler(event, _):
    print("Lambda invoked", json.dumps(event))
    try:
        path    = event["requestContext"]["path"]
        method  = event["requestContext"]["httpMethod"]
        headers = event.get("headers", {})
        params  = event.get("queryStringParameters") or {}
        body    = event.get("body", "{}")

        ## Route based on path and method
        if path == "/callback" and method == "GET":
            return serve_callback_path(params)
        elif path == "/search" and method == "POST":
            return serve_search_path(headers, body)
        elif path == "/auth-check" and method == "GET":
            return serve_authcheck_path(headers)
        else:
            return {
                "statusCode": 404,
                "body": json.dumps({"error": "Not found"})
            }

    except subprocess.CalledProcessError as e:
        return {
            "statusCode": 500,
            "body": json.dumps({
                "error": str(e),
                "stdout": e.stdout,
                "stderr": e.stderr
            })
        }
    except Exception as e:
        if "Token exchange failed" in str(e):
            return {
                "statusCode": 502,
                "body": json.dumps({"error": "Token exchange failed", "details": str(e)})
            }
        return {
            "statusCode": 500,
            "body": json.dumps({"error": str(e)})
        }
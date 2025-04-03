import json
import requests

# Define your webhook URL
WEBHOOK_URL = "https://your-webhook-url.com/lookup"

# Define your routing destinations
ROUTES = {
    "APPROVED": "https://approved-endpoint.com",
    "DENIED": "https://denied-endpoint.com",
    "PENDING": "https://pending-endpoint.com"
}

def lambda_handler(event, context):
    try:
        # Parse incoming request
        body = json.loads(event.get("body", "{}"))
        lookup_value = body.get("lookup_value")

        if not lookup_value:
            return {
                "statusCode": 400,
                "body": json.dumps({"error": "Missing lookup_value"})
            }

        # Call webhook for lookup
        response = requests.post(WEBHOOK_URL, json={"lookup_value": lookup_value}, timeout=5)
        
        if response.status_code != 200:
            return {
                "statusCode": 502,
                "body": json.dumps({"error": "Failed to retrieve lookup result"})
            }

        # Extract webhook response
        result = response.json().get("status")

        # Route based on result
        if result in ROUTES:
            return {
                "statusCode": 302,
                "headers": {
                    "Location": ROUTES[result]
                },
                "body": json.dumps({"message": f"Redirecting to {ROUTES[result]}"})
            }
        else:
            return {
                "statusCode": 400,
                "body": json.dumps({"error": "Invalid response from webhook"})
            }

    except Exception as e:
        return {
            "statusCode": 500,
            "body": json.dumps({"error": str(e)})
        }

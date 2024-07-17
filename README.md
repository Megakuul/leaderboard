# Leaderboard

Simple leaderboard system constructed on aws lowcost serverless infrastructure.

### Deploy the system

Prerequisites:
- `aws cli`
- `aws-sam cli`
- `nodejs`

First you need to build and deploy the whole system to aws.

For this first build the backend lambda functions:
```bash
sam build
```

Then deploy them to aws:
```bash
sam deploy --parameters-overrides CognitoDomainPrefix=example LeaderboardDomain=example.com
```

This initiates the CloudFormation stack. If the specified domain isn't hosted on your AWS account within Route53, manually add the dns validation entries using the cname provided on the `ACM` dashboard.

Note: The stack remains in `CREATE_IN_PROGRESS` state until certificate creation completes (-> has been validated).

In the next step you need to deploy the web application:

1. Build the sveltekit static files (it uses static-adapter to generate prerendered html pages):
   ```bash
   cd web
   npm run build
   ```

2. Upload to the S3 bucket (bucket name is provided in SAM output):
   ```bash
   aws s3 cp build s3://<bucket-name> --recursive
   ```


Finally add a dns entry pointing the application's domain to the CloudFront address (found in the SAM output).


### Remove the system



### Authentication

Authentication is managed through Cognito using the OAuth2 implicit flow. To access protected api endpoints, utilize the cognito access token (in authorization bearer) acquired via implicit flow.

**Important**: Never ever ever use implicit code flow in production applications containing user data that should be protected.
The implicit OAuth2 flow is not considered very secure for various reasons, use the code flow instead as implemented in [battleshiper](https://github.com/megakuul/battleshiper). This application does only use the implicit flow for simplicity and because I was too lazy to implement the code flow backend.


### API
---

The leaderboard api runs on top of lambda functions behind an api-gateway. For simplicity every route uses its own lambda function.


```GET /api/fetch```

**Params**: 
  - **username**: fetches entries queried by the provided username. only applies if previous pagination params are unset.
  - **elo**: fetches entries queried by the provided elo. only applies if previous pagination params are unset.
  - **lastpagekey**: fetches the next page of sorted entries using a base64-encoded json "LastEvaluatedKey" from dynamodb. only applies if previous pagination params are unset.

  - **noparam**: if no parameter is set, fetches the first page sorted desc by elo.

**Returns**:

  - **200**: application/json
    ```json
    {
      "message": "success message xy",
      "newpagekey": "BASE64ENCODEDLASTPAGEKEYORNULLIFNOTQUERIEDBYPAGE",
      "users": [
        {
          "username": "Wendelin Knack",
          "title": "Wendig",
          "iconurl": "https://urltoicon",
          "elo": 420
        }
      ]
    }
    ```
  - **400-500**: text/plain
    ```
    errormessage as plaintext
    ```


```POST /api/update```

**Headers**:
  - **Authorization**: "Bearer id_token"

**Returns**:

  - **200**: application/json
    ```json
    {
      "message": "success message xy",
    }
    ```
  - **401**: text/plain
    Provided id_token has expired or is invalid (catched by the API gateway).
    ```
    errormessage as plaintext
    ```
  - **400-500**: text/plain
    ```
    errormessage as plaintext
    ```


```POST /api/addgame```

**Headers**:
  - **Authorization**: "Bearer id_token"

**Body**:
  - ```json
    {
      "results": [
        {
          "username": "Kater Karlo",
          "placement": 1
        },
        {
          "username": "Panzerknacker",
          "placement": 2
        }
      ]
    }
    ```

**Returns**:

  - **200**: application/json
    ```json
    {
      "message": "success message xy",
    }
    ```
  - **401**: text/plain
    Provided id_token has expired or is invalid (catched by the API gateway).
    ```
    errormessage as plaintext
    ```
  - **400-500**: text/plain
    ```
    errormessage as plaintext
    ```
# Leaderboard

![leaderboard favicon](/web/public/favicon.svg "leaderboard favicon")

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

Then we need to generate the `ACM` certificate for the application (important, the certificate must be in `us-east-1` regardless of your region). This can be done via the AWS ACM dashboard or via cli:
```bash
export CERT_ARN=$(aws acm request-certificate --region us-east-1 --domain-name example.com --validation-method DNS --query 'CertificateArn' --output text)
aws acm describe-certificate --region us-east-1 --certificate-arn $CERT_ARN --query 'Certificate.DomainValidationOptions[0].ResourceRecord'
```
(the values provided by the ResourceRecord output must be entered in your dns console to validate the domain)


After validation of the certificate, deploy the system to aws and provide the cert-arn obtained in the previous step as parameter:
```bash
sam deploy --parameter-overrides CognitoDomainPrefix=example LeaderboardDomain=example.com LeaderboardDomainCertificateArn=$CERT_ARN
```

In the next step you need to deploy the web application:

1. Build the svelte application and specify the build environment variables like described in `.env.example`. You can find the required values for cognito in the `SAM` output.
   ```bash
   cd web
   VITE_COGNITO_DOMAIN=<CognitoEndpoint> VITE_COGNITO_CLIENT_ID=<CognitoClientID> VITE_LEADERBOARD_REGIONS=<DeploymentRegion> npm run build
   ```

2. Upload to the S3 bucket (bucket name is provided in SAM output):
   ```bash
   aws s3 cp dist s3://<bucket-name> --recursive
   ```


Then add a dns entry pointing the application's domain to the CloudFront address (found in the SAM output as <CNameEntries>).


In the next step, we need to configure the dns to support our SES (Simple Email Service) configuration. For this, you will need to set some dns entries on the application's domain:

1. The `dmarc` entry is used to specify the behavior for `spf` and `dkim` failures. The recommended strict configuration is provided in the SAM output as `SPFEntries`. 
2. The `spf` entry is used to specify which server is allowed to send mail from this domain. Here we set `include:amazonses.com` so that mail servers use the SPF configuration of SES (containing the respective ip's). The recommended strict configuration is provided in the SAM output as `SPFEntries`.
3. The `dkim` records are required to verify the mail signature. `dkim` needs three different entries and is also serving the purpose of validating the applications' domain on SES. The required configuration is provided in the SAM output as `DKIMEntries`.


Finally we must request production access for AWS SES. As mentioned [here](https://docs.aws.amazon.com/ses/latest/dg/request-production-access.html), AWS SES sandboxes the environment by default, which means that all recipients MUST be manually added and verified in the SES portal.

You can follow the instructions [here](https://docs.aws.amazon.com/ses/latest/dg/setting-up.html#quick-start-verify-email-addresses) or follow the wizard in the SES portal to request production access for the leaderboard domain.

After the production access has been granted, the application should be fully functional.


### Remove the system

You can delete the system from your aws account with:
```bash
sam delete --stack-name leaderboard
```

The certificate can be revoked with:
```bash
aws acm delete-certificate --region us-east-1 --certificate-arn $CERT_ARN
```


### Authentication

Authentication is managed through Cognito using the OAuth2 implicit flow. To access protected api endpoints, utilize the cognito access token (in authorization bearer) acquired via implicit flow.

**Important**: Never ever ever use implicit code flow in production applications containing user data that should be protected.
The implicit OAuth2 flow is not considered very secure for various reasons, use the code flow instead as implemented in [battleshiper](https://github.com/megakuul/battleshiper). This application does only use the implicit flow for simplicity and because I was too lazy to implement the code flow backend.


### Future outlook

In the future, this leaderboard system could be extended to work multi-regional. This would require the following adjustments:

- Update dynamodb tables to global tables and deploy them in multiple regions.
- Extract api gateway and lambda code into separate templates to make them deployable in multiple regions.
- Create s3 bucket replication policies to replicate the frontend to multiple regions.
- Create global accelerator between cloudfront and the api gateways / s3 buckets.

(this setup will, however, be more expensive to operate on a small scale)


### API
---

The leaderboard api runs on top of lambda functions behind an api-gateway. For simplicity every route uses its own lambda function.



```GET /api/user/fetch```
Fetches users from the leaderboard.

**Params**:
  - **username**: fetches entries queried by the provided username over all regions. 
  - **elo**: fetches one page of entries starting on the provided elo. only applies if username is not set.
  - **lastpagekey**: fetches the next page of sorted entries (sorted by elo) by region using a base64-encoded json "LastEvaluatedKey" from dynamodb. defaults to "" which returns the first page.
  - **region**: specifies the region from where to fetch the entries. defaults to the region where the called function operates in.
  - **pagesize**: specifies the size of the page for pagination requests. defaults to the maximum page size.

**Returns**:

  - **200**: application/json
    ```json
    {
      "message": "success message xy",
      "newpagekey": "BASE64ENCODEDLASTPAGEKEYORNULLIFNOTQUERIEDBYPAGE",
      "users": [
        {
          "username": "Wendelin Knack",
          "region": "eu-central-1",
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



```POST /api/user/update```
Updates the leaderboard user based on the data from the identity-provider (cognito).
The region is updated based on the aws region of the called function.

**Headers**:
  - **Authorization**: "Bearer id_token"

**Body**:
  - ```json
    {
      "user_updates": {
        "title": "Wendig",
        "iconurl": "https://urltoicon"
      }
    }
    ```

**Returns**:

  - **200**: application/json
    ```json
    {
      "message": "success message xy",
      "updated_user": {
        "username": "Wendelin Knack",
        "region": "eu-central-1",
        "title": "Wendig",
        "email": "wendelin@panzerknacker.org",
        "iconurl": "https://urltoicon",
        "elo": 420
      }
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



```GET /api/game/fetch```
Fetches played games.

**Params**: 
  - **gameid**: fetches games based on the gameid.
  - **date**: fetches games based on the date. only applies if previous params are unset.

**Returns**:

  - **200**: application/json
    ```json
    {
      "message": "success message xy",
      "games": [
        {
          "gameid": "550e8400-e29b-11d4-a716-446655440000",
          "date": "2006-01-02",
          "expires_in": 1721550651,
          "readonly": true,
          "participants": {
            "Panzerknacker": {
              "username": "Panzerknacker",
              "underdog": false,
              "team": 2,
              "placement": 2,
              "points": 130,
              "elo": 250,
              "elo_update": -10,
              "confirmed": true
            },
            "Kater Karlo": {
              "username": "Kater Karlo",
              "underdog": true,
              "team": 1,
              "placement": 1,
              "points": 160,
              "elo": 200,
              "elo_update": 20,
              "confirmed": true
            },
          }
        }
      ]
    }
    ```
  - **400-500**: text/plain
    ```
    errormessage as plaintext
    ```



```POST /api/game/add```
Adds a game to the leaderboard.

**Headers**:
  - **Authorization**: "Bearer id_token"

**Body**:
  - ```json
    {
      "placement_points": 100,
      "participants": [
        {
          "username": "Kater Karlo",
          "team": 1,
          "placement": 1,
          "points": 160
        },
        {
          "username": "Panzerknacker",
          "team": 2,
          "placement": 2,
          "points": 130
        }
      ]
    }
    ```

**Returns**:

  - **200**: application/json
    ```json
    {
      "message": "success message xy",
      "gameid": "game-uuid"
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



```GET /api/game/confirm```
Lets a user confirm the specified game. If all users confirmed the game, this will also finish the game and distribute the elo to all players.

**Params**: 
  - **gameid**: specifies the game by id. parameter is required.
  - **username**: identifies the user to confirm by username. parameter is required.
  - **code**: specifies the confirm secret that authorizes the user to confirm. parameter is required.

**Returns**:

  - **200**: text/plain
    ```json
    successmessage as plaintext
    ```
  - **400-500**: text/plain
    ```
    errormessage as plaintext
    ```
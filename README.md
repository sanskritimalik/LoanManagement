# Description

This is an example of API server written using [Echo](https://echo.labstack.com/) framework.
It uses [Auth0](https://auth0.com/) service to authenticate users and [mongodb](https://www.mongodb.com/) to store associated account resources.

You can also read a step-by-step guide for implementation of this API: [link](https://dev.to/mrsydar/how-to-implement-a-basic-crud-in-golang-protected-by-auth0-7je).

## How it works

### Authentication

When a user lands on Auth0 login page associated with your application, 
the service authenticates him with entered credentials or in some other way (e.g. Google account).

After that, the user receives a special code and gets redirected to API server http://localhost:9000/callback endpoint.

The callback endpoint reads the mentioned code, sends forms a request with received code and some vulnerable data like application secret,
and sends it to Auth0.

Finally, if provided data is valid, the service sends a response with JWT token which should be returned to the user client as a callback response.

From this point, the user can use the received JWT token to pass API server authentication and get access to protected resources.

### Authorization

The idea is to have an account resource associated with each authenticated user, so we can create other resources which belong to this user.

After standard JWT verification (HS256 signing algorithm), the custom middleware is used to associate user with existing account resource
stored in the database by an email value provided in the token. If there's no account resource that could be associated, the new one is created.
Finally, an account identifier is put into the request context, so API server could use this value to allow, restrict and filter resources,
which belong to the user.

For this example, there is a `notes` resource which could belong only to one user. 

## How to run

1. Register a new application on the auth0 platform, set signing algorithm as HS256 and add http://localhost:9000/callback value to Allowed Callback URLs text field.
2. Create a .env file in the root folder and fill it with all the required variables that could be found in the ./configs/constants/envnames/envnames.go file or use the template below:
```
DATABASE_NAME='<database-name>'
# e.g. mongodb+srv://user:password@examplecluster.12r45rg.mongodb.net/?retryWrites=true&w=majority
MONGO_URI='<mongo-uri>'
ACCOUNT_COLLECTION='accounts'

# e.g. 'dev-example.us.auth0.com'
AUTH0_DOMAIN='dev-tslb5vli.us.auth0.com'
AUTH0_CALLBACK_URL='http://localhost:9000/callback'
AUTH0_CALLBACK_ENDPOINT='/callback'
# e.g. 'EXAMjthvduhabfcmABCRUduydFePLE'
AUTH0_CLIENT_ID='<auth0-client-id>'
# e.g. 'Gnal9SLkijfxu0lkoif-u8i2vSclkfjsdkfsdJKJL2HofdsOKdkjfflkdflgk'
AUTH0_CLIENT_SECRET='<auth0-client-secret>'
```
3. Enjoy


https://dev-2groccvvlzjrpj7p.us.auth0.com/authorize?response_type=code&client_id=IYZJqE1fFHSKTvvmlyNPsQMim6Y552kV&redirect_uri=http://localhost:8080/callback&scope=openid%20email&state=STATE

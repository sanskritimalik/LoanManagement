# Description

This is an example of API server written using [Echo](https://echo.labstack.com/) framework.
It uses [Auth0](https://auth0.com/) service to authenticate users and [mongodb](https://www.mongodb.com/) to store associated account resources.

# Loan Management System API
About
This is a API first approach application built to manage a Loan Management System. These API's can be consumed by any front end tech stack be it in react, angular or offcourse any mobile app. Currently the Application is built using using [Echo](https://echo.labstack.com/) framework.
It uses [Auth0](https://auth0.com/) service to authenticate users and [mongodb](https://www.mongodb.com/) to store associated account resources.

It is an app that allows authenticated users to go through a loan application. It doesn’t have to contain too many fields, but at least “amount required” and “loan term.” All the loans will be assumed to have a “weekly” repayment frequency. After the loan is approved, the user must be able to submit the weekly loan repayments. It can be a simplified repay functionality, which won’t need to check if the dates are correct but will just set the weekly amount to be repaid.

- Choices I made for the application
    - Used OAuth to integrate Authentication and Authorisation using Jwt token
    - Golang for implementation as it has a simple easy to understand syntax
    - MongoDb to store data to allow flexibility to change the structure of data with minimal changes

- Created various segregation layers within the implemenation to keep the logic scalable and moduler
    - used a configs package to initialize different configs e.g. database, environment, logger etc
    - used a controllers package with various segregations specific to different usecase
    - models package to initialize various entities used across implementation
    - pkg to keep all the generic common functionalities required application wide
    - routes to initiate and manage all the different routes from a single place


      https://dev-2groccvvlzjrpj7p.us.auth0.com/authorize?response_type=code&client_id=IYZJqE1fFHSKTvvmlyNPsQMim6Y552kV&redirect_uri=http://localhost:8080/callback&scope=openid%20email&state=STATE
    - sign up and allow access. This further leads to a screen which dispalys the bearer token
    - This bearer token can be used to further access endpoints  
    - accept the request from the route
    - send it to the controller to process
    - get back the response and render to user


Github repo
https://github.com/sanskritimalik/LoanManagement

Features
- Customer signs up via the Oauth Portal and uses token to access any APIs associated with application
- Customer creates a loan
- Admin approves the loan
- Customer can only view self owned loan
- Customer can repay the loan only once Admin approves the loan
- Once customer pays all the scheduled payment the Loan is marked automatically marked as Paid
- Customer create a loan: Customer submit a loan request defining amount and term example:
- Request amount of 10.000 $ with term 3 on date 7th Feb 2022
customer will generate 3 scheduled repayments:
14th Feb 2022 with amount 3.333,33 $
21st Feb 2022 with amount 3.333,33 $
28th Feb 2022 with amount 3.333,34 $
the loan and scheduled repayments will have state PENDING
- Admin approve the loan:
- Admin change the pending loans to state APPROVED
- Customer can view loan belong to him:
- Add a policy check to make sure that the customers can view them own loan only.
- Customer add a repayments:
- Customer add a repayment with amount greater or equal to the scheduled repayment
- The scheduled repayment change the status to PAID
- If all the scheduled repayments connected to a loan are PAID automatically also the loan become PAID

Packages used
Oauth Service - https://auth0.com/
MongoDB - https://www.mongodb.com/

Auth Routes

CreateUser - http://localhost:8080/user  (Post)
GetUsers - http://localhost:8080/user (Get)

Loan Routes

Create Loan - http://localhost:8080/loans
View Loan - http://localhost:8080/loans
Approve Loan - http://localhost:8080/approve/<loan_uuid>
Repay Loan - http://localhost:8080/repayments/<loan_uuid>

## How it works

### Authentication

When a user lands on Auth0 login page associated with your application, 
the service authenticates him with entered credentials or in some other way (e.g. Google account).

After that, the user receives a special code and gets redirected to API server http://localhost:8080/callback endpoint.

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

For this example, there is a `loans` resource which could belong only to one user. 

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
3. Get the application files
download the zip or clone it from github - https://github.com/sanskritimalik/LoanManagement
put it at <your-sites-or-htdocs-folder-path> from where you can load the application
4. Create a DB on MongoDB with name "LoanManagement"
Add Collections
5. Once the above steps are completed, use the URL in below format to signup and hence generate the bearer token.

https://<AUTH0_DOMAIN>/authorize?response_type=code&client_id=<AUTH0_CLIENT_ID>&redirect_uri=http://localhost:9000/callback&scope=openid%20email&state=STATE

6. Use the route endpoints on Postman to perform various operations

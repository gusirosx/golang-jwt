# golang-jwt

// To do //

// Handle GET requests at /users/claims
// Handle POST requests at /users/claims
// Handle PUT requests at /users/claims
// Handle DELETE requests at /users/claims




Please feel free to create a new issue if you come across one or want a new feature to be added. I am looking for contributors, feel free to send pull requests.

##What is Gondalf?##
Gondalf is a ready to deploy microservice that provides user management, authentication, and role based authorization features out of box. Gondalf is built using [martini] and [gorm], and uses [postgresql] as the default database.

##Features:##

###1. User management###
- User creation
- Validating unique username
- Password change on first login
- Encrypted password storage
- Activity logs

###2. Authentication###
- User authentication
- Token-based authentication
- Custom token expiry and renewal times

###3. Authorization###
- Role-based authorization including group permissions


##Why Gondalf?##
Over the course of multiple projects I realized that there are some common features that can be packed into a single microservice and can be used right out of the box. Gondalf is the first piece in that set.


##TODO List##
- [X] Add end points for permission checking
- [X] Input timeout values and server port from config file
- [X] Refresh app properties from DB after fixed interval
- [X] Add a cron job for cleaning up and archiving expired session tokens to keep the validation request latency low
- [X] Dockerize gondalf
- [X] Refactored the API to include a consistent error payload
- [ ] Add more events to Activity Logs
- [ ] Improve documentation - add details about logging, app properties, and configuration 
- [ ] Switch to [negroni](https://github.com/codegangsta/negroni) and [gorilla mux](http://www.gorillatoolkit.org/pkg/mux)
- [ ] Add TLS support for the end point
- [ ] Provide one click deploy solution
- [ ] Add CI on checkins
- [ ] Add support for other databases


###Why call it *Gondalf* ?###

Because  <img src="http://www.reactiongifs.com/wp-content/uploads/2013/12/shall-not-pass.gif" width="150px" height="75px"/> and it is Go, so why not both?


 
##Installation Instructions:##

- Clone the repository in a local directory

- Database configuration can be set under the the config file named gondalf.config

- Gondalf creates required tables using the configuration provided in the config file. For this the 
initdb flag should be set to true when starting the app.

`$ bash startApp.sh -initdb=true` 



##Request and Response formats##

###Error Codes###

- Invalid Session Token
- Expired Session Token
- Unregistered User
- Invalid Password
- First Login Change Password
- Authentication Failed
- Encryption Failed
- Database Error
- Permission Denied
- System Error
- Duplicate Username Error

###LoginCredential###

####Request####

```javascript
{
  "username": "test2User",
  "password" : "testPassword",
  "deviceId" : 1
}
```

deviceId code 1 for web, 2 for mobile

####Response####

```javascript
{
  "sessionToken": "testSessionToken",
}
```

###ValidateUsername###

For validating unique username

####Request####

```javascript
{
	"username": "test2User"
}
```

####Response####

```javascript
{
	"valid": true
}
```

###CreateUser###

####Request####

```javascript
{
	"username": "test2User",
	"legalname": "testLegalName",
	"password": "testPassword"
}
```
####Response####

```javascript
{
	"userCreated": true
}
```

###Change password###

####Request####

```javascript
{
	"username": "test2User",
	"oldPassword": "testOldPassword",
	"newPassword": "newTestPassword",
	"deviceId": 1
}
```

####Response####

If the old credentials are correct then:

```javascript
{
	"passwordChanged": true,
	"sessionToken": "testSessionToken"
}
```

###Validate Session Token###

####Request####

```javascript
{
	"sessionToken": "testSessionToken"
}
```

####Response####

```javascript
{
	"userId": 1234
}
```

###Permission Checking###

####Request####

```javascript
{
  "userId": 123456,
  "permissionDescription" : "ADMIN"
}
````

####Response####

```javascript
{
  "permissionResult" : true
}
```

###Error Response###

```javascript
{
	"status": "Internal Server Error" / "Unauthorized" / "Conflict" / "Forbidden",
	"message": "Invalid Session Token" / "Expired Session Token" etc.,
	"description": ""
}
```

======================================================================================
<html>
    <p align="center" >
        <img width="40%" src="./static/images/lugbit_logo.png"/>
    </p>
</html>

## Golang authentication and sessions
A custom registration, authentication and session handling implemented in Golang and MySQL.

## Motivation
This project was started purely as a learning tool to aid myself with learning Golang and specifically web development with Golang. I wanted to incorporate Go with a RDBMS such as MySQL as well as implement other features such as account verification upon sign up and sessions without the use of a framework.
 
## Features
* **Registration and activation**
    - Activation link is sent to the user's registration email with a UUID token embedded upon signing up. This one time use link must be used before it is expired or it will not activate the associated account.
    - Activation links are one time use. Once clicked, it will be marked as used in the database and cannot be used again.
    - Activaton links can be resent by visiting the /send-activation route which will generate a new activation link with a new expiry date.
* **Authentication**
    - Registered users with activated accounts can login with their email and password.
    - The email address the user enters is checked against the database. If the email exists, the password entered is hashed and compared against the hashed password entry in the database and if the hash matches, the user is logged in and a session is created.
* **Sessions**
    - Once a user successfully authenticates, a session is created.
    - The UUID generated upon logging in is inserted in the sessions table with the user's unique id.
    - A cookie is sent to the client with the same UUID that was inserted in the database.
    - When a user makes a request for a secured route e.g. /my-profile, the session cookie is received by the server and verified against the sessions table. If the same UUID in the cookie is also present in the sessions table, the user is granted access to the secured route.
    - Sessions and cookies have a maximum life before they expire, once expired, the user will need to login again.
    - Sessions and cookies are automatically renewed when the user makes a request to any secured routes. This will reset the max life of the session and cookie.
    - Session entries along with the session cookie is destroyed upon logging out.
* **Input validation**
    - The registration, login and any routes with input forms are validated to make sure they are not empty, have the correct format or unique if a user is entering their email address on sign up.

## Setup
* Load the MySQL database schema located at **./static/db/userAuthDBSchema.sql**
* Set the app environment variables located at **./.env**
    - **DB_USERNAME =** _Your MySQL server username_
    - **DB_PASSWORD =** _Your MySQL server password_
    - **DB_ADDRESS =** _Your MySQL server address_
    - **DB_PORT =** _Your MySQL server port_ 
    - **DB_NAME =** _The database name, the default is "userAuthDB" as set on the schema._
    - **SENDER_ADDRESS =** _Verification link sender email address_
    - **SENDER_PASSWORD =** _Verification link sender email address password_
    - **SMTP_SERVER =** _Verification link sender email address SMTP server_
* Build the binary by running **go build** and then run the executable
```
go build

```
* Visit **localhost:8080** on your web browser

## Credits
Thanks to [Todd McLeod](https://github.com/GoesToEleven) and particularly his course on Udemy([Course Repo](https://github.com/GoesToEleven/golang-web-dev)).

## License
MIT License

Copyright (c) 2018 Marck Muñoz
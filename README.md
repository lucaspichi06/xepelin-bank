# Xepelin Bank ![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

This API is responsible for provide three endpoints:

- Account creation: this endpoint is in charge of creating new accounts
- Account balance: this endpoint will return the balance for a provided account id
- Transactions: this endpoint is in charge to process the different transactions (deposit/withdraw/transfer)

## Index
- [Run it locally](#run-it-locally)
- [Api Docs](#api-docs)
- [Test the application](#test-the-application)
- [How I would implement the solution in AWS](#solution-approach)

## Run it locally
To run the api locally it is necessary to clone this repository from Github:
```` bash
git clone https://github.com/lucaspichi06/xepelin-bank.git
````

After that, you have to move to the root of the project and run the following command:
```` bash
make up
````

It will start the `docker-compose` with all the dependencies you need to run the application.

Now you have the API running in the port :8080, so you can invoke the endpoints with the following curls (you can use Postman too):

- Account Creation
````bash
curl --location 'http://localhost:8080/accounts' \
--header 'token: my-secret-token' \
--header 'Content-Type: application/json' \
--data '{
    "name": "my-new-account"
}'
`````

- Account Balance
````bash
curl --location 'http://localhost:8080/accounts/ACC_ID/balance'
`````
_Note: replace `ACC_ID` with a valid value_

- Transactions Process

````bash
curl --location 'http://localhost:8080/transactions' \
--header 'token: my-secret-token' \
--header 'Content-Type: application/json' \
--data '{
    "account_id": "ACC_ID",
    "type": "deposit|withdraw|transfer",
    "amount": 100.00,
    "destination_id": "DEST_ID"
}'
`````
_Note: replace `ACC_ID` with a valid value (replace `DEST_ID` as well when `"type": "transfer"`)_

_Note: every transaction impacts in the account balance. Also, every transaction event is stored in the table `transactions`_

- Besides that, you have another endpoint to check the API health status. If the API is running successfully, it has to return the word ```pong```:
````bash
curl --location --request GET 'http://localhost:8080/ping'
````



- Transaction Logger: the application implements a logger to print in the stdout every transaction greater than $10000.00.

## Api Docs
The documentation has been done using `Swagger`. You can access to the documentation page here:
````
http://localhost:8080/docs/index.html
````
_Note: the application needs to be running_

Swagger provides an excellent UI to show the documentation and provides the possibility to test the application from there.

## Test the application
This API have unit tests to warranties the integrity of the application.
You can run this tests with the following command from the root of the project:
```` bash
make test
````

## Solution Approach
For an implementation of this application in AWS, I suggest to implement several EC2 instances connected to a Load Balancer (through the Elastic Load Balancing service) to distribute the traffic between the instances.

Each EC2 instance executes the Golang application to handle the requests and interacts with the DB.

For the database, it's possible to use Amazon RDS to manage the MySQL DB.

````
                            |internet gateway|
                                    |
                              |Load Balancer|
                    ________________|___________________
                    |               |                   |
               EC2 Instance     EC2 Instance     EC2 Instance (N)
               (application)    (application)    (application)
                    |               |                   |
                MySQL (RDS)      MySQL (RDS)        MySQL (RDS)

````
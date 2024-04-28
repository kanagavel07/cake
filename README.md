# Golang API Challenge - Waracle

A Golang API hosted on Google Cloud Platform (GCP) that meets the provided
specifications.

The API should allow users to interact with cake data, including listing all cakes,
searching for cakes by yumFactor and/or name, adding a new cake, and deleting an existing
cake.

Specifications:

- The API should return an OpenAPI or Swagger spec.
- Users should be able to fetch a cake by ID.
- Users should be able to list all cakes.
- Users should be able to search for a cake by yumFactor and/or name.
- Users should be able to add another cake.
- Users should be able to delete an existing cake.

## Prerequisites

- Go version 1.22
- A Google Account
- gcloud CLI installed locally
- A Firebase Project created with a Firestore Database, Collection and a Document
- Functions Framework for Go (The Functions Framework lets you write lightweight functions that run in many different environments, including: Google Cloud Functions and Our local development machine)

## Development

- API Entry points: Below HTTP API and its methods are implemented
- /cakes:
  GET Method: To list all cakes
  POST Method: To add another cake
- /cake:
  GET Method: To fetch a cake by ID
  DELETE Method: To delete an existing cake

To test the functions locally, main function with a http server is also added in cmd/main.go
To send http requests through REST client, the sample requests are added in tests/requests.http

## Testing locally

Functions Framework for Go is added to test the APIs locally
Run the below command in the terminal to start the API functions, default port is 8080 if PORT=[portnumber] is not provided
`PORT=[portnumber] LOCAL_ONLY=true go run cmd/main.go`

## Deploying Cloud Functions

Create Cloud Functions by using the Google Cloud CLI

- `gcloud auth login`
- `gcloud functions deploy CakesAPI --gen2 --region=europe-west2 --runtime=go121 --source=. --entry-point=cakes --trigger-http`
- `gcloud functions deploy CakeAPI --gen2 --region=europe-west2 --runtime=go121 --source=. --entry-point=cake --trigger-http`

## Github Project

Github Project Projects is an adaptable, flexible tool for planning and tracking work on GitHub.
https://github.com/users/kanagavel07/projects/1

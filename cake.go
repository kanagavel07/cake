package cake

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

// Type Cakes represents the list of cakes in the http responses
type Cakes []Cake

// Type Cakes defines the structure of the fields of the cake in the database collection
type Cake struct {
	ID        string `firestore:"-"`
	Name      string `firestore:"name,omitempty"`
	Comment   string `firestore:"comment,omitempty"`
	ImageURL  string `firestore:"imageUrl,omitempty"`
	YumFactor int    `firestore:"yumFactor,omitempty"`
}

// Type DeleteBody represents the expected structure of the body of a delete http request
type DeleteBody struct {
	ID string `json:"id"`
}

func init() {
	functions.HTTP("cakes", CakesAPI)
	functions.HTTP("cake", CakeAPI)
}

// CakesAPI is an HTTP Cloud Function to handle /cakes GET and POST methods
func CakesAPI(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	// Firebase configuration
	firebaseConfig := &firebase.Config{
		ProjectID: "dark-form-304415",
	}

	// Create firebase app
	firebaseApp, err := firebase.NewApp(ctx, firebaseConfig)
	if err != nil {
		log.Printf("Error initializing firebase app: %v\n", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create firestore client
	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Printf("Error initializing firestore: %v", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer firestoreClient.Close()

	switch method := request.Method; method {
	case http.MethodGet:
		getCakes(ctx, firestoreClient, responseWriter)
	case http.MethodPost:
		addCake(ctx, firestoreClient, responseWriter, request)
	default:
		http.Error(responseWriter, "Unsupported HTTP method", http.StatusNotFound)
	}
}

// CakeAPI is an HTTP Cloud Function to handle /cake GET and DELETE methods
func CakeAPI(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	// Firebase configuration
	firebaseConfig := &firebase.Config{
		ProjectID: "dark-form-304415",
	}

	// Create firebase app
	firebaseApp, err := firebase.NewApp(ctx, firebaseConfig)
	if err != nil {
		log.Printf("Error initializing firebase app: %v\n", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create firestore client
	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Printf("Error initializing firestore: %v", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer firestoreClient.Close()

	switch method := request.Method; method {
	case http.MethodGet:
		getCake(ctx, firestoreClient, responseWriter, request)
	case http.MethodDelete:
		deleteCake(ctx, firestoreClient, responseWriter, request)
	default:
		http.Error(responseWriter, "Unsupported HTTP method", http.StatusNotFound)
	}
}

func getCakes(ctx context.Context, client *firestore.Client, responseWriter http.ResponseWriter) {

	// Retrieve all the documents in an iterator
	cakesIterator := client.Collection("cakes").Documents(ctx)
	defer cakesIterator.Stop()

	var cakes Cakes

	// Loop through the documents to be added in the http response
	for {
		cakeDocument, err := cakesIterator.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}

		var cake Cake

		if err := cakeDocument.DataTo(&cake); err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Adding Document ID in the document data
		cake.ID = cakeDocument.Ref.ID

		cakes = append(cakes, cake)
	}

	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(cakes)
}

func addCake(ctx context.Context, client *firestore.Client, responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		panic(err)
	}
	defer request.Body.Close()

	var newCake Cake

	err = json.Unmarshal(body, &newCake)

	if err != nil {
		http.Error(responseWriter, "Json unmarshalling failed", http.StatusBadRequest)
		return
	}

	// Adding new document into the firestore collection
	newDocument, _, err := client.Collection("cakes").Add(ctx, &newCake)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Adding Document ID in the document data for the http response
	newCake.ID = newDocument.ID

	responseWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(responseWriter).Encode(newCake)
}

func getCake(ctx context.Context, client *firestore.Client, responseWriter http.ResponseWriter, request *http.Request) {

	// Fetch document by ID
	id := request.URL.Query().Get("id")
	cakeDocument, err := client.Collection("cakes").Doc(id).Get(ctx)

	if err != nil {
		responseWriter.WriteHeader(http.StatusNoContent)
		return
	}

	var cake Cake

	if err := cakeDocument.DataTo(&cake); err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Adding Document ID in the document data for the http response
	cake.ID = cakeDocument.Ref.ID

	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(cake)
}

func deleteCake(ctx context.Context, client *firestore.Client, responseWriter http.ResponseWriter, request *http.Request) {

	body, err := io.ReadAll(request.Body)

	if err != nil {
		panic(err)
	}
	defer request.Body.Close()

	var deleteBody DeleteBody

	err = json.Unmarshal(body, &deleteBody)

	if err != nil {
		http.Error(responseWriter, "Json unmarshalling failed", http.StatusBadRequest)
		return
	}

	// Delete Document by ID
	_, err = client.Collection("cakes").Doc(deleteBody.ID).Delete(ctx)

	if err != nil {
		responseWriter.WriteHeader(http.StatusNoContent)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}

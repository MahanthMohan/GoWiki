package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongo model for the article struct
type article struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title,omitempty"`
	Body   string             `json:"body,omitempty"`
	Image  string             `json:"image,omitempty"`
	Author string             `json:"author,omitempty"`
}

// Global constants (Exception: credential field in databaseURI)
var databaseURI = "mongodb+srv://mahanth:" + getCredential() + "@gowiki.ayckl.mongodb.net/GoWiki?retryWrites=true&w=majority"

const myDatabase = "GoWiki"
const myCollection = "Articles"

// Creates a global variable of type mongo.Collection to be used in functions
var collection *mongo.Collection

func getCredential() string {
	res, err := ioutil.ReadFile("credentials.txt")
	if err != nil {
		log.Fatalln(err)
	}
	return string(res)
}

// A function to initialize database before running CRUD operations
func init() {
	clientOptions := options.Client().ApplyURI(databaseURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("** Database connection successful **")
	// Assigns the value of the collection to var collection
	collection = client.Database(myDatabase).Collection(myCollection)
	log.Println("** Created a collection **")
}

// A limiter that has a ceiling of maximum of 60 requests per minute => 1 s/r
// bucketSize = 3 requests -> queued
var limiter = rate.NewLimiter(rate.Every(time.Second), 3)

// A rate limit function to limit the number of requests per minute
func rateLimit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Each http function has two input parameters a http.ResponseWriter and a Request, which is a struct object
// Create an article
func createArticle(w http.ResponseWriter, r *http.Request) {
	// All the headers for a post request
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var newArticle article
	_ = json.NewDecoder(r.Body).Decode(&newArticle)
	insertResult, err := collection.InsertOne(context.Background(), newArticle)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(insertResult)
	json.NewEncoder(w).Encode(newArticle)
}

// Read all the articles
func getAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mongoCursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatalln(err)
	}
	var articles []primitive.M
	// Next document in the mongo collection
	for mongoCursor.Next(context.Background()) {
		var document bson.M
		mongoCursor.Decode(&document)
		articles = append(articles, document)
	}
	mongoCursor.Close(context.Background())
	// Encode the document to json and write it as a json response
	json.NewEncoder(w).Encode(articles)
}

// Delete an article
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	args := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(args["id"])
	if err != nil {
		log.Fatalln("Error parsing the _id of document")
	}
	deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(deleteResult)
	json.NewEncoder(w).Encode(args["id"])
}

// Router function to handle all http routes using a mux router
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/create", createArticle).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/read", getAllArticles).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/delete/{id}", deleteArticle).Methods("DELETE", "OPTIONS")
	return router
}

func main() {
	r := Router()
	fmt.Println("Server running on PORT 8080")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, rateLimit(r)))
}

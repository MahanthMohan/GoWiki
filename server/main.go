package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongo model for the article struct
type article struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string             `json:"title,omitempty"`
	Body  string             `json:"body,omitempty"`
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
	args := mux.Vars(r)
	deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"title": string(args["title"])})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(deleteResult)
	json.NewEncoder(w).Encode(args["title"])
}

// Router function to handle all http routes using a mux router
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/create", createArticle).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/read", getAllArticles).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/delete/{title}", deleteArticle).Methods("DELETE", "OPTIONS")
	return router
}

func main() {
	r := Router()
	fmt.Println("Server running on PORT 8080")
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

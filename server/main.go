package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type article struct {
	Title  string `json:"title,omitempty"`
	Body   string `json:"body,omitempty"`
	Image  string `json:"image,omitempty"`
	Author string `json:"author,omitempty"`
}

const (
	myDatabase   = "GoWiki"
	myCollection = "Articles"
)

var (
	databaseURI = "mongodb+srv://mahanth:" + os.Getenv("DB_SECRET") +
		"@gowiki.ayckl.mongodb.net/GoWiki?retryWrites=true&w=majority"
	collection *mongo.Collection
)

func init() {
	clientOptions := options.Client().ApplyURI(databaseURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("** Database connection successful **")
	collection = client.Database(myDatabase).Collection(myCollection)
	log.Println("** Created a collection **")
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var newArticle article
	_ = json.NewDecoder(r.Body).Decode(&newArticle)
	insertResult, err := collection.InsertOne(context.Background(), newArticle)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(insertResult)
	w.WriteHeader(200)
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mongoCursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatalln(err)
	}

	var articles []article
	for mongoCursor.Next(context.Background()) {
		var document article
		mongoCursor.Decode(&document)
		articles = append(articles, document)
	}
	mongoCursor.Close(context.Background())
	json.NewEncoder(w).Encode(articles)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	title := strings.Split(r.URL.Path, "/")[3]
	deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"title": title})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(deleteResult)
	w.WriteHeader(200)
}

func main() {
	fmt.Println("Server running on PORT 8080")
	http.HandleFunc("/api/create", createArticle)
	http.HandleFunc("/api/read", getAllArticles)
	http.HandleFunc("/api/delete/", deleteArticle)
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

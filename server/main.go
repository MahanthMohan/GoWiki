package main

import (
	"context"
	"encoding/json"
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

var (
	databaseURI = "mongodb+srv://mahanth:" + os.Getenv("DB_SECRET") +
		"@gowiki.ayckl.mongodb.net/GoWiki?retryWrites=true&w=majority"
	collection *mongo.Collection
)

func init() {
	clientOptions := options.Client().ApplyURI(databaseURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	println("** Database connection successful **")
	collection = client.Database("GoWiki").Collection("Articles")
	println("** Created a collection **")
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var newArticle article
	_ = json.NewDecoder(r.Body).Decode(&newArticle)
	_, err := collection.InsertOne(context.Background(), newArticle)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	mongoCursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err := collection.DeleteOne(context.Background(), bson.M{"title": strings.Split(r.URL.Path, "/")[3]})
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
}

func main() {
	println("Server running on PORT 8080")
	http.HandleFunc("/api/read", getAllArticles)
	http.HandleFunc("/api/create", createArticle)
	http.HandleFunc("/api/delete/", deleteArticle)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	panic(http.ListenAndServe(":"+port, nil))
}

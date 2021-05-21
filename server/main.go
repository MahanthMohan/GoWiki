package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
	handle(err)
	println("** Database connection successful **")
	collection = client.Database("GoWiki").Collection("Articles")
	println("** Created a collection **")
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var newArticle article
	body, err := ioutil.ReadAll(r.Body)
	handle(err)
	err = json.Unmarshal(body, &newArticle)
	handle(err)
	_, err = collection.InsertOne(context.Background(), newArticle)
	handle(err)
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	mongoCursor, err := collection.Find(context.Background(), bson.M{})
	handle(err)

	var articles []article
	for mongoCursor.Next(context.Background()) {
		var document article
		err := mongoCursor.Decode(&document)
		handle(err)
		articles = append(articles, document)
	}
	err = mongoCursor.Close(context.Background())
	handle(err)
	err = json.NewEncoder(w).Encode(articles)
	handle(err)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	_, err := collection.DeleteOne(context.Background(), bson.M{"title": strings.Split(r.URL.Path, "/")[3]})
	handle(err)
}

func main() {
	println("Server running on PORT 8080")
	http.HandleFunc("/api/read", getAllArticles)
	http.HandleFunc("/api/create", createArticle)
	http.HandleFunc("/api/delete/", deleteArticle)
	port := os.Getenv("PORT")
	panic(http.ListenAndServe(":"+port, nil))
}

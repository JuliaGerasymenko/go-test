package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"context"
	"encoding/json"
	"go-test/handlers"
	"go-test/models"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
)

const (
	CONNECTIONSTRING = "mongodb://localhost:27017"
	DBNAME           = "userInfo"
	COLLECTION       = "peopleinfo"
)

type File struct {
	Objects []Person `json:"objects"` //, bson: "objects, omitempty"`
}

type Person struct {
	Email      string `json: "email,omitempty", bson: "email,omitempty"`
	Last_name  string `json: "last_name,omitempty", bson: "last_name,omitempty"`
	Country    string `json: "country,omitempty", bson:"country,omitempty" `
	City       string `json: "city,omitempty", bson: "city,omitempty"`
	Gender     string `json:"gender,omitempty", bson: "gender,omitempty"`
	Birth_date string `json: "birth_date,omitempty"  bson: "birth_date,omitempty"`
}

//func NewPerson(birthDate string) *Person {
//	return &Person{BirthDate: birthDate}
//}

//var client *mongo.Client
//
//func ConnectDB(ctx context.Context) *mongo.Collection {
//
//	// Set client options
//	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
//
//	// Connect to MongoDB
//	client, err := mongo.Connect(ctx, clientOptions)
//
//	if err != nil {
//		log.Fatal("1",err)
//	}
//
//	fmt.Println("Connected to MongoDB!")
//
//	collection := client.Database("DBNAME").Collection("COLLECTIONNAME")
//
//	return collection
//}

func init() {
	var file File

	client, err := mongo.NewClient(options.Client().ApplyURI(CONNECTIONSTRING))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database(DBNAME).Collection(COLLECTION)

	data, err := ioutil.ReadFile("../../users_go1.json")

	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(data, &file); err != nil {
		log.Fatal(err)
	}

	var people []interface{}

	for _, person := range file.Objects {

		people = append(people, person)
	}

	_, err = collection.InsertMany(context.Background(), people)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", UpdatePersonEndpoint).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}

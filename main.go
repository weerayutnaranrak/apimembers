package main

import (
	"os"
	"fmt"
	"context"
	"encoding/json"
	"log"
	"net/http"
	
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)






func ConnectDB() *mongo.Collection {

	clientOptions := options.Client().ApplyURI("mongodb+srv://appdemo:appdemo@cluster0.sdgiz.mongodb.net/demoapp?retryWrites=true&w=majority")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("go_rest_api").Collection("members")

	return collection
}


type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}


func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}


type Member struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Lastname  string             `json:"lastname" bson:"lastname,omitempty"`
	Age    string         `json:"age" bson:"age,omitempty"`
	Job    string         `json:"job" bson:"job,omitempty"`
	Status    string         `json:"status" bson:"status,omitempty"`
	Address    string         `json:"address" bson:"address,omitempty"`
}


func getMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var members []Member

	collection := ConnectDB()

	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var member Member

		err := cur.Decode(&member) 
		if err != nil {
			log.Fatal(err)
		}

		members = append(members, member)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(members) 
}

func getMember(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var member Member

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := ConnectDB()

	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&member)

	if err != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(member)
}

func createMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var member Member

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&member)

	// connect db
	collection := ConnectDB()

	// insert our member model.
	result, err := collection.InsertOne(context.TODO(), member)

	if err != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	var member Member

	collection := ConnectDB()

	filter := bson.M{"_id": id}

	_ = json.NewDecoder(r.Body).Decode(&member)


	 update := bson.D{
	 	{"$set", bson.D{
			{"name", member.Name},
	 		{"lastname", member.Lastname},
	 		{"age", member.Age},
	 		{"job", member.Job},
	 		{"status", member.Status},
	 		{"address", member.Address},
	 	}},
	 }

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&member)

	if err != nil {
		GetError(err, w)
		return
	}

	member.ID = id

	json.NewEncoder(w).Encode(member)
}

func deleteMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])

	collection := ConnectDB()

	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func getPort() string {
	var port = os.Getenv("PORT") 
	if port == "" {
	   port = "8000"
	   fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port 
}

func main()  {
	r := mux.NewRouter()

	r.HandleFunc("/api/members",getMembers).Methods("GET")
	r.HandleFunc("/api/members/{id}",getMember).Methods("GET")
	r.HandleFunc("/api/members",createMember).Methods("POST")
	r.HandleFunc("/api/members/{id}",updateMember).Methods("PUT")
	r.HandleFunc("/api/members/{id}",deleteMember).Methods("DELETE")

	log.Fatal(http.ListenAndServe(getPort(), r))
}



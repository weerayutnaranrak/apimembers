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

	"./helper"
	"./models"


)

func getMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var members []models.Member

	collection := helper.ConnectDB()

	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var member models.Member

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

	var member models.Member
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := helper.ConnectDB()

	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&member)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(member)
}

func createMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var member models.Member

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&member)

	// connect db
	collection := helper.ConnectDB()

	// insert our member model.
	result, err := collection.InsertOne(context.TODO(), member)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	var member models.Member

	collection := helper.ConnectDB()

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
		helper.GetError(err, w)
		return
	}

	member.ID = id

	json.NewEncoder(w).Encode(member)
}

func deleteMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])

	collection := helper.ConnectDB()

	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
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



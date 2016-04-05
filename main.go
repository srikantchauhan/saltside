package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var Collection *mgo.Collection

// Bird Struct
type Bird struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name"`
	Family     string        `json:"family"`
	Continents []string      `json:"continents"`
	Added      string        `json:"added"`
	Visible    bool          `json:"visible"`
}

func main() {
	Init()
	router := mux.NewRouter()
	router.HandleFunc("/birds", handleBirds).Methods("GET")
	router.HandleFunc("/birds", createBird).Methods("POST")
	router.HandleFunc("/birds/{id}", handleBird).Methods("GET", "DELETE")
	http.ListenAndServe(":8080", router)
}

//init mongo connection
func Init() {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	Collection = session.DB("saltside").C("birds")
}

//get bird
//delete bird
func handleBird(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	birdId := vars["id"]

	switch req.Method {
	case "GET":
		bird := FindById(birdId)
		if bird == nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, string("Bird not found"))
		}
		outgoingJSON, error := json.Marshal(bird)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(res, string(outgoingJSON))
	case "DELETE":
		status := DeleteById(birdId)
		if !status {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}

//get all birds
func handleBirds(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	outgoingJSON, error := json.Marshal(FindAll())
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

//create a bird
func createBird(res http.ResponseWriter, req *http.Request) {
	bird := new(Bird)
	decoder := json.NewDecoder(req.Body)
	error := decoder.Decode(&bird)
	if error != nil || (bird.Name == "" || bird.Family == "" || len(bird.Continents) == 0) {
		if error != nil {
			log.Println(error.Error())
		}
		http.Error(res, "", http.StatusBadRequest)
		return
	}
	Add(bird)
	outgoingJSON, err := json.Marshal(bird)
	if err != nil {
		log.Println(error.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, string(outgoingJSON))
}

//find all documents
func FindAll() []*Bird {
	var result []*Bird
	Collection.Find(nil).All(&result)
	return result
}

//find document by id
func FindById(id string) *Bird {
	var result *Bird
	objectId := bson.ObjectIdHex(id)
	Collection.Find(bson.M{"_id": objectId}).One(&result)
	return result
}

//add new document
func Add(bird *Bird) *Bird {
	Add.Id = bson.NewObjectId()
	err := Collection.Insert(bird)
	if err != nil {
		log.Println(error.Error())
	}
	return bird
}

//delete document by id
func DeleteById(id string) bool {
	var result *Bird
	objectId := bson.ObjectIdHex(id)
	Collection.Find(bson.M{"_id": objectId}).One(&result)
	if result == nil {
		return false
	}
	err := Collection.Remove(bson.M{"_id": objectId})
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

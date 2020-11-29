package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var dbFilename string

// Response represents a response by the server
type Response struct {
	Ok          bool        `json:"ok"`
	Descirption string      `json:"decsription"`
	Result      interface{} `json:"result"`
}

var todos []TodoItem

func save() {
	file, err := os.Create(dbFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log.Printf("Saving to file")
	err = json.NewEncoder(file).Encode(todos)
	if err != nil {
		log.Printf("Error saving DB to file: %v\n", err)
	}
}

func read() {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		return
	}

	file, err := os.Open(dbFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log.Printf("Reading from file")
	err = json.NewDecoder(file).Decode(&todos)
	if err != nil {
		log.Printf("Error reading DB from file: %v\n", err)
	}
}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Retriving all todos")

	response := Response{
		Ok:          true,
		Descirption: "Retrieved succesfully",
		Result:      todos,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error writing json to reponse: %v\n", err)
		return
	}
	log.Printf("Response: %v\n", response)
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoItem := TodoItem{}
	json.NewDecoder(r.Body).Decode(&todoItem)
	id := todoItem.ID
	log.Printf("Removing todo item #%v\n", id)

	item := internalGet(id)
	var response Response
	if item == nil {
		response = Response{
			Ok:          false,
			Descirption: "Not in DB",
		}
	} else {
		response = Response{
			Ok:          true,
			Descirption: "Retrieved succesfully",
			Result:      *item,
		}
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error writing json to reponse: %v\n", err)
		return
	}
	log.Printf("Response: %v\n", response)
}

func internalGet(id int) *TodoItem {
	var item *TodoItem

	for _, ti := range todos {
		if ti.ID == id {
			item = &ti
		}
	}

	return item
}

func add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoItem := TodoItem{}
	json.NewDecoder(r.Body).Decode(&todoItem)
	valid, err := ValidateAllFields(todoItem)
	if err != nil {
		response := Response{
			Ok:          false,
			Descirption: err.Error(),
		}
		log.Printf("Received error: %v while processing the item: %v\n", err, todoItem)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("Error writing json to reponse: %v\n", err)
			return
		}
		log.Printf("Response: %v\n", response)
		return
	}
	if !valid {
		response := Response{
			Ok:          false,
			Descirption: "Invalid request",
		}
		log.Printf("Received invalid data: %v\n", todoItem)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("Error writing json to reponse: %v\n", err)
			return
		}
		log.Printf("Response: %v\n", response)
		return
	}

	duplicateItem := internalGet(todoItem.ID)
	if duplicateItem != nil {
		response := Response{
			Ok:          false,
			Descirption: "Repeated ID",
			Result:      *duplicateItem,
		}
		log.Printf("Tried to add an item with duplicate ID, item with id=%v from DB: %v\n",
			duplicateItem.ID, *duplicateItem)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("Error writing json to reponse: %v\n", err)
			return
		}
		log.Printf("Response: %v\n", response)
		return
	}

	log.Printf("Adding todo item: %v\n", todoItem)
	todos = append(todos, todoItem)
	save()

	response := Response{
		Ok:          true,
		Descirption: "Added to the DB",
		Result:      todoItem,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error writing json to reponse: %v\n", err)
		return
	}
	log.Printf("Response: %v\n", response)
}

func remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todoItem := TodoItem{}
	json.NewDecoder(r.Body).Decode(&todoItem)
	id := todoItem.ID
	log.Printf("Removing todo item #%v\n", id)
	var removed *TodoItem
	removedIdx := -1
	for i, ti := range todos {
		if ti.ID == id {
			removed = &ti
            removedIdx = i
            break
		}
    }
    log.Printf("removedIdx: %v\n", removedIdx)

	var response Response
	if removed == nil {
		response = Response{
			Ok:          false,
			Descirption: "ID not in DB",
		}
	} else {
		todos = append(todos[:removedIdx], todos[removedIdx+1:]...)
		save()
		response = Response{
			Ok:          true,
			Descirption: "Removed from DB",
			Result:      *removed,
		}
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error writing json to reponse: %v\n", err)
		return
	}
	log.Printf("Response: %v\n", response)
}

func main() {
    log.Printf("args: %v\n", os.Args)
	router := mux.NewRouter().StrictSlash(true)
    
    portNumber := flag.Int("port", 1337, "The port number to run the API on")
    flagDbFilename := flag.String("db", "db.json", "The name of the DB file")
    flag.Parse()

    log.Printf("portNumber: %d, dbFilename: %s\n", *portNumber, *flagDbFilename)
    
	port := fmt.Sprintf(":%d", *portNumber)
    dbFilename = *flagDbFilename

	router.HandleFunc("/getAll", getAll)
	router.HandleFunc("/get", get)
	router.HandleFunc("/add", add)
	router.HandleFunc("/remove", remove)
	read()
	log.Printf("Listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

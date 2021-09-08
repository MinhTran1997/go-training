package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-mongodb-api/models"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var employeesCollection *mongo.Collection

func connectToDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	// connect to MongoDB cloud via given URI
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Ping to check if the connection is valid, successfully
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// get all the current databases in cloud just to make sure that we can retreive the data in cloud (this step is not mandatory)
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	// -------- get data in MongoDB cloud of a specific database-collection, then assgin to employeesCollection variable to store data in local code --------
	// note: any changes, update in the local code also affect the data in cloud
	gotrainingDatabase := client.Database("gotraning")
	employeesCollection = gotrainingDatabase.Collection("employees")
}

func disconnectDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client.Disconnect(ctx)
}

func HomeLink(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome to 'RESTful API with Golang and MongoDB' session!!!")
}

func GetAllEmployees(response http.ResponseWriter, request *http.Request) {

	connectToDatabase()
	defer disconnectDatabase()

	// set header for request
	response.Header().Set("content-type", "application/json")

	// var employees is a list, plural noun (s) because this is a getAll method
	var employees []models.Employee
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// https://docs.mongodb.com/manual/reference/method/db.collection.find/
	// Selects documents in a collection or view and returns a cursor to the selected documents.
	cursor, err := employeesCollection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)

	// The next document in the cursor returned by the db.collection.find() method
	// https://docs.mongodb.com/manual/reference/method/cursor.next/#mongodb-method-cursor.next
	// loop through each item in document, decode that item then assign it to pointer employee var, then append it to employees var (list)
	for cursor.Next(ctx) {
		var employee models.Employee
		cursor.Decode(&employee)
		employees = append(employees, employee)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// send employees var to the api's reponse
	json.NewEncoder(response).Encode(employees)
}

func GetEmployeeByID(response http.ResponseWriter, request *http.Request) {

	connectToDatabase()
	defer disconnectDatabase()

	// set header for request
	response.Header().Set("content-type", "application/json")

	// var employees is not a list, singular noun (s) because this is a getById method
	var employee models.Employee
	// get the ID parameter in the request
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"])

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// https://docs.mongodb.com/manual/reference/method/db.collection.findOne/#mongodb-method-db.collection.findOne
	// Returns one document that satisfies the specified query criteria on the collection or view. If multiple documents satisfy the query, t
	// his method returns the first document according to the natural order which reflects the order of documents on the disk
	// . If no document satisfies the query, the method returns null.
	// Although similar to the find() method, the findOne() method returns a document rather than a cursor.
	// --> findOne based on ID

	// find data then decode that data then assign to pointer employee
	err := employeesCollection.FindOne(ctx, models.Employee{ID: id}).Decode(&employee)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// send employee var to the api's reponse
	json.NewEncoder(response).Encode(employee)
}

func DeleteEmployeeByID(response http.ResponseWriter, request *http.Request) {

	connectToDatabase()
	defer disconnectDatabase()

	// set header for request
	response.Header().Set("content-type", "application/json")

	// get the ID parameter in the request
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"])

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// https://docs.mongodb.com/manual/reference/method/db.collection.deleteOne/
	// Removes a single document from a collection.
	// Returns:	A document containing:
	//  -- A boolean acknowledged as true if the operation ran with write concern or false if write concern was disabled
	//  -- deletedCount containing the number of deleted documents
	result, err := employeesCollection.DeleteOne(ctx, models.Employee{ID: id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// send employee var to the api's reponse
	json.NewEncoder(response).Encode(result)
}

func AddEmployee(response http.ResponseWriter, request *http.Request) {

	connectToDatabase()
	defer disconnectDatabase()

	// set header for request
	response.Header().Set("content-type", "application/json")

	// var employees is not a list, singular noun (s) because this is a addEmployee
	var employee models.Employee

	// get the request body, decode the request then assign it to pointer employee
	json.NewDecoder(request.Body).Decode(&employee)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// https://docs.mongodb.com/manual/reference/method/db.collection.insertOne/
	// Inserts a document into a collection.
	// Returns:	A document containing:
	// -- A boolean acknowledged as true if the operation ran with write concern or false if write concern was disabled.
	// -- A field insertedId with the _id value of the inserted document.
	result, _ := employeesCollection.InsertOne(ctx, employee)

	// assign the auto-generated ID in cloud database to the ID in the employeesCollection var
	id := result.InsertedID
	employee.ID = id.(primitive.ObjectID)

	// the employee parameter in Encode method is not the document in the cloud databse, it is the variable in the local code
	// the two reference each other
	json.NewEncoder(response).Encode(employee)
}

func UpdateEmployeeById(response http.ResponseWriter, request *http.Request) {

	connectToDatabase()
	defer disconnectDatabase()

	// set header for request
	response.Header().Set("content-type", "application/json")

	var employee models.Employee
	id, _ := primitive.ObjectIDFromHex(mux.Vars(request)["id"]) //get ID in request
	json.NewDecoder(request.Body).Decode(&employee)             // get request body then assign to var employee
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//https://docs.mongodb.com/manual/reference/method/db.collection.updateOne/
	// 	The method returns a document that contains:
	// -- matchedCount containing the number of matched documents
	// -- modifiedCount containing the number of modified documents
	// -- upsertedId containing the _id for the upserted document.
	// -- A boolean acknowledged as true if the operation ran with write concern or false if write concern was disabled

	result, err := employeesCollection.UpdateOne(ctx, models.Employee{ID: id}, bson.M{"$set": employee})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(result)
}

package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-echo-api-mongodb/models"

	"github.com/labstack/echo/v4"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type response struct {
	ID      primitive.ObjectID `json:"id,omitempty"`
	Message string             `json:"message,omitempty"`
}

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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(databases)

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

//-------- body-dump middleware handler --------
func BodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
	fmt.Printf("----------------------------------------\n")
}

//-------- handler functions for router --------
func GetAllEmployees(c echo.Context) error {

	connectToDatabase()
	defer disconnectDatabase()

	// var employees is a list, plural noun (s) because this is a getAll method
	var employees []models.Employee
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// https://docs.mongodb.com/manual/reference/method/db.collection.find/
	// Selects documents in a collection or view and returns a cursor to the selected documents.
	cursor, err := employeesCollection.Find(ctx, bson.M{})
	if err != nil {
		return err
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
		return err
	}

	// send employees var to the api's reponse
	return c.JSON(http.StatusOK, employees)
}

func GetEmployeeByID(c echo.Context) error {

	connectToDatabase()
	defer disconnectDatabase()

	// var employees is not a list, singular noun (s) because this is a getById method
	var employee models.Employee
	// get the ID parameter in the request
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

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
		return err
	}

	// send employee var to the api's reponse
	return c.JSON(http.StatusOK, employee)
}

func DeleteEmployeeByID(c echo.Context) error {

	connectToDatabase()
	defer disconnectDatabase()

	// get the ID parameter in the request
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// https://docs.mongodb.com/manual/reference/method/db.collection.deleteOne/
	// Removes a single document from a collection.
	// Returns:	A document containing:
	//  -- A boolean acknowledged as true if the operation ran with write concern or false if write concern was disabled
	//  -- deletedCount containing the number of deleted documents
	_, err := employeesCollection.DeleteOne(ctx, models.Employee{ID: id})
	if err != nil {
		return err
	}

	res := response{
		ID:      id,
		Message: "Employee is successfully deleted!!!",
	}

	// send employee var to the api's reponse
	return c.JSON(http.StatusOK, res)
}

func AddEmployee(c echo.Context) error {

	connectToDatabase()
	defer disconnectDatabase()

	// var employees is not a list, singular noun (s) because this is a addEmployee
	var employee models.Employee

	// get the request body, decode the request then assign it to pointer employee
	err := c.Bind(&employee)
	if err != nil {
		return err
	}

	// https://pkg.go.dev/context#WithTimeout
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
	return c.JSON(http.StatusCreated, employee)
}

func UpdateEmployeeById(c echo.Context) error {

	connectToDatabase()
	defer disconnectDatabase()

	var employee models.Employee
	id, _ := primitive.ObjectIDFromHex(c.Param("id")) //get ID in request

	err := c.Bind(&employee)
	if err != nil {
		return err
	} // get request body then assign to var employee
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
		return err
	}

	return c.JSON(http.StatusOK, result)
}

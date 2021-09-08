package handlers

// The handlers package is the bridge between APIs and Database.
// This package will handle all the db operations like Insert, Select, Update, and Delete (CRUD).

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go-postgres-api/models" // models package where User schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db (which is ElephantSQL here)
func createConnection() *sql.DB {
	// use godotenv to load the .env file
	// which loads env vars from a .env file
	// --> https://github.com/joho/godotenv --> USAGE
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read the POSTGRES_URL from the .env and connect to the db.
	// Open opens a database specified by its database driver name and a driver-specific data source name usually consisting of at least a database name and connection information.
	// func Open(driverName, dataSourceName string) (*DB, error)
	// Open may just validate its arguments without creating a connection to the database.
	// To verify that the data source name is valid, call Ping.
	// https://pkg.go.dev/database/sql#Open for more details

	// Getenv retrieves the value of the environment variable named by the key
	// func Getenv(key string) string
	// https://pkg.go.dev/os#Getenv for more details
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// The panic function is an inbuilt function which is defined under the builtin package of the Go language.
	// This function terminates the flow of control and starts panicking.

	// check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// return the connection
	return db
}

func HomeLink(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome to 'RESTful API with Golang and Postgres' session!!!")
}

// CreateUser create a user in the postgres db
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an EMPTY user of type models.User
	var user models.User

	// decode the json request to User type
	// NewDecoder decode the body that sent from request
	// then, Decode(&user) reads the next JSON-encoded value from its input and stores it in the value pointed to by &user
	// --> output of Decode method is an "error" and it also will store a json value to it's input
	// https://pkg.go.dev/encoding/json#Decoder.Decode for more details
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	// insertUser(user) returns new added ID
	insertID := insertUser(user)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	// send the response
	// the input param of Encode method is the value that we want to show in the reponse
	json.NewEncoder(w).Encode(res)
}

// GetUser will return a single user by its id
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getUser function with user id to retrieve a single user
	// getUser(int64(id)) returns user data and error (if it exists)
	user, err := getUser(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send the response
	// the input param of Encode method is the value that we want to show in the reponse
	json.NewEncoder(w).Encode(user)
}

// GetAllUser will return all the users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get all the users in the db
	// getAllUsers() returns all users data and error (if it exists)
	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	// send all the users as response
	// the input param of Encode method is the value that we want to show in the reponse
	json.NewEncoder(w).Encode(users)
}

// UpdateUser update user's detail in the postgres db
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty user of type models.User
	var user models.User

	// decode the json request to user
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update user to update the user
	updatedRows := updateUser(int64(id), user)

	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/records affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	// the input param of Encode method is the value that we want to show in the reponse
	json.NewEncoder(w).Encode(res)
}

// DeleteUser delete user's detail in the postgres db
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteUser, convert the int to int64
	deletedRows := deleteUser(int64(id))

	// format the message string
	msg := fmt.Sprintf("User deleted successfully. Total rows/records affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	// the input param of Encode method is the value that we want to show in the reponse
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions in DB -------------------------

// insert one user in the DB
func insertUser(user models.User) int64 {

	// create the postgres db connection
	// createConnection() funciton connect to the db, and returns an DB type to allow us to interact with the database via its methods
	// https://pkg.go.dev/database/sql#pkg-types --> type DB for more details
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the new inserted ID in the &id (pointer)
	// --> https://pkg.go.dev/database/sql#Rows.Scan
	// QueryRow executes a query that is expected to return at most one row.
	// func (db *DB) QueryRow(query string, args ...interface{}) *Row
	// https://pkg.go.dev/database/sql#DB.QueryRow
	err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single user with ID %v", id)

	// return the inserted id
	return id
}

// get one user from the DB by its userid
func getUser(id int64) (models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	// func in defer statement will execute when all the surroundings function returns
	defer db.Close()

	// create a user of models.User type
	var user models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE userid=$1`

	// execute the sql statement
	// QueryRow executes a query that is expected to return at most one row.
	// func (db *DB) QueryRow(query string, args ...interface{}) *Row
	// https://pkg.go.dev/database/sql#DB.QueryRow
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
	// Scan function will save the user info that we get into &user
	// --> https://pkg.go.dev/database/sql#Rows.Scan
	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	fmt.Printf("Get a user")
	// return empty user on error
	return user, err
}

// get one user from the DB by its userid
func getAllUsers() ([]models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var users []models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users`

	// execute the sql statement
	// Query executes a query that returns rows, typically a SELECT.
	// func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
	// https://pkg.go.dev/database/sql#DB.Query
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var user models.User

		// unmarshal the row object to user
		// Scan function will save the user info that we get into &user
		// --> https://pkg.go.dev/database/sql#Rows.Scan
		err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		users = append(users, user)

	}

	fmt.Printf("Get all users")
	// return empty user on error
	return users, err
}

// update user in the DB
func updateUser(id int64, user models.User) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`

	// execute the sql statement
	// Exec executes a query without returning any rows (such as Update, delete)
	// --> func (db *DB) Exec(query string, args ...interface{}) (Result, error)
	// --> https://pkg.go.dev/database/sql#DB.Exec
	res, err := db.Exec(sqlStatement, id, user.Name, user.Location, user.Age)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Update user, total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete user in the DB
func deleteUser(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM users WHERE userid=$1`

	// execute the sql statement
	// Exec executes a query without returning any rows (such as Update, delete)
	// --> func (db *DB) Exec(query string, args ...interface{}) (Result, error)
	// --> https://pkg.go.dev/database/sql#DB.Exec
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Delete user, total rows/record affected %v", rowsAffected)

	return rowsAffected
}

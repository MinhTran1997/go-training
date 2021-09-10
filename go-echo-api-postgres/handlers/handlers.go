package handlers

import (
	"database/sql"
	"fmt"
	"go-echo-api-postgres/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

//-------- body-dump middleware handler --------
func BodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
	fmt.Printf("----------------------------------------\n")
}

//-------- handler functions for routers --------
func GetAllUsers(c echo.Context) error {
	users, err := getAllUsers()
	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	user, err := getUser(int64(id))
	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return err
	}

	insertID := insertUser(user)

	res := response{
		ID:      insertID,
		Message: "User created successfully!!!",
	}

	return c.JSON(http.StatusCreated, res)
}

func UpdateUser(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(c.Param("id"))
	updatedRows := updateUser(int64(id), user)
	msg := fmt.Sprintf("User updated successfully!!! Total rows/records affected: %v", updatedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	return c.JSON(http.StatusOK, res)
}

func DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deleteUser(int64(id))
	msg := fmt.Sprintf("User with ID %v is successfully deleted!!!", id)
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	return c.JSON(http.StatusOK, res)
}

//-------- handler functions in DB --------

// create a connection to DB
func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database is successfully connected!!!")

	return db
}

// insert one user in the DB
func insertUser(user models.User) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

	var id int64

	err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	return id
}

// get one user from the DB based on its userid
func getUser(id int64) (models.User, error) {
	db := createConnection()
	defer db.Close()

	var user models.User

	sqlStatement := `SELECT * FROM users WHERE userid=$1`

	row := db.QueryRow(sqlStatement, id)

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

	return user, err
}

// get all users from the DB
func getAllUsers() ([]models.User, error) {
	db := createConnection()
	defer db.Close()

	var users []models.User

	sqlStatement := `SELECT * FROM users`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		users = append(users, user)
	}
	return users, err
}

// update user in the DB based on its userid
func updateUser(id int64, user models.User) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`

	res, err := db.Exec(sqlStatement, id, user.Name, user.Location, user.Age)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

// delete user in the DB based on its userid
func deleteUser(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM users WHERE userid=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

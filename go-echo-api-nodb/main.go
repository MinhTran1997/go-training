package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	user struct {
		// https://echo.labstack.com/guide/binding/
		// tag của struct là json nên các field sẽ đc bind dựa vào request body
		// json - source is request body. Uses Go json package for unmarshalling.
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	// tạo biến users có kiểu key là int, value là dạng struct user khai báo bên trên
	// data của user sẽ tương tự như sau: {"1":{"id":1,"name":"bocau"},"2":{"id":2,"name":"ben"}}
	users = map[int]*user{}

	// biến đếm sẽ tăng dần cho ID khi tạo mới để tránh bị trùng lặp
	seq = 1
)

//----------
// Handlers
//----------

func createUser(c echo.Context) error {
	// tạo 1 biến u pointer tới struct user và set id mới cho nó thông qua biến seq
	u := &user{
		ID: seq,
	}

	// https://echo.labstack.com/guide/binding/
	// Echo provides following method to bind data from different sources (path params, query params, request body) to structure using Context#Bind(i interface{}) method
	// --> create User là method post, nên phải có reqeust body. method bind() giúp bind data trong request body vào strcut đã được khai báo, ở đây là struct user
	err := c.Bind(u)
	if err != nil {
		return err
	}
	users[u.ID] = u // u lúc này đang chứa request body data sau khi được bind. đổ data của u vào users
	seq++           // seq + 1 để sau này có tạo mới thì id ko bị trùng

	// https://echo.labstack.com/guide/response/
	// Context#JSON(code int, i interface{}) can be used to encode a provided Go type into JSON and send it as response with status code.
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	// lấy param id trong request url
	id, _ := strconv.Atoi(c.Param("id"))

	// return user dựa vào id
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c echo.Context) error {
	// tạo ra 1 biến u có kiểu dữ liệu của struct user, nhưng lúc này u đang rỗng ko có data
	u := new(user)

	// vì là method PUT, phải có request body nên phải dùng bind() để đổ data từ request body vào 1 local var
	// bind(u), đổ data của request body vào biến u vừa tạo ở bên trên
	err := c.Bind(u)
	if err != nil {
		return err
	}

	// lấy param id trong request url
	id, _ := strconv.Atoi(c.Param("id"))
	// lấy name trong request body gán cho name hiện tại có trong data để update, dựa vào id lấy ra từ request url
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func deleteUser(c echo.Context) error {
	// lấy param id trong request url
	id, _ := strconv.Atoi(c.Param("id"))

	// delete là 1 function có sẵn trong go dùng để delete data trong map dựa vào key
	// https://golang.org/doc/go1#delete
	// khi xóa là xóa dựa vào key của users chứ ko phải key của user, xem dòng 24 để hiểu thêm
	delete(users, id)

	// Context#NoContent(code int) can be used to send empty body with status code
	// https://echo.labstack.com/guide/response/
	return c.NoContent(http.StatusNoContent)
}

func getAllUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

// bodydump handler is to captures the request and response payload and calls the registered handler. Generally used for debugging/logging purpose
// https://echo.labstack.com/middleware/body-dump/
func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
	fmt.Printf("----------------------------------------\n")
}

func main() {
	e := echo.New()

	// Middleware

	// bodydump handler is to captures the request and response payload and calls the registered handler. Generally used for debugging/logging purpose
	// https://echo.labstack.com/middleware/body-dump/
	e.Use(middleware.BodyDump(bodyDumpHandler))

	// Logger middleware logs the information about each HTTP request (ko log dc respond vs request body, phải dùng bodyDump)
	// https://echo.labstack.com/middleware/logger/
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "host=${host}, method=${method}, uri=${uri}, status=${status}, error=${error}, message=${message}\n",
	}))

	// Recover middleware recovers from panics anywhere in the chain, prints stack trace and handles the control to the centralized HTTPErrorHandler.
	// https://echo.labstack.com/middleware/recover/
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", getAllUsers)
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// Start server
	// https://echo.labstack.com/guide/http_server/
	// Echo provides following convenience methods to start HTTP server with Echo as a request handler:
	// Echo.Start is convenience method that starts http server with Echo serving requests.
	e.Logger.Fatal(e.Start(":1323"))
}

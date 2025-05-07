package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init_database() {
	var err error
	connStr := "postgres://admin:9503@localhost:5432/users?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("error opening db")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("error connecting to database")
	}
	fmt.Println("conection to db successfull")
}

type User struct {
	Id       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// registration
func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Fatal(err)
		return
	}
	query := `INSERT INTO users (id,name,email,password) VALUES ($1,$2,$3,$4)`
	_, err := db.Exec(query, user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		log.Fatal("error inserting user")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error inserting user"})
		return
	}
	fmt.Println("user inserted successfully")
}

//login auth

//logout

// get_user
func get_user(c *gin.Context) {
	var user User
	name := c.Query("name")
	query := `SELECT * FROM users WHERE name = $1`
	err := db.QueryRow(query, name).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		log.Fatal("error querying database")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error getting user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// update
func update_user(c *gin.Context) {
	name := c.Query("name")
	email := c.Query("email")
	password := c.Query("password")
	user := c.Query("user")
	query := `UPDATE users SET name = $1, email = $2, password = $3 WHERE name = $4`
	_, err := db.Exec(query, name, email, password, user)
	if err != nil {
		log.Fatal("error changing ingo")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error updating user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})

}

//preferences

func main() {
	init_database()

	router := gin.Default()
	router.POST("/register", register)
	router.GET("/get_user", get_user)
	router.POST("/update_user", update_user)
	router.Run()
}

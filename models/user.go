package models

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Database connection established")
}

type User struct {
	ClientID int
	UserName string
	FullName string
}

type Result struct {
	Name   string
	ID     int
	Email  string
	Status int
	Alert  string
}

// Fetch user list from client master
func GetUsers() ([]User, error) {
	rows, err := db.Query("SELECT client_id, username, full_name FROM client_master ORDER BY client_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ClientID, &user.UserName, &user.FullName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	//fmt.Println(users)
	return users, nil
}

func GetLogeddetails(login_username, login_password string) (Result, error) {
	//fmt.Println(login_username)
	//fmt.Println(login_password)

	sqlStatement := `SELECT client_id, full_name, password, status FROM client_master WHERE username = $1;`
	var client_id int
	var full_name string
	var password string
	var status int

	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := db.QueryRow(sqlStatement, login_username)
	switch err := row.Scan(&client_id, &full_name, &password, &status); err {
	case sql.ErrNoRows:
		//fmt.Println("Data Not Found")
		message := Result{
			Alert: "Data Not Found",
		}
		return message, nil
	case nil:
		//fmt.Println(client_id, full_name, password, status)

		if status != 1 {
			message := Result{
				Alert: "Account Not",
			}
			return message, nil
		}

		// func CompareHashAndPassword(hashedPassword, password []byte) error
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(login_password))

		// returns nill
		if err == nil {
			//fmt.Println("You have successfully logged in :")
			message := Result{
				Name:   full_name,
				ID:     client_id,
				Email:  login_username,
				Status: status,
			}
			// manage login history
			//fmt.Println(GetLocalIP())

			var ip = "192.168.29.4"
			//var ip = string(GetLocalIP())
			fmt.Println(ip)
			sqlStatement := `INSERT INTO login_history (client_id, login_ip) VALUES ($1,  $2);`
			db.QueryRow(sqlStatement, client_id, ip)

			return message, nil

		} else {

			message := Result{
				Alert: "incorrect password",
			}

			return message, nil

		}

	default:
		panic(err)
	}

}

func UsersRegistration(name, email string) (Result, error) {
	//fmt.Println(name)
	//fmt.Println(email)

	// create hash from password
	var password = "India123"

	var hash []byte
	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	sqlStatement := `
INSERT INTO client_master (username, full_name, password, status)
VALUES ($1, $2, $3, $4)
RETURNING client_id`
	id := 0
	err := db.QueryRow(sqlStatement, email, name, hash, 1).Scan(&id)
	if err != nil {
		panic(err)
	}

	message := Result{
		Name:   name,
		ID:     id,
		Email:  email,
		Status: 1,
	}
	return message, nil
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

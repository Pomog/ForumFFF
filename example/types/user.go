package types

import (
	"database/sql"
	"fmt"

	"forum-authentication/config"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Provider  string `json:"provider"`
}

func (u *User) CreateUser(user User, provider string) (int, error) {
	insertStmt := `INSERT INTO users (username, age, gender, firstname, lastname, email, password, provider) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := config.DB.Prepare(insertStmt)
	if err != nil {
		return 0, err
	}
	fmt.Println("provider")
	result, err := stmt.Exec(user.Username, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, user.Password, provider)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}

func (u *User) GetUserByUsername(username string) (User, error) {
	stmt := `SELECT * FROM users WHERE username=?`

	err := config.DB.QueryRow(stmt, username).Scan(&u.Id, &u.Username, &u.Age, &u.Gender, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Provider)
	if err != nil {
		if err == sql.ErrNoRows {
			return *u, fmt.Errorf("User not found")
		}
		return *u, err
	}
	return *u, nil
}

func (u *User) GetUserByEmail(email string) (User, error) {
	stmt := `SELECT * FROM users WHERE email=?`

	err := config.DB.QueryRow(stmt, email).Scan(&u.Id, &u.Username, &u.Age, &u.Gender, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Provider)
	if err != nil {
		if err == sql.ErrNoRows {
			return *u, fmt.Errorf("User not found")
		}
		return *u, err
	}
	return *u, nil
}

func (u *User) GetUserFromSession(value string) (User, error) {
	stmt := `
	SELECT users.id, users.username, users.age, users.gender, users.firstname, users.lastname, users.email, users.password
	FROM sessions 
	JOIN users ON sessions.user_id = users.id 
	WHERE sessions.value = ?
	`
	err := config.DB.QueryRow(stmt, value).Scan(&u.Id, &u.Username, &u.Age, &u.Gender, &u.FirstName, &u.LastName, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("User not found")
		}
		return User{}, err
	}
	return *u, nil
}

func (u *User) CheckCredentials(username, password string) (User, error) {
	user, err := u.GetUserByUsername(username)
	if err != nil {
		return User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return User{}, err
	}

	return user, nil
}

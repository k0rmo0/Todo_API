package model

import (
	"errors"
	"fmt"

	"database/sql"

	"github.com/bicom/todos/utils"
	"github.com/casbin/casbin"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("mykey")

var (
	//UserTypeAdmin ...
	UserTypeAdmin = "admin"
)

//User .
type User struct {
	Type            string                    `db:"type" json:"type"`
	ID              int                       `db:"id" json:"id"`
	FirstName       string                    `db:"firstname" json:"firstname"`
	LastName        string                    `db:"lastname" json:"lastname"`
	Username        string                    `db:"username" json:"username"`
	Password        string                    `db:"password" json:"password"`
	Email           string                    `db:"email" json:"email"`
	Token           string                    `db:"token" json:"token"`
	Issued          int64                     `db:"issued" json:"issued"`
	UserPermissions map[string]PathPermission `db:"-" json:"user_permissions"`
}

//PathPermission ...
type PathPermission struct {
	Read  bool
	Write bool
}

//IsAdmin ...
func (m *User) IsAdmin() bool {
	return m.Type == UserTypeAdmin
}

//SetPermissions ...
func (m *User) SetPermissions(rules *casbin.Enforcer) {
	var userPermissions = make(map[string]PathPermission)

	typePermissions := rules.GetImplicitPermissionsForUser(m.Type)
	emailPermissions := rules.GetImplicitPermissionsForUser(m.Email)

	userPermissions = MergePermissions(typePermissions, userPermissions)
	fmt.Println(userPermissions)
	userPermissions = MergePermissions(emailPermissions, userPermissions)

	m.UserPermissions = userPermissions
}

//MergePermissions ...
func MergePermissions(listPermissions [][]string, userPermissions map[string]PathPermission) map[string]PathPermission {
	for _, item := range listPermissions {

		var existingPermission = false
		var tmpPermissions PathPermission

		for path, permissions := range userPermissions {
			if path == item[1] {
				existingPermission = true

				if item[2] == "read" {
					permissions.Read = true
					userPermissions[path] = permissions
				} else if item[2] == "write" {
					permissions.Write = true
					userPermissions[path] = permissions

				}
			}
		}

		if !existingPermission {
			if item[2] == "read" {
				tmpPermissions.Read = true
				tmpPermissions.Write = false

				userPermissions[item[1]] = tmpPermissions

			} else if item[2] == "write" {
				tmpPermissions.Read = false
				tmpPermissions.Write = true

				userPermissions[item[1]] = tmpPermissions
			}
		}
	}
	return userPermissions
}

//Create ...
func (m *User) Create() error {

	db := utils.SQLAcc.GetSQLDB()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(m.Password), 10)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO users(username, firstname, lastname, email, password, type) values (?, ?, ?, ?, ?, ?)", m.Username, m.FirstName, m.LastName, m.Email, bytes, "user")

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

//Login ...
func (m *User) Login() (User, error) {
	var user User
	db := utils.SQLAcc.GetSQLDB()
	err := db.Get(&user, "SELECT * FROM users WHERE username = ?", m.Username)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(m.Password))
	if err != nil {
		return user, err
	}

	return user, nil
}

//ListUsers ...
func ListUsers(exclude int) ([]User, error) {
	db := utils.SQLAcc.GetSQLDB()

	var rows []User
	var err error

	query := "SELECT * FROM users"

	if exclude != 0 {
		query += ` WHERE id <> ?`

		err = db.Select(&rows, query, exclude)
	}

	if err == sql.ErrNoRows {
		return rows, errors.New("Can't get all users")
	}

	if err != nil {
		return rows, err
	}

	return rows, nil
}

//IsLoggedIn ...
func (m *User) IsLoggedIn(identity string) error {
	db := utils.SQLAcc.GetSQLDB()

	err := db.Get(m, `SELECT * FROM users WHERE email=? OR username=?`, identity, identity)

	if err == sql.ErrNoRows {
		return errors.New("User with email/username '" + identity + "' does not exist")
	}

	return err
}

//UpdateTokenInfo ...
func (m *User) UpdateTokenInfo() error {
	db := utils.SQLAcc.GetSQLDB()
	_, err := db.Exec(`UPDATE users SET token= ?, issued=? WHERE username=?`, m.Token, m.Issued, m.Username)

	if err != nil {
		return err
	}

	return nil
}

//UpdatePassword ...
func (m *User) UpdatePassword(oldpass string, newpass string) error {
	db := utils.SQLAcc.GetSQLDB()

	var user User

	err := db.Get(&user, "SELECT * FROM users WHERE username = ?", m.Username)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldpass))

	if err != nil {
		return err
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(newpass), 12)
	_, err = db.Exec(`UPDATE users SET password=? WHERE username=?`, bytes, m.Username)

	if err != nil {
		return err
	}

	return nil
}

//UpdateType ...
func (m *User) UpdateType() error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec(`UPDATE users SET type=? WHERE id=?`, m.Type, m.ID)

	if err != nil {
		return err
	}

	return nil
}

//DeleteUser ...
func (m *User) DeleteUser(UserID string) error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec("DELETE FROM users WHERE id=?", UserID)
	if err != nil {
		return err
	}
	return nil
}

//GetUser ...
func (m *User) GetUser(UserID string) (User, error) {
	db := utils.SQLAcc.GetSQLDB()

	var user User

	err := db.Get(&user, "SELECT * FROM users WHERE id=?", UserID)

	if err == sql.ErrNoRows {
		return user, errors.New("User does not exist")
	}
	if err != nil {
		return user, err
	}

	return user, nil
}

//Clear Token token data used for logout
func (m *User) Clear() error {
	db := utils.SQLAcc.GetSQLDB()

	m.Token = ""
	m.Issued = 0

	_, err := db.Exec(`UPDATE users SET token=?, issued=? WHERE username=?`, m.Token, m.Issued, m.Username)

	if err != nil {
		return err
	}

	return nil
}

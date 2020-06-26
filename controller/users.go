package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/bicom/todos/model"
	"github.com/bicom/todos/utils"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

//Users struct .
type Users struct{}

//Create ...
func (uc Users) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.WriteJSON(w, err, http.StatusBadRequest)
		return
	}

	//Checks contents of register fields ...
	if user.Username == "" || len(user.Username) < 3 || user.Password == "" || len(user.Password) < 3 {
		errMess := "Can't leave user/password field blank, or less than 3 character"
		utils.WriteJSON(w, errMess, http.StatusBadRequest)
		return
	}

	regexpEmail := regexp.MustCompile("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$")

	if user.Email == "" || !regexpEmail.MatchString(user.Email) {
		err = errors.New("Email is empty or wrong email format")
		utils.WriteJSON(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	err = user.Create()

	if err != nil {
		utils.WriteJSON(w, err, http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, "Succesfully registered !", http.StatusOK)
}

//Login ...
func (uc Users) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := context.Get(r, "user").(model.User)
	utils.WriteJSON(w, u, http.StatusAccepted)
}

//ListAll ...
func (uc Users) ListAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := context.Get(r, "user").(model.User)
	var users []model.User
	var err error

	if !user.IsAdmin() {
		utils.WriteJSON(w, "You ar not allowed to list users", http.StatusInternalServerError)
		return
	}

	users, err = model.ListUsers(user.ID)

	if err != nil {
		utils.WriteJSON(w, "Error listing users", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, users, http.StatusOK)
}

//UpdatePassword ...
func (uc Users) UpdatePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := context.Get(r, "user").(model.User)

	oldPass := r.URL.Query().Get("oldpass")
	newPass1 := r.URL.Query().Get("newpass")

	if oldPass != "" && newPass1 != "" {
		if len(newPass1) <= 4 {
			err := errors.New("Password must be longer than 4 characters")
			utils.WriteJSON(w, err, http.StatusForbidden)
			return

		}
	}
	err := user.UpdatePassword(oldPass, newPass1)

	if err != nil {
		utils.WriteJSON(w, "Old passwords do not match", http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, "Successfully changed your password", http.StatusOK)

}

//UpdatePassword2 ...
func (uc Users) UpdatePassword2(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := context.Get(r, "user").(model.User)

	type parameters struct {
		Oldpassword string `json:"oldpass"`
		Newpassword string `json:"newpass"`
	}

	var par parameters

	err := json.NewDecoder(r.Body).Decode(&par)

	if err != nil {
		utils.WriteJSON(w, err, http.StatusInternalServerError)
		return
	}

	if par.Newpassword != "" && par.Oldpassword != "" {
		if len(par.Newpassword) <= 4 {
			err = errors.New("Password must be longer than 4 characters")
			utils.WriteJSON(w, err, http.StatusForbidden)
			return
		}
	}

	err = user.UpdatePassword(par.Oldpassword, par.Newpassword)
	if err != nil {
		utils.WriteJSON(w, "Error updating password", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, "Password successfully changed", http.StatusAccepted)
}

//UpdateType ...
func (uc Users) UpdateType(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := context.Get(r, "user").(model.User)

	var u model.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		utils.WriteJSON(w, err, http.StatusInternalServerError)
		return
	}

	if user.IsAdmin() && u.Type != "" {

		err = u.UpdateType()

		if err != nil {
			utils.WriteJSON(w, "Error updating user type", http.StatusInternalServerError)
		}
		utils.WriteJSON(w, "User updated", http.StatusOK)
		return
	}

	utils.WriteJSON(w, "You have no perrmission to update user type", http.StatusUnauthorized)
}

//DeleteUser ...
func (uc Users) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := context.Get(r, "user").(model.User)

	if !user.IsAdmin() {
		err := errors.New("You are not allowed to delete user")
		utils.WriteJSON(w, err, http.StatusInternalServerError)
		return
	}

	userID := params.ByName("id")

	err := user.DeleteUser(userID)

	if err != nil {
		utils.WriteJSON(w, "Error deleting user", http.StatusInternalServerError)
	}

	utils.WriteJSON(w, "User deleted", http.StatusOK)
}

//GetUser ...
func (uc Users) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := context.Get(r, "user").(model.User)

	userID := params.ByName("id")

	if userID != userID && !user.IsAdmin() {
		err := errors.New("You are not allowed to see this user")
		utils.WriteJSON(w, err, http.StatusUnauthorized)
		return
	}

	user, err := user.GetUser(userID)
	if err != nil {
		utils.WriteJSON(w, err, http.StatusInternalServerError)
	}

	utils.WriteJSON(w, user, http.StatusOK)
}

//Logout ...
func (uc Users) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := context.Get(r, "user").(model.User)
	err := user.Clear()

	if err != nil {
		utils.WriteJSON(w, "Error logging out", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, "Logged out", http.StatusOK)
}

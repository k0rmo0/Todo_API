package middlleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/casbin/casbin"
	jwtreq "github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/bicom/todos/model"
	"github.com/bicom/todos/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

var (
	pubkey []byte
)

//Provider ...
type Provider struct {
	rules *casbin.Enforcer
}

//Middlleware ...
type Middlleware struct{}

var (
	jwtKey = []byte("mykey")
)

var (
	// ErrUserTypeNotDefine ..
	ErrUserTypeNotDefine = errors.New("User Type not set on account")
	// ErrUserNotFound user not activated
	ErrUserNotFound = errors.New("This user doesn't exists")
)

// Clear ...
func (m Middlleware) Clear(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	defer context.Clear(req)

	next(res, req)
}

// CORS ...
func (m Middlleware) CORS(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// CORS support for Preflighted requests
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

	next(res, req)
}

// Preflight ...
func (m Middlleware) Preflight(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if req.Method == "OPTIONS" {
		utils.Renderer.Render(res, http.StatusOK, map[string]string{"status": "OK"})
		return
	}

	next(res, req)
}

// SetRBAC ...
func (m *Provider) SetRBAC(model, policies string) error {
	var err error
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	m.rules, err = casbin.NewEnforcerSafe(path+model, path+policies, false)

	if err != nil {
		return err
	}

	return nil
}

//JWT ...
func (m Provider) JWT(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//Checking uri strings ...
	if strings.Contains(r.RequestURI, "/register") && r.Method == "POST" {
		next(w, r)
		return
	} else if strings.Contains(r.RequestURI, "/login") {

		username, password, ok := r.BasicAuth()

		if !ok {
			utils.WriteJSON(w, "Authorization header format must be Basic ", http.StatusUnauthorized)
			return
		}

		user := model.User{Username: username, Password: password}

		user, _, err := issueToken(r, user)

		if err != nil {
			utils.WriteJSON(w, err, http.StatusForbidden)
			return
		}

		user.SetPermissions(m.rules)

		context.Set(r, "user", user)
		next(w, r)
	} else if strings.Contains(r.RequestURI, "/logout") {
		user, _, _ := checkToken(r)

		user.SetPermissions(m.rules)

		context.Set(r, "user", user)
		next(w, r)
	} else {
		user, _, err := checkToken(r)

		if err != nil {
			utils.WriteJSON(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		var uri []string

		idx := strings.Index(r.RequestURI, "?")

		if idx == int(-1) {
			uri = strings.Split(r.RequestURI, ".")
		} else {
			uri = strings.Split(r.RequestURI[:idx], ".")
		}

		route := strings.Split(uri[0], "/")

		fullAction := uri[0]
		action := "/"

		if len(route) > 1 {
			action += strings.Split(uri[0], "/")[1]
		}

		fmt.Println(fullAction)
		fmt.Println(action)

		user.SetPermissions(m.rules)

		var userAllowedByRole = false

		if action == "/users" || action == "/user" || action == "/todo" ||
			action == "/todos" || action == "/task" || action == "/tasks" {
			if m.rules != nil {
				if m.rules.Enforce(user.Type, fullAction, requestMethod2Mode(r.Method)) {
					fmt.Println("RBAC allowed by type " + user.Type + " and email " + user.Email + " for action " + fullAction)
					userAllowedByRole = true
				}
			}
			if !userAllowedByRole {
				if m.rules != nil && m.rules.Enforce(user.Type, fullAction, requestMethod2Mode(r.Method)) {
					fmt.Println("RBAC allowed by type " + user.Type + "for action" + fullAction)
				} else if m.rules != nil && !m.rules.Enforce(user.Email, fullAction, requestMethod2Mode(r.Method)) {
					fmt.Println("Forbidden RBAC per Email for action"+fullAction, errors.New("Forbidden RBAC"))
					utils.WriteJSON(w, "Acces Forbidden", http.StatusForbidden)
					return
				}
			}
		}
		context.Set(r, "user", user)

		next(w, r)
	}

}

func checkToken(r *http.Request) (model.User, *jwt.Token, error) {
	token, err := jwtreq.ParseFromRequest(r, jwtreq.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Client is not using the correct algorithm")
		}
		return jwtKey, nil
	})

	var user model.User
	if err != nil || !token.Valid {
		return user, token, errors.New("Invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	err = user.IsLoggedIn(claims["usr"].(string))

	if err != nil {
		return user, token, errors.New("User for token not found")
	}
	if user.Issued != 0 && time.Now().After(time.Unix(user.Issued, 0)) {
		return user, token, errors.New("Token expired")
	}
	if user.Issued == 0 {
		return user, token, errors.New("User logged out")
	}
	return user, token, nil
}

func issueToken(r *http.Request, m model.User) (model.User, *jwt.Token, error) {
	if m.Username == "" || m.Password == "" {
		err := errors.New("Missing username or password")

		return m, nil, err
	}

	m, err := m.Login()

	if err != nil {
		return m, nil, err
	}

	if m.ID == 0 {
		return m, nil, ErrUserNotFound
	}
	if m.Type == "" {
		return m, nil, ErrUserTypeNotDefine
	}

	if time.Now().Before(time.Unix(m.Issued, 0)) {
		return m, nil, nil
	}

	issued := time.Now().Add(time.Hour * 48).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": m.Username,
		"exp": issued,
	})

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return m, nil, err
	}

	m.Token = tokenString
	m.Issued = issued

	err = m.UpdateTokenInfo()
	if err != nil {
		return m, nil, err
	}

	m.Password = ""

	return m, nil, nil
}

//CheckTodo ...
func (m Middlleware) CheckTodo(h httprouter.Handle) httprouter.Handle {

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		user := context.Get(req, "user").(model.User)

		todoID, err := strconv.Atoi(params.ByName("id"))

		todo, err := model.GetAnyToDo(todoID)

		if err != nil {
			fmt.Println("The list is not reachable")
			utils.WriteJSON(res, err, 400)
			return
		}

		if todo.UserID != user.ID && !user.IsAdmin() {
			utils.WriteJSON(res, "You are not allowed to make any changes to this list", 403)
			return
		}

		h(res, req, params)
	}
}

//CheckTask ...
func (m Middlleware) CheckTask(h httprouter.Handle) httprouter.Handle {

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		user := context.Get(req, "user").(model.User)

		if req.Method == "POST" || req.Method == "GET" {

			todoID, err := strconv.Atoi(params.ByName("id"))

			var todo model.ToDo

			todo, err = model.GetAnyToDo(todoID)

			if err != nil {
				fmt.Println("The list is not reachable")
				utils.WriteJSON(res, err, 400)
				return
			}

			if todo.UserID != user.ID && user.Type != "admin" {
				utils.WriteJSON(res, "You are not allowed to make any changes to this list", 403)
				return
			}
		} else if req.Method == "DELETE" {

			taskID, err := strconv.Atoi(params.ByName("id"))

			var task model.Task
			var todo model.ToDo

			task, err = model.GetAnyTask(taskID)

			if err != nil {
				fmt.Println("The task is not reachable")
				utils.WriteJSON(res, err, 400)
				return
			}

			todo, err = model.GetAnyToDo(task.ToDoID)

			if err != nil {
				fmt.Println("The list is not reachable")
				utils.WriteJSON(res, err, 400)
				return
			}

			if todo.UserID != user.ID && !user.IsAdmin() {
				utils.WriteJSON(res, "You are not allowed to make any changes to this list", 403)
				return
			}
		}

		h(res, req, params)
	}
}

func requestMethod2Mode(reqMethod string) string {
	if reqMethod == "POST" || reqMethod == "PUT" || reqMethod == "DELETE" || reqMethod == "PATCH" {
		return "write"
	}

	return "read"
}

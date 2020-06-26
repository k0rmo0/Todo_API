package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bicom/todos/controller"
	middlleware "github.com/bicom/todos/middleware"
	"github.com/bicom/todos/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

var (
	users    controller.Users
	mdlw     middlleware.Middlleware
	task     controller.ToDoController
	provider middlleware.Provider
)

func main() {
	//connecting to DB
	utils.GetSQLDB("dev", "conf/conf.yaml")

	//RBAC configuration
	err := provider.SetRBAC("/conf/rbac.conf", "/conf/policy.csv")
	if err != nil {
		fmt.Println(err)
	}

	//ROUTER
	mux := httprouter.New()

	//NEGORNI MIDDLEWARE
	n := negroni.Classic()

	n.Use(negroni.HandlerFunc(mdlw.Clear))
	n.Use(negroni.HandlerFunc(mdlw.CORS))
	n.Use(negroni.HandlerFunc(mdlw.Preflight))
	n.Use(negroni.HandlerFunc(provider.JWT))

	//USER OPTIONS
	mux.POST("/register", users.Create)
	mux.POST("/login", users.Login)
	mux.GET("/users", users.ListAll)
	mux.PUT("/user/password1", users.UpdatePassword)
	mux.PUT("/user/password2", users.UpdatePassword2)
	mux.PUT("/user/type", users.UpdateType)
	mux.DELETE("/user/:id", users.DeleteUser)
	mux.GET("/user/:id", users.GetUser)
	mux.GET("/logout", users.Logout)

	//TODO
	mux.POST("/todo", task.CreateToDo)
	mux.POST("/task/:id", mdlw.CheckTask(task.CreateTask))

	mux.DELETE("/todo/:id", mdlw.CheckTodo(task.DeleteToDo))
	mux.DELETE("/task/:id", mdlw.CheckTask(task.DeleteTask)) //task ID

	mux.PUT("/todo/name/:id", mdlw.CheckTodo(task.UpdateToDoName))
	mux.PUT("/todo/description/:id", mdlw.CheckTodo(task.UpdateToDoDescription))
	mux.PUT("/task/name/:id", mdlw.CheckTodo(task.UpdateTaskName))         //ToDo ID
	mux.PUT("/task/date/:id", mdlw.CheckTodo(task.UpdateTaskDateFinish))   //ToDo ID
	mux.PUT("/task/priority/:id", mdlw.CheckTodo(task.UpdateTaskPriority)) //ToDo ID
	mux.PUT("/task/status/:id", mdlw.CheckTodo(task.UpdateTaskStatus))     //ToDo ID

	mux.GET("/todos", task.ListAllToDos)
	mux.GET("/task/:id", mdlw.CheckTask(task.ListTasks))                     //ToDo ID
	mux.GET("/tasks/active/:id", mdlw.CheckTask(task.ListAllActiveTasks))    //ToDo ID
	mux.GET("/tasks/completed/:id", mdlw.CheckTask(task.ListCompletedTasks)) //ToDo ID

	n.UseHandler(mux)

	fmt.Println("Server started...")
	log.Fatal(http.ListenAndServe(":8000", n))
}

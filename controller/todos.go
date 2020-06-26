package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bicom/todos/utils"

	"github.com/bicom/todos/model"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

//ToDoController ...
type ToDoController struct{}

//CreateToDo ...
func (tdc ToDoController) CreateToDo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := context.Get(r, "user").(model.User)

	var todo model.ToDo

	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		fmt.Println(err)
		utils.WriteJSON(w, err.Error(), 400)
		return
	}

	err = todo.CreateToDo(user.ID)
	if err != nil {
		fmt.Println(err)
		utils.WriteJSON(w, err.Error(), 400)
		return
	}

	utils.WriteJSON(w, todo, 200)
}

//CreateTask ...
func (tdc ToDoController) CreateTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	todoID, err := strconv.Atoi(params.ByName("id"))

	var task model.Task

	err = json.NewDecoder(r.Body).Decode(&task)

	err = task.CreateTask(todoID)
	if err != nil {
		fmt.Println(err)
		utils.WriteJSON(w, err.Error(), 400)
		return
	}

	utils.WriteJSON(w, task, 200)
}

//DeleteToDo ...
func (tdc *ToDoController) DeleteToDo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := context.Get(r, "user").(model.User)

	todoID, err := strconv.Atoi(params.ByName("id"))

	var todo model.ToDo

	if user.IsAdmin() {
		err = todo.DeleteToDo(0, todoID)

		if err != nil {
			utils.WriteJSON(w, err, 500)
			return
		}
	} else {
		err = todo.DeleteToDo(user.ID, todoID)

		if err != nil {
			fmt.Println(err)
			utils.WriteJSON(w, err.Error(), 403)
			return
		}
	}

	utils.WriteJSON(w, "ToDo table deleted.", 200)
}

//DeleteTask ...
func (tdc *ToDoController) DeleteTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	taskID, err := strconv.Atoi(params.ByName("id"))

	var task model.Task

	task, err = model.GetAnyTask(taskID)

	err = task.DeleteTask(task.ToDoID, task.ID)

	if err != nil {
		fmt.Println(err)
		utils.WriteJSON(w, err, 400)
		return
	}

	utils.WriteJSON(w, "Task deleted", 200)
}

//UpdateToDoName ...
func (tdc ToDoController) UpdateToDoName(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	todoID, err := strconv.Atoi(params.ByName("id"))

	var todo model.ToDo

	err = json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		utils.WriteJSON(w, "Error Unarshal name value", 500)
		return
	}

	err = todo.UpdateToDoName(todoID)

	if err != nil {
		utils.WriteJSON(w, err, 403)
		return
	}

	utils.WriteJSON(w, "Update done!", 200)
}

//UpdateToDoDescription ...
func (tdc ToDoController) UpdateToDoDescription(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	todoID, err := strconv.Atoi(params.ByName("id"))

	var todo model.ToDo

	err = json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		utils.WriteJSON(w, "Error Unarshal name value", 500)
		return
	}

	err = todo.UpdateToDoDescription(todoID)

	if err != nil {
		utils.WriteJSON(w, err, 403)
		return
	}

	utils.WriteJSON(w, "Update done!", 200)
}

//UpdateTaskName ...
func (tdc ToDoController) UpdateTaskName(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	todoID, err := strconv.Atoi(params.ByName("id"))

	var task model.Task

	err = json.NewDecoder(r.Body).Decode(&task) //send ID and name

	err = task.UpdateTaskName(todoID, task.ID)

	if err != nil {
		utils.WriteJSON(w, err, 500)
		return
	}

	utils.WriteJSON(w, "Update done!", 200)
}

//UpdateTaskDateStart ...
func (tdc ToDoController) UpdateTaskDateStart(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var task model.Task

	err := json.NewDecoder(r.Body).Decode(&task) //send ID, ToDoID and name

	err = task.UpdateTaskDateStart(task.ToDoID, task.ID)

	if err != nil {
		utils.WriteJSON(w, err, 500)
		return
	}

	utils.WriteJSON(w, "Update done!", 200)
}

//UpdateTaskDateFinish ...
func (tdc ToDoController) UpdateTaskDateFinish(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	todoID, err := strconv.Atoi(params.ByName("id"))

	var task model.Task

	err = json.NewDecoder(r.Body).Decode(&task) //send ID, ToDoID and name

	err = task.UpdateTaskDateFinish(todoID, task.ID)

	if err != nil {
		utils.WriteJSON(w, err, 500)
		return
	}

	utils.WriteJSON(w, "Update done!", 200)
}

//UpdateTaskPriority ...
func (tdc ToDoController) UpdateTaskPriority(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var task model.Task

	err := json.NewDecoder(r.Body).Decode(&task) //send ID, ToDoID and name

	err = task.UpdateTaskPriority(task.ToDoID, task.ID)

	if err != nil {
		utils.WriteJSON(w, err, 500)
		return
	}

	utils.WriteJSON(w, "Update done!", 200)
}

//UpdateTaskStatus ...
func (tdc ToDoController) UpdateTaskStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var task model.Task

	err := json.NewDecoder(r.Body).Decode(&task) //send ID, ToDoID and name

	err = task.UpdateTaskStatus(task.ToDoID, task.ID)

	if err != nil {
		utils.WriteJSON(w, err, 500)
		return
	}

	utils.WriteJSON(w, "Update done!", 200)
}

//ListAllToDos shows all ToDo lists created by users, admin can see all...
func (tdc ToDoController) ListAllToDos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := context.Get(r, "user").(model.User)

	var todos []model.ToDo
	var err error

	if user.IsAdmin() {
		todos, err = model.ListAllToDos(0)
		if err != nil {
			fmt.Println(err)
			utils.WriteJSON(w, "Unable to show all created ToDos", 400)
		}
	} else {
		todos, err = model.ListAllToDos(user.ID)
		if err != nil {
			fmt.Println(err)
			utils.WriteJSON(w, "Unable to show all created ToDos", 400)
		}
	}

	utils.WriteJSON(w, todos, 200)
}

//ListTasks shows all tasks per user or admin request, asks for id of todo ...
func (tdc ToDoController) ListTasks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	todoID, err := strconv.Atoi(params.ByName("id"))

	var tasks []model.Task

	tasks, err = model.ListTasks(todoID)

	if err != nil {
		utils.WriteJSON(w, err, 500)
		return
	}

	utils.WriteJSON(w, tasks, 200)
}

//ListAllActiveTasks shows all active tasks ...
func (tdc ToDoController) ListAllActiveTasks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	todoID, err := strconv.Atoi(params.ByName("id"))

	var activeTasks []model.Task

	activeTasks, err = model.ListAllActiveTasks(todoID)

	if err != nil {
		utils.WriteJSON(w, err, 500)
		return
	}

	utils.WriteJSON(w, activeTasks, 200)

}

//ListCompletedTasks ...
func (tdc ToDoController) ListCompletedTasks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	todoID, err := strconv.Atoi(params.ByName("id"))

	var activeTasks []model.Task

	activeTasks, err = model.ListAllActiveTasks(todoID)

	if err != nil {
		utils.WriteJSON(w, err, 500)
		return
	}

	utils.WriteJSON(w, activeTasks, 200)

}

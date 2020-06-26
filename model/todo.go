package model

import (
	"fmt"

	"github.com/bicom/todos/utils"
)

var user User

//ToDo ...
type ToDo struct {
	ID          int    `db:"id" json:"id"`                   //auto increment
	Name        string `db:"name" json:"name"`               //name of a to-do list
	Description string `db:"description" json:"description"` //more detailed info about to-do list
	UserID      int    `db:"userID" json:"userID"`           //ID from User struct
}

//Task contains a concrete task for to-do list ...
type Task struct {
	ID          int    `db:"id" json:"id"`     //auto increment
	Name        string `db:"name" json:"name"` //task name (what is supposed to be done)
	DateCreated string `db:"dateC" json:"dateCreated"`
	DateFinish  string `db:"dateF" json:"dateFinish"`
	Priority    string `db:"priority" json:"priority"` //value between 1-5
	Status      bool   `db:"status" json:"status"`     //not completed, completed (0,1)
	ToDoID      int    `db:"ToDoID" json:"todoID"`     //ID that is the same as ID from ToDo
}

//CreateToDoTable ...
var CreateToDoTable = `CREATE TABLE ToDo(
	id INT(11) NOT NULL AUTO_INCREMENT,
	name VARCHAR(150),
	description VARCHAR(255),
	userID VARCHAR(150),
	PRIMARY KEY(id)
	);
	`

//CreateTaskTable ...
var CreateTaskTable = `CREATE TABLE task(
	id INT(11) NOT NULL AUTO_INCREMENT,
	name VARCHAR(255),
	dateC VARCHAR(255),
	dateF VARCHAR(255),
	priority INT(11),
	status INT(11) DEFAULT '0',
	ToDoID INT(11),
	PRIMARY KEY(id)
	);
	`

//CreateToDo ...
func (td *ToDo) CreateToDo(userID int) error {
	db := utils.SQLAcc.GetSQLDB()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO ToDo (name, description, userID) VALUES(?, ?, ?)", td.Name, td.Description, userID)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

//CreateTask ...
func (ts *Task) CreateTask(todoID int) error {
	db := utils.SQLAcc.GetSQLDB()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO task (name, dateC, dateF, priority, status, ToDoID) VALUES(?,?,?,?,?,?)", ts.Name, ts.DateCreated, ts.DateFinish, ts.Priority, ts.Status, todoID)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

//DeleteToDo ...
func (td *ToDo) DeleteToDo(userID int, todoID int) error {
	db := utils.SQLAcc.GetSQLDB()

	if userID != 0 {
		_, err := db.Exec("DELETE FROM ToDo WHERE id=? AND userID=?", todoID, userID)
		if err != nil {
			return err
		}
	} else {
		_, err := db.Exec("DELETE FROM ToDo WHERE id=?", todoID)
		if err != nil {
			return err
		}
	}
	return nil
}

//DeleteTask ...
func (ts *Task) DeleteTask(todoID int, taskID int) error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec("DELETE FROM task WHERE id=? AND ToDoID=?", taskID, todoID)
	if err != nil {
		return err
	}

	return nil
}

//UpdateToDoName ...
func (td *ToDo) UpdateToDoName(todoID int) error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec("UPDATE ToDo SET name=? where id=?", td.Name, todoID)

	if err != nil {
		return err
	}

	return nil
}

//UpdateToDoDescription ...
func (td *ToDo) UpdateToDoDescription(todoID int) error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec("UPDATE ToDo SET description=? where id=?", td.Description, todoID)

	if err != nil {
		return err
	}

	return nil
}

//UpdateTaskName ...
func (ts *Task) UpdateTaskName(todoID int, taskID int) error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec("UPDATE task SET name=? WHERE id=? AND ToDoID=?", ts.Name, taskID, todoID)

	if err != nil {
		return err
	}

	return nil
}

//UpdateTaskDateStart ...
func (ts Task) UpdateTaskDateStart(todoID int, taskID int) error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec("UPDATE task SET dateC=? WHERE id=? AND ToDoID=?", ts.DateCreated, taskID, todoID)
	if err != nil {
		return err
	}

	return nil
}

//UpdateTaskDateFinish ...
func (ts *Task) UpdateTaskDateFinish(todoID int, taskID int) error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec("UPDATE task SET dateF=? WHERE id=? AND ToDoID=?", ts.DateFinish, taskID, todoID)
	if err != nil {
		return err
	}

	return nil
}

//UpdateTaskPriority ...
func (ts *Task) UpdateTaskPriority(todoID int, taskID int) error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec("UPDATE task SET priority=? WHERE id=? AND ToDoID=?", ts.Priority, taskID, todoID)
	if err != nil {
		return err
	}

	return nil
}

//UpdateTaskStatus ...
func (ts *Task) UpdateTaskStatus(todoID int, taskID int) error {
	db := utils.SQLAcc.GetSQLDB()

	_, err := db.Exec("UPDATE task SET status=? WHERE id=? AND ToDoID=?", ts.Status, taskID, todoID)
	if err != nil {
		return err
	}

	return nil
}

//ListAllToDos (admin only)...
func ListAllToDos(userID int) ([]ToDo, error) {
	db := utils.SQLAcc.GetSQLDB()

	var todos []ToDo
	var err error

	if userID == 0 {
		err = db.Select(&todos, "SELECT * FROM ToDo")

		if err != nil {
			fmt.Println("Cannot show all created ToDos")
			return nil, err
		}
	} else {
		err = db.Select(&todos, "SELECT * FROM ToDo WHERE userID=?", userID)

		if err != nil {
			fmt.Println("Cannot show all created ToDos")
			return nil, err
		}
	}

	return todos, err
}

//GetAnyToDo returns ToDo using ToDoID ...
func GetAnyToDo(todoID int) (ToDo, error) {
	db := utils.SQLAcc.GetSQLDB()

	var todo ToDo

	err := db.Get(&todo, "SELECT * FROM ToDo WHERE id=?", todoID)
	if err != nil {
		return todo, err
	}

	return todo, nil

}

//GetAnyTask returns Task using taskID ...
func GetAnyTask(taskID int) (Task, error) {
	db := utils.SQLAcc.GetSQLDB()

	var task Task

	err := db.Get(&task, "SELECT * FROM task WHERE id=?", taskID)
	if err != nil {
		return task, err
	}

	return task, nil

}

//ListTasks shows all tasks per a user ...
func ListTasks(tdid int) ([]Task, error) {
	db := utils.SQLAcc.GetSQLDB()

	var tasks []Task
	var err error

	err = db.Select(&tasks, "SELECT * FROM task WHERE ToDoID=?", tdid)

	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

//ListAllActiveTasks ...
func ListAllActiveTasks(tdid int) ([]Task, error) {
	db := utils.SQLAcc.GetSQLDB()

	var activeTasks []Task

	err := db.Select(&activeTasks, "SELECT * FROM task WHERE status=0 AND ToDoID=?", tdid)
	if err != nil {
		return nil, err
	}

	return activeTasks, err
}

//ListCompletedTasks ...
func ListCompletedTasks(tdid string) ([]Task, error) {
	db := utils.SQLAcc.GetSQLDB()

	var completedTasks []Task

	err := db.Select(&completedTasks, "SELECT * FROM task WHERE status=1 AND ToDoID=?", tdid)
	if err != nil {
		return nil, err
	}

	return completedTasks, err
}

package tasksBD 

import (
	//"fmt"
	"log"
	"database/sql"
	_ "modernc.org/sqlite"
)

type Task struct {
	Id int 				`json:"id"`
	Title string		`json:"title"`
	Subject string 		`json:"subject"`
	Topic string		`json:"topic"`
	Short string		`json:"short"`
	Content string		`json:"content"`
	Instruction string	`json:"instruction"`
	Image string		`json:"image"`
	Answer_type string	`json:"answer_type"`
	Allow_photo bool	`json:"allow_photo"`
	Status int 			`json:"status"`
}

type Tasklist []Task 

type TasksDBConnection struct{
	db *sql.DB
}

func NewTasksDBConnection(db *sql.DB) *TasksDBConnection{
	return &TasksDBConnection{db: db}
}

func (connection *TasksDBConnection) FindTasks(subject string) (*Tasklist, error){
	var list Tasklist
	log.Println(subject)
	rows, err := connection.db.Query("SELECT * FROM tasks WHERE subject = ?", subject)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next(){
    	var t Task
    	if  err := rows.Scan(&t.Id, &t.Title, &t.Subject, &t.Topic, &t.Short, &t.Content, &t.Instruction, &t.Image, &t.Answer_type, &t.Status); err !=nil {
    		return nil, err
    	}
    	list = append(list, t)
    }
	log.Println(&list)
	return &list, nil
}

func (connection *TasksDBConnection) FindTasksById(id []int) (*Tasklist, error){
	var list Tasklist
	for _, value := range id {
		t, err:= connection.getTaskById(value)
		if t==nil && err ==nil{
			log.Println("task not found")
			continue
		}
		if err != nil {
			return nil, err
		}
		list = append(list, *t)
	}
	return &list, nil
}

func (connection *TasksDBConnection) getTaskById(id int) (*Task, error){
	var task Task
	err := connection.db.QueryRow(
        "SELECT * FROM tasks WHERE id = ?", id,
    ).Scan(&task.Id, &task.Title, &task.Subject, &task.Topic, &task.Short, &task.Content, &task.Instruction, &task.Image, &task.Answer_type, &task.Status)
    if err == sql.ErrNoRows {
        return nil, nil // пользователь не найден
    }
    if err != nil {
        return nil, err
    }
    return &task, nil
}
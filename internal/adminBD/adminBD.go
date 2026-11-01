package adminBD 

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

type AdminDBConnection struct{
	db *sql.DB
}

func NewAdminDBConnection(db *sql.DB) *AdminDBConnection{
	return &AdminDBConnection{db: db}
}

func (connection *AdminDBConnection) TableExists(tableName string) (bool, error) {
    var count int
    err := connection.db.QueryRow(
        "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?",
        tableName,
    ).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (connection *AdminDBConnection) CreateTable() (bool, error){
	existed, _:= connection.TableExists("tasks")

	if existed {
    	return false, nil
	} else {
    	log.Println("creating table tasks...")
	}
	_, err := connection.db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            Id INTEGER PRIMARY KEY AUTOINCREMENT,
            Title TEXT NOT NULL,
            Subject TEXT NOT NULL,
            Topic TEXT NOT NULL,
            Short TEXT NOT NULL,
            Content TEXT UNIQUE NOT NULL,
            Instruction TEXT NOT NULL,
            Image TEXT NOT NULL,
            Answer_type TEXT NOT NULL,
            Status INTEGER NOT NULL
        )
    `)
    if err != nil {
        log.Fatal(err)
        return false, err
    }
    return true, nil
}

func (connection *AdminDBConnection) AddTasks(list *Tasklist) (bool, error){
	for _, value := range *list {
		_, err:= connection.db.Exec("INSERT OR IGNORE INTO tasks(title, subject, topic, short, content, instruction, image, answer_type, status) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)", &value.Title, &value.Subject, &value.Topic, &value.Short, &value.Content, &value.Instruction, &value.Image, &value.Answer_type, &value.Status)
		if err != nil {
			log.Println(err)
		}
	}
	return true, nil
}

func (connection *AdminDBConnection) FindAllTasks() (*Tasklist, error){
	var list Tasklist
	rows, err := connection.db.Query("SELECT * FROM tasks")
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

func (connection *AdminDBConnection) DeleteTask(id int) (bool, error) {
	_, err := connection.db.Exec("DELETE FROM tasks WHERE Id = ?", id)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}


// func (connection *TasksDBConnection) getTaskById(id int) (*Task, error){
// 	var task Task
// 	err := connection.db.QueryRow(
//         "SELECT * FROM tasks WHERE id = ?", id,
//     ).Scan(&task.Id, &task.Title, &task.Subject, &task.Topic, &task.Short, &task.Content, &task.Instruction, &task.Image, &task.Answer_type, &task.Allow_photo)
//     if err == sql.ErrNoRows {
//         return nil, nil // пользователь не найден
//     }
//     if err != nil {
//         return nil, err
//     }
//     return &task, nil
// }
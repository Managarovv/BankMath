package main

import (
	"fmt"
	"net/http"
	//"time"
	"log"
	"encoding/json"
	"BankMath/internal/tasksBD"
	"BankMath/internal/adminBD"
	"database/sql"
	_ "modernc.org/sqlite"
	"strconv"   
	"strings"
)

type TaskHandler struct {
	// Сюда сохраняем подключение к вашему модулю
	// (тип должен совпадать с тем, что возвращает NewTasksDBConnection)
	tasksConn *tasksBD.TasksDBConnection 
	adminConn *adminBD.AdminDBConnection
}

func (h *TaskHandler) taskslist (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//log.Println(r.URL.Query().Get("subject"))
	if  r.URL.Query().Get("subject") == "" {
		http.Error(w, "incorrect fields", http.StatusBadRequest)
		return
	} else {
		ft, err := h.tasksConn.FindTasks(r.URL.Query().Get("subject"))
		if err != nil {
			http.NotFound(w, r)
			return 
		} else {
			log.Println(*ft)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*ft)			
		}
	}
}

func (h *TaskHandler) taskslistByID(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем query-параметр
	rawIds := r.URL.Query().Get("ids") // Ожидаем: ?ids=1,5,12,42
	if rawIds == "" {
		http.Error(w, "Missing 'ids' parameter", http.StatusBadRequest)
		return
	}

	// 2. Парсим ID
	idStrings := strings.Split(rawIds, ",")
	ids := make([]int, 0, len(idStrings)) // Оптимизируем выделение памяти

	for _, s := range idStrings {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			continue // Игнорируем лишние запятые, например "1,,2"
		}
		id, err := strconv.Atoi(trimmed)
		if err != nil {
			http.Error(w, "Invalid ID format: "+trimmed, http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}

	// Проверяем, набрались ли валидные ID после фильтрации пустых строк
	if len(ids) == 0 {
		http.Error(w, "Missing 'ids' parameter", http.StatusBadRequest)
		return
	}

	// 3. Запрос к БД / Сервису
	ftid, err := h.tasksConn.FindTasksById(ids)
	if err != nil {
		log.Printf("Error finding tasks by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError) // Ошибка сервера, а не 404
		return
	}

	// Если ничего не найдено, возвращаем пустой слайс вместо null или 404
	if ftid == nil {
		//ftid = []Task{} // Или другой ваш тип данных, инициализируем пустой список
	}

	log.Println(ftid)

	// 4. Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Явно указываем статус 200
	if err := json.NewEncoder(w).Encode(ftid); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}


/*func (h *TaskHandler) taskslistByID(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	rawIds := r.URL.Query().Get("ids") // должно прийти: ?ids=1,5,12,42
	if rawIds == "" {
    	http.Error(w, "Missing 'ids' parameter", http.StatusBadRequest)
    	return
	}
	idStrings := strings.Split(rawIds, ",")
	var ids []int
	for _, s := range idStrings {
    	id, err := strconv.Atoi(strings.TrimSpace(s))
    	if err != nil {
        	http.Error(w, "Invalid ID format", http.StatusBadRequest)
        	return
    	}
    	ids = append(ids, id)
	}
	if len(ids) == 0{
		http.Error(w, "Missing 'ids' parameter", http.StatusBadRequest)
    	return
	}
	ftid, err:= h.tasksConn.FindTasksById(ids)
	if err != nil {
		http.NotFound(w, r)
		return 
	}
	if ftid == nil {
		http.NotFound(w, r) // Отдаст 404 Not Found, так как задачи не найдены
		return
	} 
		log.Println(*ftid)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(*ftid)
	
}*/

func (h *TaskHandler) addNewTasks(w http.ResponseWriter, r *http.Request){
	var list adminBD.Tasklist
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }
	addT, err := h.adminConn.AddTasks(&list)
	if err != nil {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return 
	}
	if addT == false {
		fmt.Fprintf(w, "thats tasks already exist in db")
		return
	} else{
		fmt.Fprintf(w, "tasks added in db")
		return
	}
}

func (h *TaskHandler) findAll(w http.ResponseWriter, r *http.Request){
	ft, err := h.adminConn.FindAllTasks()
	if err != nil {
		log.Printf("Error finding tasks: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError) // Ошибка сервера, 500
		return
	}
	if ft == nil {
		log.Printf("No one tasks find")
		http.Error(w, "NotFound", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Явно указываем статус 200
	if err := json.NewEncoder(w).Encode(ft); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func (h *TaskHandler) deleteById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "incorrect fields", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	success, err := h.adminConn.DeleteTask(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !success {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Task was deleted")
}


func main() {
	db, err := sql.Open("sqlite", "./sqlite/Bank.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    log.Println("Подключено к SQLite")

	mux := http.NewServeMux()
	muxadmin := http.NewServeMux()

	tasksConnection:= tasksBD.NewTasksDBConnection(db)
	adminConnection:= adminBD.NewAdminDBConnection(db)
	handle := &TaskHandler{
		tasksConn: tasksConnection,
		adminConn: adminConnection,
	}

	tasks := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /tasks/", http.StripPrefix("/tasks/", tasks))

	mux.HandleFunc("GET /taskslist/", handle.taskslist)
	mux.HandleFunc("GET /taskslistbyid", handle.taskslistByID)

	//mux.HandleFunc("/student", editStudent)

	muxadmin.HandleFunc("GET /admin/search", handle.taskslistByID)
	muxadmin.HandleFunc("GET /admin/search/all", handle.findAll)
	muxadmin.HandleFunc("GET /admin/", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "static/admin/admin.html")
	})
	muxadmin.HandleFunc("GET /admin/createdb", func(w http.ResponseWriter, r *http.Request){
		res, err:= handle.adminConn.CreateTable()
		if err != nil {
			http.Error(w, "somethibg went wrong", http.StatusInternalServerError)
			return 
		}
		if res {
			fmt.Fprintf(w, "database succesfull created")
		} else{
			http.Error(w, "database already exist", http.StatusConflict)
		}
	}) 
	muxadmin.HandleFunc("POST /admin", handle.addNewTasks)
	muxadmin.HandleFunc("DELETE /admin/{id}", handle.deleteById)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		http.NotFound(w, r)
	})

	go func() {
		log.Println("Админка слушает на http://localhost:8081 (скрыта от внешнего мира)")
		if err := http.ListenAndServe(":8081", muxadmin); err != nil {
			fmt.Printf("Ошибка запуска админки: %v\n", err)
		}
	}()
	

	server := &http.Server{
		Addr: ":80",
		Handler: mux,
	}
	log.Println("Server listening")
	server.ListenAndServe()
}
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта

// Обработчик для получения всех задач -> метод GET
func getTasks(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GET new request!")

	// сериализуем данные из слайса tasks
	//resp, err := json.MarshalIndent(tasks, "", "    ")
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // При ошибке сервер должен вернуть статус 500 Internal Server Error
		return
	}

	// в заголовок записываем тип контента, данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// возвращаем статус 200 OK об успешном выполнении запроса
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные данные в тело ответа
	w.Write(resp)
}

// Обработчик для отправки(добавления) новой задачи сервер -> метод POST
func postTask(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("POST new request!")

	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// добавляем новую задачу в мапу
	tasks[task.ID] = task
	// указываем тип контента
	w.Header().Set("Content-Type", "application/json")
	// При успешном запросе сервер должен вернуть статус 201 Created
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTaskId(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GET new request ID!")
	// получаем значение параметра по ключу id из запроса
	id := chi.URLParam(r, "id")
	// проверяем ключ в мапе
	taskId, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена!", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик удаления задачи по ID
func deleteTaskId(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("DELETE new request!")

	id := chi.URLParam(r, "id")
	// если id с такой задачей нет
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Такой задачи не существует!", http.StatusBadRequest)
		return
	}
	// если задача есть, удаляем из мапы
	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики

	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks", getTasks)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTask`
	r.Post("/tasks", postTask)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTaskId`
	r.Get("/tasks/{id}", getTaskId)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `deleteTaskId`
	r.Delete("/tasks/{id}", deleteTaskId)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

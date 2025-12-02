package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pepecloud/go-todo-cli/internal/model"

	"github.com/k0kubun/pp/v3"
)

// 1. Main - точка входа. Создаем мапу и сканер тут.
func main() {
	// Создаем хранилище
	todoList := make(map[string]model.Task)
	history := make([]model.Event, 0)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Приложение запущено!")

	// Тут будет твой главный цикл обработки команд
	for {
		fmt.Println("Введите команду:")

		if scanner.Scan() {
			text := scanner.Text()
			LogEvent(&history, text, "")
			switch text {
			case "add":
				CreateTask(todoList, scanner, &history)
			case "list":
				pp.Println("Ваши задачи:", todoList)
			case "help":
				fmt.Printf("Команды для пользования ботом:\n add - создание новой задачи\n list - созданные задачи.\n del - удалить задачу по загаловку\n done - пометить задачу как выполненная\n exit - остановит программу\n")
			case "del":
				DelTask(todoList, scanner, &history)
			case "done":
				DoneTask(todoList, scanner, &history)
			case "exit":
				return
			case "events":
				pp.Println(history)
			default:
				LogEvent(&history, text, "Я вас не понимаю.")
				fmt.Println("Я вас не понимаю.")
			}
		}
	}
}

// 2. Функция вынесена ЗА пределы main
func CreateTask(todoList map[string]model.Task, sc *bufio.Scanner, history *[]model.Event) {
	fmt.Print("Введите заголовок задачи (ОДНО СЛОВО): ")
	if !sc.Scan() {
		return
	}
	title := sc.Text()
	LogEvent(history, title, "")

	// Проверка на одно слово
	words := strings.Fields(title) // Fields надежнее Split
	if len(words) != 1 {
		fmt.Println("Ошибка: Заголовок должен быть одним словом!")
		LogEvent(history, title, "Ошибка: Заголовок должен быть одним словом!")
		return
	}

	// Подтверждение
	if !Confirmation(sc, history) {
		return // Если юзер сказал 'n', выходим из функции создания
	}

	fmt.Print("Введите текст задачи: ")
	if !sc.Scan() {
		return
	}
	text := sc.Text()
	LogEvent(history, text, "")

	// Снова подтверждение
	if !Confirmation(sc, history) {
		return
	}

	// Создаем и сохраняем
	myNewTask := model.NewTask(title, text)
	todoList[title] = myNewTask

	fmt.Println("Задача сохранена! Всего задач:", len(todoList))
}

// Удаление задачи по загаловку
func DelTask(todoList map[string]model.Task, sc *bufio.Scanner, history *[]model.Event) {
	fmt.Println("Вы хотите удалить задачу по заггаловку!")
	fmt.Println("Введите заголовок:")

	if !sc.Scan() {
		return
	}
	title := sc.Text()
	LogEvent(history, title, "")

	// Проверка на одно слово
	words := strings.Fields(title) // Fields надежнее Split
	if len(words) != 1 {
		fmt.Println("Ошибка: Заголовок должен быть одним словом!")
		LogEvent(history, title, "Ошибка: Заголовок должен быть одним словом!")
		return
	}

	_, ok := todoList[title]

	if !ok {
		fmt.Println("Такой задачи нету!")
		LogEvent(history, title, "Такой задачи нету!")
		return
	}

	fmt.Printf("Вы уверены, что хотите удалить '%s'? ", title)
	// Подтверждение
	if !Confirmation(sc, history) {
		return // Если юзер сказал 'n', выходим из функции создания
	}

	delete(todoList, title)
	fmt.Println("Задача удалена.")
	LogEvent(history, title, "Задача удалена.")
}

func DoneTask(todoList map[string]model.Task, sc *bufio.Scanner, history *[]model.Event) bool {
	fmt.Println("Вы хотите пометить задание как выполненное!")
	fmt.Println("Введите заголовок:")

	if !sc.Scan() {
		return false
	}
	title := sc.Text()
	LogEvent(history, title, "")

	// Проверка на одно слово
	words := strings.Fields(title) // Fields надежнее Split
	if len(words) != 1 {
		fmt.Println("Ошибка: Заголовок должен быть одним словом!")
		LogEvent(history, title, "Ошибка: Заголовок должен быть одним словом!")
		return false
	}

	_, ok := todoList[title]

	if !ok {
		fmt.Println("Такой задачи нету!")
		LogEvent(history, title, "Такой задачи нету!")
	}

	fmt.Printf("Вы хотите пометить заголовок %s как выполненный?\n", title)

	if !Confirmation(sc, history) {
		return false
	}

	if task, ok := todoList[title]; ok {
		task.Complite = true
		todoList[title] = task
	}

	fmt.Println("Задача выполнена!")
	LogEvent(history, title, "Задача выполнена!")
	return true
}

func LogEvent(history *[]model.Event, input string, errText string) {
	newEvent := model.Event{
		RawInput:  input,
		ErrorText: errText,
		CreatedAt: time.Now(), // Не забудь импортировать "time"
	}

	// *history означает "взять значение по этому адресу"
	*history = append(*history, newEvent)
}

// 3. Функция Confirmation возвращает bool (true/false), а не просто печатает
func Confirmation(sc *bufio.Scanner, history *[]model.Event) bool {
	fmt.Print("Продолжить? (y/n): ")
	if !sc.Scan() {
		return false
	}

	input := strings.ToLower(strings.TrimSpace(sc.Text()))
	LogEvent(history, input, "")

	switch input {
	case "y":
		fmt.Println("Ок, продолжаем.")
		LogEvent(history, input, "Ввел y")
		return true
	case "n":
		fmt.Println("Отмена операции.")
		LogEvent(history, input, "Ввел n")
		return false
	default:
		fmt.Println("Непонятно, считаем за 'нет'.")
		LogEvent(history, input, "")
		return false
	}
}

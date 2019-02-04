package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/codegangsta/cli"
)

// Task contains content of task and some information about tasks
type Task struct {
	Content string
	Time    time.Time
}

func main() {
	app := cli.NewApp()

	app.Name = "Todo"
	app.Usage = "add, list and complete tasks"

	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "add task",
			Action: func(c *cli.Context) {
				text := Task{Content: c.Args().First(), Time: time.Now()}

				if err := AddTask(text, "tasks_data"); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:      "complete",
			ShortName: "comp",
			Usage:     "complete task",
			Action: func(c *cli.Context) {
				posOfTask, err := strconv.Atoi(c.Args().First())

				if err != nil {
					fmt.Println("You shoudld write integer")
				}

				if _, err := CompleteTask(posOfTask); err != nil {
					fmt.Println(err)
				}
			},
		},
		{
			Name:      "list",
			ShortName: "ls",
			Usage:     "list all uncompleted tasks",
			Action: func(c *cli.Context) {
				data, err := ListTasks()

				if err != nil {
					log.Fatal(err)
				}

				fmt.Print(data)
			},
		},
		{
			Name:      "clear all",
			ShortName: "clall",
			Usage:     "clear all data from data file",
			Action: func(c *cli.Context) {
				err := ClearData()

				if err != nil {
					log.Fatal(err)
				}
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

// AddTask add task in JSON file
func AddTask(task Task, dataFileName string) error {
	j, err := json.Marshal(task)

	if err != nil {
		return fmt.Errorf("problems with marshaling file")
	}

	j = append(j, "\n"...)

	file, _ := os.OpenFile(dataFileName+".json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if _, err := file.Write(j); err != nil {
		return fmt.Errorf("problems with opening file")
	}

	return nil
}

// ListTasks list all uncompleted tasks from tasks_data file
func ListTasks() (string, error) {
	var result string

	file, err := OpenAndCheckFile("tasks_data.json")

	if err != nil {
		return "List of tasks is empty\n", nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	idx := 1

	for scanner.Scan() {
		text := scanner.Text()
		task := Task{}

		err := json.Unmarshal([]byte(text), &task)

		if err != nil {
			return "", fmt.Errorf("problems with unmarshaling file")
		}

		result += CreateTask(idx, task)

		idx++
	}

	if idx == 1 {
		return "List of tasks is empty\n", nil
	}

	return result, nil
}

// CompleteTask completed task with posision equals posOfTask
func CompleteTask(posOfTask int) (string, error) {
	file, err := OpenAndCheckFile("tasks_data.json")

	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	idx := 1

	flag := 0

	for scanner.Scan() {
		text := scanner.Text()
		task := Task{}

		err := json.Unmarshal([]byte(text), &task)

		if err != nil {
			return "", fmt.Errorf("problems with unmarshaling data file")
		}

		if idx == posOfTask {
			idx++
			flag = 1
			continue
		}

		idx++

		AddTask(task, "tmp")
	}

	Swap("tmp.json", "tasks_data.json")

	if flag == 0 {
		return fmt.Sprintf("There're less than %d tasks\n", idx), nil
	}

	return "", nil
}

// OpenAndCheckFile opens file
func OpenAndCheckFile(fileName string) (*os.File, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("file with data now exists")
	}

	file, err := os.Open(fileName)

	if err != nil {
		return nil, fmt.Errorf("problems with opening data file")
	}

	return file, nil
}

// ClearData rewrite data file
func ClearData() error {
	_, err := os.OpenFile("tmp.json", os.O_CREATE, 0600)

	if err != nil {
		return fmt.Errorf("problem with creating new data file")
	}

	Swap("tmp.json", "tasks_data.json")

	return nil
}

// Swap swaps two files
func Swap(file1 string, file2 string) {
	os.Rename(file1, file2)

	os.Remove(file1)
}

// CreateTask prints task in nice format
func CreateTask(idx int, task Task) (res string) {
	res = fmt.Sprintf("[%d] ", idx) +
		fmt.Sprintf("(%02d:%02d:%02d ", task.Time.Hour(), task.Time.Minute(), task.Time.Second()) +
		fmt.Sprintf("%02d-%02d-%d) ", task.Time.Day(), task.Time.Month(), task.Time.Year()) +
		fmt.Sprintln(task.Content)

	return
}

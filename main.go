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

				AddTask(text, "tasks_data")
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

				if CompleteTask(posOfTask) != nil {
					fmt.Println(err)
				}
			},
		},
		{
			Name:      "list",
			ShortName: "ls",
			Usage:     "list all uncompleted tasks",
			Action: func(c *cli.Context) {
				err := ListTasks()

				if err != nil {
					log.Fatal(err)
				}
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
func AddTask(task Task, dataFileName string) {
	j, err := json.Marshal(task)

	if err != nil {
		log.Fatal(err)
	}

	j = append(j, "\n"...)

	file, _ := os.OpenFile(dataFileName+".json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if _, err := file.Write(j); err != nil {
		log.Fatal(err)
	}
}

// ListTasks list all uncompleted tasks from tasks_data file
func ListTasks() error {
	file, err := OpenAndCheckFile("tasks_data.json")

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	idx := 1

	for scanner.Scan() {
		text := scanner.Text()
		task := Task{}

		err := json.Unmarshal([]byte(text), &task)

		if err != nil {
			return err
		}

		PrintTask(idx, task)

		idx++
	}

	if idx == 1 {
		fmt.Println("List of tasks is empty")
	}

	return nil
}

// CompleteTask completed task with posision equals posOfTask
func CompleteTask(posOfTask int) error {
	file, err := OpenAndCheckFile("tasks_data.json")

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	idx := 1

	for scanner.Scan() {
		text := scanner.Text()
		task := Task{}

		err := json.Unmarshal([]byte(text), &task)

		if err != nil {
			return fmt.Errorf("problems with unmarshaling data file")
		}

		if idx == posOfTask {
			idx++
			continue
		}

		idx++

		AddTask(task, "tmp")
	}

	Swap("tmp.json", "tasks_data.json")

	return nil
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

// PrintTask prints task in nice format
func PrintTask(idx int, task Task) {
	fmt.Printf("[%d] ", idx)
	fmt.Printf("(%02d:%02d:%02d ", task.Time.Hour(), task.Time.Minute(), task.Time.Second())
	fmt.Printf("%02d-%02d-%d) ", task.Time.Day(), task.Time.Month(), task.Time.Year())
	fmt.Println(task.Content)
}

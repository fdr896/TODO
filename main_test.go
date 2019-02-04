package main

import (
	"fmt"
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	ClearData()

	task := Task{Content: "Hello", Time: time.Now()}

	if err := AddTask(task, "tasks_data"); err != nil {
		t.Error("AddTask doesn't work correctly", err)
	}

	res, err := ListTasks()

	if err != nil {
		t.Error("ListTask doesn't work correctly", err)
	}

	needRes := CreateTask(1, task)

	if res != needRes {
		t.Error("output doesn't match needRes: \nexpected:" + needRes + "have: " + res)
	}
}

func TestListTaskWithEmptyContent(t *testing.T) {
	ClearData()

	res, err := ListTasks()

	if err != nil {
		t.Error("ListTask doesn't work with empty content")
	}

	needRes := "List of tasks is empty\n"

	if res != needRes {
		t.Error("output doesn't mathch ndRes: \nexpected: " + needRes + "have: " + res)
	}
}

func TestClearData(t *testing.T) {
	if err := ClearData(); err != nil {
		t.Error("ClearData doesn't work correctly")
	}

	res, _ := ListTasks()

	needRes := "List of tasks is empty\n"

	if res != needRes {
		t.Error("ClearData doesn't clear all data:\n expected: ", needRes, "have: ", res)
	}
}

func TestCompleteTask(t *testing.T) {
	ClearData()

	name := "tasks_data"

	task1 := Task{Content: "Hello", Time: time.Now()}
	task2 := Task{Content: "Goodbye", Time: time.Now()}
	task3 := Task{Content: "Mazafaka", Time: time.Now()}

	AddTask(task1, name)
	AddTask(task2, name)
	AddTask(task3, name)

	CompleteTask(2)

	res, _ := ListTasks()

	needRes := CreateTask(1, task1) + CreateTask(2, task3)

	if res != needRes {
		t.Error("CompleteTask doesn't work correctly\n expected: ", needRes, "have: ", res)
	}

	ClearData()
}

func TestCompleteTaskWithFewTasks(t *testing.T) {
	ClearData()

	name := "tasks_data"

	task1 := Task{Content: "Hello", Time: time.Now()}
	task2 := Task{Content: "Goodbye", Time: time.Now()}
	task3 := Task{Content: "Mazafaka", Time: time.Now()}

	AddTask(task1, name)
	AddTask(task2, name)
	AddTask(task3, name)

	res, _ := CompleteTask(4)

	needRes := fmt.Sprintf("There're less than %d tasks\n", 4)

	if res != needRes {
		t.Errorf("CompleteTask with few tasks doesn't work correctly\n expected: %s have: %s", needRes, res)
	}
}

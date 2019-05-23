package main

import "time"

type Task struct {
	Autor int64
	//Assigned int64
	DueTime time.Time
	Message string
}

func (t *Task) Create() {

}

func (t *Task) Update(id int64) {

}

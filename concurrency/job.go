package concurrency

import (
	//"time"
)

type Job struct {
    ID        int
    //CreatedAt time.Time
    //UpdatedAt time.Time
    Input interface{}
}

type JobChannel chan Job
type JobQueue chan chan Job
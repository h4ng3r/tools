package concurrency

import (
    "sync"
)

type Worker struct {
    ID      int
    JobChan JobChannel
    Queue   JobQueue   // shared between all workers and dispatchers.
    Quit    chan struct{}
}

func (wr *Worker) Start(wg *sync.WaitGroup, jobFunc JobFuncType) {
    go func() {
        for {
            wr.Queue <- wr.JobChan
            select {
            case job := <-wr.JobChan:
                // TODO: status struct
            	jobFunc(job.Input, wr.ID, job.ID)
                wg.Done()
            case <-wr.Quit:
                close(wr.JobChan)
                return
            }
        }
    }() 
}
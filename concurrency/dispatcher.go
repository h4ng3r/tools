package concurrency

import (
    "sync"
)


type JobFuncType func(interface{}, int, int) //interface{}

type disp struct {
    Workers  []*Worker  // this is the list of workers that dispatcher tracks
    WorkChan JobChannel // client submits a job to this channel
    Queue    JobQueue   // this is the shared JobPool between the workers
    Jobs int
    wg sync.WaitGroup   // wating group
    jobFunc JobFuncType
}

func NewDispather(num int, jobFunc JobFuncType) *disp {
    return &disp{
        Workers:  make([]*Worker, num),
        WorkChan: make(JobChannel),
        Queue:    make(JobQueue),
        //wg:		  make(sync.WaitGroup),
        Jobs: 0,
        jobFunc: jobFunc,
    }
}

func (d *disp) Start() *disp {
    l := len(d.Workers)
    //wss := reflect.ValueOf(ws)
    for i := 0; i < l; i++ {
        wrk := Worker{i, make(JobChannel), d.Queue, make(chan struct{})}
        wrk.Start(&d.wg, d.jobFunc)
        d.Workers[i] = &wrk
    }
    go d.process()
    return d
}

func (d *disp) process() {
    for {
        select {
        case job := <-d.WorkChan: // listen to a submitted job on WorkChannel
            jobChan := <-d.Queue  // pull out an available jobchannel from queue
            jobChan <- job        // submit the job on the available jobchannel
        }
    }
}

func (d *disp) Wait() {
	d.wg.Wait()
}	

// TODO: add generator but handle the count
func (d *disp) Submit(jobInput interface{}) {
	d.wg.Add(1)
	d.Jobs += 1
    d.WorkChan <- Job{ID:d.Jobs, Input:jobInput} // TODO: FIX
}
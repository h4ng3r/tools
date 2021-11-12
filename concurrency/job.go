package concurrency

type JobChannel chan interface{}
type JobQueue chan chan interface{}

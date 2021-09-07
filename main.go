package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var letterRunes = []rune("abcdefg")
var minPriority = 1
var maxPriority = 5

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandPriority() int {
	return rand.Intn(maxPriority-minPriority+1) + minPriority
}

type Task struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []Task

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	task := x.(Task)
	*pq = append(*pq, task)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
func executeTask(tasks <-chan Task) {
	for task := range tasks {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("output:- %.2d:%s \n", task.priority, task.value)
	}
}
func main() {
	pq := make(PriorityQueue, 1)
	/*I have added this initial task because heap.Push operation
	        requries some initial task in order that Less function is
			invoked without any panic
	   TODO: Need to remove this initial task*/
	task := Task{
		value:    RandStringRunes(4),
		priority: RandPriority(),
	}
	fmt.Printf("input:- %.2d:%s \n", task.priority, task.value)
	pq[0] = task
	heap.Init(&pq)

	var wg sync.WaitGroup
	wg.Add(1)

	jobs := make(chan Task)

	go executeTask(jobs)
	go executeTask(jobs)

	/*Below go function is meant to generate the tasks with random priority*/
	go func() {
		for i := 1; i <= 100; i++ {
			// Insert a new item and then modify its priority.
			task := Task{
				value:    RandStringRunes(4),
				priority: RandPriority(),
			}

			fmt.Printf("input:- %.2d:%s \n", task.priority, task.value)
			heap.Push(&pq, task)

			// Popping will lead to removal of task as per the priority
			jobs <- heap.Pop(&pq).(Task)
		}
		wg.Done()
	}()

	wg.Wait()
}

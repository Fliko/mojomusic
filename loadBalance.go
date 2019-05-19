package main

import "container/heap"
// Request defines the work to do and 
type Request struct {
	fn func() int
	c chan int
}

// Worker defines server connections
type Worker struct {
	requests chan Request	// work to do
	pending int				// current running tasks
	index int				// where to find it
}

type Pool []*Worker

type Balancer struct {
	pool Pool
	done chan *Worker
}

func (b *Balancer) balance(work chan Request) {
	for {
		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

// Make Pool a heap interface that returns the Worker with the smallest # of pending tasks
func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}
func (p Pool) Len() int { return len(p)}
func (p Pool) Swap(i,j int) {
	p[i].pending, p[j].pending = p[j].pending, p[i].pending
}
func (p *Pool) Pop() interface{} {
    old := *p
    n := len(old)
    x := old[n-1]
    *p = old[0 : n-1]
    return x
}
func (p *Pool) Push(x interface{}) {
	*p = append(*p, x.(*Worker))
}


func (b *Balancer) dispatch(req Request) {
	w := heap.Pop(&b.pool).(*Worker)
	w.requests <- req
	w.pending++
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Worker) {
	w.pending--
	heap.Remove(&b.pool, w.index)
	heap.Push(&b.pool, w)
}

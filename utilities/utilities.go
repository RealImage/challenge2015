package utilities

import (
	"sync"
	"github.com/RealImage/Challenge/models"
)

type WorkRequest struct {
	Actor string
	Path  []models.Relationship
}

type Queue struct {
	Works []WorkRequest
	sync.Mutex
}

func (q *Queue) Enqueue(w ...WorkRequest) {
	q.Lock()
	defer q.Unlock()
	q.Works = append(q.Works, w...)
}

//TODO fix empty relationship
func (q *Queue) Dequeue() WorkRequest {
	q.Lock()
	defer q.Unlock()
	rel := q.Works[0]
	q.Works = q.Works[1:]
	return rel
}

func (q *Queue) Empty() bool {
	q.Lock()
	defer q.Unlock()
	return len(q.Works) == 0
}

func NewQueue() *Queue {
	return &Queue{}
}
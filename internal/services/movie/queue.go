package movie

type Queue struct {
	item_value []*PersonNetwork
}

func (q *Queue) Enqueue(item *PersonNetwork) {
	q.item_value = append(q.item_value, item) //used to add items
}

func (q *Queue) Dequeue() *PersonNetwork {
	if q.IsEmpty() {
		return nil
	}
	item := q.item_value[0]
	q.item_value = q.item_value[1:] //used to remove items
	return item
}

func (q *Queue) IsEmpty() bool {
	return len(q.item_value) == 0
}

func (q *Queue) Size() int {
	return len(q.item_value)
}

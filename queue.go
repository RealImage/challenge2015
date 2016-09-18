package main

type queue struct{
	value[]string 
}

func(q *queue)enqueue(name string){
	q.value = append(q.value, name)
}

func(q *queue)dequeue() string{	
	e := q.value[0]
	q.value = q.value[1:]
	return e
}

func(q *queue)length() int{
	return len(q.value)
}

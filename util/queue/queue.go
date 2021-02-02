package queue

type Queue struct {
	queue []int
}

func (q *Queue) Enqueue(v int) {
	q.queue = append(q.queue, v) // Enqueue
}
func (q *Queue) Dequeue() int {
	qu := q.queue[0]
	q.queue = q.queue[1:]   // Dequeue
	return qu
}

func (q *Queue) Size() int {
	return len(q.queue)
}
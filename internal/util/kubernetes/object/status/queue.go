package status

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sync"
)

const defaultBufferSize = 8192

type Queue struct {
	lock     sync.RWMutex
	channels map[string]chan event.GenericEvent
}

func NewQueue() *Queue {
	return &Queue{
		channels: make(map[string]chan event.GenericEvent),
	}
}
func (q *Queue) Publish(obj client.Object) {
	ch := q.getChanForKind(obj.GetObjectKind().GroupVersionKind())
	ch <- event.GenericEvent{Object: obj}
}

func (q *Queue) Subscribe(gvk schema.GroupVersionKind) chan event.GenericEvent {
	return q.getChanForKind(gvk)
}

func (q *Queue) getChanForKind(gvk schema.GroupVersionKind) chan event.GenericEvent {
	q.lock.Lock()
	defer q.lock.Unlock()
	ch, ok := q.channels[gvk.String()]
	if !ok {
		ch = make(chan event.GenericEvent, defaultBufferSize)
		q.channels[gvk.String()] = ch
	}
	return ch
}

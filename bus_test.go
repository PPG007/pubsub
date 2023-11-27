package pubsub

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testEvent struct {
	name string
	age  int
}

func TestSubscribe(t *testing.T) {
	b := New()
	count := 0
	b.Subscribe("test", func(args ...interface{}) {
		count++
	})
	b.Publish("test", testEvent{
		name: "PPG007",
		age:  23,
	})
	b.WaitUntilDone()
	assert.Equal(t, 1, count)
}

func TestUnsubscribe(t *testing.T) {
	b := New()
	count := 0
	id := b.Subscribe("test", func(args ...interface{}) {
		count++
	})
	b.Publish("test", testEvent{
		name: "PPG007",
		age:  23,
	})
	b.WaitUntilDone()
	assert.Equal(t, 1, count)
	b.Unsubscribe("test", id)
	b.Publish("test", testEvent{
		name: "PPG007",
		age:  23,
	})
	b.WaitUntilDone()
	assert.Equal(t, 1, count)
}

func TestArgs(t *testing.T) {
	b := New()
	b.Subscribe("test", func(args ...interface{}) {
		assert.Len(t, args, 1)
		assert.IsType(t, testEvent{}, args[0])
	})
	b.Publish("test", testEvent{
		name: "PPG007",
		age:  23,
	})
	b.WaitUntilDone()

}

package run

import "errors"

type memoryTape map[int]byte

func (t memoryTape) Get(ptr int) byte {
	v, ok := t[ptr]
	if !ok {
		return 0
	}
	return v
}

func (t memoryTape) Set(ptr int, val byte) {
	if val == 0 {
		delete(t, ptr)
	} else {
		t[ptr] = val
	}
}

func (t memoryTape) ApplyDelta(ptr, delta int) {
	t.Set(ptr, byte(int(t.Get(ptr))+delta))
}

var ErrorEmptyStack = errors.New("stack: stack is empty")

type stack struct {
	arr []int
}

func (c *stack) Push(name int) {
	c.arr = append(c.arr, name)
}

func (c *stack) Pop() (int, error) {
	len := len(c.arr)
	if len > 0 {
		v := c.arr[len-1]
		c.arr = c.arr[:len-1]
		return v, nil
	}
	return 0, ErrorEmptyStack
}

func (c *stack) Top() (int, error) {
	len := len(c.arr)
	if len > 0 {
		return c.arr[len-1], nil
	}
	return 0, ErrorEmptyStack
}

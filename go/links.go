package main

import (
	"errors"
	"fmt"
	"os"
)

type LinkQueue struct {
	queue []LinkTask
}

type LinkTask struct {
	Callback func()
}

//func (t SymlinkTask) Link() {
//}
//
//func (t HardlinkTask) Link() {
//	err := os.Link(t.Source, t.Destination)
//	check(err, fmt.Sprintf("creating hard link %s -> %s", t.Source, t.Destination))
//	t.Callback()
//}

func (q *LinkQueue) TrySymlink(source string, destination string, callback func()) {
	var task = LinkTask{Callback: func() {
		// This might be relative, cannot stat
		//if _, err := os.Stat(source); err != nil {
		//	panic(fmt.Errorf("statting %s before symlinking to it", source))
		//}
		err := os.Symlink(source, destination)
		check(err, fmt.Sprintf("creating symlink %s -> %s", destination, source))
		callback()
	}}
	if _, err := os.Stat(source); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			q.queue = append(q.queue, task)
			fmt.Printf("Enqueued symlink task because %s does not yet exist\n", source)
			return
		}
		panic(err)
	}

	task.Callback()
}

func (q *LinkQueue) TryHardlink(source string, destination string, callback func()) {
	var task = LinkTask{
		Callback: func() {
			if _, err := os.Stat(source); err != nil {
				panic(err)
			}
			err := os.Link(source, destination)
			check(err, fmt.Sprintf("creating hard link %s -> %s", destination, source))
			callback()
		},
	}
	if _, err := os.Stat(source); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			q.queue = append(
				q.queue,
				task,
			)
			fmt.Printf("Enqueued hard link task because %s does not yet exist\n", source)
			return
		}
		panic(err)
	}

	task.Callback()
}

package pararender

import (
	"html/template"

	"github.com/diamondburned/l4d2lb/pages/components/errbox"
)

type Task struct {
	rendered chan template.HTML
	renderer func(interface{}) template.HTML
}

// BroadcastError broadcasts errors to a list of tasks.
func BroadcastError(err error, tasks ...Task) {
	for _, task := range tasks {
		task.Send(nil, err)
	}
}

// EmptyTask creates an empty task.
func EmptyTask(renderer func(interface{}) template.HTML) Task {
	return Task{
		rendered: make(chan template.HTML),
		renderer: renderer,
	}
}

// Send sends the results to the task channel.
func (t Task) Send(v interface{}, err error) {
	if err != nil {
		t.rendered <- errbox.RenderHTML(err)
	} else {
		t.rendered <- t.renderer(v)
	}
}

// Render waits for the rendered HTML.
func (t Task) Render() template.HTML {
	return <-t.rendered
}

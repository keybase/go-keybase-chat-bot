package kbchat

import (
	"fmt"
	"log"
	"time"
)

func ErrToOK(err error) string {
	if err == nil {
		return "ok"
	}
	return "ERROR: " + err.Error()
}

type DebugOutput struct {
	name string
}

func NewDebugOutput(name string) *DebugOutput {
	return &DebugOutput{
		name: name,
	}
}

func (d *DebugOutput) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("%s: %s\n", d.name, msg)
}

func (d *DebugOutput) Trace(f func() error, format string, args ...interface{}) func() {
	msg := fmt.Sprintf(format, args...)
	start := time.Now()
	log.Printf("+ %s: %s\n", d.name, msg)
	return func() {
		log.Printf("- %s: %s -> %s [time=%v]\n", d.name, msg, ErrToOK(f()), time.Since(start))
	}
}

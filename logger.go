package flowchart

import (
	"fmt"
)

type FlowLogger interface {
	Log(msg string, ctx ...any)
	New(ctx ...any) FlowLogger
}

type FlowLoggerImpl struct {
	context map[any]any
}

func (f FlowLoggerImpl) New(ctx ...any) FlowLogger {
	if f.context == nil {
		f.context = map[any]any{}
	}
	newCtx := map[any]any{}
	for key, val := range f.context {
		newCtx[key] = val
	}
	for i := 0; i < len(ctx); i += 2 {
		key := ctx[i]
		val := ctx[i+1]
		newCtx[key] = val
	}
	newLogger := FlowLoggerImpl{newCtx}
	return newLogger
}

func (f FlowLoggerImpl) Log(msg string, ctx ...any) {
	fmt.Print(msg)
	for key, value := range f.context {
		fmt.Printf("; %v: %v", key, value)
	}
	fmt.Println()
}

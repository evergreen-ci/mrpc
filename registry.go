package mrpc

import (
	"sync"

	"github.com/evergreen-ci/mrpc/mongowire"
	"github.com/pkg/errors"
)

type OperationRegistry struct {
	ops map[mongowire.OpScope]HandlerFunc
	mu  sync.RWMutex
}

func (o *OperationRegistry) Add(op mongowire.OpScope, h HandlerFunc) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if err := op.Validate(); err != nil {
		return errors.Wrap(err, "invalid operation")
	}

	if h == nil {
		return errors.Errorf("cannot define nil handler function for operation %+v", op)
	}

	if _, ok := o.ops[op]; ok {
		return errors.Errorf("operation '%+v' is already defined", op)
	}

	o.ops[op] = h

	return nil
}

func (o *OperationRegistry) Get(scope *mongowire.OpScope) (HandlerFunc, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	scopeCopy := *scope
	if handler, ok := o.ops[scopeCopy]; ok {
		return handler, ok
	}

	// Default to using a handler without a context if there isn't a more
	// specific context match.
	scopeCopy.Context = ""
	handler, ok := o.ops[scopeCopy]
	return handler, ok
}

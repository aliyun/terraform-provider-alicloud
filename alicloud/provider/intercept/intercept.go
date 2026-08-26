// Package intercept provides the interceptor mechanism: Op, Call, Interceptor and
// two registries. An Interceptor sees the call only, never resource data, which is
// what lets one implementation work on both stacks.
package intercept

import (
	"context"
)

type Op string

const (
	OpCreate Op = "Create"
	OpRead   Op = "Read"
	OpUpdate Op = "Update"
	OpDelete Op = "Delete"
)

// Call describes one intercepted invocation, in fields both stacks have.
type Call struct {
	Name string // Terraform type name, e.g. "alicloud_vpc"
	Op   Op
	Meta interface{} // *connectivity.AliyunClient on both stacks
}

// Interceptor wraps one invocation: Before runs in forward order and aborts the
// operation on error, After always runs in reverse order and may reword or replace
// the error. It cannot make one disappear: Execute restores an error a hook
// returned nil for, since the operation would otherwise be reported as successful
// with whatever state the failure left behind.
type Interceptor interface {
	Before(ctx context.Context, call Call) error
	After(ctx context.Context, call Call, err error) error
}

// Literals with no setter: ChainOf snapshots chains while the provider server is
// being built, so adding an interceptor is an edit here.
var (
	registry = map[string][]Interceptor{}
	global   = []Interceptor{}
)

func Registry(name string) []Interceptor {
	return registry[name]
}

func Global() []Interceptor {
	return global
}

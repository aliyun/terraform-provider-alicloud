// One fixture per stack, each using as many optional capabilities as it offers.
package provider_test

import (
	"context"
	"slices"
	"sync"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
)

type hookLog struct {
	mu    sync.Mutex
	calls []string
}

func (l *hookLog) record(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.calls = append(l.calls, name)
}

func (l *hookLog) has(name string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return slices.Contains(l.calls, name)
}

func (l *hookLog) all() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.calls...)
}

func (l *hookLog) reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.calls = nil
}

type chainRecorder struct {
	mu     sync.Mutex
	before []intercept.Call
	after  []intercept.Call
	errs   []error
}

var _ intercept.Interceptor = (*chainRecorder)(nil)

func (c *chainRecorder) Before(ctx context.Context, call intercept.Call) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.before = append(c.before, call)
	return nil
}

func (c *chainRecorder) After(ctx context.Context, call intercept.Call, err error) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.after = append(c.after, call)
	c.errs = append(c.errs, err)
	return err
}

func (c *chainRecorder) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.before, c.after, c.errs = nil, nil, nil
}

func (c *chainRecorder) beforeOps() []string { return opsOf(c.snapshotBefore()) }

func (c *chainRecorder) afterOps() []string { return opsOf(c.snapshotAfter()) }

func (c *chainRecorder) snapshotBefore() []intercept.Call {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]intercept.Call(nil), c.before...)
}

func (c *chainRecorder) snapshotAfter() []intercept.Call {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]intercept.Call(nil), c.after...)
}

func opsOf(calls []intercept.Call) []string {
	out := make([]string, 0, len(calls))
	for _, c := range calls {
		out = append(out, string(c.Op))
	}
	return out
}

func namesOf(calls []intercept.Call) []string {
	out := make([]string, 0, len(calls))
	for _, c := range calls {
		out = append(out, c.Name)
	}
	return out
}

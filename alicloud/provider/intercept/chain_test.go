package intercept

import (
	"context"
	"reflect"
	"testing"
)

type namedStub struct {
	name string
}

func (s *namedStub) Before(ctx context.Context, call Call) error {
	return nil
}

func (s *namedStub) After(ctx context.Context, call Call, err error) error {
	return err
}

func setRegistries(t *testing.T, g []Interceptor, r map[string][]Interceptor) {
	t.Helper()
	oldGlobal, oldRegistry := global, registry
	t.Cleanup(func() { global, registry = oldGlobal, oldRegistry })
	global, registry = g, r
}

func TestChainOfEmpty(t *testing.T) {
	setRegistries(t, nil, map[string][]Interceptor{})
	if chain := ChainOf("alicloud_none", nil); len(chain) != 0 {
		t.Fatalf("ChainOf on empty registries: want empty chain, got %v", chain)
	}
}

func TestChainOfOrder(t *testing.T) {
	setRegistries(t,
		[]Interceptor{&namedStub{name: "g1"}, &namedStub{name: "g2"}},
		map[string][]Interceptor{"alicloud_x": {&namedStub{name: "r1"}}})

	chain := ChainOf("alicloud_x", []Interceptor{&namedStub{name: "d1"}})
	var got []string
	for _, i := range chain {
		got = append(got, i.(*namedStub).name)
	}
	want := []string{"g1", "g2", "d1", "r1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ChainOf order: got %v, want %v", got, want)
	}
}

func TestChainOfUnregisteredNameSeesGlobalOnly(t *testing.T) {
	setRegistries(t,
		[]Interceptor{&namedStub{name: "g1"}},
		map[string][]Interceptor{"alicloud_x": {&namedStub{name: "r1"}}})

	chain := ChainOf("alicloud_y", nil)
	if len(chain) != 1 || chain[0].(*namedStub).name != "g1" {
		t.Fatalf("ChainOf for unregistered name: want global only, got %v", chain)
	}
}

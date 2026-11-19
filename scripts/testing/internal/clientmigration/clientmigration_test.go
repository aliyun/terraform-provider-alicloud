package clientmigration

import (
	"strings"
	"testing"
)

func TestUnitClientConstructorMigration(t *testing.T) {
	const template = "package sample\nimport " + connectivityImport + "\nfunc check() {\nDECL\nCALL\nassert(1)\n}\n"
	source := func(decl, call string) string {
		return strings.NewReplacer("DECL", decl, "CALL", call).Replace(template)
	}
	for _, tc := range []struct {
		name, decl, oldCall, newCall string
		want                         bool
	}{
		{"pointer", "cfg := &connectivity.Config{}", "c, err := cfg.Client()", "c, err := connectivity.NewClient(cfg)", true},
		{"value", "cfg := connectivity.Config{}", "return cfg.Client()", "return connectivity.NewClient(&cfg)", true},
		{"inline", "", "return (&connectivity.Config{\n Region: region,\n}).Client()", "return connectivity.NewClient(&connectivity.Config{\n Region: region,\n})", true},
		{"changed receiver", "cfg := &connectivity.Config{}", "return other.Client()", "return connectivity.NewClient(cfg)", false},
		{"changed field", "", "return (&connectivity.Config{Region: old}).Client()", "return connectivity.NewClient(&connectivity.Config{Region: next})", false},
		{"unknown type", "cfg := other.Config{}", "return cfg.Client()", "return connectivity.NewClient(&cfg)", false},
		{"wrong indirection", "cfg := &connectivity.Config{}", "return cfg.Client()", "return connectivity.NewClient(&cfg)", false},
		{"different factory", "cfg := &connectivity.Config{}", "return cfg.Client()", "return other.NewClient(cfg)", false},
		{"shadowed package", "connectivity := other; cfg := &connectivity.Config{}", "return cfg.Client()", "return connectivity.NewClient(cfg)", false},
		{"different argument count", "cfg := &connectivity.Config{}", "return cfg.Client()", "return connectivity.NewClient(cfg, region)", false},
		{"changed assertion", "cfg := &connectivity.Config{}", "c := cfg.Client(); assert(c)", "c := connectivity.NewClient(cfg); assert(other)", false},
		{"changed step elsewhere", "cfg := &connectivity.Config{}", "c := cfg.Client()\nstep(first)", "c := connectivity.NewClient(cfg)\nstep(second)", false},
		{"text literal", "", `s := "cfg.Client()"`, `s := "connectivity.NewClient(cfg)"`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			old, next := source(tc.decl, tc.oldCall), source(tc.decl, tc.newCall)
			if got := equivalent([]byte(old), []byte(next), map[int]bool{5: true}); got != tc.want {
				t.Fatalf("equivalent = %v, want %v", got, tc.want)
			}
			if equivalent([]byte(old), []byte(next), map[int]bool{6: true}) {
				t.Fatal("a constructor outside the added lines was rewritten")
			}
			old = strings.Replace(old, connectivityImport, `"example.com/connectivity"`, 1)
			next = strings.Replace(next, connectivityImport, `"example.com/connectivity"`, 1)
			if equivalent([]byte(old), []byte(next), map[int]bool{5: true}) {
				t.Fatal("another package was treated as the provider's connectivity package")
			}
		})
	}
}

package coverage

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
)

func TestUnitSelection(t *testing.T) {
	src := []byte(`package alicloud; func resourceExample() {}`)
	old := map[string]string{"alicloud_example": "resourceExample"}
	current := map[string]string{"alicloud_example": "resourceExample"}
	sources := map[string][]byte{"alicloud/resource_alicloud_example.go": src}
	missing := func(string) ([]byte, error) { return nil, fmt.Errorf("missing") }
	cases := []struct {
		name          string
		paths         []string
		before, after map[string]string
		sources       map[string][]byte
		want          map[string]string
		fail          bool
	}{
		{"implementation", []string{"alicloud/resource_alicloud_example.go"}, old, current, sources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"test only or deleted test", []string{"alicloud/resource_alicloud_example_test.go"}, old, current, sources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"no resources", []string{"README.md"}, old, current, sources, map[string]string{}, false},
		{"registration", []string{"alicloud/provider.go"}, map[string]string{}, current, sources, map[string]string{"alicloud_example": "alicloud/resource_alicloud_example_test.go"}, false},
		{"unresolved registration", []string{"alicloud/provider.go"}, map[string]string{}, map[string]string{"alicloud_unknown": "unknown"}, sources, nil, true},
		{"unregistered changed file", []string{"alicloud/resource_alicloud_unknown.go"}, old, current, map[string][]byte{"alicloud/resource_alicloud_unknown.go": []byte(`package alicloud;func unknown(){}`)}, nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Select(tc.paths, tc.before, tc.after, tc.sources, missing)
			if (err != nil) != tc.fail || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v; want %v fail=%v", got, err, tc.want, tc.fail)
			}
		})
	}
	t.Run("rename", func(t *testing.T) {
		got, err := Select([]string{"alicloud/resource_alicloud_example.go", "alicloud/resource_alicloud_renamed.go"}, old, current, map[string][]byte{"alicloud/resource_alicloud_renamed.go": src}, func(string) ([]byte, error) { return src, nil })
		if err != nil || got["alicloud_example"] != "alicloud/resource_alicloud_renamed_test.go" {
			t.Fatal(got, err)
		}
	})
	t.Run("deleted resource", func(t *testing.T) {
		got, err := Select([]string{"alicloud/resource_alicloud_example.go"}, old, map[string]string{}, map[string][]byte{}, func(string) ([]byte, error) { return src, nil })
		if err != nil || len(got) != 0 {
			t.Fatal(got, err)
		}
	})
	// Framework-native resources whose flat-layout test file references an
	// implementation that moved to the service/ layout are handled by the
	// framework branch in Select. These tests stub frameworkRegistered and
	// frameworkFindByTypeName so no filesystem access is needed.
	t.Run("framework flat-layout test without ordered collections skipped", func(t *testing.T) {
		origReg := frameworkRegistered
		origResolve := frameworkFindByTypeName
		frameworkRegistered = func(typeName string) bool { return typeName == "alicloud_example" }
		frameworkFindByTypeName = func(typeName string) *splitlayout.Resolved {
			if typeName == "alicloud_example" {
				return &splitlayout.Resolved{TypeName: typeName, ImplPath: "alicloud/service/example/resource_example.go"}
			}
			return nil
		}
		defer func() { frameworkRegistered = origReg; frameworkFindByTypeName = origResolve }()
		implSrc := []byte(`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func NewExampleResource() resource.Resource { return nil }
func (r *exampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
}`)
		got, err := Select([]string{"alicloud/resource_alicloud_example_test.go"}, old, current,
			map[string][]byte{"alicloud/service/example/resource_example.go": implSrc}, missing)
		if err != nil || len(got) != 0 {
			t.Fatalf("expected framework-native resource with no ordered collections to be skipped; got %v, %v", got, err)
		}
	})
	t.Run("framework flat-layout test with ListNestedAttribute errors", func(t *testing.T) {
		origReg := frameworkRegistered
		origResolve := frameworkFindByTypeName
		frameworkRegistered = func(typeName string) bool { return typeName == "alicloud_example" }
		frameworkFindByTypeName = func(typeName string) *splitlayout.Resolved {
			if typeName == "alicloud_example" {
				return &splitlayout.Resolved{TypeName: typeName, ImplPath: "alicloud/service/example/resource_example.go"}
			}
			return nil
		}
		defer func() { frameworkRegistered = origReg; frameworkFindByTypeName = origResolve }()
		implSrc := []byte(`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func NewExampleResource() resource.Resource { return nil }
func (r *exampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"items": schema.ListNestedAttribute{NestedObject: schema.NestedAttributeObject{}},
	}}
}`)
		_, err := Select([]string{"alicloud/resource_alicloud_example_test.go"}, old, current,
			map[string][]byte{"alicloud/service/example/resource_example.go": implSrc}, missing)
		if err == nil {
			t.Fatal("expected error for framework-native resource with ordered collections; got nil")
		}
	})
	t.Run("framework flat-layout test with ListNestedBlock errors", func(t *testing.T) {
		origReg := frameworkRegistered
		origResolve := frameworkFindByTypeName
		frameworkRegistered = func(typeName string) bool { return typeName == "alicloud_example" }
		frameworkFindByTypeName = func(typeName string) *splitlayout.Resolved {
			if typeName == "alicloud_example" {
				return &splitlayout.Resolved{TypeName: typeName, ImplPath: "alicloud/service/example/resource_example.go"}
			}
			return nil
		}
		defer func() { frameworkRegistered = origReg; frameworkFindByTypeName = origResolve }()
		implSrc := []byte(`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func NewExampleResource() resource.Resource { return nil }
func (r *exampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Blocks: map[string]schema.Block{
		"items": schema.ListNestedBlock{NestedObject: schema.NestedBlockObject{}},
	}}
}`)
		_, err := Select([]string{"alicloud/resource_alicloud_example_test.go"}, old, current,
			map[string][]byte{"alicloud/service/example/resource_example.go": implSrc}, missing)
		if err == nil {
			t.Fatal("expected error for framework-native resource with ListNestedBlock; got nil")
		}
	})
	t.Run("framework flat-layout unresolvable path fails closed", func(t *testing.T) {
		origReg := frameworkRegistered
		origResolve := frameworkFindByTypeName
		frameworkRegistered = func(typeName string) bool { return typeName == "alicloud_example" }
		frameworkFindByTypeName = func(typeName string) *splitlayout.Resolved { return nil }
		defer func() { frameworkRegistered = origReg; frameworkFindByTypeName = origResolve }()
		_, err := Select([]string{"alicloud/resource_alicloud_example_test.go"}, old, current,
			map[string][]byte{"alicloud/service/example/resource_example.go": []byte(`package example`)}, missing)
		if err == nil {
			t.Fatal("expected error for framework-native resource with unresolvable implementation path; got nil")
		}
	})
	t.Run("framework flat-layout missing source fails closed", func(t *testing.T) {
		origReg := frameworkRegistered
		origResolve := frameworkFindByTypeName
		frameworkRegistered = func(typeName string) bool { return typeName == "alicloud_example" }
		frameworkFindByTypeName = func(typeName string) *splitlayout.Resolved {
			return &splitlayout.Resolved{TypeName: typeName, ImplPath: "alicloud/service/example/resource_example.go"}
		}
		defer func() { frameworkRegistered = origReg; frameworkFindByTypeName = origResolve }()
		// No source provided for the resolved implementation path.
		_, err := Select([]string{"alicloud/resource_alicloud_example_test.go"}, old, current,
			map[string][]byte{}, missing)
		if err == nil {
			t.Fatal("expected error for framework-native resource with missing implementation source; got nil")
		}
	})
}

func TestUnitFrameworkOrderedCollections(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		wantErr bool
		wantHas bool
	}{
		{
			"no ordered collections",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
}`,
			false, false,
		},
		{
			"ListAttribute present",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{"items": schema.ListAttribute{}}}
}`,
			false, true,
		},
		{
			"ListNestedAttribute present",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{"items": schema.ListNestedAttribute{NestedObject: schema.NestedAttributeObject{}}}}
}`,
			false, true,
		},
		{
			"ListNestedBlock present",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Blocks: map[string]schema.Block{"items": schema.ListNestedBlock{NestedObject: schema.NestedBlockObject{}}}}
}`,
			false, true,
		},
		{
			"SetAttribute not ordered",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{"items": schema.SetAttribute{}}}
}`,
			false, false,
		},
		{
			"SetNestedAttribute not ordered",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{"items": schema.SetNestedAttribute{NestedObject: schema.NestedAttributeObject{}}}}
}`,
			false, false,
		},
		{
			"schema from cross-file helper fails closed",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schemaHelper()
}`,
			true, false,
		},
		{
			"schema from method call fails closed",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func (r T) Schema() schema.Schema {
	return r.buildSchema()
}`,
			true, false,
		},
		{
			"schema literal with helper-built Attributes fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func resourceSchema() schema.Schema {
	return schema.Schema{Attributes: attributesFromAnotherFile()}
}`,
			true, false,
		},
		{
			"schema literal with helper-built Blocks fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func resourceSchema() schema.Schema {
	return schema.Schema{Blocks: blocksFromAnotherFile()}
}`,
			true, false,
		},
		{
			"schema literal with variable-built Attributes fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func resourceSchema() schema.Schema {
	attrs := someAttributes
	return schema.Schema{Attributes: attrs}
}`,
			true, false,
		},
		{
			// Coordinator #8 exact counter-example: a cross-file helper buried
			// inside a SingleNestedAttribute's Attributes. A top-level scan sees
			// only StringAttribute-shaped leaves and would return (false, nil);
			// the recursive descent must classify the nested helper as
			// unresolved and fail closed.
			"nested SingleNestedAttribute with cross-file helper fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func resourceSchema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{
		"nested": schema.SingleNestedAttribute{Optional: true, Attributes: nestedAttributesFromAnotherFile()},
	}}
}`,
			true, false,
		},
		{
			"nested SingleNestedAttribute with inline unordered Attributes is clean",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{
		"nested": schema.SingleNestedAttribute{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}},
	}}
}`,
			false, false,
		},
		{
			"nested SingleNestedAttribute with inline ListAttribute is ordered",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{
		"nested": schema.SingleNestedAttribute{Attributes: map[string]schema.Attribute{"items": schema.ListAttribute{}}},
	}}
}`,
			false, true,
		},
		{
			"SetNestedAttribute NestedObject with cross-file helper fails closed",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{
		"set": schema.SetNestedAttribute{NestedObject: schema.NestedAttributeObject{Attributes: nestedAttributesFromAnotherFile()}},
	}}
}`,
			true, false,
		},
		{
			"SetNestedAttribute NestedObject with inline ListAttribute is ordered",
			`package example
import "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{
		"set": schema.SetNestedAttribute{NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{"items": schema.ListAttribute{}}}},
	}}
}`,
			false, true,
		},
		{
			// Root review #6 exact counter-example: a schema literal built inline
			// is then mutated by an index-map assignment whose value is a
			// cross-file helper returning an ordered collection. The inline scan
			// sees only the inline "name" entry and would return (false, nil);
			// the post-construction write must be detected and fail closed.
			"post-construction index-map assignment to Attributes fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	s := schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{Required: true}}}
	s.Attributes["ordered"] = orderedField()
	return s
}`,
			true, false,
		},
		{
			// Post-construction index-map assignment with an inline UNORDERED
			// literal: the main scan sees no ordered selector, so without the
			// mutation guard it would return (false, nil), silently ignoring
			// the element added after construction. The guard must fail closed
			// because the schema is no longer built fully inline.
			"post-construction index-map assignment with inline unordered literal fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	s := schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
	s.Attributes["extra"] = schema.StringAttribute{}
	return s
}`,
			true, false,
		},
		{
			"post-construction index-map assignment to Blocks fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	s := schema.Schema{Blocks: map[string]schema.Block{}}
	s.Blocks["items"] = blockFromAnotherFile()
	return s
}`,
			true, false,
		},
		{
			"post-construction field replacement of Attributes fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	s := schema.Schema{}
	s.Attributes = attributesFromAnotherFile()
	return s
}`,
			true, false,
		},
		{
			"post-construction field replacement of Blocks with helper fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	s := schema.Schema{}
	s.Blocks = blocksFromAnotherFile()
	return s
}`,
			true, false,
		},
		{
			// A post-construction write to a non-collection schema field (e.g.
			// Description) is not an Attributes/Blocks mutation and must not be
			// rejected; the inline schema stays ordered-free.
			"post-construction write to non-Attributes field stays clean",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	s := schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
	s.Description = "example"
	return s
}`,
			false, false,
		},
		{
			// alias escape: the Attributes map is copied into a local that the
			// single-file analysis cannot track, then mutated through the
			// alias. The later write is invisible, so the alias itself must be
			// rejected rather than chasing the alias use.
			"post-construction alias of Attributes fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	s := schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
	m := s.Attributes
	m["ordered"] = orderedField()
	return s
}`,
			true, false,
		},
		{
			// pointer alias: &s.Attributes escapes the map into a pointer the
			// analysis cannot track; reject the escape instead of following
			// the pointer dereference.
			"post-construction pointer alias of Blocks fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	s := schema.Schema{Blocks: map[string]schema.Block{"root": schema.SingleNestedBlock{}}}
	p := &s.Blocks
	(*p)["ordered"] = blockFromAnotherFile()
	return s
}`,
			true, false,
		},
		{
			// var declaration alias: the Attributes map is captured by a var
			// declaration rather than :=, exercising the GenDecl/ValueSpec
			// path; the escape must still be rejected.
			"post-construction var declaration alias of Attributes fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
var m = s.Attributes
func Schema() schema.Schema {
	s := schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
	m["ordered"] = orderedField()
	return s
}`,
			true, false,
		},
		{
			// receiver escape: a method on the schema type mutates Attributes
			// through the receiver; the write is recognized via the .Attributes
			// selector on the receiver LHS, without enumerating receiver forms.
			"post-construction receiver write to Attributes fails closed",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
type T struct{ s schema.Schema }
func (t *T) extend() {
	t.s.Attributes["ordered"] = orderedField()
}
func Schema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
}`,
			true, false,
		},
		{
			// reading a non-Attributes selector (Description) into a local is
			// not a collection escape and must stay clean.
			"post-construction read of non-Attributes field stays clean",
			`package ecs
import schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
func Schema() schema.Schema {
	s := schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
	d := s.Description
	_ = d
	return s
}`,
			false, false,
		},
		{
			"syntax error fails closed",
			`package example this is not valid go`,
			true, false,
		},
		{
			"no schema import fails closed",
			`package example
func Schema() string { return "" }`,
			true, false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			has, err := frameworkOrderedCollections([]byte(tc.src))
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr && has != tc.wantHas {
				t.Fatalf("has=%v wantHas=%v", has, tc.wantHas)
			}
		})
	}
}

func TestUnitRegistry(t *testing.T) {
	got, err := Registry([]byte(`package alicloud;var p = Provider{ResourcesMap:map[string]*Resource{"alicloud_example":resourceExample()}}`))
	if err != nil || got["alicloud_example"] != "resourceExample" {
		t.Fatal(got, err)
	}
	for _, src := range []string{`package alicloud;var p = Provider{ResourcesMap: dynamic()}`, `package alicloud;var p = Provider{ResourcesMap:map[string]*Resource{"alicloud_example":dynamic(name)}}`} {
		if _, err := Registry([]byte(src)); err == nil {
			t.Fatal("unresolved registry passed")
		}
	}
}

func TestUnitSplitSelection(t *testing.T) {
	impl := []byte(`package cloudauth; func ResourceFaceConfig() {}`)
	implPath := "alicloud/service/cloudauth/resource_face_config.go"
	testPath := "alicloud/service/cloudauth/resource_face_config_test.go"
	imsPath := "alicloud/service/ims/resource_user.go"
	after := map[string]string{"alicloud_cloudauth_face_config": "cloudauth.ResourceFaceConfig"}
	before := map[string]string{"alicloud_cloudauth_face_config": "resourceAlicloudCloudauthFaceConfig"}
	sources := map[string][]byte{
		implPath: impl,
		imsPath:  []byte(`package ims; func resourceUser() {}`),
	}
	missing := func(string) ([]byte, error) { return nil, fmt.Errorf("missing") }
	want := map[string]string{"alicloud_cloudauth_face_config": testPath}
	cases := []struct {
		name   string
		paths  []string
		before map[string]string
		want   map[string]string
	}{
		{"split implementation", []string{implPath}, before, want},
		{"split test file", []string{testPath}, before, want},
		{"register is a registration source", []string{"alicloud/service/cloudauth/register.go"}, map[string]string{}, want},
		{"framework file without registration", []string{imsPath}, before, map[string]string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Select(tc.paths, tc.before, after, sources, missing)
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
		})
	}
	t.Run("cross-package function names do not collide", func(t *testing.T) {
		shared := map[string][]byte{
			"alicloud/service/a/resource_x.go": []byte(`package a; func Shared() {}`),
			"alicloud/service/b/resource_x.go": []byte(`package b; func Shared() {}`),
		}
		got, err := Select([]string{"alicloud/service/a/register.go"}, map[string]string{}, map[string]string{"alicloud_a_x": "a.Shared"}, shared, missing)
		if err != nil || got["alicloud_a_x"] != "alicloud/service/a/resource_x_test.go" {
			t.Fatal(got, err)
		}
	})
	t.Run("within-package duplicate fails", func(t *testing.T) {
		dup := map[string][]byte{
			"alicloud/service/a/resource_x.go": []byte(`package a; func Dup() {}`),
			"alicloud/service/a/resource_y.go": []byte(`package a; func Dup() {}`),
		}
		if _, err := Select(nil, map[string]string{}, map[string]string{}, dup, missing); err == nil {
			t.Fatal("expected duplicate function failure")
		}
	})
	t.Run("deleted split implementation with kept registration", func(t *testing.T) {
		if _, err := Select([]string{implPath}, before, after, map[string][]byte{}, func(string) ([]byte, error) { return impl, nil }); err == nil {
			t.Fatal("expected missing-implementation failure")
		}
	})
	t.Run("deleted split resource", func(t *testing.T) {
		got, err := Select([]string{implPath}, before, map[string]string{}, map[string][]byte{}, func(string) ([]byte, error) { return impl, nil })
		if err != nil || len(got) != 0 {
			t.Fatal(got, err)
		}
	})
}

func TestUnitRegistryUnion(t *testing.T) {
	providerSrc := []byte(`package alicloud;var p = Provider{ResourcesMap:map[string]*Resource{"alicloud_example":resourceExample()}}`)
	register := []byte(`package cloudauth
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cloudauth",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_cloudauth_face_config", Factory: ResourceFaceConfig},
		},
		SDKDataSources: []conns.SDKResource{
			{TypeName: "alicloud_cloudauth_face_configs", Factory: DataSourceCloudauthFaceConfigs},
		},
		Resources: []conns.FrameworkResource{
			{TypeName: "alicloud_cloudauth_framework", Factory: NewFramework},
		},
	}
}`)
	got, err := RegistryUnion(providerSrc, map[string][]byte{"alicloud/service/cloudauth/register.go": register})
	if err != nil {
		t.Fatal(err)
	}
	if got["alicloud_example"] != "resourceExample" {
		t.Fatal("flat entry lost:", got)
	}
	if got["alicloud_cloudauth_face_config"] != "cloudauth.ResourceFaceConfig" {
		t.Fatal("SDK registration missing:", got)
	}
	if _, ok := got["alicloud_cloudauth_face_configs"]; ok {
		t.Fatal("data source leaked into the resource registry:", got)
	}
	if _, ok := got["alicloud_cloudauth_framework"]; ok {
		t.Fatal("framework registration leaked into the resource registry:", got)
	}
	dup := []byte(`package cloudauth
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{SDKResources: []conns.SDKResource{{TypeName: "alicloud_example", Factory: resourceExample}}}
}`)
	if _, err := RegistryUnion(providerSrc, map[string][]byte{"alicloud/service/cloudauth/register.go": dup}); err == nil {
		t.Fatal("duplicate declaration passed")
	}
}

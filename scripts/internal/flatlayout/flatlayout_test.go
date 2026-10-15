package flatlayout

import (
	"reflect"
	"testing"
)

// A provider.go extract shaped like the real one: both map literals in the
// Provider() composite, including the mode-B entry whose constructor is
// declared in a differently named file.
const providerSrc = `package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {Type: schema.TypeString, Required: true},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alicloud_instance_types": dataSourceAlicloudInstanceTypes(),
			"alicloud_havips":         dataSourceAlicloudHavips(),
			// "alicloud_subnet" aims to match aws usage habit.
			"alicloud_subnets": dataSourceAlicloudSubnets(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_instance":                resourceAlicloudInstance(),
			"alicloud_edas_application_deployment": resourceAlicloudEdasApplicationPackageAttachment(),
		},
	}
}
`

// The shared plumbing file: schema helpers only, no constructor at all.
const sharedPlumbingSrc = `package alicloud

func outputInstancesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ids": {Type: schema.TypeList, Computed: true},
	}
}

func outputShortVpcsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vpc_ids": {Type: schema.TypeList, Computed: true},
	}
}
`

// A mode-B file: the constructor is registered under a different TypeName.
const modeBSrc = `package alicloud

func resourceAlicloudEdasApplicationPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasApplicationPackageAttachmentCreate,
	}
}

func resourceAlicloudEdasApplicationPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
`

// A normal file plus a legacy V0 helper: one registration, one
// constructor-shaped function registered nowhere.
const v0HelperSrc = `package alicloud

func resourceAlicloudInstance() *schema.Resource {
	return &schema.Resource{}
}

func resourceAlicloudInstanceV0() *schema.Resource {
	return &schema.Resource{}
}
`

// A file declaring two constructors registered under different names.
const multiCtorSrc = `package alicloud

func dataSourceAlicloudInstanceTypes() *schema.Resource {
	return &schema.Resource{}
}

func dataSourceAlicloudSubnets() *schema.Resource {
	return &schema.Resource{}
}
`

func TestProviderRegistrations(t *testing.T) {
	regs := ProviderRegistrations([]byte(providerSrc))
	want := map[string][]Registration{
		"dataSourceAlicloudInstanceTypes":                  {{TypeName: "alicloud_instance_types", Factory: "dataSourceAlicloudInstanceTypes", DataSource: true}},
		"dataSourceAlicloudHavips":                         {{TypeName: "alicloud_havips", Factory: "dataSourceAlicloudHavips", DataSource: true}},
		"dataSourceAlicloudSubnets":                        {{TypeName: "alicloud_subnets", Factory: "dataSourceAlicloudSubnets", DataSource: true}},
		"resourceAlicloudInstance":                         {{TypeName: "alicloud_instance", Factory: "resourceAlicloudInstance", DataSource: false}},
		"resourceAlicloudEdasApplicationPackageAttachment": {{TypeName: "alicloud_edas_application_deployment", Factory: "resourceAlicloudEdasApplicationPackageAttachment", DataSource: false}},
	}
	if !reflect.DeepEqual(regs, want) {
		t.Errorf("ProviderRegistrations(provider) = %+v, want %+v", regs, want)
	}
	if got := ProviderRegistrations([]byte("not go source")); got != nil {
		t.Errorf("ProviderRegistrations(malformed) = %+v, want nil", got)
	}
}

// An alias pair: one constructor registered under two TypeNames, like
// alicloud_havip and alicloud_vpc_ha_vip both calling
// resourceAliCloudVpcHaVip.
const aliasProviderSrc = `package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_havip":      resourceAliCloudVpcHaVip(),
			"alicloud_vpc_ha_vip": resourceAliCloudVpcHaVip(),
		},
	}
}
`

const aliasSrc = `package alicloud

func resourceAliCloudVpcHaVip() *schema.Resource {
	return &schema.Resource{}
}
`

func TestJoinSrcSharedConstructorYieldsEveryRegistration(t *testing.T) {
	regs := ProviderRegistrations([]byte(aliasProviderSrc))
	got, err := JoinSrc("alicloud/resource_alicloud_vpc_ha_vip.go", []byte(aliasSrc), regs)
	if err != nil {
		t.Fatalf("JoinSrc(alias) error: %s", err)
	}
	if len(got.Registrations) != 2 {
		t.Fatalf("JoinSrc(alias) returned %d registrations, want 2 (one per TypeName)", len(got.Registrations))
	}
	seen := map[string]bool{}
	for _, r := range got.Registrations {
		seen[r.TypeName] = true
	}
	for _, want := range []string{"alicloud_havip", "alicloud_vpc_ha_vip"} {
		if !seen[want] {
			t.Errorf("JoinSrc(alias) is missing the %s registration", want)
		}
	}
}

func TestJoinSrcSharedPlumbingRegistersNothing(t *testing.T) {
	regs := ProviderRegistrations([]byte(providerSrc))
	got, err := JoinSrc("alicloud/data_source_alicloud_common.go", []byte(sharedPlumbingSrc), regs)
	if err != nil {
		t.Fatalf("JoinSrc(shared plumbing) error: %s", err)
	}
	if len(got.Registrations) != 0 || len(got.Unregistered) != 0 {
		t.Errorf("JoinSrc(shared plumbing) = %+v, want empty Result", got)
	}
}

func TestJoinSrcModeBResolvesRealEntity(t *testing.T) {
	regs := ProviderRegistrations([]byte(providerSrc))
	got, err := JoinSrc("alicloud/resource_alicloud_edas_application_package_attachment.go", []byte(modeBSrc), regs)
	if err != nil {
		t.Fatalf("JoinSrc(mode B) error: %s", err)
	}
	want := Result{Registrations: []Registration{
		{TypeName: "alicloud_edas_application_deployment", Factory: "resourceAlicloudEdasApplicationPackageAttachment", DataSource: false},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("JoinSrc(mode B) = %+v, want %+v", got, want)
	}
}

func TestJoinSrcSeparatesV0HelperFromRegistration(t *testing.T) {
	regs := ProviderRegistrations([]byte(providerSrc))
	got, err := JoinSrc("alicloud/resource_alicloud_instance.go", []byte(v0HelperSrc), regs)
	if err != nil {
		t.Fatalf("JoinSrc(V0 helper) error: %s", err)
	}
	if len(got.Registrations) != 1 || got.Registrations[0].TypeName != "alicloud_instance" {
		t.Errorf("JoinSrc(V0 helper) registrations = %+v, want only alicloud_instance", got.Registrations)
	}
	if !reflect.DeepEqual(got.Unregistered, []string{"resourceAlicloudInstanceV0"}) {
		t.Errorf("JoinSrc(V0 helper) unregistered = %v, want [resourceAlicloudInstanceV0]", got.Unregistered)
	}
}

func TestJoinSrcMultipleRegistrations(t *testing.T) {
	regs := ProviderRegistrations([]byte(providerSrc))
	got, err := JoinSrc("alicloud/data_source_alicloud_instance_types.go", []byte(multiCtorSrc), regs)
	if err != nil {
		t.Fatalf("JoinSrc(multi ctor) error: %s", err)
	}
	if len(got.Registrations) != 2 {
		t.Fatalf("JoinSrc(multi ctor) returned %d registrations, want 2", len(got.Registrations))
	}
	for _, want := range []string{"alicloud_instance_types", "alicloud_subnets"} {
		found := false
		for _, r := range got.Registrations {
			if r.TypeName == want && r.DataSource {
				found = true
			}
		}
		if !found {
			t.Errorf("JoinSrc(multi ctor) is missing the %s data source registration", want)
		}
	}
}

func TestJoinSrcRejectsMalformedSource(t *testing.T) {
	if _, err := JoinSrc("alicloud/resource_alicloud_instance.go", []byte("not go source"), nil); err == nil {
		t.Error("JoinSrc(malformed source) succeeded, want error")
	}
}

func TestImplPathOf(t *testing.T) {
	cases := map[string]string{
		"alicloud/resource_alicloud_instance.go":             "alicloud/resource_alicloud_instance.go",
		"alicloud/resource_alicloud_instance_test.go":        "alicloud/resource_alicloud_instance.go",
		"alicloud/data_source_alicloud_havips_test.go":       "alicloud/data_source_alicloud_havips.go",
		"alicloud/service/cloudauth/resource_face_config.go": "alicloud/service/cloudauth/resource_face_config.go",
	}
	for in, want := range cases {
		if got := ImplPathOf(in); got != want {
			t.Errorf("ImplPathOf(%s) = %s, want %s", in, got, want)
		}
	}
}

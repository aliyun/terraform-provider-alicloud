package splitlayout

import "testing"

// A framework-only register.go, shaped like alicloud/service/ims/register.go.
const frameworkRegisterSrc = `package ims

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ims",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_ims_default_domain",
				Factory:  NewDefaultDomainDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_ims_user",
				Factory:  NewUserResource,
			},
		},
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_ims_plugin",
				Factory:  NewPluginEphemeralResource,
			},
		},
	}
}
`

// An SDKv2 register.go, shaped like alicloud/service/cloudauth/register.go.
const sdkRegisterSrc = `package cloudauth

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cloudauth",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cloudauth_face_config",
				Factory:  ResourceFaceConfig,
			},
		},
	}
}
`

// A register.go using every registration field, shaped like
// alicloud/service/cas/register.go extended with the remaining kinds.
const allKindsRegisterSrc = `package cas

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cas",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_cas_certificates",
				Factory:  DataSourceCasCertificates,
			},
		},
		ListResources: []conns.ListResource{
			{
				TypeName: "alicloud_cas_certificate",
				Factory:  NewCertificateListResource,
			},
		},
		Actions: []conns.Action{
			{
				TypeName: "alicloud_ssl_certificates_service_pca_cert_sync",
				Factory:  NewPcaCertSyncAction,
			},
		},
	}
}
`

func TestFrameworkRegisteredIn(t *testing.T) {
	for _, typeName := range []string{
		"alicloud_ims_user",
		"alicloud_ims_default_domain",
		"alicloud_ims_plugin",
	} {
		if !FrameworkRegisteredIn([]byte(frameworkRegisterSrc), typeName) {
			t.Errorf("FrameworkRegisteredIn(framework register, %s) = false, want true", typeName)
		}
	}
	for _, typeName := range []string{
		"alicloud_ims_missing",
		// SDKv2 registrations must never leak into the framework set.
		"alicloud_cloudauth_face_config",
	} {
		if FrameworkRegisteredIn([]byte(frameworkRegisterSrc), typeName) {
			t.Errorf("FrameworkRegisteredIn(framework register, %s) = true, want false", typeName)
		}
	}
	if FrameworkRegisteredIn([]byte(sdkRegisterSrc), "alicloud_cloudauth_face_config") {
		t.Error("FrameworkRegisteredIn(SDK register, alicloud_cloudauth_face_config) = true, want false")
	}
	if FrameworkRegisteredIn([]byte("not go source"), "alicloud_ims_user") {
		t.Error("FrameworkRegisteredIn(malformed source) = true, want false")
	}
}

func TestSDKRegisterEntriesIgnoresFrameworkKeys(t *testing.T) {
	if got := SDKRegisterEntries([]byte(frameworkRegisterSrc)); len(got) != 0 {
		t.Errorf("SDKRegisterEntries(framework register) = %v, want empty", got)
	}
	got := SDKRegisterEntries([]byte(sdkRegisterSrc))
	if len(got) != 1 {
		t.Fatalf("SDKRegisterEntries(SDK register) returned %d entries, want 1", len(got))
	}
	if got[0].TypeName != "alicloud_cloudauth_face_config" || got[0].Factory != "ResourceFaceConfig" || got[0].Kind != KindResource || !got[0].IsResource() {
		t.Errorf("SDKRegisterEntries(SDK register) = %+v, want the face config resource entry", got[0])
	}
}

func TestFrameworkRegisterEntries(t *testing.T) {
	got := FrameworkRegisterEntries([]byte(frameworkRegisterSrc))
	if len(got) != 3 {
		t.Fatalf("FrameworkRegisterEntries(framework register) returned %d entries, want 3", len(got))
	}
	byName := map[string]Entry{}
	for _, e := range got {
		byName[e.TypeName] = e
	}
	wantKinds := map[string]Kind{
		"alicloud_ims_user":           KindResource,
		"alicloud_ims_default_domain": KindDataSource,
		"alicloud_ims_plugin":         KindEphemeral,
	}
	for typeName, kind := range wantKinds {
		e, ok := byName[typeName]
		if !ok {
			t.Errorf("FrameworkRegisterEntries(framework register) is missing %s", typeName)
			continue
		}
		if e.Kind != kind {
			t.Errorf("FrameworkRegisterEntries(framework register): %s has kind %s, want %s", typeName, e.Kind, kind)
		}
	}
	if e, ok := byName["alicloud_ims_user"]; ok && e.Factory != "NewUserResource" {
		t.Errorf("FrameworkRegisterEntries(framework register): alicloud_ims_user factory = %s, want NewUserResource", e.Factory)
	}
	if entries := FrameworkRegisterEntries([]byte(sdkRegisterSrc)); len(entries) != 0 {
		t.Errorf("FrameworkRegisterEntries(SDK register) = %v, want empty", entries)
	}
}

func TestRegisterEntriesAllKinds(t *testing.T) {
	sdk := SDKRegisterEntries([]byte(allKindsRegisterSrc))
	if len(sdk) != 1 || sdk[0].Kind != KindDataSource || sdk[0].IsResource() {
		t.Errorf("SDKRegisterEntries(all kinds register) = %v, want the certificates data source only", sdk)
	}
	fw := FrameworkRegisterEntries([]byte(allKindsRegisterSrc))
	if len(fw) != 2 {
		t.Fatalf("FrameworkRegisterEntries(all kinds register) returned %d entries, want 2", len(fw))
	}
	byName := map[string]Entry{}
	for _, e := range fw {
		byName[e.TypeName] = e
	}
	if e, ok := byName["alicloud_cas_certificate"]; !ok || e.Kind != KindListResource {
		t.Errorf("FrameworkRegisterEntries(all kinds register): alicloud_cas_certificate = %+v, want a list resource", byName["alicloud_cas_certificate"])
	}
	if e, ok := byName["alicloud_ssl_certificates_service_pca_cert_sync"]; !ok || e.Kind != KindAction {
		t.Errorf("FrameworkRegisterEntries(all kinds register): alicloud_ssl_certificates_service_pca_cert_sync = %+v, want an action", byName["alicloud_ssl_certificates_service_pca_cert_sync"])
	}
}

func TestIsSplitSourceStems(t *testing.T) {
	cases := []struct {
		path          string
		product, stem string
		isTest, ok    bool
	}{
		{"alicloud/service/cloudauth/resource_face_config.go", "cloudauth", "resource_face_config", false, true},
		{"alicloud/service/cloudauth/resource_face_config_test.go", "cloudauth", "resource_face_config", true, true},
		{"alicloud/service/kms/ephemeral_plaintext.go", "kms", "ephemeral_plaintext", false, true},
		{"alicloud/service/kms/ephemeral_plaintext_test.go", "kms", "ephemeral_plaintext", true, true},
		{"alicloud/service/cas/action_pca_cert_sync.go", "cas", "action_pca_cert_sync", false, true},
		{"alicloud/service/cas/action_pca_cert_sync_test.go", "cas", "action_pca_cert_sync", true, true},
		{"alicloud/service/ims/data_source_default_domain.go", "ims", "data_source_default_domain", false, true},
		{"alicloud/service/vpc/list_vswitch.go", "vpc", "list_vswitch", false, true},
		{"alicloud/service/cloudauth/register.go", "", "", false, false},
		{"alicloud/service/cloudauth/service.go", "", "", false, false},
		{"alicloud/resource_alicloud_vpc.go", "", "", false, false},
		{"alicloud/service/kms/ephemeral.go", "", "", false, false},
	}
	for _, c := range cases {
		product, stem, isTest, ok := IsSplitSource(c.path)
		if product != c.product || stem != c.stem || isTest != c.isTest || ok != c.ok {
			t.Errorf("IsSplitSource(%s) = (%q, %q, %v, %v), want (%q, %q, %v, %v)",
				c.path, product, stem, isTest, ok, c.product, c.stem, c.isTest, c.ok)
		}
	}
}

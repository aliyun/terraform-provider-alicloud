package splitlayout

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

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

// A register.go that registers one type name twice — once as a managed
// resource and once as a data source, each with its own factory in its own
// implementation file — shaped like alicloud/service/sae/register.go, whose
// same-named pair is why the reverse lookup matches the kind as well as the
// name.
const sameNameResourceFirstRegisterSrc = `package sae

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sae",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBDataSource,
			},
		},
	}
}
`

// The same pair with the declaration order flipped, so a lookup can be proven
// to depend on the (kind, name) identity rather than on register order.
const sameNameDataSourceFirstRegisterSrc = `package sae

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sae",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBResource,
			},
		},
	}
}
`

// writeSameNameTree materializes a product directory whose register.go
// declares the same-named resource/data-source pair and whose two
// implementation files declare one factory each. The implementation bodies
// are never compiled — only parsed — so stubs suffice.
func writeSameNameTree(t *testing.T, registerSrc string) {
	t.Helper()
	dir := filepath.Join("alicloud", "service", "sae")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		filepath.Join(dir, "register.go"):                      registerSrc,
		filepath.Join(dir, "resource_load_balancer_nlb.go"):    "package sae\n\nfunc NewLoadBalancerNLBResource() *nlbResource { return nil }\n",
		filepath.Join(dir, "data_source_load_balancer_nlb.go"): "package sae\n\nfunc NewLoadBalancerNLBDataSource() *nlbDataSource { return nil }\n",
	}
	for path, src := range files {
		if err := os.WriteFile(path, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestUnitFrameworkFindByKindAndTypeName(t *testing.T) {
	const name = "alicloud_sae_load_balancer_nlb"
	for _, tc := range []struct {
		desc     string
		register string
	}{
		{"resource declared first", sameNameResourceFirstRegisterSrc},
		{"data source declared first", sameNameDataSourceFirstRegisterSrc},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Chdir(t.TempDir())
			writeSameNameTree(t, tc.register)

			r := FrameworkFindByKindAndTypeName(KindResource, name)
			if r == nil {
				t.Fatal("FrameworkFindByKindAndTypeName(resource, name) = nil, want the resource registration")
			}
			if r.Kind != KindResource || filepath.Base(r.ImplPath) != "resource_load_balancer_nlb.go" {
				t.Errorf("resource lookup resolved to (%s, %s), want (resource, resource_load_balancer_nlb.go)", r.Kind, r.ImplPath)
			}

			d := FrameworkFindByKindAndTypeName(KindDataSource, name)
			if d == nil {
				t.Fatal("FrameworkFindByKindAndTypeName(data source, name) = nil, want the data source registration")
			}
			if d.Kind != KindDataSource || filepath.Base(d.ImplPath) != "data_source_load_balancer_nlb.go" {
				t.Errorf("data source lookup resolved to (%s, %s), want (data_source, data_source_load_balancer_nlb.go)", d.Kind, d.ImplPath)
			}

			if got := FrameworkFindByKindAndTypeName(KindEphemeral, name); got != nil {
				t.Errorf("FrameworkFindByKindAndTypeName(ephemeral, name) = %+v, want nil: the name is not registered as an ephemeral", got)
			}
		})
	}
}

// TestUnitFrameworkFactoriesDeclaredIn pins the contract the file-granularity
// guard relies on: the entries returned are exactly the registrations whose
// factories the given source declares, scoped to the file's own product, and
// every failure to prove that is an error the caller must fail closed on
// rather than a silent "no conflict".
func TestUnitFrameworkFactoriesDeclaredIn(t *testing.T) {
	const impl = "alicloud/service/sae/data_source_load_balancer_nlb.go"
	t.Run("one factory per file resolves its own entry", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameTree(t, sameNameResourceFirstRegisterSrc)

		src, err := os.ReadFile(impl)
		if err != nil {
			t.Fatal(err)
		}
		got, err := FrameworkFactoriesDeclaredIn(impl, src)
		if err != nil {
			t.Fatalf("FrameworkFactoriesDeclaredIn(data source impl) errored: %s", err)
		}
		if len(got) != 1 || got[0].TypeName != "alicloud_sae_load_balancer_nlb" || got[0].Kind != KindDataSource || got[0].Factory != "NewLoadBalancerNLBDataSource" {
			t.Errorf("FrameworkFactoriesDeclaredIn(data source impl) = %+v, want the pair's data source entry alone", got)
		}
	})
	t.Run("a file declaring both factories of the pair reports both", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameTree(t, sameNameResourceFirstRegisterSrc)
		shared := "alicloud/service/sae/load_balancer_nlb.go"
		sharedSrc := "package sae\n\nfunc NewLoadBalancerNLBResource() *nlbResource { return nil }\n\nfunc NewLoadBalancerNLBDataSource() *nlbDataSource { return nil }\n"
		if err := os.WriteFile(shared, []byte(sharedSrc), 0644); err != nil {
			t.Fatal(err)
		}

		got, err := FrameworkFactoriesDeclaredIn(shared, []byte(sharedSrc))
		if err != nil {
			t.Fatalf("FrameworkFactoriesDeclaredIn(shared impl) errored: %s", err)
		}
		if len(got) != 2 {
			t.Errorf("FrameworkFactoriesDeclaredIn(shared impl) = %+v, want both entries of the pair", got)
		}
	})
	t.Run("another product's same-named factory is not this file's", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameTree(t, sameNameResourceFirstRegisterSrc)
		// ecs's register references the same factory name, backed by its own
		// same-named declaration in its own package — the only way such a
		// register compiles.
		ecsDir := filepath.Join("alicloud", "service", "ecs")
		if err := os.MkdirAll(ecsDir, 0755); err != nil {
			t.Fatal(err)
		}
		ecsRegister := "package ecs\n\nfunc ServicePackage() conns.ServicePackage {\n	return conns.ServicePackage{\n		Name: \"Ecs\",\n		DataSources: []conns.DataSource{\n			{\n				TypeName: \"alicloud_ecs_load_balancer_nlb\",\n				Factory:  NewLoadBalancerNLBDataSource,\n			},\n		},\n	}\n}\n"
		if err := os.WriteFile(filepath.Join(ecsDir, "register.go"), []byte(ecsRegister), 0644); err != nil {
			t.Fatal(err)
		}
		ecsImpl := "package ecs\n\nfunc NewLoadBalancerNLBDataSource() *nlbDataSource { return nil }\n"
		if err := os.WriteFile(filepath.Join(ecsDir, "data_source_load_balancer_nlb.go"), []byte(ecsImpl), 0644); err != nil {
			t.Fatal(err)
		}

		src, err := os.ReadFile(impl)
		if err != nil {
			t.Fatal(err)
		}
		got, err := FrameworkFactoriesDeclaredIn(impl, src)
		if err != nil {
			t.Fatalf("FrameworkFactoriesDeclaredIn(data source impl) errored: %s", err)
		}
		if len(got) != 1 {
			t.Errorf("FrameworkFactoriesDeclaredIn(data source impl) = %+v, want the sae data source entry alone: the ecs entry names a factory ecs declares itself", got)
		}
	})
	t.Run("a product file without its register fails closed", func(t *testing.T) {
		t.Chdir(t.TempDir())
		dir := filepath.Join("alicloud", "service", "sae")
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		implSrc := "package sae\n\nfunc NewLoadBalancerNLBDataSource() *nlbDataSource { return nil }\n"
		if err := os.WriteFile(impl, []byte(implSrc), 0644); err != nil {
			t.Fatal(err)
		}

		got, err := FrameworkFactoriesDeclaredIn(impl, []byte(implSrc))
		if err == nil {
			t.Fatalf("FrameworkFactoriesDeclaredIn(file without register) = (%+v, nil), want an error: the guard cannot prove what the file declares", got)
		}
	})
	t.Run("an unparseable source fails closed", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameTree(t, sameNameResourceFirstRegisterSrc)

		got, err := FrameworkFactoriesDeclaredIn(impl, []byte("package sae\n\nfunc ( {\n"))
		if err == nil {
			t.Fatalf("FrameworkFactoriesDeclaredIn(unparseable src) = (%+v, nil), want an error: the guard cannot prove what the file declares", got)
		}
	})
	t.Run("an unparseable register fails closed, not as no entries", func(t *testing.T) {
		t.Chdir(t.TempDir())
		// A register that fails to parse must surface as an error: the guard
		// counts factories off this result, and folding the parse failure
		// into "no entries" would turn the file-granularity guard into a
		// silent pass — exactly the fail-closed contract its error return
		// promises.
		writeSameNameTree(t, "package sae\n\nfunc ( {\n")

		src, err := os.ReadFile(impl)
		if err != nil {
			t.Fatal(err)
		}
		got, err := FrameworkFactoriesDeclaredIn(impl, src)
		if err == nil {
			t.Fatalf("FrameworkFactoriesDeclaredIn(unparseable register) = (%+v, nil), want an error: a broken register cannot prove what the file declares", got)
		}
	})
	t.Run("alias registrations return every entry that shares the factory", func(t *testing.T) {
		t.Chdir(t.TempDir())
		// Alias registrations reuse one factory for several type names
		// (alicloud_disk/alicloud_ecs_disk, the sslcertificateservice legacy
		// aliases): every entry whose factory this file declares comes back,
		// and the caller counts distinct factories rather than entries.
		aliasRegister := `package sae

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sae",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_sae_load_balancer_nlb",
				Factory:  NewLoadBalancerNLBDataSource,
			},
			{
				TypeName: "alicloud_sae_load_balancer_nlb_alias",
				Factory:  NewLoadBalancerNLBDataSource,
			},
		},
	}
}
`
		writeSameNameTree(t, aliasRegister)

		src, err := os.ReadFile(impl)
		if err != nil {
			t.Fatal(err)
		}
		got, err := FrameworkFactoriesDeclaredIn(impl, src)
		if err != nil {
			t.Fatalf("FrameworkFactoriesDeclaredIn(data source impl) errored: %s", err)
		}
		if len(got) != 2 {
			t.Fatalf("FrameworkFactoriesDeclaredIn(data source impl) = %+v, want both entries that share the one factory", got)
		}
		for _, e := range got {
			if e.Factory != "NewLoadBalancerNLBDataSource" {
				t.Errorf("entry %+v does not name the factory this file declares", e)
			}
		}
	})
	t.Run("a file outside the product directories scans no register", func(t *testing.T) {
		t.Chdir(t.TempDir())
		writeSameNameTree(t, sameNameResourceFirstRegisterSrc)

		got, err := FrameworkFactoriesDeclaredIn("alicloud/resource_widget.go", []byte("package alicloud\n\nfunc NewWidgetResource() *widget { return nil }\n"))
		if err != nil {
			t.Fatalf("FrameworkFactoriesDeclaredIn(flat file) errored: %s", err)
		}
		if len(got) != 0 {
			t.Errorf("FrameworkFactoriesDeclaredIn(flat file) = %+v, want empty: no registered framework factory can be declared outside a product package", got)
		}
	})
}

// TestUnitFrameworkFindByTypeNamePrefersFirstDeclaration pins the name-only
// lookup's behavior: it returns whichever same-named entry the register
// declares first, which is why the kind-aware caller must not fall back to it
// when the pair is registered.
func TestUnitFrameworkFindByTypeNamePrefersFirstDeclaration(t *testing.T) {
	const name = "alicloud_sae_load_balancer_nlb"
	for _, tc := range []struct {
		desc     string
		register string
		want     Kind
	}{
		{"resource declared first", sameNameResourceFirstRegisterSrc, KindResource},
		{"data source declared first", sameNameDataSourceFirstRegisterSrc, KindDataSource},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Chdir(t.TempDir())
			writeSameNameTree(t, tc.register)
			r := FrameworkFindByTypeName(name)
			if r == nil {
				t.Fatal("FrameworkFindByTypeName(name) = nil, want the first-declared registration")
			}
			if r.Kind != tc.want {
				t.Errorf("FrameworkFindByTypeName(name) resolved to kind %s, want %s", r.Kind, tc.want)
			}
		})
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

func TestProviderSourcePaths(t *testing.T) {
	want := []string{"alicloud/provider/sdkv2/provider.go", "alicloud/provider.go"}
	if got := ProviderSourcePaths(); !reflect.DeepEqual(got, want) {
		t.Fatalf("ProviderSourcePaths() = %v, want %v", got, want)
	}
	for _, path := range want {
		if !IsProviderSource(path) {
			t.Errorf("IsProviderSource(%s) = false, want true", path)
		}
	}
	if IsProviderSource("alicloud/provider/sdkv2/assemble_provider.go") {
		t.Error("IsProviderSource accepted a non-provider implementation file")
	}
}

func TestReadProviderSourceFrom(t *testing.T) {
	errMissing := errors.New("missing")
	t.Run("current path wins", func(t *testing.T) {
		var reads []string
		got, err := ReadProviderSourceFrom(func(path string) ([]byte, error) {
			reads = append(reads, path)
			return []byte("current"), nil
		})
		if err != nil || string(got) != "current" || !reflect.DeepEqual(reads, []string{ProviderSourcePath}) {
			t.Fatalf("got %q, reads=%v, err=%v", got, reads, err)
		}
	})
	t.Run("historical base falls back", func(t *testing.T) {
		var reads []string
		got, err := ReadProviderSourceFrom(func(path string) ([]byte, error) {
			reads = append(reads, path)
			if path == ProviderSourcePath {
				return nil, errMissing
			}
			return []byte("legacy"), nil
		})
		if err != nil || string(got) != "legacy" || !reflect.DeepEqual(reads, ProviderSourcePaths()) {
			t.Fatalf("got %q, reads=%v, err=%v", got, reads, err)
		}
	})
	t.Run("missing source returns error", func(t *testing.T) {
		if _, err := ReadProviderSourceFrom(func(string) ([]byte, error) { return nil, errMissing }); !errors.Is(err, errMissing) {
			t.Fatalf("ReadProviderSourceFrom error = %v, want %v", err, errMissing)
		}
	})
}

func TestReadProviderSource(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.MkdirAll(filepath.Dir(ProviderSourcePath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(LegacyProviderSourcePath, []byte("legacy"), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadProviderSource()
	if err != nil || string(got) != "legacy" {
		t.Fatalf("legacy ReadProviderSource() = %q, %v", got, err)
	}
	if err := os.WriteFile(ProviderSourcePath, []byte("current"), 0644); err != nil {
		t.Fatal(err)
	}
	got, err = ReadProviderSource()
	if err != nil || string(got) != "current" {
		t.Fatalf("current ReadProviderSource() = %q, %v", got, err)
	}
}

package alicloud

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudFCCustomDomainsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.FCCustomDomainSupportRegions)
	rand := acctest.RandInt()
	resourceId := "data.alicloud_fc_custom_domains.default"
	name := fmt.Sprintf("tf-testacc-fc-custom-domains-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceFCCustomDomainsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_fc_custom_domain.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_fc_custom_domain.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.domain_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.domain_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.domain_name}",
			"ids":        []string{"${alicloud_fc_custom_domain.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.domain_name}_fake",
			"ids":        []string{"${alicloud_fc_custom_domain.default.id}_fake"},
		}),
	}

	var existFCCustomDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"domains.#":                              "1",
			"ids.#":                                  "1",
			"names.#":                                "1",
			"domains.0.id":                           "alicloud-provider.shop",
			"domains.0.domain_name":                  "alicloud-provider.shop",
			"domains.0.protocol":                     "HTTP,HTTPS",
			"domains.0.account_id":                   CHECKSET,
			"domains.0.api_version":                  CHECKSET,
			"domains.0.created_time":                 CHECKSET,
			"domains.0.last_modified_time":           CHECKSET,
			"domains.0.route_config.0.path":          "/*",
			"domains.0.route_config.0.service_name":  name,
			"domains.0.route_config.0.function_name": name,
			"domains.0.route_config.0.qualifier":     "v1",
			"domains.0.route_config.0.methods.0":     "GET",
			"domains.0.route_config.0.methods.1":     "POST",
			"domains.0.cert_config.0.cert_name":      "test",
			"domains.0.cert_config.0.certificate":    strings.Replace(testFcCertificate, `\n`, "\n", -1),
		}
	}

	var fakeFCCustomDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"domains.#": "0",
			"ids.#":     "0",
			"names.#":   "0",
		}
	}

	var fcCustomDomainsRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existFCCustomDomainsMapFunc,
		fakeMapFunc:  fakeFCCustomDomainsMapFunc,
	}

	fcCustomDomainsRecordsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceFCCustomDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_fc_custom_domain" "default" {
	domain_name = "alicloud-provider.shop"
	protocol = "HTTP,HTTPS"
	route_config {
		path = "/*"
		service_name = "${alicloud_fc_service.default.name}"
		function_name = "${alicloud_fc_function.default.name}"
		qualifier = "v1"
		methods = ["GET","POST"]
	}
	cert_config {
		cert_name = "test"
		private_key = "%s"
		certificate = "%s"
	}
}

resource "alicloud_fc_service" "default" {
    name = "${var.name}"
    description = "${var.name}-description"
}

resource "alicloud_oss_bucket" "default" {
	bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "default" {
	bucket = "${alicloud_oss_bucket.default.id}"
	key = "fc/hello.zip"
	content = <<EOF
		# -*- coding: utf-8 -*-
	def handler(event, context):
		print("hello world")
		return 'hello world'
	EOF
}

resource "alicloud_fc_function" "default" {
	service = "${alicloud_fc_service.default.name}"
	name = "${var.name}"
	oss_bucket = "${alicloud_oss_bucket.default.id}"
	oss_key = "${alicloud_oss_bucket_object.default.key}"
	memory_size = 512
	runtime = "python3.10"
	handler = "hello.handler"
}
`, name, testFcPrivateKey, testFcCertificate)
}

// TestUnitAlicloudFCCustomDomainsDataSourceNilFields verifies the data source
// Read does not panic when the FC ListCustomDomains response omits domain-level
// and route pointer string fields — the root cause of the nil pointer panic in
// GitHub issue #9166. It also asserts empty fields surface as empty strings and
// that the existing domains/ids/names output semantics are preserved.
//
// The test is network-free: gomonkey patches AliyunClient.WithFcClient to return
// a forged *fc.ListCustomDomainsOutput. The CustomDomains slice element type is
// the SDK's unexported customDomainMetadata, so the two test elements are built
// via reflection; the exported RouteConfig/PathConfig values are constructed
// normally.
func TestUnitAlicloudFCCustomDomainsDataSourceNilFields(t *testing.T) {
	// fc-go-sdk defines ListCustomDomainsOutput.CustomDomains as
	// []*customDomainMetadata where customDomainMetadata is unexported. Build the
	// two test elements via reflection so the test compiles without naming that
	// type, while constructing the exported RouteConfig/PathConfig values directly.
	output := &fc.ListCustomDomainsOutput{}
	outVal := reflect.ValueOf(output).Elem()
	cdField := outVal.FieldByName("CustomDomains") // []*customDomainMetadata
	elemPtrType := cdField.Type().Elem()           // *customDomainMetadata
	elemType := elemPtrType.Elem()                 // customDomainMetadata

	// First domain: every domain-level pointer field is nil and both RouteConfig
	// and CertConfig are nil. Pre-fix this triggered the panic at *domain.DomainName,
	// *domain.AccountID, *domain.Protocol and *domain.APIVersion.
	allNilDomain := reflect.New(elemType) // *customDomainMetadata with all nil fields

	// Second domain: DomainName/Protocol set, RouteConfig present but its PathConfig
	// has all-nil pointer fields. Pre-fix this triggered the panic at *v.Path,
	// *v.ServiceName, *v.FunctionName and *v.Qualifier.
	routeNilFieldsDomain := reflect.New(elemType)
	rnfd := routeNilFieldsDomain.Elem()
	foo := "foo.com"
	rnfd.FieldByName("DomainName").Set(reflect.ValueOf(&foo))
	httpProto := "HTTP"
	rnfd.FieldByName("Protocol").Set(reflect.ValueOf(&httpProto))
	rnfd.FieldByName("RouteConfig").Set(reflect.ValueOf(&fc.RouteConfig{
		Routes: []fc.PathConfig{{}},
	}))

	customDomains := reflect.MakeSlice(cdField.Type(), 2, 2)
	customDomains.Index(0).Set(allNilDomain)
	customDomains.Index(1).Set(routeNilFieldsDomain)
	cdField.Set(customDomains)

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "WithFcClient",
		func(_ *connectivity.AliyunClient, _ func(*fc.Client) (interface{}, error)) (interface{}, error) {
			return output, nil
		})

	d, _ := schema.InternalMap(dataSourceAlicloudFcCustomDomains().Schema).Data(nil, nil)
	meta := &connectivity.AliyunClient{}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("dataSourceAlicloudFcCustomDomainsRead panicked on nil pointer fields: %v", r)
		}
	}()
	if err := dataSourceAlicloudFcCustomDomainsRead(d, meta); err != nil {
		t.Fatalf("dataSourceAlicloudFcCustomDomainsRead returned unexpected error: %v", err)
	}

	domains := d.Get("domains").([]interface{})
	if len(domains) != 2 {
		t.Fatalf("expected 2 domains, got %d", len(domains))
	}

	// AC1: domain-level nil pointer fields surface as empty strings; nil
	// RouteConfig / CertConfig yield empty lists instead of a panic.
	first := domains[0].(map[string]interface{})
	for _, key := range []string{"id", "domain_name", "account_id", "protocol", "api_version"} {
		if v := first[key].(string); v != "" {
			t.Errorf("domain[0].%s: expected empty string for nil pointer, got %q", key, v)
		}
	}
	if rc := first["route_config"].([]interface{}); len(rc) != 0 {
		t.Errorf("domain[0].route_config: expected empty list for nil RouteConfig, got %d", len(rc))
	}
	if cc := first["cert_config"].([]interface{}); len(cc) != 0 {
		t.Errorf("domain[0].cert_config: expected empty list for nil CertConfig, got %d", len(cc))
	}

	// AC1 (route fields): DomainName preserved; nil route pointer fields surface
	// as empty strings with no panic.
	second := domains[1].(map[string]interface{})
	if v := second["domain_name"].(string); v != "foo.com" {
		t.Errorf("domain[1].domain_name: expected foo.com, got %q", v)
	}
	rc2 := second["route_config"].([]interface{})
	if len(rc2) != 1 {
		t.Fatalf("domain[1].route_config: expected 1 route, got %d", len(rc2))
	}
	route := rc2[0].(map[string]interface{})
	for _, key := range []string{"path", "service_name", "function_name", "qualifier"} {
		if v := route[key].(string); v != "" {
			t.Errorf("route[0].%s: expected empty string for nil pointer, got %q", key, v)
		}
	}

	// AC2: ids/names output semantics preserved — with no name_regex/ids filter
	// set, both domains survive, empty domain name first then foo.com.
	ids := d.Get("ids").([]interface{})
	names := d.Get("names").([]interface{})
	if len(ids) != 2 || len(names) != 2 {
		t.Fatalf("expected 2 ids and 2 names, got ids=%d names=%d", len(ids), len(names))
	}
	if names[1].(string) != "foo.com" {
		t.Errorf("names[1]: expected foo.com, got %q", names[1])
	}
	if ids[1].(string) != "foo.com" {
		t.Errorf("ids[1]: expected foo.com, got %q", ids[1])
	}
}

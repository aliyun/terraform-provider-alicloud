package fcv3_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest/scoped"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/account"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fcv3"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sslcertificateservice"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
)

// Dependencies cover the resource and data source fixtures across the FCV3
// suite. Keep the selection here rather than repeating it in individual cases.
var fcSuite = scoped.New(
	fcv3.ServicePackage(), account.ServicePackage(), ecs.ServicePackage(),
	fc.ServicePackage(), sls.ServicePackage(), nas.ServicePackage(),
	oss.ServicePackage(), resourcemanager.ServicePackage(),
	sslcertificateservice.ServicePackage(), vpc.ServicePackage(),
)

// TestUnitFCV3ScopedDependencies prevents an innocent shared-helper import from
// silently restoring the complete provider compilation graph.
func TestUnitFCV3ScopedDependencies(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", "-test", ".")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("list test dependencies: %v", err)
	}
	forbidden := map[string]bool{
		"github.com/aliyun/terraform-provider-alicloud/alicloud":                    true,
		"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest":            true,
		"github.com/aliyun/terraform-provider-alicloud/alicloud/provider":           true,
		"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry":  true,
		"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/framework": true,
	}
	for _, dependency := range strings.Fields(string(output)) {
		if forbidden[dependency] {
			t.Errorf("FCV3 test imports full provider entry %s", dependency)
		}
	}
}

func TestUnitFCV3ScopedFixtureCoverage(t *testing.T) {
	factory, err := fcSuite.ProtoV5ProviderFactories()["alicloud"]()
	if err != nil {
		t.Fatal(err)
	}
	response, err := factory.GetProviderSchema(context.Background(), &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range response.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("provider schema: %s: %s", d.Summary, d.Detail)
		}
	}
	files, err := filepath.Glob("*_test.go")
	if err != nil {
		t.Fatal(err)
	}
	declaration := regexp.MustCompile(`(?m)^\s*(resource|data)\s+"(alicloud_[a-z0-9_]+)"\s+"`)
	for _, name := range files {
		source, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range declaration.FindAllStringSubmatch(string(source), -1) {
			registered := response.ResourceSchemas[match[2]] != nil
			if match[1] == "data" {
				registered = response.DataSourceSchemas[match[2]] != nil
			}
			if !registered {
				t.Errorf("%s: missing %s %s in the shared test provider", name, match[1], match[2])
			}
		}
	}
}

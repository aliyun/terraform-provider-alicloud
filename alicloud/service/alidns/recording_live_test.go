package alidns_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// These opt-in cloud pilots retain the production Provider Configure and real
// registered resources. Existing, unmodified HashiCorp randomness demonstrates
// the process-wide seed: resource test authors need no recording-specific API.
func runRecordingLiveLifecycle(t *testing.T, resourceType, nameAttribute string, definition func(string, string) string) {
	t.Helper()
	if os.Getenv("TF_ACC") == "" || !recording.Enabled() {
		t.Skip("requires TF_ACC=1 and record/replay mode")
	}
	terraform, err := exec.LookPath("terraform")
	if err != nil {
		t.Fatal("installed Terraform CLI required")
	}
	terraform, err = filepath.Abs(terraform)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("TF_ACC_TERRAFORM_PATH", terraform)
	t.Setenv("CHECKPOINT_DISABLE", "1")
	name := "tf-rr-" + sdkacctest.RandString(12)
	t.Logf("Recording resource: %s name=%s", resourceType, name)
	if resourceType == "alicloud_alidns_domain_group" {
		t.Logf("Recording DNS group: %s (updated name: %s-u)", name, name)
	}
	factory := func() (tfprotov5.ProviderServer, error) {
		p := alicloud.Provider()
		p.ResourcesMap = map[string]*schema.Resource{resourceType: p.ResourcesMap[resourceType]}
		p.DataSourcesMap = nil
		return schema.NewGRPCProviderServer(p), nil
	}
	config := func(phase string) string {
		return fmt.Sprintf(`
provider "alicloud" { region = "cn-hangzhou" }
resource %q "test" {
%s
}
`, resourceType, definition(name, phase))
	}
	address := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){"alicloud": factory},
		Steps: []resource.TestStep{
			{Config: config("initial"), Check: resource.TestCheckResourceAttr(address, nameAttribute, name)},
			{ResourceName: address, ImportState: true, ImportStateVerify: true},
			{Config: config("updated"), Check: resource.TestCheckResourceAttrSet(address, "id")},
		},
	})
}

func TestRecordingAlidnsGroupLiveTerraformCLI(t *testing.T) {
	runRecordingLiveLifecycle(t, "alicloud_alidns_domain_group", "domain_group_name", func(name, phase string) string {
		if phase == "updated" {
			name += "-u"
		}
		return fmt.Sprintf("domain_group_name = %q", name)
	})
}

func TestRecordingVpcLiveTerraformCLI(t *testing.T) {
	runRecordingLiveLifecycle(t, "alicloud_vpc", "vpc_name", func(name, phase string) string {
		return fmt.Sprintf("vpc_name = %q\ncidr_block = %q\ndescription = %q", name, "192.168.0.0/24", "recording pilot "+phase)
	})
}

func TestRecordingSaeNamespaceLiveTerraformCLI(t *testing.T) {
	runRecordingLiveLifecycle(t, "alicloud_sae_namespace", "namespace_name", func(name, phase string) string {
		return fmt.Sprintf("namespace_name = %q\nnamespace_short_id = %q\nnamespace_description = %q\nenable_micro_registration = false", name, strings.ReplaceAll(name, "-", ""), "recording pilot "+phase)
	})
}

func TestRecordingSlsProjectLiveTerraformCLI(t *testing.T) {
	runRecordingLiveLifecycle(t, "alicloud_log_project", "project_name", func(name, phase string) string {
		return fmt.Sprintf("project_name = %q\ndescription = %q", name, "recording pilot "+phase)
	})
}

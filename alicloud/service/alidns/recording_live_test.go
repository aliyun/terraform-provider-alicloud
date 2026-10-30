package alidns_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestRecordingAlidnsGroupLiveTerraformCLI is an explicit opt-in cloud pilot. The same
// binary runs record and externally network-isolated replay in separate processes.
// Unlike the DNS fixtures, it retains the production Provider Configure function.
func TestRecordingAlidnsGroupLiveTerraformCLI(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || !recording.Enabled() {
		t.Skip("requires TF_ACC=1 and record/replay mode")
	}
	// Draw twice from the saved session seed: unique in record, identical in replay.
	// Two base-36 draws keep both the name and its update within the 20-char limit.
	name := "tf-rr-" + strconv.FormatInt(int64(acctest.RecordingRandInt()), 36) + strconv.FormatInt(int64(acctest.RecordingRandInt()), 36)
	t.Logf("Recording DNS group: %s (updated name: %s-u)", name, name)
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
	factory := func() (tfprotov5.ProviderServer, error) {
		p := alicloud.Provider()
		p.ResourcesMap = map[string]*schema.Resource{"alicloud_alidns_domain_group": p.ResourcesMap["alicloud_alidns_domain_group"]}
		p.DataSourcesMap = nil
		return schema.NewGRPCProviderServer(p), nil
	}
	config := func(suffix string) string {
		return fmt.Sprintf(`
provider "alicloud" {
 region = "cn-hangzhou"
 account_type = "Domestic"
}
resource "alicloud_alidns_domain_group" "test" {
 domain_group_name = %q
}
`, name+suffix)
	}
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){"alicloud": factory},
		Steps: []resource.TestStep{
			{Config: config(""), Check: resource.TestCheckResourceAttr("alicloud_alidns_domain_group.test", "domain_group_name", name)},
			{ResourceName: "alicloud_alidns_domain_group.test", ImportState: true, ImportStateVerify: true},
			{Config: config("-u"), Check: resource.TestCheckResourceAttr("alicloud_alidns_domain_group.test", "domain_group_name", name+"-u")},
		},
	})
}

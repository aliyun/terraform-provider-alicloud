package itag_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The iTAG LabelUser lives under a pre-existing iTAG tenant, which no Terraform
// resource creates. The tenant id is supplied via ALICLOUD_ITAG_TENANT_ID; the
// test is skipped when it is absent, so it only runs against an account that
// has iTAG provisioned.
//
// The test uses the terraform-plugin-sdk/v2 helper/resource harness (rather than
// terraform-plugin-testing/helper/resource) because package alicloud already
// links helper/resource transitively (e.g. the polardb polar_fs resources), and
// linking both helper/resource packages redefines the -sweep flag at init and
// panics before the test runs. acctest.ProtoV5ProviderFactories() serves the
// real muxed provider under either harness; the sdk harness keeps a single
// helper/resource in the test binary so the -sweep flag is registered once.
func testAccItagLabelUserPreCheck(t *testing.T) {
	acctest.PreCheck(t)
	if os.Getenv("ALICLOUD_ITAG_TENANT_ID") == "" {
		t.Skip("ALICLOUD_ITAG_TENANT_ID must be set to run iTAG LabelUser acceptance tests")
	}
}

func TestAccAliCloudItagLabelUser(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)
	tenantId := os.Getenv("ALICLOUD_ITAG_TENANT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccItagLabelUserPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// sdk helper/resource.TestCheckFunc takes the sdk terraform.State, so
		// the testing-harness acctest.CheckDestroyNoop (plugin-testing State)
		// does not assign; a local noop keeps the sdk-harness build consistent.
		CheckDestroy: func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config: testAccItagLabelUserConfig(r, tenantId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_itag_label_user.default", "id"),
					resource.TestCheckResourceAttr("alicloud_itag_label_user.default", "tenant_id", tenantId),
					resource.TestCheckResourceAttrSet("alicloud_itag_label_user.default", "account_no"),
					resource.TestCheckResourceAttr("alicloud_itag_label_user.default", "account_type", "ALIYUN"),
					resource.TestCheckResourceAttr("alicloud_itag_label_user.default", "role", "OPERATOR"),
					resource.TestCheckResourceAttr("alicloud_itag_label_user.default", "user_name", fmt.Sprintf("tf-testacc-itag-user-%d", r)),
					resource.TestCheckResourceAttrSet("alicloud_itag_label_user.default", "user_id"),
				),
			},
			{
				Config: testAccItagLabelUserConfigUpdate(r, tenantId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_itag_label_user.default", "role", "ADMIN"),
					resource.TestCheckResourceAttr("alicloud_itag_label_user.default", "user_name", fmt.Sprintf("tf-testacc-itag-user-update-%d", r)),
				),
			},
			{
				ResourceName:      "alicloud_itag_label_user.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccItagLabelUserConfig(r int, tenantId string) string {
	return fmt.Sprintf(`
resource "alicloud_ram_user" "test" {
  name         = "tf-testacc-itag-%[1]d"
  display_name = "tf-testacc-itag-%[1]d"
}

resource "alicloud_itag_label_user" "default" {
  tenant_id    = "%[2]s"
  account_no   = alicloud_ram_user.test.unique_id
  account_type = "ALIYUN"
  role         = "OPERATOR"
  user_name    = "tf-testacc-itag-user-%[1]d"
}
`, r, tenantId)
}

func testAccItagLabelUserConfigUpdate(r int, tenantId string) string {
	return fmt.Sprintf(`
resource "alicloud_ram_user" "test" {
  name         = "tf-testacc-itag-%[1]d"
  display_name = "tf-testacc-itag-%[1]d"
}

resource "alicloud_itag_label_user" "default" {
  tenant_id    = "%[2]s"
  account_no   = alicloud_ram_user.test.unique_id
  account_type = "ALIYUN"
  role         = "ADMIN"
  user_name    = "tf-testacc-itag-user-update-%[1]d"
}
`, r, tenantId)
}

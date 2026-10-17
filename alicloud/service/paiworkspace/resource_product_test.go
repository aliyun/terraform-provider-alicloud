package paiworkspace_test

import (
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Matches the resource's fixed Create diagnostic summary
// ("Creating alicloud_pai_workspace_product: calling CreateProductOrders"),
// so the step passes whenever the invalid product_code is rejected at Create,
// regardless of the specific server error message.
var testAccPaiWorkspaceProductCreateError = regexp.MustCompile(
	`alicloud_pai_workspace_product.*CreateProductOrders`,
)

// The PAI Workspace Product resource places a real, chargeable order on
// CreateProductOrders (the success response carries an OrderId). To honour the
// no-real-paid-order test constraint, the Create success/failure, Read
// found/not-found and Delete no-op CRUD paths are covered by the no-cloud
// httptest unit tests in resource_product_internal_test.go, which exercise the
// real ROA wire layer against a mock server under ci-check-quick.
//
// The one real-cloud AccTest here exercises the Create *failure* path: an
// invalid product_code makes CreateProductOrders reject the order before any
// purchase (no OrderId, no charge), so `terraform apply` fails at Create. This
// is the "failure" path the no-real-order constraint explicitly permits for the
// order-create step. The resource's Create diagnostic summary is fixed text, so
// ExpectError matches deterministically regardless of the specific server error.
func TestAccAliCloudPaiWorkspaceProductCreateInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudPaiWorkspaceProductInvalidConfig(),
				ExpectError: testAccPaiWorkspaceProductCreateError,
			},
		},
	})
}

func testAccAliCloudPaiWorkspaceProductInvalidConfig() string {
	return `
resource "alicloud_pai_workspace_product" "default" {
  product_code = "INVALID_NOT_A_REAL_PRODUCT_CODE"
  auto_pay     = false
}
`
}

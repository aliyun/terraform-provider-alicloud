package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudCrDiagnosisTask_basic covers the create-read-import lifecycle
// of alicloud_cr_diagnosis_task. The resource has no deletion API, so destroy
// is a no-op (the remote task may remain); CheckDestroy uses CheckDestroyNoop,
// mirroring the v1 no-delete pattern (e.g. alicloud_aligreen_oss_stock_task).
//
// status is a volatile, asynchronously-transitioning attribute
// (PENDING/RUNNING -> SUCCESS/FAILED): a diagnosis task completes shortly
// after creation, so the value observed at apply may differ from the value
// observed at a later refresh. IDRefreshName (which requires the refreshed
// state to be equivalent to the applied state) is therefore omitted, and
// status is excluded from ImportStateVerify for the same reason.
func TestAccAliCloudCrDiagnosisTask_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-crdt-%d", rand)
	resourceId := "alicloud_cr_diagnosis_task.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCrDiagnosisTaskBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "diagnosis_task_id"),
					resource.TestCheckResourceAttrSet(resourceId, "instance_id"),
					resource.TestCheckResourceAttr(resourceId, "diagnosis_type", "INSTANCE_DIAGNOSIS"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

func testAccAliCloudCrDiagnosisTaskBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type    = "Subscription"
  period          = 1
  renewal_status  = "ManualRenewal"
  instance_type   = "Basic"
  instance_name   = var.name
}

resource "alicloud_cr_diagnosis_task" "default" {
  instance_id    = alicloud_cr_ee_instance.default.id
  diagnosis_type = "INSTANCE_DIAGNOSIS"
}
`, name)
}

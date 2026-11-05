package threatdetection_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccAliCloudThreatDetectionServiceLinkedRole_basic(t *testing.T) {
	rac, testAccCheck := testAccThreatDetectionServiceLinkedRoleDependence()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccThreatDetectionServiceLinkedRoleConfig(threatDetectionServiceLinkedRoleOssMfd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_linked_role": threatDetectionServiceLinkedRoleOssMfd,
					}),
				),
			},
			{
				ResourceName:      threatDetectionServiceLinkedRoleId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudThreatDetectionServiceLinkedRole_replace(t *testing.T) {
	rac, testAccCheck := testAccThreatDetectionServiceLinkedRoleDependence()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccThreatDetectionServiceLinkedRoleConfig(threatDetectionServiceLinkedRoleOssMfd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_linked_role": threatDetectionServiceLinkedRoleOssMfd,
					}),
				),
			},
			{
				Config: testAccThreatDetectionServiceLinkedRoleConfig(threatDetectionServiceLinkedRoleMssp),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(threatDetectionServiceLinkedRoleId, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_linked_role": threatDetectionServiceLinkedRoleMssp,
					}),
				),
			},
			{
				ResourceName:      threatDetectionServiceLinkedRoleId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudThreatDetectionServiceLinkedRole_alreadyExists(t *testing.T) {
	rac, testAccCheck := testAccThreatDetectionServiceLinkedRoleDependence()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccThreatDetectionServiceLinkedRoleConfig(threatDetectionServiceLinkedRoleMssp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_linked_role": threatDetectionServiceLinkedRoleMssp,
					}),
				),
			},
			{
				Config:      testAccThreatDetectionServiceLinkedRoleDuplicateConfig(threatDetectionServiceLinkedRoleMssp),
				ExpectError: regexp.MustCompile("already exists, import it with: terraform import alicloud_threat_detection_service_linked_role"),
			},
		},
	})
}

const (
	threatDetectionServiceLinkedRoleId     = "alicloud_threat_detection_service_linked_role.default"
	threatDetectionServiceLinkedRoleMssp   = "AliyunServiceRoleForAntiRansomwareMssp"
	threatDetectionServiceLinkedRoleOssMfd = "AliyunServiceRoleForOssMfd"
)

var aliCloudThreatDetectionServiceLinkedRoleMap = map[string]string{
	"role_status": acctest.CHECKSET,
}

func testAccThreatDetectionServiceLinkedRoleDependence() (*acctest.ResourceAttrCheck, acctest.ResourceAttrMapUpdate) {
	var v map[string]interface{}
	ra := acctest.ResourceAttrInit(threatDetectionServiceLinkedRoleId, aliCloudThreatDetectionServiceLinkedRoleMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(threatDetectionServiceLinkedRoleId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionServiceLinkedRole")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	return rac, rac.ResourceAttrMapUpdateSet()
}

func testAccThreatDetectionServiceLinkedRoleConfig(roleName string) string {
	return fmt.Sprintf(`
resource "alicloud_threat_detection_service_linked_role" "default" {
  service_linked_role = "%s"
}
`, roleName)
}

func testAccThreatDetectionServiceLinkedRoleDuplicateConfig(roleName string) string {
	return testAccThreatDetectionServiceLinkedRoleConfig(roleName) + fmt.Sprintf(`
resource "alicloud_threat_detection_service_linked_role" "duplicate" {
  service_linked_role = "%s"
}
`, roleName)
}

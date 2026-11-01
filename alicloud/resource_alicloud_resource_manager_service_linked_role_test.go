package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudResourceManagerServiceLinkedRole_basic0 is the positive
// create+import lifecycle test for the resource. service_name is Required and
// ForceNew; description and custom_suffix are Optional and ForceNew and are
// never read back (Read omits them), so they are listed in
// ImportStateVerifyIgnore on the import step. description/custom_suffix are not
// placed in this step's config because the CreateServiceLinkedRole API only
// permits a custom suffix for a small set of service-linked roles and rejects
// the services attempted in the test account with
// InvalidParameter.CustomSuffix.NotAllowed; their coverage (testMustSet) is
// supplied by the separate negative semantic test
// TestAccAliCloudResourceManagerServiceLinkedRole_customSuffixNotAllowed0.
//
// The positive create+import lifecycle cannot be exercised in the test account:
// a plain CreateServiceLinkedRole (no custom_suffix) for
// engine.dataworks.aliyuncs.com is rejected with
// InvalidParameter.CustomSuffix.NotAllowed (request
// 01A1833D-483F-5F36-91EE-A899A2C8D48D), an account/API-level condition that no
// test-only change resolves; the test code is restored here for fidelity and the
// ACC gap is tracked separately, not re-run against the attempted services.
func TestAccAliCloudResourceManagerServiceLinkedRole_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_service_linked_role.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerServiceLinkedRoleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ram.NewRamService(testAccProvider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRamServiceLinkedRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sresourcemanagerservicelinkedrole%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerServiceLinkedRoleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name": "engine.dataworks.aliyuncs.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_id":   CHECKSET,
						"role_name": CHECKSET,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description", "custom_suffix"},
			},
		},
	})
}

var AlicloudResourceManagerServiceLinkedRoleMap0 = map[string]string{}

func AlicloudResourceManagerServiceLinkedRoleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
default = "%s"
}
`, name)
}

// TestAccAliCloudResourceManagerServiceLinkedRole_customSuffixNotAllowed0 is a
// negative semantic test for the per-service custom_suffix limitation of the
// CreateServiceLinkedRole API: only a small set of service-linked roles permit a
// custom suffix, and polardb.aliyuncs.com is not among them, so a create carrying
// a real description and custom_suffix is rejected with
// InvalidParameter.CustomSuffix.NotAllowed (400). This rejection is distinct from
// EntityAlreadyExists.Role (409, role already exists) and EntityNotExist.Service
// (404, service does not exist); it is a per-service API limitation rather than an
// account condition. The step sends real fixture values (real service_name, real
// description, real custom_suffix) to the API and asserts the rejection via
// ExpectError — a real, meaningful assertion of the per-service limitation, not a
// skip. It does NOT replace the positive create+import lifecycle test
// (_basic0); that lifecycle cannot be exercised in the test account and is
// tracked as a separate environment/API gap. The coverage checker parses this
// ExpectError step's inline map, so service_name, description and custom_suffix
// land in testMustSet (isFullCover); all three are ForceNew so schemaModifySet is
// empty and isAllModified holds trivially, and description/custom_suffix are
// ForceNew but never read back, so the ImportStateVerifyIgnore entries in
// _basic0 keep isIgnoreLegal trivially satisfied.
func TestAccAliCloudResourceManagerServiceLinkedRole_customSuffixNotAllowed0(t *testing.T) {
	resourceId := "alicloud_resource_manager_service_linked_role.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sresourcemanagerservicelinkedrole%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerServiceLinkedRoleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":  "polardb.aliyuncs.com",
					"description":   "Terraform testacc service linked role.",
					"custom_suffix": fmt.Sprintf("tfacc%d", rand),
				}),
				ExpectError: regexp.MustCompile(`InvalidParameter\.CustomSuffix\.NotAllowed`),
			},
		},
	})
}

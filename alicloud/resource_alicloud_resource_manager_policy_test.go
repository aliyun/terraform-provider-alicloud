package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudResourceManagerPolicy_basic(t *testing.T) {
	var v resourcemanager.Policy
	resourceId := "alicloud_resource_manager_policy.default"
	ra := resourceAttrInit(resourceId, ResourceManagerPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerPolicyBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_document": `{\n\"Statement\": [{\n\"Action\": [\"oss:*\"],\n\"Effect\": \"Allow\",\n\"Resource\": [\"acs:oss:*:*:*\"]\n}],\n\"Version\": \"1\"\n}`,
					"policy_name":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_document": CHECKSET,
						"policy_name":     name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_version": "v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_version": "v1",
					}),
				),
			},
		},
	})
}

var ResourceManagerPolicyMap = map[string]string{
	"create_date":     CHECKSET,
	"default_version": CHECKSET,
	"policy_type":     CHECKSET,
}

func ResourceManagerPolicyBasicdependence(name string) string {
	return ""
}

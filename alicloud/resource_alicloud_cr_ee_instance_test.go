package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCrEEInstance_Basic(t *testing.T) {
	t.Skip("Skipping cr ee instance test case")

	var v *cr_ee.GetInstanceResponse
	resourceId := "alicloud_cr_ee_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEEInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-instance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEEInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "Subscription",
					"period":         "1",
					"renew_period":   "1",
					"renewal_status": "AutoRenewal",
					"instance_type":  "Basic",
					"instance_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         CHECKSET,
						"created_time":   CHECKSET,
						"end_time":       CHECKSET,
						"renew_period":   "1",
						"renewal_status": "AutoRenewal",
						"instance_name":  name,
						"instance_type":  "Basic",
						"payment_type":   "Subscription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "custom_oss_bucket"},
			},
		},
	})
}

func resourceCrEEInstanceConfigDependence(name string) string {
	return ""
}

package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ons_instance", &resource.Sweeper{
		Name: "alicloud_ons_instance",
		F:    testSweepOnsInstance,
		Dependencies: []string{
			"alicloud_ons_topic",
			"alicloud_ons_group",
		},
	})
}

func testSweepOnsInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
		"GID-tf-testAcc",
		"GID_tf-testacc",
		"CID-tf-testAcc",
		"CID_tf-testacc",
	}

	request := ons.CreateOnsInstanceInServiceListRequest()
	request.PreventCache = onsService.GetPreventCache()

	raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsInstanceInServiceList(request)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve ONS instance in service list: %s", err)
	}

	var response *ons.OnsInstanceInServiceListResponse
	response, _ = raw.(*ons.OnsInstanceInServiceListResponse)

	for _, v := range response.Data.InstanceVO {
		name := v.InstanceName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ons instance: %s ", name)
			continue
		}
		log.Printf("[INFO] delete ons instance: %s ", name)

		request := ons.CreateOnsInstanceDeleteRequest()
		request.InstanceId = v.InstanceId
		request.PreventCache = onsService.GetPreventCache()

		_, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsInstanceDelete(request)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete ons instance (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAlicloudOnsInstance_basic(t *testing.T) {
	var v *ons.OnsInstanceBaseInfoResponse
	resourceId := "alicloud_ons_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc%sonsinstancebasic%v.abc", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOnsInstanceConfigDependence)

	var onsInstanceBasicMap = map[string]string{
		"name":   fmt.Sprintf("tf-testacc%sonsinstancebasic%v.abc", defaultRegionToTest, rand),
		"remark": "default remark",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testacc%sonsinstancebasic%v.abc", defaultRegionToTest, rand),
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
					"name": "tf-testacc-instance-name-change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": "tf-testacc-instance-name-change"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "updated remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"remark": "updated remark"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name":   "${var.name}",
					"remark": "default remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(onsInstanceBasicMap),
				),
			},
		},
	})

}

func resourceOnsInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%v"
}
`, name)
}

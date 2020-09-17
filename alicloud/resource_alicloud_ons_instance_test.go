package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
	var v ons.InstanceBaseInfo
	resourceId := "alicloud_ons_instance.default"
	ra := resourceAttrInit(resourceId, OnsInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sOnsInstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, OnsInstanceBasicdependence)
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
					"instance_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
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
					"instance_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "Test-Remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "Test-Remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}",
					"remark":        "Test-Remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"remark":        "Test-Remark",
					}),
				),
			},
		},
	})
}

var OnsInstanceMap = map[string]string{
	"instance_type": CHECKSET,
	"release_time":  CHECKSET,
	"status":        CHECKSET,
}

func OnsInstanceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}

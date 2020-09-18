package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ons_group", &resource.Sweeper{
		Name: "alicloud_ons_group",
		F:    testSweepOnsGroup,
	})
}

func testSweepOnsGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	prefixes := []string{
		"GID-tf-testAcc",
		"GID_tf-testacc",
		"CID-tf-testAcc",
		"CID_tf-testacc",
	}

	instanceListReq := ons.CreateOnsInstanceInServiceListRequest()

	raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsInstanceInServiceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve ons instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*ons.OnsInstanceInServiceListResponse)

	var instanceIds []string
	for _, v := range instanceListResp.Data.InstanceVO {
		instanceIds = append(instanceIds, v.InstanceId)
	}

	for _, instanceId := range instanceIds {
		request := ons.CreateOnsGroupListRequest()
		request.InstanceId = instanceId

		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsGroupList(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve ons groups on instance (%s): %s", instanceId, err)
			continue
		}

		groupListResp, _ := raw.(*ons.OnsGroupListResponse)
		groups := groupListResp.Data.SubscribeInfoDo

		for _, v := range groups {
			groupId := v.GroupId
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(groupId), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ons group: %s ", groupId)
				continue
			}
			log.Printf("[INFO] delete ons group: %s ", groupId)

			request := ons.CreateOnsGroupDeleteRequest()
			request.InstanceId = instanceId
			request.GroupId = v.GroupId

			_, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
				return onsClient.OnsGroupDelete(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete ons group (%s): %s", groupId, err)
			}
		}
	}

	return nil
}

func TestAccAlicloudOnsGroup_basic(t *testing.T) {
	var v ons.SubscribeInfoDo
	resourceId := "alicloud_ons_group.default"
	ra := resourceAttrInit(resourceId, onsGroupBasicMap)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOnsGroupConfigDependence)

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
					"instance_id": "${alicloud_ons_instance.default.id}",
					"group_id":    "${var.group_id}",
					"remark":      "alicloud_ons_group_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand),
						"remark":   "alicloud_ons_group_remark",
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
					"tags": map[string]string{
						"Created": "TFM",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.Created": "TFM",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"group_id": "${var.group_id}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v_change", rand)}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"group_id": "${var.group_id}",
					"remark":   "alicloud_ons_group_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand),
						"remark":   "alicloud_ons_group_remark",
					}),
				),
			},
		},
	})

}

func resourceOnsGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ons_instance" "default" {
  name = "%s"
}

variable "group_id" {
 default = "%s"
}
`, name, name)
}

var onsGroupBasicMap = map[string]string{
	"group_id": "${var.group_id}",
	"remark":   "alicloud_ons_group_remark",
}

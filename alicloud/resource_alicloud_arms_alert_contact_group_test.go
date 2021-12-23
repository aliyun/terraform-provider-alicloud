package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_arms_alert_contact_group", &resource.Sweeper{
		Name: "alicloud_arms_alert_contact_group",
		F:    testSweepArmsAlertContactGroup,
	})
}

func testSweepArmsAlertContactGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	action := "SearchAlertContactGroup"
	request := make(map[string]interface{})
	request["IsDetail"] = false
	request["RegionId"] = client.RegionId
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		log.Printf("[ERROR] %s failed error: %v", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.ContactGroups", response)
	if err != nil {
		log.Printf("[ERROR] %s error: %v", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["ContactGroupName"])
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping arms alert contact group: %s ", name)
			continue
		}
		log.Printf("[INFO] delete arms alert contact group: %s ", name)
		action = "DeleteAlertContactGroup"
		request = map[string]interface{}{
			"ContactGroupId": fmt.Sprint(item["ContactGroupId"]),
			"RegionId":       client.RegionId,
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] %s failed error: %v", action, err)
		}
	}
	return nil
}

func TestAccAlicloudArmsAlertContactGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_alert_contact_group.default"
	ra := resourceAttrInit(resourceId, ArmsAlertContactGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsAlertContactGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccArmsAlertContactGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ArmsAlertContactGroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_group_name": "${var.name}",
					"contact_ids":              []string{"${alicloud_arms_alert_contact.default.0.id}", "${alicloud_arms_alert_contact.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_group_name": name,
						"contact_ids.#":            "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_group_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_ids": []string{"${alicloud_arms_alert_contact.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_group_name": "${var.name}",
					"contact_ids":              []string{"${alicloud_arms_alert_contact.default.0.id}", "${alicloud_arms_alert_contact.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_group_name": name,
						"contact_ids.#":            "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var ArmsAlertContactGroupMap = map[string]string{}

func ArmsAlertContactGroupBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_arms_alert_contact" "default" {
	count = 2
    alert_contact_name = "${var.name}-${count.index}"
	email = "${var.name}-${count.index}@aaa.com"
}
`, name)
}

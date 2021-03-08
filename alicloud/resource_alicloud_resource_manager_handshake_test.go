package alicloud

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_resource_manager_handshake", &resource.Sweeper{
		Name: "alicloud_resource_manager_handshake",
		F:    testSweepResourceManagerHandshake,
	})
}

func testSweepResourceManagerHandshake(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "ListHandshakesForResourceDirectory"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}

	var handshakeIds []string

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExists.ResourceDirectory"}) {
				return nil
			}
			log.Printf("[ERROR] Failed to retrieve resoure manager handshake in service list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Handshakes.Handshake", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Handshakes.Handshake", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			// Skip Invalid handshake.
			if v, ok := item["Status"].(string); ok && v == "Pending" {
				handshakeIds = append(handshakeIds, item["HandshakeId"].(string))
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, handshakeId := range handshakeIds {
		log.Printf("[INFO] Delete resource manager handshake %s ", handshakeId)

		request := map[string]interface{}{
			"HandshakeId": handshakeId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete resource manager handshake (%s): %s", handshakeId, err)
		}
	}
	return nil
}

func TestAccAlicloudResourceManagerHandshake_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_handshake.default"
	ra := resourceAttrInit(resourceId, ResourceManagerHandshakeMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerHandshake")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerHandshake%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerHandshakeBasicdependence)
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
					"target_entity": "${alicloud_resource_manager_account.example.id}",
					"target_type":   "Account",
					"note":          "test resource manager handshake",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_entity": CHECKSET,
						"target_type":   "Account",
						"note":          "test resource manager handshake",
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

var ResourceManagerHandshakeMap = map[string]string{}

func ResourceManagerHandshakeBasicdependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_resource_manager_account" "example" {
  display_name = "%s"
}
`, name)
}

package alicloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
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
	resourceManagerService := ResourcemanagerService{client}

	request := resourcemanager.CreateListHandshakesForResourceDirectoryRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var handshakeIds []string
	for {
		raw, err := resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.ListHandshakesForResourceDirectory(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve resoure manager handshake in service list: %s", err)
		}

		response, _ := raw.(*resourcemanager.ListHandshakesForResourceDirectoryResponse)

		for _, v := range response.Handshakes.Handshake {
			// Skip Invalid handshake.
			if v.Status == "Pending" {
				handshakeIds = append(handshakeIds, v.HandshakeId)
			}
		}
		if len(response.Handshakes.Handshake) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	for _, handshakeId := range handshakeIds {
		log.Printf("[INFO] Delete resource manager handshake %s ", handshakeId)

		request := resourcemanager.CreateCancelHandshakeRequest()
		request.HandshakeId = handshakeId
		_, err = resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.CancelHandshake(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete resource manager handshake (%s): %s", handshakeId, err)
		}
	}
	return nil
}

func TestAccAlicloudResourceManagerHandshake_basic(t *testing.T) {
	var v resourcemanager.Handshake
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
					"target_entity": os.Getenv("ALICLOUD_ACCOUNT_ID"),
					"target_type":   "Account",
					"note":          "test resource manager handshake",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_entity": os.Getenv("ALICLOUD_ACCOUNT_ID"),
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
	return ""
}

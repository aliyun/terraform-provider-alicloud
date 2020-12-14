package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_group", &resource.Sweeper{
		Name: "alicloud_api_gateway_group",
		F:    testSweepApiGatewayGroup,
		Dependencies: []string{
			"alicloud_api_gateway_api",
		},
	})
}

func testSweepApiGatewayGroup(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	req := cloudapi.CreateDescribeApiGroupsRequest()
	req.PageNumber = requests.NewInteger(1)
	req.PageSize = requests.NewInteger(PageSizeLarge)
	sweeped := false

	for {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeApiGroups(req)
		})
		if err != nil {
			log.Printf("[ERROR] Describe Api Groups: %s", err)
			return nil
		}
		response, _ := raw.(*cloudapi.DescribeApiGroupsResponse)

		for _, v := range response.ApiGroupAttributes.ApiGroupAttribute {
			name := v.GroupName
			id := v.GroupId
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping api group: %s", name)
				continue
			}
			sweeped = true

			log.Printf("[INFO] Deleting Api Group: %s", name)

			req := cloudapi.CreateDeleteApiGroupRequest()
			req.GroupId = id
			_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
				return cloudApiClient.DeleteApiGroup(req)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Api Group (%s): %s", name, err)
			}
		}

		if len(response.ApiGroupAttributes.ApiGroupAttribute) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		req.PageNumber = page
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudApigatewayGroup_basic(t *testing.T) {
	var v *cloudapi.DescribeApiGroupResponse

	resourceId := "alicloud_api_gateway_group.default"
	ra := resourceAttrInit(resourceId, apigatewayGroupBasicMap)

	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccGroup_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayGroupConfigDependence)

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
					"name":        "${var.name}",
					"description": "${var.description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"name": "${var.name}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.description}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf_testAcc api gateway description_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${var.name}",
					"description": "${var.description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "tf_testAcc api gateway description",
					}),
				),
			},
		},
	})
}
func TestAccAlicloudApigatewayGroup_multi(t *testing.T) {
	var v *cloudapi.DescribeApiGroupResponse
	resourceId := "alicloud_api_gateway_group.default.9"
	ra := resourceAttrInit(resourceId, apigatewayGroupBasicMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccGroup_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayGroupConfigDependence)

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
					"name":        "${var.name}${count.index}",
					"description": "${var.description}",
					"count":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func resourceApigatewayGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	variable "description" {
	  default = "tf_testAcc api gateway description"
	}
	`, name)
}

var apigatewayGroupBasicMap = map[string]string{
	"name":        CHECKSET,
	"description": "tf_testAcc api gateway description",
}

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

func TestAccAliCloudApigatewayGroup_basic(t *testing.T) {
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
					"base_path":   "${var.base_path}",
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
					"base_path": "${var.base_path}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"base_path": "/test_by_tf_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${var.name}",
					"description": "${var.description}",
					"base_path":   "${var.base_path}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "tf_testAcc api gateway description",
						"base_path":   "/test_by_tf",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_log_config": []map[string]string{{
						"request_body":     "true",
						"response_body":    "true",
						"query_string":     "*",
						"request_headers":  "*",
						"response_headers": "*",
						"jwt_claims":       "*",
					},
					}}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_log_config.0.request_body":     "true",
						"user_log_config.0.response_body":    "true",
						"user_log_config.0.query_string":     "*",
						"user_log_config.0.request_headers":  "*",
						"user_log_config.0.response_headers": "*",
						"user_log_config.0.jwt_claims":       "*",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudApigatewayGroup_basic01(t *testing.T) {
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
					"instance_id": "api-shared-vpc-001",
					"base_path":   "${var.base_path}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "tf_testAcc api gateway description",
						"instance_id": "api-shared-vpc-001",
						"base_path":   "/test_by_tf",
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

func TestAccAliCloudApigatewayGroup_multi(t *testing.T) {
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
					"base_path":   "${var.base_path}",
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

	variable "base_path" {
      default = "/test_by_tf"
    }
	`, name)
}

var apigatewayGroupBasicMap = map[string]string{
	"name":        CHECKSET,
	"description": "tf_testAcc api gateway description",
}

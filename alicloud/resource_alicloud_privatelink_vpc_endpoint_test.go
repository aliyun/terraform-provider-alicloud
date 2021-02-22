package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_privatelink_vpc_endpoint",
		&resource.Sweeper{
			Name: "alicloud_privatelink_vpc_endpoint",
			F:    testSweepPrivatelinkVpcEndpoint,
		})
}

func testSweepPrivatelinkVpcEndpoint(region string) error {
	if !testSweepPreCheckWithRegions(region, false, connectivity.PrivateLinkRegions) {
		log.Printf("[INFO] Skipping privatelink unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	action := "ListVpcEndpoints"
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoint", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Endpoints", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Endpoints", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if item["EndpointName"] == nil {
				continue
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["EndpointName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Privatelink VpcEndpoint: %s", item["EndpointName"].(string))
				continue
			}
			sweeped = true
			action = "DeleteVpcEndpoint"
			request := map[string]interface{}{
				"EndpointId": item["EndpointId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Privatelink VpcEndpoint (%s): %s", item["EndpointName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Privatelink VpcEndpoint have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Privatelink VpcEndpoint success: %s ", item["EndpointName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudPrivatelinkVpcEndpoint_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivatelinkVpcEndpointBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_id":         "${alicloud_privatelink_vpc_endpoint_service.default.id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"security_group_ids": []string{"${alicloud_security_group.default.id}"},
					"vpc_endpoint_name":  name,
					"depends_on":         []string{"alicloud_privatelink_vpc_endpoint_service.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_id":           CHECKSET,
						"vpc_id":               CHECKSET,
						"security_group_ids.#": "1",
						"vpc_endpoint_name":    name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "Terraform Test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "Terraform Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.default2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{"${alicloud_security_group.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids":   []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.default2.id}"},
					"endpoint_description": "Terraform Test Update",
					"vpc_endpoint_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "2",
						"endpoint_description": "Terraform Test Update",
						"vpc_endpoint_name":    name,
					}),
				),
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointMap = map[string]string{
	"bandwidth":                CHECKSET,
	"connection_status":        CHECKSET,
	"endpoint_business_status": CHECKSET,
	"endpoint_domain":          CHECKSET,
	"service_name":             CHECKSET,
	"status":                   CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointBasicDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  is_default = true
	}
	resource "alicloud_security_group" "default" {
	  name        = "tftest"
	  description = "privatelink test security group"
	  vpc_id      = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_security_group" "default2" {
	  name        = "%[1]s"
	  description = "privatelink test security group2"
	  vpc_id      = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	 service_description = "%[1]s"
	 connect_bandwidth = 103
     auto_accept_connection = false
	}
`, name)
}

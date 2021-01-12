package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_privatelink_vpc_endpoint_service",
		&resource.Sweeper{
			Name: "alicloud_privatelink_vpc_endpoint_service",
			F:    testSweepPrivatelinkVpcEndpointService,
		})
}

func testSweepPrivatelinkVpcEndpointService(region string) error {
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
	action := "ListVpcEndpointServices"
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Services", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Services", response)
		}
		sweeped := false
		for _, v := range resp.([]interface{}) {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["ServiceDescription"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Privatelink VpcEndpoint Service: %s", item["ServiceId"].(string))
				continue
			}
			sweeped = true
			action = "DeleteVpcEndpointService"
			request := map[string]interface{}{
				"ServiceId": item["ServiceId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Privatelink VpcEndpoint Service (%s): %s", item["ServiceId"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Privatelink VpcEndpoint Service  have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Privatelink VpcEndpoint Service  success: %s ", item["ServiceId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudPrivatelinkVpcEndpointService_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpointService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointServiceTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivatelinkVpcEndpointServiceBasicDependence)
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
					"service_description":    name,
					"connect_bandwidth":      "103",
					"auto_accept_connection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description":    name,
						"connect_bandwidth":      "103",
						"auto_accept_connection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "payer"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_accept_connection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_accept_connection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connect_bandwidth": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connect_bandwidth": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_accept_connection": "false",
					"service_description":    name,
					"connect_bandwidth":      "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_accept_connection": "false",
						"service_description":    name,
						"connect_bandwidth":      "200",
					}),
				),
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointServiceMap = map[string]string{
	"service_business_status": "Normal",
	"service_domain":          CHECKSET,
	"status":                  CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointServiceBasicDependence(name string) string {
	return ""
}

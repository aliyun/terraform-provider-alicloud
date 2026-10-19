package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_privatelink_vpc_endpoint_zone",
		&resource.Sweeper{
			Name: "alicloud_privatelink_vpc_endpoint_zone",
			F:    testSweepPrivatelinkVpcEndpointZone,
		})
}

// testSweepPrivatelinkVpcEndpointZone removes every zone added to test-owned
// VPC endpoints so the parent endpoint sweeper can delete the endpoint. The
// remove mirrors resourceAliCloudPrivateLinkVpcEndpointZoneDelete
// (RemoveZoneFromVpcEndpoint).
func testSweepPrivatelinkVpcEndpointZone(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.PrivateLinkRegions) {
		log.Printf("[INFO] Skipping privatelink unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	endpoints, err := listAllPrivatelinkVpcEndpoints(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list privatelink vpc endpoints: %s", err)
		return nil
	}
	for _, ep := range endpoints {
		endpointId := fmt.Sprint(ep["EndpointId"])
		endpointName := fmt.Sprint(ep["EndpointName"])
		// Red-line #4: skip zones under infra-managed endpoints
		// (ProtectedEnabled=true) even under SWEEP_ALL — RemoveZoneFromVpcEndpoint
		// may succeed on a protected endpoint and would disrupt live
		// infra-managed services (e.g. ADB cluster endpoints).
		if privatelinkParentIsManagedProtected(ep) {
			log.Printf("[INFO] Skipping privatelink vpc endpoint zone under managed-protected endpoint: %s", endpointId)
			continue
		}
		if !privatelinkSweepParentShouldDelete(endpointName) {
			log.Printf("[INFO] Skipping privatelink vpc endpoint zone under protected endpoint: %s", endpointId)
			continue
		}
		listAction := "ListVpcEndpointZones"
		listReq := map[string]interface{}{
			"EndpointId": endpointId,
		}
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			response, err := client.RpcPost("Privatelink", "2020-04-15", listAction, nil, listReq, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
					break
				}
				log.Printf("[ERROR] Failed to list privatelink vpc endpoint zones for endpoint %s: %s", endpointId, err)
				break
			}
			resp, err := jsonpath.Get("$.Zones", response)
			if err != nil {
				log.Printf("[ERROR] Failed to parse $.Zones for endpoint %s: %s", endpointId, err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				zoneId := fmt.Sprint(item["ZoneId"])
				deleteAction := "RemoveZoneFromVpcEndpoint"
				deleteReq := map[string]interface{}{
					"EndpointId":  endpointId,
					"ZoneId":      zoneId,
					"RegionId":    client.RegionId,
					"ClientToken": buildClientToken(deleteAction),
				}
				if _, err := client.RpcPost("Privatelink", "2020-04-15", deleteAction, nil, deleteReq, true); err != nil {
					if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
						continue
					}
					log.Printf("[ERROR] Failed to remove privatelink vpc endpoint zone %s:%s: %s", endpointId, zoneId, err)
					continue
				}
				log.Printf("[INFO] Delete privatelink vpc endpoint zone success: %s:%s", endpointId, zoneId)
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				listReq["NextToken"] = nextToken
			} else {
				break
			}
		}
	}
	return nil
}

// Currently, Private network slb can only be created through the console.
func TestAccAliCloudPrivatelinkVpcEndpointZone_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_zone.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointZoneMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpointZone")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointZoneTest%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivatelinkVpcEndpointZoneBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SlbPrivateNetSupportRegions)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.default.id}",
					"vswitch_id":  "${alicloud_vswitch.default.id}",
					"zone_id":     "${alicloud_vswitch.default.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointZoneMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointZoneBasicDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_slb_zones" "default" {}

	resource "alicloud_vpc" "default" {
	  description = "test-terraform-service"
	  cidr_block  = "10.0.0.0/8"
	  vpc_name    = "%[1]s"
	}
	
	resource "alicloud_vswitch" "default" {
	  vpc_id     = alicloud_vpc.default.id
	  zone_id    = data.alicloud_slb_zones.default.zones.0.id
	  cidr_block = "10.1.0.0/16"
	}

	resource "alicloud_slb_load_balancer" "default" {
	  load_balancer_name = "%[1]s"
	  load_balancer_spec  = "slb.s2.small"
      address_type = "intranet"
      instance_charge_type = "PayBySpec"
      delete_protection = "off"
      vswitch_id = alicloud_vswitch.default.id
      master_zone_id = data.alicloud_slb_zones.default.zones.0.id
      slave_zone_id = data.alicloud_slb_zones.default.zones.1.id
      lifecycle {
        ignore_changes = [delete_protection]
      }
	}

	data "alicloud_vswitches" "default" {
	 	 is_default = true
	}
	resource "alicloud_security_group" "default" {
	 name = "%[1]s"
	 description = "privatelink test security group"
	 vpc_id = alicloud_vpc.default.id
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	  service_description   = "test for privatelink connection"
	  service_resource_type = "slb"
      auto_accept_connection = false
	}
	resource "alicloud_privatelink_vpc_endpoint_service_resource" "default" {
	 service_id    =  "${alicloud_privatelink_vpc_endpoint_service.default.id}"
	 resource_id   =  "${alicloud_slb_load_balancer.default.id}"
	 resource_type = "slb"
	}
	resource "alicloud_privatelink_vpc_endpoint" "default" {
	 service_id = alicloud_privatelink_vpc_endpoint_service_resource.default.service_id
	 vpc_id = alicloud_vpc.default.id
	 security_group_ids = [alicloud_security_group.default.id]
	 vpc_endpoint_name = "%[1]s"
	 depends_on = [alicloud_privatelink_vpc_endpoint_service.default]
	}
`, name)
}

// lintignore: R001
func TestUnitAlicloudPrivatelinkVpcEndpointZone(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_zone"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_zone"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"endpoint_id": "AddZoneToVpcEndpointValue",
		"vswitch_id":  "AddZoneToVpcEndpointValue",
		"zone_id":     "AddZoneToVpcEndpointValue",
		"dry_run":     false,
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// ListVpcEndpointZones
		"Zones": []interface{}{
			map[string]interface{}{
				"ZoneStatus": "Connected",
				"VSwitchId":  "AddZoneToVpcEndpointValue",
				"ZoneId":     "AddZoneToVpcEndpointValue",
				"EndpointId": "AddZoneToVpcEndpointValue",
			},
		},
		"VSwitchId": "AddZoneToVpcEndpointValue",
		"ZoneId":    "AddZoneToVpcEndpointValue",
	}
	CreateMockResponse := map[string]interface{}{
		// AddZoneToVpcEndpoint

	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_privatelink_vpc_endpoint_zone", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointZoneCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// ListVpcEndpointZones Response
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "EndpointConnectionOperationDenied", "EndpointLocked", "EndpointOperationDenied", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AddZoneToVpcEndpoint" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudPrivateLinkVpcEndpointZoneCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_zone"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	err = resourceAliCloudPrivateLinkVpcEndpointZoneUpdate(dExisted, rawClient)
	assert.NotNil(t, err)

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_privatelink_vpc_endpoint_zone", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_zone"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListVpcEndpointZones" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudPrivateLinkVpcEndpointZoneRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointZoneDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_privatelink_vpc_endpoint_zone", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_zone"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "EndpointConnectionOperationDenied", "EndpointLocked", "EndpointOperationDenied", "nil", "EndpointZoneNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RemoveZoneFromVpcEndpoint" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudPrivateLinkVpcEndpointZoneDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "EndpointZoneNotFound":
			assert.Nil(t, err)
		}
	}
}

func TestAccAliCloudPrivateLinkVpcEndpointZone_basic_without_zoneid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_zone.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointZoneMap4898)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpointZone")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointzone%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointZoneBasicDependence4898)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.defaulti9F95i.id}",
					"vswitch_id":  "${alicloud_vswitch.defaultmEwUAc.id}",
					"eni_ip":      "10.1.0.245",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":     CHECKSET,
						"endpoint_id": CHECKSET,
						"vswitch_id":  CHECKSET,
						"eni_ip":      "10.1.0.245",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAliCloudPrivateLinkVpcEndpointZone_basic4898(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_zone.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointZoneMap4898)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpointZone")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointzone%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointZoneBasicDependence4898)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":     "${data.alicloud_nlb_zones.default.zones.0.id}",
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.defaulti9F95i.id}",
					"vswitch_id":  "${alicloud_vswitch.defaultmEwUAc.id}",
					"eni_ip":      "10.1.0.245",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":     CHECKSET,
						"endpoint_id": CHECKSET,
						"vswitch_id":  CHECKSET,
						"eni_ip":      "10.1.0.245",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudPrivateLinkVpcEndpointZoneMap4898 = map[string]string{
	"status": CHECKSET,
}

func AlicloudPrivateLinkVpcEndpointZoneBasicDependence4898(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_nlb_zones" "default" {}

resource "alicloud_vpc" "defaultbFzA4a" {
  description = "test-terraform"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "defaultmEwUAc" {
  vpc_id     = alicloud_vpc.defaultbFzA4a.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "10.1.0.0/16"
}

resource "alicloud_security_group" "default1FTFrP" {
  name = var.name
  vpc_id = alicloud_vpc.defaultbFzA4a.id
}

resource "alicloud_privatelink_vpc_endpoint_service" "defaultr0WBYX" {
  service_description   = "test-zejun-service"
  connect_bandwidth     = "3072"
  service_resource_type = "nlb"
}

resource "alicloud_privatelink_vpc_endpoint" "defaulti9F95i" {
  vpc_id     = alicloud_vpc.defaultbFzA4a.id
  service_id = alicloud_privatelink_vpc_endpoint_service.defaultr0WBYX.id
  security_group_ids = [alicloud_security_group.default1FTFrP.id]
  # The NLB-backed endpoint service only reports a zone as supported after the
  # service_resource (the NLB registration) is created. Without this dependency
  # the endpoint and service_resource race in parallel, and AddZoneToVpcEndpoint
  # can run before the zone is propagated, failing with
  # EndpointServiceNotSupportedZone. The service_resource Create waits on its own
  # state refresh, so depending on it here lets endpoint_zone inherit the ordering.
  depends_on = [alicloud_privatelink_vpc_endpoint_service_resource.defaultdTPOne]
}

resource "alicloud_vpc" "defaultVpcService" {
  description = "test-terraform-service"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "defaultuYH1VC" {
  vpc_id     = alicloud_vpc.defaultVpcService.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "10.1.0.0/16"
}

resource "alicloud_vswitch" "defaultTJZ8ud" {
  vpc_id     = alicloud_vpc.defaultVpcService.id
  zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  cidr_block = "10.10.0.0/16"
}

resource "alicloud_nlb_load_balancer" "defaultyuY5jZ" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultuYH1VC.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultTJZ8ud.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
  load_balancer_name = var.name

  vpc_id       = alicloud_vpc.defaultVpcService.id
  address_type = "Intranet"
}

resource "alicloud_privatelink_vpc_endpoint_service_resource" "defaultdTPOne" {
  zone_id       = data.alicloud_nlb_zones.default.zones.0.id
  resource_id   = alicloud_nlb_load_balancer.defaultyuY5jZ.id
  resource_type = "nlb"
  service_id    = alicloud_privatelink_vpc_endpoint_service.defaultr0WBYX.id
}


`, name)
}

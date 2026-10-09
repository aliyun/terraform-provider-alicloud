package alicloud

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

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

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
	p := Provider().(*schema.Provider).ResourcesMap
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
	// gomonkey cannot intercept the tea-rpc client.NewClient constructor
	// reliably (its 4-line wrapper gets inlined into the call site), so patch
	// the DoRequest method instead: a structured non-retryable SDK error must
	// propagate out of Create.
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
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
	assert.Nil(t, err)

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

	// The Read ladder's last "{}" case hits the NotFound branch in Read and
	// clears dExisted.Id(); rebuild dExisted from dInit.State() so Delete gets a
	// valid "<endpoint_id>:<zone_id>" id instead of panicking on
	// strings.Split("", ":").
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_privatelink_vpc_endpoint_zone", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_zone"].Schema).Data(dInit.State(), diff)

	// Delete
	// Same reason as the Create block: patch the DoRequest method (the
	// client.NewClient constructor seam is inlined away by tea-rpc v1.3.4); a
	// structured non-deny, non-404 SDK error must not be swallowed by Delete.
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
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
	errorCodes = []string{"NonRetryableError", "Throttling", "EndpointConnectionOperationDenied", "EndpointLocked", "EndpointOperationDenied", "nil", "EndpointZoneNotFound", "GatewayLoadBalancerZoneCountDeny"}
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
		case "EndpointZoneNotFound", "GatewayLoadBalancerZoneCountDeny":
			// GatewayLoadBalancerZoneCountDeny is expected to be treated as
			// deleted: a GWLB-backed endpoint service requires the endpoint to
			// keep its minimum zone count, so the backend refuses the removal
			// and the zone only leaves the endpoint together with the endpoint.
			assert.Nil(t, err)
		}
	}
}

// TestUnitAlicloudPrivatelinkVpcEndpointZoneDeleteGwlbZoneCountDeny locks the
// delete behavior against a local RPC stub: a GWLB-backed endpoint service
// requires the endpoint to keep its minimum zone count, so
// RemoveZoneFromVpcEndpoint answers GatewayLoadBalancerZoneCountDeny and the
// zone only leaves the endpoint together with the endpoint itself
// (DeleteVpcEndpoint). Delete must accept that error instead of failing the
// destroy.
func TestUnitAlicloudPrivatelinkVpcEndpointZoneDeleteGwlbZoneCountDeny(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.Form.Get("Action"); got != "RemoveZoneFromVpcEndpoint" {
			t.Errorf("Action = %q, want RemoveZoneFromVpcEndpoint", got)
		}
		if got := r.Form.Get("EndpointId"); got != "ep-test" {
			t.Errorf("EndpointId = %q, want ep-test", got)
		}
		if got := r.Form.Get("ZoneId"); got != "cn-hangzhou-a" {
			t.Errorf("ZoneId = %q, want cn-hangzhou-a", got)
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"RequestId":"mock","HostId":"privatelink.aliyuncs.com","Code":"GatewayLoadBalancerZoneCountDeny","Message":"The GatewayLoadBalancer zone count not support"}`)
	}))
	defer server.Close()

	host := server.Listener.Addr().String()
	t.Setenv("NO_PROXY", host)

	var endpoints sync.Map
	endpoints.Store("privatelink", host)
	var signVersions sync.Map
	cred, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	assert.Nil(t, err)
	rawClient, err := (&connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           cred,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            &endpoints,
		SignVersion:          &signVersions,
		SkipRegionValidation: true,
	}).Client()
	assert.Nil(t, err)

	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_zone"].Schema).Data(nil, nil)
	d.SetId("ep-test:cn-hangzhou-a")

	err = resourceAliCloudPrivateLinkVpcEndpointZoneDelete(d, rawClient)
	assert.Nil(t, err)
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

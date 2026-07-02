package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	testAccPrivateLinkCrossRegionServiceRegion  = "cn-hangzhou"
	testAccPrivateLinkCrossRegionEndpointRegion = "cn-beijing"
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
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Services", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Services", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
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
			_, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
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

func TestAccAliCloudPrivatelinkVpcEndpointService_base(t *testing.T) {
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
				ImportStateVerifyIgnore: []string{"dry_run"},
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

func TestAccAliCloudPrivatelinkVpcEndpointService_supportedRegionList(t *testing.T) {
	resourceId := "alicloud_privatelink_vpc_endpoint_service.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPrivateLinkCrossRegion%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPrivateLinkCrossRegion(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccPrivateLinkProviderFactories(),
		CheckDestroy:      testAccCheckPrivateLinkVpcEndpointServiceDestroyInRegion(testAccPrivateLinkCrossRegionServiceRegion),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateLinkVpcEndpointServiceSupportedRegionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description":     name,
						"service_resource_type":   "nlb",
						"auto_accept_connection":  "true",
						"supported_region_list.#": "1",
					}),
					testAccCheckPrivateLinkVpcEndpointServiceSupportedRegions(resourceId, testAccPrivateLinkCrossRegionServiceRegion),
					testAccAttachPrivateLinkVpcEndpointServiceResources(
						resourceId,
						"alicloud_nlb_load_balancer.service",
						"data.alicloud_nlb_zones.default",
					),
				),
			},
			{
				Config: testAccPrivateLinkVpcEndpointServiceSupportedAllRegionsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"supported_region_list.#": "2",
						"service_support_ipv6":    "false",
					}),
					testAccCheckPrivateLinkVpcEndpointServiceSupportedRegions(resourceId, testAccPrivateLinkCrossRegionServiceRegion, testAccPrivateLinkCrossRegionEndpointRegion),
				),
			},
			{
				Config: testAccPrivateLinkVpcEndpointServiceSupportedRegionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"supported_region_list.#": "1",
					}),
					testAccCheckPrivateLinkVpcEndpointServiceSupportedRegions(resourceId, testAccPrivateLinkCrossRegionServiceRegion),
				),
			},
			// ImportState is intentionally omitted because this test uses
			// ProviderFactories with an inline provider region; Terraform SDK v1
			// import steps can fail with "unknown provider \"alicloud\"" in this shape.
		},
	})
}

func testAccPrivateLinkVpcEndpointServiceSupportedRegionConfig(name string) string {
	return testAccPrivateLinkVpcEndpointServiceSupportedRegionListConfig(name, []string{
		testAccPrivateLinkCrossRegionServiceRegion,
	})
}

func testAccPrivateLinkVpcEndpointServiceSupportedAllRegionsConfig(name string) string {
	return testAccPrivateLinkVpcEndpointServiceSupportedRegionListConfig(name, []string{
		testAccPrivateLinkCrossRegionServiceRegion,
		testAccPrivateLinkCrossRegionEndpointRegion,
	})
}

func testAccPrivateLinkProviderFactories() map[string]terraform.ResourceProviderFactory {
	return map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			return Provider(), nil
		},
	}
}

func testAccPreCheckPrivateLinkCrossRegion(t *testing.T) {
	testAccPreCheck(t)
	if !testAccPrivateLinkRegionInList(connectivity.Hangzhou, connectivity.NLBSupportRegions) {
		t.Skipf("Skipping unsupported service region %s. NLB supported regions: %v.", connectivity.Hangzhou, connectivity.NLBSupportRegions)
		t.Skipped()
	}
}

func testAccPrivateLinkRegionInList(region connectivity.Region, regions []connectivity.Region) bool {
	for _, r := range regions {
		if region == r {
			return true
		}
	}
	return false
}

func testAccCheckPrivateLinkVpcEndpointServiceDestroyInRegion(region string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rawClient, err := sharedClientForRegion(region)
		if err != nil {
			return WrapErrorf(err, "Error getting Alicloud client.")
		}
		client := rawClient.(*connectivity.AliyunClient)
		service := PrivateLinkServiceV2{client}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_privatelink_vpc_endpoint_service" {
				continue
			}
			_, err := service.DescribePrivateLinkVpcEndpointService(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}
			return fmt.Errorf("PrivateLink VpcEndpointService %s still exists", rs.Primary.ID)
		}
		return nil
	}
}

func testAccAttachPrivateLinkVpcEndpointServiceResources(serviceResourceName, resourceResourceName, zonesDataSourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceId, err := testAccPrivateLinkStateResourceID(s, serviceResourceName)
		if err != nil {
			return err
		}
		resourceId, err := testAccPrivateLinkStateResourceID(s, resourceResourceName)
		if err != nil {
			return err
		}
		zoneIds, err := testAccPrivateLinkNlbZoneIdsFromState(s, zonesDataSourceName)
		if err != nil {
			return err
		}

		rawClient, err := sharedClientForRegion(testAccPrivateLinkCrossRegionServiceRegion)
		if err != nil {
			return WrapErrorf(err, "Error getting Alicloud client.")
		}
		client := rawClient.(*connectivity.AliyunClient)
		for _, zoneId := range zoneIds {
			if err := testAccAttachPrivateLinkVpcEndpointServiceResource(client, serviceId, resourceId, zoneId); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccPrivateLinkStateResourceID(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("resource %s not found", resourceName)
	}
	if rs.Primary == nil || rs.Primary.ID == "" {
		return "", fmt.Errorf("resource %s has empty ID", resourceName)
	}
	return rs.Primary.ID, nil
}

func testAccPrivateLinkNlbZoneIdsFromState(s *terraform.State, resourceName string) ([]string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource %s not found", resourceName)
	}
	if rs.Primary == nil {
		return nil, fmt.Errorf("resource %s has empty state", resourceName)
	}

	zoneIds := make([]string, 0, 2)
	for i := 0; i < 2; i++ {
		zoneId := rs.Primary.Attributes[fmt.Sprintf("zones.%d.id", i)]
		if zoneId == "" {
			zoneId = rs.Primary.Attributes[fmt.Sprintf("ids.%d", i)]
		}
		if zoneId == "" {
			return nil, fmt.Errorf("resource %s has empty NLB zone at index %d", resourceName, i)
		}
		zoneIds = append(zoneIds, zoneId)
	}
	return zoneIds, nil
}

func testAccAttachPrivateLinkVpcEndpointServiceResource(client *connectivity.AliyunClient, serviceId, resourceId, zoneId string) error {
	privateLinkServiceV2 := PrivateLinkServiceV2{client}
	id := fmt.Sprintf("%s:%s:%s", serviceId, resourceId, zoneId)
	if _, err := privateLinkServiceV2.DescribePrivateLinkVpcEndpointServiceResource(id); err == nil {
		return nil
	} else if !NotFoundError(err) {
		return WrapError(err)
	}

	action := "AttachResourceToVpcEndpointService"
	query := make(map[string]interface{})
	request := map[string]interface{}{
		"ResourceId":   resourceId,
		"ServiceId":    serviceId,
		"ZoneId":       zoneId,
		"RegionId":     client.RegionId,
		"ResourceType": "nlb",
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request["ClientToken"] = buildClientToken(action)
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointServiceOperationDenied", "ConcurrentCallNotSupported", "EndpointServiceLocked"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if _, describeErr := privateLinkServiceV2.DescribePrivateLinkVpcEndpointServiceResource(id); describeErr == nil {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service_resource", action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{resourceId}, 5*time.Minute, 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceResourceStateRefreshFunc(id, "ResourceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, id)
	}
	return nil
}

func testAccCheckPrivateLinkVpcEndpointServiceSupportedRegions(resourceName string, expected ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}
		actual := make(map[string]bool)
		for key, value := range rs.Primary.Attributes {
			if strings.HasPrefix(key, "supported_region_list.") && key != "supported_region_list.#" {
				actual[value] = true
			}
		}
		if len(actual) != len(expected) {
			return fmt.Errorf("expected supported_region_list %#v, got %#v", expected, actual)
		}
		for _, region := range expected {
			if !actual[region] {
				return fmt.Errorf("expected supported_region_list to contain %s, got %#v", region, actual)
			}
		}
		return nil
	}
}

var AlicloudPrivatelinkVpcEndpointServiceMap = map[string]string{
	"service_business_status": "Normal",
	"service_domain":          CHECKSET,
	"status":                  CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointServiceBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alicloud_privatelink_service" "open" {
	  enable = "On"
	}
`, name)
}

func testAccPrivateLinkVpcEndpointServiceSupportedRegionListConfig(name string, supportedRegions []string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  region = "%s"
}

data "alicloud_nlb_zones" "default" {}

resource "alicloud_vpc" "service" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "service_a" {
  vpc_id     = alicloud_vpc.service.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "10.1.0.0/16"
}

resource "alicloud_vswitch" "service_b" {
  vpc_id     = alicloud_vpc.service.id
  zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  cidr_block = "10.2.0.0/16"
}

resource "alicloud_nlb_load_balancer" "service" {
  load_balancer_name = var.name
  vpc_id             = alicloud_vpc.service.id
  address_type       = "Intranet"

  zone_mappings {
    vswitch_id = alicloud_vswitch.service_a.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }

  zone_mappings {
    vswitch_id = alicloud_vswitch.service_b.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}

resource "alicloud_privatelink_vpc_endpoint_service" "default" {
  service_description    = var.name
  service_resource_type  = "nlb"
  auto_accept_connection = true
  supported_region_list  = [
%s
  ]

  depends_on = [alicloud_nlb_load_balancer.service]
}
`, name, testAccPrivateLinkCrossRegionServiceRegion, testAccPrivateLinkQuotedList(supportedRegions))
}

func testAccPrivateLinkQuotedList(values []string) string {
	list := make([]string, 0, len(values))
	for _, value := range values {
		list = append(list, fmt.Sprintf("    %q,", value))
	}
	return strings.Join(list, "\n")
}

func TestUnitVpcEndpointServiceConnectBandwidthReadByResourceType(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap

	cases := []struct {
		name                string
		serviceResourceType string
		response            map[string]interface{}
		expected            int
	}{
		{
			name:                "slb missing connect bandwidth is authoritative",
			serviceResourceType: "slb",
			response: map[string]interface{}{
				"ServiceResourceType": "slb",
			},
			expected: 0,
		},
		{
			name:                "slb returned connect bandwidth is authoritative",
			serviceResourceType: "slb",
			response: map[string]interface{}{
				"ServiceResourceType": "slb",
				"ConnectBandwidth":    100,
			},
			expected: 100,
		},
		{
			name:                "nlb missing connect bandwidth preserves config",
			serviceResourceType: "nlb",
			response: map[string]interface{}{
				"ServiceResourceType": "nlb",
			},
			expected: 3072,
		},
		{
			name:                "alb zero connect bandwidth preserves config",
			serviceResourceType: "alb",
			response: map[string]interface{}{
				"ServiceResourceType": "alb",
				"ConnectBandwidth":    0,
			},
			expected: 3072,
		},
		{
			name:                "gwlb empty connect bandwidth preserves config",
			serviceResourceType: "gwlb",
			response: map[string]interface{}{
				"ServiceResourceType": "gwlb",
				"ConnectBandwidth":    "",
			},
			expected: 3072,
		},
	}

	for _, tc := range cases {
		d, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(nil, nil)
		assert.Nil(t, d.Set("service_resource_type", tc.serviceResourceType))
		assert.Nil(t, d.Set("connect_bandwidth", 3072))
		setVpcEndpointServiceConnectBandwidth(d, tc.response)
		assert.Equal(t, tc.expected, d.Get("connect_bandwidth"), tc.name)
	}
}

// lintignore: R001
func TestUnitAlicloudPrivatelinkVpcEndpointService(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"service_description":    "CreateVpcEndpointServiceValue",
		"connect_bandwidth":      100,
		"auto_accept_connection": false,
		"dry_run":                false,
		"payer":                  "CreateVpcEndpointServiceValue",
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
		// GetVpcEndpointServiceAttribute
		"AutoAcceptEnabled":     false,
		"ConnectBandwidth":      100,
		"ServiceBusinessStatus": "CreateVpcEndpointServiceValue",
		"ServiceDescription":    "CreateVpcEndpointServiceValue",
		"ServiceDomain":         "CreateVpcEndpointServiceValue",
		"ServiceStatus":         "Active",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateVpcEndpoint
		"ServiceId": "CreateVpcEndpointServiceValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_privatelink_vpc_endpoint_service", errorCode))
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
	err = resourceAliCloudPrivateLinkVpcEndpointServiceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetVpcEndpointServiceAttribute Response
		"ServiceId": "CreateVpcEndpointServiceValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateVpcEndpointService" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointServiceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointServiceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateVpcEndpointServiceAttribute
	attributesDiff := map[string]interface{}{
		"auto_accept_connection": true,
		"connect_bandwidth":      200,
		"service_description":    "UpdateVpcEndpointServiceAttributeValue",
		"dry_run":                true,
	}
	diff, err := newInstanceDiff("alicloud_privatelink_vpc_endpoint_service", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetVpcEndpointServiceAttribute Response
		"AutoAcceptEnabled":  true,
		"ConnectBandwidth":   200,
		"ServiceDescription": "UpdateVpcEndpointServiceAttributeValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateVpcEndpointServiceAttribute" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudPrivateLinkVpcEndpointServiceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetVpcEndpointServiceAttribute" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointServiceRead(dExisted, rawClient)
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
	err = resourceAliCloudPrivateLinkVpcEndpointServiceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "EndpointServiceConnectionDependence", "nil", "EndpointServiceNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVpcEndpointService" {
				switch errorCode {
				case "NonRetryableError", "EndpointServiceNotFound":
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
		err := resourceAliCloudPrivateLinkVpcEndpointServiceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "EndpointServiceNotFound":
			assert.Nil(t, err)
		}
	}

}

// Test PrivateLink VpcEndpointService. >>> Resource test cases, automatically generated.
// Case 生命周期测试-克隆-nlb 4837
func TestAccAliCloudPrivatelinkVpcEndpointService_case4837(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointServiceMap4837)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpointService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointServiceBasicDependence4837)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description":    "test-zejun",
					"auto_accept_connection": "false",
					"payer":                  "Endpoint",
					"service_resource_type":  "nlb",
					"zone_affinity_enabled":  "false",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"dry_run":                "false",
					"service_support_ipv6":   "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description":    "test-zejun",
						"auto_accept_connection": "false",
						"payer":                  "Endpoint",
						"service_resource_type":  "nlb",
						"zone_affinity_enabled":  "false",
						"resource_group_id":      CHECKSET,
						"dry_run":                "false",
						"service_support_ipv6":   "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description":   "test-zejun-2",
					"zone_affinity_enabled": "true",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					//"service_support_ipv6":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description":   "test-zejun-2",
						"zone_affinity_enabled": "true",
						"resource_group_id":     CHECKSET,
						//"service_support_ipv6":  "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description":    "test-zejun",
					"auto_accept_connection": "true",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description":    "test-zejun",
						"auto_accept_connection": "true",
						"resource_group_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connect_bandwidth", "dry_run"},
			},
		},
	})
}

var AlicloudPrivateLinkVpcEndpointServiceMap4837 = map[string]string{
	"vpc_endpoint_service_name": CHECKSET,
	"status":                    CHECKSET,
	"create_time":               CHECKSET,
	"service_domain":            CHECKSET,
	"service_business_status":   CHECKSET,
	"region_id":                 CHECKSET,
}

func AlicloudPrivateLinkVpcEndpointServiceBasicDependence4837(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case pvl+gwlb生命周期测试 9628
func TestAccAliCloudPrivatelinkVpcEndpointService_case9628(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointServiceMap9628)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpointService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointServiceBasicDependence9628)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.WuLanChaBu})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payer":                  "Endpoint",
					"auto_accept_connection": "false",
					"service_description":    "pvl+gwlb测试create",
					"dry_run":                "false",
					"service_resource_type":  "gwlb",
					"address_ip_version":     "IPv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payer":                  "Endpoint",
						"auto_accept_connection": "false",
						"service_description":    "pvl+gwlb测试create",
						"dry_run":                "false",
						"service_resource_type":  "gwlb",
						"address_ip_version":     "IPv4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": "测试update",
					"address_ip_version":  "DualStack",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "测试update",
						"address_ip_version":  "DualStack",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudPrivateLinkVpcEndpointServiceMap9628 = map[string]string{
	"vpc_endpoint_service_name": CHECKSET,
	"status":                    CHECKSET,
	"create_time":               CHECKSET,
	"service_business_status":   CHECKSET,
	"region_id":                 CHECKSET,
}

func AlicloudPrivateLinkVpcEndpointServiceBasicDependence9628(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}


`, name)
}

// Test PrivateLink VpcEndpointService. <<< Resource test cases, automatically generated.

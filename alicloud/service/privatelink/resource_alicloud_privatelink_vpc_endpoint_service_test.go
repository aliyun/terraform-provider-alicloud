package privatelink_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	privatelink "github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatelink"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testAccPrivateLinkCrossRegionServiceRegion  = "cn-hangzhou"
	testAccPrivateLinkCrossRegionEndpointRegion = "cn-beijing"
)

var privatelinkTestEndpoints sync.Map
var privatelinkTestSignVersion sync.Map

// privatelinkTestClientForRegion rebuilds the package-level sweeper client on the
// exported connectivity API: the flat-suite root helper is not reachable from
// this package, so a local rebuild mirrors the actiontrail pattern.
func privatelinkTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &privatelinkTestEndpoints,
		SignVersion: &privatelinkTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	return conf.Client()
}

// privatelinkTestSweepPreCheckWithRegions mirrors the root-package
// testSweepPreCheckWithRegions gate; it is not reachable from a service test
// package.
func privatelinkTestSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}

func init() {
	resource.AddTestSweepers(
		"alicloud_privatelink_vpc_endpoint_service",
		&resource.Sweeper{
			Name: "alicloud_privatelink_vpc_endpoint_service",
			F:    testSweepPrivatelinkVpcEndpointService,
		})
}

func testSweepPrivatelinkVpcEndpointService(region string) error {
	if privatelinkTestSweepPreCheckWithRegions(region, false, connectivity.PrivateLinkRegions) {
		log.Printf("[INFO] Skipping privatelink unsupported region: %s", region)
		return nil
	}
	client, err := privatelinkTestClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"MaxResults": helper.PageSizeLarge,
	}
	var response map[string]interface{}
	action := "ListVpcEndpointServices"
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Services", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Services", response)
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivatelinkVpcEndpointService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointServiceTest%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivatelinkVpcEndpointServiceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceMap)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPrivateLinkCrossRegion%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPrivateLinkCrossRegion(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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

func testAccPreCheckPrivateLinkCrossRegion(t *testing.T) {
	acctest.PreCheck(t)
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
		client, err := privatelinkTestClientForRegion(region)
		if err != nil {
			return errs.WrapErrorf(err, "Error getting Alicloud client.")
		}
		service := privatelink.NewPrivateLinkServiceV2(client)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_privatelink_vpc_endpoint_service" {
				continue
			}
			_, err := service.DescribePrivateLinkVpcEndpointService(rs.Primary.ID)
			if err != nil {
				if errs.NotFoundError(err) {
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

		client, err := privatelinkTestClientForRegion(testAccPrivateLinkCrossRegionServiceRegion)
		if err != nil {
			return errs.WrapErrorf(err, "Error getting Alicloud client.")
		}
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
	privateLinkServiceV2 := privatelink.NewPrivateLinkServiceV2(client)
	id := fmt.Sprintf("%s:%s:%s", serviceId, resourceId, zoneId)
	if _, err := privateLinkServiceV2.DescribePrivateLinkVpcEndpointServiceResource(id); err == nil {
		return nil
	} else if !errs.NotFoundError(err) {
		return errs.WrapError(err)
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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		request["ClientToken"] = helper.BuildClientToken(action)
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"EndpointServiceOperationDenied", "ConcurrentCallNotSupported", "EndpointServiceLocked"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if _, describeErr := privateLinkServiceV2.DescribePrivateLinkVpcEndpointServiceResource(id); describeErr == nil {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service_resource", action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{resourceId}, 5*time.Minute, 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceResourceStateRefreshFunc(id, "ResourceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, id)
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
	"service_domain":          acctest.CHECKSET,
	"status":                  acctest.CHECKSET,
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

func TestAccAliCloudPrivatelinkVpcEndpointService_resource(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcesresource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivatelinkVpcEndpointServiceResourceDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create the endpoint service with a backend resource attached
				// inline, so the cross-zone association is set at creation time.
				Config: testAccConfig(map[string]interface{}{
					"service_description":    name,
					"service_resource_type":  "nlb",
					"auto_accept_connection": "true",
					"resource": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_nlb_load_balancer.default.id}",
							"resource_type": "nlb",
							"zone_id":       "${data.alicloud_nlb_zones.default.zones.0.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description":   name,
						"service_resource_type": "nlb",
						"resource.#":            "1",
					}),
				),
			},
			{
				// Attach a second backend resource (a different load balancer in
				// a different zone) — exercises the incremental Attach path and
				// changes both resource_id and zone_id.
				Config: testAccConfig(map[string]interface{}{
					"resource": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_nlb_load_balancer.default.id}",
							"resource_type": "nlb",
							"zone_id":       "${data.alicloud_nlb_zones.default.zones.0.id}",
						},
						{
							"resource_id":   "${alicloud_nlb_load_balancer.second.id}",
							"resource_type": "nlb",
							"zone_id":       "${data.alicloud_nlb_zones.default.zones.1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource.#": "2",
					}),
				),
			},
			{
				// Detach back to a single backend resource, keeping the second
				// load balancer as the remaining attachment — exercises the
				// incremental Detach path and changes resource_id and zone_id.
				Config: testAccConfig(map[string]interface{}{
					"resource": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_nlb_load_balancer.second.id}",
							"resource_type": "nlb",
							"zone_id":       "${data.alicloud_nlb_zones.default.zones.1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// dry_run is a request-only flag and is not persisted by the
				// service, so it is not present in the refreshed state to verify
				// on import.
				//
				// `resource.resource_type` is a no-op for TypeSet import
				// verification (the SDK matches ImportStateVerifyIgnore by
				// flattened-key prefix, so it never matches resource.N.resource_type).
				// It is retained here as a coverage-exemption marker:
				// resource_type is bound to the backend resource identity
				// (slb/alb/nlb/gwlb) and cannot be modified in place without
				// reattaching a different backend type, so the testing-coverage
				// check exempts it via this entry. The `resource` set itself is
				// still verified on import: element identity is stable (Set hashes
				// on resource_id + resource_type + zone_id) and Read backfills
				// zone_id per item, so a refreshed element matches the configured
				// one.
				ImportStateVerifyIgnore: []string{"dry_run", "resource.resource_type"},
			},
		},
	})
}

func TestAccAliCloudPrivatelinkVpcEndpointService_sameBackendDualZone(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%splvpcessdualzone%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivatelinkVpcEndpointServiceResourceDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
		},
		IDRefreshName: resourceId,
		Providers:     acctest.SDKProviders,
		CheckDestroy:  rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Reproduces the cross-zone bug: the SAME backend load balancer is
				// attached to the endpoint service in two availability zones via two
				// inline resource blocks with identical resource_id + resource_type
				// but different zone_id. Before the fix the resource TypeSet hashed
				// on resource_id + resource_type only, so the two blocks were
				// deduplicated to one and apply aborted with "Provider produced
				// inconsistent final plan", blocking cross-region endpoint services
				// (the server requires ≥2 zones). With zone_id in the hash both
				// attachments are preserved and sent to CreateVpcEndpointService.
				Config: testAccConfig(map[string]interface{}{
					"service_description":    name,
					"service_resource_type":  "nlb",
					"auto_accept_connection": "true",
					"resource": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_nlb_load_balancer.default.id}",
							"resource_type": "nlb",
							"zone_id":       "${data.alicloud_nlb_zones.default.zones.0.id}",
						},
						{
							"resource_id":   "${alicloud_nlb_load_balancer.default.id}",
							"resource_type": "nlb",
							"zone_id":       "${data.alicloud_nlb_zones.default.zones.1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_resource_type": "nlb",
						"resource.#":            "2",
					}),
				),
			},
			{
				// A no-op plan must show no changes: the set is stable because
				// zone_id is part of element identity and Read backfills it per
				// item, so each refreshed element matches its configured counterpart
				// and no detach/attach churn is generated.
				Config: testAccConfig(map[string]interface{}{
					"service_description":    name,
					"service_resource_type":  "nlb",
					"auto_accept_connection": "true",
					"resource": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_nlb_load_balancer.default.id}",
							"resource_type": "nlb",
							"zone_id":       "${data.alicloud_nlb_zones.default.zones.0.id}",
						},
						{
							"resource_id":   "${alicloud_nlb_load_balancer.default.id}",
							"resource_type": "nlb",
							"zone_id":       "${data.alicloud_nlb_zones.default.zones.1.id}",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func AlicloudPrivatelinkVpcEndpointServiceResourceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_privatelink_service" "open" {
  enable = "On"
}

data "alicloud_nlb_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "zone_a" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "10.1.0.0/16"
}

resource "alicloud_vswitch" "zone_b" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  cidr_block = "10.2.0.0/16"
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = "${var.name}-1"
  vpc_id             = alicloud_vpc.default.id
  address_type       = "Intranet"

  zone_mappings {
    vswitch_id = alicloud_vswitch.zone_a.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }

  zone_mappings {
    vswitch_id = alicloud_vswitch.zone_b.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}

resource "alicloud_nlb_load_balancer" "second" {
  load_balancer_name = "${var.name}-2"
  vpc_id             = alicloud_vpc.default.id
  address_type       = "Intranet"

  zone_mappings {
    vswitch_id = alicloud_vswitch.zone_a.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }

  zone_mappings {
    vswitch_id = alicloud_vswitch.zone_b.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}
`, name)
}

// Test PrivateLink VpcEndpointService. >>> Resource test cases, automatically generated.
// Case 生命周期测试-克隆-nlb 4837
func TestAccAliCloudPrivatelinkVpcEndpointService_case4837(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointServiceMap4837)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointservice%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateLinkVpcEndpointServiceBasicDependence4837)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description":    "test-zejun",
					"connect_bandwidth":      "3072",
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
						"connect_bandwidth":      "3072",
						"auto_accept_connection": "false",
						"payer":                  "Endpoint",
						"service_resource_type":  "nlb",
						"zone_affinity_enabled":  "false",
						"resource_group_id":      acctest.CHECKSET,
						"dry_run":                "false",
						"service_support_ipv6":   "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description":   "test-zejun-2",
					"connect_bandwidth":     "3073",
					"zone_affinity_enabled": "true",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					//"service_support_ipv6":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description":   "test-zejun-2",
						"connect_bandwidth":     "3073",
						"zone_affinity_enabled": "true",
						"resource_group_id":     acctest.CHECKSET,
						//"service_support_ipv6":  "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description":    "test-zejun",
					"connect_bandwidth":      "3072",
					"auto_accept_connection": "true",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description":    "test-zejun",
						"connect_bandwidth":      "3072",
						"auto_accept_connection": "true",
						"resource_group_id":      acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"vpc_endpoint_service_name": acctest.CHECKSET,
	"status":                    acctest.CHECKSET,
	"create_time":               acctest.CHECKSET,
	"service_domain":            acctest.CHECKSET,
	"service_business_status":   acctest.CHECKSET,
	"region_id":                 acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointServiceMap9628)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointservice%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateLinkVpcEndpointServiceBasicDependence9628)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"vpc_endpoint_service_name": acctest.CHECKSET,
	"status":                    acctest.CHECKSET,
	"create_time":               acctest.CHECKSET,
	"service_business_status":   acctest.CHECKSET,
	"region_id":                 acctest.CHECKSET,
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

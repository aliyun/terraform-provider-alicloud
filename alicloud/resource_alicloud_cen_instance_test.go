package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cen_instance", &resource.Sweeper{
		Name: "alicloud_cen_instance",
		F:    testSweepCenInstances,
		Dependencies: []string{
			"alicloud_cen_bandwidth_package",
			"alicloud_cen_flowlog",
			"alicloud_cen_instance_attachment",
			"alicloud_cen_bandwidth_limit",
		},
	})
}

var sweepCenInstanceIds []string

func testSweepCenInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []cbn.Cen
	describeCensRequest := cbn.CreateDescribeCensRequest()
	describeCensRequest.RegionId = client.RegionId
	describeCensRequest.PageSize = requests.NewInteger(PageSizeLarge)
	describeCensRequest.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCbnClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DescribeCens(describeCensRequest)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving CEN Instances: %s", err)
		}
		describeCensResponse, _ := raw.(*cbn.DescribeCensResponse)
		if len(describeCensResponse.Cens.Cen) < 1 {
			break
		}
		insts = append(insts, describeCensResponse.Cens.Cen...)

		if len(describeCensResponse.Cens.Cen) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(describeCensRequest.PageNumber)
		if err != nil {
			return err
		}
		describeCensRequest.PageNumber = page
	}

	sweepCenInstanceIds = make([]string, 0)
	for _, cenInstance := range insts {
		name := cenInstance.Name
		cenId := cenInstance.CenId
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping CEN Instance: %s (%s)", name, cenId)
				continue
			}
		}
		sweepCenInstanceIds = append(sweepCenInstanceIds, cenId)
		describeCenAttachedChildInstancesRequest := cbn.CreateDescribeCenAttachedChildInstancesRequest()
		describeCenAttachedChildInstancesRequest.CenId = cenId
		raw, err := client.WithCbnClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DescribeCenAttachedChildInstances(describeCenAttachedChildInstancesRequest)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to Describe CEN Attached Instance (%s (%s)): %s", name, cenId, err)
		}
		describeCenAttachedChildInstancesResponse, _ := raw.(*cbn.DescribeCenAttachedChildInstancesResponse)
		for _, childInstance := range describeCenAttachedChildInstancesResponse.ChildInstances.ChildInstance {
			instanceId := childInstance.ChildInstanceId
			log.Printf("[INFO] Detaching CEN Child Instance: %s (%s %s)", name, cenId, instanceId)
			detachCenChildInstanceRequest := cbn.CreateDetachCenChildInstanceRequest()
			detachCenChildInstanceRequest.CenId = cenId
			detachCenChildInstanceRequest.ChildInstanceId = instanceId
			detachCenChildInstanceRequest.ChildInstanceType = childInstance.ChildInstanceType
			detachCenChildInstanceRequest.ChildInstanceRegionId = childInstance.ChildInstanceRegionId
			_, err := client.WithCbnClient(func(cenClient *cbn.Client) (interface{}, error) {
				return cenClient.DetachCenChildInstance(detachCenChildInstanceRequest)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to Detach CEN Attached Instance (%s (%s %s)): %s", name, cenId, instanceId, err)
			}
			cenService := CenService{client}
			err = cenService.WaitForCenInstanceAttachment(cenId+COLON_SEPARATED+instanceId, Deleted, DefaultCenTimeoutLong)
			if err != nil {
				log.Printf("[ERROR] Failed to WaitFor CEN Attached Instance Detached (%s (%s %s)): %s", name, cenId, instanceId, err)
			}
		}

		action := "ListTransitRouterVbrAttachments"
		request := make(map[string]interface{})
		request["CenId"] = cenId
		request["RegionId"] = client.RegionId
		request["PageSize"] = PageSizeLarge
		request["PageNumber"] = 1
		var response map[string]interface{}
		for {
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(2*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] %s failed: %v", action, err)
				break
			}
			resp, err := jsonpath.Get("$.TransitRouterAttachments", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterAttachments", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				name := fmt.Sprint(item["TransitRouterAttachmentName"])
				id := fmt.Sprint(item["TransitRouterAttachmentId"])
				action := "DeleteTransitRouterVbrAttachment"
				log.Printf("[DEBUG] %s %s:%s", action, id, name)

				request := map[string]interface{}{
					"TransitRouterAttachmentId": id,
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"Operation.Blocking", "InvalidOperation.ChildInstanceStatus"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					log.Printf("[ERROR] %s failed %v", action, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}

		action = "ListTransitRouterVpcAttachments"
		request = make(map[string]interface{})
		request["CenId"] = cenId
		request["RegionId"] = client.RegionId
		request["PageSize"] = PageSizeLarge
		request["PageNumber"] = 1
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(2*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] %s failed: %v", action, err)
				break
			}
			resp, err := jsonpath.Get("$.TransitRouterAttachments", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterAttachments", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				name := fmt.Sprint(item["TransitRouterAttachmentName"])
				id := fmt.Sprint(item["TransitRouterAttachmentId"])
				action := "DeleteTransitRouterVpcAttachment"
				log.Printf("[DEBUG] %s %s:%s", action, id, name)

				request := map[string]interface{}{
					"TransitRouterAttachmentId": id,
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					log.Printf("[ERROR] %s failed %v", action, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}

		action = "ListTransitRouterPeerAttachments"
		request = make(map[string]interface{})
		request["CenId"] = cenId
		request["RegionId"] = client.RegionId
		request["PageSize"] = PageSizeLarge
		request["PageNumber"] = 1
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(2*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] %s failed: %v", action, err)
				break
			}
			resp, err := jsonpath.Get("$.TransitRouterAttachments", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterAttachments", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				name := fmt.Sprint(item["TransitRouterAttachmentName"])
				id := fmt.Sprint(item["TransitRouterAttachmentId"])
				action := "DeleteTransitRouterPeerAttachment"
				log.Printf("[DEBUG] %s %s:%s", action, id, name)

				request := map[string]interface{}{
					"TransitRouterAttachmentId": id,
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					log.Printf("[ERROR] %s failed %v", action, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}

		describeCenPrivateZoneRoutesRequest := cbn.CreateDescribeCenPrivateZoneRoutesRequest()
		describeCenPrivateZoneRoutesRequest.RegionId = client.RegionId
		describeCenPrivateZoneRoutesRequest.AccessRegionId = client.RegionId
		describeCenPrivateZoneRoutesRequest.CenId = cenInstance.CenId

		raw, err = client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenPrivateZoneRoutes(describeCenPrivateZoneRoutesRequest)
		})
		if err == nil {
			response, _ := raw.(*cbn.DescribeCenPrivateZoneRoutesResponse)
			for _, resp := range response.PrivateZoneInfos.PrivateZoneInfo {
				request := cbn.CreateUnroutePrivateZoneInCenToVpcRequest()
				request.AccessRegionId = resp.AccessRegionId
				request.CenId = cenInstance.CenId
				if _, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
					return cbnClient.UnroutePrivateZoneInCenToVpc(request)
				}); err != nil {
					log.Printf("\n Failed to UnroutePrivateZoneInCenToVpc. Error: %v", err)
				}
			}
		}

		action = "ListTransitRouters"
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["CenId"] = cenId
		request["PageSize"] = PageSizeLarge
		request["PageNumber"] = 1
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] %s failed: %v", action, err)
				break
			}
			resp, err := jsonpath.Get("$.TransitRouters", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouters", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				id := fmt.Sprint(item["TransitRouterId"])
				action := "ListTransitRouterRouteTables"
				request := make(map[string]interface{})
				request["RegionId"] = client.RegionId
				request["TransitRouterId"] = id
				request["PageSize"] = PageSizeLarge
				request["PageNumber"] = 1
				var response map[string]interface{}
				for {
					response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
					if err != nil {
						log.Printf("[ERROR] %s failed: %v", action, err)
						break
					}
					resp, err := jsonpath.Get("$.TransitRouterRouteTables", response)
					if err != nil {
						log.Printf("\n jsonpath.Get $.TransitRouterRouteTables failed %v", err)
						break
					}
					result, _ := resp.([]interface{})
					for _, v := range result {
						item := v.(map[string]interface{})
						id := fmt.Sprint(item["TransitRouterRouteTableId"])
						action := "DeleteTransitRouterRouteTable"
						log.Printf("[DEBUG] %s %s", action, name)
						request := map[string]interface{}{
							"TransitRouterRouteTableId": id,
						}
						response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
						if err != nil {
							log.Printf("[ERROR] %s failed %v", action, err)
						}
					}
					if len(result) < PageSizeLarge {
						break
					}
					request["PageNumber"] = request["PageNumber"].(int) + 1
				}

				action = "DeleteTransitRouter"
				log.Printf("[DEBUG] %s %s", action, id)

				request = map[string]interface{}{
					"TransitRouterId": id,
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"Operation.Blocking"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					log.Printf("[ERROR] %s failed %v", action, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	for _, cenId := range sweepCenInstanceIds {

		log.Printf("[INFO] Deleting CEN Instance: %s ", cenId)
		deleteCenRequest := cbn.CreateDeleteCenRequest()
		deleteCenRequest.CenId = cenId
		_, err = client.WithCbnClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DeleteCen(deleteCenRequest)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CEN Instance (%s): %s", cenId, err)
		}
	}
	return nil
}

func TestAccAliCloudCenInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CenSupportRegions)
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudCenInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCenConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level": "REDUCED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level": "REDUCED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Instance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Instance",
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

func TestAccAliCloudCenInstance_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CenSupportRegions)
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudCenInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCenConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenInstanceBasicDependence0)
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
					"protection_level":  "REDUCED",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"cen_instance_name": name,
					"description":       name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Instance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":  "REDUCED",
						"resource_group_id": CHECKSET,
						"cen_instance_name": name,
						"description":       name,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Instance",
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

func TestAccAliCloudCenInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CenSupportRegions)
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudCenInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCenConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level": "REDUCED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level": "REDUCED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Instance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Instance",
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

func TestAccAliCloudCenInstance_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CenSupportRegions)
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudCenInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCenConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenInstanceBasicDependence0)
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
					"protection_level":  "REDUCED",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"name":              name,
					"description":       name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Instance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":  "REDUCED",
						"resource_group_id": CHECKSET,
						"name":              name,
						"description":       name,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Instance",
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

func TestAccAliCloudCenInstance_Multi(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CenSupportRegions)
	resourceId := "alicloud_cen_instance.default.2"
	ra := resourceAttrInit(resourceId, AliCloudCenInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCenConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenInstanceBasicDependence0)
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
					"count":             "3",
					"protection_level":  "REDUCED",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"cen_instance_name": name + "-${count.index}",
					"description":       name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Instance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":  "REDUCED",
						"resource_group_id": CHECKSET,
						"cen_instance_name": name + fmt.Sprint(-2),
						"description":       name,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Instance",
					}),
				),
			},
		},
	})
}

var AliCloudCenInstanceMap0 = map[string]string{
	"protection_level":  CHECKSET,
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
}

func AliCloudCenInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name)
}

func TestUnitAliCloudCenInstance(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cen_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cen_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"cen_instance_name": "CreateCenValue",
		"description":       "CreateCenValue",
		"protection_level":  "CreateCenValue",
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
		// DescribeCens
		"Cens": map[string]interface{}{
			"Cen": []interface{}{
				map[string]interface{}{
					"CenBandwidthPackageIds": map[string]interface{}{
						"CenBandwidthPackageId": []interface{}{},
					},
					"CenId":           "CreateCenValue",
					"Name":            "CreateCenValue",
					"CreationTime":    "DefaultValue",
					"Description":     "CreateCenValue",
					"ProtectionLevel": "CreateCenValue",
					"Status":          "Active",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateCen
		"CenId":     "CreateCenValue",
		"RequestId": "MockValue",
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cen_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudCenCenInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeCens Response
		"Cens": map[string]interface{}{
			"Cen": []interface{}{
				map[string]interface{}{
					"CenId": "CreateCenValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "Operation.Blocking", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateCen" {
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
		err := resourceAliCloudCenCenInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cen_instance"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudCenCenInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyCenAttribute
	attributesDiff := map[string]interface{}{
		"cen_instance_name": "ModifyCenAttributeValue",
		"description":       "ModifyCenAttributeValue",
		"protection_level":  "ModifyCenAttributeValue",
	}
	diff, err := newInstanceDiff("alicloud_cen_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cen_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeCens Response
		"Cens": map[string]interface{}{
			"Cen": []interface{}{
				map[string]interface{}{
					"Name":            "ModifyCenAttributeValue",
					"Description":     "ModifyCenAttributeValue",
					"ProtectionLevel": "ModifyCenAttributeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyCenAttribute" {
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
		err := resourceAliCloudCenCenInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cen_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// TagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"TagResourcesValue_1": "TagResourcesValue_1",
			"TagResourcesValue_2": "TagResourcesValue_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_cen_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cen_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeCens Response
		"Cens": map[string]interface{}{
			"Cen": []interface{}{
				map[string]interface{}{
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "TagResourcesValue_1",
								"Value": "TagResourcesValue_1",
							},
							map[string]interface{}{
								"Key":   "TagResourcesValue_2",
								"Value": "TagResourcesValue_2",
							},
						},
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "TagResources" {
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
		err := resourceAliCloudCenCenInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cen_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UntagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"UntagResourcesValue3_1": "UnTagResourcesValue3_1",
			"UntagResourcesValue3_2": "UnTagResourcesValue3_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_cen_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cen_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeCens Response
		"Cens": map[string]interface{}{
			"Cen": []interface{}{
				map[string]interface{}{
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "UntagResourcesValue3_1",
								"Value": "UnTagResourcesValue3_1",
							},
							map[string]interface{}{
								"Key":   "UntagResourcesValue3_2",
								"Value": "UnTagResourcesValue3_2",
							},
						},
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UntagResources" {
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
		err := resourceAliCloudCenCenInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cen_instance"].Schema).Data(dExisted.State(), nil)
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
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeCens" {
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
		err := resourceAliCloudCenCenInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudCenCenInstanceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "ParameterCenInstanceId"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteCen" {
				switch errorCode {
				case "NonRetryableError", "ParameterCenInstanceId":
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
		err := resourceAliCloudCenCenInstanceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "ParameterCenInstanceId":
			assert.Nil(t, err)
		}
	}
}

// Test Cen CenInstance. >>> Resource test cases, automatically generated.
// Case CenInstance_20241108_线上 8803
func TestAccAliCloudCenCenInstance_basic8803(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCenCenInstanceMap8803)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenCenInstanceBasicDependence8803)
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
					"cen_instance_name": name,
					"description":       "create",
					"protection_level":  "FULL",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name,
						"description":       "create",
						"protection_level":  "FULL",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": name + "_update",
					"description":       "update",
					"protection_level":  "REDUCED",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name + "_update",
						"description":       "update",
						"protection_level":  "REDUCED",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenCenInstanceMap8803 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudCenCenInstanceBasicDependence8803(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 依赖资源_副本1730799777882_副本1730870121932 8693
func TestAccAliCloudCenCenInstance_basic8693(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCenCenInstanceMap8693)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenCenInstanceBasicDependence8693)
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
					"cen_instance_name": name,
					"description":       "create",
					"protection_level":  "FULL",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name,
						"description":       "create",
						"protection_level":  "FULL",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": name + "_update",
					"description":       "update",
					"protection_level":  "REDUCED",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name + "_update",
						"description":       "update",
						"protection_level":  "REDUCED",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenCenInstanceMap8693 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudCenCenInstanceBasicDependence8693(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 依赖资源_副本1730799777882 8674
func TestAccAliCloudCenCenInstance_basic8674(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCenCenInstanceMap8674)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenCenInstanceBasicDependence8674)
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
					"cen_instance_name": name,
					"description":       "create",
					"protection_level":  "FULL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name,
						"description":       "create",
						"protection_level":  "FULL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": name + "_update",
					"description":       "update",
					"protection_level":  "REDUCED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name + "_update",
						"description":       "update",
						"protection_level":  "REDUCED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenCenInstanceMap8674 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudCenCenInstanceBasicDependence8674(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 依赖资源 7163
func TestAccAliCloudCenCenInstance_basic7163(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCenCenInstanceMap7163)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenCenInstanceBasicDependence7163)
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
					"cen_instance_name": name,
					"description":       "create",
					"protection_level":  "FULL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name,
						"description":       "create",
						"protection_level":  "FULL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": name + "_update",
					"description":       "update",
					"protection_level":  "REDUCED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name + "_update",
						"description":       "update",
						"protection_level":  "REDUCED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenCenInstanceMap7163 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudCenCenInstanceBasicDependence7163(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 全生命周期_可重入 4401
func TestAccAliCloudCenCenInstance_basic4401(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCenCenInstanceMap4401)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenCenInstanceBasicDependence4401)
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
					"cen_instance_name": name,
					"description":       "create",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"protection_level":  "FULL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name,
						"description":       "create",
						"resource_group_id": CHECKSET,
						"protection_level":  "FULL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": name + "_update",
					"description":       "update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"protection_level":  "REDUCED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name + "_update",
						"description":       "update",
						"resource_group_id": CHECKSET,
						"protection_level":  "REDUCED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenCenInstanceMap4401 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudCenCenInstanceBasicDependence4401(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 组播域接入RMC依赖CEN 3514
func TestAccAliCloudCenCenInstance_basic3514(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCenCenInstanceMap3514)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenCenInstanceBasicDependence3514)
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
					"cen_instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name,
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenCenInstanceMap3514 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudCenCenInstanceBasicDependence3514(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 全生命周期v2 3124
func TestAccAliCloudCenCenInstance_basic3124(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCenCenInstanceMap3124)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenCenInstanceBasicDependence3124)
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
					"protection_level":  "FULL",
					"description":       "test",
					"cen_instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":  "FULL",
						"description":       "test",
						"cen_instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":  "REDUCED",
					"description":       "testupdate",
					"cen_instance_name": name + "_update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":  "REDUCED",
						"description":       "testupdate",
						"cen_instance_name": name + "_update",
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenCenInstanceMap3124 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudCenCenInstanceBasicDependence3124(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 全生命周期v2_无资源组依赖，可重入 3141
func TestAccAliCloudCenCenInstance_basic3141(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCenCenInstanceMap3141)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenCenInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenCenInstanceBasicDependence3141)
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
					"protection_level":  "FULL",
					"description":       "test",
					"cen_instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":  "FULL",
						"description":       "test",
						"cen_instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":  "REDUCED",
					"description":       "testupdate",
					"cen_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":  "REDUCED",
						"description":       "testupdate",
						"cen_instance_name": name + "_update",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenCenInstanceMap3141 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudCenCenInstanceBasicDependence3141(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Cen CenInstance. <<< Resource test cases, automatically generated.

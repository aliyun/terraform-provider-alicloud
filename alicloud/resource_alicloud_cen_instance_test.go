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

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

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
		return fmt.Errorf("error getting Alicloud client: %s", err)
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
		raw, err := client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
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
		sweepCenInstanceIds = append(sweepCenInstanceIds, cenId)
		describeCenAttachedChildInstancesRequest := cbn.CreateDescribeCenAttachedChildInstancesRequest()
		describeCenAttachedChildInstancesRequest.CenId = cenId
		raw, err := client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
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
			_, err := client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
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
		conn, err := client.NewCbnClient()
		if err != nil {
			return WrapError(err)
		}
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(2*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
				conn, err := client.NewCbnClient()
				if err != nil {
					return WrapError(err)
				}
				for {
					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		_, err = client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DeleteCen(deleteCenRequest)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CEN Instance (%s): %s", cenId, err)
		}
	}
	return nil
}

func TestAccAlicloudCenInstance_basic(t *testing.T) {
	var cen cbn.Cen
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, cenInstanceMap)
	serviceFunc := func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &cen, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCenConfig-%d", defaultRegionToTest, rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CenNoSkipRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": name,
					"description":       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name,
						"description":       name,
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
					"tags": map[string]string{
						"Created": "TF",
						"Name":    name,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.Name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"cen_instance_name": name + "update"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": name + "update"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": name,
					"description":       name,
					"tags": map[string]string{
						"Created": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": name,
						"description":       name,
						"tags.%":            "1",
						"tags.Created":      "TF",
						"tags.Name":         REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCenInstance_basic1(t *testing.T) {
	var cen cbn.Cen
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, cenInstanceMap)
	serviceFunc := func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &cen, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCenConfig-%d", defaultRegionToTest, rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CenNoSkipRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"description":      name,
					"protection_level": "REDUCED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"description":      name,
						"protection_level": "REDUCED",
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

func TestAccAlicloudCenInstance_multi(t *testing.T) {
	var cen cbn.Cen
	resourceId := "alicloud_cen_instance.default.4"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &cen, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CenNoSkipRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceMultiConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": fmt.Sprintf("tf-testAcc%sCenConfig-%d-4", defaultRegionToTest, rand),
						"description":       "tf-testAccCenConfigDescription",
					}),
				),
			},
		},
	})
}

var cenInstanceMap = map[string]string{
	"protection_level": "REDUCED",
	"status":           "Active",
	"description":      "tf-testAccCenConfigDescription",
}

func resourceCenInstanceConfigDependence(name string) string {
	return ""
}

func testAccCenInstanceMultiConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cen_instance" "default" {
		cen_instance_name = "tf-testAcc%sCenConfig-%d-${count.index}"
		description = "tf-testAccCenConfigDescription"
		count = 5
}
`, defaultRegionToTest, rand)
}

func testAccCheckCenInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_instance" {
			continue
		}

		// Try to find the CEN
		cbnService := CbnService{client}
		instance, err := cbnService.DescribeCenInstance(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if fmt.Sprint(instance["CenId"]) != "" {
			return fmt.Errorf("CEN %s still exist", fmt.Sprint(instance["CenId"]))
		}
	}

	return nil
}

func TestUnitAlicloudCenInstance(t *testing.T) {
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
	err = resourceAlicloudCenInstanceCreate(dInit, rawClient)
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
		err := resourceAlicloudCenInstanceCreate(dInit, rawClient)
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
	err = resourceAlicloudCenInstanceUpdate(dExisted, rawClient)
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
		err := resourceAlicloudCenInstanceUpdate(dExisted, rawClient)
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
		err := resourceAlicloudCenInstanceUpdate(dExisted, rawClient)
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
		err := resourceAlicloudCenInstanceUpdate(dExisted, rawClient)
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
		err := resourceAlicloudCenInstanceRead(dExisted, rawClient)
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
	err = resourceAlicloudCenInstanceDelete(dExisted, rawClient)
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
		err := resourceAlicloudCenInstanceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "ParameterCenInstanceId":
			assert.Nil(t, err)
		}
	}
}

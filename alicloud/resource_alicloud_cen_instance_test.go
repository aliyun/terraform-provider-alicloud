package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
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

	// Regional clients are cached so each candidate region's cbn endpoint is
	// created only once across all CENs. VPC/VBR/VPN attachments and transit
	// routers are regional to the child instance's region, not the CEN's home
	// region, so sub-resource cleanup must probe each candidate region.
	regionalClients := map[string]*connectivity.AliyunClient{}

	sweepCenInstanceIds = make([]string, 0)
	if err := forEachCenForSubSweep(client, func(cenId, cenName string) {
		sweepCenInstanceIds = append(sweepCenInstanceIds, cenId)

		// (a) Detach basic child instances (VPC/VBR attached via the basic CEN
		// child-instance model, not the enterprise TR attachment model). Clears
		// InvalidOperation.AttachedChildInstanceExist on TRs that still have
		// basic child-instance attachments.
		detachCenChildInstances(client, cenId, cenName)

		// (b) Delete ALL transit-router attachments across all candidate
		// regions FIRST. VPC/VBR/VPN attachments live in the child instance's
		// region (queried via regional clients); peer attachments live in the
		// CEN home region. Force=true removes route-table associations
		// automatically. This clears AttachedChildInstanceExist on the TR.
		sweepCenTransitRouterAttachments(client, cenId, cenName,
			"ListTransitRouterVbrAttachments", "DeleteTransitRouterVbrAttachment", "VBR",
			cenSubSweepRegions(client.RegionId), regionalClients)
		sweepCenTransitRouterAttachments(client, cenId, cenName,
			"ListTransitRouterVpcAttachments", "DeleteTransitRouterVpcAttachment", "VPC",
			cenSubSweepRegions(client.RegionId), regionalClients)
		sweepCenTransitRouterAttachments(client, cenId, cenName,
			"ListTransitRouterVpnAttachments", "DeleteTransitRouterVpnAttachment", "VPN",
			cenSubSweepRegions(client.RegionId), regionalClients)
		sweepCenTransitRouterAttachments(client, cenId, cenName,
			"ListTransitRouterPeerAttachments", "DeleteTransitRouterPeerAttachment", "Peer",
			[]string{client.RegionId}, regionalClients)

		// Fallback: delete orphaned TR Peer attachments by attachment ID.
		// ListTransitRouterPeerAttachments returns IllegalParam when the
		// owning TR was already deleted (by a prior sweep run), leaving
		// orphaned Peer attachments with DataTransfer bandwidth allocations
		// that block DeleteCen with InvalidOperation.AttachedChildInstanceExist.
		sweepCenOrphanedPeerAttachments(client, cenId, cenName)

		// Unroute private zones (non-fatal; legacy flow).
		unrouteCenPrivateZones(client, cenId)

		// (c) Delete NON-system transit-router route tables, then
		// (d) delete transit routers. System route tables cannot be deleted
		// (OperationFailed.DeleteSystemRouteTable) and are skipped; route
		// tables with residual associations/propagations/route entries are
		// logged and left for the CEN delete to surface (there is no Force
		// param on DeleteTransitRouterRouteTable).
		sweepCenTransitRouterRouteTablesAndTRs(client, cenId, cenName, regionalClients)
	}); err != nil {
		return err
	}

	// (e) DeleteCen. The API has no Force param; if custom route tables or
	// route entries remain (OperationFailed.DeleteCenWithCustomRouteTable /
	// DeleteCenWithRouteEntry), the CEN is logged as a leftover with the
	// reason rather than force-deleted.
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

// detachCenChildInstances detaches all basic child instances (VPC/VBR) from a
// CEN via the DescribeCenAttachedChildInstances / DetachCenChildInstance APIs.
// This is the basic child-instance model (not the enterprise TR attachment
// model) and must run before TR/attachment cleanup so
// InvalidOperation.AttachedChildInstanceExist is cleared on the TR. Non-fatal:
// per-instance errors are logged and the loop continues.
func detachCenChildInstances(client *connectivity.AliyunClient, cenId, cenName string) {
	describeCenAttachedChildInstancesRequest := cbn.CreateDescribeCenAttachedChildInstancesRequest()
	describeCenAttachedChildInstancesRequest.CenId = cenId
	raw, err := client.WithCbnClient(func(cenClient *cbn.Client) (interface{}, error) {
		return cenClient.DescribeCenAttachedChildInstances(describeCenAttachedChildInstancesRequest)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to Describe CEN Attached Instance (%s (%s)): %s", cenName, cenId, err)
		return
	}
	describeCenAttachedChildInstancesResponse, _ := raw.(*cbn.DescribeCenAttachedChildInstancesResponse)
	for _, childInstance := range describeCenAttachedChildInstancesResponse.ChildInstances.ChildInstance {
		instanceId := childInstance.ChildInstanceId
		log.Printf("[INFO] Detaching CEN Child Instance: %s (%s %s)", cenName, cenId, instanceId)
		detachCenChildInstanceRequest := cbn.CreateDetachCenChildInstanceRequest()
		detachCenChildInstanceRequest.CenId = cenId
		detachCenChildInstanceRequest.ChildInstanceId = instanceId
		detachCenChildInstanceRequest.ChildInstanceType = childInstance.ChildInstanceType
		detachCenChildInstanceRequest.ChildInstanceRegionId = childInstance.ChildInstanceRegionId
		_, err := client.WithCbnClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DetachCenChildInstance(detachCenChildInstanceRequest)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to Detach CEN Attached Instance (%s (%s %s)): %s", cenName, cenId, instanceId, err)
			continue
		}
		cenService := CenService{client}
		if err := cenService.WaitForCenInstanceAttachment(cenId+COLON_SEPARATED+instanceId, Deleted, DefaultCenTimeoutLong); err != nil {
			log.Printf("[ERROR] Failed to WaitFor CEN Attached Instance Detached (%s (%s %s)): %s", cenName, cenId, instanceId, err)
		}
	}
}

// unrouteCenPrivateZones unroutes private-zone routes from a CEN (legacy
// flow). Non-fatal: errors are logged and the loop continues.
func unrouteCenPrivateZones(client *connectivity.AliyunClient, cenId string) {
	describeCenPrivateZoneRoutesRequest := cbn.CreateDescribeCenPrivateZoneRoutesRequest()
	describeCenPrivateZoneRoutesRequest.RegionId = client.RegionId
	describeCenPrivateZoneRoutesRequest.AccessRegionId = client.RegionId
	describeCenPrivateZoneRoutesRequest.CenId = cenId
	raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeCenPrivateZoneRoutes(describeCenPrivateZoneRoutesRequest)
	})
	if err != nil {
		return
	}
	response, _ := raw.(*cbn.DescribeCenPrivateZoneRoutesResponse)
	for _, resp := range response.PrivateZoneInfos.PrivateZoneInfo {
		request := cbn.CreateUnroutePrivateZoneInCenToVpcRequest()
		request.AccessRegionId = resp.AccessRegionId
		request.CenId = cenId
		if _, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.UnroutePrivateZoneInCenToVpc(request)
		}); err != nil {
			log.Printf("[ERROR] Failed to UnroutePrivateZoneInCenToVpc for CEN %s: %v", cenId, err)
		}
	}
}

// sweepCenTransitRouterRouteTablesAndTRs lists transit routers for a CEN across
// all candidate regions, deletes non-system route tables under each TR, then
// deletes the TRs. Transit routers are regional resources: a CEN with a
// cn-beijing home region can have TRs in cn-shanghai, so the LIST must probe
// each candidate region via a regional client. regionalClients caches
// per-region AliyunClients so each endpoint is created only once.
//
// Force-order: route tables are deleted BEFORE the owning TR (otherwise
// DeleteTransitRouter fails with InvalidOperation.AttachedChildInstanceExist /
// residual route-table references). System route tables are skipped
// (OperationFailed.DeleteSystemRouteTable). Route tables with residual
// associations/propagations/route entries are logged and left for the CEN
// delete to surface — DeleteTransitRouterRouteTable has no Force param.
func sweepCenTransitRouterRouteTablesAndTRs(client *connectivity.AliyunClient, cenId, cenName string, regionalClients map[string]*connectivity.AliyunClient) {
	seenTR := map[string]bool{}
	for _, region := range cenSubSweepRegions(client.RegionId) {
		listClient := client
		if region != client.RegionId {
			if rc, ok := regionalClients[region]; ok {
				listClient = rc
			} else {
				raw, err := sharedClientForRegion(region)
				if err != nil {
					log.Printf("[ERROR] sharedClientForRegion(%s) for CEN %s (%s): %v", region, cenName, cenId, err)
					continue
				}
				listClient = raw.(*connectivity.AliyunClient)
				regionalClients[region] = listClient
			}
		}
		// List transit routers in this region for this CEN. PageSize/PageNumber
		// pagination (max 50 per page).
		trReq := map[string]interface{}{
			"CenId":      cenId,
			"RegionId":   region,
			"PageSize":   PageSizeLarge,
			"PageNumber": 1,
		}
		for {
			var trResp map[string]interface{}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				resp, e := listClient.RpcPost("Cbn", "2017-09-12", "ListTransitRouters", nil, trReq, true)
				trResp = resp
				if e != nil {
					if IsExpectedErrors(e, []string{"IllegalParam.CenIdOrRegionId", "IllegalParam.Region"}) {
						return retry.NonRetryableError(e)
					}
					if NeedRetry(e) {
						wait()
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				addDebug("ListTransitRouters", trResp, trReq)
				return nil
			})
			if err != nil {
				// Non-fatal: many regions have no TR for this CEN.
				break
			}
			att, jerr := jsonpath.Get("$.TransitRouters", trResp)
			if jerr != nil {
				break
			}
			result, _ := att.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				trId := fmt.Sprint(item["TransitRouterId"])
				if trId == "" || seenTR[trId] {
					continue
				}
				seenTR[trId] = true
				// (c) Delete non-system route tables, then (d) delete the TR.
				sweepCenTransitRouterRouteTables(listClient, cenId, cenName, trId)
				deleteTransitRouter(listClient, cenId, cenName, trId)
			}
			if len(result) < PageSizeLarge {
				break
			}
			trReq["PageNumber"] = trReq["PageNumber"].(int) + 1
		}
	}
}

// sweepCenTransitRouterRouteTables lists route tables for a transit router and
// deletes only Custom route tables — System route tables are not deletable
// (OperationFailed.DeleteSystemRouteTable). Uses MaxResults/NextToken
// pagination (the ListTransitRouterRouteTables API contract). Route tables with
// residual associations/propagations/route entries are non-retryable and
// logged (DeleteTransitRouterRouteTable has no Force param).
func sweepCenTransitRouterRouteTables(listClient *connectivity.AliyunClient, cenId, cenName, trId string) {
	rtReq := map[string]interface{}{
		"TransitRouterId": trId,
		"MaxResults":      100,
	}
	for {
		var resp map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(2*time.Minute, func() *retry.RetryError {
			r, e := listClient.RpcPost("Cbn", "2017-09-12", "ListTransitRouterRouteTables", nil, rtReq, true)
			resp = r
			if e != nil {
				if IsExpectedErrors(e, []string{"InvalidTransitRouterId.NotFound"}) {
					return retry.NonRetryableError(e)
				}
				if NeedRetry(e) {
					wait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			addDebug("ListTransitRouterRouteTables", resp, rtReq)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] ListTransitRouterRouteTables failed for TR %s (CEN %s (%s)): %v", trId, cenName, cenId, err)
			break
		}
		att, jerr := jsonpath.Get("$.TransitRouterRouteTables", resp)
		if jerr != nil {
			log.Printf("[ERROR] jsonpath.Get $.TransitRouterRouteTables failed for TR %s (CEN %s (%s)): %v", trId, cenName, cenId, jerr)
			break
		}
		result, _ := att.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			rtId := fmt.Sprint(item["TransitRouterRouteTableId"])
			rtType := fmt.Sprint(item["TransitRouterRouteTableType"])
			if rtType == "System" {
				log.Printf("[INFO] Skipping system route table (not deletable): %s (TR %s, CEN %s (%s))", rtId, trId, cenName, cenId)
				continue
			}
			log.Printf("[INFO] Deleting CEN TR route table: %s (TR %s, CEN %s (%s))", rtId, trId, cenName, cenId)
			delReq := map[string]interface{}{
				"TransitRouterRouteTableId": rtId,
				"ClientToken":               buildClientToken("DeleteTransitRouterRouteTable"),
			}
			var delResp map[string]interface{}
			delWait := incrementalWait(3*time.Second, 5*time.Second)
			delErr := retry.Retry(1*time.Minute, func() *retry.RetryError {
				r, e := listClient.RpcPost("Cbn", "2017-09-12", "DeleteTransitRouterRouteTable", nil, delReq, false)
				delResp = r
				if e != nil {
					// NotFound → already deleted, treat as success.
					if IsExpectedErrors(e, []string{"InvalidTransitRouterRouteTableId.NotFound"}) {
						return nil
					}
					// Residual associations/propagations/route entries are
					// non-retryable (no Force param). Log + continue.
					if IsExpectedErrors(e, []string{
						"OperationFailed.DeleteRouteTableWithAssociation",
						"OperationFailed.DeleteRouteTableWithPropagation",
						"OperationFailed.DeleteRouteTableWithRouteEntry",
						"OperationFailed.DeleteRouteTableWithRouteAggregation",
						"IncorrectStatus.TransitRouterRouteTable",
					}) {
						return retry.NonRetryableError(e)
					}
					if NeedRetry(e) {
						delWait()
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				return nil
			})
			addDebug("DeleteTransitRouterRouteTable", delResp, delReq)
			if delErr != nil {
				log.Printf("[ERROR] DeleteTransitRouterRouteTable failed for %s (TR %s, CEN %s (%s)): %v", rtId, trId, cenName, cenId, delErr)
			}
		}
		// NextToken pagination (MaxResults/NextToken, not PageSize/PageNumber).
		nextToken := ""
		if v, ok := resp["NextToken"]; ok && v != nil {
			if s := fmt.Sprint(v); s != "" {
				nextToken = s
			}
		}
		if nextToken == "" {
			break
		}
		rtReq["NextToken"] = nextToken
	}
}

// deleteTransitRouter force-deletes a transit router. Works cross-region by TR
// ID from the regional client that owns the TR's region. Non-fatal: errors are
// logged so the sweeper can continue to the next TR / CEN.
func deleteTransitRouter(listClient *connectivity.AliyunClient, cenId, cenName, trId string) {
	log.Printf("[INFO] Deleting CEN transit router: %s (CEN %s (%s))", trId, cenName, cenId)
	req := map[string]interface{}{
		"TransitRouterId": trId,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		resp, e := listClient.RpcPost("Cbn", "2017-09-12", "DeleteTransitRouter", nil, req, false)
		if e != nil {
			if IsExpectedErrors(e, []string{"InvalidTransitRouterId.NotFound"}) {
				return nil
			}
			if IsExpectedErrors(e, []string{"Operation.Blocking", "InvalidCenInstanceStatusNotAvailableForTransitRouterOperation"}) || NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		addDebug("DeleteTransitRouter", resp, req)
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] DeleteTransitRouter failed for %s (CEN %s (%s)): %v", trId, cenName, cenId, err)
	}
}

// sweepCenOrphanedPeerAttachments queries DescribeCenInterRegionBandwidthLimits
// and force-deletes any remaining TR Peer attachments by attachment ID. This is
// a fallback for when ListTransitRouterPeerAttachments returns
// IllegalParam.CenIdOrRegionId (which happens when the owning TR was already
// deleted, leaving an orphaned Peer attachment with a DataTransfer bandwidth
// allocation that blocks DeleteCen with InvalidOperation.AttachedChildInstanceExist).
// Non-fatal: per-attachment errors are logged and the loop continues.
func sweepCenOrphanedPeerAttachments(client *connectivity.AliyunClient, cenId, cenName string) {
	req := map[string]interface{}{
		"CenId":      cenId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		var resp map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(2*time.Minute, func() *retry.RetryError {
			r, e := client.RpcPost("Cbn", "2017-09-12", "DescribeCenInterRegionBandwidthLimits", nil, req, true)
			resp = r
			if e != nil {
				if NeedRetry(e) {
					wait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			addDebug("DescribeCenInterRegionBandwidthLimits", resp, req)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] DescribeCenInterRegionBandwidthLimits failed for CEN %s (%s): %v", cenName, cenId, err)
			return
		}
		att, jerr := jsonpath.Get("$.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit", resp)
		if jerr != nil {
			break
		}
		result, _ := att.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			attachId := fmt.Sprint(item["TransitRouterAttachmentId"])
			if attachId == "" || attachId == "<nil>" {
				continue
			}
			log.Printf("[INFO] Deleting orphaned TR Peer attachment: %s (CEN %s (%s))", attachId, cenName, cenId)
			delReq := map[string]interface{}{
				"TransitRouterAttachmentId": attachId,
				"Force":                     true,
			}
			delWait := incrementalWait(3*time.Second, 5*time.Second)
			delErr := retry.Retry(1*time.Minute, func() *retry.RetryError {
				_, e := client.RpcPost("Cbn", "2017-09-12", "DeleteTransitRouterPeerAttachment", nil, delReq, false)
				if e != nil {
					if IsExpectedErrors(e, []string{"InvalidTransitRouterAttachment.NotFound"}) {
						return nil
					}
					if IsExpectedErrors(e, []string{"Operation.Blocking"}) || NeedRetry(e) {
						delWait()
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				return nil
			})
			if delErr != nil {
				log.Printf("[ERROR] DeleteTransitRouterPeerAttachment failed for %s (CEN %s (%s)): %v", attachId, cenName, cenId, delErr)
			}
		}
		totalCount := 0
		if tc, ok := resp["TotalCount"]; ok && tc != nil {
			if f, ok := tc.(float64); ok {
				totalCount = int(f)
			}
		}
		if len(result) < PageSizeLarge || totalCount <= len(result) {
			break
		}
		req["PageNumber"] = req["PageNumber"].(int) + 1
	}
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// lintignore: R001
func TestUnitAliCloudCenInstance(t *testing.T) {
	p := Provider().ResourcesMap
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

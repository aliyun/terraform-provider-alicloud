package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_vpc_endpoint_linked_vpc", &resource.Sweeper{
		Name: "alicloud_cr_vpc_endpoint_linked_vpc",
		F:    testSweepCrVpcEndpointLinkedVpc,
	})
}

// testSweepCrVpcEndpointLinkedVpc cleans up Container Registry VPC-endpoint
// linked VPCs left on CR EE instances by acceptance tests. Linked VPCs are
// enumerated via GetInstanceVpcEndpoint (for both the Registry and Chart
// modules) and deleted with DeleteInstanceVpcEndpointLinkedVpc (mirroring the
// resource Delete). Instances marked NODELETING are preserved; in
// non-SWEEP_ALL mode only linked VPCs on tf-* named instances are removed.
func testSweepCrVpcEndpointLinkedVpc(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting Alicloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}

	prefixes := []string{"tf-testAcc", "tf_testAcc", "tftestAcc"}
	modules := []string{"Registry", "Chart"}

	pageNo, pageSize := 1, PageSizeLarge
	for {
		resp, lerr := crService.ListCrEEInstances(pageNo, pageSize)
		if lerr != nil {
			log.Printf("[ERROR] Failed to ListCrEEInstances: %s", lerr)
			return nil
		}
		for _, inst := range resp.Instances {
			if protectedSweepVpcName(inst.InstanceName) {
				log.Printf("[INFO] Skipping CR vpc endpoint linked vpc on instance %s (%s): protected (NODELETING)", inst.InstanceId, inst.InstanceName)
				continue
			}
			if !sweepAll() {
				match := false
				for _, p := range prefixes {
					if strings.HasPrefix(strings.ToLower(inst.InstanceName), strings.ToLower(p)) {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			for _, moduleName := range modules {
				listAction := "GetInstanceVpcEndpoint"
				listReq := map[string]interface{}{
					"RegionId":   client.RegionId,
					"InstanceId": inst.InstanceId,
					"ModuleName": moduleName,
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				var response map[string]interface{}
				rerr := retry.Retry(5*time.Minute, func() *retry.RetryError {
					response, err = client.RpcPost("cr", "2018-12-01", listAction, nil, listReq, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if rerr != nil {
					log.Printf("[ERROR] Failed to GetInstanceVpcEndpoint for instance %s module %s: %s", inst.InstanceId, moduleName, rerr)
					continue
				}
				raw, jerr := jsonpath.Get("$.LinkedVpcs", response)
				if jerr != nil {
					// No LinkedVpcs path → no linked VPCs for this module; skip.
					continue
				}
				result, _ := raw.([]interface{})
				for _, v := range result {
					item := v.(map[string]interface{})
					vpcId := fmt.Sprint(item["VpcId"])
					vswitchId := fmt.Sprint(item["VswitchId"])
					mod := fmt.Sprint(item["ModuleName"])
					if mod == "" || mod == "<nil>" {
						mod = moduleName
					}
					if protectedSweepVpcName(vpcId) {
						log.Printf("[INFO] Skipping CR vpc endpoint linked vpc %s:%s: protected (NODELETING)", inst.InstanceId, vpcId)
						continue
					}
					log.Printf("[INFO] Deleting CR vpc endpoint linked vpc: %s:%s:%s:%s", inst.InstanceId, vpcId, vswitchId, mod)
					delReq := map[string]interface{}{
						"RegionId":   client.RegionId,
						"InstanceId": inst.InstanceId,
						"VpcId":      vpcId,
						"VswitchId":  vswitchId,
						"ModuleName": mod,
					}
					dwait := incrementalWait(3*time.Second, 3*time.Second)
					if derr := retry.Retry(5*time.Minute, func() *retry.RetryError {
						_, e := client.RpcPost("cr", "2018-12-01", "DeleteInstanceVpcEndpointLinkedVpc", nil, delReq, true)
						if e != nil {
							if NotFoundError(e) || IsExpectedErrors(e, []string{"INSTANCE_ACCESS_NOT_EXIST"}) {
								return nil
							}
							if NeedRetry(e) {
								dwait()
								return retry.RetryableError(e)
							}
							return retry.NonRetryableError(e)
						}
						return nil
					}); derr != nil {
						log.Printf("[ERROR] Failed to delete CR vpc endpoint linked vpc %s:%s: %s", inst.InstanceId, vpcId, derr)
						continue
					}
				}
			}
		}
		if len(resp.Instances) < pageSize {
			break
		}
		pageNo++
	}
	return nil
}

func TestAccAlicloudCrVpcEndpointLinkedVpc_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_vpc_endpoint_linked_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudCrVpcEndpointLinkedVpcMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrVpcEndpointLinkedVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCrVpcEndpointLinkedVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrVpcEndpointLinkedVpcBasicDependence0)
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
					"instance_id":                      "${data.alicloud_cr_ee_instances.default.ids.0}",
					"vpc_id":                           "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                       "${data.alicloud_vswitches.default.ids.0}",
					"module_name":                      "Registry",
					"enable_create_dns_record_in_pvzt": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"vpc_id":      CHECKSET,
						"vswitch_id":  CHECKSET,
						"module_name": "Registry",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_create_dns_record_in_pvzt"},
			},
		},
	})
}

var AlicloudCrVpcEndpointLinkedVpcMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudCrVpcEndpointLinkedVpcBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_cr_ee_instances" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
`, name)
}

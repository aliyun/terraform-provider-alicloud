package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_internet_endpoint", &resource.Sweeper{
		Name:         "alicloud_cr_internet_endpoint",
		Dependencies: []string{"alicloud_cr_endpoint_acl_policy"},
		F:            testSweepCrInternetEndpoint,
	})
}

// testSweepCrInternetEndpoint disables the internet access endpoint on CR EE
// instances left enabled by acceptance tests. alicloud_cr_internet_endpoint has
// no real Delete API — its destroy sets Enable=false via
// UpdateInstanceEndpointStatus — so the sweeper mirrors that (fire-and-forget
// disable, idempotent). Endpoint ACL policies are swept first (Dependency).
// Instances marked NODELETING are preserved; in non-SWEEP_ALL mode only tf-*
// named instances are touched.
func testSweepCrInternetEndpoint(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting Alicloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}

	prefixes := []string{"tf-testAcc", "tf_testAcc", "tftestAcc"}

	pageNo, pageSize := 1, PageSizeLarge
	for {
		resp, lerr := crService.ListCrEEInstances(pageNo, pageSize)
		if lerr != nil {
			log.Printf("[ERROR] Failed to ListCrEEInstances: %s", lerr)
			return nil
		}
		for _, inst := range resp.Instances {
			if protectedSweepVpcName(inst.InstanceName) {
				log.Printf("[INFO] Skipping CR internet endpoint on instance %s (%s): protected (NODELETING)", inst.InstanceId, inst.InstanceName)
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
			log.Printf("[INFO] Disabling CR internet endpoint on instance: %s (%s)", inst.InstanceId, inst.InstanceName)
			delReq := map[string]interface{}{
				"RegionId":     client.RegionId,
				"InstanceId":   inst.InstanceId,
				"Enable":       "false",
				"ModuleName":   "Registry",
				"EndpointType": "Internet",
			}
			dwait := incrementalWait(3*time.Second, 3*time.Second)
			if derr := retry.Retry(5*time.Minute, func() *retry.RetryError {
				_, e := client.RpcPost("cr", "2018-12-01", "UpdateInstanceEndpointStatus", map[string]interface{}{}, delReq, true)
				if e != nil {
					if NotFoundError(e) {
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
				log.Printf("[ERROR] Failed to disable CR internet endpoint on instance %s: %s", inst.InstanceId, derr)
				continue
			}
		}
		if len(resp.Instances) < pageSize {
			break
		}
		pageNo++
	}
	return nil
}

// TestAccAliCloudCrInternetEndpoint_basic0 covers the full lifecycle:
// create with two ACL entries -> update to a single entry (add + remove) ->
// update to no entries (remove all) -> import.
//
// The config is expressed as an inline map[string]interface{} (built by
// resourceTestAccConfigFunc) rather than a hand-rolled fmt.Sprintf HCL string
// so that scripts/testing/testing_coverage_rate_check.go can parse the
// `entries` set and its nested `entry`/`comment` attributes for both set and
// modify coverage. The generator-emitted HCL is semantically identical to the
// previous hand-written blocks.
func TestAccAliCloudCrInternetEndpoint_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_internet_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudCrInternetEndpointMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrInternetEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrInternetEndpointBasicDependence0)
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
					"instance_id": "${local.cr_endpoint_instance_id}",
					"entries": []map[string]interface{}{
						{
							"entry":   "192.168.1.0/24",
							"comment": "entry-1",
						},
						{
							"entry":   "10.0.0.0/8",
							"comment": "entry-2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"status":      CHECKSET,
					}),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "2"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entries": []map[string]interface{}{
						{
							"entry":   "172.16.0.0/12",
							"comment": "entry-3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"status":      CHECKSET,
					}),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "1"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entries": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"status":      CHECKSET,
					}),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "0"),
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

var AlicloudCrInternetEndpointMap0 = map[string]string{}

func AlicloudCrInternetEndpointBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Economy"
  instance_name  = var.name
  image_scanner  = "DISABLE"
}

locals {
  cr_endpoint_instance_id = alicloud_cr_ee_instance.default.id
}
`, name)
}

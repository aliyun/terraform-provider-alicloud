package directmail_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/directmail"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// directMailSweeperPrefixes lists the name prefixes of the resources created
// by the DirectMail acceptance tests and by the documented examples. An
// interrupted run can leave them behind, so the sweepers must recognise every
// one of them; the documented examples use "terraform-example".
var directMailSweeperPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"terraform-example",
}

// directMailNameShouldSweep reports whether a resource name was created by a
// test or a documented example and therefore may be swept. Matching stays
// case-insensitive and prefix-only, so unrelated names are never deleted.
func directMailNameShouldSweep(name string) bool {
	for _, prefix := range directMailSweeperPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

func init() {
	resource.AddTestSweepers(
		"alicloud_direct_mail_dedicated_ip_pool",
		&resource.Sweeper{
			Name: "alicloud_direct_mail_dedicated_ip_pool",
			F:    testSweepDirectMailDedicatedIpPool,
		})
}

// TestUnitDirectMailSweeperPrefixes pins the sweeper name matching used by both
// DirectMail sweepers: leftovers of the documented examples (whose default name
// is "terraform-example") must be recognised after an interrupted run, while
// unrelated production names stay untouched.
func TestUnitDirectMailSweeperPrefixes(t *testing.T) {
	swept := []string{
		"terraform-example",
		"terraform-example-pool",
		"Terraform-Example",
		"tf-testAcc-direct-mail-pool-12345",
		"tf_testAcc_direct_mail_config_set",
	}
	kept := []string{
		"prod-config-set",
		"customer-dedicated-pool",
		"",
	}
	for _, name := range swept {
		if !directMailNameShouldSweep(name) {
			t.Fatalf("expected %q to be swept with prefixes %v", name, directMailSweeperPrefixes)
		}
	}
	for _, name := range kept {
		if directMailNameShouldSweep(name) {
			t.Fatalf("expected %q to be kept with prefixes %v", name, directMailSweeperPrefixes)
		}
	}
}

func testSweepDirectMailDedicatedIpPool(region string) error {
	rawClient, err := sweep.SharedClientForRegionWithBackendRegions(region, true, connectivity.DmSupportRegions)
	if err != nil {
		log.Printf("error getting Alicloud client: %s", err)
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "DedicatedIpPoolList"
	request := map[string]interface{}{
		"PageIndex": 1,
		"PageSize":  helper.PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			if errs.NotFoundError(err) {
				log.Printf("[INFO] DirectMail Dedicated Ip Pool list reported not found; nothing to sweep: %s", err)
				return nil
			}
			log.Printf("[ERROR] Failed to fetch DirectMail Dedicated Ip Pool: %s", err)
			return nil
		}
		v, err := jsonpath.Get("$.IpPools", response)
		if err != nil {
			return nil
		}
		result, _ := v.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["Name"])
			if !directMailNameShouldSweep(name) {
				log.Printf("[INFO] Skipping DirectMail Dedicated Ip Pool: %s", name)
				continue
			}
			action := "DedicatedIpPoolDelete"
			deleteRequest := map[string]interface{}{
				"Id": item["Id"],
			}
			_, err = client.RpcPost("Dm", "2015-11-23", action, nil, deleteRequest, true)
			if err != nil {
				if errs.NotFoundError(err) {
					log.Printf("[INFO] DirectMail Dedicated Ip Pool %s was already removed: %s", name, err)
				} else {
					log.Printf("[ERROR] Failed to delete DirectMail Dedicated Ip Pool (%s): %s", name, err)
				}
			} else {
				log.Printf("[INFO] Delete DirectMail Dedicated Ip Pool success: %s", name)
			}
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageIndex"] = request["PageIndex"].(int) + 1
	}
	return nil
}

// Creating a dedicated IP pool requires purchased dedicated IP instances in
// the account (pool quota derives from purchased IPs); the shared tf
// acceptance test account owns none, so this case skips at runtime unless
// ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN is set. The steps below only exercise
// an empty pool; attaching purchased instances is covered by
// TestAccAliCloudDirectMailDedicatedIpPool_singleMember and
// TestAccAliCloudDirectMailDedicatedIpPool_buyResourceIds because member IDs
// cannot be fabricated.
func TestAccAliCloudDirectMailDedicatedIpPool_basic(t *testing.T) {
	if os.Getenv("ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN") == "" {
		t.Skip("requires purchased dedicated IPs; set ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN=1 to run")
	}
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_dedicated_ip_pool.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDirectMailDedicatedIpPoolMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return directmail.NewDirectMailServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDirectMailDedicatedIpPool")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-pool-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDirectMailDedicatedIpPoolBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		// No IDRefreshName here: the SDK v2 ID-only refresh rebuilds a minimal
		// provider config, which would drop the assume_role block and query the
		// base test account instead of the target role account.
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     name,
						"ip_count": "0",
						"ips.#":    "0",
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

// A pool created with a single purchased member; runs with one or more IDs.
func TestAccAliCloudDirectMailDedicatedIpPool_singleMember(t *testing.T) {
	if os.Getenv("ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN") == "" {
		t.Skip("requires purchased dedicated IPs; set ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN=1 to run")
	}
	buyResourceIds := directmail.DirectMailBuyResourceIdSet(os.Getenv("ALICLOUD_DIRECT_MAIL_DEDICATED_IP_IDS"))
	if len(buyResourceIds) == 0 {
		t.Skip("requires a purchased dedicated IP instance ID; set ALICLOUD_DIRECT_MAIL_DEDICATED_IP_IDS to run")
	}
	member := buyResourceIds[0]
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_dedicated_ip_pool.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDirectMailDedicatedIpPoolMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return directmail.NewDirectMailServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDirectMailDedicatedIpPool")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-pool-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDirectMailDedicatedIpPoolBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": member,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"buy_resource_ids": member,
						"ip_count":         "1",
						"ips.#":            "1",
						"ips.0.id":         member,
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

// Membership updates need at least two distinct purchased IP instance IDs:
// with a single ID every step would carry the same member and the update path
// would never run, so the case is skipped (unverified) instead.
func TestAccAliCloudDirectMailDedicatedIpPool_buyResourceIds(t *testing.T) {
	if os.Getenv("ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN") == "" {
		t.Skip("requires purchased dedicated IPs; set ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN=1 to run")
	}
	buyResourceIds := directmail.DirectMailBuyResourceIdSet(os.Getenv("ALICLOUD_DIRECT_MAIL_DEDICATED_IP_IDS"))
	buyResourceIds = directMailDistinct(buyResourceIds)
	if len(buyResourceIds) < 2 {
		t.Skip("requires at least two distinct purchased dedicated IP instance IDs in ALICLOUD_DIRECT_MAIL_DEDICATED_IP_IDS; membership update is unverified")
	}
	expectedFirst := buyResourceIds[0]
	expectedSecond := buyResourceIds[1]
	expected := expectedFirst + "," + expectedSecond
	// The same two members in reverse order with surrounding spaces exercise
	// the order-insensitive normalization of buy_resource_ids.
	reversedSpaced := expectedSecond + ", " + expectedFirst
	var poolId string
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_dedicated_ip_pool.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDirectMailDedicatedIpPoolMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return directmail.NewDirectMailServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDirectMailDedicatedIpPool")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-pool-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDirectMailDedicatedIpPoolBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		// No IDRefreshName here: the SDK v2 ID-only refresh rebuilds a minimal
		// provider config, which would drop the assume_role block and query the
		// base test account instead of the target role account.
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": expectedFirst,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"buy_resource_ids": expectedFirst,
						"ip_count":         "1",
						"ips.#":            "1",
						"ips.0.id":         expectedFirst,
					}),
					testAccCheckDirectMailDedicatedIpPoolSameId(resourceId, &poolId),
				),
			},
			{
				// Swapping the only member distinguishes replacement from
				// append: append would leave both members in the pool.
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": expectedSecond,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"buy_resource_ids": expectedSecond,
						"ip_count":         "1",
						"ips.#":            "1",
						"ips.0.id":         expectedSecond,
					}),
					testAccCheckDirectMailDedicatedIpPoolSameId(resourceId, &poolId),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": reversedSpaced,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"buy_resource_ids": expected,
						"ip_count":         "2",
						"ips.#":            "2",
						"ips.0.id":         acctest.REMOVEKEY,
					}),
					testAccCheckDirectMailDedicatedIpPoolSameId(resourceId, &poolId),
				),
			},
			{
				// Re-planning the reversed, space-separated configuration must
				// produce an empty plan (no permanent update loop).
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": reversedSpaced,
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Shrinking to one member must remove the other one.
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": expectedFirst,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"buy_resource_ids": expectedFirst,
						"ip_count":         "1",
						"ips.#":            "1",
						"ips.0.id":         expectedFirst,
					}),
					testAccCheckDirectMailDedicatedIpPoolSameId(resourceId, &poolId),
				),
			},
			{
				// All members are removed outside Terraform: refresh must record
				// the empty membership instead of failing, and the plan must
				// propose restoring the configured member.
				PreConfig: func() {
					testAccDirectMailDedicatedIpPoolClearMembersExternally(t, &poolId)
				},
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": expectedFirst,
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// The regular update reconciles the pool without replacing it.
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": expectedFirst,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"buy_resource_ids": expectedFirst,
						"ip_count":         "1",
						"ips.#":            "1",
					}),
					testAccCheckDirectMailDedicatedIpPoolSameId(resourceId, &poolId),
				),
			},
			{
				// Removing buy_resource_ids from the configuration clears the
				// membership.
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"buy_resource_ids": "",
						"ip_count":         "0",
						"ips.#":            "0",
						"ips.0.id":         acctest.REMOVEKEY,
					}),
					testAccCheckDirectMailDedicatedIpPoolSameId(resourceId, &poolId),
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

func directMailDistinct(items []string) []string {
	result := make([]string, 0, len(items))
	for i, item := range items {
		if i > 0 && item == items[i-1] {
			continue
		}
		result = append(result, item)
	}
	return result
}

// testAccCheckDirectMailDedicatedIpPoolSameId records the pool ID on first use
// and fails if a later step replaced the pool.
func testAccCheckDirectMailDedicatedIpPoolSameId(resourceId string, poolId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if *poolId == "" {
			*poolId = rs.Primary.ID
			return nil
		}
		if rs.Primary.ID != *poolId {
			return fmt.Errorf("pool %s was replaced by %s", *poolId, rs.Primary.ID)
		}
		return nil
	}
}

// testAccDirectMailDedicatedIpPoolClearMembersExternally removes every member
// with the provider's client (the same target-role credentials) and waits
// until the pool reports no members.
func testAccDirectMailDedicatedIpPoolClearMembersExternally(t *testing.T, poolId *string) {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"Id":             *poolId,
		"RegionId":       client.RegionId,
		"BuyResourceIds": "",
		"UpdateResource": true,
	}
	if _, err := client.RpcPost("Dm", "2015-11-23", "DedicatedIpPoolUpdate", nil, request, true); err != nil {
		t.Fatalf("clearing members of pool %s outside Terraform failed: %s", *poolId, err)
	}
	service := directmail.NewDirectMailServiceV2(client)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		object, err := service.DescribeDirectMailDedicatedIpPool(*poolId)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if ips, _ := object["Ips"].([]interface{}); len(ips) > 0 {
			return retry.RetryableError(fmt.Errorf("pool %s still has %d members", *poolId, len(ips)))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("waiting for pool %s members to be cleared failed: %s", *poolId, err)
	}
}

// testAccDirectMailAssumeRoleProviderConfig returns the provider block used by
// the DirectMail dedicated-IP guard cases. The base credentials still come from
// the ambient ALICLOUD_ACCESS_KEY/ALICLOUD_SECRET_KEY, so the shared test
// identity is untouched; the provider then assumes the target role and every
// client (create, read, import and destroy) runs with that role's STS
// credentials, including the security token. When the role ARN is not provided
// the block is empty and behaviour is unchanged.
func testAccDirectMailAssumeRoleProviderConfig() string {
	roleArn := os.Getenv("ALICLOUD_DIRECT_MAIL_ASSUME_ROLE_ARN")
	if roleArn == "" {
		return ""
	}
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = acctest.DefaultRegionToTest
	}
	return fmt.Sprintf(`
provider "alicloud" {
  region = "%s"
  assume_role {
    role_arn     = "%s"
    session_name = "terraform-acc"
  }
}
`, region, roleArn)
}

var AlicloudDirectMailDedicatedIpPoolMap0 = map[string]string{
	"create_time": acctest.CHECKSET,
	"ip_count":    acctest.CHECKSET,
}

func AlicloudDirectMailDedicatedIpPoolBasicDependence0(name string) string {
	return testAccDirectMailAssumeRoleProviderConfig()
}

// directMailTestClient returns a client whose DirectMail RPC calls go to a
// local httptest server, so request parameters and responses are observable
// without cloud credentials.
func directMailTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("dm", endpoint)
	return client
}

func directMailWriteJSON(t *testing.T, w http.ResponseWriter, body map[string]interface{}) {
	t.Helper()
	body["RequestId"] = "test-request"
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Error(err)
	}
}

func TestUnitDirectMailDedicatedIpPoolReadEmptyMembership(t *testing.T) {
	client := directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if action := r.Form.Get("Action"); action != "DedicatedIpPoolList" {
			t.Errorf("unexpected action %s", action)
		}
		directMailWriteJSON(t, w, map[string]interface{}{
			"IpPools": []interface{}{
				map[string]interface{}{"Id": "pool-1", "Name": "pool", "IpCount": 0, "Ips": []interface{}{}},
			},
		})
	})
	d := schema.TestResourceDataRaw(t, directmail.ResourceDedicatedIpPool().Schema, map[string]interface{}{
		"name":             "pool",
		"buy_resource_ids": "ip-a",
	})
	d.SetId("pool-1")
	if err := directmail.ResourceDedicatedIpPoolRead(d, client); err != nil {
		t.Fatalf("refresh of a pool whose members were removed must succeed, got: %s", err)
	}
	if d.Id() != "pool-1" {
		t.Fatalf("pool ID must be kept, got %q", d.Id())
	}
	if got := d.Get("buy_resource_ids").(string); got != "" {
		t.Fatalf("buy_resource_ids must record the empty membership, got %q", got)
	}
}

func TestUnitDirectMailDedicatedIpPoolReadNormalizesMembership(t *testing.T) {
	client := directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		directMailWriteJSON(t, w, map[string]interface{}{
			"IpPools": []interface{}{
				map[string]interface{}{"Id": "pool-1", "Name": "pool", "IpCount": 2, "Ips": []interface{}{
					map[string]interface{}{"Id": "ip-b", "Ip": "192.0.2.2"},
					map[string]interface{}{"Id": "ip-a", "Ip": "192.0.2.1"},
				}},
			},
		})
	})
	d := schema.TestResourceDataRaw(t, directmail.ResourceDedicatedIpPool().Schema, map[string]interface{}{"name": "pool"})
	d.SetId("pool-1")
	if err := directmail.ResourceDedicatedIpPoolRead(d, client); err != nil {
		t.Fatal(err)
	}
	if got := d.Get("buy_resource_ids").(string); got != "ip-a,ip-b" {
		t.Fatalf("expected sorted membership ip-a,ip-b, got %q", got)
	}
}

func TestUnitDirectMailDescribeDedicatedIpPoolExactId(t *testing.T) {
	pools := []interface{}{
		map[string]interface{}{"Id": "pool-other", "Name": "other"},
	}
	client := directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if poolId := r.Form.Get("PoolId"); poolId == "" {
			t.Errorf("PoolId must be sent")
		}
		directMailWriteJSON(t, w, map[string]interface{}{"IpPools": pools})
	})
	service := directmail.NewDirectMailServiceV2(client)
	if _, err := service.DescribeDirectMailDedicatedIpPool("pool-1"); !errs.NotFoundError(err) {
		t.Fatalf("another pool in the response must not be adopted; expected NotFound, got %v", err)
	}
	pools = append(pools, map[string]interface{}{"Id": "pool-1", "Name": "target"})
	object, err := service.DescribeDirectMailDedicatedIpPool("pool-1")
	if err != nil {
		t.Fatal(err)
	}
	if object["Name"] != "target" {
		t.Fatalf("expected the exact pool, got %v", object)
	}
}

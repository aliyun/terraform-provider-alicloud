package alicloud

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccDirectMailDedicatedIpPoolsOutputFile = "./tf-testAcc-direct-mail-pools-output.json"

func TestAccAliCloudDirectMailDedicatedIpPoolsDataSource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_dedicated_ip_pool.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DirectMailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailDedicatedIpPool")
	ra := resourceAttrInit(resourceId, map[string]string{})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-pool-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_direct_mail_dedicated_ip_pools.default"),
					resource.TestCheckResourceAttrSet("data.alicloud_direct_mail_dedicated_ip_pools.default", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#"),
				),
			},
			{
				// A pool_id that matches no pool must return an empty result
				// instead of any existing pool in the same account.
				Config: fmt.Sprintf(`
data "alicloud_direct_mail_dedicated_ip_pools" "default" {
  pool_id = "%s"
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "ids.#", "0"),
				),
			},
			{
				Config: fmt.Sprintf(`
data "alicloud_direct_mail_dedicated_ip_pools" "default" {
  keyword = "%s"
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#", "0"),
				),
			},
			{
				// all is sent to the list API and output_file must write the
				// (empty) result list to disk.
				PreConfig: func() {
					os.Remove(testAccDirectMailDedicatedIpPoolsOutputFile)
				},
				Config: fmt.Sprintf(`
data "alicloud_direct_mail_dedicated_ip_pools" "default" {
  keyword     = "%s"
  all         = true
  output_file = "%s"
}
`, name, testAccDirectMailDedicatedIpPoolsOutputFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "all", "true"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "output_file", testAccDirectMailDedicatedIpPoolsOutputFile),
					testAccCheckDirectMailOutputFileContent(testAccDirectMailDedicatedIpPoolsOutputFile, "[]"),
				),
			},
			{
				// A non-matching ID must go through the client-side ids
				// filter and yield an empty result.
				Config: fmt.Sprintf(`
data "alicloud_direct_mail_dedicated_ip_pools" "default" {
  ids = ["tf-testAcc-direct-mail-pool-missing-%d"]
}
`, rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#", "0"),
				),
			},
			{
				// An unwritable output path must fail the read instead of
				// reporting success without the export.
				Config: fmt.Sprintf(`
data "alicloud_direct_mail_dedicated_ip_pools" "default" {
  keyword     = "%s"
  output_file = "./tf-testAcc-direct-mail-missing-dir-%d/pools.json"
}
`, name, rand),
				ExpectError: regexp.MustCompile("no such file or directory"),
			},
		},
	})
}

// Discriminating two pools in the same scope requires real pool objects;
// pool creation is rejected in accounts without purchased dedicated IPs, so
// this case skips at runtime unless ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN is set.
func TestAccAliCloudDirectMailDedicatedIpPoolsDataSource_poolId(t *testing.T) {
	if os.Getenv("ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN") == "" {
		t.Skip("requires purchased dedicated IPs; set ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN=1 to run")
	}
	var firstId string
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_dedicated_ip_pool.first"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DirectMailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailDedicatedIpPool")
	ra := resourceAttrInit(resourceId, map[string]string{})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-pool-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourcePoolIdConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.0.name", name+"-a"),
					resource.TestCheckResourceAttrPair("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.0.id", resourceId, "id"),
					testAccCheckDirectMailDedicatedIpPoolSameId(resourceId, &firstId),
				),
			},
			{
				Config: testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourcePoolIdConfig(name) + `
data "alicloud_direct_mail_dedicated_ip_pools" "second" {
  pool_id = alicloud_direct_mail_dedicated_ip_pool.second.id
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.second", "pools.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.second", "pools.0.name", name+"-b"),
				),
			},
			{
				// ids filters client-side by exact ID.
				Config: testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceIdsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.0.name", name+"-a"),
				),
			},
			{
				// A non-matching ID must exclude the two real pools.
				Config: testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceIdsMissingConfig(name, rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#", "0"),
				),
			},
			{
				// The first pool is deleted while the second one remains: the
				// deleted ID must be reported as not found by the resource
				// lookup and the data source, and must not resolve to the
				// remaining pool.
				Config: testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceRemainingConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.0.name", name+"-b"),
					testAccCheckDirectMailDedicatedIpPoolGone(t, &firstId),
				),
			},
		},
	})
}

// testAccCheckDirectMailOutputFileContent verifies that output_file was really
// written by the data source in the test process working directory.
func testAccCheckDirectMailOutputFileContent(path string, want ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("output_file %s was not written: %s", path, err)
		}
		defer os.Remove(path)
		for _, item := range want {
			if !strings.Contains(string(content), item) {
				return fmt.Errorf("output_file %s does not contain %q: %s", path, item, content)
			}
		}
		return nil
	}
}

func testAccCheckDirectMailDedicatedIpPoolGone(t *testing.T, poolId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		service := DirectMailServiceV2{client}
		if _, err := service.DescribeDirectMailDedicatedIpPool(*poolId); !NotFoundError(err) {
			return fmt.Errorf("deleted pool %s must be not found, got: %v", *poolId, err)
		}
		d := schema.TestResourceDataRaw(t, dataSourceAliCloudDirectMailDedicatedIpPools().Schema, map[string]interface{}{
			"pool_id": *poolId,
		})
		if err := dataSourceAliCloudDirectMailDedicatedIpPoolRead(d, client); err != nil {
			return err
		}
		if pools := d.Get("pools").([]interface{}); len(pools) != 0 {
			return fmt.Errorf("pool_id %s of a deleted pool must return no pools, got %v", *poolId, pools)
		}
		return nil
	}
}

func testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourcePoolIdConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_direct_mail_dedicated_ip_pool" "first" {
  name = "%s-a"
}

resource "alicloud_direct_mail_dedicated_ip_pool" "second" {
  name = "%s-b"
}

data "alicloud_direct_mail_dedicated_ip_pools" "default" {
  pool_id = alicloud_direct_mail_dedicated_ip_pool.first.id
}
`, testAccDirectMailAssumeRoleProviderConfig(), name, name)
}

func testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceIdsConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_direct_mail_dedicated_ip_pool" "first" {
  name = "%s-a"
}

resource "alicloud_direct_mail_dedicated_ip_pool" "second" {
  name = "%s-b"
}

data "alicloud_direct_mail_dedicated_ip_pools" "default" {
  ids = [alicloud_direct_mail_dedicated_ip_pool.first.id]
}
`, testAccDirectMailAssumeRoleProviderConfig(), name, name)
}

func testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceIdsMissingConfig(name string, rand int) string {
	return fmt.Sprintf(`
%s

resource "alicloud_direct_mail_dedicated_ip_pool" "first" {
  name = "%s-a"
}

resource "alicloud_direct_mail_dedicated_ip_pool" "second" {
  name = "%s-b"
}

data "alicloud_direct_mail_dedicated_ip_pools" "default" {
  ids = ["tf-testAcc-direct-mail-pool-missing-%d"]
}
`, testAccDirectMailAssumeRoleProviderConfig(), name, name, rand)
}

func testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceRemainingConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_direct_mail_dedicated_ip_pool" "second" {
  name = "%s-b"
}

data "alicloud_direct_mail_dedicated_ip_pools" "default" {
  pool_id = alicloud_direct_mail_dedicated_ip_pool.second.id
}
`, testAccDirectMailAssumeRoleProviderConfig(), name)
}

const testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceConfig = `
data "alicloud_direct_mail_dedicated_ip_pools" "default" {
}
`

func TestUnitDirectMailDedicatedIpPoolsDataSourcePoolIdAndOutputFile(t *testing.T) {
	client := directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Form.Get("PoolId") != "pool-2" {
			t.Errorf("PoolId must be sent, got %q", r.Form.Get("PoolId"))
		}
		directMailWriteJSON(t, w, map[string]interface{}{
			"IpPools": []interface{}{
				map[string]interface{}{"Id": "pool-1", "Name": "first"},
				map[string]interface{}{"Id": "pool-2", "Name": "second"},
			},
		})
	})
	output := filepath.Join(t.TempDir(), "pools.json")
	d := schema.TestResourceDataRaw(t, dataSourceAliCloudDirectMailDedicatedIpPools().Schema, map[string]interface{}{
		"pool_id":     "pool-2",
		"output_file": output,
	})
	if err := dataSourceAliCloudDirectMailDedicatedIpPoolRead(d, client); err != nil {
		t.Fatal(err)
	}
	pools := d.Get("pools").([]interface{})
	if len(pools) != 1 || pools[0].(map[string]interface{})["id"] != "pool-2" {
		t.Fatalf("pool_id must select only the exact pool, got %v", pools)
	}
	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("output_file was not written: %s", err)
	}
	if !strings.Contains(string(content), `"pool-2"`) || strings.Contains(string(content), `"pool-1"`) {
		t.Fatalf("unexpected output_file content: %s", content)
	}

	unwritable := schema.TestResourceDataRaw(t, dataSourceAliCloudDirectMailDedicatedIpPools().Schema, map[string]interface{}{
		"pool_id":     "pool-2",
		"output_file": filepath.Join(t.TempDir(), "missing-dir", "pools.json"),
	})
	if err := dataSourceAliCloudDirectMailDedicatedIpPoolRead(unwritable, client); err == nil {
		t.Fatal("an unwritable output_file must return an error")
	}
}

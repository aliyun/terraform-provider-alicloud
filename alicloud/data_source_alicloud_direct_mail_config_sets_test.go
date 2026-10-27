package alicloud

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDirectMailConfigSetsOutputFile = "./tf-testAcc-direct-mail-config-sets-output.json"

func TestAccAliCloudDirectMailConfigSetsDataSource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_config_set.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-cs-%d", rand)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DirectMailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailConfigSet")
	ra := resourceAttrInit(resourceId, map[string]string{})
	rac := resourceAttrCheckInit(rc, ra)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDirectMailConfigSetsDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_direct_mail_config_sets.default"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.0.name", name),
					resource.TestCheckResourceAttrSet("data.alicloud_direct_mail_config_sets.default", "ids.#"),
				),
			},
			{
				// all is sent to the list API, and output_file must write the
				// matching configuration set to disk.
				PreConfig: func() {
					os.Remove(testAccDirectMailConfigSetsOutputFile)
				},
				Config: testAccCheckAlicloudDirectMailConfigSetsDataSourceAllConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_direct_mail_config_sets.default"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "all", "true"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "output_file", testAccDirectMailConfigSetsOutputFile),
					testAccCheckDirectMailOutputFileContent(testAccDirectMailConfigSetsOutputFile, name),
				),
			},
			{
				// An unwritable output path must fail the read instead of
				// reporting success without the export.
				Config:      testAccCheckAlicloudDirectMailConfigSetsDataSourceUnwritableConfig(name, rand),
				ExpectError: regexp.MustCompile("no such file or directory"),
			},
			{
				// ids filters client-side by exact ID.
				Config: testAccCheckAlicloudDirectMailConfigSetsDataSourceIdsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_direct_mail_config_sets.default"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.0.name", name),
				),
			},
			{
				// A non-matching ID must yield an empty result.
				Config: testAccCheckAlicloudDirectMailConfigSetsDataSourceIdsMissingConfig(name, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_direct_mail_config_sets.default"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.#", "0"),
				),
			},
			{
				// enable_details makes the data source call ConfigSetDetail per
				// item, so this step exercises that branch against the cloud.
				// The list and detail responses expose the same fields, so the
				// branch itself is distinguished by request observation in
				// TestUnitDirectMailConfigSetsDataSourceEnableDetails. Keep the
				// rich configuration last: removing validation_option would be
				// a real in-place update, not a test reset.
				Config: testAccCheckAlicloudDirectMailConfigSetsDataSourceDetailsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_direct_mail_config_sets.default"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "enable_details", "true"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.0.name", name),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.0.description", "tf-testAcc-description"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.0.validation_option.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.0.validation_option.0.enabled", "true"),
				),
			},
			{
				// A keyword that matches no configuration set must return an
				// empty result instead of any existing set in the same account.
				Config: fmt.Sprintf(`
data "alicloud_direct_mail_config_sets" "default" {
  keyword = "tf-testAcc-direct-mail-cs-none-%d"
}
`, rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_direct_mail_config_sets.default", "sets.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDirectMailConfigSetsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "alicloud_direct_mail_config_set" "default" {
  name = "%s"
}

data "alicloud_direct_mail_config_sets" "default" {
  keyword = alicloud_direct_mail_config_set.default.name
}
`, name)
}

func testAccCheckAlicloudDirectMailConfigSetsDataSourceAllConfig(name string) string {
	return fmt.Sprintf(`
resource "alicloud_direct_mail_config_set" "default" {
  name = "%s"
}

data "alicloud_direct_mail_config_sets" "default" {
  keyword     = alicloud_direct_mail_config_set.default.name
  all         = true
  output_file = "%s"
}
`, name, testAccDirectMailConfigSetsOutputFile)
}

func testAccCheckAlicloudDirectMailConfigSetsDataSourceUnwritableConfig(name string, rand int) string {
	return fmt.Sprintf(`
resource "alicloud_direct_mail_config_set" "default" {
  name = "%s"
}

data "alicloud_direct_mail_config_sets" "default" {
  keyword     = alicloud_direct_mail_config_set.default.name
  output_file = "./tf-testAcc-direct-mail-missing-dir-%d/sets.json"
}
`, name, rand)
}

func testAccCheckAlicloudDirectMailConfigSetsDataSourceIdsConfig(name string) string {
	return fmt.Sprintf(`
resource "alicloud_direct_mail_config_set" "default" {
  name = "%s"
}

data "alicloud_direct_mail_config_sets" "default" {
  ids = [alicloud_direct_mail_config_set.default.id]
}
`, name)
}

func testAccCheckAlicloudDirectMailConfigSetsDataSourceIdsMissingConfig(name string, rand int) string {
	return fmt.Sprintf(`
resource "alicloud_direct_mail_config_set" "default" {
  name = "%s"
}

data "alicloud_direct_mail_config_sets" "default" {
  ids = ["tf-testAcc-direct-mail-cs-missing-%d"]
}
`, name, rand)
}

func testAccCheckAlicloudDirectMailConfigSetsDataSourceDetailsConfig(name string) string {
	return fmt.Sprintf(`
resource "alicloud_direct_mail_config_set" "default" {
  name        = "%s"
  description = "tf-testAcc-description"

  validation_option {
    enabled               = true
    forbidden_status_list = ["INVALID"]
  }
}

data "alicloud_direct_mail_config_sets" "default" {
  keyword        = alicloud_direct_mail_config_set.default.name
  enable_details = true

  depends_on = [alicloud_direct_mail_config_set.default]
}
`, name)
}

// The list and detail responses carry different descriptions here, so the
// assertions show which API produced the attributes and that ConfigSetDetail is
// called only when enable_details is true.
func TestUnitDirectMailConfigSetsDataSourceEnableDetails(t *testing.T) {
	for _, enableDetails := range []bool{false, true} {
		var detailCalls atomic.Int32
		client := directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			switch r.Form.Get("Action") {
			case "ConfigSetList":
				directMailWriteJSON(t, w, map[string]interface{}{
					"ConfigSets": []interface{}{
						map[string]interface{}{"Id": "cs-1", "Name": "cs", "Description": "from-list"},
					},
				})
			case "ConfigSetDetail":
				detailCalls.Add(1)
				if r.Form.Get("Id") != "cs-1" {
					t.Errorf("unexpected detail Id %q", r.Form.Get("Id"))
				}
				directMailWriteJSON(t, w, map[string]interface{}{
					"Detail": map[string]interface{}{"Id": "cs-1", "Name": "cs", "Description": "from-detail",
						"ValidationOption": map[string]interface{}{"Enabled": true, "ForbiddenStatusList": []interface{}{"INVALID"}}},
				})
			default:
				t.Errorf("unexpected action %s", r.Form.Get("Action"))
			}
		})
		d := schema.TestResourceDataRaw(t, dataSourceAliCloudDirectMailConfigSets().Schema, map[string]interface{}{
			"enable_details": enableDetails,
		})
		if err := dataSourceAliCloudDirectMailConfigSetRead(d, client); err != nil {
			t.Fatal(err)
		}
		set := d.Get("sets").([]interface{})[0].(map[string]interface{})
		wantDescription, wantCalls, wantOptions := "from-list", 0, 0
		if enableDetails {
			wantDescription, wantCalls, wantOptions = "from-detail", 1, 1
		}
		if got := detailCalls.Load(); set["description"] != wantDescription || got != int32(wantCalls) || len(set["validation_option"].([]interface{})) != wantOptions {
			t.Fatalf("enable_details=%t: description=%v detailCalls=%d validation_option=%v", enableDetails, set["description"], got, set["validation_option"])
		}
	}
}

func TestUnitDirectMailConfigSetsDataSourceOutputFile(t *testing.T) {
	client := directMailTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		directMailWriteJSON(t, w, map[string]interface{}{
			"ConfigSets": []interface{}{map[string]interface{}{"Id": "cs-1", "Name": "cs"}},
		})
	})
	output := filepath.Join(t.TempDir(), "sets.json")
	d := schema.TestResourceDataRaw(t, dataSourceAliCloudDirectMailConfigSets().Schema, map[string]interface{}{"output_file": output})
	if err := dataSourceAliCloudDirectMailConfigSetRead(d, client); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("output_file was not written: %s", err)
	}
	if !strings.Contains(string(content), `"cs-1"`) {
		t.Fatalf("unexpected output_file content: %s", content)
	}
	unwritable := schema.TestResourceDataRaw(t, dataSourceAliCloudDirectMailConfigSets().Schema, map[string]interface{}{
		"output_file": filepath.Join(t.TempDir(), "missing-dir", "sets.json"),
	})
	if err := dataSourceAliCloudDirectMailConfigSetRead(unwritable, client); err == nil {
		t.Fatal("an unwritable output_file must return an error")
	}
}

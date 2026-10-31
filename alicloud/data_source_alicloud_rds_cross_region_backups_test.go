package alicloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

/*
	Because there are two lines between enabling remote disaster recovery and generating remote disaster recovery, the generation time of the remote disaster recovery set cannot be determined, the query cannot be determined to have a value, and the OpenAPI will return normally without a value. All users skip the 'Test' and have passed the 'Test' offline simulation.
*/

func SkipTestAccAliCloudRdsCrossRegionBackupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	testConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsCrossRegionBackupsDataSourceName(rand, map[string]string{
			"db_instance_id": `"${alicloud_rds_instance_cross_backup_policy.policy.instance_id}"`,
			"start_time":     `"${var.startTime}"`,
			"end_time":       `"${var.endTime}"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsCrossRegionBackupsDataSourceName(rand, map[string]string{
			"db_instance_id": CHECKSET,
			"start_time":     `"${var.startTime}"`,
			"end_time":       `"${var.errEndTime}"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            `["7190112"]`,
			"db_instance_id": `"${alicloud_rds_instance_cross_backup_policy.policy.instance_id}"`,
			"start_time":     `"${var.startTime}"`,
			"end_time":       `"${var.endTime}"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            CHECKSET,
			"db_instance_id": `"${alicloud_rds_instance_cross_backup_policy.policy.instance_id}"`,
			"start_time":     `"${var.startTime}"`,
			"end_time":       `"${var.endTime}"`,
		}),
	}
	var existAlicloudRdsCrossRegionBackupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "0",
			"backups.#":                            "0",
			"backups.0.cross_backup_download_link": CHECKSET,
			"backups.0.backup_end_time":            CHECKSET,
			"backups.0.id":                         CHECKSET,
			"backups.0.cross_backup_id":            CHECKSET,
			"backups.0.backup_method":              CHECKSET,
			"backups.0.cross_backup_set_size":      CHECKSET,
			"backups.0.backup_start_time":          CHECKSET,
			"backups.0.backup_type":                CHECKSET,
			"backups.0.consistent_time":            CHECKSET,
			"backups.0.instance_id":                CHECKSET,
			"backups.0.backup_set_status":          CHECKSET,
			"backups.0.db_instance_storage_type":   CHECKSET,
			"backups.0.backup_set_scale":           CHECKSET,
			"backups.0.category":                   CHECKSET,
			"backups.0.cross_backup_region":        CHECKSET,
			"backups.0.cross_backup_set_file":      CHECKSET,
			"backups.0.cross_backup_set_location":  CHECKSET,
			"backups.0.engine":                     CHECKSET,
			"backups.0.engine_version":             CHECKSET,
			"backups.0.recovery_begin_time":        CHECKSET,
			"backups.0.recovery_end_time":          CHECKSET,
			"backups.0.restore_regions.#":          "2",
		}
	}
	var fakeAlicloudRdsCrossRegionBackupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"backups.#": "0",
		}
	}
	var alicloudRdsCrossRegionBackupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_rds_cross_region_backups.default",
		existMapFunc: existAlicloudRdsCrossRegionBackupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudRdsCrossRegionBackupsDataSourceNameMapFunc,
	}
	alicloudRdsCrossRegionBackupsCheckInfo.dataSourceTestCheck(t, rand, testConf, idsConf)
}
func testAccCheckAlicloudRdsCrossRegionBackupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	startTime := time.Now().AddDate(0, 0, -2).Format("2006-01-02T15:04Z")
	endTime := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04Z")
	errEndTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04Z")
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-rds-cross-backup"
}

variable "startTime" {
 default = "%v"
}

variable "endTime" {
 default = "%v"
}

variable "errEndTime" {
 default = "%v"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.5.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "local_ssd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_db_zones.default.ids.5
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_db_zones.default.ids.5
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

data "alicloud_rds_cross_regions" "regions" {
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id               = local.vswitch_id
  instance_name            = var.name
}

resource "alicloud_rds_instance_cross_backup_policy" "policy" {
  instance_id         = alicloud_db_instance.default.id
  cross_backup_region = data.alicloud_rds_cross_regions.regions.ids.0
}

data "alicloud_rds_cross_region_backups" "default" {
  %s
}`, startTime, endTime, errEndTime, strings.Join(pairs, "\n"))
}

// rdsCrossRegionBackupsTestClient returns a client whose Rds RPC calls go to a
// local httptest server, so request parameters and responses are observable
// without cloud credentials.
func rdsCrossRegionBackupsTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
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
	endpoints.Store("rds", endpoint)
	return client
}

func writeCrossRegionBackupsJSON(t *testing.T, w http.ResponseWriter, body map[string]interface{}) {
	t.Helper()
	body["RequestId"] = "test-request"
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Error(err)
	}
}

// TestUnitAlicloudRdsCrossRegionBackupsDataSourceResourceGroupId locks the fix
// for the misspelled schema key (resource__group_id): the Read must forward the
// configured resource_group_id to DescribeCrossRegionBackups as
// ResourceGroupId, and must omit it when the field is not set. DBInstanceId is
// always forwarded so the query stays bounded to the configured instance.
func TestUnitAlicloudRdsCrossRegionBackupsDataSourceResourceGroupId(t *testing.T) {
	const (
		dbInstanceID    = "rm-test-instance"
		resourceGroupID = "rg-test-filter"
	)

	assertRequest := func(t *testing.T, raw map[string]interface{}, wantResourceGroupID string) {
		var seenResourceGroupID, seenDBInstanceID string
		var resourceGroupIDPresent bool
		client := rdsCrossRegionBackupsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			switch r.Form.Get("Action") {
			case "DescribeCrossRegionBackups":
				if _, ok := r.Form["ResourceGroupId"]; ok {
					resourceGroupIDPresent = true
					seenResourceGroupID = r.Form.Get("ResourceGroupId")
				} else {
					resourceGroupIDPresent = false
				}
				seenDBInstanceID = r.Form.Get("DBInstanceId")
				writeCrossRegionBackupsJSON(t, w, map[string]interface{}{
					"Items": map[string]interface{}{"Item": []interface{}{}},
				})
			default:
				t.Errorf("unexpected action %s", r.Form.Get("Action"))
			}
		})
		d := schema.TestResourceDataRaw(t, dataSourceAlicloudRdsCrossRegionBackups().Schema, raw)
		if err := dataSourceAlicloudRdsCrossRegionBackupsRead(d, client); err != nil {
			t.Fatalf("read returned error: %s", err)
		}
		if seenDBInstanceID != dbInstanceID {
			t.Errorf("DBInstanceId = %q, want %q", seenDBInstanceID, dbInstanceID)
		}
		if wantResourceGroupID == "" {
			if resourceGroupIDPresent {
				t.Errorf("ResourceGroupId should be omitted, got %q", seenResourceGroupID)
			}
		} else if !resourceGroupIDPresent || seenResourceGroupID != wantResourceGroupID {
			t.Errorf("ResourceGroupId = %q (present=%v), want %q", seenResourceGroupID, resourceGroupIDPresent, wantResourceGroupID)
		}
	}

	t.Run("with resource_group_id", func(t *testing.T) {
		assertRequest(t, map[string]interface{}{
			"db_instance_id":    dbInstanceID,
			"resource_group_id": resourceGroupID,
		}, resourceGroupID)
	})
	t.Run("without resource_group_id", func(t *testing.T) {
		assertRequest(t, map[string]interface{}{
			"db_instance_id": dbInstanceID,
		}, "")
	})
}

package alicloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudMongodbShardingAuditFiltersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbShardingAuditFiltersDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}
	var existAlicloudMongodbShardingAuditFiltersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"filters.#":                "1",
			"filters.0.id":             CHECKSET,
			"filters.0.db_instance_id": CHECKSET,
			"filters.0.audit_status":   "enable",
			"filters.0.filter":         "admin,slow",
			"filters.0.service_type":   "V2_Standard",
			// 180 is the V2 cold storage duration configured below and differs from the Standard-edition
			// default (30), so these assertions fail if the data source reads TtlForStandard instead of
			// TtlForV2Standard. The data source reads every node role at once, so the API reports the
			// merged-view marker `logic`.
			"filters.0.storage_period":     "180",
			"filters.0.hot_storage_period": "7",
			"filters.0.role_type":          "logic",
		}
	}
	var fakeAlicloudMongodbShardingAuditFiltersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"filters.#": "0",
		}
	}
	var alicloudMongodbShardingAuditFiltersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_sharding_audit_filters.default",
		existMapFunc: existAlicloudMongodbShardingAuditFiltersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMongodbShardingAuditFiltersDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMongodbShardingAuditFiltersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckAlicloudMongodbShardingAuditFiltersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccShardingAuditFilter-%d"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id     = data.alicloud_vswitches.default.ids[0]
  engine_version = "4.2"
  name           = var.name
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
}

resource "alicloud_mongodb_sharding_audit_filter" "default" {
  db_instance_id     = alicloud_mongodb_sharding_instance.default.id
  audit_status       = "enable"
  filter             = "admin,slow"
  service_type       = "V2_Standard"
  storage_period     = 180
  hot_storage_period = 7
}

data "alicloud_mongodb_sharding_audit_filters" "default" {
  db_instance_id = alicloud_mongodb_sharding_audit_filter.default.db_instance_id
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

// TestUnitAliCloudMongodbShardingAuditFiltersDataSourceRetention pins the edition-specific retention
// fields the data source exports. Each fixture keeps the two editions on different values, so reading
// TtlForStandard for a V2_Standard audit log fails the assertion instead of coincidentally passing.
func TestUnitAliCloudMongodbShardingAuditFiltersDataSourceRetention(t *testing.T) {
	cases := []struct {
		name          string
		logConfig     string
		storagePeriod int
		hotPeriod     int
	}{
		{
			name:          "v2 standard exports the cold storage duration",
			logConfig:     `{"RequestId":"mock","ServiceType":"V2_Standard","TtlForStandard":30,"TtlForV2Standard":180,"HotTtlForV2Standard":5}`,
			storagePeriod: 180,
			hotPeriod:     5,
		},
		{
			name:          "standard exports its retention and no hot storage",
			logConfig:     `{"RequestId":"mock","ServiceType":"Standard","TtlForStandard":30,"TtlForV2Standard":180,"HotTtlForV2Standard":5}`,
			storagePeriod: 30,
			hotPeriod:     0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request form: %v", err)
				}
				switch r.Form.Get("Action") {
				case "DescribeMongoDBLogConfig":
					fmt.Fprint(w, c.logConfig)
				case "DescribeAuditPolicy":
					fmt.Fprint(w, `{"RequestId":"mock","LogAuditStatus":"Enable"}`)
				case "DescribeAuditLogFilter":
					fmt.Fprint(w, `{"RequestId":"mock","Filter":"mongos@admin,slow-db@admin,slow","RoleType":"logic"}`)
				default:
					fmt.Fprint(w, `{"RequestId":"mock"}`)
				}
			}))
			t.Cleanup(server.Close)

			endpoint := strings.TrimPrefix(server.URL, "http://")
			t.Setenv("NO_PROXY", endpoint)
			credential, err := credentials.NewCredential(new(credentials.Config).
				SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
			if err != nil {
				t.Fatal(err)
			}
			endpoints := new(sync.Map)
			config := &connectivity.Config{
				AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
				RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
				Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
			}
			client, err := config.Client()
			if err != nil {
				t.Fatal(err)
			}
			endpoints.Store("dds", endpoint)

			data := schema.TestResourceDataRaw(t, dataSourceAlicloudMongodbShardingAuditFilters().Schema, map[string]interface{}{
				"db_instance_id": "dds-mock",
			})
			assert.NoError(t, dataSourceAlicloudMongodbShardingAuditFiltersRead(data, client))
			assert.Equal(t, 1, data.Get("filters.#"))
			assert.Equal(t, c.storagePeriod, data.Get("filters.0.storage_period"))
			assert.Equal(t, c.hotPeriod, data.Get("filters.0.hot_storage_period"))
			assert.Equal(t, "logic", data.Get("filters.0.role_type"))
			assert.Equal(t, "admin,slow", data.Get("filters.0.filter"))
		})
	}
}

package ddosbgp_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	ddosbgp "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ddosbgp"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDdosbgpInstanceDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_ddosbgp_instances.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceDdosbgpInstanceConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddosbgp_instance.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddosbgp_instance.default.name}-fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ddosbgp_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ddosbgp_instance.default.id}-fake"},
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddosbgp_instance.default.name}",
			"ids":        []string{"${alicloud_ddosbgp_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddosbgp_instance.default.name}-fake",
			"ids":        []string{"${alicloud_ddosbgp_instance.default.id}"},
		}),
	}

	var existDdosbgpInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        acctest.CHECKSET,
			"names.#":                      "1",
			"names.0":                      fmt.Sprintf("tf_testAcc%d", rand),
			"instances.#":                  "1",
			"instances.0.name":             fmt.Sprintf("tf_testAcc%d", rand),
			"instances.0.type":             string(ddosbgp.Enterprise),
			"instances.0.base_bandwidth":   "20",
			"instances.0.bandwidth":        "-1",
			"instances.0.ip_count":         "100",
			"instances.0.ip_type":          "IPv4",
			"instances.0.region":           "cn-hangzhou",
			"instances.0.normal_bandwidth": "100",
		}
	}

	var fakeDdosbgpInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var ddosbgpInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDdosbgpInstanceMapFunc,
		FakeMapFunc:  fakeDdosbgpInstanceMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithTime(t, []int{15})
		acctest.PreCheckWithRegions(t, true, connectivity.DdosbgpSupportedRegions)
	}
	ddosbgpInstanceCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func dataSourceDdosbgpInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`

resource "alicloud_ddosbgp_instance" "default" {
      name                    = "%s"
      bandwidth               = -1
      base_bandwidth          = 20
      normal_bandwidth        = 100
      ip_count                = 100
      ip_type                 = "IPv4"
      type                    = "Enterprise"
	}`, name)
}

// TestAccAliCloudDdosbgpInstanceDataSource_readOnly validates the instances data
// source read-only query (DescribeInstanceList + DescribeInstanceSpecs) and the
// name_regex filter against pre-existing ddosbgp instances owned by the test
// account, without creating any prepaid fixture. Unlike the legacy _basic case
// it is not gated by PreCheckWithTime([]int{15}); that day-15 gate only bounds
// the prepaid-fixture variant (which creates alicloud_ddosbgp_instance to keep
// subscription cost off non-scheduled days) and does not reflect a limitation of
// the read-only Describe API. PreCheckWithRegions(true, DdosbgpSupportedRegions)
// lands the test in cn-hangzhou, where the account already owns reusable prepaid
// ddosbgp instances; no new prepaid instance is subscribed and no balance-charged
// request is triggered. Coverage: DS list query + name_regex filter (exist/fake).
// Not covered here: ids filter and per-instance spec mapping for Enterprise
// instances (the account only has Professional instances); those remain covered
// by the legacy _basic case on its scheduled day.
// lintignore: AT001
func TestAccAliCloudDdosbgpInstanceDataSource_readOnly(t *testing.T) {
	resourceId := "data.alicloud_ddosbgp_instances.default"

	// configDependence returns no resources: pure read-only query against the
	// account's existing ddosbgp instances.
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "readonly", func(string) string { return "" })

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "tfacc.*",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "tfacc.*-nonexistent-zzz",
		}),
	}

	// exist: the account owns reusable ddosbgp instances whose Remark starts with
	// "tfacc", so name_regex "tfacc.*" must match a non-empty set. ids/names/
	// instances counts are asserted non-zero via regex match (account instance
	// count may vary, so an exact value is not asserted).
	var existReadOnlyMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       acctest.REGEXMATCH + "^[1-9][0-9]*$",
			"names.#":     acctest.REGEXMATCH + "^[1-9][0-9]*$",
			"instances.#": acctest.REGEXMATCH + "^[1-9][0-9]*$",
		}
	}
	// fake: a non-matching name_regex must yield an empty result set.
	var fakeReadOnlyMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}
	checkInfo := acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existReadOnlyMapFunc,
		FakeMapFunc:  fakeReadOnlyMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.DdosbgpSupportedRegions)
	}
	checkInfo.DataSourceTestCheckWithPreCheck(t, tfacctest.RandInt(), preCheck, nameRegexConf)
}

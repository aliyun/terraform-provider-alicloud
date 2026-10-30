package rds_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRdsCharacterSetNamesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_rds_character_set_names.default"
	rdsCharacterSetNamesConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: dataSourceRdsCharacterSetNamesConfigDependence(rand, map[string]string{
			"engine": `"MySQL"`,
		}),
		FakeConfig: dataSourceRdsCharacterSetNamesConfigDependence(rand, map[string]string{
			"engine": `"MySQL"`,
		}),
	}

	var existAlicloudRdsCharacterSetNamesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#": acctest.CHECKSET,
			"names.0": acctest.CHECKSET,
		}
	}

	var fakeAlicloudRdsCharacterSetNamesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#": acctest.CHECKSET,
		}
	}

	var alicloudRdsCharacterSetNamesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existAlicloudRdsCharacterSetNamesMapFunc,
		FakeMapFunc:  fakeAlicloudRdsCharacterSetNamesMapFunc,
	}

	alicloudRdsCharacterSetNamesCheckInfo.DataSourceTestCheck(t, rand, rdsCharacterSetNamesConfig)
}

func dataSourceRdsCharacterSetNamesConfigDependence(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
data "alicloud_rds_character_set_names" "default" {	
  %s
}`, strings.Join(pairs, "\n"))
}

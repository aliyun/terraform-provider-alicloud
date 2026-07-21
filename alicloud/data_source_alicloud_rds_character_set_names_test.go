package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsCharacterSetNamesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_rds_character_set_names.default"
	rdsCharacterSetNamesConfig := dataSourceTestAccConfig{
		existConfig: dataSourceRdsCharacterSetNamesConfigDependence(rand, map[string]string{
			"engine": `"MySQL"`,
		}),
		fakeConfig: dataSourceRdsCharacterSetNamesConfigDependence(rand, map[string]string{
			"engine": `"MySQL"`,
		}),
	}

	var existAlicloudRdsCharacterSetNamesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#": CHECKSET,
			"names.0": CHECKSET,
		}
	}

	var fakeAlicloudRdsCharacterSetNamesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#": CHECKSET,
		}
	}

	var alicloudRdsCharacterSetNamesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAlicloudRdsCharacterSetNamesMapFunc,
		fakeMapFunc:  fakeAlicloudRdsCharacterSetNamesMapFunc,
	}

	alicloudRdsCharacterSetNamesCheckInfo.dataSourceTestCheck(t, rand, rdsCharacterSetNamesConfig)
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

package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDataWorksExtensionsDataSource_basic(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_data_works_extensions.default"
	name := fmt.Sprintf("tf-testacc-dataworksextension%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataWorksExtensionsConfigDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"output_file": "data_works_extensions.json",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "___non_existent_extension_regex___12345",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"extensions.#": CHECKSET,
			"ids.#":        CHECKSET,
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"extensions.#": "0",
			"ids.#":        "0",
		}
	}

	checkInfo := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	checkInfo.dataSourceTestCheck(t, rand, basicConf)
}

func dataSourceDataWorksExtensionsConfigDependence(name string) string {
	return ""
}

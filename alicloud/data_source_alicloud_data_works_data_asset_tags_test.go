package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudDataWorksDataAssetTagsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_data_works_data_asset_tags.default"
	name := fmt.Sprintf("tf-testacc-dataworksdataassettag%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataWorksDataAssetTagsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"key": "${alicloud_data_works_data_asset_tag.default.key}",
			"ids": []string{"${alicloud_data_works_data_asset_tag.default.key}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"key": "${alicloud_data_works_data_asset_tag.default.key}_fake",
			"ids": []string{"${alicloud_data_works_data_asset_tag.default.key}_fake"},
		}),
	}

	var existDataWorksDataAssetTagsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "1",
			"tags.#":     "1",
			"tags.0.key": CHECKSET,
		}
	}
	var fakeDataWorksDataAssetTagsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"tags.#": "0",
			"ids.#":  "0",
		}
	}

	var DataWorksDataAssetTagsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDataWorksDataAssetTagsMapFunc,
		fakeMapFunc:  fakeDataWorksDataAssetTagsMapFunc,
	}

	DataWorksDataAssetTagsCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceDataWorksDataAssetTagsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_data_works_data_asset_tag" "default" {
	key         = var.name
	value_type  = "String"
	description = "tf_desc"
}
`, name)
}

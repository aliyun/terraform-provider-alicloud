package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// TestAccAlicloudDataWorksMetaCategoriesDataSource verifies the alicloud_data_works_meta_categories
// data source lists the category created by the companion resource under the same parent.
func TestAccAlicloudDataWorksMetaCategoriesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_data_works_meta_categories.default"
	name := fmt.Sprintf("tf-testacc-meta-cat-ds-%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataWorksMetaCategoriesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"parent_category_id": "${split(\":\", alicloud_data_works_meta_category.parent.id)[0]}",
			"ids":                []string{"${alicloud_data_works_meta_category.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"parent_category_id": "${split(\":\", alicloud_data_works_meta_category.parent.id)[0]}",
			"ids":                []string{"${alicloud_data_works_meta_category.default.id}_fake"},
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"categories.#":                    "1",
			"categories.0.category_id":        CHECKSET,
			"categories.0.name":               CHECKSET,
			"categories.0.parent_category_id": CHECKSET,
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"categories.#": "0",
			"ids.#":        "0",
		}
	}

	var checkInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	checkInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceDataWorksMetaCategoriesConfigDependence(name string) string {
	// Randomize the parent category name so each run creates a distinct parent
	// instead of colliding with a residual same-named category (e.g. a category
	// protected by a foreign-key table that the API refuses to delete with
	// Invalid.Meta.CategoryForbidden 403, which would make every rerun fail at
	// create with Invalid.Meta.CategoryNameDuplicate).
	parentName := fmt.Sprintf("tf-meta-cat-parent-%d", acctest.RandInt())
	return fmt.Sprintf(`
variable "name" {
    default = "%v"
}

resource "alicloud_data_works_meta_category" "parent" {
  name              = "%v"
  comment           = "parent category"
  parent_category_id = 0
}

resource "alicloud_data_works_meta_category" "default" {
  name              = var.name
  comment           = "child category"
  parent_category_id = alicloud_data_works_meta_category.parent.category_id
}
`, name, parentName)
}

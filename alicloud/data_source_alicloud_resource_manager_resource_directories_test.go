package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerResourceDirectoriesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(1000000, 9999999)

	conf := dataSourceTestAccConfig{
		existConfig: fmt.Sprintf(`data "alicloud_resource_manager_resource_directories" "default"{}`),
		fakeConfig:  "",
	}

	var existResourceManagerResourceDirectoriesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"directories.#":                       "1",
			"directories.0.id":                    CHECKSET,
			"directories.0.resource_directory_id": CHECKSET,
			"directories.0.master_account_id":     CHECKSET,
			"directories.0.master_account_name":   CHECKSET,
			"directories.0.root_folder_id":        CHECKSET,
			"directories.0.status":                CHECKSET,
		}
	}

	var fakeResourceManagerResourceDirectoriesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"directories.#": "0",
		}
	}

	var ResourceDirectoriesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_resource_directories.default",
		existMapFunc: existResourceManagerResourceDirectoriesRecordsMapFunc,
		fakeMapFunc:  fakeResourceManagerResourceDirectoriesRecordsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckEnterpriseAccountEnabled(t)
	}

	ResourceDirectoriesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, conf)

}

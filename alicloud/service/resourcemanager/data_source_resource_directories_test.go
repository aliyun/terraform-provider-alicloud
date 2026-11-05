package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudResourceManagerResourceDirectoriesDataSource(t *testing.T) {

	rand := tfacctest.RandIntRange(1000000, 9999999)

	conf := acctest.DataSourceTestAccConfig{
		ExistConfig: fmt.Sprintf(`data "alicloud_resource_manager_resource_directories" "default"{}`),
		FakeConfig:  "",
	}

	var existResourceManagerResourceDirectoriesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"directories.#":                       "1",
			"directories.0.id":                    acctest.CHECKSET,
			"directories.0.resource_directory_id": acctest.CHECKSET,
			"directories.0.master_account_id":     acctest.CHECKSET,
			"directories.0.master_account_name":   acctest.CHECKSET,
			"directories.0.root_folder_id":        acctest.CHECKSET,
			"directories.0.status":                acctest.CHECKSET,
		}
	}

	var fakeResourceManagerResourceDirectoriesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"directories.#": "0",
		}
	}

	var ResourceDirectoriesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_resource_directories.default",
		ExistMapFunc: existResourceManagerResourceDirectoriesRecordsMapFunc,
		FakeMapFunc:  fakeResourceManagerResourceDirectoriesRecordsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckEnterpriseAccountEnabled(t)
	}

	ResourceDirectoriesRecordsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, conf)

}

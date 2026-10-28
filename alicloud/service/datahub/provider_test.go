package datahub_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	datahubsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/datahub"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// providerCommon mirrors the root alicloud/provider_test.go provider block
// that drives an independent assume_role provider scenario. It is kept local
// to this split-package test so the datahub tests no longer depend on the
// unexported root helper, while preserving the original configuration exactly.
const providerCommon = `
provider "alicloud" {
	assume_role {}
}

`

// TestAccAliCloudProviderDatahub restores the assume_role provider scenario
// deleted from alicloud/provider_test.go during the datahub split-package
// migration. Unlike Project_basic it exercises the assume_role provider
// configuration, so it cannot be treated as a duplicate of the basic case.
func TestAccAliCloudProviderDatahub(t *testing.T) {
	var v *datahub.GetProjectResult

	resourceId := "alicloud_datahub_project.default"
	ra := acctest.ResourceAttrInit(resourceId, datahubProjectBasicMap)

	serviceFunc := func() interface{} {
		return datahubsvc.NewDatahubService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testaccdatahubproject%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, func(name string) string {
		return providerCommon + resourceDatahubProjectConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
		},
	})
}

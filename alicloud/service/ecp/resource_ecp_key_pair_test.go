package ecp_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecp"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEcpKeyPair_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecp_key_pair.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEcpKeyPairMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcpKeyPair")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcpKeyPair%d", rand)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEcpKeyPairBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcpSupportRegions)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name":   name,
					"public_key_body": "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"public_key_body"},
			},
		},
	})
}

var AlicloudEcpKeyPairMap = map[string]string{}

func AlicloudEcpKeyPairBasicDependence(name string) string {
	return ""
}

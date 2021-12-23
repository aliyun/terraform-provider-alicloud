package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEcpKeyPair_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecp_key_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudEcpKeyPairMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpKeyPair")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcpKeyPair%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcpKeyPairBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EcpSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

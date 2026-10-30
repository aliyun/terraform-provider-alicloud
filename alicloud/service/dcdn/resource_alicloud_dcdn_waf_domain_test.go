package dcdn_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dcdn"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudDcdnWafDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_waf_domain.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDCDNWafDomainMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnWafDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDCDNWafDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_dcdn_domain.default.domain_name}",
					"client_ip_tag": "X-Forwarded-For",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":   acctest.CHECKSET,
						"client_ip_tag": "X-Forwarded-For",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudDCDNWafDomainMap0 = map[string]string{
	"domain_name": acctest.CHECKSET,
}

func AlicloudDCDNWafDomainBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "domain_name" {
  default = "%s"
}
resource "alicloud_dcdn_domain" "default" {
  domain_name = "${var.domain_name}"
  sources {
    content = "1.1.1.1"
    port = "80"
    priority = "20"
    type = "ipaddr"
  }
}
`, name)
}

func TestAccAliCloudDcdnWafDomain_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_waf_domain.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDCDNWafDomainMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnWafDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDCDNWafDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${alicloud_dcdn_domain.default.domain_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ip_tag": "X-Forwarded-For",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ip_tag": "X-Forwarded-For",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ip_tag": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ip_tag": acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

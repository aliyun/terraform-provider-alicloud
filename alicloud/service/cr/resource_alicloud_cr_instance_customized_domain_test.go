// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	crsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Cr InstanceCustomizedDomain. >>> Resource test cases, automatically generated.
// Case InstanceCustomizedDomain 11888
func TestAccAliCloudCrInstanceCustomizedDomain_basic11888(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_instance_customized_domain.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCrInstanceCustomizedDomainMap11888)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrInstanceCustomizedDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCrInstanceCustomizedDomainBasicDependence11888)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
					"module_name": "Registry",
					"domain":      "alicloud-provider.cn",
					"cert_id":     "${data.alicloud_ssl_certificates_service_certificates.default.certificates.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
						"module_name": "Registry",
						"domain":      "alicloud-provider.cn",
						"cert_id":     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cert_id":        "${data.alicloud_ssl_certificates_service_certificates.default.certificates.1.id}",
					"cert_region_id": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert_region_id"},
			},
		},
	})
}

func TestAccAliCloudCrInstanceCustomizedDomain_basic11888_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_instance_customized_domain.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCrInstanceCustomizedDomainMap11888)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrInstanceCustomizedDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCrInstanceCustomizedDomainBasicDependence11888)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_cr_ee_instance.default.id}",
					"module_name":    "Registry",
					"domain":         "alicloud-provider.cn",
					"cert_id":        "${data.alicloud_ssl_certificates_service_certificates.default.certificates.0.id}",
					"cert_region_id": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
						"module_name": "Registry",
						"domain":      "alicloud-provider.cn",
						"cert_id":     acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert_region_id"},
			},
		},
	})
}

var AliCloudCrInstanceCustomizedDomainMap11888 = map[string]string{
	"region_id":     acctest.CHECKSET,
	"create_time":   acctest.CHECKSET,
	"modified_time": acctest.CHECKSET,
}

func AliCloudCrInstanceCustomizedDomainBasicDependence11888(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_ssl_certificates_service_certificates" "default" {
  keyword = "alicloud-provider.cn"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 1
  renewal_status = "AutoRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}
`, name)
}

// Test Cr InstanceCustomizedDomain. <<< Resource test cases, automatically generated.

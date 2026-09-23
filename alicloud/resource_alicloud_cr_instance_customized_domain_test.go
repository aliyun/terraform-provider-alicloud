// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cr InstanceCustomizedDomain. >>> Resource test cases, automatically generated.
// Case InstanceCustomizedDomain 11888
func TestAccAliCloudCrInstanceCustomizedDomain_basic11888(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_instance_customized_domain.default"
	ra := resourceAttrInit(resourceId, AliCloudCrInstanceCustomizedDomainMap11888)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrInstanceCustomizedDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCrInstanceCustomizedDomainBasicDependence11888)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"instance_id": CHECKSET,
						"module_name": "Registry",
						"domain":      "alicloud-provider.cn",
						"cert_id":     CHECKSET,
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
						"cert_id": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCrInstanceCustomizedDomainMap11888)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrInstanceCustomizedDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCrInstanceCustomizedDomainBasicDependence11888)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"instance_id": CHECKSET,
						"module_name": "Registry",
						"domain":      "alicloud-provider.cn",
						"cert_id":     CHECKSET,
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
	"region_id":     CHECKSET,
	"create_time":   CHECKSET,
	"modified_time": CHECKSET,
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

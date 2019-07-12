package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCdnDomainConfig_ip_allow_list(t *testing.T) {
	var v *cdn.DomainConfig

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "ip_allow_list_set",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "ip_list",
							"arg_value": "110.110.110.110",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":               name,
						"function_name":             "ip_allow_list_set",
						"function_args.0.arg_name":  "ip_list",
						"function_args.0.arg_value": "110.110.110.110",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_name"},
			},
		},
	})
}

func TestAccAlicloudCdnDomainConfig_referer_white_list(t *testing.T) {
	var v *cdn.DomainConfig

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "referer_white_list_set",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "refer_domain_allow_list",
							"arg_value": "110.110.110.110",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":               name,
						"function_name":             "referer_white_list_set",
						"function_args.0.arg_name":  "refer_domain_allow_list",
						"function_args.0.arg_value": "110.110.110.110",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_name"},
			},
		},
	})
}

func resourceCdnDomainConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "default" {
	  domain_name = "%s"
	  cdn_type = "web"
      scope = "overseas"
      sources {
         content = "www.aliyuntest.com"
         type = "domain"
         priority = 20
         port = 80
         weight = 10
      }
	}
`, name)
}

func resourceCdnDomainConfigDependence_oss(name string) string {
	return fmt.Sprintf(`
	
	resource "alicloud_cdn_domain_new" "default" {
	  domain_name = "tf-testacc%s-oss.xiaozhu.com"
	  cdn_type = "web"
      scope = "overseas"
      sources {
         content = "www.aliyuntest.com"
         type = "oss"
         priority = 20
         port = 80
         weight = 10
      }
	}

	resource "alicloud_oss_bucket" "default" {
	  bucket = "tf-test-domain-config-%s"
	}
`, name, name)
}

var cdnDomainConfigBasicMap = map[string]string{
	"domain_name":               CHECKSET,
	"function_name":             CHECKSET,
	"function_args.0.arg_name":  CHECKSET,
	"function_args.0.arg_value": CHECKSET,
}

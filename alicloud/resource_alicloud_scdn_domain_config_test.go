package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudScdnDomainConfig_ip_allow_list(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_scdn_domain_config.default"
	ra := resourceAttrInit(resourceId, scdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &ScdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceScdnDomainConfigDependence)
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
					"domain_name":   "${alicloud_scdn_domain.default.domain_name}",
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
						"domain_name":     name,
						"function_name":   "ip_allow_list_set",
						"function_args.#": "1",
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

func TestAccAlicloudScdnDomainConfig_referer_white_list(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_scdn_domain_config.default"
	ra := resourceAttrInit(resourceId, scdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &ScdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceScdnDomainConfigDependence)
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
					"domain_name":   "${alicloud_scdn_domain.default.domain_name}",
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
						"domain_name":     name,
						"function_name":   "referer_white_list_set",
						"function_args.#": "1",
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

func TestAccAlicloudScdnDomainConfig_filetype_based_ttl_set(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_scdn_domain_config.default"
	ra := resourceAttrInit(resourceId, scdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &ScdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceScdnDomainConfigDependence)

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
					"domain_name":   "${alicloud_scdn_domain.default.domain_name}",
					"function_name": "filetype_based_ttl_set",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "ttl",
							"arg_value": "300",
						},
						{
							"arg_name":  "file_type",
							"arg_value": "jpg",
						},
						{
							"arg_name":  "weight",
							"arg_value": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "filetype_based_ttl_set",
						"function_args.#": "3",
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

func resourceScdnDomainConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	resource "alicloud_scdn_domain" "default" {
	 domain_name = var.name
	 status = "online"
     sources {
        content = "1.1.1.1"
        type = "ipaddr"
        priority = 20
        port = 80
     }
	}
`, name)
}

var scdnDomainConfigBasicMap = map[string]string{
	"config_id": CHECKSET,
	"status":    "success",
}

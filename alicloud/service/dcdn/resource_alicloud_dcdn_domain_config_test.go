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

func TestAccAliCloudDcdnDomainConfig_ip_allow_list(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_domain_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDCDNDomainConfigMap0)
	serviceFunc := func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDCDNDomainConfigBasicDependence0)
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
					"function_name": "ip_allow_list_set",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "ip_list",
							"arg_value": "192.168.0.1",
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
				Config: testAccConfig(map[string]interface{}{
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "ip_list",
							"arg_value": "192.168.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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

func TestAccAliCloudDcdnDomainConfig_referer_white_list(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_domain_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDCDNDomainConfigMap0)
	serviceFunc := func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDCDNDomainConfigBasicDependence0)
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
					"function_name": "referer_white_list_set",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "refer_domain_allow_list",
							"arg_value": "192.168.0.1",
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
				Config: testAccConfig(map[string]interface{}{
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "refer_domain_allow_list",
							"arg_value": "192.168.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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

func TestAccAliCloudDcdnDomainConfig_filetype_based_ttl_set(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_domain_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDCDNDomainConfigMap0)
	serviceFunc := func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDCDNDomainConfigBasicDependence1)
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
				Config: testAccConfig(map[string]interface{}{
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "ttl",
							"arg_value": "500000",
						},
						{
							"arg_name":  "file_type",
							"arg_value": "txt",
						},
						{
							"arg_name":  "weight",
							"arg_value": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_args.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_id": "${alicloud_dcdn_domain_config.parent.config_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_id": acctest.CHECKSET,
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

func TestAccAliCloudDcdnDomainConfig_filetype_based_ttl_set_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_domain_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDCDNDomainConfigMap0)
	serviceFunc := func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDCDNDomainConfigBasicDependence1)
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
					"function_name": "filetype_based_ttl_set",
					"parent_id":     "${alicloud_dcdn_domain_config.parent.config_id}",
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
						"parent_id":       acctest.CHECKSET,
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

var AliCloudDCDNDomainConfigMap0 = map[string]string{
	"config_id": acctest.CHECKSET,
	"status":    acctest.CHECKSET,
}

func AliCloudDCDNDomainConfigBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_dcdn_domain" "default" {
  		domain_name = var.name
  		scope       = "overseas"
  		status      = "online"
  		sources {
    		content  = "1.1.1.1"
    		type     = "ipaddr"
    		priority = 20
    		port     = 80
    		weight   = 10
  		}
	}
`, name)
}

func AliCloudDCDNDomainConfigBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_dcdn_domain" "default" {
  		domain_name = var.name
  		scope       = "overseas"
  		status      = "online"
  		sources {
    		content  = "1.1.1.1"
    		type     = "ipaddr"
    		priority = 20
    		port     = 80
    		weight   = 10
  		}
	}

	resource "alicloud_dcdn_domain_config" "parent" {
  		domain_name   = alicloud_dcdn_domain.default.domain_name
  		function_name = "condition"
  		function_args {
    		arg_name  = "rule"
    		arg_value = "{\"match\":{\"logic\":\"and\",\"criteria\":[{\"matchType\":\"clientipVer\",\"matchObject\":\"CONNECTING_IP\",\"matchOperator\":\"equals\",\"matchValue\":\"v6\",\"negate\":false}]},\"name\":\"example\",\"status\":\"enable\"}"
  		}
	}
`, name)
}

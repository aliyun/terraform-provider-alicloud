package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudCDNDomainConfig_filetype_based_ttl_set_new(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
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
							"arg_value": "200",
						},
						{
							"arg_name":  "file_type",
							"arg_value": "png",
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
					"parent_id": "${alicloud_cdn_domain_config.parent.config_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_id":       CHECKSET,
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

func TestAccAliCloudCDNDomainConfig_filetype_based_ttl_set_new_twin(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "filetype_based_ttl_set",
					"parent_id":     "${alicloud_cdn_domain_config.parent.config_id}",
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
						"parent_id":       CHECKSET,
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

func TestAccAliCloudCDNDomainConfig_ip_allow_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

func TestAccAliCloudCDNDomainConfig_referer_white_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

func TestAccAliCloudCDNDomainConfig_referer_black_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "referer_black_list_set",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "refer_domain_deny_list",
							"arg_value": "110.110.110.110",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "referer_black_list_set",
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

func TestAccAliCloudCDNDomainConfig_filetype_based_ttl_set(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
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

func TestAccAliCloudCDNDomainConfig_oss_auth(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := defaultRegionToTest + strconv.Itoa(rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence2)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "oss_auth",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "oss_bucket_id",
							"arg_value": "${data.alicloud_oss_buckets.default.buckets.0.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(

					testAccCheck(map[string]string{
						"domain_name":     fmt.Sprintf("tf-testacc%s-oss.alicloud-provider.cn", name),
						"function_name":   "oss_auth",
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

func TestAccAliCloudCDNDomainConfig_ip_black_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "ip_black_list_set",
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
						"function_name":   "ip_black_list_set",
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

func TestAccAliCloudCDNDomainConfig_ip_white_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "ip_white_list_set",
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
						"function_name":   "ip_white_list_set",
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

func TestAccAliCloudCDNDomainConfig_error_page(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "error_page",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "error_code",
							"arg_value": "502",
						},
						{
							"arg_name":  "rewrite_page",
							"arg_value": "http://www.alicloud-provider.cn",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "error_page",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_set_req_host_header(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "set_req_host_header",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "domain_name",
							"arg_value": "alicloud-provider.cn",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "set_req_host_header",
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

func TestAccAliCloudCDNDomainConfig_set_hashkey_args(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "set_hashkey_args",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "disable",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "set_hashkey_args",
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

func TestAccAliCloudCDNDomainConfig_aliauth(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "aliauth",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "auth_type",
							"arg_value": "no_auth",
						},
						{
							"arg_name":  "ali_auth_delta",
							"arg_value": "1800",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "aliauth",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_set_resp_header(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "set_resp_header",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "key",
							"arg_value": "expires",
						},
						{
							"arg_name":  "value",
							"arg_value": "null",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "set_resp_header",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_https_force(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "https_force",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enable",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "https_force",
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

func TestAccAliCloudCDNDomainConfig_http_force(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "http_force",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enable",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "http_force",
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

func TestAccAliCloudCDNDomainConfig_https_option(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "https_option",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "http2",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "https_option",
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

func TestAccAliCloudCDNDomainConfig_l2_oss_key(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "l2_oss_key",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "private_oss_auth",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "l2_oss_key",
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

func TestAccAliCloudCDNDomainConfig_forward_scheme(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "forward_scheme",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enable",
							"arg_value": "on",
						},
						{
							"arg_name":  "scheme_origin",
							"arg_value": "https",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "forward_scheme",
						"function_args.#": "2",
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

func SkipTestAccAliCloudCDNDomainConfig_green_manager(t *testing.T) {
	// the function: green_manager has been deleted
	t.Skip()
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "green_manager",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enable",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "green_manager",
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

func TestAccAliCloudCDNDomainConfig_tmd_signature(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "tmd_signature",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "ttl",
							"arg_value": "10",
						},
						{
							"arg_name":  "name",
							"arg_value": "tmd_signature_test",
						},
						{
							"arg_name":  "path",
							"arg_value": "/tftest",
						},
						{
							"arg_name":  "pathType",
							"arg_value": "0",
						},
						{
							"arg_name":  "interval",
							"arg_value": "20",
						},
						{
							"arg_name":  "count",
							"arg_value": "500",
						},
						{
							"arg_name":  "action",
							"arg_value": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "tmd_signature",
						"function_args.#": "7",
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

// Skip test due to product function conflict.
//func TestAccAliCloudCDNDomainConfig_dynamic(t *testing.T) {
//	var v *cdn.DomainConfigInDescribeCdnDomainConfigs
//
//	resourceId := "alicloud_cdn_domain_config.default"
//	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)
//
//	serviceFunc := func() interface{} {
//		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
//	}
//	rc := resourceCheckInit(resourceId, &v, serviceFunc)
//
//	rac := resourceAttrCheckInit(rc, ra)
//
//	testAccCheck := rac.resourceAttrMapUpdateSet()
//	rand := acctest.RandIntRange(1000000, 9999999)
//	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
//		"arg_name":  "enable",
//		"arg_value": "on",
//	}))
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//		},
//		// module name
//		IDRefreshName: resourceId,
//		ProviderFactories: testAccProviderFactory,
//		CheckDestroy:  rac.checkResourceDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
//					"function_name": "dynamic",
//					"function_args": []map[string]interface{}{
//						{
//							"arg_name":  "enable",
//							"arg_value": "on",
//						},
//					},
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"domain_name":   name,
//						"function_name": "dynamic",

//					}),
//				),
//			},
//			{
//				ResourceName:            resourceId,
//				ImportState:             true,
//				ImportStateVerify:       true,
//			},
//		},
//	})
//}

func TestAccAliCloudCDNDomainConfig_set_req_header(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "set_req_header",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "key",
							"arg_value": "tftest",
						},
						{
							"arg_name":  "value",
							"arg_value": "tftest",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "set_req_header",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_range(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "range",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enable",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "range",
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

func TestAccAliCloudCDNDomainConfig_video_seek(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "video_seek",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enable",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "video_seek",
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

func TestAccAliCloudCDNDomainConfig_https_tls_version(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "https_tls_version",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "tls10",
							"arg_value": "on",
						},
						{
							"arg_name":  "tls13",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "https_tls_version",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_HSTS(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "HSTS",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enabled",
							"arg_value": "on",
						},
						{
							"arg_name":  "https_hsts_max_age",
							"arg_value": "5184000",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "HSTS",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_filetype_force_ttl_code(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "filetype_force_ttl_code",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "file_type",
							"arg_value": "jpg",
						},
						{
							"arg_name":  "code_string",
							"arg_value": "302=0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "filetype_force_ttl_code",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_path_force_ttl_code(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "path_force_ttl_code",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "code_string",
							"arg_value": "302=0",
						},
						{
							"arg_name":  "path",
							"arg_value": "/tftest",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "path_force_ttl_code",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_gzip(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "gzip",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enable",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "gzip",
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

func TestAccAliCloudCDNDomainConfig_tesla(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "tesla",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enable",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "tesla",
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

func TestAccAliCloudCDNDomainConfig_https_origin_sni(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "https_origin_sni",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enabled",
							"arg_value": "on",
						},
						{
							"arg_name":  "https_origin_sni",
							"arg_value": "alicloud-provider.cn",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "https_origin_sni",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_brotli(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "brotli",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "brotli_level",
							"arg_value": "1",
						},
						{
							"arg_name":  "enable",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "brotli",
						"function_args.#": "2",
					}),
				),
			},
			/*{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
			},*/
		},
	})
}

func TestAccAliCloudCDNDomainConfig_ali_ua(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "ali_ua",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "ua",
							"arg_value": "User-Agent",
						},
						{
							"arg_name":  "type",
							"arg_value": "black",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "ali_ua",
						"function_args.#": "2",
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

func TestAccAliCloudCDNDomainConfig_host_redirect(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, AliCloudCDNDomainConfigMap0)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCDNDomainConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":   "${alicloud_cdn_domain_new.default.domain_name}",
					"function_name": "host_redirect",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "regex",
							"arg_value": "/$",
						},
						{
							"arg_name":  "replacement",
							"arg_value": "/go/act/sale/tbzlsy.php",
						},
						{
							"arg_name":  "flag",
							"arg_value": "redirect",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":     name,
						"function_name":   "host_redirect",
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

var AliCloudCDNDomainConfigMap0 = map[string]string{
	"parent_id": CHECKSET,
	"config_id": CHECKSET,
	"status":    CHECKSET,
}

func AliCloudCDNDomainConfigBasicDependence0(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "default" {
  		domain_name = "%s"
  		cdn_type    = "web"
  		scope       = "overseas"
  		sources {
    		content  = "www.aliyuntest.com"
    		type     = "domain"
    		priority = 20
    		port     = 80
    		weight   = 10
  		}
	}
`, name)
}

func AliCloudCDNDomainConfigBasicDependence1(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "default" {
  		domain_name = "%s"
  		cdn_type    = "web"
  		scope       = "overseas"
  		sources {
    		content  = "www.aliyuntest.com"
    		type     = "domain"
    		priority = 20
    		port     = 80
    		weight   = 10
  		}
	}

	resource "alicloud_cdn_domain_config" "parent" {
  		domain_name   = alicloud_cdn_domain_new.default.domain_name
  		function_name = "condition"
  		function_args {
    		arg_name  = "rule"
    		arg_value = "{\"match\":{\"logic\":\"and\",\"criteria\":[{\"matchType\":\"clientipVer\",\"matchObject\":\"CONNECTING_IP\",\"matchOperator\":\"equals\",\"matchValue\":\"v6\",\"negate\":false}]},\"name\":\"example\",\"status\":\"enable\"}"
  		}
	}
`, name)
}

func AliCloudCDNDomainConfigBasicDependence2(name string) string {
	return fmt.Sprintf(`
	data "alicloud_oss_buckets" "default" {
	}

	resource "alicloud_cdn_domain_new" "default" {
  		domain_name = "tf-testacc%s-oss.alicloud-provider.cn"
  		cdn_type    = "web"
  		scope       = "overseas"
  		sources {
    		content  = "${data.alicloud_oss_buckets.default.buckets.0.name}.${data.alicloud_oss_buckets.default.buckets.0.extranet_endpoint}"
    		type     = "oss"
    		priority = 20
    		port     = 80
    		weight   = 10
  		}
	}
`, name)
}

package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCDNDomainConfig_filetype_based_ttl_set_new(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudCDNDomainConfig_ip_allow_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "ip_list",
		"arg_value": "110.110.110.110",
	}))
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
						"domain_name":   name,
						"function_name": "ip_allow_list_set",
						"function_args." + hashcode1 + ".arg_name":  "ip_list",
						"function_args." + hashcode1 + ".arg_value": "110.110.110.110",
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

func TestAccAlicloudCDNDomainConfig_referer_white_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "refer_domain_allow_list",
		"arg_value": "110.110.110.110",
	}))
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
						"domain_name":   name,
						"function_name": "referer_white_list_set",
						"function_args." + hashcode1 + ".arg_name":  "refer_domain_allow_list",
						"function_args." + hashcode1 + ".arg_value": "110.110.110.110",
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

func TestAccAlicloudCDNDomainConfig_referer_black_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "refer_domain_deny_list",
		"arg_value": "110.110.110.110",
	}))
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
						"domain_name":   name,
						"function_name": "referer_black_list_set",
						"function_args." + hashcode1 + ".arg_name":  "refer_domain_deny_list",
						"function_args." + hashcode1 + ".arg_value": "110.110.110.110",
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

func TestAccAlicloudCDNDomainConfig_filetype_based_ttl_set(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "ttl",
		"arg_value": "300",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "file_type",
		"arg_value": "jpg",
	}))
	hashcode3 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "weight",
		"arg_value": "1",
	}))

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
						"domain_name":   name,
						"function_name": "filetype_based_ttl_set",
						"function_args." + hashcode1 + ".arg_name":  "ttl",
						"function_args." + hashcode1 + ".arg_value": "300",
						"function_args." + hashcode2 + ".arg_name":  "file_type",
						"function_args." + hashcode2 + ".arg_value": "jpg",
						"function_args." + hashcode3 + ".arg_name":  "weight",
						"function_args." + hashcode3 + ".arg_value": "1",
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

func TestAccAlicloudCDNDomainConfig_oss_auth(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := defaultRegionToTest + strconv.Itoa(rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence_oss)

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
					"function_name": "oss_auth",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "oss_bucket_id",
							"arg_value": "${alicloud_oss_bucket.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(

					testAccCheck(map[string]string{
						"domain_name":     fmt.Sprintf("tf-testacc%s-oss.xiaozhu.com", name),
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

func TestAccAlicloudCDNDomainConfig_ip_black_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "ip_list",
		"arg_value": "110.110.110.110",
	}))

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
						"domain_name":   name,
						"function_name": "ip_black_list_set",
						"function_args." + hashcode1 + ".arg_name":  "ip_list",
						"function_args." + hashcode1 + ".arg_value": "110.110.110.110",
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

func TestAccAlicloudCDNDomainConfig_ip_white_list(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "ip_list",
		"arg_value": "110.110.110.110",
	}))

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
						"domain_name":   name,
						"function_name": "ip_white_list_set",
						"function_args." + hashcode1 + ".arg_name":  "ip_list",
						"function_args." + hashcode1 + ".arg_value": "110.110.110.110",
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

func TestAccAlicloudCDNDomainConfig_error_page(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "error_code",
		"arg_value": "502",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "rewrite_page",
		"arg_value": "http://www.xiaozhu.com",
	}))
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
					"function_name": "error_page",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "error_code",
							"arg_value": "502",
						},
						{
							"arg_name":  "rewrite_page",
							"arg_value": "http://www.xiaozhu.com",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":   name,
						"function_name": "error_page",
						"function_args." + hashcode1 + ".arg_name":  "error_code",
						"function_args." + hashcode1 + ".arg_value": "502",
						"function_args." + hashcode2 + ".arg_name":  "rewrite_page",
						"function_args." + hashcode2 + ".arg_value": "http://www.xiaozhu.com",
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

func TestAccAlicloudCDNDomainConfig_set_req_host_header(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "domain_name",
		"arg_value": "xiaozhu.com",
	}))

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
					"function_name": "set_req_host_header",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "domain_name",
							"arg_value": "xiaozhu.com",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":   name,
						"function_name": "set_req_host_header",
						"function_args." + hashcode1 + ".arg_name":  "domain_name",
						"function_args." + hashcode1 + ".arg_value": "xiaozhu.com",
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

func TestAccAlicloudCDNDomainConfig_set_hashkey_args(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "disable",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "set_hashkey_args",
						"function_args." + hashcode1 + ".arg_name":  "disable",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_aliauth(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "auth_type",
		"arg_value": "no_auth",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "ali_auth_delta",
		"arg_value": "1800",
	}))

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
						"domain_name":   name,
						"function_name": "aliauth",
						"function_args." + hashcode1 + ".arg_name":  "auth_type",
						"function_args." + hashcode1 + ".arg_value": "no_auth",
						"function_args." + hashcode2 + ".arg_name":  "ali_auth_delta",
						"function_args." + hashcode2 + ".arg_value": "1800",
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

func TestAccAlicloudCDNDomainConfig_set_resp_header(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "key",
		"arg_value": "expires",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "value",
		"arg_value": "null",
	}))

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
						"domain_name":   name,
						"function_name": "set_resp_header",
						"function_args." + hashcode1 + ".arg_name":  "key",
						"function_args." + hashcode1 + ".arg_value": "expires",
						"function_args." + hashcode2 + ".arg_name":  "value",
						"function_args." + hashcode2 + ".arg_value": "null",
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

func TestAccAlicloudCDNDomainConfig_https_force(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enable",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "https_force",
						"function_args." + hashcode1 + ".arg_name":  "enable",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_http_force(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enable",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "http_force",
						"function_args." + hashcode1 + ".arg_name":  "enable",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_https_option(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "http2",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "https_option",
						"function_args." + hashcode1 + ".arg_name":  "http2",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_l2_oss_key(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "private_oss_auth",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "l2_oss_key",
						"function_args." + hashcode1 + ".arg_name":  "private_oss_auth",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_forward_scheme(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enable",
		"arg_value": "on",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "scheme_origin",
		"arg_value": "https",
	}))

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
						"domain_name":   name,
						"function_name": "forward_scheme",
						"function_args." + hashcode1 + ".arg_name":  "enable",
						"function_args." + hashcode1 + ".arg_value": "on",
						"function_args." + hashcode2 + ".arg_name":  "scheme_origin",
						"function_args." + hashcode2 + ".arg_value": "https",
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

func SkipTestAccAlicloudCDNDomainConfig_green_manager(t *testing.T) {
	// the function: green_manager has been deleted
	t.Skip()
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enable",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "green_manager",
						"function_args." + hashcode1 + ".arg_name":  "enable",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_tmd_signature(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "ttl",
		"arg_value": "10",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "name",
		"arg_value": "tmd_signature_test",
	}))
	hashcode3 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "path",
		"arg_value": "/tftest",
	}))
	hashcode4 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "pathType",
		"arg_value": "0",
	}))
	hashcode5 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "interval",
		"arg_value": "20",
	}))
	hashcode6 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "count",
		"arg_value": "500",
	}))
	hashcode7 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "action",
		"arg_value": "1",
	}))
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
						"domain_name":   name,
						"function_name": "tmd_signature",
						"function_args." + hashcode1 + ".arg_name":  "ttl",
						"function_args." + hashcode1 + ".arg_value": "10",
						"function_args." + hashcode2 + ".arg_name":  "name",
						"function_args." + hashcode2 + ".arg_value": "tmd_signature_test",
						"function_args." + hashcode3 + ".arg_name":  "path",
						"function_args." + hashcode3 + ".arg_value": "/tftest",
						"function_args." + hashcode4 + ".arg_name":  "pathType",
						"function_args." + hashcode4 + ".arg_value": "0",
						"function_args." + hashcode5 + ".arg_name":  "interval",
						"function_args." + hashcode5 + ".arg_value": "20",
						"function_args." + hashcode6 + ".arg_name":  "count",
						"function_args." + hashcode6 + ".arg_value": "500",
						"function_args." + hashcode7 + ".arg_name":  "action",
						"function_args." + hashcode7 + ".arg_value": "1",
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
//func TestAccAlicloudCDNDomainConfig_dynamic(t *testing.T) {
//	var v *cdn.DomainConfigInDescribeCdnDomainConfigs
//
//	resourceId := "alicloud_cdn_domain_config.default"
//	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)
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
//	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
//	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
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
//		Providers:     testAccProviders,
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
//						"function_args." + hashcode1 + ".arg_name":  "enable",
//						"function_args." + hashcode1 + ".arg_value": "on",
//					}),
//				),
//			},
//			{
//				ResourceName:            resourceId,
//				ImportState:             true,
//				ImportStateVerify:       true,
//				ImportStateVerifyIgnore: []string{"domain_name"},
//			},
//		},
//	})
//}

func TestAccAlicloudCDNDomainConfig_set_req_header(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "key",
		"arg_value": "tftest",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "value",
		"arg_value": "tftest",
	}))

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
						"domain_name":   name,
						"function_name": "set_req_header",
						"function_args." + hashcode1 + ".arg_name":  "key",
						"function_args." + hashcode1 + ".arg_value": "tftest",
						"function_args." + hashcode2 + ".arg_name":  "value",
						"function_args." + hashcode2 + ".arg_value": "tftest",
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

func TestAccAlicloudCDNDomainConfig_range(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enable",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "range",
						"function_args." + hashcode1 + ".arg_name":  "enable",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_video_seek(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enable",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "video_seek",
						"function_args." + hashcode1 + ".arg_name":  "enable",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_https_tls_version(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "tls10",
		"arg_value": "on",
	}))

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
					"function_name": "https_tls_version",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "tls10",
							"arg_value": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":   name,
						"function_name": "https_tls_version",
						"function_args." + hashcode1 + ".arg_name":  "tls10",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_HSTS(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enabled",
		"arg_value": "on",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "https_hsts_max_age",
		"arg_value": "5184000",
	}))

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
						"domain_name":   name,
						"function_name": "HSTS",
						"function_args." + hashcode1 + ".arg_name":  "enabled",
						"function_args." + hashcode1 + ".arg_value": "on",
						"function_args." + hashcode2 + ".arg_name":  "https_hsts_max_age",
						"function_args." + hashcode2 + ".arg_value": "5184000",
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

func TestAccAlicloudCDNDomainConfig_filetype_force_ttl_code(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "file_type",
		"arg_value": "jpg",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "code_string",
		"arg_value": "302=0",
	}))

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
						"domain_name":   name,
						"function_name": "filetype_force_ttl_code",
						"function_args." + hashcode1 + ".arg_name":  "file_type",
						"function_args." + hashcode1 + ".arg_value": "jpg",
						"function_args." + hashcode2 + ".arg_name":  "code_string",
						"function_args." + hashcode2 + ".arg_value": "302=0",
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

func TestAccAlicloudCDNDomainConfig_path_force_ttl_code(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "code_string",
		"arg_value": "302=0",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "path",
		"arg_value": "/tftest",
	}))

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
						"domain_name":   name,
						"function_name": "path_force_ttl_code",
						"function_args." + hashcode1 + ".arg_name":  "code_string",
						"function_args." + hashcode1 + ".arg_value": "302=0",
						"function_args." + hashcode2 + ".arg_name":  "path",
						"function_args." + hashcode2 + ".arg_value": "/tftest",
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

func TestAccAlicloudCDNDomainConfig_gzip(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enable",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "gzip",
						"function_args." + hashcode1 + ".arg_name":  "enable",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_tesla(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enable",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "tesla",
						"function_args." + hashcode1 + ".arg_name":  "enable",
						"function_args." + hashcode1 + ".arg_value": "on",
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

func TestAccAlicloudCDNDomainConfig_https_origin_sni(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enabled",
		"arg_value": "on",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "https_origin_sni",
		"arg_value": "xiaozhu.com",
	}))
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
					"function_name": "https_origin_sni",
					"function_args": []map[string]interface{}{
						{
							"arg_name":  "enabled",
							"arg_value": "on",
						},
						{
							"arg_name":  "https_origin_sni",
							"arg_value": "xiaozhu.com",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":   name,
						"function_name": "https_origin_sni",
						"function_args." + hashcode1 + ".arg_name":  "enabled",
						"function_args." + hashcode1 + ".arg_value": "on",
						"function_args." + hashcode2 + ".arg_name":  "https_origin_sni",
						"function_args." + hashcode2 + ".arg_value": "xiaozhu.com",
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

func TestAccAlicloudCDNDomainConfig_brotli(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "brotli_level",
		"arg_value": "1",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "enable",
		"arg_value": "on",
	}))

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
						"domain_name":   name,
						"function_name": "brotli",
						"function_args." + hashcode1 + ".arg_name":  "brotli_level",
						"function_args." + hashcode1 + ".arg_value": "1",
						"function_args." + hashcode2 + ".arg_name":  "enable",
						"function_args." + hashcode2 + ".arg_value": "on",
					}),
				),
			},
			/*{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_name"},
			},*/
		},
	})
}

func TestAccAlicloudCDNDomainConfig_ali_ua(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "ua",
		"arg_value": "User-Agent",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "type",
		"arg_value": "black",
	}))

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
						"domain_name":   name,
						"function_name": "ali_ua",
						"function_args." + hashcode1 + ".arg_name":  "ua",
						"function_args." + hashcode1 + ".arg_value": "User-Agent",
						"function_args." + hashcode2 + ".arg_name":  "type",
						"function_args." + hashcode2 + ".arg_value": "black",
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

func TestAccAlicloudCDNDomainConfig_host_redirect(t *testing.T) {
	var v *cdn.DomainConfigInDescribeCdnDomainConfigs

	resourceId := "alicloud_cdn_domain_config.default"
	ra := resourceAttrInit(resourceId, cdnDomainConfigBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainConfigDependence)
	hashcode1 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "regex",
		"arg_value": "/$",
	}))
	hashcode2 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "replacement",
		"arg_value": "/go/act/sale/tbzlsy.php",
	}))
	hashcode3 := strconv.Itoa(expirationCdnDomainConfigHash(map[string]interface{}{
		"arg_name":  "flag",
		"arg_value": "redirect",
	}))

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
						"domain_name":   name,
						"function_name": "host_redirect",
						"function_args." + hashcode1 + ".arg_name":  "regex",
						"function_args." + hashcode1 + ".arg_value": "/$",
						"function_args." + hashcode2 + ".arg_name":  "replacement",
						"function_args." + hashcode2 + ".arg_value": "/go/act/sale/tbzlsy.php",
						"function_args." + hashcode3 + ".arg_name":  "flag",
						"function_args." + hashcode3 + ".arg_value": "redirect",
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
         content = "${alicloud_oss_bucket.default.bucket}.${alicloud_oss_bucket.default.extranet_endpoint}"
         type = "oss"
         priority = 20
         port = 80
         weight = 10
      }
	}

	resource "alicloud_oss_bucket" "default" {
	  bucket = "tf-testacc-domain-config-%s"
	}
`, name, name)
}

var cdnDomainConfigBasicMap = map[string]string{
	"domain_name":   CHECKSET,
	"function_name": CHECKSET,
}

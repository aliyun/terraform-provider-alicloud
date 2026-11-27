package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEsaErrorPagesRedirectRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_error_pages_redirect_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaErrorPagesRedirectRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaErrorPagesRedirectRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaErrorPagesRedirectRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id": "${data.alicloud_esa_sites.default.sites.0.id}",
					"error_pages_redirect": []map[string]interface{}{
						{
							"status_code": "400",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "403",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "404",
							"target_url":  "https://alicloud.com/tf1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":                acctest.CHECKSET,
						"error_pages_redirect.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"error_pages_redirect": []map[string]interface{}{
						{
							"status_code": "403",
							"target_url":  "https://alicloud.com/tf1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"error_pages_redirect.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"error_pages_redirect": []map[string]interface{}{
						{
							"status_code": "400",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "403",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "404",
							"target_url":  "https://alicloud.com/tf1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"error_pages_redirect.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sequence": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sequence": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func TestAccAliCloudEsaErrorPagesRedirectRule_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_error_pages_redirect_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaErrorPagesRedirectRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaErrorPagesRedirectRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaErrorPagesRedirectRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"site_version": "0",
					"sequence":     "0",
					"error_pages_redirect": []map[string]interface{}{
						{
							"status_code": "400",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "403",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "404",
							"target_url":  "https://alicloud.com/tf1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":                acctest.CHECKSET,
						"site_version":           "0",
						"sequence":               "0",
						"error_pages_redirect.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func TestAccAliCloudEsaErrorPagesRedirectRule_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_error_pages_redirect_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaErrorPagesRedirectRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaErrorPagesRedirectRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaErrorPagesRedirectRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.id}",
					"rule_enable": "on",
					"rule_name":   name,
					"rule":        true,
					"error_pages_redirect": []map[string]interface{}{
						{
							"status_code": "400",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "403",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "404",
							"target_url":  "https://alicloud.com/tf1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":                acctest.CHECKSET,
						"rule_enable":            "on",
						"rule_name":              name,
						"rule":                   "true",
						"error_pages_redirect.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": `(http.host eq \"video.example.com\")`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"error_pages_redirect": []map[string]interface{}{
						{
							"status_code": "403",
							"target_url":  "https://alicloud.com/tf1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"error_pages_redirect.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"error_pages_redirect": []map[string]interface{}{
						{
							"status_code": "400",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "403",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "404",
							"target_url":  "https://alicloud.com/tf1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"error_pages_redirect.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sequence": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sequence": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func TestAccAliCloudEsaErrorPagesRedirectRule_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_error_pages_redirect_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaErrorPagesRedirectRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaErrorPagesRedirectRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaErrorPagesRedirectRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"rule_enable":  "on",
					"rule_name":    name,
					"rule":         true,
					"site_version": "0",
					"sequence":     "1",
					"error_pages_redirect": []map[string]interface{}{
						{
							"status_code": "400",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "403",
							"target_url":  "https://alicloud.com/tf1",
						},
						{
							"status_code": "404",
							"target_url":  "https://alicloud.com/tf1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":                acctest.CHECKSET,
						"rule_enable":            "on",
						"rule_name":              name,
						"rule":                   "true",
						"site_version":           "0",
						"sequence":               "1",
						"error_pages_redirect.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

var AliCloudEsaErrorPagesRedirectRuleMap0 = map[string]string{
	"id":          acctest.CHECKSET,
	"config_id":   acctest.CHECKSET,
	"config_type": acctest.CHECKSET,
}

func AliCloudEsaErrorPagesRedirectRuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}
`, name)
}

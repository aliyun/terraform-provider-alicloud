package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Esa ImageTransform. >>> Resource test cases, automatically generated.
// Case resource_ImageTransform_test 13068
func TestAccAliCloudEsaImageTransform_basic13068(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_image_transform.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaImageTransformMap13068)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaImageTransform")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaImageTransformBasicDependence13068)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id": "${data.alicloud_esa_sites.default.sites.0.id}",
					"enable":  "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id": acctest.CHECKSET,
						"enable":  "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_avif": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_avif": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_avif": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_avif": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_webp": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_webp": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_webp": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_webp": "off",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudEsaImageTransform_basic13069(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_image_transform.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaImageTransformMap13068)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaImageTransform")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaImageTransformBasicDependence13068)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.id}",
					"enable":      "on",
					"rule_enable": "on",
					"rule_name":   name,
					"rule":        true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":     acctest.CHECKSET,
						"enable":      "on",
						"rule_enable": "on",
						"rule_name":   name,
						"rule":        "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_avif": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_avif": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_avif": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_avif": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_webp": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_webp": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_webp": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_webp": "off",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudEsaImageTransform_basic13069_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_image_transform.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaImageTransformMap13068)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaImageTransform")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaImageTransformBasicDependence13068)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"enable":       "on",
					"rule_enable":  "on",
					"rule_name":    name,
					"rule":         true,
					"sequence":     "1",
					"site_version": "0",
					"auto_avif":    "on",
					"auto_webp":    "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":      acctest.CHECKSET,
						"enable":       "on",
						"rule_enable":  "on",
						"rule_name":    name,
						"rule":         "true",
						"sequence":     "1",
						"site_version": "0",
						"auto_avif":    "on",
						"auto_webp":    "on",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudEsaImageTransformMap13068 = map[string]string{
	"config_id":   acctest.CHECKSET,
	"config_type": acctest.CHECKSET,
}

func AliCloudEsaImageTransformBasicDependence13068(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "standardplan"
}
`, name)
}

// Test Esa ImageTransform. <<< Resource test cases, automatically generated.

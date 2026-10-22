// Package esa. This file is generated automatically. Please do not modify it manually, thank you!
package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Test Esa ImageTransform. >>> Resource test cases, automatically generated.
// Case resource_ImageTransform_test 13069
func TestAccAliCloudEsaImageTransform_basic13069(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccEsaImageTransform%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			if factory := acctest.ProtoV5ProviderFactories()["alicloud"]; factory == nil {
				t.Fatal("SDK acceptance tests require the real alicloud protocol 5 factory")
			}
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEsaImageTransform_basic13069_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_image_transform.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_avif", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_enable", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_webp", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "enable", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule", "http.host eq \"video.example.com\""),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_name", "rule_example"),

					resource.TestCheckResourceAttrSet("alicloud_esa_image_transform.default", "config_id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_esa_image_transform.default", tfjsonpath.New("config_type"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudEsaImageTransform_basic13069_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_image_transform.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_avif", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_webp", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule", "http.host eq \"videoo.example.com\""),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_name", "rule_viedoo"),
				),
			},
			{
				Config: testAccAliCloudEsaImageTransform_basic13069_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_image_transform.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_avif", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_webp", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule", "http.host eq \"videoo.example.com\""),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_name", "rule_viedoo"),
				),
			},
			{
				Config: testAccAliCloudEsaImageTransform_basic13069_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_image_transform.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_avif", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_webp", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule", "http.host eq \"image.example.com\""),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_name", "rule_viedoo"),
				),
			},
			{
				Config: testAccAliCloudEsaImageTransform_basic13069_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_image_transform.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_avif", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_webp", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule", "http.host eq \"image.example.com\""),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_name", "rule_image"),
				),
			},
			{
				Config: testAccAliCloudEsaImageTransform_basic13069_step5(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_image_transform.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_avif", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_webp", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule", "http.host eq \"image.example.com\""),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_name", "rule_image"),
				),
			},
			{
				Config: testAccAliCloudEsaImageTransform_basic13069_step6(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_image_transform.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_avif", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_webp", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule", "http.host eq \"image.example.com\""),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_name", "rule_image"),
				),
			},
			{
				Config: testAccAliCloudEsaImageTransform_basic13069_step7(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_image_transform.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_avif", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_enable", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "auto_webp", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "enable", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule", "http.host eq \"video.example.com\""),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_image_transform.default", "rule_name", "rule_example"),
				),
			},
			{
				ResourceName:            "alicloud_esa_image_transform.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudEsaImageTransform_basic13069_base(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_test_ImageTransform" {
  type         = "NS"
  auto_renew   = false
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_test_ImageTransform" {
  site_name   = "${lower(var.name)}.com"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_test_ImageTransform.id
  coverage    = "overseas"
  access_type = "NS"
  flatten_mode = "flatten_at_root"
}


`, name)
}

func testAccAliCloudEsaImageTransform_basic13069_step0(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_image_transform" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ImageTransform.id}"

        auto_avif = "on"

        rule_enable = "off"

        auto_webp = "on"

        enable = "off"

        rule = "http.host eq \"video.example.com\""

        sequence = 1

        site_version = 0

        rule_name = "rule_example"

}
`, name, testAccAliCloudEsaImageTransform_basic13069_base(name))
}

func testAccAliCloudEsaImageTransform_basic13069_step1(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_image_transform" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ImageTransform.id}"

        auto_avif = "off"

        rule_enable = "on"

        auto_webp = "off"

        enable = "on"

        rule = "http.host eq \"videoo.example.com\""

        sequence = 1

        site_version = 0

        rule_name = "rule_viedoo"

}
`, name, testAccAliCloudEsaImageTransform_basic13069_base(name))
}

func testAccAliCloudEsaImageTransform_basic13069_step2(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_image_transform" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ImageTransform.id}"

        auto_avif = "off"

        rule_enable = "on"

        auto_webp = "off"

        enable = "on"

        rule = "http.host eq \"videoo.example.com\""

        sequence = 1

        site_version = 0

        rule_name = "rule_viedoo"

}
`, name, testAccAliCloudEsaImageTransform_basic13069_base(name))
}

func testAccAliCloudEsaImageTransform_basic13069_step3(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_image_transform" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ImageTransform.id}"

        auto_avif = "off"

        rule_enable = "on"

        auto_webp = "off"

        enable = "on"

        rule = "http.host eq \"image.example.com\""

        sequence = 1

        site_version = 0

        rule_name = "rule_viedoo"

}
`, name, testAccAliCloudEsaImageTransform_basic13069_base(name))
}

func testAccAliCloudEsaImageTransform_basic13069_step4(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_image_transform" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ImageTransform.id}"

        auto_avif = "off"

        rule_enable = "on"

        auto_webp = "off"

        enable = "on"

        rule = "http.host eq \"image.example.com\""

        sequence = 1

        site_version = 0

        rule_name = "rule_image"

}
`, name, testAccAliCloudEsaImageTransform_basic13069_base(name))
}

func testAccAliCloudEsaImageTransform_basic13069_step5(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_image_transform" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ImageTransform.id}"

        auto_avif = "off"

        rule_enable = "on"

        auto_webp = "off"

        enable = "on"

        rule = "http.host eq \"image.example.com\""

        sequence = 1

        site_version = 0

        rule_name = "rule_image"

}
`, name, testAccAliCloudEsaImageTransform_basic13069_base(name))
}

func testAccAliCloudEsaImageTransform_basic13069_step6(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_image_transform" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ImageTransform.id}"

        auto_avif = "off"

        rule_enable = "on"

        auto_webp = "off"

        enable = "on"

        rule = "http.host eq \"image.example.com\""

        sequence = 1

        site_version = 0

        rule_name = "rule_image"

}
`, name, testAccAliCloudEsaImageTransform_basic13069_base(name))
}

func testAccAliCloudEsaImageTransform_basic13069_step7(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_image_transform" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ImageTransform.id}"

        auto_avif = "off"

        rule_enable = "off"

        auto_webp = "off"

        enable = "off"

        rule = "http.host eq \"video.example.com\""

        sequence = 1

        site_version = 0

        rule_name = "rule_example"

}
`, name, testAccAliCloudEsaImageTransform_basic13069_base(name))
}

// Test Esa ImageTransform. <<< Resource test cases, automatically generated.

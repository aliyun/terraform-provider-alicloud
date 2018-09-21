package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudOssBucketsDataSource_basic(t *testing.T) {
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOssBucketsDataSourceBasic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_oss_buckets.buckets"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.name", fmt.Sprintf("tf-testacc-bucket-ds-basic-%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.acl", "public-read"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.extranet_endpoint"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.intranet_endpoint"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.location"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.owner"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.storage_class", "Standard"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.creation_date"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.website.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.logging.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.referer_config.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.referer_config.0.allow_empty", "true"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.referer_config.0.referers.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudOssBucketsDataSource_full(t *testing.T) {
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOssBucketsDataSourceFull(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_oss_buckets.buckets"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.name", fmt.Sprintf("tf-testacc-bucket-ds-full-%d-sample", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.acl", "public-read"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.extranet_endpoint"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.intranet_endpoint"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.location"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.owner"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.storage_class", "Standard"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.buckets", "buckets.0.creation_date"),

					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.0.allowed_headers.0", "authorization"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.0.allowed_methods.0", "PUT"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.0.allowed_methods.1", "GET"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.0.allowed_origins.0", "*"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.0.expose_headers.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.0.max_age_seconds", "0"),

					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.1.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.1.allowed_headers.0", "authorization"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.1.allowed_methods.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.1.allowed_methods.0", "GET"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.1.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.1.allowed_origins.0", "http://www.a.com"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.1.expose_headers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.1.expose_headers.0", "x-oss-test"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.cors_rules.1.max_age_seconds", "100"),

					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.website.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.website.0.index_document", "index.html"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.website.0.error_document", "error.html"),

					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.logging.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.logging.0.target_bucket", fmt.Sprintf("tf-testacc-bucket-ds-full-%d-log", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.logging.0.target_prefix", "log/"),

					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.referer_config.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.referer_config.0.allow_empty", "false"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.referer_config.0.referers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.referer_config.0.referers.0", "http://www.aliyun.com"),

					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.0.id", "rule1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.0.expiration.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.0.expiration.0.days", "365"),

					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.1.id", "rule2"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.1.enabled", "true"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.1.expiration.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.buckets", "buckets.0.lifecycle_rule.1.expiration.0.date", "2018-01-12"),
				),
			},
		},
	})
}

func testAccCheckAlicloudOssBucketsDataSourceBasic(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-bucket-ds-basic-%d"
}

resource "alicloud_oss_bucket" "sample_bucket" {
	bucket = "${var.name}"
	acl = "public-read"
}

data "alicloud_oss_buckets" "buckets" {
    name_regex = "${alicloud_oss_bucket.sample_bucket.bucket}"
}
`, randInt)
}

func testAccCheckAlicloudOssBucketsDataSourceFull(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-bucket-ds-full-%d"
}

resource "alicloud_oss_bucket" "log_bucket"{
	bucket = "${var.name}-log"
}

resource "alicloud_oss_bucket" "sample_bucket" {
	bucket = "${var.name}-sample"
	acl = "public-read"

    cors_rule = [
    	{
			allowed_origins=["*"]
			allowed_methods=["PUT","GET"]
			allowed_headers=["authorization"]
		},
		{
			allowed_origins=["http://www.a.com"]
			allowed_methods=["GET"]
			allowed_headers=["authorization"]
			expose_headers=["x-oss-test"]
			max_age_seconds=100
		}
    ]

	website = {
		index_document = "index.html"
		error_document = "error.html"
	}

	logging {
		target_bucket = "${alicloud_oss_bucket.log_bucket.id}"
		target_prefix = "log/"
	}

	referer_config {
		allow_empty = false
		referers = ["http://www.aliyun.com"]
	}

	lifecycle_rule = [
		{
			id = "rule1"
			prefix = "path1/"
			enabled = true
			expiration {
				days = 365
			}
		},
		{
			id = "rule2"
			prefix = "path2/"
			enabled = true
			expiration {
				date = "2018-01-12"
			}
		}
	]
}

data "alicloud_oss_buckets" "buckets" {
    name_regex = "${alicloud_oss_bucket.sample_bucket.bucket}"
}
`, randInt)
}

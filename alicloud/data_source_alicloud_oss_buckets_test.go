package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOssBucketsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_oss_buckets.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-bucket-%d", rand),
		dataSourceOssBucketsConfigDependence_basic)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket.default.bucket}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket.default.bucket}-fake",
		}),
	}
	var existOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#":                   "1",
			"names.#":                     "1",
			"buckets.0.name":              fmt.Sprintf("tf-testacc-bucket-%d-default", rand),
			"buckets.0.acl":               "public-read",
			"buckets.0.extranet_endpoint": CHECKSET,
			"buckets.0.intranet_endpoint": CHECKSET,
			"buckets.0.location":          CHECKSET,
			"buckets.0.owner":             CHECKSET,
			"buckets.0.storage_class":     "Standard",
			"buckets.0.redundancy_type":   "LRS",
			"buckets.0.creation_date":     CHECKSET,

			"buckets.0.cors_rules.#":                   "2",
			"buckets.0.cors_rules.0.allowed_headers.#": "1",
			"buckets.0.cors_rules.0.allowed_headers.0": "authorization",
			"buckets.0.cors_rules.0.allowed_methods.#": "2",
			"buckets.0.cors_rules.0.allowed_methods.0": "PUT",
			"buckets.0.cors_rules.0.allowed_methods.1": "GET",
			"buckets.0.cors_rules.0.allowed_origins.#": "1",
			"buckets.0.cors_rules.0.allowed_origins.0": "*",
			"buckets.0.cors_rules.0.expose_headers.#":  "0",
			"buckets.0.cors_rules.0.max_age_seconds":   "0",
			"buckets.0.cors_rules.1.allowed_headers.#": "1",
			"buckets.0.cors_rules.1.allowed_headers.0": "authorization",
			"buckets.0.cors_rules.1.allowed_methods.#": "1",
			"buckets.0.cors_rules.1.allowed_methods.0": "GET",
			"buckets.0.cors_rules.1.allowed_origins.#": "1",
			"buckets.0.cors_rules.1.allowed_origins.0": "http://www.a.com",
			"buckets.0.cors_rules.1.expose_headers.#":  "1",
			"buckets.0.cors_rules.1.expose_headers.0":  "x-oss-test",
			"buckets.0.cors_rules.1.max_age_seconds":   "100",

			"buckets.0.website.#":                "1",
			"buckets.0.website.0.index_document": "index.html",
			"buckets.0.website.0.error_document": "error.html",

			"buckets.0.logging.#":               "1",
			"buckets.0.logging.0.target_bucket": fmt.Sprintf("tf-testacc-bucket-%d-target", rand),
			"buckets.0.logging.0.target_prefix": "log/",

			"buckets.0.referer_config.#":             "1",
			"buckets.0.referer_config.0.allow_empty": "false",
			"buckets.0.referer_config.0.referers.#":  "1",
			"buckets.0.referer_config.0.referers.0":  "http://www.aliyun.com",

			"buckets.0.lifecycle_rule.#":                   "2",
			"buckets.0.lifecycle_rule.0.id":                "rule1",
			"buckets.0.lifecycle_rule.0.prefix":            "path1/",
			"buckets.0.lifecycle_rule.0.enabled":           "true",
			"buckets.0.lifecycle_rule.0.expiration.#":      "1",
			"buckets.0.lifecycle_rule.0.expiration.0.days": "365",
			"buckets.0.lifecycle_rule.1.id":                "rule2",
			"buckets.0.lifecycle_rule.1.prefix":            "path2/",
			"buckets.0.lifecycle_rule.1.enabled":           "true",
			"buckets.0.lifecycle_rule.1.expiration.#":      "1",
			"buckets.0.lifecycle_rule.1.expiration.0.date": "2018-01-12",

			"buckets.0.policy": "{\"Statement\":[{\"Action\":[\"oss:*\"],\"Effect\":\"Allow\",\"Resource\":[\"acs:oss:*:*:*\"]}],\"Version\":\"1\"}",

			"buckets.0.tags.key1": "value1",
			"buckets.0.tags.key2": "value2",
		}
	}

	var fakeOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#": "0",
			"names.#":   "0",
		}
	}

	var ossBucketsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOssBucketsMapFunc,
		fakeMapFunc:  fakeOssBucketsMapFunc,
	}

	ossBucketsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func TestAccAlicloudOssBucketsDataSource_sserule(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_oss_buckets.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-bucket-%d", rand),
		dataSourceOssBucketsConfigDependence_sserule)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket.default.bucket}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket.default.bucket}-fake",
		}),
	}
	var existOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#":                   "1",
			"names.#":                     "1",
			"buckets.0.name":              fmt.Sprintf("tf-testacc-bucket-%d-default", rand),
			"buckets.0.acl":               "public-read",
			"buckets.0.extranet_endpoint": CHECKSET,
			"buckets.0.intranet_endpoint": CHECKSET,
			"buckets.0.location":          CHECKSET,
			"buckets.0.owner":             CHECKSET,
			"buckets.0.storage_class":     "Standard",
			"buckets.0.redundancy_type":   "LRS",
			"buckets.0.creation_date":     CHECKSET,

			"buckets.0.server_side_encryption_rule.0.sse_algorithm": "AES256",
		}
	}

	var fakeOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#": "0",
			"names.#":   "0",
		}
	}

	var ossBucketsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOssBucketsMapFunc,
		fakeMapFunc:  fakeOssBucketsMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	ossBucketsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func TestAccAlicloudOssBucketsDataSource_sserule_with_kmsid(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_oss_buckets.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-bucket-%d", rand),
		dataSourceOssBucketsConfigDependence_sserule_with_kmsid)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket.default.bucket}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket.default.bucket}-fake",
		}),
	}
	var existOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#":                   "1",
			"names.#":                     "1",
			"buckets.0.name":              fmt.Sprintf("tf-testacc-bucket-%d-default", rand),
			"buckets.0.acl":               "public-read",
			"buckets.0.extranet_endpoint": CHECKSET,
			"buckets.0.intranet_endpoint": CHECKSET,
			"buckets.0.location":          CHECKSET,
			"buckets.0.owner":             CHECKSET,
			"buckets.0.storage_class":     "Standard",
			"buckets.0.creation_date":     CHECKSET,

			"buckets.0.server_side_encryption_rule.0.sse_algorithm":     "KMS",
			"buckets.0.server_side_encryption_rule.0.kms_master_key_id": "kms-id",
		}
	}

	var fakeOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#": "0",
			"names.#":   "0",
		}
	}

	var ossBucketsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOssBucketsMapFunc,
		fakeMapFunc:  fakeOssBucketsMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	ossBucketsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func TestAccAlicloudOssBucketsDataSource_versioning(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_oss_buckets.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-bucket-%d", rand),
		dataSourceOssBucketsConfigDependence_versioning)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket.default.bucket}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oss_bucket.default.bucket}-fake",
		}),
	}
	var existOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#":                   "1",
			"names.#":                     "1",
			"buckets.0.name":              fmt.Sprintf("tf-testacc-bucket-%d-default", rand),
			"buckets.0.acl":               "public-read",
			"buckets.0.extranet_endpoint": CHECKSET,
			"buckets.0.intranet_endpoint": CHECKSET,
			"buckets.0.location":          CHECKSET,
			"buckets.0.owner":             CHECKSET,
			"buckets.0.storage_class":     "Standard",
			"buckets.0.redundancy_type":   "LRS",
			"buckets.0.creation_date":     CHECKSET,

			"buckets.0.versioning.0.status": "Enabled",
		}
	}

	var fakeOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"buckets.#": "0",
			"names.#":   "0",
		}
	}

	var ossBucketsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOssBucketsMapFunc,
		fakeMapFunc:  fakeOssBucketsMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	ossBucketsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceOssBucketsConfigDependence_basic(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_oss_bucket" "target"{
	bucket = "${var.name}-target"
}

resource "alicloud_oss_bucket" "default" {
	bucket = "${var.name}-default"
	acl = "public-read"

    cors_rule {
			allowed_origins=["*"]
			allowed_methods=["PUT","GET"]
			allowed_headers=["authorization"]
		}
	cors_rule {
			allowed_origins=["http://www.a.com"]
			allowed_methods=["GET"]
			allowed_headers=["authorization"]
			expose_headers=["x-oss-test"]
			max_age_seconds=100
		}

	website {
		index_document = "index.html"
		error_document = "error.html"
	}

	logging {
		target_bucket = "${alicloud_oss_bucket.target.id}"
		target_prefix = "log/"
	}

	referer_config {
		allow_empty = false
		referers = ["http://www.aliyun.com"]
	}

	lifecycle_rule {
			id = "rule1"
			prefix = "path1/"
			enabled = true
			expiration {
				days = 365
			}
		}
	lifecycle_rule {
			id = "rule2"
			prefix = "path2/"
			enabled = true
			expiration {
				date = "2018-01-12"
			}
		}
    policy = "{\"Statement\":[{\"Action\":[\"oss:*\"],\"Effect\":\"Allow\",\"Resource\":[\"acs:oss:*:*:*\"]}],\"Version\":\"1\"}"
	tags = {
		key1 = "value1",
		key2 = "value2",
	}
}

`, name)
}
func dataSourceOssBucketsConfigDependence_sserule(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_oss_bucket" "default" {
	bucket = "${var.name}-default"
	acl = "public-read"

 	server_side_encryption_rule {
		sse_algorithm = "AES256"
	}
}
`, name)
}

func dataSourceOssBucketsConfigDependence_sserule_with_kmsid(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_oss_bucket" "default" {
	bucket = "${var.name}-default"
	acl = "public-read"

 	server_side_encryption_rule {
		sse_algorithm = "KMS"
		kms_master_key_id = "kms-id"
	}
}
`, name)
}

func dataSourceOssBucketsConfigDependence_versioning(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_oss_bucket" "default" {
	bucket = "${var.name}-default"
	acl = "public-read"

	versioning {
		status = "Enabled"
	}
}
`, name)
}

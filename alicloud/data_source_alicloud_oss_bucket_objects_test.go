package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOssBucketObjectsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_oss_bucket_objects.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-bucket-object-%d", rand),
		dataSourceOssBucketObjectsConfigDependence)

	bucketNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
		}),
	}

	keyRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}-fake",
		}),
	}

	keyPrefixConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_prefix":  "tf-sample/",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_prefix":  "tf-sample-fake/",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
			"key_prefix":  "tf-sample/",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
			"key_prefix":  "tf-sample-fake/",
		}),
	}

	var existOssBucketObjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"objects.#":                        "1",
			"objects.0.key":                    fmt.Sprintf("tf-sample/%s-object", fmt.Sprintf("tf-testacc-bucket-object-%d", rand)),
			"objects.0.acl":                    "public-read",
			"objects.0.content_type":           "text/plain",
			"objects.0.content_length":         CHECKSET,
			"objects.0.cache_control":          "max-age=0",
			"objects.0.content_disposition":    "attachment; filename=\"my-object\"",
			"objects.0.content_encoding":       "gzip",
			"objects.0.expires":                "Wed, 06 May 2020 00:00:00 GMT",
			"objects.0.content_md5":            "1STMBJqp4X5QEQsYTbRmkQ==",
			"objects.0.etag":                   CHECKSET,
			"objects.0.storage_class":          CHECKSET,
			"objects.0.last_modification_time": CHECKSET,
		}
	}

	var fakeOssBucketObjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"objects.#": "0",
		}
	}

	var ossBucketObjectsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOssBucketObjectsMapFunc,
		fakeMapFunc:  fakeOssBucketObjectsMapFunc,
	}

	ossBucketObjectsCheckInfo.dataSourceTestCheck(t, rand, bucketNameConf, keyRegexConf, keyPrefixConf, allConf)
}

func TestAccAlicloudOssBucketObjectsDataSource_versioning(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_oss_bucket_objects.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-bucket-object-%d", rand),
		dataSourceOssBucketObjectsConfigDependenceVersioning)

	bucketNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
		}),
	}

	keyRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}-fake",
		}),
	}

	keyPrefixConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_prefix":  "tf-sample/",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_prefix":  "tf-sample-fake/",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
			"key_prefix":  "tf-sample/",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
			"key_prefix":  "tf-sample-fake/",
		}),
	}

	var existOssBucketObjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"objects.#":                        "1",
			"objects.0.key":                    fmt.Sprintf("tf-sample/%s-object", fmt.Sprintf("tf-testacc-bucket-object-%d", rand)),
			"objects.0.acl":                    "private",
			"objects.0.content_type":           "text/plain",
			"objects.0.content_length":         CHECKSET,
			"objects.0.cache_control":          "max-age=0",
			"objects.0.content_disposition":    "attachment; filename=\"my-object\"",
			"objects.0.content_encoding":       "gzip",
			"objects.0.expires":                "Wed, 06 May 2020 00:00:00 GMT",
			"objects.0.content_md5":            "1STMBJqp4X5QEQsYTbRmkQ==",
			"objects.0.etag":                   CHECKSET,
			"objects.0.storage_class":          CHECKSET,
			"objects.0.last_modification_time": CHECKSET,
		}
	}

	var fakeOssBucketObjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"objects.#": "0",
		}
	}

	var ossBucketObjectsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOssBucketObjectsMapFunc,
		fakeMapFunc:  fakeOssBucketObjectsMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	ossBucketObjectsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, bucketNameConf, keyRegexConf, keyPrefixConf, allConf)
}

func dataSourceOssBucketObjectsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_oss_bucket" "default" {
	bucket = "${var.name}"
	acl = "private"
}

resource "alicloud_oss_bucket_object" "default" {
	bucket = "${alicloud_oss_bucket.default.bucket}"
	key = "tf-sample/${var.name}-object"
	content = "sample content"
	content_type = "text/plain"
	cache_control = "max-age=0"
	content_disposition = "attachment; filename=\"my-object\""
	content_encoding = "gzip"
	expires = "Wed, 06 May 2020 00:00:00 GMT"
	acl = "public-read"
}

`, name)
}
func dataSourceOssBucketObjectsConfigDependenceVersioning(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_oss_bucket" "default" {
	bucket = "${var.name}"
	acl = "private"
	force_destroy = true
	versioning {
		status = "Enabled"
	}
}

resource "alicloud_oss_bucket_object" "default" {
	bucket = "${alicloud_oss_bucket.default.bucket}"
	key = "tf-sample/${var.name}-object"
	content = "sample content"
	content_type = "text/plain"
	cache_control = "max-age=0"
	content_disposition = "attachment; filename=\"my-object\""
	content_encoding = "gzip"
	expires = "Wed, 06 May 2020 00:00:00 GMT"
}

`, name)
}

package oss_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudOssBucketObjectsDataSource_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_oss_bucket_objects.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-bucket-object-%d", rand),
		dataSourceOssBucketObjectsConfigDependence)

	bucketNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
		}),
	}

	keyRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}-fake",
		}),
	}

	keyPrefixConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_prefix":  "tf-sample/",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_prefix":  "tf-sample-fake/",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
			"key_prefix":  "tf-sample/",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"objects.0.content_length":         acctest.CHECKSET,
			"objects.0.cache_control":          "max-age=0",
			"objects.0.content_disposition":    "attachment; filename=\"my-object\"",
			"objects.0.content_encoding":       "gzip",
			"objects.0.expires":                "Wed, 06 May 2020 00:00:00 GMT",
			"objects.0.content_md5":            "1STMBJqp4X5QEQsYTbRmkQ==",
			"objects.0.etag":                   acctest.CHECKSET,
			"objects.0.storage_class":          acctest.CHECKSET,
			"objects.0.last_modification_time": acctest.CHECKSET,
			"objects.0.x_oss_hash_crc64ecma":   "653038730487567532",
		}
	}

	var fakeOssBucketObjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"objects.#": "0",
		}
	}

	var ossBucketObjectsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existOssBucketObjectsMapFunc,
		FakeMapFunc:  fakeOssBucketObjectsMapFunc,
	}

	ossBucketObjectsCheckInfo.DataSourceTestCheck(t, rand, bucketNameConf, keyRegexConf, keyPrefixConf, allConf)
}

func TestAccAliCloudOssBucketObjectsDataSource_versioning(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_oss_bucket_objects.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-bucket-object-%d", rand),
		dataSourceOssBucketObjectsConfigDependenceVersioning)

	bucketNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
		}),
	}

	keyRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}-fake",
		}),
	}

	keyPrefixConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_prefix":  "tf-sample/",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_prefix":  "tf-sample-fake/",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"bucket_name": "${alicloud_oss_bucket_object.default.bucket}",
			"key_regex":   "${alicloud_oss_bucket_object.default.key}",
			"key_prefix":  "tf-sample/",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"objects.0.content_length":         acctest.CHECKSET,
			"objects.0.cache_control":          "max-age=0",
			"objects.0.content_disposition":    "attachment; filename=\"my-object\"",
			"objects.0.content_encoding":       "gzip",
			"objects.0.expires":                "Wed, 06 May 2020 00:00:00 GMT",
			"objects.0.content_md5":            "1STMBJqp4X5QEQsYTbRmkQ==",
			"objects.0.etag":                   acctest.CHECKSET,
			"objects.0.storage_class":          acctest.CHECKSET,
			"objects.0.last_modification_time": acctest.CHECKSET,
			"objects.0.x_oss_hash_crc64ecma":   "653038730487567532",
		}
	}

	var fakeOssBucketObjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"objects.#": "0",
		}
	}

	var ossBucketObjectsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existOssBucketObjectsMapFunc,
		FakeMapFunc:  fakeOssBucketObjectsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	ossBucketObjectsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, bucketNameConf, keyRegexConf, keyPrefixConf, allConf)
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

resource "alicloud_oss_bucket_public_access_block" "default" {
	bucket              = alicloud_oss_bucket.default.bucket
	block_public_access = false
}

resource "alicloud_oss_bucket_object" "default" {
	bucket = "${alicloud_oss_bucket_public_access_block.default.bucket}"
	key = "tf-sample/${var.name}-object"
	content = "sample content"
	content_type = "text/plain"
	cache_control = "max-age=0"
	content_disposition = "attachment; filename=\"my-object\""
	content_encoding = "gzip"
	expires = "Wed, 06 May 2020 00:00:00 GMT"
	acl = "private"
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

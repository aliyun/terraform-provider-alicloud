package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case HttpsCofig测试 6361
func TestAccAliCloudOssBucketHttpsConfig_basic6361(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_https_config.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketHttpsConfigMap6361)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketHttpsConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckethttpsconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketHttpsConfigBasicDependence6361)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enable": "true",
					"tls_versions": []string{
						"TLSv1.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":         CHECKSET,
						"enable":         "true",
						"tls_versions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_versions": []string{
						"TLSv1.1", "TLSv1.2", "TLSv1.3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_versions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_versions": []string{
						"TLSv1.2"},
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_versions.#": "1",
						"bucket":         CHECKSET,
						"enable":         "true",
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

var AlicloudOssBucketHttpsConfigMap6361 = map[string]string{}

func AlicloudOssBucketHttpsConfigBasicDependence6361(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
}


`, name)
}

// Case HttpsCofig测试 6361  twin
func TestAccAliCloudOssBucketHttpsConfig_basic6361_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_https_config.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketHttpsConfigMap6361)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketHttpsConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckethttpsconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketHttpsConfigBasicDependence6361)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"enable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_versions": []string{
						"TLSv1.2"},
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_versions.#": "1",
						"bucket":         CHECKSET,
						"enable":         "true",
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

// Case CipherSuit测试
func TestAccAliCloudOssBucketHttpsConfig_cipherSuit(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_https_config.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketHttpsConfigMap6361)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketHttpsConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckethttpsconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketHttpsConfigBasicDependence6361)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enable": "true",
					"tls_versions": []string{
						"TLSv1.2", "TLSv1.3"},
					"cipher_suit": []map[string]interface{}{
						{
							"enable":              "true",
							"strong_cipher_suite": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                            CHECKSET,
						"enable":                            "true",
						"tls_versions.#":                    "2",
						"cipher_suit.#":                     "1",
						"cipher_suit.0.enable":              "true",
						"cipher_suit.0.strong_cipher_suite": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cipher_suit": []map[string]interface{}{
						{
							"enable":              "true",
							"custom_cipher_suite": []string{"TLS_AES_128_GCM_SHA256"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cipher_suit.#":                       "1",
						"cipher_suit.0.enable":                "true",
						"cipher_suit.0.custom_cipher_suite.#": "1",
						"cipher_suit.0.strong_cipher_suite":   "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cipher_suit": []map[string]interface{}{
						{
							"enable":                    "true",
							"tls13_custom_cipher_suite": []string{"TLS_AES_256_GCM_SHA384"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cipher_suit.#":                             "1",
						"cipher_suit.0.enable":                      "true",
						"cipher_suit.0.tls13_custom_cipher_suite.#": "1",
						"cipher_suit.0.custom_cipher_suite.#":       "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cipher_suit": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cipher_suit.#":                             "1",
						"cipher_suit.0.enable":                      "false",
						"cipher_suit.0.tls13_custom_cipher_suite.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cipher_suit": []map[string]interface{}{
						{
							"enable":              "true",
							"strong_cipher_suite": "false",
							"custom_cipher_suite": []string{"TLS_AES_128_GCM_SHA256"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cipher_suit.#":                       "1",
						"cipher_suit.0.enable":                "true",
						"cipher_suit.0.strong_cipher_suite":   "false",
						"cipher_suit.0.custom_cipher_suite.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cipher_suit": []map[string]interface{}{
						{
							"enable":              "true",
							"strong_cipher_suite": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cipher_suit.#":                       "1",
						"cipher_suit.0.enable":                "true",
						"cipher_suit.0.strong_cipher_suite":   "true",
						"cipher_suit.0.custom_cipher_suite.#": "0",
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

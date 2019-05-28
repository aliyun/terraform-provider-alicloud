package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudOssBucket_importBasic(t *testing.T) {
	resourceName := "alicloud_oss_bucket.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketBasicConfig(acctest.RandInt()),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccAlicloudOssBucket_importCors(t *testing.T) {
	resourceName := "alicloud_oss_bucket.cors"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketCorsConfig(acctest.RandInt()),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccAlicloudOssBucket_importWebsite(t *testing.T) {
	resourceName := "alicloud_oss_bucket.website"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketWebsiteConfig(acctest.RandInt()),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}
func TestAccAlicloudOssBucket_importLogging(t *testing.T) {
	resourceName := "alicloud_oss_bucket.logging"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketLoggingConfig(acctest.RandInt()),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccAlicloudOssBucket_importReferer(t *testing.T) {
	resourceName := "alicloud_oss_bucket.referer"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketRefererConfig(acctest.RandInt()),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}
func TestAccAlicloudOssBucket_importLifecycle(t *testing.T) {
	resourceName := "alicloud_oss_bucket.lifecycle"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketLifecycleConfig(acctest.RandInt()),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

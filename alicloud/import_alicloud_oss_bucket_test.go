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
			resource.TestStep{
				Config: testAccAlicloudOssBucketBasicConfig(acctest.RandInt()),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
			resource.TestStep{
				Config: testAccAlicloudOssBucketCorsConfig(acctest.RandInt()),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
			resource.TestStep{
				Config: testAccAlicloudOssBucketWebsiteConfig(acctest.RandInt()),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
			resource.TestStep{
				Config: testAccAlicloudOssBucketLoggingConfig(acctest.RandInt()),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
			resource.TestStep{
				Config: testAccAlicloudOssBucketRefererConfig(acctest.RandInt()),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
			resource.TestStep{
				Config: testAccAlicloudOssBucketLifecycleConfig(acctest.RandInt()),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

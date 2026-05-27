package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAliCloudOssBucketReplicationBasic(t *testing.T) {
	resourceId := "alicloud_oss_bucket_replication.test"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-bucket-replication-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ossBucketReplicationCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclOssBucketReplicationBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name+"-1"),
					resource.TestCheckResourceAttr(resourceId, "action", "PUT,DELETE"),
					resource.TestCheckResourceAttr(resourceId, "destination.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.bucket", name+"-2"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.location", "oss-"+os.Getenv("ALICLOUD_REGION")),
					resource.TestCheckResourceAttr(resourceId, "destination.0.transfer_type", ""),
					resource.TestCheckResourceAttr(resourceId, "historical_object_replication", "enabled"),
					resource.TestCheckResourceAttrSet(resourceId, "rule_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttr(resourceId, "sync_role", name+"-ramrole"),
					resource.TestCheckResourceAttr(resourceId, "encryption_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "encryption_configuration.0.replica_kms_key_id"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.0.sse_kms_encrypted_objects.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.0.sse_kms_encrypted_objects.0.status", "Enabled"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.#", "3"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.0", "1230"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.1", "456"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.2", "789"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclOssBucketReplicationUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name+"-1"),
					resource.TestCheckResourceAttr(resourceId, "action", "PUT,DELETE"),
					resource.TestCheckResourceAttr(resourceId, "destination.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.bucket", name+"-2"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.location", "oss-"+os.Getenv("ALICLOUD_REGION")),
					resource.TestCheckResourceAttr(resourceId, "destination.0.transfer_type", ""),
					resource.TestCheckResourceAttr(resourceId, "historical_object_replication", "enabled"),
					resource.TestCheckResourceAttrSet(resourceId, "rule_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttr(resourceId, "sync_role", name+"-ramrole"),
					resource.TestCheckResourceAttr(resourceId, "encryption_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "encryption_configuration.0.replica_kms_key_id"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.0.sse_kms_encrypted_objects.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.0.sse_kms_encrypted_objects.0.status", "Enabled"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.#", "3"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.0", "1230"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.1", "456"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.2", "789"),
					resource.TestCheckResourceAttr(resourceId, "progress.#", "1"),
				),
			},
		},
	})
}

// Test case for Cross Region Replication
// test region: cn-hangzhou <--> cn-shanghai (Hard-coded in the HCL)
func TestAccAliCloudOssBucketReplicationCrossRegionReplication(t *testing.T) {
	resourceId := "alicloud_oss_bucket_replication.test"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-bucket-replication-%d", rand)
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
		"alicloudshanghai": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      ossBucketReplicationCheckDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: hclOssBucketReplicationCrossRegionReplication(name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name+"-1"),
					resource.TestCheckResourceAttr(resourceId, "action", "PUT,DELETE"),
					resource.TestCheckResourceAttr(resourceId, "destination.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.bucket", name+"-2"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.location", "oss-cn-shanghai"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.transfer_type", ""),
					resource.TestCheckResourceAttr(resourceId, "historical_object_replication", "enabled"),
					resource.TestCheckResourceAttrSet(resourceId, "rule_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttr(resourceId, "sync_role", name+"-ramrole"),
					resource.TestCheckResourceAttr(resourceId, "encryption_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "encryption_configuration.0.replica_kms_key_id"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.0.sse_kms_encrypted_objects.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.0.sse_kms_encrypted_objects.0.status", "Enabled"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.0", "prefix1/"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.1", "prefix2/"),
					resource.TestCheckResourceAttr(resourceId, "rtc.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "rtc.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceId, "rtc.0.status", "enabling"),
				),
			},
			{
				Config: hclOssBucketReplicationCrossRegionReplication(name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name+"-1"),
					resource.TestCheckResourceAttr(resourceId, "action", "PUT,DELETE"),
					resource.TestCheckResourceAttr(resourceId, "destination.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.bucket", name+"-2"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.location", "oss-cn-shanghai"),
					resource.TestCheckResourceAttr(resourceId, "destination.0.transfer_type", ""),
					resource.TestCheckResourceAttr(resourceId, "historical_object_replication", "enabled"),
					resource.TestCheckResourceAttrSet(resourceId, "rule_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttr(resourceId, "sync_role", name+"-ramrole"),
					resource.TestCheckResourceAttr(resourceId, "encryption_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "encryption_configuration.0.replica_kms_key_id"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.0.sse_kms_encrypted_objects.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "source_selection_criteria.0.sse_kms_encrypted_objects.0.status", "Enabled"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.0", "prefix1/"),
					resource.TestCheckResourceAttr(resourceId, "prefix_set.0.prefixes.1", "prefix2/"),
					resource.TestCheckResourceAttr(resourceId, "rtc.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "rtc.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceId, "rtc.0.status", ""),
				),
			},
		},
	})
}

func ossBucketReplicationCheckDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, p := range *providers {
			if p.Meta() == nil {
				continue
			}
			if err := ossBucketReplicationCheckDestroyWithProvider(s, p); err != nil {
				return err
			}
		}
		return nil
	}
}

func ossBucketReplicationCheckDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	ossClient := OssService{provider.Meta().(*connectivity.AliyunClient)}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket_replication" {
			continue
		}
		_, err := ossClient.DescribeOssBucketReplication(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func ossBucketReplicationCheckDestroy(s *terraform.State) error {
	clients := testAccProvider.Meta().(*connectivity.AliyunClient)
	ossClient := OssService{clients}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket_replication" {
			continue
		}

		_, err := ossClient.DescribeOssBucketReplication(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}
	return nil
}

func hclOssBucketReplicationTemplate(name string) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "bucket_src" {
  bucket = "%[1]s-1"
}

resource "alicloud_oss_bucket" "bucket_dest" {
  bucket = "%[1]s-2"
}

resource "alicloud_ram_role" "test" {
  name        = "%[1]s-ramrole"
  document    = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "oss.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  	EOF
  description = "Terraform AccTest"
  force       = true
}

resource "alicloud_ram_policy" "test" {
  policy_name     = "%[1]s-rampolicy"
  policy_document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"*"
			  ],
			  "Effect": "Allow",
			  "Resource": [
				"*"
			  ]
			}
		  ],
			"Version": "1"
		}
		EOF
  description     = "Terraform AccTest"
  force           = true
}

resource "alicloud_ram_role_policy_attachment" "test" {
  policy_name = alicloud_ram_policy.test.name
  policy_type = alicloud_ram_policy.test.type
  role_name   = alicloud_ram_role.test.name
}

resource "alicloud_kms_key" "test" {
  description            = "Hello KMS"
  pending_window_in_days = "7"
  status                 = "Enabled"
}

locals {
  bucket_src  = alicloud_oss_bucket.bucket_src.id
  bucket_dest = alicloud_oss_bucket.bucket_dest.id
  location    = alicloud_oss_bucket.bucket_dest.location
  role_name   = alicloud_ram_role.test.name
  kms_key_id  = alicloud_kms_key.test.id
}
`, name)
}

func hclOssBucketReplicationBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_oss_bucket_replication" "test" {
  prefix_set {
    prefixes = ["1230", "456", "789"]
  }

  destination {
    bucket        = local.bucket_dest
    location      = local.location
    transfer_type = ""
  }

  sync_role = local.role_name
  encryption_configuration {
    replica_kms_key_id = local.kms_key_id
  }

  source_selection_criteria {
    sse_kms_encrypted_objects {
      status = "Enabled"
    }
  }

  bucket                        = local.bucket_src
  action                        = "PUT,DELETE"
  historical_object_replication = "enabled"
}`, hclOssBucketReplicationTemplate(name))
}

func hclOssBucketReplicationUpdate(name string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_oss_bucket_replication" "test" {
  prefix_set {
    prefixes = ["1230", "456", "789"]
  }

  destination {
    bucket        = local.bucket_dest
    location      = local.location
    transfer_type = ""
  }

  sync_role = local.role_name
  encryption_configuration {
    replica_kms_key_id = local.kms_key_id
  }

  source_selection_criteria {
    sse_kms_encrypted_objects {
      status = "Enabled"
    }
  }

  bucket                        = local.bucket_src
  action                        = "PUT,DELETE"
  historical_object_replication = "enabled"
  progress {}
}`, hclOssBucketReplicationTemplate(name))
}

func hclOssBucketReplicationCrossRegionReplication(name string, rtcEnabled bool) string {
	return fmt.Sprintf(`
provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}

provider "alicloudshanghai" {
  alias  = "sh"
  region = "cn-shanghai"
}

resource "alicloud_oss_bucket" "bucket_src" {
  provider = alicloud.hz
  bucket   = "%[1]s-1"
}

resource "alicloud_oss_bucket" "bucket_dest" {
  provider = alicloudshanghai.sh
  bucket   = "%[1]s-2"
}

resource "alicloud_ram_role" "test" {
  provider    = alicloud.hz
  name        = "%[1]s-ramrole"
  document    = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "oss.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
		EOF
  description = "Terraform AccTest"
  force       = true
}

resource "alicloud_ram_policy" "test" {
  provider        = alicloud.hz
  policy_name     = "%[1]s-rampolicy"
  policy_document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"*"
			  ],
			  "Effect": "Allow",
			  "Resource": [
				"*"
			  ]
			}
		  ],
			"Version": "1"
		}
		EOF
  description     = "Terraform AccTest"
  force           = true
}

resource "alicloud_ram_role_policy_attachment" "test" {
  provider    = alicloud.hz
  policy_name = alicloud_ram_policy.test.policy_name
  policy_type = alicloud_ram_policy.test.type
  role_name   = alicloud_ram_role.test.role_name
}

resource "alicloud_kms_key" "test" {
  provider               = alicloud.hz
  description            = "Hello KMS"
  pending_window_in_days = "7"
  status                 = "Enabled"
}

resource "alicloud_oss_bucket_replication" "test" {
  provider                      = alicloud.hz
  bucket                        = alicloud_oss_bucket.bucket_src.id
  action                        = "PUT,DELETE"
  historical_object_replication = "enabled"

  prefix_set {
    prefixes = ["prefix1/", "prefix2/"]
  }

  destination {
    bucket        = alicloud_oss_bucket.bucket_dest.id
    location      = alicloud_oss_bucket.bucket_dest.location
    transfer_type = ""
  }

  sync_role = alicloud_ram_role.test.role_name
  encryption_configuration {
    replica_kms_key_id = alicloud_kms_key.test.id
  }

  source_selection_criteria {
    sse_kms_encrypted_objects {
      status = "Enabled"
    }
  }
  rtc {
    enabled = %[2]t
  }
}`, name, rtcEnabled)
}

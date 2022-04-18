package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type OssServiceWrapper struct {
	service *OssService
}

func (s *OssServiceWrapper) DescribeOssBucket(id string) (response string, err error) {
	_, err = s.service.DescribeOssBucket(id)
	if err != nil && IsExpectedErrors(err, []string{"AccessDenied"}) {
		return response, WrapErrorf(err, NotFoundMsg, AliyunOssGoSdk)
	}
	return response, err
}

func (s *OssServiceWrapper) DescribeOssBucketReplication(id string) (response string, err error) {
	_, err = s.service.DescribeOssBucketReplication(id)
	if err != nil && IsExpectedErrors(err, []string{"AccessDenied"}) {
		return response, WrapErrorf(err, NotFoundMsg, AliyunOssGoSdk)
	}
	return response, err
}

func TestAccAlicloudOssBucketReplicationBasic(t *testing.T) {
	var v string

	resourceId := "alicloud_oss_bucket_replication.default"
	ra := resourceAttrInit(resourceId, ossBucketReplicationMap)

	serviceFunc := func() interface{} {
		return &OssServiceWrapper{&OssService{testAccProvider.Meta().(*connectivity.AliyunClient)}}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-bucket-replication-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketReplicationDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":                        "${local.bucket_src}",
					"action":                        "PUT,DELETE",
					"historical_object_replication": "enabled",
					"prefix_set": []map[string]interface{}{
						{
							"prefixes": []string{
								"1230",
								"456",
								"789",
							},
						},
					},
					"destination": []map[string]interface{}{
						{
							"bucket":   "${local.bucket_dest}",
							"location": "${local.location}",
						},
					},
					"sync_role": "${local.role_name}",
					"encryption_configuration": []map[string]interface{}{
						{
							"replica_kms_key_id": "${local.kms_key_id}",
						},
					},
					"source_selection_criteria": []map[string]interface{}{
						{
							"sse_kms_encrypted_objects": []map[string]interface{}{
								{
									"status": "Enabled",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                        name + "-t-1",
						"action":                        "PUT,DELETE",
						"destination.#":                 "1",
						"destination.0.bucket":          name + "-t-2",
						"destination.0.location":        "oss-" + os.Getenv("ALICLOUD_REGION"),
						"destination.0.transfer_type":   "",
						"historical_object_replication": "enabled",
						"rule_id":                       CHECKSET,
						"status":                        CHECKSET,
						"sync_role":                     name + "-t-ramrole",
						"encryption_configuration.#":    "1",
						"encryption_configuration.0.replica_kms_key_id":                  CHECKSET,
						"source_selection_criteria.#":                                    "1",
						"source_selection_criteria.0.sse_kms_encrypted_objects.#":        "1",
						"source_selection_criteria.0.sse_kms_encrypted_objects.0.status": "Enabled",
						"prefix_set.#":            "1",
						"prefix_set.0.prefixes.#": "3",
						"prefix_set.0.prefixes.0": "1230",
						"prefix_set.0.prefixes.1": "456",
						"prefix_set.0.prefixes.2": "789",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceOssBucketReplicationDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "%s-t"
	}

	resource "alicloud_oss_bucket" "bucket_src" {
		bucket = "${var.name}-1"
	}

	resource "alicloud_oss_bucket" "bucket_dest" {
		bucket = "${var.name}-2"
	}

	resource "alicloud_ram_role" "role" {
		name     = "${var.name}-ramrole"
		document = <<EOF
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
		description = "this is a test"
		force       = true
	}

	resource "alicloud_ram_policy" "policy" {
		policy_name        = "${var.name}-rampolicy"
		policy_document    = <<EOF
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
		description = "this is a policy test"
		force       = true
	}

	resource "alicloud_ram_role_policy_attachment" "attach" {
		policy_name = alicloud_ram_policy.policy.name
		policy_type = alicloud_ram_policy.policy.type
		role_name   = alicloud_ram_role.role.name
	}

	resource "alicloud_kms_key" "key" {
		description             = "Hello KMS"
		pending_window_in_days  = "7"
		status                  = "Enabled"
	}

	locals {
		bucket_src = alicloud_oss_bucket.bucket_src.id
		bucket_dest = alicloud_oss_bucket.bucket_dest.id
		location = alicloud_oss_bucket.bucket_dest.location
		role_name = alicloud_ram_role.role.name
		kms_key_id = alicloud_kms_key.key.id
	}
`, name)
}

var ossBucketReplicationMap = map[string]string{}

package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsDiskDefaultKmsKeyId() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsDiskDefaultKmsKeyIdCreate,
		Read:   resourceAliCloudEcsDiskDefaultKmsKeyIdRead,
		Update: resourceAliCloudEcsDiskDefaultKmsKeyIdUpdate,
		Delete: resourceAliCloudEcsDiskDefaultKmsKeyIdDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"kms_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The KMS key ID used for ECS disk encryption by default.",
			},
		},
	}
}

func resourceAliCloudEcsDiskDefaultKmsKeyIdCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// Use region as the resource ID since this is a region-scoped setting
	d.SetId(client.RegionId)

	return resourceAliCloudEcsDiskDefaultKmsKeyIdUpdate(d, meta)
}

func resourceAliCloudEcsDiskDefaultKmsKeyIdRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}
	object, err := ecsServiceV2.DescribeEcsDiskDefaultKMSKeyId(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_disk_default_kms_key_id ecsServiceV2.DescribeEcsDiskDefaultKMSKeyId Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if kmsKeyId, ok := object["KMSKeyId"]; ok && kmsKeyId != nil {
		d.Set("kms_key_id", kmsKeyId)
	}
	return nil
}

func resourceAliCloudEcsDiskDefaultKmsKeyIdUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var query map[string]interface{}
	action := "ModifyDiskDefaultKMSKeyId"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["KMSKeyId"] = d.Get("kms_key_id")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
		_, err := client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, nil, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudEcsDiskDefaultKmsKeyIdRead(d, meta)
}

func resourceAliCloudEcsDiskDefaultKmsKeyIdDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var query map[string]interface{}
	action := "ResetDiskDefaultKMSKeyId"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		_, err := client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, nil, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

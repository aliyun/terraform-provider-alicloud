package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsDiskEncryptionByDefault() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsDiskEncryptionByDefaultCreate,
		Read:   resourceAliCloudEcsDiskEncryptionByDefaultRead,
		Update: resourceAliCloudEcsDiskEncryptionByDefaultUpdate,
		Delete: resourceAliCloudEcsDiskEncryptionByDefaultDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable ECS disk encryption by default.",
			},
		},
	}
}

func resourceAliCloudEcsDiskEncryptionByDefaultCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// Use region as the resource ID since this is a region-scoped setting
	d.SetId(client.RegionId)

	enabled := d.Get("enabled").(bool)
	if enabled {
		return resourceAliCloudEcsDiskEncryptionByDefaultEnableEncryption(d, meta)
	} else {
		// If we want to disable encryption, we need to explicitly disable it
		// in case it's currently enabled in the region
		return resourceAliCloudEcsDiskEncryptionByDefaultDisableEncryption(d, meta)
	}
}

func resourceAliCloudEcsDiskEncryptionByDefaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	object, err := ecsServiceV2.DescribeEcsDiskEncryptionByDefaultStatus(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_disk_encryption_by_default ecsServiceV2.DescribeEcsDiskEncryptionByDefaultStatus Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enabled", object["Encrypted"])
	return nil
}

func resourceAliCloudEcsDiskEncryptionByDefaultUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		if enabled {
			return resourceAliCloudEcsDiskEncryptionByDefaultEnableEncryption(d, meta)
		} else {
			return resourceAliCloudEcsDiskEncryptionByDefaultDisableEncryption(d, meta)
		}
	}
	return resourceAliCloudEcsDiskEncryptionByDefaultRead(d, meta)
}

func resourceAliCloudEcsDiskEncryptionByDefaultDelete(d *schema.ResourceData, meta interface{}) error {
	// When deleting the resource, disable encryption by default
	return resourceAliCloudEcsDiskEncryptionByDefaultDisableEncryption(d, meta)
}

func resourceAliCloudEcsDiskEncryptionByDefaultEnableEncryption(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var query map[string]interface{}
	action := "EnableDiskEncryptionByDefault"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		_, err := client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			// If encryption is already enabled, treat it as success
			if IsExpectedErrors(err, []string{"InvalidOperation.DefaultEncryptionAlreadyEnabled"}) {
				return nil
			}
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

	return resourceAliCloudEcsDiskEncryptionByDefaultRead(d, meta)
}

func resourceAliCloudEcsDiskEncryptionByDefaultDisableEncryption(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var query map[string]interface{}
	action := "DisableDiskEncryptionByDefault"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		_, err := client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			// If encryption is already disabled, treat it as success
			if IsExpectedErrors(err, []string{"InvalidOperation.DefaultEncryptionAlreadyDisabled"}) {
				return nil
			}
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

	return resourceAliCloudEcsDiskEncryptionByDefaultRead(d, meta)
}

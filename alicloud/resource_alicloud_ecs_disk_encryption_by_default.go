package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEcsDiskEncryptionByDefaultCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	isEncrypted := false
	if v, ok := d.GetOkExists("encrypted"); ok {
		isEncrypted = v.(bool)

		if isEncrypted {
			action := "EnableDiskEncryptionByDefault"

			wait := incrementalWait(3*time.Second, 0*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_disk_encryption_by_default", action, AlibabaCloudSdkGoERROR)
			}
		} else {
			action := "DisableDiskEncryptionByDefault"

			wait := incrementalWait(3*time.Second, 0*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_disk_encryption_by_default", action, AlibabaCloudSdkGoERROR)
			}
		}

	}

	d.SetId(fmt.Sprint(request["RegionId"]))

	return resourceAliCloudEcsDiskEncryptionByDefaultRead(d, meta)
}

func resourceAliCloudEcsDiskEncryptionByDefaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsDiskEncryptionByDefault(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_disk_encryption_by_default DescribeEcsDiskEncryptionByDefault Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("encrypted", objectRaw["Encrypted"])

	return nil
}

func resourceAliCloudEcsDiskEncryptionByDefaultUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	ecsServiceV2 := EcsServiceV2{client}
	objectRaw, _ := ecsServiceV2.DescribeEcsDiskEncryptionByDefault(d.Id())

	initedEncrypted := false
	if _, ok := d.GetOkExists("encrypted"); ok && d.IsNewResource() {
		initedEncrypted = true
	}
	if initedEncrypted || d.HasChange("encrypted") {
		var err error
		target := d.Get("encrypted").(bool)

		currentStatus, err := jsonpath.Get("Encrypted", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "Encrypted", objectRaw)
		}
		if formatBool(currentStatus) != target {
			if target == true {
				action := "EnableDiskEncryptionByDefault"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["RegionId"] = d.Id()
				wait := incrementalWait(3*time.Second, 0*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == false {
				action := "DisableDiskEncryptionByDefault"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["RegionId"] = d.Id()
				wait := incrementalWait(3*time.Second, 0*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	return resourceAliCloudEcsDiskEncryptionByDefaultRead(d, meta)
}

func resourceAliCloudEcsDiskEncryptionByDefaultDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Disk Encryption By Default. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

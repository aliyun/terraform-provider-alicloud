// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudResourceManagerResourceDirectorySharing() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerResourceDirectorySharingCreate,
		Read:   resourceAliCloudResourceManagerResourceDirectorySharingRead,
		Delete: resourceAliCloudResourceManagerResourceDirectorySharingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enable_sharing_with_rd": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudResourceManagerResourceDirectorySharingCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "EnableSharingWithResourceDirectory"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)
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
		// EnableSharingWithResourceDirectory is idempotent at the account level: a 409
		// "AlreadyEnabled" means the desired state is already in place, which Terraform
		// should treat as a successful create.
		if !IsExpectedErrors(err, []string{"AlreadyEnabled"}) {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_resource_directory_sharing", action, AlibabaCloudSdkGoERROR)
		}
	}

	accountId, err := client.AccountId()
	if err != nil {
		return WrapError(err)
	}
	d.SetId(accountId)

	return resourceAliCloudResourceManagerResourceDirectorySharingRead(d, meta)
}

func resourceAliCloudResourceManagerResourceDirectorySharingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerResourceDirectorySharing(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_directory_sharing DescribeResourceManagerResourceDirectorySharing Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enable_sharing_with_rd", objectRaw["EnableSharingWithRd"])

	return nil
}

func resourceAliCloudResourceManagerResourceDirectorySharingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Resource Directory Sharing. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

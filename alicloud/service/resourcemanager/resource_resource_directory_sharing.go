// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package resourcemanager

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceResourceDirectorySharing() *schema.Resource {
	return &schema.Resource{
		Create: ResourceResourceDirectorySharingCreate,
		Read:   ResourceResourceDirectorySharingRead,
		Delete: ResourceResourceDirectorySharingDelete,
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

func ResourceResourceDirectorySharingCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "EnableSharingWithResourceDirectory"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		// EnableSharingWithResourceDirectory is idempotent at the account level: a 409
		// "AlreadyEnabled" means the desired state is already in place, which Terraform
		// should treat as a successful create.
		if !errs.IsExpectedErrors(err, []string{"AlreadyEnabled"}) {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_resource_directory_sharing", action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	accountId, err := client.AccountId()
	if err != nil {
		return errs.WrapError(err)
	}
	d.SetId(accountId)

	return ResourceResourceDirectorySharingRead(d, meta)
}

func ResourceResourceDirectorySharingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerResourceDirectorySharing(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_directory_sharing DescribeResourceManagerResourceDirectorySharing Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("enable_sharing_with_rd", objectRaw["EnableSharingWithRd"])

	return nil
}

func ResourceResourceDirectorySharingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Resource Directory Sharing. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

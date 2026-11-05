package dbfs

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create:             resourceAlicloudDbfsServiceLinkedRoleCreate,
		Read:               resourceAlicloudDbfsServiceLinkedRoleRead,
		Delete:             resourceAlicloudDbfsServiceLinkedRoleDelete,
		DeprecationMessage: "This resource has been deprecated since v1.279.0 and will be removed in the future. See: https://help.aliyun.com/en/dbfs/",
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AliyunServiceRoleForDbfs",
				}, false),
			},
			"status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
func resourceAlicloudDbfsServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateServiceLinkedRole"
	request := map[string]interface{}{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	request["ClientToken"] = helper.BuildClientToken("CreateServiceLinkedRole")
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, request, true)
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
	d.SetId(d.Get("product_name").(string))
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityAlreadyExist.Role"}) {
			return resourceAlicloudDbfsServiceLinkedRoleRead(d, meta)
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_dbfs_service_linked_role", action, errs.AlibabaCloudSdkGoERROR)
	}
	return resourceAlicloudDbfsServiceLinkedRoleRead(d, meta)
}

func resourceAlicloudDbfsServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := Service{client}
	object, err := dbfsService.DescribeDbfsServiceLinkedRole(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dbfs_event_bus dbfsService.DescribeDbfsServiceLinkedRole Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("product_name", d.Id())
	d.Set("status", object["DbfsLinkedRole"])
	return nil
}

func resourceAlicloudDbfsServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServiceLinkedRole"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"RoleName": d.Id(),
	}
	wait := helper.IncrementalWait(3*time.Second, 0*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, request, nil, false)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

package cloudsso

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudCloudSsoGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudSsoGroupCreate,
		Read:   resourceAlicloudCloudSsoGroupRead,
		Update: resourceAlicloudCloudSsoGroupUpdate,
		Delete: resourceAlicloudCloudSsoGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-.]{1,128}$`), "The name of the resource. The name must be 1 to 128 characters in length and  can contain letters, digits, periods (.), underscores (_), and hyphens (-)."),
			},
		},
	}
}

func resourceAlicloudCloudSsoGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["DirectoryId"] = d.Get("directory_id")
	request["GroupName"] = d.Get("group_name")
	var response map[string]interface{}
	action := "CreateGroup"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		resp, err := client.RpcPost("cloudsso", "2021-05-15", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_sso_group", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseGroup := response["Group"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["DirectoryId"], ":", responseGroup["GroupId"]))

	return resourceAlicloudCloudSsoGroupRead(d, meta)
}
func resourceAlicloudCloudSsoGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	object, err := cloudssoService.DescribeCloudSsoGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_group cloudssoService.DescribeCloudSsoGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("directory_id", parts[0])
	d.Set("group_id", parts[1])
	d.Set("description", object["Description"])
	d.Set("group_name", object["GroupName"])
	return nil
}
func resourceAlicloudCloudSsoGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"GroupId":     parts[1],
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["NewDescription"] = v
		}
	}
	if d.HasChange("group_name") {
		update = true
		if v, ok := d.GetOk("group_name"); ok {
			request["NewGroupName"] = v
		}
	}
	if update {
		action := "UpdateGroup"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			resp, err := client.RpcPost("cloudsso", "2021-05-15", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCloudSsoGroupRead(d, meta)
}
func resourceAlicloudCloudSsoGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"GroupId":     parts[1],
	}

	action := "DeleteGroup"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		resp, err := client.RpcPost("cloudsso", "2021-05-15", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"DeletionConflict.Group.AccessAssigment"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Group"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

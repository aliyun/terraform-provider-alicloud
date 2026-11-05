package ons

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceOnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceOnsGroupCreate,
		Read:   resourceOnsGroupRead,
		Update: resourceOnsGroupUpdate,
		Delete: resourceOnsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validateOnsGroupId,
				ConflictsWith: []string{"group_id"},
			},
			"group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validateOnsGroupId,
				Deprecated:    "Field 'group_id' has been deprecated from version 1.98.0. Use 'group_name' instead.",
				ConflictsWith: []string{"group_name"},
			},
			"group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"http", "tcp"}, false),
				Default:      "tcp",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"read_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceOnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "OnsGroupCreate"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("group_name"); ok {
		request["GroupId"] = v
	} else if v, ok := d.GetOk("group_id"); ok {
		request["GroupId"] = v
	} else {
		return errs.WrapError(errs.Error(`[ERROR] Argument "group_id" or "group_name" must be set one!`))
	}

	if v, ok := d.GetOk("group_type"); ok {
		request["GroupType"] = v
	}

	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Throttling.User"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ons_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["GroupId"]))

	return resourceOnsGroupUpdate(d, meta)
}
func resourceOnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := Service{client: client}
	object, err := onsService.DescribeOnsGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ons_group onsService.DescribeOnsGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("group_name", parts[1])
	d.Set("group_id", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("group_type", object["GroupType"])
	d.Set("remark", object["Remark"])
	d.Set("tags", tags.TagsToMap(object["Tags"].(map[string]interface{})["Tag"]))
	return nil
}
func resourceOnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := Service{client: client}
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := onsService.SetResourceTags(d, "GROUP"); err != nil {
			return errs.WrapError(err)
		}
	}
	if d.HasChange("read_enable") {
		request := map[string]interface{}{
			"GroupId":    parts[1],
			"InstanceId": parts[0],
		}
		request["ReadEnable"] = d.Get("read_enable")
		action := "OnsGroupConsumerUpdate"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)
	return resourceOnsGroupRead(d, meta)
}
func resourceOnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "OnsGroupDelete"
	var response map[string]interface{}
	request := map[string]interface{}{
		"GroupId":    parts[1],
		"InstanceId": parts[0],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Throttling.User"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

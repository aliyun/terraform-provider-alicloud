package ons

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceOnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOnsInstanceCreate,
		Read:   resourceOnsInstanceRead,
		Update: resourceOnsInstanceUpdate,
		Delete: resourceOnsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringLenBetween(3, 64),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from version 1.97.0. Use 'instance_name' instead.",
				ConflictsWith: []string{"instance_name"},
				ValidateFunc:  validation.StringLenBetween(3, 64),
			},
			"instance_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"release_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceOnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "OnsInstanceCreate"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["InstanceName"] = v
	} else {
		return errs.WrapError(errs.Error(`[ERROR] Argument "name" or "instance_name" must be set one!`))
	}

	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 10*time.Second)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ons_instance", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))

	return resourceOnsInstanceUpdate(d, meta)
}
func resourceOnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := Service{client: client}
	object, err := onsService.DescribeOnsInstance(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ons_instance onsService.DescribeOnsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("instance_name", object["InstanceName"])
	d.Set("name", object["InstanceName"])
	instance_typeVal, err := convert.FormatInt(object["InstanceType"])
	if err != nil {
		return err
	}
	d.Set("instance_type", instance_typeVal)
	releaseTime, err := convert.FormatInt(object["ReleaseTime"])
	if err != nil {
		return err
	}
	d.Set("release_time", fmt.Sprint(releaseTime))
	d.Set("remark", object["Remark"])
	statusVal, err := convert.FormatInt(object["InstanceStatus"])
	if err != nil {
		return err
	}
	d.Set("status", statusVal)
	instance_statusVal, err := convert.FormatInt(object["InstanceStatus"])
	if err != nil {
		return err
	}
	d.Set("instance_status", instance_statusVal)

	listTagResourcesObject, err := onsService.ListTagResources(d.Id(), "INSTANCE")
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("tags", tags.TagsToMap(listTagResourcesObject))
	return nil
}
func resourceOnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := Service{client: client}
	var err error
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := onsService.SetResourceTags(d, "INSTANCE"); err != nil {
			return errs.WrapError(err)
		}
	}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		request["InstanceName"] = d.Get("instance_name")
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["InstanceName"] = d.Get("name")
	}
	if !d.IsNewResource() && d.HasChange("remark") {
		update = true
		request["Remark"] = d.Get("remark")
	}
	if update {
		action := "OnsInstanceUpdate"
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
	return resourceOnsInstanceRead(d, meta)
}
func resourceOnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "OnsInstanceDelete"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 10*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EMPTY", "Throttling.User"}) {
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
	return nil
}

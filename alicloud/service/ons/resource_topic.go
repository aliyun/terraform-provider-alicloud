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

func ResourceOnsTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceOnsTopicCreate,
		Read:   resourceOnsTopicRead,
		Update: resourceOnsTopicUpdate,
		Delete: resourceOnsTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"message_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 4, 5}),
			},
			"perm": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 4, 6}),
				Computed:     true,
				Deprecated:   "Attribute perm has been deprecated and suggest removing it from your template.",
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"tags": tags.TagsSchema(),
			"topic_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"topic"},
				ValidateFunc:  validation.StringLenBetween(1, 64),
			},
			"topic": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'topic' has been deprecated from version 1.97.0. Use 'topic_name' instead.",
				ConflictsWith: []string{"topic_name"},
				ValidateFunc:  validation.StringLenBetween(1, 64),
			},
		},
	}
}

func resourceOnsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "OnsTopicCreate"
	request := make(map[string]interface{})
	var err error
	request["InstanceId"] = d.Get("instance_id")
	request["MessageType"] = d.Get("message_type")
	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request["Topic"] = v
	} else if v, ok := d.GetOk("topic"); ok {
		request["Topic"] = v
	} else {
		return errs.WrapError(errs.Error(`[ERROR] Argument "topic" or "topic_name" must be set one!`))
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ons_topic", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Topic"]))

	onsService := Service{client: client}
	stateConf := helper.BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, onsService.OceanOnsTopicStateRefreshFunc(d.Id(), "$.ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceOnsTopicUpdate(d, meta)
}
func resourceOnsTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := Service{client: client}
	object, err := onsService.DescribeOnsTopic(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ons_topic onsService.DescribeOnsTopic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("topic_name", parts[1])
	d.Set("topic", parts[1])
	message_typeVal, err := convert.FormatInt(object["MessageType"])
	if err != nil {
		return err
	}
	d.Set("message_type", message_typeVal)
	d.Set("remark", object["Remark"])
	d.Set("tags", tags.TagsToMap(object["Tags"].(map[string]interface{})["Tag"]))

	onsTopicStatusObject, err := onsService.OnsTopicStatus(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	permVal, err := convert.FormatInt(onsTopicStatusObject["Perm"])
	if err != nil {
		return err
	}
	d.Set("perm", permVal)
	return nil
}
func resourceOnsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := Service{client: client}
	if d.HasChange("tags") {
		if err := onsService.SetResourceTags(d, "TOPIC"); err != nil {
			return errs.WrapError(err)
		}
	}
	return resourceOnsTopicRead(d, meta)
}
func resourceOnsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "OnsTopicDelete"
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"Topic":      parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 10*time.Second)
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

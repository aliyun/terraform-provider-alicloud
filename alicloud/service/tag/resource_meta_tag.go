package tag

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceMetaTag() *schema.Resource {
	return &schema.Resource{
		Create: ResourceMetaTagCreate,
		Read:   ResourceMetaTagRead,
		Delete: ResourceMetaTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringLenBetween(1, 128),
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validators.StringLenBetween(1, 128),
				},
			},
		},
	}
}

func ResourceMetaTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTags"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var err error
	tagKeyValueParamLists := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("key"); ok {
		tagKeyValueParamList := make(map[string]interface{})
		tagKeyValueParamList["Key"] = v.(string)
		if v, ok := d.GetOk("values"); ok {
			tagValueParamListsMaps := make([]map[string]interface{}, 0)
			strings := v.([]interface{})
			for _, w := range strings {
				tagValueParamListsMap := make(map[string]interface{})
				tagValueParamListsMap["Value"] = w.(string)
				tagValueParamListsMaps = append(tagValueParamListsMaps, tagValueParamListsMap)
			}
			tagKeyValueParamList["TagValueParamList"] = tagValueParamListsMaps
		}
		tagKeyValueParamLists = append(tagKeyValueParamLists, tagKeyValueParamList)
	}
	request["TagKeyValueParamList"] = tagKeyValueParamLists

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_tag_meta_tag", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v:%v", client.RegionId, d.Get("key")))
	return ResourceMetaTagRead(d, meta)
}
func ResourceMetaTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagService := NewTagService(client)
	parts, err := helper.ParseResourceIdN(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	o, err := tagService.DescribeTagValue(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("key", parts[1])
	d.Set("values", o)
	return nil
}

func ResourceMetaTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTag"
	var response map[string]interface{}
	var err error
	parts, err := helper.ParseResourceIdN(d.Id(), 2)
	w, _ := d.GetOk("values")
	strings := w.([]interface{})
	for _, v := range strings {
		request := map[string]interface{}{
			"Key":   parts[1],
			"Value": v.(string),
		}

		request["RegionId"] = client.RegionId
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
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

	requestAll := map[string]interface{}{
		"Key": parts[1],
	}

	requestAll["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, requestAll, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, requestAll)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

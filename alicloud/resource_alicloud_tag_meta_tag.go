package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudTagMetaTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudTagMetaTagCreate,
		Read:   resourceAlicloudTagMetaTagRead,
		Delete: resourceAlicloudTagMetaTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(1, 128),
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringLenBetween(1, 128),
				},
			},
		},
	}
}

func resourceAlicloudTagMetaTagCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_tag_meta_tag", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v:%v", client.RegionId, d.Get("key")))
	return resourceAlicloudTagMetaTagRead(d, meta)
}
func resourceAlicloudTagMetaTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagService := TagService{client}
	parts, err := ParseResourceIdN(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	o, err := tagService.DescribeTagValue(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("key", parts[1])
	d.Set("values", o)
	return nil
}

func resourceAlicloudTagMetaTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTag"
	var response map[string]interface{}
	var err error
	parts, err := ParseResourceIdN(d.Id(), 2)
	w, _ := d.GetOk("values")
	strings := w.([]interface{})
	for _, v := range strings {
		request := map[string]interface{}{
			"Key":   parts[1],
			"Value": v.(string),
		}

		request["RegionId"] = client.RegionId
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	requestAll := map[string]interface{}{
		"Key": parts[1],
	}

	requestAll["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, requestAll, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, requestAll)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

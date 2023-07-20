package alicloud

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"value_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
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
	conn, err := client.NewTagClient()
	if err != nil {
		return WrapError(err)
	}
	tagKeyValueParamLists := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("key_name"); ok {
		tagKeyValueParamList := make(map[string]interface{})
		tagKeyValueParamList["Key"] = v.(string)
		if w, ok := d.GetOk("value_name"); ok {
			tagValueParamListsMaps := make([]map[string]interface{}, 0)
			tagValueParamListsMap := make(map[string]interface{})
			tagValueParamListsMap["Value"] = w
			tagValueParamListsMaps = append(tagValueParamListsMaps, tagValueParamListsMap)
			tagKeyValueParamList["TagValueParamList"] = tagValueParamListsMaps
		}
		tagKeyValueParamLists = append(tagKeyValueParamLists, tagKeyValueParamList)
	}
	request["TagKeyValueParamList"] = tagKeyValueParamLists

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	d.SetId(fmt.Sprintf("%v:%v", d.Get("key_name"), d.Get("value_name")))
	return resourceAlicloudTagMetaTagRead(d, meta)
}
func resourceAlicloudTagMetaTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagService := TagService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := tagService.DescribeTagValue(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("key_name", parts[0])
	d.Set("value_name", object)
	return nil
}

func resourceAlicloudTagMetaTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTag"
	var response map[string]interface{}
	conn, err := client.NewTagClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	request := map[string]interface{}{
		"Key":   parts[0],
		"Value": parts[1],
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	return nil
}

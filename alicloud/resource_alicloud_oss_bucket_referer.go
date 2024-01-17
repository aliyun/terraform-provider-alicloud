// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketReferer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketRefererCreate,
		Read:   resourceAliCloudOssBucketRefererRead,
		Update: resourceAliCloudOssBucketRefererUpdate,
		Delete: resourceAliCloudOssBucketRefererDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allow_empty_referer": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"allow_truncate_query_string": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"referer_blacklist": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"referer": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"referer_list": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"referer": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"truncate_path": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudOssBucketRefererCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?referer")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	request["AllowEmptyReferer"] = d.Get("allow_empty_referer")
	if v, ok := d.GetOkExists("allow_truncate_query_string"); ok {
		request["AllowTruncateQueryString"] = v
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("referer_list"); v != nil {
		nodeNative, _ := jsonpath.Get("$[0].referer", v)
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["Referer"] = nodeNative.(*schema.Set).List()
		}
		request["RefererList"] = objectDataLocalMap
	}

	if v, ok := d.GetOkExists("truncate_path"); ok {
		request["TruncatePath"] = v
	}
	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("referer_blacklist"); !IsNil(v) {
		nodeNative1, _ := jsonpath.Get("$[0].referer", v)
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap1["Referer"] = nodeNative1.(*schema.Set).List()
		}
		request["RefererBlacklist"] = objectDataLocalMap1
	}

	body["RefererConfiguration"] = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.Execute(genXmlParam("PutBucketReferer", "PUT", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_referer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketRefererRead(d, meta)
}

func resourceAliCloudOssBucketRefererRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketReferer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_referer DescribeOssBucketReferer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("allow_empty_referer", convertStringToBool(objectRaw["AllowEmptyReferer"].(string)))
	d.Set("allow_truncate_query_string", convertStringToBool(objectRaw["AllowTruncateQueryString"].(string)))
	d.Set("truncate_path", convertStringToBool(objectRaw["TruncatePath"].(string)))

	refererBlacklistMaps := make([]map[string]interface{}, 0)
	refererBlacklistMap := make(map[string]interface{})
	referer2RawObj, _ := jsonpath.Get("$.RefererBlacklist.Referer", objectRaw)
	referer2Raw := make([]interface{}, 0)
	if referer2RawObj != nil {
		referer2Raw = referer2RawObj.([]interface{})
	}
	if len(referer2Raw) > 0 {

		refererBlacklistMap["referer"] = referer2Raw
		refererBlacklistMaps = append(refererBlacklistMaps, refererBlacklistMap)
	}
	d.Set("referer_blacklist", refererBlacklistMaps)
	refererListMaps := make([]map[string]interface{}, 0)
	refererListMap := make(map[string]interface{})
	referer3RawObj, _ := jsonpath.Get("$.RefererList.Referer", objectRaw)
	referer3Raw := make([]interface{}, 0)
	if referer3RawObj != nil {
		referer3Raw = referer3RawObj.([]interface{})
	}
	if len(referer3Raw) > 0 {

		refererListMap["referer"] = referer3Raw
		refererListMaps = append(refererListMaps, refererListMap)
	}
	d.Set("referer_list", refererListMaps)

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketRefererUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	var hostMap map[string]*string
	update := false
	action := fmt.Sprintf("/?referer")
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap = make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())
	if d.HasChange("allow_empty_referer") {
		update = true
	}
	request["AllowEmptyReferer"] = d.Get("allow_empty_referer")
	if d.HasChange("allow_truncate_query_string") {
		update = true
	}
	request["AllowTruncateQueryString"] = d.Get("allow_truncate_query_string")
	if d.HasChange("referer_list") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("referer_list"); v != nil {
		nodeNative, _ := jsonpath.Get("$[0].referer", d.Get("referer_list"))
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["Referer"] = nodeNative.(*schema.Set).List()
		}
		request["RefererList"] = objectDataLocalMap
	}

	if d.HasChange("truncate_path") {
		update = true
	}
	request["TruncatePath"] = d.Get("truncate_path")
	if d.HasChange("referer_blacklist") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("referer_blacklist"); v != nil {
		nodeNative1, _ := jsonpath.Get("$[0].referer", d.Get("referer_blacklist"))
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap1["Referer"] = nodeNative1.(*schema.Set).List()
		}
		request["RefererBlacklist"] = objectDataLocalMap1
	}

	body["RefererConfiguration"] = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.Execute(genXmlParam("PutBucketReferer", "PUT", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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

	return resourceAliCloudOssBucketRefererRead(d, meta)
}

func resourceAliCloudOssBucketRefererDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?referer")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Id())

	request["AllowEmptyReferer"] = "true"
	request["RefererList"] = " "
	body["RefererConfiguration"] = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.Execute(genXmlParam("PutBucketReferer", "PUT", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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

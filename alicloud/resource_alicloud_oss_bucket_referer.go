// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
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
				Required: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"referer_list": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"referer": {
							Type:     schema.TypeList,
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
	hostMap["bucket"] = StringPointer(d.Get("bucket_name").(string))

	request["AllowEmptyReferer"] = d.Get("allow_empty_referer")
	request["AllowTruncateQueryString"] = d.Get("allow_truncate_query_string")
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("referer_list"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].referer", v)
		if nodeNative != "" {
			objectDataLocalMap["Referer"] = nodeNative
		}
		request["RefererList"] = objectDataLocalMap
	}

	if v, ok := d.GetOkExists("truncate_path"); ok {
		request["TruncatePath"] = v
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.Execute(genRoaParam("PutBucketReferer", "PUT", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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

	d.SetId(fmt.Sprint(request["bucket"]))

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

	d.Set("allow_empty_referer", objectRaw["AllowEmptyReferer"])
	d.Set("allow_truncate_query_string", objectRaw["AllowTruncateQueryString"])
	d.Set("truncate_path", objectRaw["TruncatePath"])

	refererListMaps := make([]map[string]interface{}, 0)
	refererListMap := make(map[string]interface{})
	referer1RawObj, _ := jsonpath.Get("$.RefererList.Referer", objectRaw)
	referer1Raw := make(map[string]interface{})
	if referer1RawObj != nil {
		referer1Raw = referer1RawObj.(map[string]interface{})
	}
	if len(referer1Raw) > 0 {
		refererListMap["referer"] = referer1Raw
		refererListMaps = append(refererListMaps, refererListMap)
	}
	d.Set("referer_list", refererListMaps)

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
	if v := d.Get("referer_list"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].referer", v)
		if nodeNative != "" {
			objectDataLocalMap["Referer"] = nodeNative
		}
		request["RefererList"] = objectDataLocalMap
	}

	if d.HasChange("truncate_path") {
		update = true
	}
	request["TruncatePath"] = d.Get("truncate_path")
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.Execute(genRoaParam("PutBucketReferer", "PUT", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Bucket Referer. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

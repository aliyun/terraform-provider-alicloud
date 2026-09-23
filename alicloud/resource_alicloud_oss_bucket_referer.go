package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"referer_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})

	objectDataLocalMap["AllowEmptyReferer"] = d.Get("allow_empty_referer")
	objectDataLocalMap["AllowTruncateQueryString"] = d.Get("allow_truncate_query_string")
	objectDataLocalMap["TruncatePath"] = d.Get("truncate_path")

	refererList := make(map[string]interface{})
	nodeNative3, _ := jsonpath.Get("$", d.Get("referer_list").(*schema.Set).List())
	if nodeNative3 != nil && nodeNative3 != "" {
		refererList["Referer"] = nodeNative3
	}
	objectDataLocalMap["RefererList"] = refererList

	refererBlacklist := make(map[string]interface{})
	nodeNative4, _ := jsonpath.Get("$", d.Get("referer_blacklist").(*schema.Set).List())
	if nodeNative4 != nil && nodeNative4 != "" {
		refererBlacklist["Referer"] = nodeNative4
	}
	objectDataLocalMap["RefererBlacklist"] = refererBlacklist

	request["RefererConfiguration"] = objectDataLocalMap

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketReferer", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_referer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	ossServiceV2 := OssServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketRefererStateRefreshFunc(d.Id(), "#$.RefererList", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

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

	referer2RawObj, _ := jsonpath.Get("$.RefererBlacklist.Referer", objectRaw)
	referer2Raw := make([]interface{}, 0)
	if referer2RawObj != nil {
		referer2Raw = referer2RawObj.([]interface{})
	}
	if len(referer2Raw) > 0 {
		d.Set("referer_blacklist", referer2Raw)
	}
	referer3RawObj, _ := jsonpath.Get("$.RefererList.Referer", objectRaw)
	referer3Raw := make([]interface{}, 0)
	if referer3RawObj != nil {
		referer3Raw = referer3RawObj.([]interface{})
	}
	if len(referer3Raw) > 0 {
		d.Set("referer_list", referer3Raw)
	}

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketRefererUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/?referer")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())

	objectDataLocalMap := make(map[string]interface{})

	if d.HasChange("allow_empty_referer") {
		update = true
	}
	objectDataLocalMap["AllowEmptyReferer"] = d.Get("allow_empty_referer")
	if d.HasChange("allow_truncate_query_string") {
		update = true
	}
	objectDataLocalMap["AllowTruncateQueryString"] = d.Get("allow_truncate_query_string")
	if d.HasChange("truncate_path") {
		update = true
	}
	objectDataLocalMap["TruncatePath"] = d.Get("truncate_path")
	if d.HasChange("referer_list") {
		update = true
	}

	refererList := make(map[string]interface{})
	nodeNative3, _ := jsonpath.Get("$", d.Get("referer_list").(*schema.Set).List())
	if nodeNative3 != nil && nodeNative3 != "" {
		refererList["Referer"] = nodeNative3
	}
	objectDataLocalMap["RefererList"] = refererList

	if d.HasChange("referer_blacklist") {
		update = true
	}

	refererBlacklist := make(map[string]interface{})
	nodeNative4, _ := jsonpath.Get("$", d.Get("referer_blacklist").(*schema.Set).List())
	if nodeNative4 != nil && nodeNative4 != "" {
		refererBlacklist["Referer"] = nodeNative4
	}
	objectDataLocalMap["RefererBlacklist"] = refererBlacklist

	request["RefererConfiguration"] = objectDataLocalMap

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketReferer", action), query, body, nil, hostMap, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		ossServiceV2 := OssServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketRefererStateRefreshFunc(d.Id(), "#RefererList", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Id())

	objectDataLocalMap := make(map[string]interface{})

	objectDataLocalMap["AllowEmptyReferer"] = d.Get("allow_empty_referer")
	objectDataLocalMap["AllowTruncateQueryString"] = d.Get("allow_truncate_query_string")
	objectDataLocalMap["TruncatePath"] = d.Get("truncate_path")

	refererList := make(map[string]interface{})
	objectDataLocalMap["RefererList"] = refererList

	refererBlacklist := make(map[string]interface{})
	objectDataLocalMap["RefererBlacklist"] = refererBlacklist

	request["RefererConfiguration"] = objectDataLocalMap

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketReferer", action), query, body, nil, hostMap, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package fcv3

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLayerVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayerVersionCreate,
		Read:   resourceLayerVersionRead,
		Update: resourceLayerVersionUpdate,
		Delete: resourceLayerVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"1", "0"}, false),
			},
			"code": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zip_file": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"checksum": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"oss_object_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"oss_bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"code_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compatible_runtime": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"layer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"layer_version_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"license": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"public": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLayerVersionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	layerName := d.Get("layer_name")
	action := fmt.Sprintf("/2023-03-30/layers/%s/versions", layerName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("code"); !helper.IsNil(v) {
		checksum1, _ := jsonpath.Get("$[0].checksum", d.Get("code"))
		if checksum1 != nil && checksum1 != "" {
			objectDataLocalMap["checksum"] = checksum1
		}
		ossBucketName1, _ := jsonpath.Get("$[0].oss_bucket_name", d.Get("code"))
		if ossBucketName1 != nil && ossBucketName1 != "" {
			objectDataLocalMap["ossBucketName"] = ossBucketName1
		}
		ossObjectName1, _ := jsonpath.Get("$[0].oss_object_name", d.Get("code"))
		if ossObjectName1 != nil && ossObjectName1 != "" {
			objectDataLocalMap["ossObjectName"] = ossObjectName1
		}
		zipFile1, _ := jsonpath.Get("$[0].zip_file", d.Get("code"))
		if zipFile1 != nil && zipFile1 != "" {
			objectDataLocalMap["zipFile"] = zipFile1
		}

		request["code"] = objectDataLocalMap
	}

	if v, ok := d.GetOk("compatible_runtime"); ok {
		compatibleRuntimeMapsArray := v.([]interface{})
		request["compatibleRuntime"] = compatibleRuntimeMapsArray
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("license"); ok {
		request["license"] = v
	}
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("FC", "2023-03-30", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_fcv3_layer_version", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", response["layerName"], response["version"]))

	return resourceLayerVersionUpdate(d, meta)
}

func resourceLayerVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeLayerVersion(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv3_layer_version DescribeLayerVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["acl"] != nil {
		d.Set("acl", objectRaw["acl"])
	}
	if objectRaw["codeSize"] != nil {
		d.Set("code_size", objectRaw["codeSize"])
	}
	if objectRaw["createTime"] != nil {
		d.Set("create_time", objectRaw["createTime"])
	}
	if objectRaw["description"] != nil {
		d.Set("description", objectRaw["description"])
	}
	if objectRaw["layerVersionArn"] != nil {
		d.Set("layer_version_arn", objectRaw["layerVersionArn"])
	}
	if objectRaw["license"] != nil {
		d.Set("license", objectRaw["license"])
	}
	if objectRaw["layerName"] != nil {
		d.Set("layer_name", objectRaw["layerName"])
	}
	if objectRaw["version"] != nil {
		d.Set("version", objectRaw["version"])
	}

	compatibleRuntime1Raw := make([]interface{}, 0)
	if objectRaw["compatibleRuntime"] != nil {
		compatibleRuntime1Raw = objectRaw["compatibleRuntime"].([]interface{})
	}

	d.Set("compatible_runtime", compatibleRuntime1Raw)

	return nil
}

func resourceLayerVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	layerName := parts[0]
	action := fmt.Sprintf("/2023-03-30/layers/%s/acl", layerName)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("public") {
		update = true
		if v, ok := d.GetOk("public"); ok {
			query["public"] = helper.StringPointer(v.(string))
		}
	}

	if d.HasChange("acl") {
		update = true
		if v, ok := d.GetOk("acl"); ok {
			query["acl"] = helper.StringPointer(v.(string))
		}
	}

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("FC", "2023-03-30", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceLayerVersionRead(d, meta)
}

func resourceLayerVersionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	layerName := parts[0]
	version := parts[1]
	action := fmt.Sprintf("/2023-03-30/layers/%s/versions/%s", layerName, version)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("FC", "2023-03-30", action, query, nil, nil, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"LayerNotFound", "LayerVersionNotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

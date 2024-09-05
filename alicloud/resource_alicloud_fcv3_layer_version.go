package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudFcv3LayerVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudFcv3LayerVersionCreate,
		Read:   resourceAliCloudFcv3LayerVersionRead,
		Update: resourceAliCloudFcv3LayerVersionUpdate,
		Delete: resourceAliCloudFcv3LayerVersionDelete,
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
				ValidateFunc: StringInSlice([]string{"1", "0"}, false),
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
			"license": {
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

func resourceAliCloudFcv3LayerVersionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	layerName := d.Get("layer_name")
	action := fmt.Sprintf("/2023-03-30/layers/%s/versions", layerName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("code"); !IsNil(v) {
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
		compatibleRuntimeMaps := v.([]interface{})
		request["compatibleRuntime"] = compatibleRuntimeMaps
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("license"); ok {
		request["license"] = v
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fcv3_layer_version", action, AlibabaCloudSdkGoERROR)
	}

	layerNameVar, _ := jsonpath.Get("$.body.layerName", response)
	versionVar, _ := jsonpath.Get("$.body.version", response)
	d.SetId(fmt.Sprintf("%v:%v", layerNameVar, versionVar))

	return resourceAliCloudFcv3LayerVersionUpdate(d, meta)
}

func resourceAliCloudFcv3LayerVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcv3ServiceV2 := Fcv3ServiceV2{client}

	objectRaw, err := fcv3ServiceV2.DescribeFcv3LayerVersion(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv3_layer_version DescribeFcv3LayerVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["acl"] != nil {
		d.Set("acl", objectRaw["acl"])
	}
	if objectRaw["createTime"] != nil {
		d.Set("create_time", objectRaw["createTime"])
	}
	if objectRaw["description"] != nil {
		d.Set("description", objectRaw["description"])
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

	parts := strings.Split(d.Id(), ":")
	d.Set("layer_name", parts[0])
	d.Set("version", parts[1])

	return nil
}

func resourceAliCloudFcv3LayerVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	layerName := parts[0]
	action := fmt.Sprintf("/2023-03-30/layers/%s/acl", layerName)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("public") {
		update = true
	}
	if v, ok := d.GetOk("public"); ok {
		query["public"] = StringPointer(v.(string))
	}

	if d.HasChange("acl") {
		update = true
	}
	if v, ok := d.GetOk("acl"); ok {
		query["acl"] = StringPointer(v.(string))
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	return resourceAliCloudFcv3LayerVersionRead(d, meta)
}

func resourceAliCloudFcv3LayerVersionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	layerName := parts[0]
	version := parts[1]
	action := fmt.Sprintf("/2023-03-30/layers/%s/versions/%s", layerName, version)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
		if IsExpectedErrors(err, []string{"LayerNotFound", "LayerVersionNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

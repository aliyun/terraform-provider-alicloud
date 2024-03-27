package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMSENacosConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMSENacosConfigCreate,
		Read:   resourceAlicloudMSENacosConfigRead,
		Update: resourceAlicloudMSENacosConfigUpdate,
		Delete: resourceAlicloudMSENacosConfigDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en"}, false),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"text", "json", "xml", "yaml", "html", "properties"}, false),
				Default:      "text",
			},
			"desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudMSENacosConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNacosConfig"
	request := make(map[string]interface{})
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}
	if v, ok := d.GetOk("desc"); ok {
		request["Desc"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value,
			})
		}
		request["Tags"] = tags
	}
	request["InstanceId"] = d.Get("cluster_id")
	request["DataId"] = d.Get("name")
	request["NamespaceId"] = d.Get("namespace_id")
	request["Group"] = d.Get("group")
	request["Content"] = d.Get("content")
	request["Type"] = d.Get("type")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mse_engine_namespace", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["NamespaceId"], ":", request["Group"], ":", request["DataId"]))

	return resourceAlicloudMSENacosConfigRead(d, meta)
}
func resourceAlicloudMSENacosConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	object, err := mseService.DescribeMseNacosConfig(d.Id())
	log.Println(object)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_engine_namespace mseService.DescribeMseEngineConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	_ = d.Set("namespace_id", parts[1])
	_ = d.Set("cluster_id", parts[0])
	_ = d.Set("name", object["DataId"])
	_ = d.Set("group", object["Group"])
	_ = d.Set("desc", object["Desc"])
	_ = d.Set("content", object["Content"])
	_ = d.Set("type", object["Type"])
	tagsBytes, err := json.Marshal(object["Tags"])
	_ = d.Set("tags", string(tagsBytes))
	return nil
}
func resourceAlicloudMSENacosConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId":  parts[0],
		"NamespaceId": parts[1],
		"Group":       parts[2],
		"DataId":      parts[3],
	}

	request["Content"] = d.Get("content")
	request["Desc"] = d.Get("desc")
	request["Type"] = d.Get("type")

	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}
	if v, ok := d.GetOk("desc"); ok {
		request["Desc"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value,
			})
		}
		request["Tags"] = tags
	}
	action := "UpdateNacosConfig"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	return resourceAlicloudMSENacosConfigRead(d, meta)
}
func resourceAlicloudMSENacosConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteNacosConfig"
	var response map[string]interface{}
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId":  parts[0],
		"NamespaceId": parts[1],
		"Group":       parts[2],
		"DataId":      parts[3],
	}

	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}
	if v, ok := d.GetOk("desc"); ok {
		request["Desc"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		request["Tags"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

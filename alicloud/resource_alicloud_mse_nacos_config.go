package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMseNacosConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMseNacosConfigCreate,
		Read:   resourceAlicloudMseNacosConfigRead,
		Update: resourceAlicloudMseNacosConfigUpdate,
		Delete: resourceAlicloudMseNacosConfigDelete,
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
				ValidateFunc: StringInSlice([]string{"zh", "en"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"encrypted_data_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"beta_ips": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudMseNacosConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNacosConfig"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}

	request["InstanceId"] = d.Get("instance_id")
	request["DataId"] = d.Get("data_id")
	request["Group"] = d.Get("group")

	if v, ok := d.GetOk("app_name"); ok {
		request["AppName"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		request["Tags"] = v
	}
	if v, ok := d.GetOk("desc"); ok {
		request["Desc"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if v, ok := d.GetOk("content"); ok {
		request["Content"] = v
	}
	if v, ok := d.GetOk("namespace_id"); ok {
		request["NamespaceId"] = v
	}
	if v, ok := d.GetOk("beta_ips"); ok {
		request["BetaIps"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mse_nacos_config", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	var namespaceId = request["NamespaceId"]
	if namespaceId == nil {
		namespaceId = ""
	}

	dataId := request["DataId"].(string)
	group := request["Group"].(string)
	instanceId := request["InstanceId"].(string)

	d.SetId(fmt.Sprint(EscapeColons(instanceId), ":", EscapeColons(namespaceId.(string)), ":", EscapeColons(dataId), ":", EscapeColons(group)))

	return resourceAlicloudMseNacosConfigRead(d, meta)
}

func resourceAlicloudMseNacosConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}

	object, err := mseService.DescribeMseNacosConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_nacos_config mseService.DescribeMseNacosConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceIdWithEscaped(d.Id(), 4)

	err = d.Set("instance_id", parts[0])
	if err != nil {
		return err
	}
	d.Set("namespace_id", parts[1])
	d.Set("app_name", object["AppName"])
	d.Set("data_id", object["DataId"])
	d.Set("group", object["Group"])
	d.Set("type", object["Type"])
	d.Set("tags", object["Tags"])
	d.Set("content", object["Content"])
	d.Set("desc", object["Desc"])
	d.Set("id", object["Id"])

	return nil
}
func resourceAlicloudMseNacosConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	parts, err := ParseResourceIdWithEscaped(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId":  parts[0],
		"NamespaceId": parts[1],
		"DataId":      parts[2],
		"Group":       parts[3],
	}

	if d.HasChanges("content", "app_name", "tags", "desc", "type", "beta_ips") {
		update = true
		request["Content"] = d.Get("content")
		request["AppName"] = d.Get("app_name")
		request["Tags"] = d.Get("tags")
		request["Desc"] = d.Get("desc")
		request["Type"] = d.Get("type")
		request["BetaIps"] = d.Get("beta_ips")

	}

	if update {
		if v, ok := d.GetOk("accept_language"); ok {
			request["AcceptLanguage"] = v
		}
		action := "UpdateNacosConfig"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
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
	}
	return resourceAlicloudMseNacosConfigRead(d, meta)
}
func resourceAlicloudMseNacosConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteNacosConfig"
	var response map[string]interface{}
	var err error
	parts, err := ParseResourceIdWithEscaped(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId":  parts[0],
		"NamespaceId": parts[1],
		"DataId":      parts[2],
		"Group":       parts[3],
	}
	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
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

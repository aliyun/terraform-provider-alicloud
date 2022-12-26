package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAlbAscript() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbAscriptCreate,
		Read:   resourceAlicloudAlbAscriptRead,
		Update: resourceAlicloudAlbAscriptUpdate,
		Delete: resourceAlicloudAlbAscriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"position": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ascript_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"script_content": {
				Required: true,
				Type:     schema.TypeString,
			},
			"enabled": {
				Required: true,
				Type:     schema.TypeBool,
			},
			"ext_attribute_enabled": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeBool,
			},
			"ext_attributes": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute_key": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeString,
						},
						"attribute_value": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"load_balancer_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudAlbAscriptCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("ascript_name"); ok {
		request["AScripts.1.AScriptName"] = v
	}

	if v, ok := d.GetOk("position"); ok {
		request["AScripts.1.Position"] = v
	}

	if v, ok := d.GetOk("enabled"); ok {
		request["AScripts.1.Enabled"] = v
	}

	if v, ok := d.GetOk("script_content"); ok {
		request["AScripts.1.ScriptContent"] = v
	}

	if v, ok := d.GetOk("ext_attributes"); ok {
		for index, value0 := range v.(*schema.Set).List() {
			extAttributes := value0.(map[string]interface{})
			request["AScripts.1.ExtAttributes."+fmt.Sprint(index+1)+".AttributeKey"] = extAttributes["attribute_key"]
			request["AScripts.1.ExtAttributes."+fmt.Sprint(index+1)+".AttributeValue"] = extAttributes["attribute_value"]
		}
	}

	if v, ok := d.GetOk("ext_attribute_enabled"); ok {
		request["AScripts.1.ExtAttributeEnabled"] = v
	}

	if v, ok := d.GetOk("listener_id"); ok {
		request["ListenerId"] = v
	}

	request["ClientToken"] = buildClientToken("CreateAScripts")
	var response map[string]interface{}
	action := "CreateAScripts"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_ascript", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.AScriptIds[0].AScriptId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_alb_ascript")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbAscriptStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudAlbAscriptUpdate(d, meta)
}

func resourceAlicloudAlbAscriptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}

	object, err := albService.DescribeAlbAscript(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_ascript albService.DescribeAlbAscript Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ascript_name", object["AScriptName"])
	d.Set("enabled", object["Enabled"])
	d.Set("ext_attribute_enabled", object["ExtAttributeEnabled"])
	extAttributes19Maps := make([]map[string]interface{}, 0)
	extAttributes19Raw := object["ExtAttributes"]
	for _, value0 := range extAttributes19Raw.([]interface{}) {
		extAttributes19 := value0.(map[string]interface{})
		extAttributes19Map := make(map[string]interface{})
		extAttributes19Map["attribute_key"] = extAttributes19["AttributeKey"]
		extAttributes19Map["attribute_value"] = extAttributes19["AttributeValue"]
		extAttributes19Maps = append(extAttributes19Maps, extAttributes19Map)
	}
	d.Set("ext_attributes", extAttributes19Maps)
	d.Set("listener_id", object["ListenerId"])
	d.Set("load_balancer_id", object["LoadBalancerId"])
	d.Set("position", object["Position"])
	d.Set("script_content", object["ScriptContent"])
	d.Set("status", object["AScriptStatus"])

	return nil
}

func resourceAlicloudAlbAscriptUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	albService := AlbService{client}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"AScripts.1.AScriptId": d.Id(),
	}

	if d.HasChange("ascript_name") {
		update = true
		if v, ok := d.GetOk("ascript_name"); ok {
			request["AScripts.1.AScriptName"] = v
		}
	}

	if d.HasChange("enabled") {
		update = true
		request["AScripts.1.Enabled"] = d.Get("enabled")
	}
	if d.HasChange("ext_attribute_enabled") {
		update = true
		request["AScripts.1.ExtAttributeEnabled"] = d.Get("ext_attribute_enabled")
	}
	if d.HasChange("ext_attributes") {
		update = true
		if v, ok := d.GetOk("ext_attributes"); ok {
			for index, value0 := range v.([]interface{}) {
				extAttributes := value0.(map[string]interface{})
				request["AScripts.1.ExtAttributes."+fmt.Sprint(index+1)+".AttributeKey"] = extAttributes["attribute_key"]
				request["AScripts.1.ExtAttributes."+fmt.Sprint(index+1)+".AttributeValue"] = extAttributes["attribute_value"]
			}
		}
	}
	if d.HasChange("script_content") {
		update = true
		if v, ok := d.GetOk("script_content"); ok {
			request["AScripts.1.ScriptContent"] = v
		}
	}

	if update {
		action := "UpdateAScripts"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbAscriptStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudAlbAscriptRead(d, meta)
}

func resourceAlicloudAlbAscriptDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{

		"AScriptIds.1": d.Id(),
	}

	request["ClientToken"] = buildClientToken("DeleteAScripts")
	action := "DeleteAScripts"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

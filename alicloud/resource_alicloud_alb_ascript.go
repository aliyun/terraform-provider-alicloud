// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudAlbAScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlbAScriptCreate,
		Read:   resourceAliCloudAlbAScriptRead,
		Update: resourceAliCloudAlbAScriptUpdate,
		Delete: resourceAliCloudAlbAScriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ascript_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ext_attribute_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ext_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute_key": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"EsDebug"}, false),
						},
						"attribute_value": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9]{1,127}$"), "The value of the extended attribute"),
						},
					},
				},
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"position": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"RequestHead", "RequestFoot", "ResponseHead"}, false),
			},
			"script_content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAlbAScriptCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAScripts"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ListenerId"] = d.Get("listener_id")
	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOk("position"); ok {
		objectDataLocalMap["Position"] = v
	}

	if v, ok := d.GetOk("enabled"); ok {
		objectDataLocalMap["Enabled"] = v
	}

	if v, ok := d.GetOk("script_content"); ok {
		objectDataLocalMap["ScriptContent"] = v
	}

	if v := d.Get("ext_attributes"); !IsNil(v) {
		if v, ok := d.GetOk("ext_attributes"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["AttributeKey"] = dataLoopTmp["attribute_key"]
				dataLoopMap["AttributeValue"] = dataLoopTmp["attribute_value"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["ExtAttributes"] = localMaps
		}

	}

	if v, ok := d.GetOk("ext_attribute_enabled"); ok {
		objectDataLocalMap["ExtAttributeEnabled"] = v
	}

	if v, ok := d.GetOk("ascript_name"); ok {
		objectDataLocalMap["AScriptName"] = v
	}

	AScriptsMap := make([]interface{}, 0)
	AScriptsMap = append(AScriptsMap, objectDataLocalMap)
	request["AScripts"] = AScriptsMap
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_ascript", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.AScriptIds[0].AScriptId", response)
	d.SetId(fmt.Sprint(id))

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albServiceV2.AlbAScriptStateRefreshFunc(d.Id(), "AScriptStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlbAScriptRead(d, meta)
}

func resourceAliCloudAlbAScriptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	objectRaw, err := albServiceV2.DescribeAlbAScript(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_ascript DescribeAlbAScript Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AScriptName"] != nil {
		d.Set("ascript_name", objectRaw["AScriptName"])
	}
	if objectRaw["Enabled"] != nil {
		d.Set("enabled", objectRaw["Enabled"])
	}
	if objectRaw["ExtAttributeEnabled"] != nil {
		d.Set("ext_attribute_enabled", objectRaw["ExtAttributeEnabled"])
	}
	if objectRaw["ListenerId"] != nil {
		d.Set("listener_id", objectRaw["ListenerId"])
	}
	if objectRaw["Position"] != nil {
		d.Set("position", objectRaw["Position"])
	}
	if objectRaw["ScriptContent"] != nil {
		d.Set("script_content", objectRaw["ScriptContent"])
	}
	if objectRaw["AScriptStatus"] != nil {
		d.Set("status", objectRaw["AScriptStatus"])
	}

	extAttributes1Raw := objectRaw["ExtAttributes"]
	extAttributesMaps := make([]map[string]interface{}, 0)
	if extAttributes1Raw != nil {
		for _, extAttributesChild1Raw := range extAttributes1Raw.([]interface{}) {
			extAttributesMap := make(map[string]interface{})
			extAttributesChild1Raw := extAttributesChild1Raw.(map[string]interface{})
			extAttributesMap["attribute_key"] = extAttributesChild1Raw["AttributeKey"]
			extAttributesMap["attribute_value"] = extAttributesChild1Raw["AttributeValue"]

			extAttributesMaps = append(extAttributesMaps, extAttributesMap)
		}
	}
	if objectRaw["ExtAttributes"] != nil {
		if err := d.Set("ext_attributes", extAttributesMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAliCloudAlbAScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	action := "UpdateAScripts"
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "AScripts.0.AScriptId", d.Id())
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	objectDataLocalMap["AScriptId"] = d.Id()
	if d.HasChange("enabled") {
		update = true
		objectDataLocalMap["Enabled"] = d.Get("enabled")
	}

	if d.HasChange("script_content") {
		update = true
		objectDataLocalMap["ScriptContent"] = d.Get("script_content")
	}

	if d.HasChange("ext_attributes") {
		update = true
		if v := d.Get("ext_attributes"); v != nil {
			if v, ok := d.GetOk("ext_attributes"); ok {
				localData, err := jsonpath.Get("$", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range localData.([]interface{}) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["AttributeKey"] = dataLoopTmp["attribute_key"]
					dataLoopMap["AttributeValue"] = dataLoopTmp["attribute_value"]
					localMaps = append(localMaps, dataLoopMap)
				}
				objectDataLocalMap["ExtAttributes"] = localMaps
			}

		}
	}

	if d.HasChange("ext_attribute_enabled") {
		update = true
		objectDataLocalMap["ExtAttributeEnabled"] = d.Get("ext_attribute_enabled")
	}

	if d.HasChange("ascript_id") {
		update = true
		objectDataLocalMap["AScriptId"] = d.Get("ascript_id")
	}

	if d.HasChange("ascript_name") {
		update = true
		objectDataLocalMap["AScriptName"] = d.Get("ascript_name")
	}

	AScriptsMap := make([]interface{}, 0)
	AScriptsMap = append(AScriptsMap, objectDataLocalMap)
	request["AScripts"] = AScriptsMap
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbAScriptStateRefreshFunc(d.Id(), "AScriptStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudAlbAScriptRead(d, meta)
}

func resourceAliCloudAlbAScriptDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAScripts"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["AScriptIds.1"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albServiceV2.AlbAScriptStateRefreshFunc(d.Id(), "AScriptStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

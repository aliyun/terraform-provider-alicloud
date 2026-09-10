// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudPaiWorkspaceUserConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceUserConfigCreate,
		Read:   resourceAliCloudPaiWorkspaceUserConfigRead,
		Update: resourceAliCloudPaiWorkspaceUserConfigUpdate,
		Delete: resourceAliCloudPaiWorkspaceUserConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"category_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceUserConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/userconfigs")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})
	objectDataLocalMap["CategoryName"] = d.Get("category_name")
	objectDataLocalMap["ConfigKey"] = d.Get("config_key")
	objectDataLocalMap["ConfigValue"] = d.Get("config_value")

	if v, ok := d.GetOk("scope"); ok {
		objectDataLocalMap["Scope"] = v
	}

	ConfigsMap := make([]interface{}, 0)
	ConfigsMap = append(ConfigsMap, objectDataLocalMap)
	request["Configs"] = ConfigsMap

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "Configs.0.ConfigKey", d.Get("config_key"))
	jsonString, _ = sjson.Set(jsonString, "Configs.0.CategoryName", d.Get("category_name"))
	_ = json.Unmarshal([]byte(jsonString), &request)

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPut("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_user_config", action, AlibabaCloudSdkGoERROR)
	}

	ConfigsCategoryNameVar, _ := jsonpath.Get("Configs[0].CategoryName", request)
	ConfigsConfigKeyVar, _ := jsonpath.Get("Configs[0].ConfigKey", request)
	d.SetId(fmt.Sprintf("%v:%v", ConfigsCategoryNameVar, ConfigsConfigKeyVar))

	return resourceAliCloudPaiWorkspaceUserConfigRead(d, meta)
}

func resourceAliCloudPaiWorkspaceUserConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceUserConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_user_config DescribePaiWorkspaceUserConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_value", objectRaw["ConfigValue"])
	d.Set("scope", objectRaw["Scope"])
	d.Set("category_name", objectRaw["CategoryName"])
	d.Set("config_key", objectRaw["ConfigKey"])

	return nil
}

func resourceAliCloudPaiWorkspaceUserConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/api/v1/userconfigs")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	objectDataLocalMap["CategoryName"] = d.Get("category_name")
	objectDataLocalMap["ConfigKey"] = d.Get("config_key")

	if d.HasChange("config_value") {
		update = true
	}
	objectDataLocalMap["ConfigValue"] = d.Get("config_value")

	if v, ok := d.GetOk("scope"); ok {
		objectDataLocalMap["Scope"] = v
	}

	ConfigsMap := make([]interface{}, 0)
	ConfigsMap = append(ConfigsMap, objectDataLocalMap)
	request["Configs"] = ConfigsMap

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "Configs.0.ConfigKey", parts[1])
	jsonString, _ = sjson.Set(jsonString, "Configs.0.CategoryName", parts[0])
	_ = json.Unmarshal([]byte(jsonString), &request)

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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

	return resourceAliCloudPaiWorkspaceUserConfigRead(d, meta)
}

func resourceAliCloudPaiWorkspaceUserConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	CategoryName := parts[0]
	action := fmt.Sprintf("/api/v1/userconfigs/%s", CategoryName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["ConfigKey"] = StringPointer(parts[1])

	if v, ok := d.GetOk("scope"); ok {
		query["Scope"] = StringPointer(v.(string))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AIWorkSpace", "2021-02-04", action, query, nil, nil, true)

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

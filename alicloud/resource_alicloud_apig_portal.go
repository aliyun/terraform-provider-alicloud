// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudApigPortal() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigPortalCreate,
		Read:   resourceAliCloudApigPortalRead,
		Update: resourceAliCloudApigPortalUpdate,
		Delete: resourceAliCloudApigPortalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"portal_setting_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_approve_subscriptions": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auto_approve_developers": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"builtin_auth_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"portal_domain_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudApigPortalCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/himarket/v1/portals")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["name"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_portal", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.portalId", response)
	d.SetId(fmt.Sprint(id))

	if _, ok := d.GetOk("portal_setting_config"); ok {
		if err := resourceAliCloudApigPortalUpdate(d, meta); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliCloudApigPortalRead(d, meta)
}

func resourceAliCloudApigPortalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigPortal(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_portal DescribeApigPortal Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["description"])
	d.Set("name", objectRaw["name"])

	if cfgRaw, ok := objectRaw["portalSettingConfig"]; ok && cfgRaw != nil {
		if cfg, ok := cfgRaw.(map[string]interface{}); ok {
			settingConfig := make([]map[string]interface{}, 0)
			settingConfig = append(settingConfig, map[string]interface{}{
				"builtin_auth_enabled":       cfg["builtinAuthEnabled"],
				"auto_approve_developers":    cfg["autoApproveDevelopers"],
				"auto_approve_subscriptions": cfg["autoApproveSubscriptions"],
			})
			if err := d.Set("portal_setting_config", settingConfig); err != nil {
				return WrapError(err)
			}
		}
	}

	if domainsRaw, ok := objectRaw["portalDomainConfig"]; ok && domainsRaw != nil {
		if domains, ok := domainsRaw.([]interface{}); ok {
			domainConfigs := make([]map[string]interface{}, 0)
			for _, item := range domains {
				if m, ok := item.(map[string]interface{}); ok {
					domainConfigs = append(domainConfigs, map[string]interface{}{
						"domain":   m["domain"],
						"type":     m["type"],
						"protocol": m["protocol"],
					})
				}
			}
			if err := d.Set("portal_domain_config", domainConfigs); err != nil {
				return WrapError(err)
			}
		}
	}

	return nil
}

func resourceAliCloudApigPortalUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	portalId := d.Id()
	action := fmt.Sprintf("/himarket/v1/portals/%s", portalId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if v, ok := d.GetOk("portal_setting_config"); ok {
		portalSettingConfig := make(map[string]interface{})
		if configList, ok := v.([]interface{}); ok && len(configList) > 0 && configList[0] != nil {
			if configMap, ok := configList[0].(map[string]interface{}); ok {
				portalSettingConfig["builtinAuthEnabled"] = configMap["builtin_auth_enabled"]
				portalSettingConfig["autoApproveDevelopers"] = configMap["auto_approve_developers"]
				portalSettingConfig["autoApproveSubscriptions"] = configMap["auto_approve_subscriptions"]
			}
		}
		if len(portalSettingConfig) > 0 {
			request["portalSettingConfig"] = portalSettingConfig
		}
	}
	if d.HasChange("portal_setting_config") {
		update = true
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("name") {
		update = true
	}
	if v, ok := d.GetOk("name"); ok || d.HasChange("name") {
		request["name"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
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

	return resourceAliCloudApigPortalRead(d, meta)
}

func resourceAliCloudApigPortalDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	portalId := d.Id()
	action := fmt.Sprintf("/himarket/v1/portals/%s", portalId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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

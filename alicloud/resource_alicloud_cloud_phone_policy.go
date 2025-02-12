// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudPhonePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudPhonePolicyCreate,
		Read:   resourceAliCloudCloudPhonePolicyRead,
		Update: resourceAliCloudCloudPhonePolicyUpdate,
		Delete: resourceAliCloudCloudPhonePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"camera_redirect": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"clipboard": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"read", "write", "readwrite", "off"}, false),
			},
			"lock_resolution": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"net_redirect_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_user_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("^\\d+$"), "Transparent proxy port. The Port value range is 1\\~ 65535."),
						},
						"net_redirect": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"on", "off"}, false),
						},
						"proxy_password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"proxy_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"socks5"}, false),
						},
						"host_addr": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("^(\\d{0,3}\\.){3}\\d{0,3}$"), "The transparent proxy IP address. The format is IPv4 address."),
						},
						"custom_proxy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"on", "off"}, false),
						},
					},
				},
			},
			"policy_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resolution_height": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 640, 720, 940, 1280, 1920}),
			},
			"resolution_width": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 480, 536, 720, 800, 1080}),
			},
		},
	}
}

func resourceAliCloudCloudPhonePolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePolicyGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOkExists("resolution_width"); ok {
		request["ResolutionWidth"] = v
	}
	if v, ok := d.GetOk("policy_group_name"); ok {
		request["PolicyGroupName"] = v
	}
	if v, ok := d.GetOk("clipboard"); ok {
		request["Clipboard"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("net_redirect_policy"); !IsNil(v) {
		hostAddr1, _ := jsonpath.Get("$[0].host_addr", v)
		if hostAddr1 != nil && hostAddr1 != "" {
			objectDataLocalMap["HostAddr"] = hostAddr1
		}
		netRedirect1, _ := jsonpath.Get("$[0].net_redirect", v)
		if netRedirect1 != nil && netRedirect1 != "" {
			objectDataLocalMap["NetRedirect"] = netRedirect1
		}
		proxyPassword1, _ := jsonpath.Get("$[0].proxy_password", v)
		if proxyPassword1 != nil && proxyPassword1 != "" {
			objectDataLocalMap["ProxyPassword"] = proxyPassword1
		}
		proxyUserName1, _ := jsonpath.Get("$[0].proxy_user_name", v)
		if proxyUserName1 != nil && proxyUserName1 != "" {
			objectDataLocalMap["ProxyUserName"] = proxyUserName1
		}
		customProxy1, _ := jsonpath.Get("$[0].custom_proxy", v)
		if customProxy1 != nil && customProxy1 != "" {
			objectDataLocalMap["CustomProxy"] = customProxy1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && port1 != "" {
			objectDataLocalMap["Port"] = port1
		}
		proxyType1, _ := jsonpath.Get("$[0].proxy_type", v)
		if proxyType1 != nil && proxyType1 != "" {
			objectDataLocalMap["ProxyType"] = proxyType1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["NetRedirectPolicy"] = string(objectDataLocalMapJson)
	}

	if v, ok := d.GetOk("camera_redirect"); ok {
		request["CameraRedirect"] = v
	}
	if v, ok := d.GetOk("lock_resolution"); ok {
		request["LockResolution"] = v
	}
	if v, ok := d.GetOkExists("resolution_height"); ok {
		request["ResolutionHeight"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_phone_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PolicyGroupId"]))

	return resourceAliCloudCloudPhonePolicyRead(d, meta)
}

func resourceAliCloudCloudPhonePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudPhoneServiceV2 := CloudPhoneServiceV2{client}

	objectRaw, err := cloudPhoneServiceV2.DescribeCloudPhonePolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_phone_policy DescribeCloudPhonePolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("camera_redirect", objectRaw["CameraRedirect"])
	d.Set("clipboard", objectRaw["Clipboard"])
	d.Set("lock_resolution", objectRaw["LockResolution"])
	d.Set("policy_group_name", objectRaw["PolicyGroupName"])
	d.Set("resolution_height", objectRaw["SessionResolutionHeight"])
	d.Set("resolution_width", objectRaw["SessionResolutionWidth"])

	netRedirectPolicyMaps := make([]map[string]interface{}, 0)
	netRedirectPolicyMap := make(map[string]interface{})
	netRedirectPolicyRaw := make(map[string]interface{})
	if objectRaw["NetRedirectPolicy"] != nil {
		netRedirectPolicyRaw = objectRaw["NetRedirectPolicy"].(map[string]interface{})
	}
	if len(netRedirectPolicyRaw) > 0 {
		netRedirectPolicyMap["custom_proxy"] = netRedirectPolicyRaw["CustomProxy"]
		netRedirectPolicyMap["host_addr"] = netRedirectPolicyRaw["HostAddr"]
		netRedirectPolicyMap["net_redirect"] = netRedirectPolicyRaw["NetRedirect"]
		netRedirectPolicyMap["port"] = netRedirectPolicyRaw["Port"]
		netRedirectPolicyMap["proxy_password"] = netRedirectPolicyRaw["ProxyPassword"]
		netRedirectPolicyMap["proxy_type"] = netRedirectPolicyRaw["ProxyType"]
		netRedirectPolicyMap["proxy_user_name"] = netRedirectPolicyRaw["ProxyUserName"]

		netRedirectPolicyMaps = append(netRedirectPolicyMaps, netRedirectPolicyMap)
	}
	if err := d.Set("net_redirect_policy", netRedirectPolicyMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCloudPhonePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyPolicyGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PolicyGroupId"] = d.Id()

	if d.HasChange("resolution_width") {
		update = true
	}
	if v, ok := d.GetOk("resolution_width"); ok || d.HasChange("resolution_width") {
		request["ResolutionWidth"] = v
	}
	if d.HasChange("policy_group_name") {
		update = true
	}
	if v, ok := d.GetOk("policy_group_name"); ok || d.HasChange("policy_group_name") {
		request["PolicyGroupName"] = v
	}
	if d.HasChange("clipboard") {
		update = true
	}
	if v, ok := d.GetOk("clipboard"); ok || d.HasChange("clipboard") {
		request["Clipboard"] = v
	}
	if d.HasChange("net_redirect_policy") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("net_redirect_policy"); v != nil {
		hostAddr1, _ := jsonpath.Get("$[0].host_addr", v)
		if hostAddr1 != nil && (d.HasChange("net_redirect_policy.0.host_addr") || hostAddr1 != "") {
			objectDataLocalMap["HostAddr"] = hostAddr1
		}
		netRedirect1, _ := jsonpath.Get("$[0].net_redirect", v)
		if netRedirect1 != nil && (d.HasChange("net_redirect_policy.0.net_redirect") || netRedirect1 != "") {
			objectDataLocalMap["NetRedirect"] = netRedirect1
		}
		proxyPassword1, _ := jsonpath.Get("$[0].proxy_password", v)
		if proxyPassword1 != nil && (d.HasChange("net_redirect_policy.0.proxy_password") || proxyPassword1 != "") {
			objectDataLocalMap["ProxyPassword"] = proxyPassword1
		}
		proxyUserName1, _ := jsonpath.Get("$[0].proxy_user_name", v)
		if proxyUserName1 != nil && (d.HasChange("net_redirect_policy.0.proxy_user_name") || proxyUserName1 != "") {
			objectDataLocalMap["ProxyUserName"] = proxyUserName1
		}
		customProxy1, _ := jsonpath.Get("$[0].custom_proxy", v)
		if customProxy1 != nil && (d.HasChange("net_redirect_policy.0.custom_proxy") || customProxy1 != "") {
			objectDataLocalMap["CustomProxy"] = customProxy1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && (d.HasChange("net_redirect_policy.0.port") || port1 != "") {
			objectDataLocalMap["Port"] = port1
		}
		proxyType1, _ := jsonpath.Get("$[0].proxy_type", v)
		if proxyType1 != nil && (d.HasChange("net_redirect_policy.0.proxy_type") || proxyType1 != "") {
			objectDataLocalMap["ProxyType"] = proxyType1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["NetRedirectPolicy"] = string(objectDataLocalMapJson)
	}

	if d.HasChange("camera_redirect") {
		update = true
	}
	if v, ok := d.GetOk("camera_redirect"); ok || d.HasChange("camera_redirect") {
		request["CameraRedirect"] = v
	}
	if d.HasChange("lock_resolution") {
		update = true
	}
	if v, ok := d.GetOk("lock_resolution"); ok || d.HasChange("lock_resolution") {
		request["LockResolution"] = v
	}
	if d.HasChange("resolution_height") {
		update = true
	}
	if v, ok := d.GetOk("resolution_height"); ok || d.HasChange("resolution_height") {
		request["ResolutionHeight"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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

	return resourceAliCloudCloudPhonePolicyRead(d, meta)
}

func resourceAliCloudCloudPhonePolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicyGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PolicyGroupIds.1"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)

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

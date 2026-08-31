// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudArmsAddonRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsAddonReleaseCreate,
		Read:   resourceAliCloudArmsAddonReleaseRead,
		Update: resourceAliCloudArmsAddonReleaseUpdate,
		Delete: resourceAliCloudArmsAddonReleaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"addon_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"addon_release_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"addon_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aliyun_lang": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"values": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudArmsAddonReleaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "InstallAddon"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["EnvironmentId"] = d.Get("environment_id")
	query["ReleaseName"] = d.Get("addon_release_name")
	request["RegionId"] = client.RegionId

	request["AddonVersion"] = d.Get("addon_version")
	request["Values"] = d.Get("values")
	if v, ok := d.GetOk("values"); ok {
		request["Values"] = v
	}
	request["Name"] = d.Get("addon_name")
	if v, ok := d.GetOk("aliyun_lang"); ok {
		request["AliyunLang"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_addon_release", action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.Code", response)
	if fmt.Sprint(code) != "200" {
		log.Printf("[DEBUG] Resource alicloud_arms_addon_release InstallAddon Failed!!! %s", response)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_addon_release", action, AlibabaCloudSdkGoERROR, response)
	}

	environmentId, _ := jsonpath.Get("$.Data.EnvironmentId", response)
	releaseName, _ := jsonpath.Get("$.Data.ReleaseName", response)
	d.SetId(fmt.Sprintf("%v:%v", environmentId, releaseName))

	return resourceAliCloudArmsAddonReleaseRead(d, meta)
}

func resourceAliCloudArmsAddonReleaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ArmsServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsAddonRelease(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_addon_release DescribeArmsAddonRelease Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("values", objectRaw["Config"])
	release1RawObj, _ := jsonpath.Get("$.Release", objectRaw)
	release1Raw := make(map[string]interface{})
	if release1RawObj != nil {
		release1Raw = release1RawObj.(map[string]interface{})
	}
	d.Set("addon_name", release1Raw["AddonName"])
	d.Set("addon_version", release1Raw["Version"])
	d.Set("aliyun_lang", release1Raw["Language"])
	d.Set("create_time", release1Raw["CreateTime"])
	d.Set("addon_release_name", release1Raw["ReleaseName"])
	d.Set("environment_id", release1Raw["EnvironmentId"])

	return nil
}

func resourceAliCloudArmsAddonReleaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpgradeAddonRelease"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["EnvironmentId"] = parts[0]
	query["ReleaseName"] = parts[1]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("addon_version") {
		update = true
	}
	request["AddonVersion"] = d.Get("addon_version")
	request["Values"] = d.Get("values")
	if !d.IsNewResource() && d.HasChange("values") {
		update = true
		request["Values"] = d.Get("values")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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

	return resourceAliCloudArmsAddonReleaseRead(d, meta)
}

func resourceAliCloudArmsAddonReleaseDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAddonRelease"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["EnvironmentId"] = parts[0]
	query["ReleaseName"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"404"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	armsServiceV2 := ArmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, armsServiceV2.ArmsAddonReleaseStateRefreshFunc(d.Id(), "$.Release.ReleaseName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

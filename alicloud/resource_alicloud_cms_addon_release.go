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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudCmsAddonRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAddonReleaseCreate,
		Read:   resourceAliCloudCmsAddonReleaseRead,
		Update: resourceAliCloudCmsAddonReleaseUpdate,
		Delete: resourceAliCloudCmsAddonReleaseDelete,
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
			"config": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"env_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"integration_policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsAddonReleaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	policyId := d.Get("integration_policy_id")
	action := fmt.Sprintf("/integration-policies/%s/addon-releases", policyId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("addon_release_name"); ok {
		request["releaseName"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		request["values"] = v
	}
	if v, ok := d.GetOk("env_type"); ok {
		request["envType"] = v
	}
	if v, ok := d.GetOk("workspace"); ok {
		request["workspace"] = v
	}
	if v, ok := d.GetOk("aliyun_lang"); ok {
		request["aliyunLang"] = v
	}
	request["version"] = d.Get("addon_version")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["dryRun"] = v
	}
	request["addonName"] = d.Get("addon_name")
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_addon_release", action, AlibabaCloudSdkGoERROR)
	}

	releaseenvironmentIdVar, _ := jsonpath.Get("$.release.environmentId", response)
	releasereleaseNameVar, _ := jsonpath.Get("$.release.releaseName", response)
	d.SetId(fmt.Sprintf("%v:%v", releaseenvironmentIdVar, releasereleaseNameVar))

	return resourceAliCloudCmsAddonReleaseRead(d, meta)
}

func resourceAliCloudCmsAddonReleaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsAddonRelease(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_addon_release DescribeCmsAddonRelease Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("addon_name", objectRaw["addonName"])
	d.Set("addon_version", objectRaw["version"])
	d.Set("aliyun_lang", objectRaw["language"])
	d.Set("config", objectRaw["config"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("env_type", objectRaw["envType"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("workspace", objectRaw["workspace"])
	d.Set("addon_release_name", objectRaw["releaseName"])
	d.Set("integration_policy_id", objectRaw["policyId"])

	return nil
}

func resourceAliCloudCmsAddonReleaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	releaseName := parts[1]
	policyId := parts[0]
	action := fmt.Sprintf("/integration-policies/%s/addon-releases/%s", policyId, releaseName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("config") {
		update = true
	}
	if v, ok := d.GetOk("config"); ok || d.HasChange("config") {
		request["values"] = v
	}
	if d.HasChange("addon_version") {
		update = true
	}
	request["addonVersion"] = d.Get("addon_version")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["dryRun"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
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

	return resourceAliCloudCmsAddonReleaseRead(d, meta)
}

func resourceAliCloudCmsAddonReleaseDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	policyId := parts[0]
	action := fmt.Sprintf("/integration-policies/%s/addon-releases", policyId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["releaseName"] = StringPointer(parts[1])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cmsServiceV2 := CmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cmsServiceV2.CmsAddonReleaseStateRefreshFunc(d.Id(), "releaseName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsSlsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsSlsGroupCreate,
		Read:   resourceAlicloudCmsSlsGroupRead,
		Update: resourceAlicloudCmsSlsGroupUpdate,
		Delete: resourceAlicloudCmsSlsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"sls_group_config": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 25,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sls_logstore": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sls_project": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sls_region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sls_user_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"sls_group_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sls_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9_]{1,31}$`), "The name must be `2` to `32` characters in length, and can contain letters, digits and underscores (_). It must start with a letter."),
			},
		},
	}
}

func resourceAlicloudCmsSlsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHybridMonitorSLSGroup"
	request := make(map[string]interface{})
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	slsGroupConfigMaps := make([]map[string]interface{}, 0)
	for _, slsGroupConfig := range d.Get("sls_group_config").(*schema.Set).List() {
		slsGroupConfigArg := slsGroupConfig.(map[string]interface{})
		slsGroupConfigMaps = append(slsGroupConfigMaps, map[string]interface{}{
			"SLSLogstore": slsGroupConfigArg["sls_logstore"],
			"SLSProject":  slsGroupConfigArg["sls_project"],
			"SLSRegion":   slsGroupConfigArg["sls_region"],
			"SLSUserId":   slsGroupConfigArg["sls_user_id"],
		})
	}
	request["SLSGroupConfig"] = slsGroupConfigMaps

	if v, ok := d.GetOk("sls_group_description"); ok {
		request["SLSGroupDescription"] = v
	}
	request["SLSGroupName"] = d.Get("sls_group_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_sls_group", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["SLSGroupName"]))

	return resourceAlicloudCmsSlsGroupRead(d, meta)
}
func resourceAlicloudCmsSlsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsSlsGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_sls_group cmsService.DescribeCmsSlsGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("sls_group_name", d.Id())
	if sLSGroupConfigList, ok := object["SLSGroupConfig"]; ok && sLSGroupConfigList != nil {
		slsGroupConfigMaps := make([]map[string]interface{}, 0)
		for _, sLSGroupConfigListItem := range sLSGroupConfigList.([]interface{}) {
			if sLSGroupConfigListItemMap, ok := sLSGroupConfigListItem.(map[string]interface{}); ok {
				sLSGroupConfigListItemMap["sls_logstore"] = sLSGroupConfigListItemMap["SLSLogstore"]
				sLSGroupConfigListItemMap["sls_project"] = sLSGroupConfigListItemMap["SLSProject"]
				sLSGroupConfigListItemMap["sls_region"] = sLSGroupConfigListItemMap["SLSRegion"]
				sLSGroupConfigListItemMap["sls_user_id"] = sLSGroupConfigListItemMap["SLSUserId"]
				slsGroupConfigMaps = append(slsGroupConfigMaps, sLSGroupConfigListItemMap)
			}
			d.Set("sls_group_config", slsGroupConfigMaps)
		}
	}

	d.Set("sls_group_description", object["SLSGroupDescription"])
	return nil
}
func resourceAlicloudCmsSlsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"SLSGroupName": d.Id(),
	}

	slsGroupConfigMaps := make([]map[string]interface{}, 0)
	for _, slsGroupConfig := range d.Get("sls_group_config").(*schema.Set).List() {
		slsGroupConfigArg := slsGroupConfig.(map[string]interface{})
		slsGroupConfigMaps = append(slsGroupConfigMaps, map[string]interface{}{
			"SLSLogstore": slsGroupConfigArg["sls_logstore"],
			"SLSProject":  slsGroupConfigArg["sls_project"],
			"SLSRegion":   slsGroupConfigArg["sls_region"],
			"SLSUserId":   slsGroupConfigArg["sls_user_id"],
		})
	}
	request["SLSGroupConfig"] = slsGroupConfigMaps

	if d.HasChange("sls_group_config") {
		update = true
	}

	if d.HasChange("sls_group_description") {
		update = true
	}
	if v, ok := d.GetOk("sls_group_description"); ok {
		request["SLSGroupDescription"] = v
	}
	if update {
		action := "ModifyHybridMonitorSLSGroup"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return resourceAlicloudCmsSlsGroupRead(d, meta)
}
func resourceAlicloudCmsSlsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHybridMonitorSLSGroup"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"SLSGroupName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"404"}) {
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}

package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsDynamicTagGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsDynamicTagGroupCreate,
		Read:   resourceAliCloudCmsDynamicTagGroupRead,
		Delete: resourceAliCloudCmsDynamicTagGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tag_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_express_filter_relation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"contact_group_list": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"template_id_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"match_express": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"tag_value_match_function": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsDynamicTagGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	var response map[string]interface{}
	action := "CreateDynamicTagGroup"
	request := make(map[string]interface{})
	var err error

	request["TagRegionId"] = client.RegionId
	request["TagKey"] = d.Get("tag_key")
	request["ContactGroupList"] = d.Get("contact_group_list")

	if v, ok := d.GetOk("match_express_filter_relation"); ok {
		request["MatchExpressFilterRelation"] = v
	}

	if v, ok := d.GetOk("template_id_list"); ok {
		request["TemplateIdList"] = v
	}

	matchExpress := d.Get("match_express")
	matchExpressMaps := make([]map[string]interface{}, 0)
	for _, matchExpressList := range matchExpress.(*schema.Set).List() {
		matchExpressMap := map[string]interface{}{}
		matchExpressArg := matchExpressList.(map[string]interface{})

		matchExpressMap["TagValue"] = matchExpressArg["tag_value"]
		matchExpressMap["TagValueMatchFunction"] = matchExpressArg["tag_value_match_function"]

		matchExpressMaps = append(matchExpressMaps, matchExpressMap)
	}

	request["MatchExpress"] = matchExpressMaps
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_dynamic_tag_group", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["Id"]))

	stateConf := BuildStateConf([]string{}, []string{"FINISH"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cmsService.CmsDynamicTagGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCmsDynamicTagGroupRead(d, meta)
}

func resourceAliCloudCmsDynamicTagGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	object, err := cmsService.DescribeCmsDynamicTagGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_dynamic_tag_group cmsService.DescribeCmsDynamicTagGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("tag_key", object["TagKey"])
	d.Set("match_express_filter_relation", object["MatchExpressFilterRelation"])

	if contactGroupList, ok := object["ContactGroupList"]; ok {
		if contactGroupLists, ok := contactGroupList.(map[string]interface{})["ContactGroupList"]; ok {

			d.Set("contact_group_list", contactGroupLists.([]interface{}))
		}
	}

	if templateIdList, ok := object["TemplateIdList"]; ok {
		if templateIdLists, ok := templateIdList.(map[string]interface{})["TemplateIdList"]; ok {

			d.Set("template_id_list", templateIdLists.([]interface{}))
		}
	}

	if matchExpress, ok := object["MatchExpress"]; ok {
		if matchExpressList, ok := matchExpress.(map[string]interface{})["MatchExpress"]; ok {
			matchExpressMaps := make([]map[string]interface{}, 0)
			for _, matchExpresses := range matchExpressList.([]interface{}) {
				matchExpressArg := matchExpresses.(map[string]interface{})
				matchExpressMap := map[string]interface{}{}

				if tagValue, ok := matchExpressArg["TagValue"]; ok {
					matchExpressMap["tag_value"] = tagValue
				}

				if tagValueMatchFunction, ok := matchExpressArg["TagValueMatchFunction"]; ok {
					matchExpressMap["tag_value_match_function"] = tagValueMatchFunction
				}

				matchExpressMaps = append(matchExpressMaps, matchExpressMap)
			}

			d.Set("match_express", matchExpressMaps)
		}
	}

	d.Set("status", object["Status"])

	return nil
}

func resourceAliCloudCmsDynamicTagGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	action := "DeleteDynamicTagGroup"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"DynamicTagRuleId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cmsService.CmsDynamicTagGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

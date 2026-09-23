// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudApigAiModelProviders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApigAiModelProviderRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"model_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"model_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"model_provider_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bound_services": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"express_type":       {Type: schema.TypeString, Computed: true},
									"group_name":         {Type: schema.TypeString, Computed: true},
									"name":               {Type: schema.TypeString, Computed: true},
									"namespace":          {Type: schema.TypeString, Computed: true},
									"pai_workspace_id":   {Type: schema.TypeString, Computed: true},
									"pai_workspace_name": {Type: schema.TypeString, Computed: true},
									"qualifier":          {Type: schema.TypeString, Computed: true},
									"service_id":         {Type: schema.TypeString, Computed: true},
									"source_type":        {Type: schema.TypeString, Computed: true},
									"status":             {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"model_cards": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_id":     {Type: schema.TypeString, Computed: true},
									"model_card_id":  {Type: schema.TypeString, Computed: true},
									"model_name":     {Type: schema.TypeString, Computed: true},
									"model_provider": {Type: schema.TypeString, Computed: true},
									"source":         {Type: schema.TypeString, Computed: true},
									"update_time":    {Type: schema.TypeString, Computed: true},
								},
							},
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudApigAiModelProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	// ListAiModelProviders
	action := fmt.Sprintf("/v1/ai-model-providers")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("gateway_id"); ok {
		query["gatewayId"] = StringPointer(v.(string))
	}

	query["pageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	query["pageNumber"] = StringPointer("1")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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

		resp, _ := jsonpath.Get("$.data.items[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["modelProviderId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["pageNumber"])
		query["pageNumber"] = StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["modelProviderId"]

		mapping["model_count"] = objectRaw["modelCount"]
		mapping["source"] = objectRaw["source"]
		mapping["update_time"] = objectRaw["updateTime"]
		mapping["display_name"] = objectRaw["displayName"]
		mapping["gateway_id"] = objectRaw["gatewayId"]
		mapping["model_provider"] = objectRaw["provider"]
		mapping["model_provider_id"] = objectRaw["modelProviderId"]

		boundServicesRaw := objectRaw["boundServices"]
		boundServicesMaps := make([]map[string]interface{}, 0)
		if boundServicesRaw != nil {
			for _, boundServicesChildRaw := range convertToInterfaceArray(boundServicesRaw) {
				boundServicesMap := make(map[string]interface{})
				boundServicesChildRaw := boundServicesChildRaw.(map[string]interface{})
				boundServicesMap["express_type"] = boundServicesChildRaw["expressType"]
				boundServicesMap["group_name"] = boundServicesChildRaw["groupName"]
				boundServicesMap["name"] = boundServicesChildRaw["name"]
				boundServicesMap["namespace"] = boundServicesChildRaw["namespace"]
				boundServicesMap["pai_workspace_id"] = boundServicesChildRaw["paiWorkspaceId"]
				boundServicesMap["pai_workspace_name"] = boundServicesChildRaw["paiWorkspaceName"]
				boundServicesMap["qualifier"] = boundServicesChildRaw["qualifier"]
				boundServicesMap["service_id"] = boundServicesChildRaw["serviceId"]
				boundServicesMap["source_type"] = boundServicesChildRaw["sourceType"]
				boundServicesMap["status"] = boundServicesChildRaw["status"]

				boundServicesMaps = append(boundServicesMaps, boundServicesMap)
			}
		}
		mapping["bound_services"] = boundServicesMaps
		modelCardsRaw := objectRaw["modelCards"]
		modelCardsMaps := make([]map[string]interface{}, 0)
		if modelCardsRaw != nil {
			for _, modelCardsChildRaw := range convertToInterfaceArray(modelCardsRaw) {
				modelCardsMap := make(map[string]interface{})
				modelCardsChildRaw := modelCardsChildRaw.(map[string]interface{})
				modelCardsMap["gateway_id"] = modelCardsChildRaw["gatewayId"]
				modelCardsMap["model_card_id"] = modelCardsChildRaw["modelCardId"]
				modelCardsMap["model_name"] = modelCardsChildRaw["modelName"]
				modelCardsMap["model_provider"] = modelCardsChildRaw["modelProvider"]
				modelCardsMap["source"] = modelCardsChildRaw["source"]
				modelCardsMap["update_time"] = modelCardsChildRaw["updateTime"]

				modelCardsMaps = append(modelCardsMaps, modelCardsMap)
			}
		}
		mapping["model_cards"] = modelCardsMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("providers", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

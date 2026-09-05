package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionCustomCheckItems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionCustomCheckItemsRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of Custom Check Item IDs.",
			},
			"check_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the custom check item.",
			},
			"check_show_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the check item.",
			},
			"check_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The types of the check item to query.",
			},
			"lang": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The language of the check item.",
			},
			"statuses": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The statuses of the check item to query.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     PageSizeLarge,
				Description: "The page size of the list query.",
			},
			"current_page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The current page number of the list query.",
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the custom check item.",
						},
						"check_show_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the check item.",
						},
						"vendor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vendor to which the check item belongs.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The asset type to which the check item belongs.",
						},
						"instance_sub_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The asset subtype to which the check item belongs.",
						},
						"risk_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The risk level of the check item.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the check item.",
						},
						"description": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"solution": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"assist_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"check_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The check rule of the check item.",
						},
						"check_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the check item.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The remarks of the check item.",
						},
					},
				},
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Save the result to a file.",
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionCustomCheckItemsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ListCheckItems"

	pageSize := PageSizeLarge
	if v, ok := d.GetOk("page_size"); ok {
		pageSize = v.(int)
	}
	currentPage := 1
	singlePage := false
	if v, ok := d.GetOk("current_page"); ok {
		currentPage = v.(int)
		singlePage = true
	}

	request := map[string]interface{}{
		"PageSize":    pageSize,
		"CurrentPage": currentPage,
	}
	query := map[string]interface{}{}
	var response map[string]interface{}
	var err error

	if v, ok := d.GetOk("check_id"); ok {
		request["CheckId"] = v
	}
	if v, ok := d.GetOk("check_show_name"); ok {
		request["CheckShowName"] = v
	}
	if v, ok := d.GetOk("check_types"); ok {
		request["CheckTypes"] = convertToInterfaceArray(v)
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("statuses"); ok {
		request["Statuses"] = convertToInterfaceArray(v)
	}

	var objects []map[string]interface{}
	ids := make([]string, 0)

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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

		resp, _ := jsonpath.Get("$.CheckItems[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
			ids = append(ids, fmt.Sprint(item["CheckId"]))
		}

		if singlePage || len(result) < pageSize {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{
			"check_id":          objectRaw["CheckId"],
			"check_show_name":   objectRaw["CheckShowName"],
			"vendor":            objectRaw["Vendor"],
			"instance_type":     objectRaw["InstanceType"],
			"instance_sub_type": objectRaw["InstanceSubType"],
			"risk_level":        objectRaw["RiskLevel"],
			"status":            objectRaw["Status"],
			"check_rule":        objectRaw["CheckRule"],
			"check_type":        objectRaw["CheckType"],
			"remark":            objectRaw["Remark"],
			"description":       flattenCustomCheckItemStruct(objectRaw["Description"]),
			"solution":          flattenCustomCheckItemStruct(objectRaw["Solution"]),
			"assist_info":       flattenCustomCheckItemStruct(objectRaw["AssistInfo"]),
		}
		s = append(s, mapping)
	}

	if err := d.Set("items", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

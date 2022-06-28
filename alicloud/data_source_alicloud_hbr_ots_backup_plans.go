package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudHbrOtsBackupPlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrOtsBackupPlansRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"plan_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"plan_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plans": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"retention": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ots_backup_plan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ots_backup_plan_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"schedule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ots_detail": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"table_names": {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHbrOtsBackupPlansRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeBackupPlans"
	request := make(map[string]interface{})
	request["SourceType"] = "OTS"
	filtersMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("plan_id"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"Key":    "planId",
			"Values": []string{v.(string)},
		})
	}
	if v, ok := d.GetOk("plan_name"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"Key":    "planName",
			"Values": []string{v.(string)},
		})
	}
	if v, ok := d.GetOk("vault_id"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"Key":    "vaultId",
			"Values": []string{v.(string)},
		})
	}
	request["Filters"] = filtersMapList
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var otsBackupPlanNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		otsBackupPlanNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_ots_backup_plans", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.BackupPlans.BackupPlan", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if otsBackupPlanNameRegex != nil && !otsBackupPlanNameRegex.MatchString(fmt.Sprint(item["PlanName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PlanId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"vault_id":             object["VaultId"],
			"backup_type":          object["BackupType"],
			"source_type":          object["SourceType"],
			"disabled":             object["Disabled"],
			"retention":            fmt.Sprint(object["Retention"]),
			"created_time":         fmt.Sprint(object["CreatedTime"]),
			"id":                   fmt.Sprint(object["PlanId"]),
			"ots_backup_plan_id":   fmt.Sprint(object["PlanId"]),
			"ots_backup_plan_name": object["PlanName"],
			"schedule":             object["Schedule"],
			"updated_time":         fmt.Sprint(object["UpdatedTime"]),
		}
		otsDetails := make([]map[string]interface{}, 0)
		if v, ok := object["OtsDetail"].(map[string]interface{}); ok {
			otsDetail := make(map[string]interface{}, 0)
			otsDetail["table_names"] = v["TableNames"].(map[string]interface{})["TableName"]
			otsDetails = append(otsDetails, otsDetail)
		}
		mapping["ots_detail"] = otsDetails
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["PlanName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("plans", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

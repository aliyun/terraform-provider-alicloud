package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCmsHybridMonitorSlsTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsHybridMonitorSlsTaskCreate,
		Read:   resourceAlicloudCmsHybridMonitorSlsTaskRead,
		Update: resourceAlicloudCmsHybridMonitorSlsTaskUpdate,
		Delete: resourceAlicloudCmsHybridMonitorSlsTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"attach_labels": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"collect_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{15, 60}),
			},
			"collect_target_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sls_process_config": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"express": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alias": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"express": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"filter": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filters": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{">", ">=", "=", "<=", "<", "!=", "contain", "notContain"}, false),
												},
												"sls_key_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"relation": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"and", "or"}, false),
									},
								},
							},
						},
						"group_by": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alias": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"sls_key_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"statistics": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alias": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"function": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"count", "sum", "avg", "max", "min", "value", "countps", "sumps", "distinct", "distribution", "percentile"}, false),
									},
									"parameter_one": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter_two": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"sls_key_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCmsHybridMonitorSlsTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHybridMonitorTask"
	request := make(map[string]interface{})
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("attach_labels"); ok {
		for attachLabelsPtr, attachLabels := range v.(*schema.Set).List() {
			attachLabelsArg := attachLabels.(map[string]interface{})
			request["AttachLabels."+fmt.Sprint(attachLabelsPtr+1)+".Name"] = attachLabelsArg["name"]
			request["AttachLabels."+fmt.Sprint(attachLabelsPtr+1)+".Value"] = attachLabelsArg["value"]
		}
	}
	if v, ok := d.GetOk("collect_interval"); ok {
		request["CollectInterval"] = v
	}
	request["CollectTargetType"] = d.Get("collect_target_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Namespace"] = d.Get("namespace")
	if v, ok := d.GetOk("sls_process_config"); ok {
		slsProcessConfigList := v.(*schema.Set).List()
		if len(slsProcessConfigList) > 0 {
			slsProcessConfig := slsProcessConfigList[0]
			if slsProcessConfigArg, ok := slsProcessConfig.(map[string]interface{}); ok {
				slsProcessConfigMap := make(map[string]interface{})
				expressSli := make([]map[string]interface{}, 0)
				if v, ok := slsProcessConfigArg["express"]; ok && len(v.(*schema.Set).List()) > 0 {
					for _, express := range v.(*schema.Set).List() {
						expressMap := make(map[string]interface{})
						expressMap["Alias"] = express.(map[string]interface{})["alias"]
						expressMap["Express"] = express.(map[string]interface{})["express"]
						expressSli = append(expressSli, expressMap)
					}
				}
				slsProcessConfigMap["Express"] = expressSli

				if v, ok := slsProcessConfigArg["filter"]; ok && len(v.(*schema.Set).List()) > 0 {
					filterMap := make(map[string]interface{})
					filterArg := v.(*schema.Set).List()[0].(map[string]interface{})
					filtersSli := make([]map[string]interface{}, 0)
					if v, ok := filterArg["filters"]; ok && len(v.(*schema.Set).List()) > 0 {
						for _, filters := range v.(*schema.Set).List() {
							filtersMap := make(map[string]interface{})
							filtersMap["Operator"] = filters.(map[string]interface{})["operator"]
							filtersMap["SLSKeyName"] = filters.(map[string]interface{})["sls_key_name"]
							filtersMap["Value"] = filters.(map[string]interface{})["value"]
							filtersSli = append(filtersSli, filtersMap)
						}
					}
					filterMap["Filters"] = filtersSli
					filterMap["Relation"] = filterArg["relation"]
					slsProcessConfigMap["Filter"] = filterMap
				}

				groupBySli := make([]map[string]interface{}, 0)
				if v, ok := slsProcessConfigArg["group_by"]; ok && len(v.(*schema.Set).List()) > 0 {

					for _, groupBy := range v.(*schema.Set).List() {
						groupByMap := make(map[string]interface{})
						groupByMap["Alias"] = groupBy.(map[string]interface{})["alias"]
						groupByMap["SLSKeyName"] = groupBy.(map[string]interface{})["sls_key_name"]
						groupBySli = append(groupBySli, groupByMap)
					}
				}
				slsProcessConfigMap["GroupBy"] = groupBySli

				statisticsSli := make([]map[string]interface{}, 0)
				if v, ok := slsProcessConfigArg["statistics"]; ok && len(v.(*schema.Set).List()) > 0 {
					for _, statistics := range v.(*schema.Set).List() {
						statisticsMap := make(map[string]interface{})
						statisticsMap["Alias"] = statistics.(map[string]interface{})["alias"]
						statisticsMap["Function"] = statistics.(map[string]interface{})["function"]
						statisticsMap["Parameter1"] = statistics.(map[string]interface{})["parameter_one"]
						statisticsMap["Parameter2"] = statistics.(map[string]interface{})["parameter_two"]
						statisticsMap["SLSKeyName"] = statistics.(map[string]interface{})["sls_key_name"]
						statisticsSli = append(statisticsSli, statisticsMap)
					}
				}
				slsProcessConfigMap["Statistics"] = statisticsSli
				request["SLSProcessConfig"] = slsProcessConfigMap
			}
		}
	}
	request["TaskName"] = d.Get("task_name")
	request["TaskType"] = "aliyun_sls"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_hybrid_monitor_sls_task", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["TaskId"]))

	return resourceAlicloudCmsHybridMonitorSlsTaskRead(d, meta)
}
func resourceAlicloudCmsHybridMonitorSlsTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsHybridMonitorSlsTask(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_hybrid_monitor_sls_task cmsService.DescribeCmsHybridMonitorSlsTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if attachLabelsList, ok := object["AttachLabels"]; ok && attachLabelsList != nil {
		attachLabelsMaps := make([]map[string]interface{}, 0)
		for _, attachLabelsListItem := range attachLabelsList.([]interface{}) {
			if attachLabelsListItemMap, ok := attachLabelsListItem.(map[string]interface{}); ok {
				attachLabelsMaps = append(attachLabelsMaps, map[string]interface{}{
					"name":  attachLabelsListItemMap["Name"],
					"value": attachLabelsListItemMap["Value"],
				})
			}
			d.Set("attach_labels", attachLabelsMaps)
		}
	}

	if v, ok := object["CollectInterval"]; ok && fmt.Sprint(v) != "0" {
		d.Set("collect_interval", formatInt(v))
	}
	d.Set("collect_target_type", object["CollectTargetType"])
	d.Set("description", object["Description"])
	d.Set("namespace", object["Namespace"])

	slsProcessConfigSli := make([]map[string]interface{}, 0)
	if len(object["SLSProcessConfig"].(map[string]interface{})) > 0 {
		slsProcessConfig := object["SLSProcessConfig"]
		slsProcessConfigMap := make(map[string]interface{})

		expressSli := make([]map[string]interface{}, 0)
		if len(slsProcessConfig.(map[string]interface{})["Express"].([]interface{})) > 0 {
			for _, express := range slsProcessConfig.(map[string]interface{})["Express"].([]interface{}) {
				expressMap := make(map[string]interface{})
				expressMap["alias"] = express.(map[string]interface{})["Alias"]
				expressMap["express"] = express.(map[string]interface{})["Express"]
				expressSli = append(expressSli, expressMap)
			}
		}
		slsProcessConfigMap["express"] = expressSli

		filterSli := make([]map[string]interface{}, 0)
		if len(slsProcessConfig.(map[string]interface{})["Filter"].(map[string]interface{})) > 0 {
			filter := slsProcessConfig.(map[string]interface{})["Filter"]
			filterMap := make(map[string]interface{})

			filtersSli := make([]map[string]interface{}, 0)
			if len(filter.(map[string]interface{})["Filters"].([]interface{})) > 0 {
				for _, filters := range filter.(map[string]interface{})["Filters"].([]interface{}) {
					filtersMap := make(map[string]interface{})
					filtersMap["operator"] = filters.(map[string]interface{})["Operator"]
					filtersMap["sls_key_name"] = filters.(map[string]interface{})["SLSKeyName"]
					filtersMap["value"] = filters.(map[string]interface{})["Value"]
					filtersSli = append(filtersSli, filtersMap)
				}
			}
			filterMap["filters"] = filtersSli
			filterMap["relation"] = filter.(map[string]interface{})["Relation"]
			filterSli = append(filterSli, filterMap)
		}
		slsProcessConfigMap["filter"] = filterSli

		groupBySli := make([]map[string]interface{}, 0)
		if len(slsProcessConfig.(map[string]interface{})["GroupBy"].([]interface{})) > 0 {
			for _, groupBy := range slsProcessConfig.(map[string]interface{})["GroupBy"].([]interface{}) {
				groupByMap := make(map[string]interface{})
				groupByMap["alias"] = groupBy.(map[string]interface{})["Alias"]
				groupByMap["sls_key_name"] = groupBy.(map[string]interface{})["SLSKeyName"]
				groupBySli = append(groupBySli, groupByMap)
			}
		}
		slsProcessConfigMap["group_by"] = groupBySli

		statisticsSli := make([]map[string]interface{}, 0)
		if len(slsProcessConfig.(map[string]interface{})["Statistics"].([]interface{})) > 0 {
			for _, statistics := range slsProcessConfig.(map[string]interface{})["Statistics"].([]interface{}) {
				statisticsMap := make(map[string]interface{})
				statisticsMap["alias"] = statistics.(map[string]interface{})["Alias"]
				statisticsMap["function"] = statistics.(map[string]interface{})["Function"]
				statisticsMap["parameter_one"] = statistics.(map[string]interface{})["Parameter1"]
				statisticsMap["parameter_two"] = statistics.(map[string]interface{})["Parameter2"]
				statisticsMap["sls_key_name"] = statistics.(map[string]interface{})["SLSKeyName"]
				statisticsSli = append(statisticsSli, statisticsMap)
			}
		}
		slsProcessConfigMap["statistics"] = statisticsSli
		slsProcessConfigSli = append(slsProcessConfigSli, slsProcessConfigMap)
	}
	d.Set("sls_process_config", slsProcessConfigSli)
	d.Set("task_name", object["TaskName"])
	return nil
}
func resourceAlicloudCmsHybridMonitorSlsTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"TaskId": d.Id(),
	}
	if d.HasChange("attach_labels") {
		update = true
		if v, ok := d.GetOk("attach_labels"); ok {
			for attachLabelsPtr, attachLabels := range v.(*schema.Set).List() {
				attachLabelsArg := attachLabels.(map[string]interface{})
				request["AttachLabels."+fmt.Sprint(attachLabelsPtr+1)+".Name"] = attachLabelsArg["name"]
				request["AttachLabels."+fmt.Sprint(attachLabelsPtr+1)+".Value"] = attachLabelsArg["value"]
			}
		}
	}
	if d.HasChange("description") {
		update = true
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if d.HasChange("collect_interval") {
		update = true
	}

	if v, ok := d.GetOk("collect_interval"); ok {
		request["CollectInterval"] = v
	}

	if d.HasChange("sls_process_config") {
		update = true
	}
	if v, ok := d.GetOk("sls_process_config"); ok {
		slsProcessConfigList := v.(*schema.Set).List()
		if len(slsProcessConfigList) > 0 {
			slsProcessConfig := slsProcessConfigList[0]
			if slsProcessConfigArg, ok := slsProcessConfig.(map[string]interface{}); ok {
				slsProcessConfigMap := make(map[string]interface{})
				expressSli := make([]map[string]interface{}, 0)
				if v, ok := slsProcessConfigArg["express"]; ok && len(v.(*schema.Set).List()) > 0 {
					for _, express := range v.(*schema.Set).List() {
						expressMap := make(map[string]interface{})
						expressMap["Alias"] = express.(map[string]interface{})["alias"]
						expressMap["Express"] = express.(map[string]interface{})["express"]
						expressSli = append(expressSli, expressMap)
					}
				}
				slsProcessConfigMap["Express"] = expressSli

				if v, ok := slsProcessConfigArg["filter"]; ok && len(v.(*schema.Set).List()) > 0 {
					filterMap := make(map[string]interface{})
					filterArg := v.(*schema.Set).List()[0].(map[string]interface{})
					filtersSli := make([]map[string]interface{}, 0)
					if vv, ok := filterArg["filters"]; ok && len(vv.(*schema.Set).List()) > 0 {
						for _, filters := range vv.(*schema.Set).List() {
							filtersMap := make(map[string]interface{})
							filtersMap["Operator"] = filters.(map[string]interface{})["operator"]
							filtersMap["SLSKeyName"] = filters.(map[string]interface{})["sls_key_name"]
							filtersMap["Value"] = filters.(map[string]interface{})["value"]
							filtersSli = append(filtersSli, filtersMap)
						}
					}
					filterMap["Filters"] = filtersSli
					filterMap["Relation"] = filterArg["relation"]
					slsProcessConfigMap["Filter"] = filterMap
				}

				groupBySli := make([]map[string]interface{}, 0)
				if v, ok := slsProcessConfigArg["group_by"]; ok && len(v.(*schema.Set).List()) > 0 {
					for _, groupBy := range v.(*schema.Set).List() {
						groupByMap := make(map[string]interface{})
						groupByMap["Alias"] = groupBy.(map[string]interface{})["alias"]
						groupByMap["SLSKeyName"] = groupBy.(map[string]interface{})["sls_key_name"]
						groupBySli = append(groupBySli, groupByMap)
					}
				}
				slsProcessConfigMap["GroupBy"] = groupBySli

				statisticsSli := make([]map[string]interface{}, 0)
				if v, ok := slsProcessConfigArg["statistics"]; ok && len(v.(*schema.Set).List()) > 0 {
					for _, statistics := range v.(*schema.Set).List() {
						statisticsMap := make(map[string]interface{})
						statisticsMap["Alias"] = statistics.(map[string]interface{})["alias"]
						statisticsMap["Function"] = statistics.(map[string]interface{})["function"]
						statisticsMap["Parameter1"] = statistics.(map[string]interface{})["parameter_one"]
						statisticsMap["Parameter2"] = statistics.(map[string]interface{})["parameter_two"]
						statisticsMap["SLSKeyName"] = statistics.(map[string]interface{})["sls_key_name"]
						statisticsSli = append(statisticsSli, statisticsMap)
					}
				}
				slsProcessConfigMap["Statistics"] = statisticsSli
				request["SLSProcessConfig"] = slsProcessConfigMap
			}
		}
	}

	if update {
		action := "ModifyHybridMonitorTask"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
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
	return resourceAlicloudCmsHybridMonitorSlsTaskRead(d, meta)
}
func resourceAlicloudCmsHybridMonitorSlsTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHybridMonitorTask"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TaskId": d.Id(),
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
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
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}

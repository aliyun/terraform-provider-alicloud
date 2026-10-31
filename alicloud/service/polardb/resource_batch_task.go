package polardb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBatchTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceBatchTaskCreate,
		Read:   resourceBatchTaskRead,
		Delete: resourceBatchTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"task_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"task_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"polarclaw_install_skills", "polarclaw_uninstall_skills"}, false),
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"task_params": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"skill_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"batch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBatchTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}
	action := "CreateBatchTask"
	request := map[string]interface{}{}
	if v, ok := d.GetOk("region_id"); ok {
		request["RegionId"] = v.(string)
	}
	if v, ok := d.GetOk("task_type"); ok {
		request["TaskType"] = v.(string)
	}
	if v, ok := d.GetOk("task_name"); ok {
		request["TaskName"] = v.(string)
	}
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIds := make([]string, 0)
		ids := v.([]interface{})
		for _, id := range ids {
			instanceIds = append(instanceIds, id.(string))
		}
		jsonData, err := json.Marshal(instanceIds)
		if err != nil {
			return errs.WrapError(err)
		}
		request["InstanceIds"] = string(jsonData)
	}
	if v, ok := d.GetOk("task_params"); ok {
		params := v.([]interface{})
		if request["TaskType"] == "polarclaw_install_skills" {
			parameters := make([]map[string]interface{}, 0)
			for _, param := range params {
				item := param.(map[string]interface{})
				param := map[string]interface{}{
					"skillName": item["skill_name"].(string),
				}
				if v, ok := item["version"]; ok && v != nil && v.(string) != "" {
					param["version"] = v.(string)
				}
				parameters = append(parameters, param)
			}
			jsonData, err := json.Marshal(parameters)
			if err != nil {
				return errs.WrapError(err)
			}
			request["Param"] = string(jsonData)
		} else {
			names := make([]string, 0)
			for _, param := range params {
				item := param.(map[string]interface{})
				names = append(names, item["skill_name"].(string))
			}
			jsonData, err := json.Marshal(names)
			if err != nil {
				return errs.WrapError(err)
			}
			request["Param"] = string(jsonData)
		}
	}

	response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		logs.AddDebug(action, response, request)
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_batch_task", action, errs.AlibabaCloudSdkGoERROR)
	}
	batchId, ok := response["BatchId"].(string)
	if !ok || batchId == "" {
		return errs.WrapError(fmt.Errorf("CreateBatchTask returned empty BatchId"))
	}
	d.Set("batch_id", batchId)
	// wait batch task status change from DISPATCHING/RUNNING to COMPLETED
	stateConf := helper.BuildStateConf([]string{"DISPATCHING", "RUNNING"}, []string{"COMPLETED"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, polarDBService.PolarDBBatchTaskStateRefreshFunc(batchId, []string{"FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, IdMsg, batchId)
	}
	d.SetId(batchId)
	return resourceBatchTaskRead(d, meta)
}

func resourceBatchTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}
	applicationAttribute, err := polarDBService.DescribePolarDBBatchTaskAttribute(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	if v, ok := applicationAttribute["TaskType"].(string); ok {
		d.Set("task_type", v)
	}
	if v, ok := applicationAttribute["Status"].(string); ok {
		d.Set("task_status", v)
	}
	return nil
}

func resourceBatchTaskDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

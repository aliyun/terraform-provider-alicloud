package alicloud

import (
	"encoding/json"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBBatchTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBBatchTaskCreate,
		Read:   resourceAlicloudPolarDBBatchTaskRead,
		Delete: resourceAlicloudPolarDBBatchTaskDelete,
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
				ValidateFunc: StringInSlice([]string{"polarclaw_install_skills", "polarclaw_uninstall_skills"}, false),
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

func resourceAlicloudPolarDBBatchTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
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
			return WrapError(err)
		}
		request["InstanceIds"] = string(jsonData)
	}
	if v, ok := d.GetOk("task_params"); ok {
		params := v.([]interface{})
		if request["TaskType"] == "polarclaw_install_skills" {
			parameters := make([]map[string]interface{}, 0)
			for _, param := range params {
				item := param.(map[string]interface{})
				parameters = append(parameters, map[string]interface{}{
					"skillName": item["skill_name"].(string),
					"version":   item["version"].(string),
				})
			}
			jsonData, err := json.Marshal(parameters)
			if err != nil {
				return WrapError(err)
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
				return WrapError(err)
			}
			request["Param"] = string(jsonData)
		}
	}

	response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		addDebug(action, response, request)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_application", action, AlibabaCloudSdkGoERROR)
	}
	if v, ok := response["BatchId"]; ok {
		d.Set("batch_id", v.(string))
	}
	// wait application status change from ClassChanging to running
	stateConf := BuildStateConf([]string{"DISPATCHING", "RUNNING"}, []string{"COMPLETED"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, polarDBService.PolarDBBatchTaskStateRefreshFunc(d.Get("batch_id").(string), []string{"FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Get("batch_id").(string))
	}
	d.SetId(d.Get("batch_id").(string))
	return resourceAlicloudPolarDBBatchTaskRead(d, meta)
}

func resourceAlicloudPolarDBBatchTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	applicationAttribute, err := polarDBService.DescribePolarDBBatchTaskAttribute(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("task_type", applicationAttribute["TaskType"].(string))
	d.Set("task_status", applicationAttribute["Status"].(string))
	return nil
}

func resourceAlicloudPolarDBBatchTaskDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

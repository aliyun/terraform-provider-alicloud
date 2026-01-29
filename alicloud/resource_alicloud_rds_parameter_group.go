package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRdsParameterGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRdsParameterGroupCreate,
		Read:   resourceAliCloudRdsParameterGroupRead,
		Update: resourceAliCloudRdsParameterGroupUpdate,
		Delete: resourceAliCloudRdsParameterGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"modify_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameter_detail": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"param_detail"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"param_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"param_detail": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"param_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Deprecated:    "Field 'param_detail' has been deprecated from provider version 1.270.0. New field 'parameter_detail' instead.",
				ConflictsWith: []string{"parameter_detail"},
			},
			"parameter_group_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameter_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRdsParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateParameterGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["EngineVersion"] = d.Get("engine_version")
	request["Engine"] = d.Get("engine")
	if v, ok := d.GetOk("parameter_group_desc"); ok {
		request["ParameterGroupDesc"] = v
	}
	request["ParameterGroupName"] = d.Get("parameter_group_name")

	parameterMap := map[string]interface{}{}
	if v, ok := d.GetOk("parameter_detail"); ok {
		list := v.(*schema.Set).List()
		for _, v := range list {
			v := v.(map[string]interface{})
			parameterMap[v["param_name"].(string)] = v["param_value"]
		}
	} else if v, ok := d.GetOk("param_detail"); ok {
		list := v.(*schema.Set).List()
		for _, v := range list {
			v := v.(map[string]interface{})
			parameterMap[v["param_name"].(string)] = v["param_value"]
		}
	}
	parameterStr, _ := convertMaptoJsonString(parameterMap)
	request["Parameters"] = parameterStr

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_parameter_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ParameterGroupId"]))

	return resourceAliCloudRdsParameterGroupUpdate(d, meta)
}

func resourceAliCloudRdsParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsServiceV2 := RdsServiceV2{client}

	objectRaw, err := rdsServiceV2.DescribeRdsParameterGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_parameter_group DescribeRdsParameterGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("engine", objectRaw["Engine"])
	d.Set("engine_version", objectRaw["EngineVersion"])
	d.Set("parameter_group_desc", objectRaw["ParameterGroupDesc"])
	d.Set("parameter_group_name", objectRaw["ParameterGroupName"])

	if v, ok := objectRaw["ParamDetail"].(map[string]interface{})["ParameterDetail"].([]interface{}); ok {
		parameterDetail := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			parameterDetail = append(parameterDetail, map[string]interface{}{
				"param_name":  item["ParamName"],
				"param_value": item["ParamValue"],
			})
		}
		d.Set("parameter_detail", parameterDetail)
		d.Set("param_detail", parameterDetail)
	}
	return nil
}

func resourceAliCloudRdsParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyParameterGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ParameterGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("modify_mode"); ok {
		request["ModifyMode"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if !d.IsNewResource() && d.HasChange("parameter_group_desc") {
		update = true
		request["ParameterGroupDesc"] = d.Get("parameter_group_desc")
	}

	if !d.IsNewResource() && d.HasChange("parameter_group_name") {
		update = true
	}
	request["ParameterGroupName"] = d.Get("parameter_group_name")
	if !d.IsNewResource() && d.HasChange("parameter_detail") {
		update = true
	}

	parameterMap := map[string]interface{}{}
	if v, ok := d.GetOk("parameter_detail"); ok {
		list := v.(*schema.Set).List()
		for _, v := range list {
			v := v.(map[string]interface{})
			parameterMap[v["param_name"].(string)] = v["param_value"]
		}
	} else if v, ok := d.GetOk("param_detail"); ok {
		list := v.(*schema.Set).List()
		for _, v := range list {
			v := v.(map[string]interface{})
			parameterMap[v["param_name"].(string)] = v["param_value"]
		}
	}
	parameterStr, _ := convertMaptoJsonString(parameterMap)
	request["Parameters"] = parameterStr

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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

	return resourceAliCloudRdsParameterGroupRead(d, meta)
}

func resourceAliCloudRdsParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteParameterGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ParameterGroupId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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

	return nil
}

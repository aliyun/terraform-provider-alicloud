package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPolarDBParameterGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBParameterGroupCreate,
		Read:   resourceAlicloudPolarDBParameterGroupRead,
		Delete: resourceAlicloudPolarDBParameterGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL"}, false),
			},
			"db_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"5.6", "5.7", "8.0"}, false),
			},
			"parameters": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"param_value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPolarDBParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateParameterGroup"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ParameterGroupName"] = d.Get("name")
	request["DBType"] = d.Get("db_type")
	request["DBVersion"] = d.Get("db_version")
	parametersMap := map[string]interface{}{}
	for _, parametersList := range d.Get("parameters").(*schema.Set).List() {
		parametersArg := parametersList.(map[string]interface{})
		parametersMap[fmt.Sprint(parametersArg["param_name"])] = parametersArg["param_value"]
	}
	parametersJson, err := convertMaptoJsonString(parametersMap)
	if err != nil {
		return WrapError(err)
	}
	request["Parameters"] = parametersJson

	if v, ok := d.GetOk("description"); ok {
		request["ParameterGroupDesc"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_parameter_group", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["ParameterGroupId"]))

	return resourceAlicloudPolarDBParameterGroupRead(d, meta)
}

func resourceAlicloudPolarDBParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	object, err := polarDBService.DescribePolarDBParameterGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object["ParameterGroupName"])
	d.Set("db_type", object["DBType"])
	d.Set("db_version", object["DBVersion"])

	if parameterDetailList, ok := object["ParameterDetail"].([]interface{}); ok {
		parameterDetailMaps := make([]map[string]interface{}, 0)
		for _, parameterDetail := range parameterDetailList {
			parameterDetailMap := map[string]interface{}{}
			parameterDetailArg := parameterDetail.(map[string]interface{})
			parameterDetailMap["param_name"] = parameterDetailArg["ParamName"]
			parameterDetailMap["param_value"] = parameterDetailArg["ParamValue"]
			parameterDetailMaps = append(parameterDetailMaps, parameterDetailMap)
		}
		d.Set("parameters", parameterDetailMaps)
	}

	d.Set("description", object["ParameterGroupDesc"])

	return nil
}

func resourceAlicloudPolarDBParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteParameterGroup"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"ParameterGroupId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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

	return nil
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPolarDbParameterGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPolarDbParameterGroupCreate,
		Read:   resourceAliCloudPolarDbParameterGroupRead,
		Delete: resourceAliCloudPolarDbParameterGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"parameter_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"parameter_group_name", "name"},
				Computed:     true,
				ForceNew:     true,
			},
			"parameters": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"param_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.263.0. New field 'parameter_group_name' instead.",
				ForceNew:   true,
			},
		},
	}
}

func resourceAliCloudPolarDbParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateParameterGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("description"); ok {
		request["ParameterGroupDesc"] = v
	}
	request["DBVersion"] = d.Get("db_version")
	request["DBType"] = d.Get("db_type")
	_, err = jsonpath.Get("$", d.Get("parameters"))
	if err == nil {
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
	}

	if v, ok := d.GetOk("name"); ok || d.HasChange("name") {
		request["ParameterGroupName"] = v
	}

	if v, ok := d.GetOk("parameter_group_name"); ok {
		request["ParameterGroupName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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

	return resourceAliCloudPolarDbParameterGroupRead(d, meta)
}

func resourceAliCloudPolarDbParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}

	objectRaw, err := polarDbServiceV2.DescribePolarDbParameterGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_parameter_group DescribePolarDbParameterGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("db_type", objectRaw["DBType"])
	d.Set("db_version", objectRaw["DBVersion"])
	d.Set("description", objectRaw["ParameterGroupDesc"])
	d.Set("parameter_group_name", objectRaw["ParameterGroupName"])

	parameterDetailRaw := objectRaw["ParameterDetail"]
	parametersMaps := make([]map[string]interface{}, 0)
	if parameterDetailRaw != nil {
		for _, parameterDetailChildRaw := range convertToInterfaceArray(parameterDetailRaw) {
			parametersMap := make(map[string]interface{})
			parameterDetailChildRaw := parameterDetailChildRaw.(map[string]interface{})
			parametersMap["param_name"] = parameterDetailChildRaw["ParamName"]
			parametersMap["param_value"] = parameterDetailChildRaw["ParamValue"]

			parametersMaps = append(parametersMaps, parametersMap)
		}
	}
	if err := d.Set("parameters", parametersMaps); err != nil {
		return err
	}

	d.Set("name", d.Get("parameter_group_name"))
	return nil
}

func resourceAliCloudPolarDbParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {

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
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"ParamGroupsNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

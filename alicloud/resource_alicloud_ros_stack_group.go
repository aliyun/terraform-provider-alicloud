// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRosStackGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRosStackGroupCreate,
		Read:   resourceAliCloudRosStackGroupRead,
		Update: resourceAliCloudRosStackGroupUpdate,
		Delete: resourceAliCloudRosStackGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"administration_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_deployment": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"retain_stacks_on_account_removal": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execution_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameter_key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"permission_model": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"SELF_MANAGED", "SERVICE_MANAGED"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stack_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stack_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"template_body": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_ids": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `account_ids` has been removed from provider version 1.257.0.",
			},
			"operation_description": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `operation_description` has been removed from provider version 1.257.0.",
			},
			"operation_preferences": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `operation_preferences` has been removed from provider version 1.257.0.",
			},
			"region_ids": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `region_ids` has been removed from provider version 1.257.0.",
			},
		},
	}
}

func resourceAliCloudRosStackGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateStackGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("stack_group_name"); ok {
		request["StackGroupName"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("template_id"); ok {
		request["TemplateId"] = v
	}
	if v, ok := d.GetOk("execution_role_name"); ok {
		request["ExecutionRoleName"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("template_version"); ok {
		request["TemplateVersion"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("auto_deployment"); !IsNil(v) {
		retainStacksOnAccountRemoval1, _ := jsonpath.Get("$[0].retain_stacks_on_account_removal", v)
		if retainStacksOnAccountRemoval1 != nil && retainStacksOnAccountRemoval1 != "" {
			objectDataLocalMap["RetainStacksOnAccountRemoval"] = retainStacksOnAccountRemoval1
		}
		enabled1, _ := jsonpath.Get("$[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			objectDataLocalMap["Enabled"] = enabled1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["AutoDeployment"] = string(objectDataLocalMapJson)
	}

	if v, ok := d.GetOk("parameters"); ok {
		parametersMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			if dataLoop1 == nil {
				continue
			}
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["ParameterValue"] = dataLoop1Tmp["parameter_value"]
			dataLoop1Map["ParameterKey"] = dataLoop1Tmp["parameter_key"]
			parametersMapsArray = append(parametersMapsArray, dataLoop1Map)
		}
		request["Parameters"] = parametersMapsArray
	}

	if v, ok := d.GetOk("capabilities"); ok {
		capabilitiesMapsArray := v.([]interface{})
		request["Capabilities"] = capabilitiesMapsArray
	}

	if v, ok := d.GetOk("template_body"); ok {
		request["TemplateBody"] = v
	}
	if v, ok := d.GetOk("permission_model"); ok {
		request["PermissionModel"] = v
	}
	if v, ok := d.GetOk("template_url"); ok {
		request["TemplateURL"] = v
	}
	if v, ok := d.GetOk("administration_role_name"); ok {
		request["AdministrationRoleName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ROS", "2019-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ros_stack_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["StackGroupName"]))

	return resourceAliCloudRosStackGroupUpdate(d, meta)
}

func resourceAliCloudRosStackGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosServiceV2 := RosServiceV2{client}

	objectRaw, err := rosServiceV2.DescribeRosStackGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ros_stack_group DescribeRosStackGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("administration_role_name", objectRaw["AdministrationRoleName"])
	d.Set("description", objectRaw["Description"])
	d.Set("execution_role_name", objectRaw["ExecutionRoleName"])
	d.Set("permission_model", objectRaw["PermissionModel"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("stack_group_id", objectRaw["StackGroupId"])
	d.Set("status", objectRaw["Status"])

	autoDeploymentMaps := make([]map[string]interface{}, 0)
	autoDeploymentMap := make(map[string]interface{})
	autoDeploymentRaw := make(map[string]interface{})
	if objectRaw["AutoDeployment"] != nil {
		autoDeploymentRaw = objectRaw["AutoDeployment"].(map[string]interface{})
	}
	if len(autoDeploymentRaw) > 0 {
		autoDeploymentMap["enabled"] = autoDeploymentRaw["Enabled"]
		autoDeploymentMap["retain_stacks_on_account_removal"] = autoDeploymentRaw["RetainStacksOnAccountRemoval"]

		autoDeploymentMaps = append(autoDeploymentMaps, autoDeploymentMap)
	}
	if err := d.Set("auto_deployment", autoDeploymentMaps); err != nil {
		return err
	}
	parametersRaw := objectRaw["Parameters"]
	parametersMaps := make([]map[string]interface{}, 0)
	if parametersRaw != nil {
		for _, parametersChildRaw := range parametersRaw.([]interface{}) {
			parametersMap := make(map[string]interface{})
			parametersChildRaw := parametersChildRaw.(map[string]interface{})
			parametersMap["parameter_key"] = parametersChildRaw["ParameterKey"]
			parametersMap["parameter_value"] = parametersChildRaw["ParameterValue"]

			parametersMaps = append(parametersMaps, parametersMap)
		}
	}
	if err := d.Set("parameters", parametersMaps); err != nil {
		return err
	}

	objectRaw, err = rosServiceV2.DescribeStackGroupListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = rosServiceV2.DescribeStackGroupGetTemplate(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("template_body", objectRaw["TemplateBody"])

	d.Set("stack_group_name", d.Id())

	return nil
}

func resourceAliCloudRosStackGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateStackGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["StackGroupName"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("template_id") {
		update = true
		request["TemplateId"] = d.Get("template_id")
	}

	if !d.IsNewResource() && d.HasChange("execution_role_name") {
		update = true
		request["ExecutionRoleName"] = d.Get("execution_role_name")
	}

	if !d.IsNewResource() && d.HasChange("template_version") {
		update = true
		request["TemplateVersion"] = d.Get("template_version")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("auto_deployment") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("auto_deployment"); v != nil {
			retainStacksOnAccountRemoval1, _ := jsonpath.Get("$[0].retain_stacks_on_account_removal", v)
			if retainStacksOnAccountRemoval1 != nil && (d.HasChange("auto_deployment.0.retain_stacks_on_account_removal") || retainStacksOnAccountRemoval1 != "") {
				objectDataLocalMap["RetainStacksOnAccountRemoval"] = retainStacksOnAccountRemoval1
			}
			enabled1, _ := jsonpath.Get("$[0].enabled", v)
			if enabled1 != nil && (d.HasChange("auto_deployment.0.enabled") || enabled1 != "") {
				objectDataLocalMap["Enabled"] = enabled1
			}

			objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
			if err != nil {
				return WrapError(err)
			}
			request["AutoDeployment"] = string(objectDataLocalMapJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("parameters") {
		update = true
		if v, ok := d.GetOk("parameters"); ok || d.HasChange("parameters") {
			parametersMapsArray := make([]interface{}, 0)
			for _, dataLoop := range v.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["ParameterValue"] = dataLoopTmp["parameter_value"]
				dataLoopMap["ParameterKey"] = dataLoopTmp["parameter_key"]
				parametersMapsArray = append(parametersMapsArray, dataLoopMap)
			}
			request["Parameters"] = parametersMapsArray
		}
	}

	if !d.IsNewResource() && d.HasChange("template_body") {
		update = true
		request["TemplateBody"] = d.Get("template_body")
	}

	if !d.IsNewResource() && d.HasChange("permission_model") {
		update = true
		request["PermissionModel"] = d.Get("permission_model")
	}

	if !d.IsNewResource() && d.HasChange("template_url") {
		update = true
		request["TemplateURL"] = d.Get("template_url")
	}

	if !d.IsNewResource() && d.HasChange("administration_role_name") {
		update = true
		request["AdministrationRoleName"] = d.Get("administration_role_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ROS", "2019-09-10", action, query, request, true)
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
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "stackgroup"
	if d.HasChange("stack_group_id") {
		update = true
	}
	request["ResourceId"] = d.Get("stack_group_id")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ROS", "2019-09-10", action, query, request, true)
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

	if d.HasChange("tags") {
		rosServiceV2 := RosServiceV2{client}
		if err := rosServiceV2.SetResourceTags(d, "stackgroup"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudRosStackGroupRead(d, meta)
}

func resourceAliCloudRosStackGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteStackGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["StackGroupName"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ROS", "2019-09-10", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"StackGroupNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

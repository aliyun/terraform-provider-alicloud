package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBApplicationCreate,
		Read:   resourceAlicloudPolarDBApplicationRead,
		Update: resourceAlicloudPolarDBApplicationUpdate,
		Delete: resourceAlicloudPolarDBApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"application_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"components": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_type": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"component_class": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"component_replica": {
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},
			"pay_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				ForceNew:     true,
				Optional:     true,
				Default:      PostPaid,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: polardbPostPaidDiffSuppressFunc,
			},
			"used_time": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"model_from": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ai_db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"model_api_key": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"model_base_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"model_api": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"model_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_version": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"component_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"security_ip_list": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_ip_array_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modify_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudPolarDBApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	request, err := buildApplicationCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	action := "CreateApplication"

	// Execute CreateApplication only once, no retry
	response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		addDebug(action, response, request)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_application", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["ApplicationId"]))
	// wait application status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDBService.PolarDBApplicationStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudPolarDBApplicationRead(d, meta)
}

func resourceAlicloudPolarDBApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	//wait := incrementalWait(3*time.Second, 3*time.Second)
	if d.HasChange("upgrade_version") {
		action := "UpgradeApplicationVersion"
		request := map[string]interface{}{
			"ApplicationId": d.Id(),
		}
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				addDebug(action, response, request)
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"MinorVersionUpgrading"}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, polarDBService.PolarDBApplicationStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	if d.HasChanges("parameters") {
		action := "ModifyApplicationParameter"
		request := map[string]interface{}{
			"ApplicationId": d.Id(),
		}
		if v, ok := d.GetOk("parameters"); ok {
			parameters := make([]map[string]interface{}, 0)
			params := v.([]interface{})
			for _, param := range params {

				paramMap := make(map[string]interface{})
				item := param.(map[string]interface{})

				paramMap["ParameterName"] = item["parameter_name"].(string)
				paramMap["ParameterValue"] = item["parameter_value"].(string)

				parameters = append(parameters, paramMap)
			}
			jsonData, err := json.Marshal(parameters)
			if err != nil {
				return WrapError(err)
			}
			request["Parameters"] = string(jsonData)
		}

		response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			addDebug(action, response, request)
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_application", action, AlibabaCloudSdkGoERROR)
		}
		// wait application status change from ClassChanging to running
		stateConf := BuildStateConf([]string{"ClassChanging", "Restarting"}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, polarDBService.PolarDBApplicationStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChanges("security_ip_list", "security_ip_array_name", "modify_mode", "component_id", "security_groups") {
		action := "ModifyApplicationWhitelist"
		modifyRequest := map[string]interface{}{
			"ApplicationId": d.Id(),
		}
		if _, ok := d.GetOk("modify_mode"); ok {
			modifyRequest["ModifyMode"] = d.Get("modify_mode").(string)
		}
		if _, ok := d.GetOk("security_ip_array_name"); ok {
			modifyRequest["SecurityIPArrayName"] = d.Get("security_ip_array_name").(string)
		}
		if _, ok := d.GetOk("security_ip_list"); ok {
			modifyRequest["SecurityIPList"] = d.Get("security_ip_list").(string)
		}
		if _, ok := d.GetOk("security_groups"); ok {
			modifyRequest["SecurityGroups"] = d.Get("security_groups").(string)
		}
		if _, ok := d.GetOk("component_id"); ok {
			modifyRequest["ComponentId"] = d.Get("component_id").(string)
		}

		response, err := client.RpcPost("polardb", "2017-08-01", action, nil, modifyRequest, false)

		if err != nil {
			addDebug(action, response, modifyRequest)
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)
	return resourceAlicloudPolarDBApplicationRead(d, meta)
}

func resourceAlicloudPolarDBApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	applicationAttribute, err := polarDBService.DescribePolarDBApplicationAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("region_id", applicationAttribute["RegionId"].(string))
	d.Set("zone_id", applicationAttribute["ZoneId"].(string))
	d.Set("vswitch_id", applicationAttribute["VSwitchId"].(string))
	d.Set("pay_type", getChargeType(applicationAttribute["PayType"].(string)))
	d.Set("vpc_id", applicationAttribute["VPCId"].(string))
	if v, ok := applicationAttribute["Description"]; ok {
		d.Set("description", v.(string))
	}
	d.Set("application_type", applicationAttribute["ApplicationType"].(string))
	components := applicationAttribute["Components"].([]interface{})
	componentsList := make([]map[string]interface{}, 0, len(components))
	for i, _ := range components {
		componentsMap := components[i].(map[string]interface{})
		component := make(map[string]interface{})
		if v, ok := componentsMap["ComponentType"]; ok {
			componentType := v.(string)
			component["component_type"] = convertComponentTypeReadResponse(componentType)
		}
		if v, ok := componentsMap["ComponentClass"]; ok {
			component["component_class"] = v.(string)
		}
		if v, ok := componentsMap["ComponentReplica"]; ok {
			if num, ok := v.(json.Number); ok {
				if replica, err := num.Int64(); err == nil {
					component["component_replica"] = int(replica)
				}
			} else if replica, ok := v.(int); ok {
				component["component_replica"] = replica
			}
		}
		componentsList = append(componentsList, component)
	}
	d.Set("components", componentsList)
	d.Set("status", applicationAttribute["Status"].(string))
	if v, ok := applicationAttribute["DBClusterId"]; ok {
		d.Set("db_cluster_id", v.(string))
	}
	return nil
}

func resourceAlicloudPolarDBApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	app, err := polarDBService.DescribePolarDBApplicationAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	// Pre paid application can not be release.
	if PayType(app["PayType"].(string)) == Prepaid {
		return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}
	_, err = polarDBService.DeletePolarDBApplication(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteApplication", AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, polarDBService.PolarDBApplicationStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		if strings.HasPrefix(err.Error(), "couldn't find resource") {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, IdMsg, d.Id())
	}
	d.SetId("")
	return nil
}

func buildApplicationCreateRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"Description":     Trim(d.Get("description").(string)),
		"ApplicationType": Trim(d.Get("application_type").(string)),
		"Architecture":    Trim(d.Get("architecture").(string)),
		"ClientToken":     buildClientToken("CreateApplication"),
	}

	if v, ok := d.GetOk("db_cluster_id"); ok && v.(string) != "" {
		request["DBClusterId"] = v.(string)
	}

	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request["ResourceGroupId"] = v.(string)
	}

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request["ZoneId"] = Trim(zone.(string))
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v.(string)
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v.(string)
	}

	if request["VSwitchId"] != nil {
		if request["ZoneId"] == nil {
			// check vswitchId in zone
			vsw, err := vpcService.DescribeVSwitch(request["VSwitchId"].(string))
			if err != nil {
				return nil, WrapError(err)
			}

			if v, ok := request["ZoneId"].(string); !ok || v == "" {
				request["ZoneId"] = vsw.ZoneId
			} else if request["ZoneId"] != vsw.ZoneId {
				return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request["ZoneId"]))
			}
		}
	}

	if v, ok := d.GetOk("components"); ok {
		components := make([]map[string]interface{}, 0)
		comps := v.([]interface{})
		for _, comp := range comps {
			item := comp.(map[string]interface{})
			components = append(components, map[string]interface{}{
				"ComponentType":    item["component_type"].(string),
				"ComponentClass":   item["component_class"].(string),
				"ComponentReplica": item["component_replica"].(int),
			})
		}
		jsonData, err := json.Marshal(components)
		if err != nil {
			return nil, WrapError(err)
		}
		request["Components"] = string(jsonData)
	}

	payType := Trim(d.Get("pay_type").(string))
	request["PayType"] = string(Postpaid)
	if payType == string(PrePaid) {
		request["PayType"] = string(Prepaid)
	}
	if PayType(request["PayType"].(string)) == Prepaid {
		period := d.Get("period").(int)
		request["UsedTime"] = strconv.Itoa(period)
		request["Period"] = string(Month)
		if period > 9 {
			request["UsedTime"] = strconv.Itoa(period / 12)
			request["Period"] = string(Year)
		}
		request["AutoRenew"] = d.Get("auto_renew").(bool)
	}

	if v, ok := d.GetOk("security_ip_list"); ok {
		request["SecurityIPList"] = v.(string)
	}

	if v, ok := d.GetOk("model_from"); ok {
		request["ModelFrom"] = v.(string)
	}

	if v, ok := d.GetOk("ai_db_cluster_id"); ok {
		request["AIDBClusterId"] = v.(string)
	}

	if v, ok := d.GetOk("model_api_key"); ok {
		request["ModelApiKey"] = v.(string)
	}

	if v, ok := d.GetOk("model_base_url"); ok {
		request["ModelBaseUrl"] = v.(string)
	}

	if v, ok := d.GetOk("model_api"); ok {
		request["ModelApi"] = v.(string)
	}

	if v, ok := d.GetOk("model_name"); ok {
		request["ModelName"] = v.(string)
	}

	return request, nil
}

func convertComponentTypeReadResponse(source string) string {
	switch source {
	case "polarclaw":
		return "polarclaw_comp"
	}
	return ""
}

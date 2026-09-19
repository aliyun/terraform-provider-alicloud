package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudDataWorksTenantResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksTenantResourceGroupCreate,
		Read:   resourceAliCloudDataWorksTenantResourceGroupRead,
		Update: resourceAliCloudDataWorksTenantResourceGroupUpdate,
		Delete: resourceAliCloudDataWorksTenantResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"tenant_resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_resource_group_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aliyun_resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"payment_duration_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			},
			"auto_renew_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"spec": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchemaForceNew(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_amount": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"spec_standard": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDataWorksTenantResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	request["Name"] = d.Get("tenant_resource_group_name")
	request["PaymentType"] = d.Get("payment_type")
	request["VpcId"] = d.Get("vpc_id")
	request["VswitchId"] = d.Get("vswitch_id")

	if v, ok := d.GetOk("tenant_resource_group_description"); ok {
		request["Remark"] = v
	}
	if v, ok := d.GetOk("aliyun_resource_group_id"); ok {
		request["AliyunResourceGroupId"] = v
	}
	if v, ok := d.GetOk("payment_duration"); ok {
		request["PaymentDuration"] = v
	}
	if v, ok := d.GetOk("payment_duration_unit"); ok {
		request["PaymentDurationUnit"] = v
	}
	if v, ok := d.GetOk("auto_renew_enabled"); ok {
		request["AutoRenewEnabled"] = v
	}
	if v, ok := d.GetOk("spec"); ok {
		request["Spec"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["AliyunResourceTags"] = tagsMap
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, query, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_tenant_resource_group", action, AlibabaCloudSdkGoERROR)
	}

	rgOrder, ok := response["ResourceGroupOrder"].(map[string]interface{})
	if !ok {
		return WrapError(Error("CreateResourceGroup response missing ResourceGroupOrder"))
	}
	id, ok := rgOrder["Id"].(string)
	if !ok || id == "" {
		return WrapError(Error("CreateResourceGroup response missing ResourceGroupOrder.Id"))
	}
	d.SetId(id)

	// Wait for the resource group to leave the Creating state so that Read
	// can populate the full set of attributes (DefaultVpcId, AliyunResourceGroupId,
	// AliyunResourceTags, etc. are only returned after the resource reaches Normal).
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Creating"},
		Target:     []string{"Normal"},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Refresh:    dataWorksTenantResourceGroupStatusRefreshFunc(client, d.Id()),
	}
	if _, e := stateConf.WaitForState(); e != nil {
		return WrapErrorf(e, DefaultErrorMsg, d.Id(), "WaitForNormalState", AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudDataWorksTenantResourceGroupRead(d, meta)
}

func resourceAliCloudDataWorksTenantResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksTenantResourceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_tenant_resource_group DescribeDataWorksTenantResourceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if v, ok := objectRaw["Name"]; ok && v != nil {
		d.Set("tenant_resource_group_name", v)
	}
	if v, ok := objectRaw["Remark"]; ok && v != nil {
		d.Set("tenant_resource_group_description", v)
	}
	if v, ok := objectRaw["PaymentType"]; ok && v != nil {
		d.Set("payment_type", v)
	}
	if v, ok := objectRaw["DefaultVpcId"]; ok && v != nil {
		d.Set("vpc_id", v)
	}
	if v, ok := objectRaw["DefaultVswitchId"]; ok && v != nil {
		d.Set("vswitch_id", v)
	}
	if v, ok := objectRaw["AliyunResourceGroupId"]; ok && v != nil {
		d.Set("aliyun_resource_group_id", v)
	}
	if v, ok := objectRaw["Status"]; ok && v != nil {
		d.Set("status", fmt.Sprint(v))
	}
	if v, ok := objectRaw["CreateTime"]; ok && v != nil {
		d.Set("create_time", fmt.Sprint(v))
	}
	if v, ok := objectRaw["ResourceGroupType"]; ok && v != nil {
		d.Set("resource_group_type", fmt.Sprint(v))
	}
	if v, ok := objectRaw["OrderInstanceId"]; ok && v != nil {
		d.Set("order_instance_id", fmt.Sprint(v))
	}
	if v, ok := objectRaw["CreateUser"]; ok && v != nil {
		d.Set("create_user", fmt.Sprint(v))
	}

	if specRaw, ok := objectRaw["Spec"]; ok && specRaw != nil {
		if specMap, ok := specRaw.(map[string]interface{}); ok {
			if v, ok := specMap["Amount"]; ok && v != nil {
				if amount, err := convertToInt(v); err == nil {
					d.Set("spec_amount", amount)
				}
			}
			if v, ok := specMap["Standard"]; ok && v != nil {
				d.Set("spec_standard", fmt.Sprint(v))
			}
		}
	}

	if tagsRaw, ok := objectRaw["AliyunResourceTags"]; ok && tagsRaw != nil {
		d.Set("tags", tagsToMapFromKeyValues(tagsRaw))
	}

	return nil
}

func resourceAliCloudDataWorksTenantResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	if d.HasChange("tenant_resource_group_description") {
		action := "UpdateResourceGroup"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["Id"] = d.Id()

		if !d.IsNewResource() && d.HasChange("tenant_resource_group_description") {
			request["Remark"] = d.Get("tenant_resource_group_description")
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2020-05-18", action, query, request, false)
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
		d.SetPartial("tenant_resource_group_description")
	}

	d.Partial(false)
	return resourceAliCloudDataWorksTenantResourceGroupRead(d, meta)
}

func resourceAliCloudDataWorksTenantResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["Id"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, query, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"704203"}) {
				// 704203: resource group is still in CREATING state, retry until it reaches NORMAL.
				wait()
				return resource.RetryableError(err)
			}
			if NotFoundError(err) || IsExpectedErrors(err, []string{"1203110072"}) {
				return nil
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

func dataWorksTenantResourceGroupStatusRefreshFunc(client *connectivity.AliyunClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		dataWorksServiceV2 := DataWorksServiceV2{client}
		object, err := dataWorksServiceV2.DescribeDataWorksTenantResourceGroup(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "Gone", nil
			}
			return nil, "", WrapError(err)
		}
		status, ok := object["Status"].(string)
		if !ok {
			return nil, "", WrapError(Error("DescribeDataWorksTenantResourceGroup response missing Status"))
		}
		return object, status, nil
	}
}

func convertToInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int64:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		i, e := strconv.Atoi(val)
		return i, e
	default:
		i, e := strconv.Atoi(fmt.Sprint(v))
		return i, e
	}
}

func tagsToMapFromKeyValues(tags interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if tags == nil {
		return result
	}
	tagList, ok := tags.([]interface{})
	if !ok {
		return result
	}
	for _, t := range tagList {
		tagMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		key := ""
		if k, ok := tagMap["Key"]; ok {
			key = fmt.Sprint(k)
		}
		value := ""
		if v, ok := tagMap["Value"]; ok {
			value = fmt.Sprint(v)
		}
		if key != "" {
			result[key] = value
		}
	}
	return result
}

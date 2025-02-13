package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudReservedInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudReservedInstanceCreate,
		Read:   resourceAliCloudReservedInstanceRead,
		Update: resourceAliCloudReservedInstanceUpdate,
		Delete: resourceAliCloudReservedInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allocation_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"renewal_status": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(RenewAutoRenewal), string(RenewNormal)}, false),
			},
			"auto_renew_period": {
				Computed:     true,
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 12, 36, 60}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("renewal_status").(string) == string(RenewAutoRenewal) {
						return false
					}
					return true
				},
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"expired_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"operation_locks": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lock_reason": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"reserved_instance_name": {
				Optional:      true,
				Computed:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"name"},
			},
			"start_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"tags": tagsSchema(),
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Region", "Zone"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Windows", "Linux"}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 5}),
			},
			"offering_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"No Upfront", "Partial Upfront", "All Upfront"}, false),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.194.0. New field 'reserved_instance_name' instead.",
				ConflictsWith: []string{"reserved_instance_name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}
func resourceAliCloudReservedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("instance_amount"); ok {
		request["InstanceAmount"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("offering_type"); ok {
		request["OfferingType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}
	if v, ok := d.GetOk("reserved_instance_name"); ok {
		request["ReservedInstanceName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["ReservedInstanceName"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if scope, ok := d.GetOk("scope"); ok {
		request["Scope"] = scope
		if v, ok := d.GetOk("zone_id"); ok {
			if scope == "Zone" && v == "" {
				return WrapError(Error("Required when Scope is Zone."))
			}
			request["ZoneId"] = v
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}

	if v, ok := d.GetOk("renewal_status"); ok {
		request["AutoRenew"] = v.(string) == string(RenewAutoRenewal)
	}
	if v, ok := d.GetOkExists("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}

	request["ClientToken"] = buildClientToken("PurchaseReservedInstancesOffering")
	var response map[string]interface{}
	action := "PurchaseReservedInstancesOffering"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp

		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_reserved_instance", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.ReservedInstanceIdSets.ReservedInstanceId[0]", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ecs_reserved_instance")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsReservedInstanceStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudReservedInstanceRead(d, meta)
}
func resourceAliCloudReservedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	ecsService := EcsService{client}
	d.Partial(true)
	if d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "reservedinstance"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	request := ecs.CreateModifyReservedInstanceAttributeRequest()
	request.ReservedInstanceId = d.Id()
	request.RegionId = client.RegionId
	if d.HasChange("name") {
		update = true
		if v, ok := d.GetOk("name"); ok {
			request.ReservedInstanceName = v.(string)
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request.Description = v.(string)
		}
	}
	if d.HasChange("reserved_instance_name") {
		update = true
		if v, ok := d.GetOk("reserved_instance_name"); ok {
			request.ReservedInstanceName = v.(string)
		}
	}

	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyReservedInstanceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		d.SetPartial("name")
		d.SetPartial("reserved_instance_name")
		d.SetPartial("description")
	}

	if d.HasChanges("auto_renew_period", "renewal_status") {
		request := map[string]interface{}{
			"ReservedInstanceId": []string{d.Id()},
			"RegionId":           client.RegionId,
		}
		if v, ok := d.GetOk("auto_renew_period"); ok {
			if formatInt(v) == 1 {
				request["Period"] = v
				request["PeriodUnit"] = "Month"
			} else {
				request["Period"] = formatInt(v) / 12
				request["PeriodUnit"] = "Year"
			}
		}
		if v, ok := d.GetOk("renewal_status"); ok {
			request["RenewalStatus"] = v
		}

		action := "ModifyReservedInstanceAutoRenewAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("auto_renew_period")
		d.SetPartial("renewal_status")
	}

	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsReservedInstanceStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudReservedInstanceRead(d, meta)
}
func resourceAliCloudReservedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	reservedInstances, err := ecsService.DescribeReservedInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_type", reservedInstances.InstanceType)
	d.Set("allocation_status", reservedInstances.AllocationStatus)
	d.Set("scope", reservedInstances.Scope)
	d.Set("zone_id", reservedInstances.ZoneId)
	d.Set("instance_amount", reservedInstances.InstanceAmount)
	d.Set("platform", strings.Title(reservedInstances.Platform))
	d.Set("offering_type", reservedInstances.OfferingType)
	d.Set("name", reservedInstances.ReservedInstanceName)
	d.Set("reserved_instance_name", reservedInstances.ReservedInstanceName)
	d.Set("description", reservedInstances.Description)
	d.Set("resource_group_id", reservedInstances.ReservedInstanceId)
	d.Set("status", reservedInstances.Status)
	d.Set("create_time", reservedInstances.CreationTime)
	d.Set("expired_time", reservedInstances.ExpiredTime)
	d.Set("start_time", reservedInstances.StartTime)

	operationLocks := make([]map[string]interface{}, 0)
	for _, lock := range reservedInstances.OperationLocks.OperationLock {
		operationLocks = append(operationLocks, map[string]interface{}{
			"lock_reason": lock.LockReason,
		})
	}
	d.Set("operation_locks", operationLocks)
	tags := make(map[string]string)
	for _, t := range reservedInstances.Tags.Tag {
		tags[t.TagKey] = t.TagValue
	}
	d.Set("tags", tags)

	object, err := ecsService.DescribeReservedInstanceAutoRenewAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("renewal_status", object["RenewalStatus"])

	if v, ok := object["Duration"]; ok && formatInt(v) != 0 {
		renewPeriod := formatInt(v)
		if object["PeriodUnit"] == string(Year) {
			renewPeriod = renewPeriod * 12
		}
		d.Set("auto_renew_period", renewPeriod)
	}

	return nil
}
func resourceAliCloudReservedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// PurchaseReservedInstancesOffering can not be release.
	return nil
}

package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEcsDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsDedicatedHostCreate,
		Read:   resourceAlicloudEcsDedicatedHostRead,
		Update: resourceAlicloudEcsDedicatedHostUpdate,
		Delete: resourceAlicloudEcsDedicatedHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action_on_maintenance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Stop", "Migrate"}, false),
				Default:      "Stop",
			},
			"auto_placement": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"false", "on"}, false),
				Default:      "on",
			},
			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"dedicated_host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detail_fee": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"expired_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_attributes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"udp_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"slb_udp_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				Default:      "PostPaid",
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sale_cycle": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcsDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateAllocateDedicatedHostsRequest()
	if v, ok := d.GetOk("action_on_maintenance"); ok {
		request.ActionOnMaintenance = v.(string)
	}
	if v, ok := d.GetOk("auto_placement"); ok {
		request.AutoPlacement = v.(string)
	}
	if v, ok := d.GetOk("auto_release_time"); ok {
		request.AutoReleaseTime = v.(string)
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request.AutoRenewPeriod = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("dedicated_host_name"); ok {
		request.DedicatedHostName = v.(string)
	}
	request.DedicatedHostType = d.Get("dedicated_host_type").(string)
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if _, ok := d.GetOk("expired_time"); ok {
		if v, err := strconv.Atoi(d.Get("expired_time").(string)); err == nil {
			request.Period = requests.NewInteger(v)
		} else {
			return WrapError(err)
		}
	}
	if v, ok := d.GetOk("network_attributes"); ok {
		network_attributes := v.(*schema.Set).List()
		for _, value := range network_attributes {
			arg := value.(map[string]interface{})
			request.NetworkAttributesSlbUdpTimeout = requests.NewInteger(arg["slb_udp_timeout"].(int))
			request.NetworkAttributesUdpTimeout = requests.NewInteger(arg["udp_timeout"].(int))
		}
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request.ChargeType = v.(string)
	}
	request.Quantity = requests.NewInteger(1)
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("sale_cycle"); ok {
		request.PeriodUnit = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		addTags := make([]ecs.AllocateDedicatedHostsTag, 0)
		for key, value := range v.(map[string]interface{}) {
			addTags = append(addTags, ecs.AllocateDedicatedHostsTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &addTags
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = v.(string)
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.AllocateDedicatedHosts(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_dedicated_host", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ecs.AllocateDedicatedHostsResponse)
	d.SetId(fmt.Sprintf("%v", response.DedicatedHostIdSets.DedicatedHostId[0]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{"PermanentFailure"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcsDedicatedHostUpdate(d, meta)
}
func resourceAlicloudEcsDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDedicatedHost(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("action_on_maintenance", object.ActionOnMaintenance)
	d.Set("auto_placement", object.AutoPlacement)
	d.Set("auto_release_time", object.AutoReleaseTime)
	d.Set("dedicated_host_name", object.DedicatedHostName)
	d.Set("dedicated_host_type", object.DedicatedHostType)
	d.Set("description", object.Description)
	d.Set("expired_time", object.ExpiredTime)
	d.Set("payment_type", object.ChargeType)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("sale_cycle", object.SaleCycle)
	d.Set("status", object.Status)
	d.Set("zone_id", object.ZoneId)
	return nil
}
func resourceAlicloudEcsDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "ddh"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := ecs.CreateModifyDedicatedHostAutoReleaseTimeRequest()
	request.DedicatedHostId = d.Id()
	request.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("auto_release_time") {
		update = true
		request.AutoReleaseTime = d.Get("auto_release_time").(string)
	}
	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyDedicatedHostAutoReleaseTime(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("auto_release_time")
	}
	update = false
	joinResourceGroupReq := ecs.CreateJoinResourceGroupRequest()
	joinResourceGroupReq.ResourceId = d.Id()
	joinResourceGroupReq.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		joinResourceGroupReq.ResourceGroupId = d.Get("resource_group_id").(string)
	}
	joinResourceGroupReq.ResourceType = "ddh"
	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.JoinResourceGroup(joinResourceGroupReq)
		})
		addDebug(joinResourceGroupReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), joinResourceGroupReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}
	update = false
	renewDedicatedHostsReq := ecs.CreateRenewDedicatedHostsRequest()
	renewDedicatedHostsReq.DedicatedHostIds = d.Id()
	if !d.IsNewResource() && d.HasChange("expired_time") {
		update = true
		if v, err := strconv.Atoi(d.Get("expired_time").(string)); err == nil {
			renewDedicatedHostsReq.Period = requests.NewInteger(v)
		} else {
			return WrapError(err)
		}
	}
	renewDedicatedHostsReq.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("sale_cycle") {
		update = true
		renewDedicatedHostsReq.PeriodUnit = d.Get("sale_cycle").(string)
	}
	if update && d.Get("charge_type").(string) == string(PrePaid) && !d.HasChange("charge_type") {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.RenewDedicatedHosts(renewDedicatedHostsReq)
		})
		addDebug(renewDedicatedHostsReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), renewDedicatedHostsReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("expired_time")
		d.SetPartial("sale_cycle")
	}
	update = false
	modifyDedicatedHostAttributeReq := ecs.CreateModifyDedicatedHostAttributeRequest()
	modifyDedicatedHostAttributeReq.DedicatedHostId = d.Id()
	modifyDedicatedHostAttributeReq.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("action_on_maintenance") {
		update = true
		modifyDedicatedHostAttributeReq.ActionOnMaintenance = d.Get("action_on_maintenance").(string)
	}
	if !d.IsNewResource() && d.HasChange("auto_placement") {
		update = true
		modifyDedicatedHostAttributeReq.AutoPlacement = d.Get("auto_placement").(string)
	}
	if !d.IsNewResource() && d.HasChange("dedicated_host_name") {
		update = true
		modifyDedicatedHostAttributeReq.DedicatedHostName = d.Get("dedicated_host_name").(string)
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		modifyDedicatedHostAttributeReq.Description = d.Get("description").(string)
	}
	if !d.IsNewResource() && d.HasChange("network_attributes") {
		update = true
		for _, value := range d.Get("network_attributes").(*schema.Set).List() {
			arg := value.(map[string]interface{})
			modifyDedicatedHostAttributeReq.NetworkAttributesUdpTimeout = requests.NewInteger(arg["udp_timeout"].(int))
			modifyDedicatedHostAttributeReq.NetworkAttributesSlbUdpTimeout = requests.NewInteger(arg["slb_udp_timeout"].(int))
		}
	}
	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyDedicatedHostAttribute(modifyDedicatedHostAttributeReq)
		})
		addDebug(modifyDedicatedHostAttributeReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyDedicatedHostAttributeReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("action_on_maintenance")
		d.SetPartial("auto_placement")
		d.SetPartial("dedicated_host_name")
		d.SetPartial("description")
		d.SetPartial("network_attributes")
	}
	update = false
	modifyDedicatedHostsChargeTypeReq := ecs.CreateModifyDedicatedHostsChargeTypeRequest()
	modifyDedicatedHostsChargeTypeReq.DedicatedHostIds = d.Id()
	modifyDedicatedHostsChargeTypeReq.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("auto_renew") {
		modifyDedicatedHostsChargeTypeReq.AutoPay = requests.NewBoolean(d.Get("auto_renew").(bool))
	}
	if d.HasChange("detail_fee") {
		modifyDedicatedHostsChargeTypeReq.DetailFee = requests.NewBoolean(d.Get("detail_fee").(bool))
	}
	if d.HasChange("dry_run") {
		modifyDedicatedHostsChargeTypeReq.DryRun = requests.NewBoolean(d.Get("dry_run").(bool))
	}
	if !d.IsNewResource() && d.HasChange("expired_time") {
		if v, err := strconv.Atoi(d.Get("expired_time").(string)); err == nil {
			modifyDedicatedHostsChargeTypeReq.Period = requests.NewInteger(v)
		} else {
			return WrapError(err)
		}
	}
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		modifyDedicatedHostsChargeTypeReq.DedicatedHostChargeType = d.Get("payment_type").(string)
	}
	if !d.IsNewResource() && d.HasChange("sale_cycle") {
		modifyDedicatedHostsChargeTypeReq.PeriodUnit = d.Get("sale_cycle").(string)
	}
	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyDedicatedHostsChargeType(modifyDedicatedHostsChargeTypeReq)
		})
		addDebug(modifyDedicatedHostsChargeTypeReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyDedicatedHostsChargeTypeReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("auto_renew")
		d.SetPartial("detail_fee")
		d.SetPartial("dry_run")
		d.SetPartial("expired_time")
		d.SetPartial("payment_type")
		d.SetPartial("sale_cycle")
	}
	d.Partial(false)
	return resourceAlicloudEcsDedicatedHostRead(d, meta)
}
func resourceAlicloudEcsDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateReleaseDedicatedHostRequest()
	request.DedicatedHostId = d.Id()
	request.RegionId = client.RegionId
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ReleaseDedicatedHost(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDedicatedHostId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

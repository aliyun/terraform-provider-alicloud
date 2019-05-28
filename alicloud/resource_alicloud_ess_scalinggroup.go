package alicloud

import (
	"fmt"

	"time"

	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEssScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingGroupCreate,
		Read:   resourceAliyunEssScalingGroupRead,
		Update: resourceAliyunEssScalingGroupUpdate,
		Delete: resourceAliyunEssScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 1000),
			},
			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 1000),
			},
			"scaling_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_cooldown": {
				Type:         schema.TypeInt,
				Default:      300,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 86400),
			},
			"vswitch_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'vswitch_id' has been deprecated from provider version 1.7.1, and new field 'vswitch_ids' can replace it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"removal_policies": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MaxItems: 2,
				MinItems: 1,
			},
			"db_instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
				MinItems: 1,
			},
			"loadbalancer_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
				MinItems: 0,
			},
			"multi_az_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      Priority,
				ValidateFunc: validateAllowedStringValue([]string{string(Priority), string(Balance)}),
				ForceNew:     true,
			},
		},
	}
}

func resourceAliyunEssScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {

	request, err := buildAlicloudEssScalingGroupArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateScalingGroup(request)
		})
		if err != nil {
			if IsExceptedError(err, EssThrottling) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*ess.CreateScalingGroupResponse)
		d.SetId(response.ScalingGroupId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scalinggroup", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	if err := essService.WaitForEssScalingGroup(d.Id(), Inactive, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunEssScalingGroupUpdate(d, meta)
}

func resourceAliyunEssScalingGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("min_size", object.MinSize)
	d.Set("max_size", object.MaxSize)
	d.Set("scaling_group_name", object.ScalingGroupName)
	d.Set("default_cooldown", object.DefaultCooldown)
	d.Set("multi_az_policy", object.MultiAZPolicy)
	var polices []string
	if len(object.RemovalPolicies.RemovalPolicy) > 0 {
		for _, v := range object.RemovalPolicies.RemovalPolicy {
			polices = append(polices, v)
		}
	}
	d.Set("removal_policies", polices)
	var dbIds []string
	if len(object.DBInstanceIds.DBInstanceId) > 0 {
		for _, v := range object.DBInstanceIds.DBInstanceId {
			dbIds = append(dbIds, v)
		}
	}
	d.Set("db_instance_ids", dbIds)

	var slbIds []string
	if len(object.LoadBalancerIds.LoadBalancerId) > 0 {
		for _, v := range object.LoadBalancerIds.LoadBalancerId {
			slbIds = append(slbIds, v)
		}
	}
	d.Set("loadbalancer_ids", slbIds)

	var vswitchIds []string
	if len(object.VSwitchIds.VSwitchId) > 0 {
		for _, v := range object.VSwitchIds.VSwitchId {
			vswitchIds = append(vswitchIds, v)
		}
	}
	d.Set("vswitch_ids", vswitchIds)

	return nil
}

func resourceAliyunEssScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateModifyScalingGroupRequest()
	request.ScalingGroupId = d.Id()

	if d.HasChange("scaling_group_name") {
		request.ScalingGroupName = d.Get("scaling_group_name").(string)
	}

	if d.HasChange("min_size") {
		request.MinSize = requests.NewInteger(d.Get("min_size").(int))
	}

	if d.HasChange("max_size") {
		request.MaxSize = requests.NewInteger(d.Get("max_size").(int))
	}

	if d.HasChange("default_cooldown") {
		request.DefaultCooldown = requests.NewInteger(d.Get("default_cooldown").(int))
	}

	if d.HasChange("removal_policies") {
		policyies := d.Get("removal_policies").(*schema.Set).List()
		s := reflect.ValueOf(request).Elem()
		for i, p := range policyies {
			s.FieldByName(fmt.Sprintf("RemovalPolicy%d", i+1)).Set(reflect.ValueOf(p.(string)))
		}
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return resourceAliyunEssScalingGroupRead(d, meta)
}

func resourceAliyunEssScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	request := ess.CreateDeleteScalingGroupRequest()
	request.ScalingGroupId = d.Id()
	request.ForceDelete = requests.NewBoolean(true)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingGroup(request)
	})

	if err != nil {
		if IsExceptedErrors(err, []string{InvalidScalingGroupIdNotFound}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	return WrapError(essService.WaitForEssScalingGroup(d.Id(), Deleted, DefaultTimeout))
}

func buildAlicloudEssScalingGroupArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingGroupRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	request := ess.CreateCreateScalingGroupRequest()

	request.MinSize = requests.NewInteger(d.Get("min_size").(int))
	request.MaxSize = requests.NewInteger(d.Get("max_size").(int))
	request.DefaultCooldown = requests.NewInteger(d.Get("default_cooldown").(int))

	if v, ok := d.GetOk("scaling_group_name"); ok && v.(string) != "" {
		request.ScalingGroupName = v.(string)
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		ids := expandStringList(v.(*schema.Set).List())
		request.VSwitchIds = &ids
	}

	if dbs, ok := d.GetOk("db_instance_ids"); ok {
		request.DBInstanceIds = convertListToJsonString(dbs.(*schema.Set).List())
	}

	if lbs, ok := d.GetOk("loadbalancer_ids"); ok {
		for _, lb := range lbs.(*schema.Set).List() {
			if err := slbService.WaitForSlb(lb.(string), Active, DefaultTimeout); err != nil {
				return nil, WrapError(err)
			}
		}
		request.LoadBalancerIds = convertListToJsonString(lbs.(*schema.Set).List())
	}

	if v, ok := d.GetOk("multi_az_policy"); ok && v.(string) != "" {
		request.MultiAZPolicy = v.(string)
	}

	return request, nil
}

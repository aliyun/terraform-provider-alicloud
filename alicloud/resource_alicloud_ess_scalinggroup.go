package alicloud

import (
	"fmt"

	"time"

	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"min_size": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 100),
			},
			"max_size": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 100),
			},
			"scaling_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_cooldown": &schema.Schema{
				Type:         schema.TypeInt,
				Default:      300,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 86400),
			},
			"vswitch_id": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'vswitch_id' has been deprecated from provider version 1.7.1, and new field 'vswitch_ids' can replace it.",
			},
			"vswitch_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"removal_policies": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MaxItems: 2,
				MinItems: 1,
			},
			"db_instance_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
				MinItems: 1,
			},
			"loadbalancer_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
				MinItems: 1,
			},
		},
	}
}

func resourceAliyunEssScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {

	args, err := buildAlicloudEssScalingGroupArgs(d, meta)
	if err != nil {
		return err
	}

	essconn := meta.(*AliyunClient).essconn

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		scaling, err := essconn.CreateScalingGroup(args)
		if err != nil {
			if IsExceptedError(err, EssThrottling) {
				return resource.RetryableError(fmt.Errorf("CreateScalingGroup timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("CreateScalingGroup got an error: %#v.", err))
		}
		d.SetId(scaling.ScalingGroupId)
		return nil
	}); err != nil {
		return err
	}

	return resourceAliyunEssScalingGroupUpdate(d, meta)
}

func resourceAliyunEssScalingGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	scaling, err := client.DescribeScalingGroupById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe ESS scaling group Attribute: %#v", err)
	}

	d.Set("min_size", scaling.MinSize)
	d.Set("max_size", scaling.MaxSize)
	d.Set("scaling_group_name", scaling.ScalingGroupName)
	d.Set("default_cooldown", scaling.DefaultCooldown)
	var polices []string
	if len(scaling.RemovalPolicies.RemovalPolicy) > 0 {
		for _, v := range scaling.RemovalPolicies.RemovalPolicy {
			polices = append(polices, v)
		}
	}
	d.Set("removal_policies", polices)
	var dbIds []string
	if len(scaling.DBInstanceIds.DBInstanceId) > 0 {
		for _, v := range scaling.DBInstanceIds.DBInstanceId {
			dbIds = append(dbIds, v)
		}
	}
	d.Set("db_instance_ids", dbIds)

	var slbIds []string
	if len(scaling.LoadBalancerIds.LoadBalancerId) > 0 {
		for _, v := range scaling.LoadBalancerIds.LoadBalancerId {
			slbIds = append(slbIds, v)
		}
	}
	d.Set("loadbalancer_ids", slbIds)

	var vswitchIds []string
	if len(scaling.VSwitchIds.VSwitchId) > 0 {
		for _, v := range scaling.VSwitchIds.VSwitchId {
			vswitchIds = append(vswitchIds, v)
		}
	}
	d.Set("vswitch_ids", vswitchIds)

	return nil
}

func resourceAliyunEssScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).essconn
	args := ess.CreateModifyScalingGroupRequest()
	args.ScalingGroupId = d.Id()

	d.Partial(true)

	if d.HasChange("scaling_group_name") {
		args.ScalingGroupName = d.Get("scaling_group_name").(string)
		d.SetPartial("scaling_group_name")
	}

	if d.HasChange("min_size") {
		args.MinSize = requests.NewInteger(d.Get("min_size").(int))
		d.SetPartial("min_size")
	}

	if d.HasChange("max_size") {
		args.MaxSize = requests.NewInteger(d.Get("max_size").(int))
		d.SetPartial("max_size")
	}

	if d.HasChange("default_cooldown") {
		args.DefaultCooldown = requests.NewInteger(d.Get("default_cooldown").(int))
		d.SetPartial("default_cooldown")
	}

	if d.HasChange("removal_policies") {
		policyies := d.Get("removal_policies").(*schema.Set).List()
		s := reflect.ValueOf(args).Elem()
		for i, p := range policyies {
			s.FieldByName(fmt.Sprintf("RemovalPolicy%d", i+1)).Set(reflect.ValueOf(p.(string)))
		}
		d.SetPartial("removal_policies")
	}

	if _, err := conn.ModifyScalingGroup(args); err != nil {
		return err
	}

	d.Partial(false)

	return resourceAliyunEssScalingGroupRead(d, meta)
}

func resourceAliyunEssScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {

	return meta.(*AliyunClient).DeleteScalingGroupById(d.Id())
}

func buildAlicloudEssScalingGroupArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingGroupRequest, error) {
	client := meta.(*AliyunClient)
	args := ess.CreateCreateScalingGroupRequest()

	args.MinSize = requests.NewInteger(d.Get("min_size").(int))
	args.MaxSize = requests.NewInteger(d.Get("max_size").(int))
	args.DefaultCooldown = requests.NewInteger(d.Get("default_cooldown").(int))

	if v := d.Get("scaling_group_name").(string); v != "" {
		args.ScalingGroupName = v
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		ids := expandStringList(v.(*schema.Set).List())
		args.VSwitchIds = &ids
	}

	if dbs, ok := d.GetOk("db_instance_ids"); ok {
		args.DBInstanceIds = convertListToJsonString(dbs.(*schema.Set).List())
	}

	if lbs, ok := d.GetOk("loadbalancer_ids"); ok {
		for _, lb := range lbs.(*schema.Set).List() {
			if err := client.slbconn.WaitForLoadBalancerAsyn(lb.(string), slb.ActiveStatus, DefaultTimeout); err != nil {
				return nil, fmt.Errorf("WaitForLoadbalancer %s %s got error: %#v", lb.(string), slb.ActiveStatus, err)
			}
		}
		args.LoadBalancerIds = convertListToJsonString(lbs.(*schema.Set).List())
	}

	return args, nil
}

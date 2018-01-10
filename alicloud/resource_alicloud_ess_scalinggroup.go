package alicloud

import (
	"fmt"

	"github.com/denverdino/aliyungo/ess"
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"removal_policies": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MaxItems: 2,
			},
			"db_instance_ids": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MaxItems: 3,
			},
			"loadbalancer_ids": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
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

	scaling, err := essconn.CreateScalingGroup(args)
	if err != nil {
		return err
	}

	d.SetId(scaling.ScalingGroupId)

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
	d.Set("removal_policies", scaling.RemovalPolicies.RemovalPolicy)
	d.Set("db_instance_ids", scaling.DBInstanceIds)
	d.Set("loadbalancer_ids", scaling.LoadBalancerId)

	return nil
}

func resourceAliyunEssScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).essconn
	args := &ess.ModifyScalingGroupArgs{
		ScalingGroupId: d.Id(),
	}
	d.Partial(true)

	if d.HasChange("scaling_group_name") {
		args.ScalingGroupName = d.Get("scaling_group_name").(string)
		d.SetPartial("scaling_group_name")
	}

	if d.HasChange("min_size") {
		minsize := d.Get("min_size").(int)
		args.MinSize = &minsize
		d.SetPartial("min_size")
	}

	if d.HasChange("max_size") {
		maxsize := d.Get("max_size").(int)
		args.MaxSize = &maxsize
		d.SetPartial("max_size")
	}

	if d.HasChange("default_cooldown") {
		cooldown := d.Get("default_cooldown").(int)
		args.DefaultCooldown = &cooldown
		d.SetPartial("default_cooldown")
	}

	if d.HasChange("removal_policies") {
		policyStrings := d.Get("removal_policies").([]interface{})
		args.RemovalPolicy = expandStringList(policyStrings)
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

func buildAlicloudEssScalingGroupArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingGroupArgs, error) {
	client := meta.(*AliyunClient)
	args := &ess.CreateScalingGroupArgs{
		RegionId: getRegion(d, meta),
	}

	minsize := d.Get("min_size").(int)
	maxsize := d.Get("max_size").(int)
	cooldown := d.Get("default_cooldown").(int)
	args.MinSize = &minsize
	args.MaxSize = &maxsize
	args.DefaultCooldown = &cooldown

	if v := d.Get("scaling_group_name").(string); v != "" {
		args.ScalingGroupName = v
	}

	if v := d.Get("vswitch_id").(string); v != "" {
		args.VSwitchId = v

		// get vpcId
		vpcId, err := client.GetVpcIdByVSwitchId(v)

		if err != nil {
			return nil, fmt.Errorf("VswitchId %s is not valid of current region", v)
		}
		// fill vpcId by vswitchId
		args.VpcId = vpcId

	}

	dbs, ok := d.GetOk("db_instance_ids")
	if ok {
		dbsStrings := dbs.([]interface{})
		args.DBInstanceId = expandStringList(dbsStrings)
	}

	lbs, ok := d.GetOk("loadbalancer_ids")
	if ok {
		args.LoadBalancerIds = convertListToJsonString(lbs.([]interface{}))
	}

	return args, nil
}

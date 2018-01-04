package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRouterInterfaceCreate,
		Read:   resourceAlicloudRouterInterfaceRead,
		Update: resourceAlicloudRouterInterfaceUpdate,
		Delete: resourceAlicloudRouterInterfaceDelete,

		Schema: map[string]*schema.Schema{
			"opposite_region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"router_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ecs.VRouter), string(ecs.VBR)}),
				ForceNew: true,
			},
			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ecs.InitiatingSide), string(ecs.AcceptingSide)}),
				ForceNew: true,
			},
			"specification": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(GetAllRouterInterfaceSpec()),
			},
			"access_point_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"opposite_access_point_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"opposite_router_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ecs.VRouter), string(ecs.VBR)}),
				Default:  ecs.VRouter,
				ForceNew: true,
			},
			"opposite_router_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"opposite_interface_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"opposite_interface_owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRouterInterfaceDescription,
			},
			"health_check_source_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_target_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	args, err := buildAlicloudRouterInterfaceCreateArgs(d, meta)
	if err != nil {
		return err
	}

	response, err := conn.CreateRouterInterface(args)
	if err != nil {
		return fmt.Errorf("CreateRouterInterface got an error: %#v", err)
	}

	d.SetId(response.RouterInterfaceId)

	if err := conn.WaitForRouterInterfaceAsyn(getRegion(d, meta), d.Id(), ecs.Idle, 300); err != nil {
		return fmt.Errorf("WaitForRouterInterface %s got error: %#v", ecs.Idle, err)
	}

	return resourceAlicloudRouterInterfaceUpdate(d, meta)
}

func resourceAlicloudRouterInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	d.Partial(true)

	args, attributeUpdate, err := buildAlicloudRouterInterfaceModifyAttrArgs(d, meta)
	if err != nil {
		return err
	}

	if attributeUpdate {
		if _, err := conn.ModifyRouterInterfaceAttribute(args); err != nil {
			return fmt.Errorf("ModifyRouterInterfaceAttribute got an error: %#v", err)
		}
	}

	if d.HasChange("specification") && !d.IsNewResource() {
		d.SetPartial("specification")
		if _, err := conn.ModifyRouterInterfaceSpec(&ecs.ModifyRouterInterfaceSpecArgs{
			RouterInterfaceId: d.Id(),
			RegionId:          getRegion(d, meta),
			Spec:              ecs.Spec(d.Get("specification").(string)),
		}); err != nil {
			return fmt.Errorf("ModifyRouterInterfaceSpec got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRouterInterfaceRead(d, meta)
}

func resourceAlicloudRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	filter := ecs.Filter{Key: "RouterInterfaceId", Value: []string{d.Id()}}
	args := &ecs.DescribeRouterInterfacesArgs{
		RegionId: getRegion(d, meta),
		Filter:   []ecs.Filter{filter},
	}
	resp, err := conn.DescribeRouterInterfaces(args)
	if err != nil {
		return fmt.Errorf("DescribeRouterInterfaces got an error: %#v", err)
	}

	routerInterface := resp.RouterInterfaceSet.RouterInterfaceType
	if len(routerInterface) == 0 {
		return fmt.Errorf("No router interface found.")
	}

	for _, ri := range routerInterface {
		if ri.RouterInterfaceId == d.Id() {
			d.Set("role", ri.Role)
			d.Set("specification", ri.Spec)
			d.Set("name", ri.Name)
			d.Set("router_id", ri.RouterId)
			d.Set("router_type", ri.RouterType)
			d.Set("description", ri.Description)
			d.Set("access_point_id", ri.AccessPointId)
			d.Set("opposite_access_point_id", ri.OppositeAccessPointId)
			d.Set("opposite_router_type", ri.OppositeRouterType)
			d.Set("opposite_router_id", ri.OppositeRouterId)
			d.Set("opposite_interface_id", ri.OppositeInterfaceId)
			d.Set("opposite_interface_owner_id", ri.OppositeInterfaceOwnerId)
			d.Set("health_check_source_ip", ri.HealthCheckSourceIp)
			d.Set("health_check_target_ip", ri.HealthCheckTargetIp)
			return nil
		}
	}

	d.SetId("")
	return fmt.Errorf("No router interface found.")
}

func resourceAlicloudRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.OperateRouterInterfaceArgs{
		RegionId:          getRegion(d, meta),
		RouterInterfaceId: d.Id(),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DeleteRouterInterface(args); err != nil {
			if IsExceptedError(err, RouterInterfaceIncorrectStatus) || IsExceptedError(err, DependencyViolationRouterInterfaceReferedByRouteEntry) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(fmt.Errorf("Delete router interface timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting interface %s: %#v", d.Id(), err))
		}
		return nil
	})
}

func buildAlicloudRouterInterfaceCreateArgs(d *schema.ResourceData, meta interface{}) (*ecs.CreateRouterInterfaceArgs, error) {
	client := meta.(*AliyunClient)

	oppositeRegion := common.Region(d.Get("opposite_region").(string))
	if err := client.JudgeRegionValidation("opposite_region", oppositeRegion); err != nil {
		return nil, err
	}

	args := &ecs.CreateRouterInterfaceArgs{
		RegionId:           getRegion(d, meta),
		RouterType:         ecs.RouterType(d.Get("router_type").(string)),
		RouterId:           d.Get("router_id").(string),
		Role:               ecs.Role(d.Get("role").(string)),
		Spec:               ecs.Spec(d.Get("specification").(string)),
		OppositeRegionId:   oppositeRegion,
		OppositeRouterType: ecs.RouterType(d.Get("opposite_router_type").(string)),
	}

	if args.RouterType == ecs.VBR {
		if args.Role != ecs.InitiatingSide {
			return nil, fmt.Errorf("'role': valid value is only 'InitiatingSide' when 'router_type' is 'VBR'.")
		}

		if args.OppositeRouterType != ecs.VRouter {
			return nil, fmt.Errorf("'opposite_router_type': valid value is only 'VRouter' when 'router_type' is 'VBR'.")
		}

		v, ok := d.GetOk("access_point_id")
		if !ok {
			return nil, fmt.Errorf("'access_point_id': required field is not set when 'router_type' is 'VBR'.")
		}
		args.AccessPointId = v.(string)
	} else if args.OppositeRouterType == ecs.VBR {
		if args.Role != ecs.AcceptingSide {
			return nil, fmt.Errorf("'role': valid value is only 'AcceptingSide' when 'opposite_router_type' is 'VBR'.")
		}

		v, ok := d.GetOk("opposite_access_point_id")
		if !ok {
			return nil, fmt.Errorf("'opposite_access_point_id':required field is not set when 'opposite_router_type' is 'VBR'.")
		}
		args.OppositeAccessPointId = v.(string)
	}

	if args.Role == ecs.AcceptingSide {
		if args.Spec == "" {
			args.Spec = Negative
		} else if args.Spec != Negative {
			return nil, fmt.Errorf("'specification': valid value is only '%s' when 'role' is 'AcceptingSide'.", Negative)
		}
	} else if oppositeRegion == getRegion(d, meta) {
		if args.RouterType == ecs.VRouter {
			if args.Spec != ecs.Large2 {
				return nil, fmt.Errorf("'specification': valid value is only '%s' when 'role' is 'InitiatingSide' and 'region' is equal to 'opposite_region' and 'router_type' is 'VRouter'.", ecs.Large2)
			}
		} else {
			if args.Spec != ecs.Middle1 && args.Spec != ecs.Middle2 && args.Spec != ecs.Middle5 && args.Spec != ecs.Large1 {
				return nil, fmt.Errorf("'specification': valid values are '%s', '%s', '%s' and '%s' when 'role' is 'InitiatingSide' and 'region' is equal to 'opposite_region' and 'router_type' is 'VBR'.", ecs.Large1, ecs.Middle1, ecs.Middle2, ecs.Middle5)
			}
		}
	} else if args.Spec == ecs.Large2 {
		return nil, fmt.Errorf("The 'specification' can not be '%s' when 'role' is 'InitiatingSide' and 'region' is not equal to 'opposite_region'.", ecs.Large2)
	}

	return args, nil
}

func buildAlicloudRouterInterfaceModifyAttrArgs(d *schema.ResourceData, meta interface{}) (*ecs.ModifyRouterInterfaceAttributeArgs, bool, error) {

	sourceIp, sourceOk := d.GetOk("health_check_source_ip")
	targetIp, targetOk := d.GetOk("health_check_target_ip")
	if sourceOk && !targetOk || !sourceOk && targetOk {
		return nil, false, fmt.Errorf("The 'health_check_source_ip' and 'health_check_target_ip' should be specified or not at one time.")
	}

	args := &ecs.ModifyRouterInterfaceAttributeArgs{
		RegionId:          getRegion(d, meta),
		RouterInterfaceId: d.Id(),
	}
	attributeUpdate := false

	if d.HasChange("health_check_source_ip") {
		d.SetPartial("health_check_source_ip")
		args.HealthCheckSourceIp = sourceIp.(string)
		args.HealthCheckTargetIp = targetIp.(string)
		attributeUpdate = true
	}

	if d.HasChange("health_check_target_ip") {
		d.SetPartial("health_check_target_ip")
		args.HealthCheckTargetIp = targetIp.(string)
		args.HealthCheckSourceIp = sourceIp.(string)
		attributeUpdate = true
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		args.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		args.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	if d.HasChange("opposite_router_id") {
		d.SetPartial("opposite_router_id")
		args.OppositeRouterId = d.Get("opposite_router_id").(string)
		attributeUpdate = true
	}

	if d.HasChange("opposite_interface_id") {
		d.SetPartial("opposite_interface_id")
		args.OppositeInterfaceId = d.Get("opposite_interface_id").(string)
		attributeUpdate = true
	}

	if d.HasChange("opposite_interface_owner_id") {
		d.SetPartial("opposite_interface_owner_id")
		args.OppositeInterfaceOwnerId = d.Get("opposite_interface_owner_id").(string)
		attributeUpdate = true
	}

	return args, attributeUpdate, nil
}

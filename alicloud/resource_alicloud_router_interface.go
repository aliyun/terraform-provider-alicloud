package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
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
					string(VRouter), string(VBR)}),
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
					string(InitiatingSide), string(AcceptingSide)}),
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
					string(VRouter), string(VBR)}),
				Default:  VRouter,
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
	client := meta.(*AliyunClient)
	args, err := buildAlicloudRouterInterfaceCreateArgs(d, meta)
	if err != nil {
		return err
	}

	response, err := client.vpcconn.CreateRouterInterface(args)
	if err != nil {
		return fmt.Errorf("CreateRouterInterface got an error: %#v", err)
	}

	d.SetId(response.RouterInterfaceId)

	if err := client.WaitForRouterInterface(d.Id(), Idle, 300); err != nil {
		return fmt.Errorf("WaitForRouterInterface %s got error: %#v", Idle, err)
	}

	return resourceAlicloudRouterInterfaceUpdate(d, meta)
}

func resourceAlicloudRouterInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

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
		request := vpc.CreateModifyRouterInterfaceSpecRequest()
		request.RegionId = string(getRegion(d, meta))
		request.RouterInterfaceId = d.Id()
		request.Spec = d.Get("specification").(string)
		if _, err := conn.ModifyRouterInterfaceSpec(request); err != nil {
			return fmt.Errorf("ModifyRouterInterfaceSpec got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRouterInterfaceRead(d, meta)
}

func resourceAlicloudRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {

	ri, err := meta.(*AliyunClient).DescribeRouterInterface(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

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

func resourceAlicloudRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateDeleteRouterInterfaceRequest()
	args.RegionId = string(getRegion(d, meta))
	args.RouterInterfaceId = d.Id()

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

func buildAlicloudRouterInterfaceCreateArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateRouterInterfaceRequest, error) {
	client := meta.(*AliyunClient)

	oppositeRegion := d.Get("opposite_region").(string)
	if err := client.JudgeRegionValidation("opposite_region", oppositeRegion); err != nil {
		return nil, err
	}

	request := vpc.CreateCreateRouterInterfaceRequest()
	request.RegionId = string(getRegion(d, meta))
	request.RouterType = d.Get("router_type").(string)
	request.RouterId = d.Get("router_id").(string)
	request.Role = d.Get("role").(string)
	request.Spec = d.Get("specification").(string)
	request.OppositeRegionId = oppositeRegion
	request.OppositeRouterType = d.Get("opposite_router_type").(string)

	if request.RouterType == string(VBR) {
		if request.Role != string(InitiatingSide) {
			return nil, fmt.Errorf("'role': valid value is only 'InitiatingSide' when 'router_type' is 'VBR'.")
		}

		if request.OppositeRouterType != string(VRouter) {
			return nil, fmt.Errorf("'opposite_router_type': valid value is only 'VRouter' when 'router_type' is 'VBR'.")
		}

		v, ok := d.GetOk("access_point_id")
		if !ok {
			return nil, fmt.Errorf("'access_point_id': required field is not set when 'router_type' is 'VBR'.")
		}
		request.AccessPointId = v.(string)
	} else if request.OppositeRouterType == string(VBR) {
		if request.Role != string(AcceptingSide) {
			return nil, fmt.Errorf("'role': valid value is only 'AcceptingSide' when 'opposite_router_type' is 'VBR'.")
		}

		v, ok := d.GetOk("opposite_access_point_id")
		if !ok {
			return nil, fmt.Errorf("'opposite_access_point_id':required field is not set when 'opposite_router_type' is 'VBR'.")
		}
		request.OppositeAccessPointId = v.(string)
	}

	if request.Role == string(AcceptingSide) {
		if request.Spec == "" {
			request.Spec = string(Negative)
		} else if request.Spec != string(Negative) {
			return nil, fmt.Errorf("'specification': valid value is only '%s' when 'role' is 'AcceptingSide'.", Negative)
		}
	} else if oppositeRegion == getRegionId(d, meta) {
		if request.RouterType == string(VRouter) {
			if request.Spec != string(Large2) {
				return nil, fmt.Errorf("'specification': valid value is only '%s' when 'role' is 'InitiatingSide' and 'region' is equal to 'opposite_region' and 'router_type' is 'VRouter'.", Large2)
			}
		} else {
			if request.Spec != string(Middle1) && request.Spec != string(Middle2) && request.Spec != string(Middle5) && request.Spec != string(Large1) {
				return nil, fmt.Errorf("'specification': valid values are '%s', '%s', '%s' and '%s' when 'role' is 'InitiatingSide' and 'region' is equal to 'opposite_region' and 'router_type' is 'VBR'.", Large1, Middle1, Middle2, Middle5)
			}
		}
	} else if request.Spec == string(Large2) {
		return nil, fmt.Errorf("The 'specification' can not be '%s' when 'role' is 'InitiatingSide' and 'region' is not equal to 'opposite_region'.", Large2)
	}

	return request, nil
}

func buildAlicloudRouterInterfaceModifyAttrArgs(d *schema.ResourceData, meta interface{}) (*vpc.ModifyRouterInterfaceAttributeRequest, bool, error) {

	sourceIp, sourceOk := d.GetOk("health_check_source_ip")
	targetIp, targetOk := d.GetOk("health_check_target_ip")
	if sourceOk && !targetOk || !sourceOk && targetOk {
		return nil, false, fmt.Errorf("The 'health_check_source_ip' and 'health_check_target_ip' should be specified or not at one time.")
	}

	args := vpc.CreateModifyRouterInterfaceAttributeRequest()
	args.RegionId = string(getRegion(d, meta))
	args.RouterInterfaceId = d.Id()

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
		args.OppositeInterfaceOwnerId = requests.Integer(d.Get("opposite_interface_owner_id").(string))
		attributeUpdate = true
	}

	return args, attributeUpdate, nil
}

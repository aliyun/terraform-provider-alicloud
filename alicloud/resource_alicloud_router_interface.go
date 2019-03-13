package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRouterInterfaceCreate,
		Read:   resourceAlicloudRouterInterfaceRead,
		Update: resourceAlicloudRouterInterfaceUpdate,
		Delete: resourceAlicloudRouterInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"opposite_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"router_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(VRouter), string(VBR)}),
				ForceNew:         true,
				DiffSuppressFunc: routerInterfaceAcceptsideDiffSuppressFunc,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(InitiatingSide), string(AcceptingSide)}),
				ForceNew: true,
			},
			"specification": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validateAllowedStringValue(GetAllRouterInterfaceSpec()),
				DiffSuppressFunc: routerInterfaceAcceptsideDiffSuppressFunc,
			},
			"access_point_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_access_point_id' has been deprecated from version 1.11.0.",
			},
			"opposite_access_point_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Attribute 'opposite_access_point_id' has been deprecated from version 1.11.0.",
			},
			"opposite_router_type": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_router_type' has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_router_type' instead.",
			},
			"opposite_router_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_router_id' has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_router_id' instead.",
			},
			"opposite_interface_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_interface_id' has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_interface_id' instead.",
			},
			"opposite_interface_owner_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_interface_owner_id' has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_interface_owner_id' instead.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRouterInterfaceDescription,
			},
			"health_check_source_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: routerInterfaceVBRTypeDiffSuppressFunc,
			},
			"health_check_target_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: routerInterfaceVBRTypeDiffSuppressFunc,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validateInstanceChargeType,
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Default:          1,
				DiffSuppressFunc: ecsPostPaidDiffSuppressFunc,
				ValidateFunc:     validateRouterInterfaceChargeTypePeriod,
			},
		},
	}
}

func resourceAlicloudRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	args, err := buildAlicloudRouterInterfaceCreateArgs(d, meta)
	if err != nil {
		return err
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateRouterInterface(args)
	})
	if err != nil {
		return fmt.Errorf("CreateRouterInterface got an error: %#v", err)
	}
	response, _ := raw.(*vpc.CreateRouterInterfaceResponse)
	d.SetId(response.RouterInterfaceId)

	if err := vpcService.WaitForRouterInterface(client.RegionId, d.Id(), Idle, 300); err != nil {
		return fmt.Errorf("WaitForRouterInterface %s got error: %#v", Idle, err)
	}

	return resourceAlicloudRouterInterfaceUpdate(d, meta)
}

func resourceAlicloudRouterInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	args, attributeUpdate, err := buildAlicloudRouterInterfaceModifyAttrArgs(d, meta)
	if err != nil {
		return err
	}

	if attributeUpdate {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyRouterInterfaceAttribute(args)
		})
		if err != nil {
			return fmt.Errorf("ModifyRouterInterfaceAttribute got an error: %#v", err)
		}
	}

	if d.HasChange("specification") && !d.IsNewResource() {
		d.SetPartial("specification")
		request := vpc.CreateModifyRouterInterfaceSpecRequest()
		request.RegionId = string(client.Region)
		request.RouterInterfaceId = d.Id()
		request.Spec = d.Get("specification").(string)
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyRouterInterfaceSpec(request)
		})
		if err != nil {
			return fmt.Errorf("ModifyRouterInterfaceSpec got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRouterInterfaceRead(d, meta)
}

func resourceAlicloudRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	ri, err := vpcService.DescribeRouterInterface(client.RegionId, d.Id())
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
	d.Set("opposite_region", ri.OppositeRegionId)
	d.Set("opposite_router_type", ri.OppositeRouterType)
	d.Set("opposite_router_id", ri.OppositeRouterId)
	d.Set("opposite_interface_id", ri.OppositeInterfaceId)
	d.Set("opposite_interface_owner_id", ri.OppositeInterfaceOwnerId)
	d.Set("health_check_source_ip", ri.HealthCheckSourceIp)
	d.Set("health_check_target_ip", ri.HealthCheckTargetIp)
	if ri.ChargeType == "Prepaid" {
		d.Set("instance_charge_type", PrePaid)
	} else {
		d.Set("instance_charge_type", PostPaid)
	}
	return nil

}

func resourceAlicloudRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	if ri, err := vpcService.DescribeRouterInterface(client.RegionId, d.Id()); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("When deleting router interface %s, describing it got an error: %#v.", d.Id(), err)
	} else if ri.Status == string(Active) {
		if err := vpcService.DeactivateRouterInterface(d.Id()); err != nil {
			return fmt.Errorf("When deleting router interface %s, deactiving it got an error: %#v.", d.Id(), err)
		}
	}

	args := vpc.CreateDeleteRouterInterfaceRequest()
	args.RegionId = string(client.Region)
	args.RouterInterfaceId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouterInterface(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidInstanceIdNotFound}) {
				return nil
			}
			if IsExceptedErrors(err, []string{RouterInterfaceIncorrectStatus, DependencyViolationRouterInterfaceReferedByRouteEntry}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(fmt.Errorf("Delete router interface timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting interface %s: %#v", d.Id(), err))
		}
		if _, err := vpcService.DescribeRouterInterface(client.RegionId, d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("When deleting router interface %s, describing it got an error: %#v.", d.Id(), err))
		}
		return resource.RetryableError(fmt.Errorf("Deleting interface %s timeout.", d.Id()))
	})
}

func buildAlicloudRouterInterfaceCreateArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateRouterInterfaceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	oppositeRegion := d.Get("opposite_region").(string)
	if err := ecsService.JudgeRegionValidation("opposite_region", oppositeRegion); err != nil {
		return nil, err
	}

	request := vpc.CreateCreateRouterInterfaceRequest()
	request.RouterType = d.Get("router_type").(string)
	request.RouterId = d.Get("router_id").(string)
	request.Role = d.Get("role").(string)
	request.Spec = d.Get("specification").(string)
	request.InstanceChargeType = d.Get("instance_charge_type").(string)
	if request.InstanceChargeType == string(PrePaid) {
		period := d.Get("period").(int)
		request.Period = requests.NewInteger(period)
		request.PricingCycle = string(Month)
		if period > 9 {
			request.Period = requests.NewInteger(period / 12)
			request.PricingCycle = string(Year)
		}
		request.AutoPay = requests.NewBoolean(true)
	}
	request.OppositeRegionId = oppositeRegion
	// Accepting side router interface spec only be Negative and router type only be VRouter.
	if request.Role == string(AcceptingSide) {
		request.Spec = string(Negative)
		request.RouterType = string(VRouter)
	} else {
		if request.Spec == "" {
			return request, fmt.Errorf("'specification': required field is not set when role is %s.", InitiatingSide)
		}
	}

	// Get VBR access point
	if request.RouterType == string(VBR) {
		req := vpc.CreateDescribeVirtualBorderRoutersRequest()
		values := []string{request.RouterId}
		filters := []vpc.DescribeVirtualBorderRoutersFilter{{
			Key:   "VbrId",
			Value: &values,
		}}
		req.Filter = &filters
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVirtualBorderRouters(req)
		})
		if err != nil {
			return request, fmt.Errorf("Describing VBR %s got an error: %#v.", request.RouterId, err)
		}
		resp, _ := raw.(*vpc.DescribeVirtualBorderRoutersResponse)
		if resp != nil && resp.TotalCount > 0 {
			request.AccessPointId = resp.VirtualBorderRouterSet.VirtualBorderRouterType[0].AccessPointId
		}
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	return request, nil
}

func buildAlicloudRouterInterfaceModifyAttrArgs(d *schema.ResourceData, meta interface{}) (*vpc.ModifyRouterInterfaceAttributeRequest, bool, error) {

	sourceIp, sourceOk := d.GetOk("health_check_source_ip")
	targetIp, targetOk := d.GetOk("health_check_target_ip")
	if sourceOk && !targetOk || !sourceOk && targetOk {
		return nil, false, fmt.Errorf("The 'health_check_source_ip' and 'health_check_target_ip' should be specified or not at one time.")
	}

	args := vpc.CreateModifyRouterInterfaceAttributeRequest()
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

	return args, attributeUpdate, nil
}

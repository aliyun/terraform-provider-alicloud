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

func resourceAlicloudRouterInterfaceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRouterInterfaceConnectionCreate,
		Read:   resourceAlicloudRouterInterfaceConnectionRead,
		Delete: resourceAlicloudRouterInterfaceConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"interface_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"opposite_interface_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"opposite_router_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(VRouter), string(VBR)}),
				Default:  VRouter,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("opposite_interface_owner_id")
				},
			},
			"opposite_router_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("opposite_interface_owner_id")
				},
			},
			"opposite_interface_owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRouterInterfaceConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	oppsiteId := d.Get("opposite_interface_id").(string)
	interfaceId := d.Get("interface_id").(string)
	ri, err := vpcService.DescribeRouterInterface(client.RegionId, interfaceId)
	if err != nil {
		return err
	}

	// At present, the interface with "active/inactive" status can not be modify opposite connection information
	// and it is RouterInterface product limitation.
	if ri.OppositeInterfaceId == oppsiteId {
		if ri.Status == string(Active) {
			return fmt.Errorf("The specified router interface connection has existed, and please import it using id %s.", interfaceId)
		}
		if ri.Status == string(Inactive) {
			if err := vpcService.ActivateRouterInterface(interfaceId); err != nil {
				return err
			}
			if err := vpcService.WaitForRouterInterface(client.RegionId, interfaceId, Active, DefaultTimeout); err != nil {
				return fmt.Errorf("When activing router interface %s got an error: %#v.", interfaceId, err)
			}
			d.SetId(interfaceId)
			return resourceAlicloudRouterInterfaceConnectionRead(d, meta)
		}
	}

	req := vpc.CreateModifyRouterInterfaceAttributeRequest()
	req.RouterInterfaceId = interfaceId
	req.OppositeInterfaceId = oppsiteId

	if owner_id, ok := d.GetOk("opposite_interface_owner_id"); ok && owner_id.(string) != "" {
		req.OppositeInterfaceOwnerId = requests.Integer(owner_id.(string))
		if v, o := d.GetOk("opposite_router_type"); !o || v.(string) == "" {
			return fmt.Errorf("'opposite_router_type' is required when 'opposite_interface_owner_id' is set.")
		} else {
			req.OppositeRouterType = v.(string)
		}

		if v, o := d.GetOk("opposite_router_id"); !o || v.(string) == "" {
			return fmt.Errorf("'opposite_router_id' is required when 'opposite_interface_owner_id' is set.")
		} else {
			req.OppositeRouterId = v.(string)
		}
	} else {
		owner := ri.OppositeInterfaceOwnerId
		if owner == "" {
			owner, err = client.AccountId()
			if err != nil {
				return err
			}
		}
		if owner == "" {
			return fmt.Errorf("Opposite router interface owner id is empty. Please use field 'opposite_interface_owner_id' or globle field 'account_id' to set.")
		}
		oppositeRi, err := vpcService.DescribeRouterInterface(ri.OppositeRegionId, oppsiteId)
		if err != nil {
			return err
		}
		req.OppositeRouterId = oppositeRi.RouterId
		req.OppositeRouterType = oppositeRi.RouterType
		req.OppositeInterfaceOwnerId = requests.Integer(owner)
	}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {

		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyRouterInterfaceAttribute(req)
		})
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Modifying RouterInterface %s Connection got an error: %#v.", interfaceId, err))
		}

		ri, err := vpcService.DescribeRouterInterface(client.RegionId, interfaceId)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("When modifying RouterInterface %s Connection, describing it got an error: %#v.", interfaceId, err))

		}
		if ri.OppositeInterfaceId == "" || ri.OppositeRouterType == "" ||
			ri.OppositeRouterId == "" || ri.OppositeInterfaceOwnerId == "" {
			return resource.RetryableError(fmt.Errorf("Modifying RouterInterface %s Connection timeout with opposite interface id is %s, "+
				"opposite router id is %s, opposite router type is %s and opposite interface owner id is %s.", interfaceId, ri.OppositeInterfaceId,
				ri.OppositeRouterId, ri.OppositeRouterType, ri.OppositeInterfaceOwnerId))
		}
		return nil
	}); err != nil {
		return err
	}
	if ri.Role == string(InitiatingSide) {
		connReq := vpc.CreateConnectRouterInterfaceRequest()
		connReq.RouterInterfaceId = interfaceId

		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {

			_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.ConnectRouterInterface(connReq)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{IncorrectOppositeInterfaceInfoNotSet}) {
					return resource.RetryableError(fmt.Errorf("Connecting router interface %s timeout.", interfaceId))
				}
				return resource.NonRetryableError(fmt.Errorf("Connecting router interface %s got an error: %#v.", interfaceId, err))
			}

			return nil
		}); err != nil {
			return err
		}

		if err := vpcService.WaitForRouterInterface(client.RegionId, interfaceId, Active, DefaultTimeout); err != nil {
			return fmt.Errorf("Connecting router interface %s got an error: %#v.", interfaceId, err)
		}
	}
	d.SetId(interfaceId)

	return resourceAlicloudRouterInterfaceConnectionRead(d, meta)
}

func resourceAlicloudRouterInterfaceConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	ri, err := vpcService.DescribeRouterInterface(client.RegionId, d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}
	if ri.Status == string(Inactive) {
		if err := vpcService.ActivateRouterInterface(d.Id()); err != nil {
			return err
		}
		if err := vpcService.WaitForRouterInterface(client.RegionId, d.Id(), Active, DefaultTimeout); err != nil {
			return fmt.Errorf("When activing router interface %s got an error: %#v.", d.Id(), err)
		}
	}

	d.Set("interface_id", ri.RouterInterfaceId)
	d.Set("opposite_interface_id", ri.OppositeInterfaceId)
	d.Set("opposite_router_type", ri.OppositeRouterType)
	d.Set("opposite_router_id", ri.OppositeRouterId)
	d.Set("opposite_interface_owner_id", ri.OppositeInterfaceOwnerId)

	return nil

}

func resourceAlicloudRouterInterfaceConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	ri, err := vpcService.DescribeRouterInterface(client.RegionId, d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	if ri.Status == string(Idle) {
		d.SetId("")
		return nil
	}

	// At present, the interface with "active/inactive" status can not be modify opposite connection information
	// and it is RouterInterface product limitation. So, the connection delete action is only modifying it to inactive.
	if ri.Status == string(Active) {
		if err := vpcService.DeactivateRouterInterface(d.Id()); err != nil {
			return err
		}
	}

	if err := vpcService.WaitForRouterInterface(client.RegionId, d.Id(), Inactive, DefaultTimeoutMedium); err != nil {
		return fmt.Errorf("Deleting routerinterface %s connection got an error: %#v.", d.Id(), err)
	}

	return nil
}

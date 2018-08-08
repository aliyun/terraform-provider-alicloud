package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
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
	client := meta.(*AliyunClient)

	oppsiteId := d.Get("opposite_interface_id").(string)
	interfaceId := d.Get("interface_id").(string)
	ri, err := client.DescribeRouterInterface(interfaceId)
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
			if err := client.ActivateRouterInterface(interfaceId); err != nil {
				return err
			}
			if err := client.WaitForRouterInterface(interfaceId, Active, DefaultTimeout); err != nil {
				return fmt.Errorf("When activing router interface %s got an error: %#v.", interfaceId, err)
			}
			d.SetId(interfaceId)
			return resourceAlicloudRouterInterfaceConnectionRead(d, meta)
		}
	}

	req := vpc.CreateModifyRouterInterfaceAttributeRequest()
	req.RouterInterfaceId = interfaceId
	req.OppositeInterfaceId = oppsiteId

	if owner_id, ok := d.GetOk("opposite_interface_owner_id"); ok && owner_id.(string) != ri.OppositeInterfaceOwnerId {
		req.OppositeInterfaceOwnerId = requests.Integer(owner_id.(string))
		if v, o := d.GetOk("opposite_router_type"); !o || v.(string) == "" {
			return fmt.Errorf("'opposite_router_type' is required when 'opposite_interface_owner_id' is not the current account.")
		} else {
			req.OppositeRouterType = v.(string)
		}

		if v, o := d.GetOk("opposite_router_id"); !o || v.(string) == "" {
			return fmt.Errorf("'opposite_router_id' is required when 'opposite_interface_owner_id' is not the current account.")
		} else {
			req.OppositeRouterId = v.(string)
		}
	} else {
		oppositeRi, err := client.DescribeRouterInterfaceInSpecifiedRegion(ri.OppositeRegionId, oppsiteId)
		if err != nil {
			return err
		}
		req.OppositeRouterId = oppositeRi["RouterId"].(string)
		req.OppositeRouterType = oppositeRi["RouterType"].(string)
		owner := ri.OppositeInterfaceOwnerId
		if owner == "" {
			owner = client.AccountId
		}
		req.OppositeInterfaceOwnerId = requests.Integer(owner)
	}

	if _, err := client.vpcconn.ModifyRouterInterfaceAttribute(req); err != nil {
		return fmt.Errorf("Modifying RouterInterface %s Connection got an error: %#v.", interfaceId, err)
	}

	if ri.Role == string(InitiatingSide) {
		connReq := vpc.CreateConnectRouterInterfaceRequest()
		connReq.RouterInterfaceId = interfaceId

		if _, err := client.vpcconn.ConnectRouterInterface(connReq); err != nil {
			return fmt.Errorf("Connecting router interface %s got an error: %#v.", interfaceId, err)
		}

		if err := client.WaitForRouterInterface(interfaceId, Active, DefaultTimeout); err != nil {
			return fmt.Errorf("Connecting router interface %s got an error: %#v.", interfaceId, err)
		}
	}
	d.SetId(interfaceId)

	return resourceAlicloudRouterInterfaceConnectionRead(d, meta)
}

func resourceAlicloudRouterInterfaceConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	ri, err := client.DescribeRouterInterface(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}
	if ri.Status == string(Inactive) {
		if err := client.ActivateRouterInterface(d.Id()); err != nil {
			return err
		}
		if err := client.WaitForRouterInterface(d.Id(), Active, DefaultTimeout); err != nil {
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
	client := meta.(*AliyunClient)

	ri, err := client.DescribeRouterInterface(d.Id())
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
		if err := client.DeactivateRouterInterface(d.Id()); err != nil {
			return err
		}
	}

	if err := client.WaitForRouterInterface(d.Id(), Inactive, DefaultTimeoutMedium); err != nil {
		return fmt.Errorf("Deleting routerinterface %s connection got an error: %#v.", d.Id(), err)
	}

	return nil
}

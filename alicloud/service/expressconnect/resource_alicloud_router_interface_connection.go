package expressconnect

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudRouterInterfaceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRouterInterfaceConnectionCreate,
		Read:   resourceAlicloudRouterInterfaceConnectionRead,
		Delete: resourceAlicloudRouterInterfaceConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		DeprecationMessage: "This resource has been deprecated since v1.199.0 and will be removed in the future. Please use 'alicloud_express_connect_router_interface' instead.",

		Schema: map[string]*schema.Schema{
			"interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"opposite_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"opposite_router_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(VRouter), string(VBR)}, false),
				Default:  VRouter,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("opposite_interface_owner_id")
				},
			},
			"opposite_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("opposite_interface_owner_id")
				},
			},
			"opposite_interface_owner_id": {
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
	vpcService := ExpressConnectService{client}

	oppositeId := d.Get("opposite_interface_id").(string)
	interfaceId := d.Get("interface_id").(string)
	object, err := vpcService.DescribeRouterInterface(interfaceId, client.RegionId)
	if err != nil {
		return errs.WrapError(err)
	}

	// At present, the interface with "active/inactive" status can not be modify opposite connection information
	// and it is RouterInterface product limitation.
	if object.OppositeInterfaceId == oppositeId {
		if object.Status == string(Active) {
			return errs.WrapError(errs.Error("The specified router interface connection has existed, and please import it using id %s.", interfaceId))
		}
		if object.Status == string(Inactive) {
			if err = vpcService.ActivateRouterInterface(interfaceId); err != nil {
				return errs.WrapError(err)
			}
			d.SetId(object.RouterInterfaceId)
			if err = vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Active, DefaultTimeout); err != nil {
				return errs.WrapError(err)
			}
			return resourceAlicloudRouterInterfaceConnectionRead(d, meta)
		}
	}

	request := vpc.CreateModifyRouterInterfaceAttributeRequest()
	request.RegionId = client.RegionId
	request.RouterInterfaceId = interfaceId
	request.OppositeInterfaceId = oppositeId

	if owner_id, ok := d.GetOk("opposite_interface_owner_id"); ok && owner_id.(string) != "" {
		request.OppositeInterfaceOwnerId = requests.Integer(owner_id.(string))
		if v, o := d.GetOk("opposite_router_type"); !o || v.(string) == "" {
			return errs.WrapError(errs.Error("'opposite_router_type' is required when 'opposite_interface_owner_id' is set."))
		} else {
			request.OppositeRouterType = v.(string)
		}

		if v, o := d.GetOk("opposite_router_id"); !o || v.(string) == "" {
			return errs.WrapError(errs.Error("'opposite_router_id' is required when 'opposite_interface_owner_id' is set."))
		} else {
			request.OppositeRouterId = v.(string)
		}
	} else {
		owner := object.OppositeInterfaceOwnerId
		if owner == "" {
			owner, err = client.AccountId()
			if err != nil {
				return errs.WrapError(err)
			}
		}
		if owner == "" {
			return errs.WrapError(errs.Error("Opposite router interface owner id is empty. Please use field 'opposite_interface_owner_id' or globle field 'account_id' to set."))
		}
		oppositeRi, err := vpcService.DescribeRouterInterface(oppositeId, object.OppositeRegionId)
		if err != nil {
			return errs.WrapError(err)
		}
		request.OppositeRouterId = oppositeRi.RouterId
		request.OppositeRouterType = oppositeRi.RouterType
		request.OppositeInterfaceOwnerId = requests.Integer(owner)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyRouterInterfaceAttribute(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_router_interface_connection", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(interfaceId)

	if err = vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Idle, DefaultTimeout); err != nil {
		return errs.WrapError(err)
	}

	if object.Role == string(InitiatingSide) {
		connectRequest := vpc.CreateConnectRouterInterfaceRequest()
		connectRequest.RegionId = client.RegionId
		connectRequest.RouterInterfaceId = interfaceId
		if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.ConnectRouterInterface(connectRequest)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"IncorrectOppositeInterfaceInfo.NotSet"}) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(connectRequest.GetActionName(), raw, connectRequest.RpcRequest, connectRequest)
			return nil
		}); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), connectRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		if err := vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Active, DefaultTimeout); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), connectRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudRouterInterfaceConnectionRead(d, meta)
}

func resourceAlicloudRouterInterfaceConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := ExpressConnectService{client}
	object, err := vpcService.DescribeRouterInterfaceConnection(d.Id(), client.RegionId)

	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	if object.Status == string(Inactive) {
		if err := vpcService.ActivateRouterInterface(d.Id()); err != nil {
			return errs.WrapError(err)
		}
		if err := vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Active, DefaultTimeout); err != nil {
			return errs.WrapError(err)
		}
	}

	d.Set("interface_id", object.RouterInterfaceId)
	d.Set("opposite_interface_id", object.OppositeInterfaceId)
	d.Set("opposite_router_type", object.OppositeRouterType)
	d.Set("opposite_router_id", object.OppositeRouterId)
	d.Set("opposite_interface_owner_id", object.OppositeInterfaceOwnerId)

	return nil

}

func resourceAlicloudRouterInterfaceConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := ExpressConnectService{client}
	object, err := vpcService.DescribeRouterInterfaceConnection(d.Id(), client.RegionId)
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	if object.Status == string(Idle) {
		d.SetId("")
		return nil
	}

	// At present, the interface with "active/inactive" status can not be modify opposite connection information
	// and it is RouterInterface product limitation. So, the connection delete action is only modifying it to inactive.
	if object.Status == string(Active) {
		if err := vpcService.DeactivateRouterInterface(d.Id()); err != nil {
			return errs.WrapError(err)
		}
	}

	return errs.WrapError(vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Inactive, DefaultTimeoutMedium))
}

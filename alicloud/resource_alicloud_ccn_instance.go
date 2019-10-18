package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCcnInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCcnInstanceCreate,
		Read:   resourceAlicloudCcnInstanceRead,
		Update: resourceAlicloudCcnInstanceUpdate,
		Delete: resourceAlicloudCcnInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceDescription,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateVpnCIDRNetworkAddress,
			},
			"total_count": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cen_uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCcnInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateCreateCloudConnectNetworkRequest()

	request.Name = d.Get("name").(string)
	request.IsDefault = requests.NewBoolean(d.Get("is_default").(bool))

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("cidr_block"); ok && v.(string) != "" {
		request.CidrBlock = v.(string)
	}

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.CreateCloudConnectNetwork(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ccn_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartag.CreateCloudConnectNetworkResponse)
	d.SetId(response.CcnId)

	return resourceAlicloudCcnInstanceRead(d, meta)
}

func resourceAlicloudCcnInstanceRead(d *schema.ResourceData, meta interface{}) error {
	ccnService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := ccnService.DescribeCcnInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("cidr_block", object.CidrBlock)
	d.Set("is_default", object.IsDefault)

	if d.Get("total_count") == "1" {
		object, err := ccnService.DescribeCcnGrantRule(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}
		d.Set("cen_id", object.CenInstanceId)
		d.Set("cen_uid", object.CenUid)
	}

	return nil
}

func resourceAlicloudCcnInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := smartag.CreateModifyCloudConnectNetworkRequest()
	request.CcnId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}
	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if d.HasChange("cidr_block") {
		request.CidrBlock = d.Get("cidr_block").(string)
		update = true
	}
	if update {
		raw, err := client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
			return ccnClient.ModifyCloudConnectNetwork(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.Get("total_count") == "1" {
		requestGrant := smartag.CreateGrantInstanceToCbnRequest()
		requestGrant.CcnInstanceId = d.Id()
		requestGrant.CenInstanceId = d.Get("cen_id").(string)
		requestGrant.CenUid = d.Get("cen_uid").(string)
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.GrantInstanceToCbn(requestGrant)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), requestGrant.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(requestGrant.GetActionName(), raw, requestGrant.RpcRequest, requestGrant)

	} else if d.Get("total_count") == "0" {
		requestRevoke := smartag.CreateRevokeInstanceFromCbnRequest()
		requestRevoke.CcnInstanceId = d.Id()
		requestRevoke.CenInstanceId = d.Get("cen_id").(string)
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.RevokeInstanceFromCbn(requestRevoke)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), requestRevoke.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(requestRevoke.GetActionName(), raw, requestRevoke.RpcRequest, requestRevoke)
	}
	return resourceAlicloudCcnInstanceRead(d, meta)
}

func resourceAlicloudCcnInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateDeleteCloudConnectNetworkRequest()
	request.CcnId = d.Id()

	raw, err := client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
		return ccnClient.DeleteCloudConnectNetwork(request)
	})

	if err != nil {
		if IsExceptedError(err, "ParameterCcnInstanceId") {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

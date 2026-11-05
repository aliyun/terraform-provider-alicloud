package edas

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudEdasSlbAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasSlbAttachmentCreate,
		Read:   resourceAlicloudEdasSlbAttachmentRead,
		Delete: resourceAlicloudEdasSlbAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slb_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slb_ip": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsIPAddress,
				Required:     true,
				ForceNew:     true,
			},
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"listener_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Optional:     true,
				ForceNew:     true,
			},
			"vserver_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"slb_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEdasSlbAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	slbId := d.Get("slb_id").(string)

	request := edas.CreateBindSlbRequest()
	request.RegionId = client.RegionId
	request.Type = d.Get("type").(string)
	request.AppId = appId
	request.SlbId = slbId
	request.SlbIp = d.Get("slb_ip").(string)
	if v, ok := d.GetOk("listener_port"); ok {
		request.ListenerPort = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("vserver_group_id"); ok {
		request.VServerGroupId = v.(string)
	}

	if err := edasService.SyncResource("slb"); err != nil {
		return err
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.BindSlb(request)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_edas_slb_attachment", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.BindSlbResponse)
	if response.Code != 200 {
		return errs.WrapError(errs.Error("bind slb failed for %s", response.Message))
	}
	d.SetId(appId + ":" + slbId)
	return resourceAlicloudEdasSlbAttachmentRead(d, meta)
}

func resourceAlicloudEdasSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	id := d.Id()
	strs, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}

	regionId := client.RegionId
	slbId := strs[1]
	appId := strs[0]

	rq := edas.CreateGetApplicationRequest()
	rq.RegionId = regionId
	rq.AppId = appId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetApplication(rq)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_edas_slb_attachment", rq.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(rq.GetActionName(), raw, rq.RoaRequest, rq)

	rs := raw.(*edas.GetApplicationResponse)
	if rs.Applcation.SlbId != slbId && rs.Applcation.ExtSlbId != slbId {
		return errs.WrapError(errs.Error("can not find slb: %s", slbId))
	}

	request := edas.CreateListSlbRequest()
	request.RegionId = regionId

	raw, err = edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListSlb(request)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_edas_slb_attachment", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response := raw.(*edas.ListSlbResponse)
	if response.Code != 200 {
		return errs.WrapError(errs.Error("List Slb failed for %s", response.Message))
	}

	for _, slb := range response.SlbList.SlbEntity {
		if slb.SlbId == slbId {
			d.Set("slb_status", slb.SlbStatus)
			d.Set("vswitch_id", slb.VswitchId)
			return nil
		}
	}

	return nil
}

func resourceAlicloudEdasSlbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateUnbindSlbRequest()
	request.RegionId = client.RegionId
	request.AppId = d.Get("app_id").(string)
	request.SlbId = d.Get("slb_id").(string)
	request.Type = d.Get("type").(string)

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.UnbindSlb(request)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_edas_slb_attachment", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response := raw.(*edas.UnbindSlbResponse)
	if response.Code != 200 {
		return errs.WrapError(errs.Error("unbind slb failed, %s", response.Message))
	}

	return nil
}

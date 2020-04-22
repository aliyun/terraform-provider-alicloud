package alicloud

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
				ValidateFunc: validation.SingleIP(),
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
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEdasSlbAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	regionId := client.RegionId
	slbId := d.Get("slb_id").(string)
	slbIp := d.Get("slb_ip").(string)
	slbType := d.Get("type").(string)
	listenerPort := d.Get("listener_port").(int)
	vserverGroupId := d.Get("vserver_group_id").(string)

	request := edas.CreateBindSlbRequest()
	request.RegionId = regionId
	request.Type = slbType
	request.AppId = appId
	request.SlbId = slbId
	request.SlbIp = slbIp
	request.ListenerPort = requests.NewInteger(listenerPort)
	request.VServerGroupId = vserverGroupId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.BindSlb(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_slb_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.BindSlbResponse)
	if response.Code != 200 {
		return Error("bind slb failed for " + response.Message)
	}
	d.SetId(appId + ":" + slbId)
	return resourceAlicloudEdasInstanceApplicationAttachmentRead(d, meta)
}

func resourceAlicloudEdasSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	id := d.Id()
	strs := strings.Split(id, ":")
	if len(strs) != 2 {
		return WrapError(Error("resource id decode failed: " + id))
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_slb_attachment", rq.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(rq.GetActionName(), raw, rq.RoaRequest, rq)

	rs := raw.(*edas.GetApplicationResponse)
	if rs.Applcation.SlbId != slbId && rs.Applcation.ExtSlbId != slbId {
		return Error("can not find slb:" + slbId)
	}

	request := edas.CreateListSlbRequest()
	request.RegionId = regionId

	raw, err = edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListSlb(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_slb_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response := raw.(*edas.ListSlbResponse)
	if response.Code != 200 {
		return Error("List Slb failed for " + response.Message)
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

	appId := d.Get("app_id").(string)
	regionId := client.RegionId
	slbId := d.Get("slb_id").(string)
	slbType := d.Get("type").(string)

	request := edas.CreateUnbindSlbRequest()
	request.RegionId = regionId
	request.AppId = appId
	request.SlbId = slbId
	request.Type = slbType

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.UnbindSlb(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_slb_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response := raw.(*edas.UnbindSlbResponse)
	if response.Code != 200 {
		return Error("unbind slb failed," + response.Message)
	}

	return nil
}

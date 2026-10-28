package edas

import (
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudEdasApplicationPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasApplicationPackageAttachmentCreate,
		Read:   resourceAlicloudEdasApplicationPackageAttachmentRead,
		Delete: resourceAlicloudEdasApplicationPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"war_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"last_package_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEdasApplicationPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	groupId := d.Get("group_id").(string)

	request := edas.CreateDeployApplicationRequest()
	request.RegionId = client.RegionId
	request.AppId = appId
	request.GroupId = groupId
	request.WarUrl = d.Get("war_url").(string)
	var packageVersion string
	if v, ok := d.GetOk("package_version"); ok {
		packageVersion = v.(string)
	} else {
		packageVersion = strconv.FormatInt(time.Now().Unix(), 10)
	}
	request.DeployType = "url"
	request.PackageVersion = packageVersion

	if version, err := edasService.GetLastPackgeVersion(appId, groupId); err != nil {
		return errs.WrapError(err)
	} else {
		d.Set("last_package_version", version)
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.DeployApplication(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.DeployApplicationResponse)
	changeOrderId := response.ChangeOrderId
	if response.Code != 200 {
		return errs.WrapError(errs.Error("deploy application failed for %s", response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := helper.BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	d.SetId(appId + ":" + packageVersion)

	return resourceAlicloudEdasApplicationPackageAttachmentRead(d, meta)
}

func resourceAlicloudEdasApplicationPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := strings.Split(d.Id(), ":")[0]

	request := edas.CreateQueryApplicationStatusRequest()
	request.RegionId = regionId
	request.AppId = appId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.QueryApplicationStatus(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_edas_application_package_attachment", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)
	response, _ := raw.(*edas.QueryApplicationStatusResponse)

	if response.Code != 200 {
		return errs.WrapError(errs.Error("QueryApplicationStatus failed for %s", response.Message))
	}

	groupId := d.Get("group_id").(string)
	for _, group := range response.AppInfo.GroupList.Group {
		if group.GroupId == groupId {
			d.SetId(appId + ":" + group.PackageVersionId)
		}
	}

	return nil
}

func resourceAlicloudEdasApplicationPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := d.Get("app_id").(string)
	packageVersion := d.Get("last_package_version").(string)
	groupId := d.Get("group_id").(string)

	if len(packageVersion) == 0 {
		return nil
	}

	request := edas.CreateRollbackApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId
	request.HistoryVersion = packageVersion
	request.GroupId = groupId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.RollbackApplication(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.RollbackApplicationResponse)
	changeOrderId := response.ChangeOrderId
	if response.Code != 200 && !strings.Contains(response.Message, "ex.app.deploy.group.empty") {
		return errs.WrapError(errs.Error("deploy application failed for %s", response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := helper.BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return nil
}

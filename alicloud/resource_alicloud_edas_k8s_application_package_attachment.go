package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEdasK8sApplicationPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasK8sApplicationPackageAttachmentCreate,
		Read:   resourceAlicloudEdasK8sApplicationPackageAttachmentRead,
		Delete: resourceAlicloudEdasK8sApplicationPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pre_stop": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"envs": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"image_tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"batch_wait_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"command": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"post_start": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"liveness": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"readiness": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"args": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
			},
			"replicas": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"image_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"limit_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"limit_memory": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"request_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"request_mem": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"nas_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mount_descs": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"local_volume": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"package_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"jdk": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"web_container": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"edas_container_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"uri_encoding": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"use_body_encoding": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"update_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"requests_m_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"limit_m_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"volumes_str": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"package_version_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"change_order_desc": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEdasK8sApplicationPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	var packageVersion string
	request := edas.CreateDeployK8sApplicationRequest()
	request.RegionId = client.RegionId
	request.AppId = appId
	request.Replicas = requests.NewInteger(d.Get("replicas").(int))
	packageType, err := edasService.QueryK8sAppPackageType(appId)
	if err != nil {
		return WrapError(err)
	}
	if strings.ToLower(packageType) == "image" {
		var image string
		if v, ok := d.GetOk("image_url"); ok {
			image = v.(string)
		}
		if len(image) == 0 {
			return WrapError(Error("image_url needed for image type application"))
		}
		request.Image = image
	} else {
		if v, ok := d.GetOk("package_url"); !ok {
			return WrapError(Error("package_url is needed for creating fatjar k8s application"))
		} else {
			request.PackageUrl = v.(string)
		}
		if v, ok := d.GetOk("package_version"); !ok {
			packageVersion = strconv.FormatInt(time.Now().Unix(), 10)
		} else {
			packageVersion = v.(string)
		}
		request.PackageVersion = packageVersion
		if v, ok := d.GetOk("jdk"); !ok {
			return WrapError(Error("jdk is needed for creating non-image k8s application"))
		} else {
			request.JDK = v.(string)
		}
		if strings.ToLower(packageType) == "war" {
			var webContainer string
			var edasContainer string
			if v, ok := d.GetOk("web_container"); ok {
				webContainer = v.(string)
			}
			if v, ok := d.GetOk("edas_container_version"); ok {
				edasContainer = v.(string)
			}
			if len(webContainer) == 0 && len(edasContainer) == 0 {
				return WrapError(Error("web_container or edas_container_version is needed for creating war k8s application"))
			}
			request.WebContainer = webContainer
			request.EdasContainerVersion = edasContainer
		}
	}

	if v, ok := d.GetOk("pre_stop"); ok {
		request.PreStop = v.(string)
	}

	if v, ok := d.GetOk("envs"); ok {
		envs, err := edasService.GetK8sEnvs(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.Envs = envs
	}

	if v, ok := d.GetOk("batch_wait_time"); ok {
		request.BatchWaitTime = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("command"); ok {
		request.Command = v.(string)
	}

	if v, ok := d.GetOk("post_start"); ok {
		request.PostStart = v.(string)
	}

	if v, ok := d.GetOk("liveness"); ok {
		request.Liveness = v.(string)
	}

	if v, ok := d.GetOk("readiness"); ok {
		request.Readiness = v.(string)
	}

	if v, ok := d.GetOk("args"); ok {
		commands, err := edasService.GetK8sCommandArgs(v.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.Args = commands
	}

	if v, ok := d.GetOk("limit_cpu"); ok {
		request.CpuLimit = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("limit_mem"); ok {
		request.MemoryLimit = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("requests_cpu"); ok {
		request.CpuRequest = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("requests_mem"); ok {
		request.MemoryRequest = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("nas_id"); ok {
		request.NasId = v.(string)
	}

	if v, ok := d.GetOk("mount_descs"); ok {
		request.MountDescs = v.(string)
	}

	if v, ok := d.GetOk("local_volume"); ok {
		request.LocalVolume = v.(string)
	}

	if v, ok := d.GetOk("uri_encoding"); ok {
		request.UriEncoding = v.(string)
	}

	if v, ok := d.GetOk("use_body_encoding"); ok {
		request.UseBodyEncoding = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("update_strategy"); ok {
		request.UpdateStrategy = v.(string)
	}

	if v, ok := d.GetOk("requests_m_cpu"); ok {
		request.McpuRequest = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("limit_m_cpu"); ok {
		request.McpuLimit = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("volumes_str"); ok {
		request.VolumesStr = v.(string)
	}

	if v, ok := d.GetOk("package_version_id"); ok {
		request.PackageVersionId = v.(string)
	}

	if v, ok := d.GetOk("change_order_desc"); ok {
		request.ChangeOrderDesc = v.(string)
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.DeployK8sApplication(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.DeployK8sApplicationResponse)
	changeOrderId := response.ChangeOrderId
	if response.Code != 200 {
		return WrapError(Error("deploy k8s application failed for " + response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	d.SetId(appId + ":" + packageVersion)

	return resourceAlicloudEdasK8sApplicationPackageAttachmentRead(d, meta)
}

func resourceAlicloudEdasK8sApplicationPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_application_package_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	response, _ := raw.(*edas.QueryApplicationStatusResponse)

	if response.Code != 200 {
		return WrapError(Error("QueryApplicationStatus failed for " + response.Message))
	}

	for _, group := range response.AppInfo.GroupList.Group {
		d.SetId(appId + ":" + group.PackageVersionId)
	}

	return nil
}

func resourceAlicloudEdasK8sApplicationPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	// do nothing
	return nil
}

package alicloud

import (
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEdasK8sApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasK8sApplicationCreate,
		Read:   resourceAlicloudEdasK8sApplicationRead,
		Update: resourceAlicloudEdasK8sApplicationUpdate,
		Delete: resourceAlicloudEdasK8sApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"image_url": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"package_url"},
			},
			"package_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"FatJar", "War", "Image"}, false),
			},
			"intranet_target_port": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"intranet_slb_port": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"application_descriotion": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"repo_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"limit_cpu": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"limit_mem": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"requests_cpu": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"requests_mem": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"command": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"command_args": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"intranet_slb_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
			},
			"intranet_slb_id": {
				Optional: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"internet_slb_id": {
				Optional: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"internet_slb_protocol": {
				Optional:     true,
				Type:         schema.TypeString,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
			},
			"internet_slb_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"internet_target_port": {
				Type:     schema.TypeInt,
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
			"pre_stop": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"post_start": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"liveness": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"readiness": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"package_url": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"image_url"},
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Computed: true,
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
			"requests_m_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"limit_m_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEdasK8sApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	var packageVersion string

	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}
	request := edas.CreateDeployK8sApplicationRequest()
	request.AppId = d.Id()
	request.RegionId = client.RegionId

	packageType, err := edasService.QueryK8sAppPackageType(d.Id())
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

	if v, ok := d.GetOk("replicas"); ok {
		request.Replicas = requests.NewInteger(v.(int))
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
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudEdasK8sApplicationRead(d, meta)
}

func resourceAlicloudEdasK8sApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}
	request := edas.CreateInsertK8sApplicationRequest()

	packageType := d.Get("package_type").(string)

	request.AppName = d.Get("application_name").(string)
	request.RegionId = client.RegionId
	request.PackageType = packageType
	request.ClusterId = d.Get("cluster_id").(string)

	if strings.ToLower(packageType) == "image" {
		if v, ok := d.GetOk("image_url"); !ok {
			return WrapError(Error("image_url is needed for creating image k8s application"))
		} else {
			request.ImageUrl = v.(string)
		}
		if v, ok := d.GetOk("repo_id"); ok {
			request.RepoId = v.(string)
		}
	} else {
		if v, ok := d.GetOk("package_url"); !ok {
			return WrapError(Error("package_url is needed for creating fatjar k8s application"))
		} else {
			request.PackageUrl = v.(string)
		}
		if v, ok := d.GetOk("package_version"); !ok {
			request.PackageVersion = strconv.FormatInt(time.Now().Unix(), 10)
		} else {
			request.PackageVersion = v.(string)
		}
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

	request.Replicas = requests.NewInteger(d.Get("replicas").(int))

	if v, ok := d.GetOk("intranet_target_port"); ok {
		request.IntranetTargetPort = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("intranet_slb_port"); ok {
		request.InternetSlbPort = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("application_descriotion"); ok {
		request.ApplicationDescription = v.(string)
	}

	if v, ok := d.GetOk("limit_cpu"); ok {
		request.LimitCpu = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("limit_mem"); ok {
		request.LimitMem = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("requests_cpu"); ok {
		request.RequestsCpu = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("requests_mem"); ok {
		request.RequestsMem = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("command"); ok {
		request.Command = v.(string)
	}

	if v, ok := d.GetOk("command_args"); ok {
		commands, err := edasService.GetK8sCommandArgs(v.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.CommandArgs = commands
	}

	if v, ok := d.GetOk("intranet_slb_protocol"); ok {
		request.IntranetSlbProtocol = v.(string)
	}

	if v, ok := d.GetOk("intranet_slb_id"); ok {
		request.IntranetSlbId = v.(string)
	}

	if v, ok := d.GetOk("internet_slb_id"); ok {
		request.InternetSlbId = v.(string)
	}

	if v, ok := d.GetOk("internet_slb_protocol"); ok {
		request.InternetSlbProtocol = v.(string)
	}

	if v, ok := d.GetOk("internet_slb_port"); ok {
		request.InternetSlbPort = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("internet_target_port"); ok {
		request.InternetTargetPort = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("envs"); ok {
		envs, err := edasService.GetK8sEnvs(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.Envs = envs
	}

	if v, ok := d.GetOk("pre_stop"); ok {
		request.PreStop = v.(string)
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

	if v, ok := d.GetOk("nas_id"); ok {
		request.NasId = v.(string)
	}

	if v, ok := d.GetOk("mount_descs"); ok {
		request.MountDescs = v.(string)
	}

	if v, ok := d.GetOk("local_volume"); ok {
		request.LocalVolume = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = v.(string)
	}

	if v, ok := d.GetOk("logical_region_id"); ok {
		request.LogicalRegionId = v.(string)
	}

	if v, ok := d.GetOk("uri_encoding"); ok {
		request.UriEncoding = v.(string)
	}

	if v, ok := d.GetOk("use_body_encoding"); ok {
		request.UseBodyEncoding = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("requests_m_cpu"); ok {
		request.RequestsmCpu = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("limit_m_cpu"); ok {
		request.LimitmCpu = requests.NewInteger(v.(int))
	}

	var appId string
	var changeOrderId string

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.InsertK8sApplication(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_k8s_application", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.InsertK8sApplicationResponse)
	appId = response.ApplicationInfo.AppId
	changeOrderId = response.ApplicationInfo.ChangeOrderId
	d.SetId(appId)
	if response.Code != 200 {
		return WrapError(Error("Create k8s application failed for " + response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudEdasK8sApplicationRead(d, meta)
}

func resourceAlicloudEdasK8sApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	response, err := edasService.DescribeEdasK8sApplication(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_edas_k8s_application ecsService.DescribeEdasK8sApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("application_name", response.App.ApplicationName)
	d.Set("application_descriotion", response.App.Description)
	d.Set("cluster_id", response.App.ClusterId)
	d.Set("replicas", response.App.Instances)
	d.Set("package_type", response.App.ApplicationType)
	d.Set("image_url", response.ImageInfo.ImageUrl)

	d.Set("command", response.App.Cmd)
	d.Set("command_args", response.App.CmdArgs.CmdArg)

	d.Set("repo_id", response.ImageInfo.RepoId)
	d.Set("package_url", "")
	d.Set("package_version", "")
	allDeploy := response.DeployGroups.DeployGroup
	for v := range allDeploy {
		d.Set("package_url", allDeploy[v].PackageUrl)
		d.Set("package_version", allDeploy[v].PackageVersion)
		d.Set("limit_cpu", allDeploy[v].CpuLimit)
		d.Set("limit_mem", allDeploy[v].MemoryLimit)

		d.Set("requests_mem", allDeploy[v].MemoryRequest)

		cpu, err := strconv.Atoi(allDeploy[v].CpuRequest)

		if err != nil {
			if NotFoundError(err) {
				log.Printf("[DEBUG] Resource alicloud_edas_k8s_application ecsService.DescribeEdasK8sApplication Failed!!! %s", err)
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}
		d.Set("requests_m_cpu", cpu*1000)
		d.Set("requests_cpu", cpu)
	}

	d.Set("edas_container_version", response.App.EdasContainerVersion)
	d.Set("pre_stop", response.Conf.PreStop)
	d.Set("post_start", response.Conf.PostStart)
	d.Set("liveness", response.Conf.Liveness)
	d.Set("readiness", response.Conf.Readiness)
	d.Set("namespace", response.NameSpace)
	d.Set("logical_region_id", response.RegionId)
	return nil
}

func resourceAlicloudEdasK8sApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateDeleteK8sApplicationRequest()
	request.RegionId = client.RegionId
	request.AppId = d.Id()

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteK8sApplication(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response := raw.(*edas.DeleteK8sApplicationResponse)
		if response.Code != 200 {
			return resource.NonRetryableError(Error("Delete k8s application failed for " + response.Message))
		}
		changeOrderId := response.ChangeOrderId

		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return resource.NonRetryableError(WrapErrorf(err, IdMsg, d.Id()))
			}
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

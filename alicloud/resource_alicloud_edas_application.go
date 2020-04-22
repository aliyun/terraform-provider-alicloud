package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEdasApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasApplicationCreate,
		Update: resourceAlicloudEdasApplicationUpdate,
		Read:   resourceAlicloudEdasApplicationRead,
		Delete: resourceAlicloudEdasApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"package_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"JAR", "WAR"}, false),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"build_pack_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"descriotion": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"component_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ecu_info": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},

			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"war_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEdasApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	applicationName := d.Get("application_name").(string)
	regionId := client.RegionId
	clusterId := d.Get("cluster_id").(string)
	ecuInfo := d.Get("ecu_info").([]interface{})
	aString := make([]string, len(ecuInfo))
	for i, v := range ecuInfo {
		aString[i] = v.(string)
	}
	packageType := d.Get("package_type").(string)
	buildPackId := d.Get("build_pack_id").(int)
	description := d.Get("descriotion").(string)
	healthCheckUrl := d.Get("health_check_url").(string)
	logicalRegionId := d.Get("logical_region_id").(string)
	//componentIds := d.Get("component_ids").(string)

	request := edas.CreateInsertApplicationRequest()
	request.RegionId = regionId
	request.ApplicationName = applicationName
	request.ClusterId = clusterId
	request.PackageType = packageType
	request.EcuInfo = strings.Join(aString, ",")
	request.BuildPackId = requests.NewInteger(buildPackId)
	request.Description = description
	request.HealthCheckURL = healthCheckUrl
	request.LogicalRegionId = logicalRegionId

	var appId string
	var changeOrderId string

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.InsertApplication(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_application", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.InsertApplicationResponse)
	appId = response.ApplicationInfo.AppId
	changeOrderId = response.ApplicationInfo.ChangeOrderId
	d.SetId(appId)
	if response.Code != 200 {
		return Error("create application failed for " + response.Message)
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudEdasApplicationUpdate(d, meta)
}

func resourceAlicloudEdasApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := d.Id()
	packageVersion := d.Get("package_version").(string)
	groupId := d.Get("group_id").(string)
	warUlr := d.Get("war_url").(string)

	if len(warUlr) == 0 || len(groupId) == 0 {
		return nil
	}

	if len(packageVersion) == 0 {
		packageVersion = strconv.Itoa(time.Now().Second())
	}
	request := edas.CreateDeployApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId
	request.GroupId = groupId
	request.PackageVersion = packageVersion
	request.DeployType = "url"

	request.WarUrl = warUlr

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.DeployApplication(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.DeployApplicationResponse)
	changeOrderId := response.ChangeOrderId
	if response.Code != 200 {
		return Error("deploy application failed for " + response.Message)
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudEdasApplicationRead(d, meta)
}

func resourceAlicloudEdasApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := d.Id()

	request := edas.CreateGetApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetApplication(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response, _ := raw.(*edas.GetApplicationResponse)
		d.Set("application_name", response.Applcation.Name)
		d.Set("cluster_id", response.Applcation.ClusterId)
		d.Set("build_pack_id", response.Applcation.BuildPackageId)
		d.Set("descriotion", response.Applcation.Description)
		d.Set("health_check_url", response.Applcation.HealthCheckUrl)
		if len(response.Applcation.ApplicationType) > 0 && response.Applcation.ApplicationType == "FatJar" {
			d.Set("package_type", "JAR")
		} else {
			d.Set("package_type", response.Applcation.ApplicationType)
		}

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_application", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}

func resourceAlicloudEdasApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := d.Id()

	//packageVersion := d.Get("package_version").(string)
	//groupId := d.Get("group_id").(string)
	//warUlr := d.Get("war_url").(string)
	if true {
		request := edas.CreateStopApplicationRequest()
		request.RegionId = regionId
		request.AppId = appId

		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.StopApplication(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response, _ := raw.(*edas.StopApplicationResponse)
		changeOrderId := response.ChangeOrderId

		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}

	request := edas.CreateDeleteApplicationRequest()
	request.RegionId = regionId
	request.AppId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteApplication(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		rsp := raw.(*edas.DeleteApplicationResponse)
		if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
			err = Error("Operation cannot be processed because there are running instances.")
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

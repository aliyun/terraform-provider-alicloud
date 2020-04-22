package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudEdasApplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEdasApplicationsRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build_package_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"running_instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"slb_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEdasApplicationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateListApplicationRequest()
	request.RegionId = client.RegionId

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			idsMap[Trim(id.(string))] = Trim(id.(string))
		}
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListApplication(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_applications", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ListApplicationResponse)
	if response.Code != 200 {
		return WrapError(Error(response.Message))
	}
	var filteredApps []edas.Application
	if len(idsMap) > 0 {
		for _, app := range response.ApplicationList.Application {
			if _, ok := idsMap[app.AppId]; ok {
				filteredApps = append(filteredApps, app)
			}
		}
	} else {
		filteredApps = response.ApplicationList.Application
	}

	return edasApplicationAttributes(d, filteredApps)
}

func edasApplicationAttributes(d *schema.ResourceData, apps []edas.Application) error {
	var appIds []string
	var s []map[string]interface{}

	for _, app := range apps {
		mapping := map[string]interface{}{
			"app_name":               app.Name,
			"app_id":                 app.AppId,
			"application_type":       app.ApplicationType,
			"build_package_id":       app.BuildPackageId,
			"cluster_id":             app.ClusterId,
			"cluster_type":           app.ClusterType,
			"instance_count":         app.InstanceCount,
			"running_instance_count": app.RunningInstanceCount,
			"health_check_url":       app.HealthCheckUrl,
			"create_time":            app.CreateTime,
			"slb_id":                 app.SlbId,
			"slb_ip":                 app.SlbIp,
			"slb_port":               app.SlbPort,
			"region_id":              app.RegionId,
		}
		appIds = append(appIds, app.AppId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(appIds))
	if err := d.Set("applications", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

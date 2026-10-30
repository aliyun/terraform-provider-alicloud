package edas

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudEdasApplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEdasApplicationsRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			if id == nil {
				continue
			}
			idsMap[helper.Trim(id.(string))] = helper.Trim(id.(string))
		}
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListApplication(request)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_edas_applications", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ListApplicationResponse)
	if response.Code != 200 {
		return errs.WrapError(errs.Error("%s", response.Message))
	}
	var filteredApps []edas.ApplicationInListApplication
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return errs.WrapError(err)
			}
		}
		for _, app := range response.ApplicationList.Application {
			if r != nil && !r.MatchString(app.Name) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[app.AppId]; !ok {
					continue
				}
			}
			filteredApps = append(filteredApps, app)
		}
	} else {
		filteredApps = response.ApplicationList.Application
	}

	return edasApplicationAttributes(d, filteredApps)
}

func edasApplicationAttributes(d *schema.ResourceData, apps []edas.ApplicationInListApplication) error {
	var appIds []string
	var s []map[string]interface{}
	var names []string

	for _, app := range apps {
		mapping := map[string]interface{}{
			"app_name":         app.Name,
			"app_id":           app.AppId,
			"application_type": app.ApplicationType,
			"build_package_id": app.BuildPackageId,
			"cluster_id":       app.ClusterId,
			"cluster_type":     app.ClusterType,
			"region_id":        app.RegionId,
		}
		appIds = append(appIds, app.AppId)
		names = append(names, app.Name)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(appIds))
	if err := d.Set("ids", appIds); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("applications", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

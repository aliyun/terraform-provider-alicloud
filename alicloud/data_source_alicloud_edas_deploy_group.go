package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudEdasDeployGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEdasDeployGroupsRead,

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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

			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEdasDeployGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	regionId := client.RegionId
	appId := d.Get("app_id").(string)
	request := edas.CreateListDeployGroupRequest()
	request.RegionId = regionId
	request.AppId = appId

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			idsMap[Trim(id.(string))] = Trim(id.(string))
		}
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListDeployGroup(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_deploy_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ListDeployGroupResponse)

	if response.Code != 200 {
		return WrapError(Error(response.Message))
	}

	var filteredGroups []edas.DeployGroup
	if len(idsMap) > 0 {
		for _, group := range response.DeployGroupList.DeployGroup {
			if _, ok := idsMap[group.GroupId]; ok {
				filteredGroups = append(filteredGroups, group)
			}
		}
	} else {
		filteredGroups = response.DeployGroupList.DeployGroup
	}

	return edasDeployGroupAttributes(d, filteredGroups)
}

func edasDeployGroupAttributes(d *schema.ResourceData, groups []edas.DeployGroup) error {
	var Ids []string
	var s []map[string]interface{}

	for _, group := range groups {
		mapping := map[string]interface{}{
			"group_id":           group.GroupId,
			"group_name":         group.GroupName,
			"group_type":         group.GroupType,
			"create_time":        group.CreateTime,
			"update_time":        group.UpdateTime,
			"app_id":             group.AppId,
			"cluster_id":         group.ClusterId,
			"package_version_id": group.PackageVersionId,
			"app_version_id":     group.AppVersionId,
		}
		Ids = append(Ids, group.GroupId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(Ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

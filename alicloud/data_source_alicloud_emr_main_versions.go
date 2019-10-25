package alicloud

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudEmrMainVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEmrMainVersionsRead,

		Schema: map[string]*schema.Schema{
			"emr_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"main_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emr_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlicloudEmrMainVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := emr.CreateListEmrMainVersionRequest()
	if emrVersion, ok := d.GetOk("emr_version"); ok {
		request.EmrVersion = strings.TrimSpace(emrVersion.(string))
	}

	raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
		return emrClient.ListEmrMainVersion(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emr_main_versions", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	var (
		mainVersions []emr.EmrMainVersion
		clusterTypes = make(map[string][]string)
	)
	response, _ := raw.(*emr.ListEmrMainVersionResponse)
	if response != nil {
		// get clusterInfo of specific emr version
		var (
			versionRequest  = emr.CreateDescribeEmrMainVersionRequest()
			versionResponse *emr.DescribeEmrMainVersionResponse
			versionRaw      interface{}
		)
		for _, v := range response.EmrMainVersionList.EmrMainVersion {
			if v.EmrVersion == "" {
				continue
			}
			versionRequest.EmrVersion = v.EmrVersion
			versionRaw, err = client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
				return emrClient.DescribeEmrMainVersion(versionRequest)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emr_main_versions", request.GetActionName(), AlibabaCloudSdkGoERROR)
			}

			versionResponse, _ = versionRaw.(*emr.DescribeEmrMainVersionResponse)
			if versionResponse == nil {
				continue
			}
			types := []string{}
			for _, c := range versionResponse.EmrMainVersion.ClusterTypeInfoList.ClusterTypeInfo {
				types = append(types, c.ClusterType)
			}
			if len(types) == 0 {
				continue
			}
			clusterTypes[v.EmrVersion] = types
		}

		mainVersions = response.EmrMainVersionList.EmrMainVersion
	}

	return emrClusterMainVersionAttributes(d, clusterTypes, mainVersions)
}

func emrClusterMainVersionAttributes(d *schema.ResourceData, clusterTypes map[string][]string, mainVersions []emr.EmrMainVersion) error {
	var (
		ids []string
		s   []map[string]interface{}
	)

	for _, version := range mainVersions {
		// if display is false, ignore it
		if !version.Display {
			continue
		}

		mapping := map[string]interface{}{
			"image_id":      version.ImageId,
			"emr_version":   version.EmrVersion,
			"cluster_types": clusterTypes[version.EmrVersion],
		}

		s = append(s, mapping)
		ids = append(ids, version.EmrVersion)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("main_versions", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEmrMainVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEmrMainVersionsRead,

		Schema: map[string]*schema.Schema{
			"emr_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	emrService := EmrService{client}
	action := "ListEmrMainVersion"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("emr_version"); ok {
		request["EmrVersion"] = v
	}

	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	clusterType, clusterTypeOk := d.GetOk("cluster_type")
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var response map[string]interface{}
	var objects []map[string]interface{}
	conn, err := client.NewEmrClient()
	if err != nil {
		return WrapError(err)
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emr_main_versions", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.EmrMainVersionList.EmrMainVersion", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.EmrMainVersionList.EmrMainVersion", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["EmrVersion"])]; !ok {
					continue
				}
			}

			clusterTypeFilter := func(filter []interface{}, source []interface{}) (result []string) {
				if len(source) == 0 {
					return
				}
				if len(filter) == 0 {
					for _, v := range source {
						clusterType := fmt.Sprint(v.(map[string]interface{})["ClusterType"])
						if "CLICKHOUSE" == clusterType { // emr cluster 'CLICKHOUSE' not supported, ignore it.
							continue
						}
						result = append(result, clusterType)
					}
					return
				}

				sourceMapping := make(map[string]bool, 0)
				for _, v := range source {
					clusterType := fmt.Sprint(v.(map[string]interface{})["ClusterType"])
					if "CLICKHOUSE" == clusterType { // emr cluster 'CLICKHOUSE' not supported, ignore it.
						continue
					}
					sourceMapping[clusterType] = true
				}

				for _, f := range filter {
					if v, _ := sourceMapping[f.(string)]; !v {
						return nil
					}
					result = append(result, f.(string))
				}
				return
			}

			source, err := emrService.DescribeEmrMainVersionClusterTypes(fmt.Sprint(item["EmrVersion"]))
			if err != nil {
				return WrapError(err)
			}
			clusterTypes := clusterTypeFilter(clusterType.([]interface{}), source)
			if clusterTypeOk && len(clusterType.([]interface{})) > 0 && len(clusterTypes) == 0 {
				continue
			}

			item["ClusterTypes"] = clusterTypes
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"emr_version":   fmt.Sprint(object["EmrVersion"]),
			"image_id":      fmt.Sprint(object["ImageId"]),
			"cluster_types": object["ClusterTypes"],
		}
		ids = append(ids, fmt.Sprint(mapping["emr_version"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("main_versions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

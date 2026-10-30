package emr

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEmrMainVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEmrMainVersionsRead,

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

func dataSourceEmrMainVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	emrService := NewService(client)
	action := "ListEmrMainVersion"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("emr_version"); ok {
		request["EmrVersion"] = v
	}

	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
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
	var err error

	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Emr", "2016-04-08", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_emr_main_versions", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.EmrMainVersionList.EmrMainVersion", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.EmrMainVersionList.EmrMainVersion", response)
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
				return errs.WrapError(err)
			}
			clusterTypes := clusterTypeFilter(clusterType.([]interface{}), source)
			if clusterTypeOk && len(clusterType.([]interface{})) > 0 && len(clusterTypes) == 0 {
				continue
			}

			item["ClusterTypes"] = clusterTypes
			objects = append(objects, item)
		}
		if len(result) < helper.PageSizeLarge {
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

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("main_versions", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudArmsRemoteWrites() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudArmsRemoteWritesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"remote_writes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_write_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_write_yaml": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudArmsRemoteWritesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPrometheusRemoteWrites"
	request := make(map[string]interface{})

	request["RegionId"] = client.RegionId
	request["ClusterId"] = d.Get("cluster_id")

	var objects []map[string]interface{}
	var remoteWriteNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		remoteWriteNameRegex = r
	}

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
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_arms_remote_writes", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if remoteWriteNameRegex != nil && !remoteWriteNameRegex.MatchString(fmt.Sprint(item["RemoteWriteName"])) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprintf("%v:%v", item["ClusterId"], item["RemoteWriteName"])]; !ok {
				continue
			}
		}

		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprintf("%v:%v", object["ClusterId"], object["RemoteWriteName"]),
			"cluster_id":        fmt.Sprint(object["ClusterId"]),
			"remote_write_name": fmt.Sprint(object["RemoteWriteName"]),
			"remote_write_yaml": object["RemoteWriteYaml"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["RemoteWriteName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("remote_writes", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

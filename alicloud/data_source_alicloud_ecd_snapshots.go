package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcdSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"desktop_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"snapshot_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"snapshots": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"desktop_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"progress": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"remain_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"snapshot_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"snapshot_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"snapshot_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_disk_size": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_disk_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcdSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	request["MaxResults"] = PageSizeLarge
	if v, ok := d.GetOk("desktop_id"); ok {
		request["DesktopId"] = v
	}

	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
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

	var snapshotNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		snapshotNameRegex = r
	}

	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeSnapshots"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_snapshots", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Snapshots", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Snapshots", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SnapshotId"])]; !ok {
					continue
				}
			}

			if snapshotNameRegex != nil && !snapshotNameRegex.MatchString(fmt.Sprint(item["SnapshotName"])) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":               fmt.Sprint(object["SnapshotId"]),
			"create_time":      object["CreationTime"],
			"description":      object["Description"],
			"desktop_id":       object["DesktopId"],
			"progress":         object["Progress"],
			"remain_time":      object["RemainTime"],
			"snapshot_id":      object["SnapshotId"],
			"snapshot_name":    object["SnapshotName"],
			"snapshot_type":    object["SnapshotType"],
			"source_disk_size": object["SourceDiskSize"],
			"source_disk_type": object["SourceDiskType"],
			"status":           object["Status"],
		}
		ids = append(ids, fmt.Sprint(object["SnapshotId"]))
		names = append(names, object["SnapshotName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("snapshots", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

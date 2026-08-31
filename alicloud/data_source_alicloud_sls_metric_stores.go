package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudSlsMetricStores() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSlsMetricStoresRead,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_stores": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_store_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"shard_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_split": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_split_shard_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"append_meta": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hot_ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_modify_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudSlsMetricStoresRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	projectName := d.Get("project_name").(string)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	action := "/logstores"
	query := make(map[string]*string)
	telemetryType := "Metrics"
	query["telemetryType"] = &telemetryType
	hostMap := map[string]*string{
		"project": StringPointer(projectName),
	}

	var response map[string]interface{}
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, err := client.Do("Sls", roaParam("GET", "2020-12-30", "ListLogStores", action), query, nil, nil, hostMap, true)
		if err != nil {
			if NeedRetry(err) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_metric_stores", "ListLogStores", AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response)

	logstoresRaw, _ := response["logstores"].([]interface{})

	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	slsServiceV2 := SlsServiceV2{client}
	for _, v := range logstoresRaw {
		name, ok := v.(string)
		if !ok || name == "" {
			continue
		}
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		id := fmt.Sprintf("%s:%s", projectName, name)
		if len(idsMap) > 0 {
			_, hitID := idsMap[id]
			_, hitName := idsMap[name]
			if !hitID && !hitName {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":                id,
			"metric_store_name": name,
		}
		objectRaw, err := slsServiceV2.DescribeSlsLogStore(id)
		if err == nil {
			if v, ok := objectRaw["ttl"]; ok {
				mapping["ttl"] = v
			}
			if v, ok := objectRaw["shardCount"]; ok {
				mapping["shard_count"] = v
			}
			if v, ok := objectRaw["autoSplit"]; ok {
				mapping["auto_split"] = v
			}
			if v, ok := objectRaw["maxSplitShard"]; ok {
				mapping["max_split_shard_count"] = v
			}
			if v, ok := objectRaw["appendMeta"]; ok {
				mapping["append_meta"] = v
			}
			if v, ok := objectRaw["hot_ttl"]; ok {
				mapping["hot_ttl"] = v
			}
			if v, ok := objectRaw["mode"]; ok {
				mapping["mode"] = v
			}
			if v, ok := objectRaw["createTime"]; ok {
				mapping["create_time"] = v
			}
			if v, ok := objectRaw["lastModifyTime"]; ok {
				mapping["last_modify_time"] = v
			}
		}
		ids = append(ids, id)
		names = append(names, name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("metric_stores", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

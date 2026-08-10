package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudSlsMetricStores() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSlsMetricStoreRead,
		Schema: map[string]*schema.Schema{
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
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_store_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ValidateFunc: validation.IntBetween(1, 500),
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
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_store_name": {
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
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"shard_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_split_shard": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_split": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"metric_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudSlsMetricStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

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

	projectName := d.Get("project_name").(string)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/metricstores")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(projectName)
	if v, ok := d.GetOk("metric_store_name"); ok {
		query["name"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("mode"); ok {
		query["mode"] = StringPointer(v.(string))
	}
	query["offset"] = StringPointer(strconv.Itoa(d.Get("offset").(int)))
	query["size"] = StringPointer(strconv.Itoa(d.Get("size").(int)))

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("GET", "2020-12-30", "ListMetricStores", action), query, nil, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"ProjectNotExist"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	resp, _ := jsonpath.Get("$.metricstores", response)
	result, _ := resp.([]interface{})

	var names []string
	for _, v := range result {
		name := fmt.Sprint(v)
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprintf("%s:%s", projectName, name)]; !ok {
				continue
			}
		}
		names = append(names, name)
	}

	ids := make([]string, 0)
	namesOut := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, name := range names {
		id := fmt.Sprintf("%s:%s", projectName, name)
		objectRaw, err := slsServiceV2.DescribeSlsMetricStore(id)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		mapping := map[string]interface{}{}
		mapping["id"] = id
		mapping["project_name"] = projectName
		mapping["metric_store_name"] = objectRaw["name"]
		mapping["create_time"] = objectRaw["createTime"]
		mapping["last_modify_time"] = objectRaw["lastModifyTime"]
		mapping["ttl"] = objectRaw["ttl"]
		mapping["shard_count"] = objectRaw["shardCount"]
		mapping["max_split_shard"] = objectRaw["maxSplitShard"]
		mapping["auto_split"] = objectRaw["autoSplit"]
		mapping["metric_type"] = objectRaw["metricType"]
		mapping["mode"] = objectRaw["mode"]

		ids = append(ids, id)
		namesOut = append(namesOut, name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", namesOut); err != nil {
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

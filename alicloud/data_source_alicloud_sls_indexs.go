package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudSlsIndexs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSlsIndexRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"logstore_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"indexs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keys": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include_keys": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"chn": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"case_sensitive": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"exclude_keys": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"token": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"log_reduce_black_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"log_reduce_white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_text_len": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudSlsIndexRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	logstore := d.Get("logstore_name")
	action := fmt.Sprintf("/logstores/%s/index", logstore)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	request["logstore"] = d.Get("logstore_name")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("GET", "2020-12-30", "GetIndex", action), query, nil, nil, hostMap, true)

		if err != nil {
			if NeedRetry(err) {
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

	resp, _ := jsonpath.Get("$", response)
	if resp != nil {
		objects = append(objects, resp.(map[string]interface{}))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(*hostMap["project"], ":", logstore)

		if _, ok := idsMap[mapping["id"].(string)]; !ok {
			continue
		}

		if objectRaw["keys"] != nil {
			keys, err := convertToJsonWithoutEscapeHTML(objectRaw["keys"].(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			mapping["keys"] = keys
		}
		mapping["max_text_len"] = objectRaw["max_text_len"]
		mapping["ttl"] = objectRaw["ttl"]

		lineMaps := make([]map[string]interface{}, 0)
		lineMap := make(map[string]interface{})
		lineRaw := make(map[string]interface{})
		if objectRaw["line"] != nil {
			lineRaw = objectRaw["line"].(map[string]interface{})
		}
		if len(lineRaw) > 0 {
			lineMap["case_sensitive"] = lineRaw["caseSensitive"]
			lineMap["chn"] = lineRaw["chn"]

			exclude_keysRaw := make([]interface{}, 0)
			if lineRaw["exclude_keys"] != nil {
				exclude_keysRaw = convertToInterfaceArray(lineRaw["exclude_keys"])
			}

			lineMap["exclude_keys"] = exclude_keysRaw
			include_keysRaw := make([]interface{}, 0)
			if lineRaw["include_keys"] != nil {
				include_keysRaw = convertToInterfaceArray(lineRaw["include_keys"])
			}

			lineMap["include_keys"] = include_keysRaw
			tokenRaw := make([]interface{}, 0)
			if lineRaw["token"] != nil {
				tokenRaw = convertToInterfaceArray(lineRaw["token"])
			}

			lineMap["token"] = tokenRaw
			lineMaps = append(lineMaps, lineMap)
		}
		mapping["line"] = lineMaps
		log_reduce_black_listRaw := make([]interface{}, 0)
		if objectRaw["log_reduce_black_list"] != nil {
			log_reduce_black_listRaw = convertToInterfaceArray(objectRaw["log_reduce_black_list"])
		}

		mapping["log_reduce_black_list"] = log_reduce_black_listRaw
		log_reduce_white_listRaw := make([]interface{}, 0)
		if objectRaw["log_reduce_white_list"] != nil {
			log_reduce_white_listRaw = convertToInterfaceArray(objectRaw["log_reduce_white_list"])
		}

		mapping["log_reduce_white_list"] = log_reduce_white_listRaw

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw[""])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("indexs", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudMseNacosConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMseNacosConfigsRead,
		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"zh", "en"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"request_pars": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"md5": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted_data_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"beta_ips": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudMseNacosConfigsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListNacosConfigs"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}

	request["PageSize"] = PageSizeLarge
	request["PageNum"] = 1

	request["InstanceId"] = d.Get("instance_id")

	var requestNamespaceId = d.Get("namespace_id")
	if requestNamespaceId == nil {
		requestNamespaceId = ""
	}

	request["NamespaceId"] = requestNamespaceId

	if v, ok := d.GetOk("request_pars"); ok {
		request["RequestPars"] = v
	}

	if v, ok := d.GetOk("data_id"); ok {
		request["DataId"] = v
	}

	if v, ok := d.GetOk("group"); ok {
		request["Group"] = v
	}

	if v, ok := d.GetOk("app_name"); ok {
		request["AppName"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		request["Tags"] = v
	}

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

	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mse_nacos_configs", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Configurations", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Configurations", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				var namespaceId = request["NamespaceId"]
				if namespaceId == nil {
					namespaceId = ""
				}
				if _, ok := idsMap[fmt.Sprint(request["InstanceId"], ":", namespaceId, ":", item["DataId"], ":", item["Group"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"group":    object["Group"],
			"data_id":  object["DataId"],
			"app_name": object["AppName"],
		}

		var namespaceId = request["NamespaceId"]
		if namespaceId == nil {
			namespaceId = ""
		}

		id := fmt.Sprint(request["InstanceId"], ":", namespaceId, ":", object["DataId"], ":", object["Group"])
		mapping["id"] = id

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, id)
			s = append(s, mapping)
			continue
		}

		mseService := MseService{client}
		getResp, err := mseService.DescribeMseNacosConfig(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["type"] = getResp["Type"]
		mapping["tags"] = getResp["Tags"]
		mapping["md5"] = getResp["Md5"]
		mapping["content"] = getResp["Content"]
		mapping["desc"] = getResp["Desc"]
		mapping["encrypted_data_key"] = getResp["EncryptedDataKey"]
		mapping["beta_ips"] = getResp["BetaIps"]

		ids = append(ids, id)

		s = append(s, mapping)

	}
	//标准代码
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("configs", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

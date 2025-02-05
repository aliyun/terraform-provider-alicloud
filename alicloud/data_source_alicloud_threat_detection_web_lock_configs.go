package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionWebLockConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionWebLockConfigsRead,
		Schema: map[string]*schema.Schema{
			"lang": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"remark": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"source_ip": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"ids": {
				Optional: true,
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
			"page_number": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  20,
			},
			"configs": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"defence_mode": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dir": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"exclusive_dir": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"exclusive_file": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"exclusive_file_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"inclusive_file_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"local_backup_dir": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"mode": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"uuid": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionWebLockConfigsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}
	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["CurrentPage"] = v.(int)
	} else {
		request["CurrentPage"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
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

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeWebLockBindList"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_web_lock_configs", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.BindList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.BindList", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Uuid"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	sasService := SasService{client}
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":   fmt.Sprint(object["Uuid"]),
			"uuid": object["Uuid"],
		}

		ids = append(ids, fmt.Sprint(object["Uuid"]))

		id := fmt.Sprint(object["Uuid"])
		object, err = sasService.DescribeThreatDetectionWebLockConfig(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["defence_mode"] = object["DefenceMode"]
		mapping["dir"] = object["Dir"]
		mapping["exclusive_dir"] = object["ExclusiveDir"]
		mapping["exclusive_file"] = object["ExclusiveFile"]
		mapping["exclusive_file_type"] = object["ExclusiveFileType"]
		mapping["inclusive_file_type"] = object["InclusiveFileType"]
		mapping["local_backup_dir"] = object["LocalBackupDir"]
		mapping["mode"] = object["Mode"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("configs", s); err != nil {
		return WrapError(err)
	}
	return nil
}

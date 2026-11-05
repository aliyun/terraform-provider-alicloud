package ons

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceOnsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOnsInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 2, 5, 7}),
			},
			"tags": tags.TagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_internal_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_internet_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_internet_secure_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"release_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tcp_endpoint": {
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

func dataSourceOnsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "OnsInstanceInServiceList"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	var objects []map[string]interface{}
	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		instanceNameRegex = r
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
	status, statusOk := d.GetOkExists("status")
	var response map[string]interface{}
	var err error
	response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, true)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ons_instances", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	resp, err := jsonpath.Get("$.Data.InstanceVO", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.InstanceVO", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if instanceNameRegex != nil {
			if !instanceNameRegex.MatchString(fmt.Sprint(item["InstanceName"])) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
				continue
			}
		}
		if statusOk {
			instanceStatusVal, instanceStatusErr := convert.FormatInt(item["InstanceStatus"])
			if instanceStatusErr != nil {
				return instanceStatusErr
			}
			if status.(int) != instanceStatusVal {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		instance_typeVal, err := convert.FormatInt(object["InstanceType"])
		if err != nil {
			return err
		}
		releaseTime, err := convert.FormatInt(object["ReleaseTime"])
		if err != nil {
			return err
		}
		statusVal, err := convert.FormatInt(object["InstanceStatus"])
		if err != nil {
			return err
		}
		instance_statusVal, err := convert.FormatInt(object["InstanceStatus"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"independent_naming": object["IndependentNaming"],
			"id":                 fmt.Sprint(object["InstanceId"]),
			"instance_id":        fmt.Sprint(object["InstanceId"]),
			"instance_name":      object["InstanceName"],
			"instance_type":      instance_typeVal,
			"release_time":       fmt.Sprint(releaseTime),
			"status":             statusVal,
			"instance_status":    instance_statusVal,
		}

		tagMap := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !tags.IgnoredTags(key, value) {
					tagMap[key] = value
				}
			}
		}
		mapping["tags"] = tagMap
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["InstanceId"]))
			names = append(names, object["InstanceName"])
			s = append(s, mapping)
			continue
		}

		onsService := Service{client: client}
		id := fmt.Sprint(object["InstanceId"])
		getResp, err := onsService.DescribeOnsInstance(id)
		if err != nil {
			return errs.WrapError(err)
		}
		mapping["http_internal_endpoint"] = getResp["Endpoints"].(map[string]interface{})["HttpInternalEndpoint"]
		mapping["http_internet_endpoint"] = getResp["Endpoints"].(map[string]interface{})["HttpInternetEndpoint"]
		mapping["http_internet_secure_endpoint"] = getResp["Endpoints"].(map[string]interface{})["HttpInternetSecureEndpoint"]
		mapping["remark"] = getResp["Remark"]
		mapping["tcp_endpoint"] = getResp["Endpoints"].(map[string]interface{})["TcpEndpoint"]
		ids = append(ids, fmt.Sprint(object["InstanceId"]))
		names = append(names, object["InstanceName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

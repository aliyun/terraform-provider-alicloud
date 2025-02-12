package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddosbgp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDdosbgpInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDdosbgpInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"base_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"normal_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDdosbgpInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ddosbgp.CreateDescribeInstanceListRequest()
	request.PageSize = requests.Integer(strconv.Itoa(PageSizeLarge))
	request.PageNo = "1"
	request.RegionId = client.RegionId
	var instances []ddosbgp.Instance

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	if v, ok := d.GetOk("ids"); ok {
		idsStr, _ := json.Marshal(v)
		request.InstanceIdList = string(idsStr)
	}

	// describe ddosbgp instance filtered by name_regex
	for {
		raw, err := client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
			return ddosbgpClient.DescribeInstanceList(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ddosbgp_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		resp, _ := raw.(*ddosbgp.DescribeInstanceListResponse)
		if len(resp.InstanceList) < 1 {
			break
		}

		for _, item := range resp.InstanceList {
			if nameRegex != nil && !nameRegex.MatchString(item.Remark) {
				continue
			}

			instances = append(instances, item)
		}

		if len(resp.InstanceList) < PageSizeLarge {
			break
		}

		currentPageNo, err := strconv.Atoi(string(request.PageNo))
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ddosbgp_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if page, err := getNextpageNumber(requests.NewInteger(currentPageNo)); err != nil {
			return WrapError(err)
		} else {
			request.PageNo = requests.Integer(page)
		}
	}

	// describe instance spec filtered by instanceids
	nameMap := make(map[string]string)
	instanceIds := make([]string, 0)
	ipTypeMap := make(map[string]string)
	instanceTypeMap := make(map[string]string)
	for _, instance := range instances {
		instanceIds = append(instanceIds, instance.InstanceId)
		nameMap[instance.InstanceId] = instance.Remark
		ipTypeMap[instance.InstanceId] = instance.IpType
		instanceTypeMap[instance.InstanceId] = instance.InstanceType
	}

	if len(instanceIds) < 1 {
		return WrapError(extractDdosbgpInstance(d, nameMap, ipTypeMap, instanceTypeMap, []interface{}{}))
	}

	action := "DescribeInstanceSpecs"
	instanceIdsStr, _ := json.Marshal(instanceIds)

	describeInstanceSpecsRequest := map[string]interface{}{
		"InstanceIdList": string(instanceIdsStr),
		"RegionId":       client.RegionId,
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(6*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ddosbgp", "2018-07-20", action, nil, describeInstanceSpecsRequest, true)
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
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ddosbgp_instances", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.InstanceSpecs", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceSpecs", response)
	}
	result, _ := resp.([]interface{})

	return WrapError(extractDdosbgpInstance(d, nameMap, ipTypeMap, instanceTypeMap, result))
}

func extractDdosbgpInstance(d *schema.ResourceData, nameMap map[string]string, ipTypeMap map[string]string, instanceTypeMap map[string]string, instanceSpecs []interface{}) error {
	var instanceIds []string
	var names []string
	var s []map[string]interface{}

	for _, v := range instanceSpecs {

		item := v.(map[string]interface{})
		ddosbgpInstanceType := string(Enterprise)
		if instanceTypeMap[item["InstanceId"].(string)] == "0" {
			ddosbgpInstanceType = string(Professional)
		}

		mapping := map[string]interface{}{
			"id":               fmt.Sprint(item["InstanceId"]),
			"name":             nameMap[item["InstanceId"].(string)],
			"region":           item["Region"],
			"bandwidth":        item["PackConfig"].(map[string]interface{})["Bandwidth"],
			"base_bandwidth":   item["PackConfig"].(map[string]interface{})["PackBasicThre"],
			"normal_bandwidth": item["PackConfig"].(map[string]interface{})["NormalBandwidth"],
			"ip_type":          ipTypeMap[item["InstanceId"].(string)],
			"ip_count":         item["PackConfig"].(map[string]interface{})["IpSpec"],
			"type":             ddosbgpInstanceType,
		}
		instanceIds = append(instanceIds, fmt.Sprint(item["InstanceId"]))
		names = append(names, nameMap[fmt.Sprint(item["InstanceId"])])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(instanceIds))
	if err := d.Set("ids", instanceIds); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

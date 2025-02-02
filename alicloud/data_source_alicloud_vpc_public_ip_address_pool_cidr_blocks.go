package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpcPublicIpAddressPoolCidrBlocks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcPublicIpAddressPoolCidrBlocksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_ip_address_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Created", "Modifying", "Deleting"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"blocks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_address_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used_ip_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_ip_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpcPublicIpAddressPoolCidrBlocksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPublicIpAddressPoolCidrBlocks"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PublicIpAddressPoolId"] = d.Get("public_ip_address_pool_id")
	request["MaxResults"] = PageSizeLarge

	if v, ok := d.GetOk("cidr_block"); ok {
		request["CidrBlock"] = v
	}

	status, statusOk := d.GetOk("status")

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
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_public_ip_address_pool_cidr_blocks", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.PublicIpPoolCidrBlockList", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PublicIpPoolCidrBlockList", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v", item["PublicIpAddressPoolId"], item["CidrBlock"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                        fmt.Sprintf("%v:%v", object["PublicIpAddressPoolId"], object["CidrBlock"]),
			"public_ip_address_pool_id": fmt.Sprint(object["PublicIpAddressPoolId"]),
			"cidr_block":                fmt.Sprint(object["CidrBlock"]),
			"status":                    object["Status"],
			"used_ip_num":               formatInt(object["UsedIpNum"]),
			"total_ip_num":              formatInt(object["TotalIpNum"]),
			"create_time":               object["CreationTime"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("blocks", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

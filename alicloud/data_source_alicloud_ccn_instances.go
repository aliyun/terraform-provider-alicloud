package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudCcnInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCcnInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
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
						"ccn_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCcnInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateDescribeCloudConnectNetworksRequest()
	request.RegionId = client.RegionId
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	raw, err := client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
		return ccnClient.DescribeCloudConnectNetworks(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ccn_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartag.DescribeCloudConnectNetworksResponse)
	var filteredCcnInstancesTemp []smartag.CloudConnectNetwork
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, ccn := range response.CloudConnectNetworks.CloudConnectNetwork {
			if r != nil && !r.MatchString(ccn.Name) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[ccn.CcnId]; !ok {
					continue
				}
			}

			filteredCcnInstancesTemp = append(filteredCcnInstancesTemp, ccn)
		}
	} else {
		filteredCcnInstancesTemp = response.CloudConnectNetworks.CloudConnectNetwork
	}

	return ccnInstancesDescriptionAttributes(d, filteredCcnInstancesTemp, client)
}

func ccnInstancesDescriptionAttributes(d *schema.ResourceData, ccns []smartag.CloudConnectNetwork, client *connectivity.AliyunClient) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	request := smartag.CreateDescribeCloudConnectNetworksRequest()
	for _, item := range ccns {
		request.CcnId = item.CcnId
		raw, err := client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
			return ccnClient.DescribeCloudConnectNetworks(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ccn_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*smartag.DescribeCloudConnectNetworksResponse)
		mapping := map[string]interface{}{
			"ccn_id":       response.CloudConnectNetworks.CloudConnectNetwork[0].CcnId,
			"name":         response.CloudConnectNetworks.CloudConnectNetwork[0].Name,
			"description":  response.CloudConnectNetworks.CloudConnectNetwork[0].Description,
			"cidr_block":   response.CloudConnectNetworks.CloudConnectNetwork[0].CidrBlock,
			"is_default":   response.CloudConnectNetworks.CloudConnectNetwork[0].IsDefault,
		}

		ids = append(ids, response.CloudConnectNetworks.CloudConnectNetwork[0].CcnId)
		names = append(names, response.CloudConnectNetworks.CloudConnectNetwork[0].Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
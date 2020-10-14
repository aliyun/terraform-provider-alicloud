package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudAccessRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAccessRulesRead,

		Schema: map[string]*schema.Schema{
			"source_cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_group_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user_access": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rw_access": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"access_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_access": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rw_access": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAccessRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateDescribeAccessRulesRequest()
	request.AccessGroupName = d.Get("access_group_name").(string)
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allArs []nas.AccessRule
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			rsp, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
				return nasClient.DescribeAccessRules(request)
			})
			raw = rsp
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_access_rules", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*nas.DescribeAccessRulesResponse)
		if len(response.AccessRules.AccessRule) < 1 {
			break
		}
		for _, rule := range response.AccessRules.AccessRule {
			if v, ok := d.GetOk("source_cidr_ip"); ok && rule.SourceCidrIp != Trim(v.(string)) {
				continue
			}
			if v, ok := d.GetOk("user_access"); ok && rule.UserAccess != Trim(v.(string)) {
				continue
			}
			if v, ok := d.GetOk("rw_access"); ok && rule.RWAccess != Trim(v.(string)) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[rule.AccessRuleId]; !ok {
					continue
				}
			}
			allArs = append(allArs, rule)
		}

		if len(response.AccessRules.AccessRule) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return accessRulesDecriptionAttributes(d, allArs, meta)
}

func accessRulesDecriptionAttributes(d *schema.ResourceData, nasSetTypes []nas.AccessRule, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, ag := range nasSetTypes {
		mapping := map[string]interface{}{
			"source_cidr_ip": ag.SourceCidrIp,
			"priority":       ag.Priority,
			"access_rule_id": ag.AccessRuleId,
			"user_access":    ag.UserAccess,
			"rw_access":      ag.RWAccess,
		}
		ids = append(ids, ag.AccessRuleId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

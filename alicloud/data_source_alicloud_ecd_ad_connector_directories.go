package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcdAdConnectorDirectories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdAdConnectorDirectoriesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"REGISTERING", "REGISTERED", "DEREGISTERING", "NEEDCONFIGTRUST", "CONFIGTRUSTFAILED", "DEREGISTERED", "ERROR", "CONFIGTRUSTING", "NEEDCONFIGUSER"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ad_connectors": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ad_connector_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"connector_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_interface_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"trust_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"specification": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ad_connector_directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dns_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_admin_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"mfa_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_dns_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"sub_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trust_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcdAdConnectorDirectoriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDirectories"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("status"); ok {
		request["DirectoryStatus"] = v
	}
	request["DirectoryType"] = "AD_CONNECTOR"
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var adDirectoryNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		adDirectoryNameRegex = r
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

	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_ad_connector_directories", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Directories", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Directories", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if adDirectoryNameRegex != nil && !adDirectoryNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DirectoryId"])]; !ok {
					continue
				}
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                        fmt.Sprint(object["DirectoryId"]),
			"ad_connector_directory_id": fmt.Sprint(object["DirectoryId"]),
			"create_time":               object["CreationTime"],
			"directory_name":            object["Name"],
			"directory_type":            object["DirectoryType"],
			"dns_address":               object["DnsAddress"],
			"dns_user_name":             object["DnsUserName"],
			"domain_name":               object["DomainName"],
			"domain_user_name":          object["DomainUserName"],
			"status":                    object["Status"],
			"sub_dns_address":           object["SubDnsAddress"],
			"trust_password":            object["TrustPassword"],
			"vswitch_ids":               object["VSwitchIds"],
			"vpc_id":                    object["VpcId"],
			"sub_domain_name":           object["SubDomainName"],
			"enable_admin_access":       object["EnableAdminAccess"],
			"mfa_enabled":               object["MfaEnabled"],
			"custom_security_group_id":  object["CustomSecurityGroupId"],
		}
		aDConnectors := make([]map[string]interface{}, 0)
		if aDConnectorsList, ok := object["ADConnectors"].([]interface{}); ok {
			for _, v := range aDConnectorsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"ad_connector_address": m1["ADConnectorAddress"],
						"connector_status":     m1["ConnectorStatus"],
						"network_interface_id": m1["NetworkInterfaceId"],
						"vswitch_id":           m1["VSwitchId"],
						"specification":        m1["Specification"],
						"trust_key":            m1["TrustKey"],
					}
					aDConnectors = append(aDConnectors, temp1)
				}
			}
		}
		mapping["ad_connectors"] = aDConnectors
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("directories", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

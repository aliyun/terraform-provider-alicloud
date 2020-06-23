package alicloud

import (
	"regexp"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	dms_enterprise "github.com/aliyun/alibaba-cloud-sdk-go/services/dms-enterprise"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudDmsEnterpriseInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDmsEnterpriseInstancesRead,
		Schema: map[string]*schema.Schema{
			"env_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_source": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"net_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"search_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DELETED", "DISABLE", "NORMAL", "UNAVAILABLE"}, false),
			},
			"tid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"instance_alias_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_link_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dba_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dba_nick_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ddl_online": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ecs_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"env_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"export_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"query_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"safe_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"use_dsql": {
							Type:     schema.TypeInt,
							Computed: true,
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

func dataSourceAlicloudDmsEnterpriseInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := dms_enterprise.CreateListInstancesRequest()
	if v, ok := d.GetOk("env_type"); ok {
		request.EnvType = v.(string)
	}
	if v, ok := d.GetOk("instance_source"); ok {
		request.InstanceSource = v.(string)
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request.DbType = v.(string)
	}
	if v, ok := d.GetOk("net_type"); ok {
		request.NetType = v.(string)
	}
	if v, ok := d.GetOk("search_key"); ok {
		request.SearchKey = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		request.InstanceState = v.(string)
	}
	if v, ok := d.GetOk("tid"); ok {
		request.Tid = requests.NewInteger(v.(int))
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []dms_enterprise.Instance
	var instance_aliasRegex *regexp.Regexp
	if v, ok := d.GetOk("instance_alias_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		instance_aliasRegex = r
	}
	for {
		raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
			return dms_enterpriseClient.ListInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*dms_enterprise.ListInstancesResponse)

		for _, item := range response.InstanceList.Instance {
			if instance_aliasRegex != nil {
				if !instance_aliasRegex.MatchString(item.InstanceAlias) {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.InstanceList.Instance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"data_link_name":    object.DataLinkName,
			"database_password": object.DatabasePassword,
			"database_user":     object.DatabaseUser,
			"dba_id":            object.DbaId,
			"dba_nick_name":     object.DbaNickName,
			"ddl_online":        object.DdlOnline,
			"ecs_instance_id":   object.EcsInstanceId,
			"ecs_region":        object.EcsRegion,
			"env_type":          object.EnvType,
			"export_timeout":    object.ExportTimeout,
			"host":              object.Host,
			"instance_alias":    object.InstanceAlias,
			"instance_id":       object.InstanceId,
			"instance_source":   object.InstanceSource,
			"instance_type":     object.InstanceType,
			"port":              object.Port,
			"query_timeout":     object.QueryTimeout,
			"safe_rule_id":      object.SafeRuleId,
			"sid":               object.Sid,
			"status":            object.State,
			"use_dsql":          object.UseDsql,
			"vpc_id":            object.VpcId,
		}
		s[i] = mapping
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

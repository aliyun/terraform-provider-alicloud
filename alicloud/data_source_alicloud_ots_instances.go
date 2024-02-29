package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
)

func dataSourceAlicloudOtsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOtsInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"tags": tagsSchema(),
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
							Removed:  "Field 'network' has been removed from provider version v1.221.0. Please Use the 'network_type_acl' and 'network_source_acl'",
						},
						"network_type_acl": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"network_source_acl": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_quota": {
							Type:     schema.TypeInt,
							Computed: true,
							Removed:  "Field 'entity_quota' has been removed from provider version v1.221.0. Please Use the 'table_quota'",
						},
						"table_quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOtsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	allInstanceNames, err := otsService.ListOtsInstance(PageSizeLarge)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_instances", "ListOtsInstance", AlibabaCloudSdkGoERROR)
	}

	idsMap := make(map[string]bool)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, x := range v.([]interface{}) {
			if x == nil {
				continue
			}
			idsMap[x.(string)] = true
		}
	}

	var nameReg *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		nameReg = regexp.MustCompile(v.(string))
	}

	var filteredInstanceNames []string
	for _, instanceName := range allInstanceNames {
		// name_regex mismatch
		if nameReg != nil && !nameReg.MatchString(instanceName) {
			continue
		}
		// ids mismatch
		if len(idsMap) != 0 {
			if _, ok := idsMap[instanceName]; !ok {
				continue
			}
		}
		filteredInstanceNames = append(filteredInstanceNames, instanceName)
	}

	// get full instance info via GetInstance
	var allInstances []RestOtsInstanceInfo
	for _, instanceName := range filteredInstanceNames {
		instanceInfo, err := otsService.DescribeOtsInstance(instanceName)
		if err != nil {
			return WrapError(err)
		}
		allInstances = append(allInstances, instanceInfo)
	}

	// filter by tag.
	var filteredInstances []RestOtsInstanceInfo
	if v, ok := d.GetOk("tags"); ok {
		if vmap, ok := v.(map[string]interface{}); ok && len(vmap) > 0 {
			for _, instance := range allInstances {
				if tagsMapEqual(vmap, otsRestTagsToMap(instance.Tags)) {
					filteredInstances = append(filteredInstances, instance)
				}
			}
		} else {
			filteredInstances = allInstances[:]
		}
	} else {
		filteredInstances = allInstances[:]
	}
	return otsInstancesDecriptionAttributes(d, filteredInstances, meta)
}

func otsInstancesDecriptionAttributes(d *schema.ResourceData, instances []RestOtsInstanceInfo, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, instance := range instances {
		mapping := map[string]interface{}{
			"id":                 instance.InstanceName,
			"name":               instance.InstanceName,
			"status":             toInstanceOuterStatus(instance.InstanceStatus),
			"cluster_type":       instance.InstanceSpecification,
			"create_time":        instance.CreateTime,
			"user_id":            instance.UserId,
			"resource_group_id":  instance.ResourceGroupId,
			"network_type_acl":   instance.NetworkTypeACL,
			"network_source_acl": instance.NetworkSourceACL,
			"policy":             instance.Policy,
			"policy_version":     instance.PolicyVersion,
			"description":        instance.InstanceDescription,
			"table_quota":        instance.Quota.TableQuota,
			"tags":               otsRestTagsToMap(instance.Tags),
		}
		names = append(names, instance.InstanceName)
		ids = append(ids, instance.InstanceName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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

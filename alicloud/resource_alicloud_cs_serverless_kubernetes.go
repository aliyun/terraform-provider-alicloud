package alicloud

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSServerlessKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSServerlessKubernetesCreate,
		Read:   resourceAlicloudCSServerlessKubernetesRead,
		Update: resourceAlicloudCSServerlessKubernetesUpdate,
		Delete: resourceAlicloudCSServerlessKubernetesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  StringLenBetween(1, 63),
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  StringLenBetween(0, 37),
				ConflictsWith: []string{"name"},
			},
			"vpc_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"zone_id"},
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Field 'vswitch_id' has been removed from provider version 1.229.1. New field 'vswitch_ids' replace it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:      1,
				ConflictsWith: []string{"zone_id"},
			},
			"service_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_rrsa": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"private_zone": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"service_discovery_types"},
				Deprecated:    "Field 'private_zone' has been deprecated from provider version 1.123.1. New field 'service_discovery_types' replace it.",
			},
			"service_discovery_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringInSlice([]string{"CoreDNS", "PrivateZone"}, false),
				},
				ConflictsWith: []string{"private_zone"},
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_id", "vswitch_ids", "vswitch_id"},
			},
			"endpoint_public_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"kube_config": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'kube_config' has been deprecated from provider version 1.187.0. New DataSource 'alicloud_cs_cluster_credential' manage your cluster's kube config.",
			},
			"client_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_ca_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "Field 'force_update' has been removed from provider version 1.229.1.",
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"load_balancer_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
				Deprecated:   "Field 'load_balancer_spec' has been deprecated from provider version 1.229.1. The load balancer has been changed to PayByCLCU so that the spec is no need anymore.",
			},
			"logging_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Default:    "SLS",
				Deprecated: "Field 'logging_type' has been deprecated from provider version 1.229.1. Please use addons `alibaba-log-controller` to enable logging.",
			},
			"sls_project_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'sls_project_name' has been deprecated from provider version 1.229.1. Please use the field `config` of addons `alibaba-log-controller` to specify log project name.",
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"retain_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ack.standard", "ack.pro.small"}, false),
			},
			"create_v2_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Removed:  "Field 'create_v2_cluster' has been removed from provider version 1.229.1.",
			},
			"rrsa_metadata": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rrsa_oidc_issuer_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_oidc_provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_oidc_provider_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"custom_san": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_options": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"SLB", "ALB", "SLS_Data", "SLS_ControlPlane", "PrivateZone"}, false),
						},
						"delete_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"delete", "retain"}, false),
						},
					},
				},
			},
			"maintenance_window": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"maintenance_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"weekly_period": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"operation_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_auto_upgrade": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"channel": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCSServerlessKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	roa, _ := client.NewRoaCsClient()
	csClient := CsClient{roa}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	tags := make([]*roacs.Tag, 0)
	tagsMap, ok := d.Get("tags").(map[string]interface{})
	if ok {
		for key, value := range tagsMap {
			if value != nil {
				if v, ok := value.(string); ok {
					tags = append(tags, &roacs.Tag{
						Key:   tea.String(key),
						Value: tea.String(v),
					})
				}
			}
		}
	}

	addons := make([]*roacs.Addon, 0)
	if v, ok := d.GetOk("addons"); ok {
		all, ok := v.([]interface{})
		if ok {
			for _, a := range all {
				addon, ok := a.(map[string]interface{})
				if ok {
					addons = append(addons, &roacs.Addon{
						Name:     tea.String(addon["name"].(string)),
						Config:   tea.String(addon["config"].(string)),
						Version:  tea.String(addon["version"].(string)),
						Disabled: tea.Bool(addon["disabled"].(bool)),
					})
				}
			}
		}
	}
	if v, ok := d.GetOk("sls_project_name"); ok {
		exist := false
		for _, addon := range addons {
			if tea.StringValue(addon.Name) == "alibaba-log-controller" {
				var config map[string]interface{}
				err := json.Unmarshal([]byte(tea.StringValue(addon.Config)), &config)
				if err == nil {
					if project, ok := config["sls_project_name"].(string); !ok || project == "" {
						config["sls_project_name"] = v.(string)
					}
					if configRaw, err := json.Marshal(config); err == nil {
						addon.Config = tea.String(string(configRaw))
					}
				}
				exist = true
				break
			}
		}
		if exist == false {
			config := map[string]interface{}{
				"sls_project_name": v.(string),
			}
			if configRaw, err := json.Marshal(config); err == nil {
				addons = append(addons, &roacs.Addon{
					Name:   tea.String("alibaba-log-controller"),
					Config: tea.String(string(configRaw)),
				})
			}

		}
	}

	request := &roacs.CreateClusterRequest{
		Name:                 tea.String(clusterName),
		KubernetesVersion:    tea.String(d.Get("version").(string)),
		ClusterType:          tea.String("ManagedKubernetes"),
		Profile:              tea.String("Serverless"),
		Tags:                 tags,
		Addons:               addons,
		DeletionProtection:   tea.Bool(d.Get("deletion_protection").(bool)),
		RegionId:             tea.String(client.RegionId),
		ResourceGroupId:      tea.String(d.Get("resource_group_id").(string)),
		Vpcid:                tea.String(d.Get("vpc_id").(string)),
		SnatEntry:            tea.Bool(d.Get("new_nat_gateway").(bool)),
		NatGateway:           tea.Bool(d.Get("new_nat_gateway").(bool)),
		EndpointPublicAccess: tea.Bool(d.Get("endpoint_public_access_enabled").(bool)),
		SecurityGroupId:      tea.String(d.Get("security_group_id").(string)),
		LoggingType:          tea.String(d.Get("logging_type").(string)),
	}
	if v, ok := d.GetOk("time_zone"); ok {
		request.Timezone = tea.String(v.(string))
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = tea.String(v.(string))
	}

	if v, ok := d.GetOk("service_cidr"); ok {
		request.ServiceCidr = tea.String(v.(string))
	}

	if v, ok := d.GetOk("service_discovery_types"); ok {
		r := make([]*string, 0)
		for _, vv := range v.([]interface{}) {
			r = append(r, tea.String(vv.(string)))
		}
		request.ServiceDiscoveryTypes = r
	}

	if v, ok := d.GetOkExists("private_zone"); ok {
		request.ServiceDiscoveryTypes = []*string{}
		if v.(bool) == true {
			request.ServiceDiscoveryTypes = []*string{tea.String("PrivateZone")}
		}
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		r := make([]*string, 0)
		for _, vv := range v.([]interface{}) {
			r = append(r, tea.String(vv.(string)))
		}
		request.VswitchIds = r
	}

	if v, ok := d.GetOk("load_balancer_spec"); ok {
		request.LoadBalancerSpec = tea.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_spec"); ok {
		request.ClusterSpec = tea.String(v.(string))
	}

	if v, ok := d.GetOk("enable_rrsa"); ok {
		request.EnableRrsa = tea.Bool(v.(bool))
	}

	if v, ok := d.GetOk("custom_san"); ok {
		request.CustomSan = tea.String(v.(string))
	}

	if v, ok := d.GetOk("maintenance_window"); ok {
		request.MaintenanceWindow = expandMaintenanceWindowConfigRoa(v.([]interface{}))
	}

	if v, ok := d.GetOk("operation_policy"); ok {
		request.OperationPolicy = &roacs.CreateClusterRequestOperationPolicy{}
		m := v.([]interface{})[0].(map[string]interface{})
		if vv, ok := m["cluster_auto_upgrade"]; ok {
			policy := vv.([]interface{})[0].(map[string]interface{})
			request.OperationPolicy.ClusterAutoUpgrade = &roacs.CreateClusterRequestOperationPolicyClusterAutoUpgrade{
				Enabled: tea.Bool(policy["enabled"].(bool)),
				Channel: tea.String(policy["channel"].(string)),
			}
		}
	}

	//set tags
	if len(tags) > 0 {
		request.Tags = tags
	}
	var err error
	var resp *roacs.CreateClusterResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = csClient.client.CreateCluster(request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_serverless_kubernetes", "CreateServerlessKubernetesCluster", AlibabaCloudSdkGoERROR)
	}
	d.SetId(tea.StringValue(resp.Body.ClusterId))
	taskId := tea.StringValue(resp.Body.TaskId)
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), "createCluster", jobDetail)
	}

	csService := CsService{client}
	stateConf = BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, csService.CsServerlessKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCSServerlessKubernetesRead(d, meta)
}

func resourceAlicloudCSServerlessKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosClient, err := client.NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}
	csClient := CsClient{rosClient}

	object, err := csClient.DescribeClusterDetail(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	vswitchIds := []string{}
	request := &roacs.DescribeClusterResourcesRequest{}
	resources, _ := rosClient.DescribeClusterResources(tea.String(d.Id()), request)
	for _, resource := range resources.Body {
		if tea.StringValue(resource.ResourceType) == "VSWITCH" {
			vswitchIds = append(vswitchIds, tea.StringValue(resource.InstanceId))
		}
	}

	d.Set("name", object.Name)
	d.Set("vpc_id", object.VpcId)
	d.Set("vswitch_ids", vswitchIds)
	d.Set("security_group_id", object.SecurityGroupId)
	d.Set("deletion_protection", object.DeletionProtection)
	d.Set("version", object.CurrentVersion)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("cluster_spec", object.ClusterSpec)

	if object.Timezone != nil {
		d.Set("timezone", object.Timezone)
	}

	if err := d.Set("tags", flattenTags(object.Tags)); err != nil {
		return WrapError(err)
	}
	if d.Get("load_balancer_spec") == "" {
		d.Set("load_balancer_spec", "slb.s2.small")
	}
	if d.Get("logging_type") == "" {
		d.Set("logging_type", "SLS")
	}

	if object.ServiceCidr != nil {
		d.Set("service_cidr", object.ServiceCidr)
	} else {
		if v, ok := object.Parameters["ServiceCIDR"]; ok {
			d.Set("service_cidr", v)
		}
	}
	capabilities := fetchClusterCapabilities(tea.StringValue(object.MetaData))
	if v, ok := capabilities["PublicSLB"]; ok {
		d.Set("endpoint_public_access_enabled", Interface2Bool(v))
	}
	metadata := fetchClusterMetaDataMap(tea.StringValue(object.MetaData))
	if v, ok := metadata["ExtraCertSAN"]; ok && v != nil {
		l := expandStringList(v.([]interface{}))
		d.Set("custom_san", strings.Join(l, ","))
	}

	if data, err := flattenRRSAMetadata(tea.StringValue(object.MetaData)); err != nil {
		return WrapError(err)
	} else {
		d.Set("rrsa_metadata", data)
		if len(data) > 0 {
			d.Set("enable_rrsa", data[0]["enabled"].(bool))
		}
	}

	if object.MaintenanceWindow != nil {
		d.Set("maintenance_window", flattenMaintenanceWindowConfigRoa(object.MaintenanceWindow))
	}

	if object.OperationPolicy != nil {
		m := make([]map[string]interface{}, 0)
		if object.OperationPolicy.ClusterAutoUpgrade != nil {
			m = append(m, map[string]interface{}{
				"cluster_auto_upgrade": []map[string]interface{}{
					{
						"enabled": tea.BoolValue(object.OperationPolicy.ClusterAutoUpgrade.Enabled),
						"channel": tea.StringValue(object.OperationPolicy.ClusterAutoUpgrade.Channel),
					},
				},
			})
		}
		d.Set("operation_policy", m)
	}

	// get cluster conn certs
	// If the cluster is failed, there is no need to get cluster certs
	if tea.StringValue(object.State) == "failed" || tea.StringValue(object.State) == "delete_failed" || tea.StringValue(object.State) == "deleting" {
		return nil
	}

	if err = setCerts(d, meta, false); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudCSServerlessKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	invoker := NewInvoker()
	// modifyCluster
	if !d.IsNewResource() && d.HasChanges("resource_group_id", "name", "name_prefix", "deletion_protection", "custom_san", "maintenance_window", "operation_policy", "enable_rrsa") {
		if err := modifyCluster(d, meta, &invoker); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyCluster", AlibabaCloudSdkGoERROR)
		}
	}

	// modify cluster tag
	if d.HasChange("tags") {
		err := updateKubernetesClusterTag(d, meta)
		if err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "ModifyClusterTags", AlibabaCloudSdkGoERROR)
		}
	}
	// migrate cluster to pro form standard
	if d.HasChange("cluster_spec") {
		err := migrateCluster(d, meta)
		if err != nil {
			return WrapError(err)
		}
	}

	// upgrade cluster
	err := UpgradeAlicloudKubernetesCluster(d, meta)
	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudCSServerlessKubernetesRead(d, meta)
}

func resourceAlicloudCSServerlessKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAlicloudCSKubernetesDelete(d, meta)
}

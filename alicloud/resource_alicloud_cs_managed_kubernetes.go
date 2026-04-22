package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	roacs "github.com/alibabacloud-go/cs-20151215/v7/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSManagedKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSManagedKubernetesCreate,
		Read:   resourceAlicloudCSManagedKubernetesRead,
		Update: resourceAlicloudCSManagedKubernetesUpdate,
		Delete: resourceAlicloudCSKubernetesDelete, // TODO Refactor delete from k8s resources
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
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
				Default:       "Terraform-Creation",
				ValidateFunc:  StringLenBetween(0, 37),
				ConflictsWith: []string{"name"},
			},
			"profile": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"worker_vswitch_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				Optional:   true,
				Deprecated: "Field 'worker_vswitch_ids' has been deprecated from provider version 1.241.0. Please use 'vswitch_ids' to managed control plane vswtiches",
			},
			// global configurations
			"zone_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 5,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:     1,
				MaxItems:     5,
				ExactlyOneOf: []string{"worker_vswitch_ids", "vswitch_ids", "zone_ids"},
			},
			"pod_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
			},
			"pod_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"service_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_cidr_mask": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc: IntBetween(24, 28),
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"user_ca": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ipvs",
				ValidateFunc: StringInSlice([]string{"iptables", "ipvs"}, false),
			},
			"ip_stack": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ipv4", "dual"}, false),
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
			"slb_internet_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"load_balancer_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
				Computed:     true,
				Deprecated:   "Field 'load_balancer_spec' has been deprecated from provider version 1.232.0. The spec will not take effect because the charge of the load balancer has been changed to PayByCLCU",
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_rrsa": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "cluster.local",
				ForceNew:    true,
				Description: "cluster local domain ",
			},
			"custom_san": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_provider_key": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsEncryptionDiffSuppressFunc,
				Description:      "The ID of the Key Management Service (KMS) key that is used to encrypt Kubernetes Secrets.",
			},
			"disable_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"client_cert": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'client_cert' has been deprecated from provider version 1.248.0. From version 1.248.0, new DataSource 'alicloud_cs_cluster_credential' is recommended to manage cluster's kubeconfig, you can also save the 'certificate_authority.client_cert' attribute content of new DataSource 'alicloud_cs_cluster_credential' to an appropriate path(like ~/.kube/client-cert.pem) for replace it.",
			},
			"client_key": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'client_key' has been deprecated from provider version 1.248.0. From version 1.248.0, new DataSource 'alicloud_cs_cluster_credential' is recommended to manage cluster's kubeconfig, you can also save the 'certificate_authority.client_key' attribute content of new DataSource 'alicloud_cs_cluster_credential' to an appropriate path(like ~/.kube/client-key.pem) for replace it.",
			},
			"cluster_ca_cert": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'cluster_ca_cert' has been deprecated from provider version 1.248.0. From version 1.248.0, new DataSource 'alicloud_cs_cluster_credential' is recommended to manage cluster's kubeconfig, you can also save the 'certificate_authority.cluster_cert' attribute content of new DataSource 'alicloud_cs_cluster_credential' to an appropriate path(like ~/.kube/cluster-ca-cert.pem) for replace it.",
			},
			"certificate_authority": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Deprecated: "Field 'certificate_authority' has been deprecated from provider version 1.248.0. Please use the attribute 'certificate_authority' of new DataSource 'alicloud_cs_cluster_credential' to replace it.",
			},
			"skip_set_certificate_authority": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_server_internet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_server_intranet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"slb_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slb_internet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slb_intranet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_enterprise_security_group": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_group_id"},
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"worker_ram_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account_issuer": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"api_audiences": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			// lintignore: S006
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ack.standard", "ack.pro.small"}, false),
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
			"control_plane_log_ttl": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"control_plane_log_project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"control_plane_log_components": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"audit_log_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"sls_project_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"auto_mode": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},
			"retain_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"rrsa_metadata": {
				Type:     schema.TypeList,
				Computed: true,
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
			"upgrade_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"control_plane_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCSManagedKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
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

	vpcService := VpcService{client}
	var vSwitchIds []string
	if v, ok := d.GetOk("vswitch_ids"); ok {
		vSwitchIds = expandStringList(v.([]interface{}))
	} else {
		if v, ok := d.GetOk("worker_vswitch_ids"); ok {
			vSwitchIds = expandStringList(v.([]interface{}))
		}
	}
	var vpcId string
	if len(vSwitchIds) > 0 {
		vsw, err := vpcService.DescribeVSwitch(vSwitchIds[0])
		if err != nil {
			return err
		}
		vpcId = vsw.VpcId
	}

	request := &roacs.CreateClusterRequest{
		Name:        tea.String(clusterName),
		RegionId:    tea.String(client.RegionId),
		ClusterType: tea.String("ManagedKubernetes"),
		Profile:     tea.String("Default"),
		Tags:        tags,
		Addons:      addons,
		Vpcid:       tea.String(vpcId),
		VswitchIds:  tea.StringSlice(vSwitchIds),
	}
	if v, ok := d.GetOk("profile"); ok {
		request.SetProfile(v.(string))
	}

	if v, ok := d.GetOk("version"); ok {
		request.SetKubernetesVersion(v.(string))
	}

	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request.SetDeletionProtection(v.(bool))
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.SetResourceGroupId(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.SetVpcid(v.(string))
	}

	if v, ok := d.GetOk("zone_ids"); ok {
		request.SetZoneIds(tea.StringSlice(expandStringList(v.([]interface{}))))
	}

	if v, ok := d.GetOk("new_nat_gateway"); ok {
		request.SetSnatEntry(v.(bool))
	}

	if v, ok := d.GetOk("slb_internet_enabled"); ok {
		request.SetEndpointPublicAccess(v.(bool))
	}

	if v, ok := d.GetOk("load_balancer_spec"); ok {
		request.SetLoadBalancerSpec(v.(string))
	}

	if v, ok := d.GetOk("is_enterprise_security_group"); ok {
		request.SetIsEnterpriseSecurityGroup(v.(bool))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request.SetSecurityGroupId(v.(string))
	}

	if v, ok := d.GetOk("service_cidr"); ok {
		request.SetServiceCidr(v.(string))
	}

	if v, ok := d.GetOk("proxy_mode"); ok {
		request.SetProxyMode(v.(string))
	}

	if v, ok := d.GetOk("ip_stack"); ok {
		request.SetIpStack(v.(string))
	}

	if v, ok := d.GetOk("timezone"); ok {
		request.SetTimezone(v.(string))
	}

	if v, ok := d.GetOk("pod_vswitch_ids"); ok {
		request.SetPodVswitchIds(tea.StringSlice(expandStringList(v.([]interface{}))))
	}

	if v, ok := d.GetOk("pod_cidr"); ok {
		request.SetContainerCidr(v.(string))
	}

	if v, ok := d.GetOk("node_cidr_mask"); ok {
		request.SetNodeCidrMask(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("cluster_spec"); ok {
		request.SetClusterSpec(v.(string))
	}

	if v, ok := d.GetOk("cluster_domain"); ok {
		request.SetClusterDomain(v.(string))
	}

	if v, ok := d.GetOk("service_account_issuer"); ok {
		request.SetServiceAccountIssuer(v.(string))
	}

	if v, ok := d.GetOk("api_audiences"); ok {
		if list := expandStringList(v.([]interface{})); len(list) > 0 {
			request.SetApiAudiences(strings.Join(list, ","))
		}
	}

	if v, ok := d.GetOk("enable_rrsa"); ok {
		request.SetEnableRrsa(v.(bool))
	}
	if v, ok := d.GetOk("custom_san"); ok {
		request.SetCustomSan(v.(string))
	}

	if v, ok := d.GetOk("encryption_provider_key"); ok {
		request.SetEncryptionProviderKey(v.(string))
	}

	// Configure control plane log. Effective only in the professional managed cluster
	if v, ok := d.GetOk("control_plane_log_components"); ok {
		request.SetControlplaneLogComponents(tea.StringSlice(expandStringList(v.([]interface{}))))
		// ttl default is 30 days
		request.SetControlplaneLogTtl("30")
	}
	if v, ok := d.GetOk("control_plane_log_ttl"); ok {
		request.SetControlplaneLogTtl(v.(string))
	}
	if v, ok := d.GetOk("control_plane_log_project"); ok {
		request.SetControlplaneLogProject(v.(string))
	}

	if v, ok := d.GetOk("maintenance_window"); ok {
		request.SetMaintenanceWindow(expandMaintenanceWindowConfigRoa(v.([]interface{})))
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
	if v, ok := d.GetOk("audit_log_config"); ok {
		m := v.([]interface{})[0].(map[string]interface{})
		if vv, ok := m["enabled"]; ok {
			request.AuditLogConfig = &roacs.CreateClusterRequestAuditLogConfig{
				Enabled: tea.Bool(vv.(bool)),
			}
		}
		if vv, ok := m["sls_project_name"]; ok {
			request.AuditLogConfig.SlsProjectName = tea.String(vv.(string))
		}
	}

	if v, ok := d.GetOk("auto_mode"); ok {
		m := v.([]interface{})[0].(map[string]interface{})
		if vv, ok := m["enabled"]; ok {
			request.AutoMode = &roacs.CreateClusterRequestAutoMode{
				Enable: tea.Bool(vv.(bool)),
			}
		}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_managed_kubernetes", "CreateManagedKubernetesCluster", AlibabaCloudSdkGoERROR)
	}
	d.SetId(tea.StringValue(resp.Body.ClusterId))
	taskId := tea.StringValue(resp.Body.TaskId)
	roaCsClient, err := client.NewRoaCsClient()
	if err == nil {
		csClient := CsClient{client: roaCsClient}
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "createCluster", jobDetail)
		}
	}

	csService := CsService{client}
	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCSManagedKubernetesRead(d, meta)
}

func resourceAlicloudCSManagedKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosClient, err := client.NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}
	csClient := CsClient{rosClient}

	object, err := csClient.DescribeClusterDetail(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cs_managed_kubernetes DescribeClusterDetail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	// compat for default value
	if spec := d.Get("load_balancer_spec").(string); spec != "" {
		d.Set("load_balancer_spec", spec)
	}

	if object.Name != nil {
		d.Set("name", object.Name)
	}

	if object.Profile != nil {
		d.Set("profile", object.Profile)
	}

	if object.VpcId != nil {
		d.Set("vpc_id", object.VpcId)
	}

	if object.VswitchIds != nil {
		d.Set("vswitch_ids", tea.StringSliceValue(object.VswitchIds))
	}

	// compat for old value
	if v := d.Get("worker_vswitch_ids"); v != nil {
		d.Set("worker_vswitch_ids", v)
	}

	if object.SecurityGroupId != nil {
		d.Set("security_group_id", object.SecurityGroupId)
	}

	if object.DeletionProtection != nil {
		d.Set("deletion_protection", object.DeletionProtection)
	}

	if object.CurrentVersion != nil {
		d.Set("version", object.CurrentVersion)
	}

	if object.ResourceGroupId != nil {
		d.Set("resource_group_id", object.ResourceGroupId)
	}

	if object.ClusterSpec != nil {
		d.Set("cluster_spec", object.ClusterSpec)
	}

	if object.Timezone != nil {
		d.Set("timezone", object.Timezone)
	}

	if object.WorkerRamRoleName != nil {
		d.Set("worker_ram_role_name", object.WorkerRamRoleName)
	}

	d.Set("cluster_domain", "cluster.local")
	if object.ClusterDomain != nil {
		d.Set("cluster_domain", object.ClusterDomain)
	}

	if err := d.Set("tags", flattenTags(object.Tags)); err != nil {
		return WrapError(err)
	}

	slbId, err := getApiServerSlbID(d, meta)
	if err != nil {
		log.Printf(DefaultErrorMsg, d.Id(), "DescribeClusterResources", err.Error())
	}
	d.Set("slb_id", slbId)

	if object.ServiceCidr != nil {
		d.Set("service_cidr", object.ServiceCidr)
	} else {
		if v, ok := object.Parameters["ServiceCIDR"]; ok {
			d.Set("service_cidr", v)
		}
	}
	if object.ProxyMode != nil {
		d.Set("proxy_mode", object.ProxyMode)
	} else {
		if v, ok := object.Parameters["ProxyMode"]; ok {
			d.Set("proxy_mode", v)
		}
	}

	if object.IpStack != nil {
		d.Set("ip_stack", object.IpStack)
	}

	if object.ContainerCidr != nil {
		d.Set("pod_cidr", object.ContainerCidr)
	} else {
		if v, ok := object.Parameters["ContainerCIDR"]; ok {
			d.Set("pod_cidr", v)
		}
	}

	capabilities := fetchClusterCapabilities(tea.StringValue(object.MetaData))
	if v, ok := capabilities["DisableEncryption"]; ok {
		d.Set("disable_encryption", Interface2Bool(v))
	}
	if kmsKeyId, ok := capabilities["EncryptionKMSKeyId"]; ok {
		d.Set("encryption_provider_key", Interface2String(kmsKeyId))
	}

	if object.NodeCidrMask != nil {
		d.Set("node_cidr_mask", formatInt(tea.StringValue(object.NodeCidrMask)))
	} else {
		// node_cidr_mask
		if v, ok := capabilities["NodeCIDRMask"]; ok {
			d.Set("node_cidr_mask", formatInt(v))
		}
	}

	metadata := fetchClusterMetaDataMap(tea.StringValue(object.MetaData))
	if v, ok := metadata["ExtraCertSAN"]; ok && v != nil {
		l := expandStringList(v.([]interface{}))
		d.Set("custom_san", strings.Join(l, ","))
	}
	// rrsa metadata only for managed, ignore attributes error
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

	if object.AutoMode != nil {
		m := make(map[string]interface{})
		if object.AutoMode.Enable != nil {
			m["enabled"] = tea.BoolValue(object.AutoMode.Enable)
		}
		d.Set("auto_mode", []map[string]interface{}{m})
	}

	// Get slb information and set connect
	connection := make(map[string]string)
	masterURL := tea.StringValue(object.MasterUrl)
	endPoint := make(map[string]string)
	_ = json.Unmarshal([]byte(masterURL), &endPoint)
	connection["api_server_internet"] = endPoint["api_server_endpoint"]
	connection["api_server_intranet"] = endPoint["intranet_api_server_endpoint"]
	if endPoint["api_server_endpoint"] != "" {
		connection["master_public_ip"] = strings.Split(strings.Split(endPoint["api_server_endpoint"], ":")[1], "/")[2]
	}
	connection["service_domain"] = fmt.Sprintf("*.%s.%s.alicontainer.com", d.Id(), tea.StringValue(object.RegionId))

	d.Set("connections", []interface{}{connection})
	d.Set("slb_internet", connection["master_public_ip"])
	if endPoint["intranet_api_server_endpoint"] != "" {
		d.Set("slb_intranet", strings.Split(strings.Split(endPoint["intranet_api_server_endpoint"], ":")[1], "/")[2])
	}

	// set nat gateway
	natRequest := vpc.CreateDescribeNatGatewaysRequest()
	natRequest.VpcId = tea.StringValue(object.VpcId)
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeNatGateways(natRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), natRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(natRequest.GetActionName(), raw, natRequest.RpcRequest, natRequest)
	nat, _ := raw.(*vpc.DescribeNatGatewaysResponse)
	if nat != nil && len(nat.NatGateways.NatGateway) > 0 {
		d.Set("nat_gateway_id", nat.NatGateways.NatGateway[0].NatGatewayId)
	}

	// get cluster conn certs
	// If the cluster is failed, there is no need to get cluster certs
	if tea.StringValue(object.State) == "failed" || tea.StringValue(object.State) == "delete_failed" || tea.StringValue(object.State) == "deleting" {
		return nil
	}

	if err = setCerts(d, meta, d.Get("skip_set_certificate_authority").(bool)); err != nil {
		return WrapError(err)
	}

	if err = checkControlPlaneLogEnable(d, meta); err != nil {
		return WrapError(err)
	}

	if err = getClusterAuditProject(d, meta); err != nil {
		return WrapError(err)
	}

	return nil

}

func resourceAlicloudCSManagedKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	invoker := NewInvoker()
	// modifyCluster
	if !d.IsNewResource() && d.HasChanges("resource_group_id", "name", "name_prefix", "deletion_protection", "maintenance_window", "operation_policy",
		"custom_san", "vswitch_ids", "timezone", "security_group_id", "enable_rrsa") {
		if err := modifyManagedK8sCluster(d, meta, &invoker); err != nil {
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

	// update control plane config
	if d.HasChanges([]string{"control_plane_log_ttl", "control_plane_log_project", "control_plane_log_components"}...) {
		if err := updateControlPlaneLog(d, meta); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateControlPlaneLog", AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("audit_log_config") {
		if err := updateClusterAuditConfig(d, meta); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateClusterAuditConfig", AlibabaCloudSdkGoERROR)
		}
	}

	// migrate cluster to pro from standard
	if d.HasChange("cluster_spec") {
		err := migrateCluster(d, meta)
		if err != nil {
			return WrapError(err)
		}
	}

	// update kms encryption
	if d.HasChange("encryption_provider_key") || d.HasChange("disable_encryption") {
		if err := updateClusterKMSEncryption(d, meta); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateKMSEncryption", AlibabaCloudSdkGoERROR)
		}
	}

	err := UpgradeAlicloudKubernetesCluster(d, meta)
	if err != nil {
		return WrapError(err)
	}

	d.Partial(false)
	return resourceAlicloudCSManagedKubernetesRead(d, meta)
}

func UpgradeAlicloudKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("version") {
		return nil
	}

	var controlPlaneOnly *bool
	if v, ok := d.GetOk("upgrade_policy"); ok {
		val, err := jsonpath.Get("$[0].control_plane_only", v)
		if err != nil {
			return WrapError(err)
		}
		if val != nil && val != "" {
			controlPlaneOnly = tea.Bool(val.(bool))
		}
	}

	clusterId := d.Id()
	version := d.Get("version").(string)
	action := "UpgradeCluster"
	c := meta.(*connectivity.AliyunClient)
	rosCsClient, err := c.NewRoaCsClient()
	if err != nil {
		return err
	}
	args := &roacs.UpgradeClusterRequest{
		NextVersion: tea.String(version),
	}
	if controlPlaneOnly != nil {
		args.MasterOnly = controlPlaneOnly
	}
	// upgrade cluster
	var resp *roacs.UpgradeClusterResponse
	err = resource.Retry(UpgradeClusterTimeout, func() *resource.RetryError {
		resp, err = rosCsClient.UpgradeCluster(tea.String(clusterId), args)
		if NeedRetry(err) || resp == nil {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), action, err)
	}

	taskId := tea.StringValue(resp.Body.TaskId)
	if taskId == "" {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), action, resp)
	}

	csClient := CsClient{client: rosCsClient}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		// try to cancel task
		wait := incrementalWait(3*time.Second, 3*time.Second)
		_ = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _err := rosCsClient.CancelTask(tea.String(taskId))
			if _err != nil {
				if NeedRetry(_err) {
					wait()
					return resource.RetryableError(_err)
				}
				log.Printf("[WARN] %s ACK Cluster cancel upgrade error: %#v", clusterId, err)
			}
			return nil
		})
		// output error message
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), action, jobDetail)
	}
	// ensure cluster state is running
	csService := CsService{client: c}
	stateConf = BuildStateConf([]string{}, []string{"running"}, UpgradeClusterTimeout, 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(clusterId, []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	d.SetPartial("version")
	return nil
}

func migrateAlicloudManagedKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	action := "MigrateCluster"
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	migrateClusterRequest := map[string]string{
		"type": "ManagedKubernetes",
		"spec": d.Get("cluster_spec").(string),
	}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenAckService)
	if err != nil {
		return WrapError(err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2015-12-15"), nil, StringPointer("POST"), StringPointer("AK"), String(fmt.Sprintf("/clusters/%s/migrate", d.Id())), nil, nil, migrateClusterRequest, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	stateConf := BuildStateConf([]string{"migrating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	d.SetPartial("cluster_spec")

	return nil
}

func updateKubernetesClusterTag(d *schema.ResourceData, meta interface{}) error {
	action := "ModifyClusterTags"
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	var modifyClusterTagsRequest []cs.Tag
	if tags, err := ConvertCsTags(d); err == nil {
		modifyClusterTagsRequest = tags
	}
	d.SetPartial("tags")
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenAckService)
	if err != nil {
		return WrapError(err)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2015-12-15"), nil, StringPointer("POST"), StringPointer("AK"), String(fmt.Sprintf("/clusters/%s/tags", d.Id())), nil, nil, modifyClusterTagsRequest, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	stateConf := BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func updateControlPlaneLog(d *schema.ResourceData, meta interface{}) error {
	request := &roacs.UpdateControlPlaneLogRequest{}
	client := meta.(*connectivity.AliyunClient)
	csClient, err := client.NewRoaCsClient()
	if err != nil {
		return err
	}
	csService := CsService{client}
	if d.HasChange("control_plane_log_ttl") {
		if v, ok := d.GetOk("control_plane_log_ttl"); ok {
			request.LogTtl = tea.String(v.(string))
		}
	}
	if d.HasChange("control_plane_log_project") {
		if v, ok := d.GetOk("control_plane_log_project"); ok {
			request.LogProject = tea.String(v.(string))
		}
	}
	if d.HasChange("control_plane_log_components") {
		if v, ok := d.GetOk("control_plane_log_components"); ok {
			list := v.([]interface{})
			components := make([]*string, len(list))
			for i, c := range list {
				components[i] = tea.String(c.(string))
			}
			request.Components = components
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = csClient.UpdateControlPlaneLog(tea.String(d.Id()), request)
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
		return err
	}
	stateConf := BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func checkControlPlaneLogEnable(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return err
	}
	var response *roacs.CheckControlPlaneLogEnableResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.CheckControlPlaneLogEnable(tea.String(d.Id()))
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
		return err
	}
	if response.Body != nil {
		if response.Body.LogTtl != nil {
			d.Set("control_plane_log_ttl", response.Body.LogTtl)
		}
		if response.Body.LogProject != nil {
			d.Set("control_plane_log_project", response.Body.LogProject)
		}
		components := make([]string, len(response.Body.Components))
		for i, c := range response.Body.Components {
			components[i] = *c
		}
		d.Set("control_plane_log_components", components)
	}

	return nil
}

func updateClusterAuditConfig(d *schema.ResourceData, meta interface{}) error {
	request := &roacs.UpdateClusterAuditLogConfigRequest{}
	client := meta.(*connectivity.AliyunClient)
	csClient, err := client.NewRoaCsClient()
	if err != nil {
		return err
	}

	if d.HasChange("audit_log_config") {
		v, ok := d.GetOk("audit_log_config")
		if ok && len(v.([]interface{})) > 0 {
			m := v.([]interface{})[0].(map[string]interface{})
			if vv, ok := m["enabled"].(bool); ok {
				request.Disable = tea.Bool(!vv)
			}
			if vv, ok := m["sls_project_name"].(string); ok {
				request.SlsProjectName = tea.String(vv)
			}
		}
	}
	csService := CsService{client}
	var response *roacs.UpdateClusterAuditLogConfigResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = csClient.UpdateClusterAuditLogConfig(tea.String(d.Id()), request)
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
		return err
	}
	taskId := tea.StringValue(response.Body.TaskId)
	c := CsClient{client: csClient}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, c.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), "UpdateClusterAuditLogConfig", jobDetail)
	}

	stateConf = BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func getClusterAuditProject(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return err
	}
	var response *roacs.GetClusterAuditProjectResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.GetClusterAuditProject(tea.String(d.Id()))
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
		return err
	}
	if response.Body != nil {
		m := make(map[string]interface{})
		if response.Body.AuditEnabled != nil {
			m["enabled"] = tea.BoolValue(response.Body.AuditEnabled)
		}
		if response.Body.SlsProjectName != nil {
			m["sls_project_name"] = tea.StringValue(response.Body.SlsProjectName)
		}
		d.Set("audit_log_config", []map[string]interface{}{m})
	}

	return nil
}

func updateClusterKMSEncryption(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csClient, err := client.NewRoaCsClient()
	if err != nil {
		return err
	}

	clusterId := d.Id()

	request := &roacs.UpdateKMSEncryptionRequest{}

	if d.HasChange("disable_encryption") {
		if v, ok := d.GetOkExists("disable_encryption"); ok {
			request.DisableEncryption = tea.Bool(v.(bool))
		}
	}
	if d.HasChange("encryption_provider_key") {
		if v, ok := d.GetOk("encryption_provider_key"); ok {
			log.Printf("[DEBUG] HasChange encryption_provider_key: %s", v.(string))
			request.KmsKeyId = tea.String(v.(string))
		}
	}

	csService := CsService{client}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = csClient.UpdateKMSEncryption(tea.String(clusterId), request)
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
		return err
	}

	// Wait for cluster to return to running state
	stateConf := BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(clusterId, []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	return nil
}

func kmsEncryptionDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("disable_encryption").(bool) {
		return true
	}
	return false
}

func modifyManagedK8sCluster(d *schema.ResourceData, meta interface{}, invoker *Invoker) error {
	updated := false
	request := &roacs.ModifyClusterRequest{}
	client := meta.(*connectivity.AliyunClient)
	csClient, err := client.NewRoaCsClient()
	csService := CsService{client}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request.SetResourceGroupId(d.Get("resource_group_id").(string))
		updated = true
	}

	if !d.IsNewResource() && d.HasChange("name") {
		var clusterName string
		if v, ok := d.GetOk("name"); ok {
			clusterName = v.(string)
			request.SetClusterName(clusterName)
			updated = true
		}
	}

	// modify cluster deletion protection
	if !d.IsNewResource() && d.HasChange("deletion_protection") {
		v := d.Get("deletion_protection")
		request.SetDeletionProtection(v.(bool))
		updated = true
	}

	// modify cluster maintenance window
	if !d.IsNewResource() && d.HasChange("maintenance_window") {
		if v := d.Get("maintenance_window").([]interface{}); len(v) > 0 {
			request.MaintenanceWindow = expandMaintenanceWindowConfigRoa(v)
			updated = true
		}
		d.SetPartial("maintenance_window")
	}

	// modify cluster maintenance window
	if !d.IsNewResource() && d.HasChange("operation_policy") {
		if v := d.Get("operation_policy").([]interface{}); len(v) > 0 {
			request.OperationPolicy = &roacs.ModifyClusterRequestOperationPolicy{}
			if vv := d.Get("operation_policy.0.cluster_auto_upgrade").([]interface{}); len(vv) > 0 {
				policy := vv[0].(map[string]interface{})
				request.OperationPolicy.ClusterAutoUpgrade = &roacs.ModifyClusterRequestOperationPolicyClusterAutoUpgrade{
					Enabled: tea.Bool(policy["enabled"].(bool)),
					Channel: tea.String(policy["channel"].(string)),
				}
			}
			updated = true
		}
	}

	// modify cluster rrsa policy
	if d.HasChange("enable_rrsa") {
		enableRRSA := false
		if v, ok := d.GetOk("enable_rrsa"); ok {
			enableRRSA = v.(bool)
		}
		// it's not allowed to disable rrsa
		if !enableRRSA {
			return fmt.Errorf("It's not supported to disable RRSA! " +
				"If your cluster has enabled this function, please manually modify your tf file and add the rrsa configuration to the file.")
		}

		// version check
		clusterVersion := d.Get("version").(string)
		if res, err := versionCompare(KubernetesClusterRRSASupportedVersion, clusterVersion); res < 0 || err != nil {
			return fmt.Errorf("RRSA is not supported in current version: %s", clusterVersion)
		}
		request.SetEnableRrsa(enableRRSA)
		updated = true
		d.SetPartial("enable_rrsa")
	}

	if d.HasChange("custom_san") {
		customSan := d.Get("custom_san").(string)
		request.SetApiServerCustomCertSans(
			&roacs.ModifyClusterRequestApiServerCustomCertSans{
				SubjectAlternativeNames: tea.StringSlice(strings.Split(customSan, ",")),
				Action:                  tea.String("overwrite"),
			},
		)
		updated = true
	}

	if d.HasChange("vswitch_ids") {
		vSwitchIds := expandStringList(d.Get("vswitch_ids").([]interface{}))
		request.SetVswitchIds(tea.StringSlice(vSwitchIds))
		updated = true
	}

	if d.HasChange("timezone") {
		request.SetTimezone(d.Get("timezone").(string))
		updated = true
	}

	if d.HasChange("security_group_id") {
		request.SetSecurityGroupId(d.Get("security_group_id").(string))
		updated = true
	}

	if updated == false {
		return nil
	}

	var resp *roacs.ModifyClusterResponse
	if err := invoker.Run(func() error {
		resp, err = csClient.ModifyCluster(tea.String(d.Id()), request)
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyCluster", AlibabaCloudSdkGoERROR)
	}

	if resp == nil || resp.Body == nil {
		return nil
	}
	taskId := tea.StringValue(resp.Body.TaskId)
	c := CsClient{client: csClient}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, c.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), "ModifyCluster", jobDetail)
	}

	stateConf = BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

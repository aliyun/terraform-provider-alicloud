package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	roacs "github.com/alibabacloud-go/cs-20151215/v4/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudCSManagedKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSManagedKubernetesCreate,
		Read:   resourceAlicloudCSKubernetesRead,
		Update: resourceAlicloudCSKubernetesUpdate,
		Delete: resourceAlicloudCSKubernetesDelete,
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
				ValidateFunc:  validation.StringLenBetween(1, 63),
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ValidateFunc:  validation.StringLenBetween(0, 37),
				ConflictsWith: []string{"name"},
			},
			// worker configurations，TODO: name issue
			"worker_vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 1,
			},
			// global configurations
			"pod_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MaxItems: 10,
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
				ValidateFunc: validation.IntBetween(24, 28),
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
				ValidateFunc: validation.StringInSlice([]string{"iptables", "ipvs"}, false),
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
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
				Default:      "slb.s1.small",
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
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				ForceNew: true,
			},
			"encryption_provider_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "disk encryption key, only in ack-pro",
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
			},
			// NOTICE: 这里需要Review
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
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'slb_id' has been deprecated from provider version 1.9.2. New field 'slb_internet' replaces it.",
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
				ForceNew: true,
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
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ack.standard", "ack.pro.small"}, false),
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
							Required: true,
						},
						"maintenance_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weekly_period": {
							Type:     schema.TypeString,
							Required: true,
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
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"retain_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
		},
	}
}

func resourceAlicloudCSManagedKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()
	csService := CsService{client}
	args, err := buildKubernetesArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *cs.Client
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			args.RegionId = common.Region(client.RegionId)
			args.ClusterType = cs.ManagedKubernetes
			return csClient.CreateManagedKubernetesCluster(&cs.ManagedKubernetesClusterCreationRequest{
				ClusterArgs: args.ClusterArgs,
				WorkerArgs:  args.WorkerArgs,
			})
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_managed_kubernetes", "CreateKubernetesCluster", response)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Args"] = args
		addDebug("CreateKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ := response.(*cs.ClusterCommonResponse)
	d.SetId(cluster.ClusterID)

	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCSKubernetesRead(d, meta)
}

func UpgradeAlicloudKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("version") {
		nextVersion := d.Get("version").(string)
		args := &cs.UpgradeClusterArgs{
			Version: nextVersion,
		}

		csService := CsService{client}
		err := csService.UpgradeCluster(d.Id(), args)

		if err != nil {
			return WrapError(err)
		}

	}
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

// versionCompare check version,
// if cueVersion is newer than neededVersion return 1
// if curVersion is equal neededVersion return 0
// if curVersion is older than neededVersion return -1
// example: neededVersion = 1.20.11-aliyun.1, curVersion = 1.22.3-aliyun.1, it will return 1
func versionCompare(neededVersion, curVersion string) (int, error) {
	if neededVersion == "" || curVersion == "" {
		if neededVersion == "" && curVersion == "" {
			return 0, nil
		} else {
			if neededVersion == "" {
				return 1, nil
			} else {
				return -1, nil
			}
		}
	}

	// 取出版本号
	regx := regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)
	neededVersion = regx.FindString(neededVersion)
	curVersion = regx.FindString(curVersion)

	currentVersions := strings.Split(neededVersion, ".")
	newVersions := strings.Split(curVersion, ".")

	compare := 0

	for index, val := range currentVersions {
		newVal := newVersions[index]
		v1, err1 := strconv.Atoi(val)
		v2, err2 := strconv.Atoi(newVal)

		if err1 != nil || err2 != nil {
			return -2, fmt.Errorf("NotSupport, current cluster version is not support: %s", curVersion)
		}

		if v1 > v2 {
			compare = -1
		} else if v1 == v2 {
			compare = 0
		} else {
			compare = 1
		}

		if compare != 0 {
			break
		}
	}

	return compare, nil
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

	_, err = csClient.UpdateControlPlaneLog(tea.String(d.Id()), request)
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
	response, err := client.CheckControlPlaneLogEnable(tea.String(d.Id()))
	if err != nil {
		return err
	}
	if response.Body != nil {
		if response.Body.LogTtl != nil {
			d.Set("control_plane_log_ttl", *response.Body.LogTtl)
		}
		if response.Body.LogProject != nil {
			d.Set("control_plane_log_project", *response.Body.LogProject)
		}
		components := make([]string, len(response.Body.Components))
		for i, c := range response.Body.Components {
			components[i] = *c
		}
		d.Set("control_plane_log_components", components)
	}

	return nil
}

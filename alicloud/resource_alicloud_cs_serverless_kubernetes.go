package alicloud

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSServerlessKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSServerlessKubernetesCreate,
		Read:   resourceAlicloudCSServerlessKubernetesRead,
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
				ValidateFunc:  validation.StringLenBetween(1, 63),
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(0, 37),
				ConflictsWith: []string{"name"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'vswitch_id' has been deprecated from provider version 1.91.0. New field 'vswitch_ids' replace it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:         1,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				ConflictsWith:    []string{"vswitch_id"},
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"private_zone": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"endpoint_public_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"kube_config": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"client_cert": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"cluster_ca_cert": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				//ValidateFunc: validateCSClusterTags,
			},
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCSServerlessKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	csService := CsService{client}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	tags := make([]cs.Tag, 0)
	tagsMap, ok := d.Get("tags").(map[string]interface{})
	if ok {
		for key, value := range tagsMap {
			if value != nil {
				if v, ok := value.(string); ok {
					tags = append(tags, cs.Tag{
						Key:   key,
						Value: v,
					})
				}
			}
		}
	}

	addons := make([]cs.Addon, 0)
	if v, ok := d.GetOk("addons"); ok {
		all, ok := v.([]interface{})
		if ok {
			for _, a := range all {
				addon, ok := a.(map[string]interface{})
				if ok {
					addons = append(addons, cs.Addon{
						Name:     addon["name"].(string),
						Config:   addon["config"].(string),
						Disabled: addon["disabled"].(bool),
					})
				}
			}
		}
	}

	args := &cs.ServerlessCreationArgs{
		Name:                 clusterName,
		ClusterType:          cs.ClusterTypeServerlessKubernetes,
		RegionId:             client.RegionId,
		VpcId:                d.Get("vpc_id").(string),
		EndpointPublicAccess: d.Get("endpoint_public_access_enabled").(bool),
		PrivateZone:          d.Get("private_zone").(bool),
		NatGateway:           d.Get("new_nat_gateway").(bool),
		SecurityGroupId:      d.Get("security_group_id").(string),
		Addons:               addons,
		KubernetesVersion:    d.Get("version").(string),
		DeletionProtection:   d.Get("deletion_protection").(bool),
		ResourceGroupId:      d.Get("resource_group_id").(string),
	}

	if v := d.Get("vswitch_id").(string); v != "" {
		args.VSwitchId = v
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		args.VswitchIds = expandStringList(v.([]interface{}))
	}

	//set tags
	if len(tags) > 0 {
		args.Tags = tags
	}

	var requestInfo *cs.Client
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.CreateServerlessKubernetesCluster(args)
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_serverless_kubernetes", "CreateServerlessKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Args"] = args
		addDebug("CreateServerlessKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ := response.(*cs.ClusterCommonResponse)
	d.SetId(cluster.ClusterID)

	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, csService.CsServerlessKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCSServerlessKubernetesRead(d, meta)
}

func resourceAlicloudCSServerlessKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()
	object, err := csService.DescribeCsServerlessKubernetes(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	_ = d.Set("name", object.Name)
	_ = d.Set("vpc_id", object.VpcId)
	_ = d.Set("vswitch_id", object.VSwitchId)
	_ = d.Set("security_group_id", object.SecurityGroupId)
	_ = d.Set("deletion_protection", object.DeletionProtection)
	_ = d.Set("tags", object.Tags)
	_ = d.Set("version", object.CurrentVersion)
	_ = d.Set("resource_group_id", object.ResourceGroupId)
	_ = d.Set("cluster_spec", object.ClusterSpec)

	var requestInfo *cs.Client
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.GetClusterCerts(d.Id())
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetClusterCerts", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = d.Id()
		addDebug("GetClusterCerts", response, requestInfo, requestMap)
	}
	cert, _ := response.(cs.ClusterCerts)
	if ce, ok := d.GetOk("client_cert"); ok && ce.(string) != "" {
		if err := writeToFile(ce.(string), cert.Cert); err != nil {
			return WrapError(err)
		}
	}
	if key, ok := d.GetOk("client_key"); ok && key.(string) != "" {
		if err := writeToFile(key.(string), cert.Key); err != nil {
			return WrapError(err)
		}
	}
	if ca, ok := d.GetOk("cluster_ca_cert"); ok && ca.(string) != "" {
		if err := writeToFile(ca.(string), cert.CA); err != nil {
			return WrapError(err)
		}
	}

	var config *cs.ClusterConfig
	if file, ok := d.GetOk("kube_config"); ok && file.(string) != "" {
		var requestInfo *cs.Client

		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.DescribeClusterUserConfig(d.Id(), !d.Get("endpoint_public_access_enabled").(bool))
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetClusterConfig", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			addDebug("GetClusterConfig", response, requestInfo, requestMap)
		}
		config, _ = response.(*cs.ClusterConfig)

		if err := writeToFile(file.(string), config.Config); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAlicloudCSServerlessKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	var requestInfo *cs.Client
	invoker := NewInvoker()
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.DeleteCluster(d.Id())
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = d.Id()
		addDebug("DeleteCluster", response, requestInfo, requestMap)
	}

	stateConf := BuildStateConf([]string{"running", "deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, csService.CsServerlessKubernetesInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

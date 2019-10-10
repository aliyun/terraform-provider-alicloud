package alicloud

import (
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCSServelessKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSServelessKubernetesCreate,
		Read:   resourceAlicloudCSServelessKubernetesRead,
		//todo The serveless cluster dose not support scaling,so Update method is not provided
		//Update: resourceAlicloudCSServelessKubernetesUpdate,
		Delete: resourceAlicloudCSServelessKubernetesDelete,
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
				ValidateFunc:  validateContainerName,
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ForceNew:      true,
				ValidateFunc:  validateContainerNamePrefix,
				ConflictsWith: []string{"name"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"enndpoint_public_access_enabled": {
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
		},
	}
}

func resourceAlicloudCSServelessKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

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

	args := &cs.ServerlessCreationArgs{
		Name:                 clusterName,
		ClusterType:          cs.ClusterTypeServelessKubernetes,
		RegionId:             client.RegionId,
		VpcId:                d.Get("vpc_id").(string),
		VSwitchId:            d.Get("vswitch_id").(string),
		EndpointPublicAccess: d.Get("enndpoint_public_access_enabled").(bool),
		PrivateZone:          d.Get("private_zone").(bool),
		NatGateway:           d.Get("new_nat_gateway").(bool),
		DeletionProtection:   d.Get("deletion_protection").(bool),
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
			return csClient.CreateServelessKubernetesCluster(args)
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_serveless_kubernetes", "CreateServelessKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Args"] = args
		addDebug("CreateServelessKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ := response.(*cs.ClusterCreationResponse)
	d.SetId(cluster.ClusterID)

	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return nil, csClient.WaitForClusterAsyn(d.Id(), cs.Running, 3600)
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "WaitForClusterAsyn", DenverdinoAliyungo)
	}
	if debugOn() {
		waitRequestMap := make(map[string]interface{})
		waitRequestMap["ClusterId"] = d.Id()
		waitRequestMap["Status"] = cs.Running
		waitRequestMap["TimeOut"] = 3600
		addDebug("WaitForClusterAsyn", response, requestInfo, waitRequestMap)
	}

	return resourceAlicloudCSServelessKubernetesRead(d, meta)
}

func resourceAlicloudCSServelessKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()
	object, err := csService.DescribeCsServelessKubernetes(d.Id())
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
	_ = d.Set("private_zone", object.PrivateZone)
	_ = d.Set("deletion_protection", object.DeletionProtection)
	_ = d.Set("tags", object.Tags)

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

	var config cs.ClusterConfig
	if file, ok := d.GetOk("kube_config"); ok && file.(string) != "" {
		var requestInfo *cs.Client

		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.DescribeClusterUserConfig(d.Id(), d.Get("enndpoint_public_access_enabled").(bool))
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
		config, _ = response.(cs.ClusterConfig)

		if err := writeToFile(file.(string), config.Config); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAlicloudCSServelessKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlicloudCSServelessKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	var requestInfo *cs.Client
	invoker := NewInvoker()
	var response interface{}

	err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.DeleteCluster(d.Id())
			})
			response = raw
			return err
		}); err != nil {
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			addDebug("DeleteCluster", response, requestInfo, requestMap)
		}
		return nil
	})
	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteCluster", DenverdinoAliyungo)
	}
	return WrapError(csService.WaitForCSServelessKubernetes(d.Id(), Deleted, DefaultLongTimeout))
}

package alicloud

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	KubernetesClusterNetworkTypeFlannel = "flannel"
	KubernetesClusterNetworkTypeTerway  = "terway"

	KubernetesClusterLoggingTypeSLS = "SLS"
)

var (
	KubernetesClusterNodeCIDRMasksByDefault = 24
)

func resourceAlicloudCSKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesCreate,
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
				Deprecated:    "Field 'name_prefix' has been deprecated from provider version 1.75.0.",
			},
			// master configurations
			"master_vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:         3,
				MaxItems:         5,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems:         3,
				MaxItems:         5,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(40, 500),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_instance_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:          PostPaid,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validation.StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				// must be a valid period, expected [1-9], 12, 24, 36, 48 or 60,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			// worker configurations
			"worker_vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:         1,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 10,
			},
			"worker_number": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"worker_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(20, 32768),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_data_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(20, 32768),
				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
			},
			"worker_data_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_data_disks": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"worker_instance_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:          PostPaid,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validation.StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"exclude_autoscaler_nodes": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			// global configurations
			// Terway network
			"pod_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MaxItems:         10,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			// Flannel network
			"pod_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"service_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"node_cidr_mask": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc:     validation.IntBetween(24, 28),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"key_name", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"password", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"user_ca": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"enable_ssh": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"image_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: imageIdSuppressFunc,
			},
			"install_cloud_monitor": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// cpu policy options of kubelet
			"cpu_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "static"}, false),
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
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
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"slb_internet_enabled": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			// computed parameters
			"kube_config": {
				Type:     schema.TypeString,
				Optional: true,
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
			"connections": {
				Type:     schema.TypeMap,
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
			"master_nodes": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"worker_nodes": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// remove parameters below
			// mix vswitch_ids between master and worker is not a good guidance to create cluster
			"worker_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_instance_type' has been removed from provider version 1.75.0. New field 'worker_instance_types' replaces it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:         3,
				MaxItems:         5,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				Removed:          "Field 'vswitch_ids' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replace it.",
			},
			// single instance type would cause extra troubles
			"master_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'master_instance_type' has been removed from provider version 1.75.0. New field 'master_instance_types' replaces it.",
			},
			// force update is a high risk operation
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Removed:  "Field 'force_update' has been removed from provider version 1.75.0.",
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// single az would be never supported.
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'vswitch_id' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replaces it.",
			},
			// worker_numbers in array is a hell of management
			"worker_numbers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:    schema.TypeInt,
					Default: 3,
				},
				MinItems:         1,
				MaxItems:         3,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				Removed:          "Field 'worker_numbers' has been removed from provider version 1.75.0. New field 'worker_number' replaces it.",
			},
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Removed:  "Field 'nodes' has been removed from provider version 1.9.4. New field 'master_nodes' replaces it.",
			},
			// too hard to use this config
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{KubernetesClusterLoggingTypeSLS}, false),
							Required:     true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				Removed:          "Field 'log_config' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			"cluster_network_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				Removed:          "Field 'cluster_network_type' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
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
		},
	}
}

func resourceAlicloudCSKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var requestInfo *cs.Client
	var raw interface{}

	// prepare args and set default value
	args, err := buildKubernetesArgs(d, meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "PrepareKubernetesClusterArgs", err)
	}

	if err = invoker.Run(func() error {
		raw, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			args.RegionId = common.Region(client.RegionId)
			args.ClusterType = cs.DelicatedKubernetes
			return csClient.CreateDelicatedKubernetesCluster(args)
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "CreateKubernetesCluster", raw)
	}

	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Params"] = args
		addDebug("CreateKubernetesCluster", raw, requestInfo, requestMap)
	}

	cluster, ok := raw.(*cs.ClusterCommonResponse)
	if ok != true {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "ParseKubernetesClusterResponse", raw)
	}

	d.SetId(cluster.ClusterID)

	// reset interval to 10s
	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCSKubernetesUpdate(d, meta)
}

func resourceAlicloudCSKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	d.Partial(true)
	invoker := NewInvoker()
	if d.HasChange("worker_number") && !d.IsNewResource() {
		password := d.Get("password").(string)
		if password == "" {
			if v := d.Get("kms_encrypted_password").(string); v != "" {
				kmsService := KmsService{client}
				decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
				if err != nil {
					return WrapError(err)
				}
				password = decryptResp.Plaintext
			}
		}

		keyPair := d.Get("key_name").(string)

		oldV, newV := d.GetChange("worker_number")

		oldValue, ok := oldV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("worker_number old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("worker_number new value can not be parsed"), "parseError %d", newValue)
		}

		if newValue < oldValue {
			return WrapErrorf(fmt.Errorf("worker_number can not be less than before"), "scaleOutFailed %d:%d", newValue, oldValue)
		}

		args := &cs.ScaleOutKubernetesClusterRequest{
			KeyPair:             keyPair,
			LoginPassword:       password,
			ImageId:             d.Get("image_id").(string),
			UserData:            d.Get("user_data").(string),
			Count:               int64(newValue) - int64(oldValue),
			WorkerVSwitchIds:    expandStringList(d.Get("worker_vswitch_ids").([]interface{})),
			WorkerInstanceTypes: expandStringList(d.Get("worker_instance_types").([]interface{})),
		}

		if v := d.Get("user_data").(string); v != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v)
			if base64DecodeError == nil {
				args.UserData = v
			} else {
				args.UserData = base64.StdEncoding.EncodeToString([]byte(v))
			}
		}

		if v, ok := d.GetOk("worker_instance_charge_type"); ok {
			args.WorkerInstanceChargeType = v.(string)
			if args.WorkerInstanceChargeType == string(PrePaid) {
				args.WorkerAutoRenew = d.Get("worker_auto_renew").(bool)
				args.WorkerAutoRenewPeriod = d.Get("worker_auto_renew_period").(int)
				args.WorkerPeriod = d.Get("worker_period").(int)
				args.WorkerPeriodUnit = d.Get("worker_period_unit").(string)
			}
		}

		if d.HasChange("worker_data_disks") {
			if dds, ok := d.GetOk("worker_data_disks"); ok {
				disks := dds.([]interface{})
				createDataDisks := make([]cs.DataDisk, 0, len(disks))
				for _, e := range disks {
					pack := e.(map[string]interface{})
					dataDisk := cs.DataDisk{
						Size:                 pack["size"].(string),
						DiskName:             pack["name"].(string),
						Category:             pack["category"].(string),
						Device:               pack["device"].(string),
						AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
						KMSKeyId:             pack["kms_key_id"].(string),
						Encrypted:            pack["encrypted"].(string),
					}
					createDataDisks = append(createDataDisks, dataDisk)
				}
				args.WorkerDataDisks = createDataDisks
			}
			d.SetPartial("worker_data_disks")
		}

		var resoponse interface{}
		if err := invoker.Run(func() error {
			var err error
			resoponse, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				resp, err := csClient.ScaleOutKubernetesCluster(d.Id(), args)
				return resp, err
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ResizeKubernetesCluster", DenverdinoAliyungo)
		}
		if debugOn() {
			resizeRequestMap := make(map[string]interface{})
			resizeRequestMap["ClusterId"] = d.Id()
			resizeRequestMap["Args"] = args
			addDebug("ResizeKubernetesCluster", resoponse, resizeRequestMap)
		}

		stateConf := BuildStateConf([]string{"scaling"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("worker_number")
	}

	if !d.IsNewResource() && (d.HasChange("name") || d.HasChange("name_prefix")) {
		var clusterName string
		if v, ok := d.GetOk("name"); ok {
			clusterName = v.(string)
		} else {
			clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
		}
		var requestInfo *cs.Client
		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return nil, csClient.ModifyClusterName(d.Id(), clusterName)
			})
			response = raw
			return err
		}); err != nil && !IsExpectedErrors(err, []string{"ErrorClusterNameAlreadyExist"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyClusterName", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			requestMap["ClusterName"] = clusterName
			addDebug("ModifyClusterName", response, requestInfo, requestMap)
		}
		d.SetPartial("name")
		d.SetPartial("name_prefix")
	}

	UpgradeAlicloudKubernetesCluster(d, meta)
	d.Partial(false)
	return resourceAlicloudCSKubernetesRead(d, meta)
}

func resourceAlicloudCSKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()
	object, err := csService.DescribeCsKubernetes(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("vpc_id", object.VpcId)
	d.Set("security_group_id", object.SecurityGroupId)
	d.Set("version", object.CurrentVersion)
	d.Set("worker_ram_role_name", object.WorkerRamRoleName)

	var masterNodes []map[string]interface{}
	var workerNodes []map[string]interface{}

	pageNumber := 1
	for {
		var result []cs.KubernetesNodeType
		var pagination *cs.PaginationResult
		var requestInfo *cs.Client
		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				nodes, paginationResult, err := csClient.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
				return []interface{}{nodes, paginationResult}, err
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}
			addDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
		}
		result, _ = response.([]interface{})[0].([]cs.KubernetesNodeType)
		pagination, _ = response.([]interface{})[1].(*cs.PaginationResult)

		if pageNumber == 1 && (len(result) == 0 || result[0].InstanceId == "") {
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				if err := invoker.Run(func() error {
					raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
						requestInfo = csClient
						nodes, _, err := csClient.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
						return nodes, err
					})
					response = raw
					return err
				}); err != nil {
					return resource.NonRetryableError(err)
				}
				tmp, _ := response.([]cs.KubernetesNodeType)
				if len(tmp) > 0 && tmp[0].InstanceId != "" {
					result = tmp
				}

				for _, stableState := range cs.NodeStableClusterState {
					// If cluster is in NodeStableClusteState, node list will not change
					if object.State == string(stableState) {
						if debugOn() {
							requestMap := make(map[string]interface{})
							requestMap["ClusterId"] = d.Id()
							requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}
							addDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
						}
						return nil
					}
				}
				time.Sleep(5 * time.Second)
				return resource.RetryableError(Error("[ERROR] There is no any nodes in kubernetes cluster %s.", d.Id()))
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
			}

		}

		if d.Get("exclude_autoscaler_nodes").(bool) {
			result, err = knockOffAutoScalerNodes(result, meta)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", AlibabaCloudSdkGoERROR)
			}
		}

		for _, node := range result {
			mapping := map[string]interface{}{
				"id":         node.InstanceId,
				"name":       node.InstanceName,
				"private_ip": node.IpAddress[0],
			}
			if node.InstanceRole == "Master" {
				masterNodes = append(masterNodes, mapping)
			} else {
				workerNodes = append(workerNodes, mapping)
			}
		}

		if len(result) < pagination.PageSize {
			break
		}
		pageNumber += 1
	}

	d.Set("master_nodes", masterNodes)
	d.Set("worker_nodes", workerNodes)
	d.Set("worker_number", int64(len(workerNodes)))

	// Get slb information
	connection := make(map[string]string)

	if len(masterNodes) != 0 {
		request := slb.CreateDescribeLoadBalancersRequest()
		request.ServerId = masterNodes[0]["id"].(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeLoadBalancers(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		lbs, _ := raw.(*slb.DescribeLoadBalancersResponse)
		for _, lb := range lbs.LoadBalancers.LoadBalancer {
			if strings.ToLower(lb.AddressType) == strings.ToLower(string(Internet)) {
				d.Set("slb_internet", lb.LoadBalancerId)
				connection["api_server_internet"] = fmt.Sprintf("https://%s:6443", lb.Address)
				connection["master_public_ip"] = lb.Address
			} else {
				d.Set("slb_intranet", lb.LoadBalancerId)
				connection["api_server_intranet"] = fmt.Sprintf("https://%s:6443", lb.Address)

				reqVpc := vpc.CreateDescribeEipAddressesRequest()
				reqVpc.AssociatedInstanceId = lb.LoadBalancerId
				reqVpc.AssociatedInstanceType = "SlbInstance"
				raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
					return vpcClient.DescribeEipAddresses(reqVpc)
				})
				eip, _ := raw.(*vpc.DescribeEipAddressesResponse)
				if eip != nil && len(eip.EipAddresses.EipAddress) > 0 {
					eipAddr := eip.EipAddresses.EipAddress[0].IpAddress
					connection["master_public_ip"] = eipAddr
					connection["api_server_internet"] = fmt.Sprintf("https://%s:6443", eipAddr)
				}
			}
		}
	}

	connection["service_domain"] = fmt.Sprintf("*.%s.%s.alicontainer.com", d.Id(), object.RegionId)

	d.Set("connections", connection)
	natRequest := vpc.CreateDescribeNatGatewaysRequest()
	natRequest.VpcId = object.VpcId
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
		requestMap["Id"] = d.Id()
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
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.GetClusterConfig(d.Id())
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetClusterConfig", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["Id"] = d.Id()
			addDebug("GetClusterConfig", response, requestInfo, requestMap)
		}
		config, _ = response.(cs.ClusterConfig)

		if err := writeToFile(file.(string), config.Config); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func resourceAlicloudCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var response interface{}
	err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.DeleteKubernetesCluster(d.Id())
			})
			response = raw
			return err
		}); err != nil {
			return resource.RetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			addDebug("DeleteCluster", response, d.Id(), requestMap)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteCluster", DenverdinoAliyungo)
	}

	stateConf := BuildStateConf([]string{"running", "deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildKubernetesArgs(d *schema.ResourceData, meta interface{}) (*cs.DelicatedKubernetesClusterCreationRequest, error) {
	client := meta.(*connectivity.AliyunClient)

	vpcService := VpcService{client}

	var vswitchID string
	if list := expandStringList(d.Get("worker_vswitch_ids").([]interface{})); len(list) > 0 {
		vswitchID = list[0]
	} else {
		vswitchID = ""
	}

	var vpcId string
	if vswitchID != "" {
		vsw, err := vpcService.DescribeVSwitch(vswitchID)
		if err != nil {
			return nil, err
		}
		vpcId = vsw.VpcId
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
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

	var apiAudiences string
	if list := expandStringList(d.Get("api_audiences").([]interface{})); len(list) > 0 {
		apiAudiences = strings.Join(list, ",")
	}

	creationArgs := &cs.DelicatedKubernetesClusterCreationRequest{
		ClusterArgs: cs.ClusterArgs{
			DisableRollback: true,
			Name:            clusterName,
			// TODO add windows and other platform support
			OsType:   "Linux",
			Platform: "CentOS",
			VpcId:    vpcId,

			// the params below is ok to be empty
			KubernetesVersion:         d.Get("version").(string),
			NodeCidrMask:              strconv.Itoa(d.Get("node_cidr_mask").(int)),
			ImageId:                   d.Get("image_id").(string),
			KeyPair:                   d.Get("key_name").(string),
			ServiceCidr:               d.Get("service_cidr").(string),
			CloudMonitorFlags:         d.Get("install_cloud_monitor").(bool),
			SecurityGroupId:           d.Get("security_group_id").(string),
			IsEnterpriseSecurityGroup: d.Get("is_enterprise_security_group").(bool),
			EndpointPublicAccess:      d.Get("slb_internet_enabled").(bool),
			SnatEntry:                 d.Get("new_nat_gateway").(bool),
			NodeNameMode:              d.Get("node_name_mode").(string),
			Addons:                    addons,
			ServiceAccountIssuer:      d.Get("service_account_issuer").(string),
			ApiAudiences:              apiAudiences,
		},
	}

	if dds, ok := d.GetOk("worker_data_disks"); ok {
		disks := dds.([]interface{})
		createDataDisks := make([]cs.DataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := cs.DataDisk{
				Size:                 pack["size"].(string),
				DiskName:             pack["name"].(string),
				Category:             pack["category"].(string),
				Device:               pack["device"].(string),
				AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
				KMSKeyId:             pack["kms_key_id"].(string),
				Encrypted:            pack["encrypted"].(string),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		creationArgs.WorkerDataDisks = createDataDisks
	}

	if v := d.Get("user_data").(string); v != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v)
		if base64DecodeError == nil {
			creationArgs.UserData = v
		} else {
			creationArgs.UserData = base64.StdEncoding.EncodeToString([]byte(v))
		}
	}

	if _, ok := d.GetOk("pod_vswitch_ids"); ok {
		creationArgs.PodVswitchIds = expandStringList(d.Get("pod_vswitch_ids").([]interface{}))
	} else {
		creationArgs.ContainerCidr = d.Get("pod_cidr").(string)
	}

	if password := d.Get("password").(string); password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return nil, WrapError(err)
			}
			password = decryptResp.Plaintext
		}
		creationArgs.LoginPassword = password
	} else {
		creationArgs.LoginPassword = password
	}

	// CA default is empty
	if userCa, ok := d.GetOk("user_ca"); ok {
		userCaContent, err := loadFileContent(userCa.(string))
		if err != nil {
			return nil, fmt.Errorf("reading user_ca file failed %s", err)
		}
		creationArgs.UserCa = string(userCaContent)
	}

	// set proxy mode and default is ipvs
	if proxyMode := d.Get("proxy_mode").(string); proxyMode != "" {
		creationArgs.ProxyMode = cs.ProxyMode(proxyMode)
	} else {
		creationArgs.ProxyMode = cs.ProxyMode(cs.IPVS)
	}

	// dedicated kubernetes must provide master_vswitch_ids
	if _, ok := d.GetOk("master_vswitch_ids"); ok {
		creationArgs.MasterArgs = cs.MasterArgs{
			MasterCount:              len(d.Get("master_vswitch_ids").([]interface{})),
			MasterVSwitchIds:         expandStringList(d.Get("master_vswitch_ids").([]interface{})),
			MasterInstanceTypes:      expandStringList(d.Get("master_instance_types").([]interface{})),
			MasterSystemDiskCategory: d.Get("master_disk_category").(string),
			MasterSystemDiskSize:     int64(d.Get("master_disk_size").(int)),
			// TODO support other params
		}
	}

	if v, ok := d.GetOk("master_instance_charge_type"); ok {
		creationArgs.MasterInstanceChargeType = v.(string)
		if creationArgs.MasterInstanceChargeType == string(PrePaid) {
			creationArgs.MasterAutoRenew = d.Get("master_auto_renew").(bool)
			creationArgs.MasterAutoRenewPeriod = d.Get("master_auto_renew_period").(int)
			creationArgs.MasterPeriod = d.Get("master_period").(int)
			creationArgs.MasterPeriodUnit = d.Get("master_period_unit").(string)
		}
	}

	creationArgs.WorkerArgs = cs.WorkerArgs{
		WorkerVSwitchIds:         expandStringList(d.Get("worker_vswitch_ids").([]interface{})),
		WorkerInstanceTypes:      expandStringList(d.Get("worker_instance_types").([]interface{})),
		NumOfNodes:               int64(d.Get("worker_number").(int)),
		WorkerSystemDiskCategory: d.Get("worker_disk_category").(string),
		WorkerSystemDiskSize:     int64(d.Get("worker_disk_size").(int)),
		// TODO support other params
	}

	if v, ok := d.GetOk("worker_instance_charge_type"); ok {
		creationArgs.WorkerInstanceChargeType = v.(string)
		if creationArgs.WorkerInstanceChargeType == string(PrePaid) {
			creationArgs.WorkerAutoRenew = d.Get("worker_auto_renew").(bool)
			creationArgs.WorkerAutoRenewPeriod = d.Get("worker_auto_renew_period").(int)
			creationArgs.WorkerPeriod = d.Get("worker_period").(int)
			creationArgs.WorkerPeriodUnit = d.Get("worker_period_unit").(string)
		}
	}
	return creationArgs, nil
}

func knockOffAutoScalerNodes(nodes []cs.KubernetesNodeType, meta interface{}) ([]cs.KubernetesNodeType, error) {
	log.Printf("[DEBUG] start to knock off auto scaler nodes %++v\n", nodes)
	client := meta.(*connectivity.AliyunClient)
	scaleNoddesMap := make(map[string]ecs.Instance)
	result := make([]cs.KubernetesNodeType, 0)
	instanceIds := make([]interface{}, 0)

	if len(nodes) == 0 {
		return result, nil
	}

	for _, node := range nodes {
		instanceIds = append(instanceIds, node.InstanceId)
	}

	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(len(nodes))
	request.InstanceIds = convertListToJsonString(instanceIds)
	tags := []ecs.DescribeInstancesTag{{Key: defaultScalingGroupTag, Value: "true"}}
	request.Tag = &tags

	// filter ecs by tags, find all ecs created by autoscaler
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstances(request)
	})

	if err != nil {
		return result, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*ecs.DescribeInstancesResponse)
	for _, instance := range response.Instances.Instance {
		scaleNoddesMap[instance.InstanceId] = instance
	}

	for _, node := range nodes {
		if _, ok := scaleNoddesMap[node.InstanceId]; !ok {
			result = append(result, node)
		}
	}

	return result, nil
}

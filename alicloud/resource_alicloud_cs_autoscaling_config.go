package alicloud

import (
	"time"

	"github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	cs "github.com/alibabacloud-go/cs-20151215/v3/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const resourceName = "resource_alicloud_cs_autoscaling_config"

func resourceAlicloudCSAutoscalingConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSAutoscalingConfigCreate,
		Read:   resourceAlicloudCSAutoscalingConfigRead,
		Update: resourceAlicloudCSAutoscalingConfigUpdate,
		Delete: resourceAlicloudCSAutoscalingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cool_down_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10m",
			},
			"unneeded_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10m",
			},
			"utilization_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.5",
			},
			"gpu_utilization_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.5",
			},
			"scan_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "30s",
			},
			"scale_down_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"expander": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "least-waste",
			},
			"skip_nodes_with_system_pods": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"skip_nodes_with_local_storage": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"daemonset_eviction_for_nodes": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"max_graceful_termination_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_replica_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"recycle_node_deletion_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scale_up_from_zero": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceAlicloudCSAutoscalingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlicloudCSAutoscalingConfigUpdate(d, meta)
}

func resourceAlicloudCSAutoscalingConfigRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlicloudCSAutoscalingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	// cluster id
	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	// auto scaling config
	updateAutoscalingConfigRequest := &cs.CreateAutoscalingConfigRequest{}
	if v, ok := d.GetOk("cool_down_duration"); ok {
		updateAutoscalingConfigRequest.CoolDownDuration = tea.String(v.(string))
	}
	if v, ok := d.GetOk("unneeded_duration"); ok {
		updateAutoscalingConfigRequest.UnneededDuration = tea.String(v.(string))
	}
	if v, ok := d.GetOk("utilization_threshold"); ok {
		updateAutoscalingConfigRequest.UtilizationThreshold = tea.String(v.(string))
	}
	if v, ok := d.GetOk("gpu_utilization_threshold"); ok {
		updateAutoscalingConfigRequest.GpuUtilizationThreshold = tea.String(v.(string))
	}
	if v, ok := d.GetOk("scan_interval"); ok {
		updateAutoscalingConfigRequest.ScanInterval = tea.String(v.(string))
	}
	enableScaleDown := d.Get("scale_down_enabled").(bool)
	updateAutoscalingConfigRequest.ScaleDownEnabled = tea.Bool(enableScaleDown)
	if v, ok := d.GetOk("expander"); ok {
		updateAutoscalingConfigRequest.Expander = tea.String(v.(string))
	}
	skipNodesWithPods := d.Get("skip_nodes_with_system_pods").(bool)
	updateAutoscalingConfigRequest.SkipNodesWithSystemPods = tea.Bool(skipNodesWithPods)
	skipNodesWithLocalStorage := d.Get("skip_nodes_with_local_storage").(bool)
	updateAutoscalingConfigRequest.SkipNodesWithLocalStorage = tea.Bool(skipNodesWithLocalStorage)
	evictDaemonset := d.Get("daemonset_eviction_for_nodes").(bool)
	updateAutoscalingConfigRequest.DaemonsetEvictionForNodes = tea.Bool(evictDaemonset)
	if v, ok := d.GetOk("max_graceful_termination_sec"); ok {
		updateAutoscalingConfigRequest.MaxGracefulTerminationSec = tea.Int32(v.(int32))
	}
	if v, ok := d.GetOk("min_replica_count"); ok {
		updateAutoscalingConfigRequest.MinReplicaCount = tea.Int32(v.(int32))
	}
	enableDeleteRecycleNode := d.Get("recycle_node_deletion_enabled").(bool)
	updateAutoscalingConfigRequest.RecycleNodeDeletionEnabled = tea.Bool(enableDeleteRecycleNode)
	scaleUpFromZero := d.Get("scale_up_from_zero").(bool)
	updateAutoscalingConfigRequest.ScaleUpFromZero = tea.Bool(scaleUpFromZero)

	_, err = client.CreateAutoscalingConfig(tea.String(clusterId), updateAutoscalingConfigRequest)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "CreateAutoscalingConfig", AliyunTablestoreGoSdk)
	}

	addDebug("CreateAutoscalingConfig", updateAutoscalingConfigRequest, err)
	d.SetId(clusterId)
	d.Partial(false)

	return resourceAlicloudCSAutoscalingConfigRead(d, meta)
}

func resourceAlicloudCSAutoscalingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

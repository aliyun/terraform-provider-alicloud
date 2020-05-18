package alicloud

import (
	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEdasK8sCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasK8sClusterCreate,
		Read:   resourceAlicloudEdasK8sClusterRead,
		Delete: resourceAlicloudEdasK8sClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cs_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"cluster_type": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
				//ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3, 4, 5}),
			},
			"network_mode": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
				//ValidateFunc: validation.IntInSlice([]int{1, 2}),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEdasK8sClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateImportK8sClusterRequest()
	request.RegionId = client.RegionId
	request.ClusterId = d.Get("cs_cluster_id").(string)
	if v, ok := d.GetOk("namespace_id"); ok {
		request.NamespaceId = v.(string)
	}
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ImportK8sCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_k8s_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ImportK8sClusterResponse)
	if response.Code != 200 {
		return WrapError(Error("create k8s cluster failed for " + response.Message))
	}
	d.SetId(response.Data)

	return resourceAlicloudEdasK8sClusterRead(d, meta)
}

func resourceAlicloudEdasK8sClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	clusterId := d.Id()
	regionId := client.RegionId

	request := edas.CreateGetClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetClusterResponse)
	if response.Code != 200 {
		return WrapError(Error("create cluster failed for " + response.Message))
	}

	d.Set("cluster_name", response.Cluster.ClusterName)
	d.Set("cluster_type", response.Cluster.ClusterType)
	d.Set("network_mode", response.Cluster.NetworkMode)
	d.Set("vpc_id", response.Cluster.VpcId)
	d.Set("region_id", response.Cluster.RegionId)

	return nil
}

func resourceAlicloudEdasK8sClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	clusterId := d.Id()
	regionId := client.RegionId

	request := edas.CreateDeleteClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteCluster(request)
		})
		response, _ := raw.(*edas.DeleteClusterResponse)
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if response.Code != 200 {
			if strings.Contains(response.Message, "there are still instances in it") {
				return resource.RetryableError(Error("delete cluster failed for " + response.Message))
			}
			return resource.NonRetryableError(Error("delete cluster failed for " + response.Message))
		}

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}

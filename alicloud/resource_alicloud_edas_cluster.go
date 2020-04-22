package alicloud

import (
	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEdasCluster() *schema.Resource {
	return &schema.Resource{
		Create: rresourceAlicloudEdasClusterCreate,
		Read:   resourceAlicloudEdasClusterRead,
		Delete: resourceAlicloudEdasClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"network_mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func rresourceAlicloudEdasClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	clusterName := d.Get("cluster_name").(string)
	regionId := client.RegionId
	clusterType := d.Get("cluster_type").(int)
	networkMode := d.Get("network_mode").(int)
	vpcId := d.Get("vpc_id").(string)
	if networkMode == 2 && len(vpcId) == 0 {
		return Error("vpcId is required for vpc network mode")
	}
	logicalRegionId := d.Get("logical_region_id").(string)

	request := edas.CreateInsertClusterRequest()
	request.RegionId = regionId
	request.ClusterName = clusterName
	request.ClusterType = requests.NewInteger(clusterType)
	request.NetworkMode = requests.NewInteger(networkMode)
	request.LogicalRegionId = logicalRegionId
	request.VpcId = vpcId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.InsertCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.InsertClusterResponse)
	if response.Code != 200 {
		return Error("create cluster failed for " + response.Message)
	}
	d.SetId(response.Cluster.ClusterId)

	return resourceAlicloudEdasClusterRead(d, meta)
}

func resourceAlicloudEdasClusterRead(d *schema.ResourceData, meta interface{}) error {
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
		return Error("create cluster failed for " + response.Message)
	}

	d.Set("cluster_name", response.Cluster.ClusterName)
	d.Set("cluster_type", response.Cluster.ClusterType)
	d.Set("network_mode", response.Cluster.NetworkMode)
	d.Set("vpc_id", response.Cluster.VpcId)

	return nil
}

func resourceAlicloudEdasClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	clusterId := d.Id()
	regionId := client.RegionId

	request := edas.CreateDeleteClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteCluster(request)
		})
		response, _ := raw.(*edas.DeleteClusterResponse)
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(10 * time.Second)
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

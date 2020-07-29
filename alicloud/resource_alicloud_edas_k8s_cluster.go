package alicloud

import (
	"strings"
	"time"

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
			},
			"network_mode": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
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
			"cluster_import_status": {
				Type:     schema.TypeInt,
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
		return WrapError(Error("import k8s cluster failed for " + response.Message))
	}
	if len(response.Data) == 0 {
		return WrapError(Error("null cluster id after import k8s cluster"))
	}
	d.SetId(response.Data)
	// 需要获取集群直到导入成功
	req := edas.CreateGetClusterRequest()
	req.ClusterId = response.Data
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetCluster(req)
		})
		response, _ := raw.(*edas.GetClusterResponse)
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if response.Code != 200 {
			return resource.NonRetryableError(Error("Get cluster failed for " + response.Message))
		}

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		if response.Cluster.ClusterImportStatus == 3 {
			return resource.RetryableError(Error("cluster is importing"))
		}
		if response.Cluster.ClusterImportStatus == 1 {
			return nil
		}

		return resource.NonRetryableError(Error("cluster status abnormal"))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

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
	d.Set("cluster_import_status", response.Cluster.ClusterImportStatus)
	d.Set("cs_cluster_id", response.Cluster.CsClusterId)

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
			return resource.NonRetryableError(Error("delete cluster failed for " + response.Message))
		}

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	//等待集群删除成功
	reqGet := edas.CreateGetClusterRequest()
	reqGet.RegionId = regionId
	reqGet.ClusterId = clusterId
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetCluster(reqGet)
		})
		response, _ := raw.(*edas.GetClusterResponse)
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)

		if response.Code == 200 {
			return resource.RetryableError(Error("cluster deleting"))
		} else if response.Code == 601 && strings.Contains(response.Message, "does not exist") {
			return nil
		} else {
			return resource.NonRetryableError(Error("check cluster status failed for " + response.Message))
		}
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}

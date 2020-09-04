package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/mse"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMseCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMseClusterCreate,
		Read:   resourceAlicloudMseClusterRead,
		Update: resourceAlicloudMseClusterUpdate,
		Delete: resourceAlicloudMseClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_entry_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_alias_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_specification": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MSE_SC_1_2_200_c", "MSE_SC_2", "MSE_SC_4_8_200_c_4_200_c", "MSE_SC_8_16_200_c"}, false),
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Eureka", "Nacos-Ans", "ZooKeeper"}, false),
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"net_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"privatenet", "pubnet"}, false),
			},
			"private_slb_specification": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pub_network_flow": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pub_slb_specification": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudMseClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}

	request := mse.CreateCreateClusterRequest()
	request.ClusterSpecification = d.Get("cluster_specification").(string)
	request.ClusterType = d.Get("cluster_type").(string)
	request.ClusterVersion = d.Get("cluster_version").(string)
	if v, ok := d.GetOk("disk_type"); ok {
		request.DiskType = v.(string)
	}

	request.InstanceCount = requests.NewInteger(d.Get("instance_count").(int))
	request.NetType = d.Get("net_type").(string)
	if v, ok := d.GetOk("private_slb_specification"); ok {
		request.PrivateSlbSpecification = v.(string)
	}

	if v, ok := d.GetOk("pub_network_flow"); ok {
		request.PubNetworkFlow = v.(string)
	}

	if v, ok := d.GetOk("pub_slb_specification"); ok {
		request.PubSlbSpecification = v.(string)
	}

	request.Region = client.RegionId
	if v, ok := d.GetOk("vswitch_id"); ok {
		request.VSwitchId = v.(string)
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request.VpcId = vsw.VpcId

	}
	raw, err := client.WithMseClient(func(mseClient *mse.Client) (interface{}, error) {
		return mseClient.CreateCluster(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mse_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*mse.CreateClusterResponse)
	d.SetId(fmt.Sprintf("%v", response.InstanceId))
	stateConf := BuildStateConf([]string{}, []string{"INIT_SUCCESS"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, mseService.MseClusterStateRefreshFunc(d.Id(), []string{"INIT_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudMseClusterUpdate(d, meta)
}
func resourceAlicloudMseClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	object, err := mseService.DescribeMseCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_cluster mseService.DescribeMseCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_type", object.ClusterType)
	d.Set("instance_count", object.InstanceCount)
	d.Set("pub_network_flow", object.PubNetworkFlow)
	d.Set("status", object.InitStatus)
	return nil
}
func resourceAlicloudMseClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	if d.HasChange("acl_entry_list") {
		request := mse.CreateUpdateAclRequest()
		request.InstanceId = d.Id()
		request.AclEntryList = convertListToCommaSeparate(d.Get("acl_entry_list").(*schema.Set).List())

		raw, err := client.WithMseClient(func(mseClient *mse.Client) (interface{}, error) {
			return mseClient.UpdateAcl(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("acl_entry_list")
	}
	update := false
	request := mse.CreateUpdateClusterRequest()
	request.InstanceId = d.Id()
	if d.HasChange("cluster_alias_name") {
		update = true
		request.ClusterAliasName = d.Get("cluster_alias_name").(string)
	}
	if update {
		raw, err := client.WithMseClient(func(mseClient *mse.Client) (interface{}, error) {
			return mseClient.UpdateCluster(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("cluster_alias_name")
	}
	d.Partial(false)
	return resourceAlicloudMseClusterRead(d, meta)
}
func resourceAlicloudMseClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	request := mse.CreateDeleteClusterRequest()
	request.InstanceId = d.Id()
	raw, err := client.WithMseClient(func(mseClient *mse.Client) (interface{}, error) {
		return mseClient.DeleteCluster(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"DESTROY_SUCCESS"}, d.Timeout(schema.TimeoutDelete), 60*time.Second, mseService.MseClusterStateRefreshFunc(d.Id(), []string{"DESTROY_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

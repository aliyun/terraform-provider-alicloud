package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/cs"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func resourceAlicloudContainerCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudContainerClusterCreate,
		Read:   resourceAlicloudContainerClusterRead,
		Update: resourceAlicloudContainerClusterUpdate,
		Delete: resourceAlicloudContainerClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validateContainerClusterName,
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateContainerClusterNamePrefix,
			},
			"size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateIntegerInRange(1, 20),
			},
			"cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      20,
				ValidateFunc: validateIntegerInRange(20, 32768),
			},
			"disk_category": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ecs.DiskCategoryCloudEfficiency,
				ForceNew:     true,
				ValidateFunc: validateDiskCategory,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudContainerClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.csconn

	// Ensure instance_type is generation three
	_, err := meta.(*AliyunClient).CheckParameterValidity(d, meta)
	if err != nil {
		return err
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		clusterName = resource.PrefixedUniqueId(v.(string))
	} else {
		clusterName = resource.UniqueId()
	}

	args := &cs.ClusterCreationArgs{
		Name:             clusterName,
		InstanceType:     d.Get("instance_type").(string),
		Password:         d.Get("password").(string),
		Size:             int64(d.Get("size").(int)),
		IOOptimized:      ecs.IoOptimized("true"),
		DataDiskCategory: ecs.DiskCategory(d.Get("disk_category").(string)),
		DataDiskSize:     int64(d.Get("disk_size").(int)),
	}

	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		cidr, cidr_ok := d.GetOk("cidr_block")
		if !cidr_ok || cidr.(string) == "" {
			return fmt.Errorf("When launching container in the VPC, the 'cidr_block' must be specified.")
		}
		args.NetworkMode = cs.VPCNetwork
		args.VSwitchID = v.(string)
		args.SubnetCIDR = cidr.(string)

		vswInfo, _, err := client.vpcconn.DescribeVSwitches(&ecs.DescribeVSwitchesArgs{
			RegionId:  getRegion(d, meta),
			VSwitchId: v.(string),
		})
		if err != nil {
			return fmt.Errorf("Error DescribeVSwitches: %#v", err)
		}
		if len(vswInfo) < 1 {
			return fmt.Errorf("There is not found specified vswitch: %s, please check and try again.", v.(string))
		}
		args.VPCID = vswInfo[0].VpcId
	} else {
		args.NetworkMode = cs.ClassicNetwork
	}

	if imageId, ok := d.GetOk("image_id"); ok {
		connection := client.ecsconn
		argsImage := &ecs.DescribeImagesArgs{
			RegionId: getRegion(d, meta),
			ImageId:  imageId.(string),
		}
		if _, _, err := connection.DescribeImages(argsImage); err != nil {
			return err
		}

		args.ECSImageID = imageId.(string)
	}

	region := getRegion(d, meta)
	cluster, err := conn.CreateCluster(region, args)

	if err != nil {
		return fmt.Errorf("Creating container Cluster got an error: %#v", err)
	}

	err = conn.WaitForClusterAsyn(cluster.ClusterID, cs.Running, 500)

	if err != nil {
		return fmt.Errorf("Waitting for container Cluster %#v got an error: %#v", cs.Running, err)
	}

	d.SetId(cluster.ClusterID)

	return resourceAlicloudContainerClusterUpdate(d, meta)
}

func resourceAlicloudContainerClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).csconn
	d.Partial(true)
	if d.HasChange("size") && !d.IsNewResource() {
		o, n := d.GetChange("size")
		oi := o.(int)
		ni := n.(int)
		if ni <= oi {
			return fmt.Errorf("The new size of clusters must greater than the current. The cluster's current size is %d.", oi)
		}
		d.SetPartial("size")
		err := conn.ResizeCluster(d.Id(), &cs.ClusterResizeArgs{
			Size:             int64(ni),
			InstanceType:     d.Get("instance_type").(string),
			Password:         d.Get("password").(string),
			DataDiskCategory: ecs.DiskCategory(d.Get("disk_category").(string)),
			DataDiskSize:     int64(d.Get("disk_size").(int)),
			ECSImageID:       d.Get("image_id").(string),
			IOOptimized:      ecs.IoOptimized("true"),
		})
		if err != nil {
			return fmt.Errorf("Resize Cluster got an error: %#v", err)
		}

		err = conn.WaitForClusterAsyn(d.Id(), cs.Running, 500)

		if err != nil {
			return fmt.Errorf("Waitting for container Cluster %#v got an error: %#v", cs.Running, err)
		}
	}
	d.Partial(false)

	return resourceAlicloudContainerClusterRead(d, meta)
}

func resourceAlicloudContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).csconn

	cluster, err := conn.DescribeCluster(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", cluster.Name)
	d.Set("size", cluster.Size)
	d.Set("network_mode", cluster.NetworkMode)
	if cluster.VPCID != "" {
		d.Set("vpc_id", cluster.VPCID)
	}
	if cluster.VSwitchID != "" {
		d.Set("vswitch_id", cluster.VSwitchID)
	}

	return nil
}

func resourceAlicloudContainerClusterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).csconn

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		err := conn.DeleteCluster(d.Id())
		if err != nil {
			if IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Cluster in use 1- trying again while it is deleted."))
		}

		resp, err := conn.DescribeCluster(d.Id())
		if err != nil {
			if IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting container cluster got an error: %#v", err))
		}
		if resp.ClusterID == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Cluster in use 2- trying again while it is deleted."))
	})
}

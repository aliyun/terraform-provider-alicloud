package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/cs"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validateContainerClusterName,
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ValidateFunc:  validateContainerClusterNamePrefix,
				ConflictsWith: []string{"name"},
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"new_nat_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"master_instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateInstanceType,
			},
			"worker_instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateInstanceType,
			},
			"worker_number": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validateIntegerInRange(1, 50),
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"pod_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ssh": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"master_disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validateIntegerInRange(40, 500),
			},
			"master_disk_category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  ecs.DiskCategoryCloudEfficiency,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ecs.DiskCategoryCloudEfficiency), string(ecs.DiskCategoryCloudSSD)}),
			},
			"worker_disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validateIntegerInRange(20, 32768),
			},
			"worker_disk_category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  ecs.DiskCategoryCloudEfficiency,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ecs.DiskCategoryCloudEfficiency), string(ecs.DiskCategoryCloudSSD)}),
			},
			"install_cloud_monitor": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"docker_version": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCSKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.csconn

	// Ensure instance_type is generation three
	args, err := buildKunernetesArgs(d, meta)
	if err != nil {
		return err
	}

	cluster, err := conn.CreateKubernetesCluster(getRegion(d, meta), args)

	if err != nil {
		return fmt.Errorf("Creating Kubernetes Cluster got an error: %#v", err)
	}

	d.SetId(cluster.ClusterID)

	if err := conn.WaitForClusterAsyn(cluster.ClusterID, cs.Running, 1800); err != nil {
		return fmt.Errorf("Waitting for kubernetes cluster %#v got an error: %#v", cs.Running, err)
	}

	return resourceAlicloudCSKubernetesUpdate(d, meta)
}

func resourceAlicloudCSKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).csconn
	d.Partial(true)
	if d.HasChange("worker_number") && !d.IsNewResource() {
		// Ensure instance_type is generation three
		args, err := buildKunernetesArgs(d, meta)
		if err != nil {
			return err
		}
		if err := conn.ResizeKubernetes(d.Id(), args); err != nil {
			return fmt.Errorf("Resize Cluster got an error: %#v", err)
		}

		err = conn.WaitForClusterAsyn(d.Id(), cs.Running, 500)

		if err != nil {
			return fmt.Errorf("Waitting for container Cluster %#v got an error: %#v", cs.Running, err)
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
		if err := conn.ModifyClusterName(d.Id(), clusterName); err != nil && !IsExceptedError(err, ErrorClusterNameAlreadyExist) {
			return fmt.Errorf("Modify Cluster Name got an error: %#v", err)
		}
		d.SetPartial("name")
		d.SetPartial("name_prefix")
	}
	d.Partial(false)

	return resourceAlicloudCSKubernetesRead(d, meta)
}

func resourceAlicloudCSKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).csconn

	cluster, err := conn.DescribeCluster(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", cluster.Name)
	// Each k8s cluster contains 3 master nodes
	d.Set("worker_number", cluster.Size-KubernetesMasterNumber)
	d.Set("vswitch_id", cluster.VSwitchID)
	d.Set("docker_version", cluster.DockerVersion)

	return nil
}

func resourceAlicloudCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).csconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DeleteCluster(d.Id())
		if err != nil {
			if IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete Kubernetes Cluster timeout and get an error: %#v.", err))
		}

		resp, err := conn.DescribeCluster(d.Id())
		if err != nil {
			if IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describing Kubernetes Cluster got an error: %#v", err))
		}
		if resp.ClusterID == "" {
			return nil
		}

		if string(resp.State) == string(Deleting) {
			time.Sleep(5 * time.Second)
		}

		return resource.RetryableError(fmt.Errorf("Delete Kubernetes Cluster timeout and get an error: %#v.", err))
	})
}

func buildKunernetesArgs(d *schema.ResourceData, meta interface{}) (*cs.KubernetesCreationArgs, error) {
	client := meta.(*AliyunClient)

	// Ensure instance_type is generation three
	_, err := meta.(*AliyunClient).CheckParameterValidity(d, meta)
	if err != nil {
		return nil, err
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	stackArgs := &cs.KubernetesStackArgs{
		MasterInstanceType:       d.Get("master_instance_type").(string),
		WorkerInstanceType:       d.Get("worker_instance_type").(string),
		Password:                 d.Get("password").(string),
		NumOfNodes:               int64(d.Get("worker_number").(int)),
		MasterSystemDiskCategory: ecs.DiskCategory(d.Get("master_disk_category").(string)),
		MasterSystemDiskSize:     int64(d.Get("master_disk_size").(int)),
		WorkerSystemDiskCategory: ecs.DiskCategory(d.Get("worker_disk_category").(string)),
		WorkerSystemDiskSize:     int64(d.Get("worker_disk_size").(int)),
		SNatEntry:                d.Get("new_nat_gateway").(bool),
		KubernetesVersion:        KubernetesVersion,
		DockerVersion:            KubernetesDockerVersion,
		ContainerCIDR:            d.Get("pod_cidr").(string),
		ServiceCIDR:              d.Get("service_cidr").(string),
		SSHFlags:                 d.Get("enable_ssh").(bool),
		ImageID:                  KubernetesImageId,
		CloudMonitorFlags:        d.Get("install_cloud_monitor").(bool),
	}
	if v, ok := d.GetOk("availability_zone"); ok && len(Trim(v.(string))) > 0 {
		stackArgs.ZoneId = Trim(v.(string))
	}
	if v, ok := d.GetOk("vswitch_id"); ok && len(Trim(v.(string))) > 0 {
		stackArgs.VSwitchID = Trim(v.(string))
		vsw, err := client.DescribeVswitch(stackArgs.VSwitchID)
		if err != nil {
			return nil, err
		}
		stackArgs.VPCID = vsw.VpcId
		if stackArgs.ZoneId != "" && vsw.ZoneId != vsw.ZoneId {
			return nil, fmt.Errorf("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, stackArgs.ZoneId)
		}
		stackArgs.ZoneId = vsw.ZoneId
	} else if !stackArgs.SNatEntry {
		return nil, fmt.Errorf("The automatic created VPC and VSwitch must set 'new_nat_gateway' to 'true'.")
	}

	return &cs.KubernetesCreationArgs{
		Name:              clusterName,
		ClusterType:       "Kubernetes",
		DisableRollback:   true,
		TimeoutMins:       30,
		KubernetesVersion: stackArgs.KubernetesVersion,
		StackParams:       *stackArgs,
	}, nil
}

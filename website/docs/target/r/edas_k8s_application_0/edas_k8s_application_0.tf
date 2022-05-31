resource "alicloud_edas_k8s_application" "default" {
  // package type is Image / FatJar / War
  package_type            = "Image"
  application_name        = "DemoApplication"
  application_descriotion = "This is description of application"
  cluster_id              = var.cluster_id
  replicas                = 2

  // set 'image_url' and 'repo_id' when package_type is 'image'
  image_url = "registry-vpc.cn-beijing.aliyuncs.com/edas-demo-image/consumer:1.0"

  // set 'package_url','package_version' and 'jdk' when package_type is not 'image'
  package_url     = var.package_url
  package_version = var.package_version
  jdk             = var.jdk

  // set 'web_container' and 'edas_container' when package_type is 'war'
  web_container          = var.web_container
  edas_container_version = var.edas_container_version

  internet_target_port  = var.internet_target_port
  internet_slb_port     = var.internet_slb_port
  internet_slb_protocol = var.internet_slb_protocol
  internet_slb_id       = var.internet_slb_id
  limit_cpu             = 4
  limit_mem             = 2048
  requests_cpu          = 0
  requests_mem          = 0
  requests_m_cpu        = 0
  limit_m_cpu           = 4000
  command               = var.command
  command_args          = var.command_args
  envs                  = var.envs
  pre_stop              = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  post_start            = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  liveness              = var.liveness
  readiness             = var.readiness
  nas_id                = var.nas_id
  mount_descs           = var.mount_descs
  local_volume          = var.local_volume
  namespace             = "default"
  logical_region_id     = cn-beijing
}

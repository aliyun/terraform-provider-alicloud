package alicloud

import "github.com/denverdino/aliyungo/ecs"

type GroupRuleNicType string

const (
	GroupRuleInternet = GroupRuleNicType("internet")
	GroupRuleIntranet = GroupRuleNicType("intranet")
)

type GroupRulePolicy string

const (
	GroupRulePolicyAccept = GroupRulePolicy("accept")
	GroupRulePolicyDrop   = GroupRulePolicy("drop")
)

const (
	EcsApiVersion20160314 = "2016-03-14"
	EcsApiVersion20140526 = "2014-05-26"
)

const GenerationOne = "ecs-1"
const GenerationTwo = "ecs-2"
const GenerationThree = "ecs-3"
const GenerationFour = "ecs-4"

var NoneIoOptimizedFamily = map[string]string{"ecs.t1": "", "ecs.t2": "", "ecs.s1": ""}
var NoneIoOptimizedInstanceType = map[string]string{"ecs.s2.small": ""}
var HalfIoOptimizedFamily = map[string]string{"ecs.s2": "", "ecs.s3": "", "ecs.m1": "", "ecs.m2": "", "ecs.c1": "", "ecs.c2": ""}

var OutdatedDiskCategory = map[ecs.DiskCategory]ecs.DiskCategory{
	ecs.DiskCategoryCloud: ecs.DiskCategoryCloud}

var SupportedDiskCategory = map[ecs.DiskCategory]ecs.DiskCategory{
	ecs.DiskCategoryCloudSSD:        ecs.DiskCategoryCloudSSD,
	ecs.DiskCategoryCloudEfficiency: ecs.DiskCategoryCloudEfficiency,
	ecs.DiskCategoryCloud:           ecs.DiskCategoryCloud}

const AllPortRange = "-1/-1"

const (
	KubernetesImageId       = "centos_7"
	KubernetesMasterNumber  = 3
	KubernetesVersion       = "1.9.3"
	KubernetesDockerVersion = "17.06.2-ce-1"
)

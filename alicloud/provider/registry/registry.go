// Package registry aggregates the service package declarations consumed by the
// provider assembly loop, and validates uniqueness, this being the only place both
// provider servers agree on.
package registry

import (
	"fmt"
	"strings"
	"sync"

	alicloudfunction "github.com/aliyun/terraform-provider-alicloud/alicloud/function"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/actiontrail"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/adb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/adcp"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/agentloop"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/aiworkspace"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/aligreen"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/amqpopen"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/antiddospublic"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/apig"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/bpstudio"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/brainindustrial"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/bssopenapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cddc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/chatbot"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/clickhouse"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudauth"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudcontrol"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfw"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudphone"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudsso"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/computenest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/das"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/databasegateway"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/datahub"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworkspublic"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dbfs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dbs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dcdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ddosbgp"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ddoscoo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/devopsrdc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dfs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dm"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dmsenterprise"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/drds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dysmsapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eais"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ebs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecd"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eci"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/edsaic"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/edsuser"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ehpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eipanycast"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/elasticsearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/emr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ens"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eventbridge"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnectrouter"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fcopen"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fnf"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ga"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/gdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/governance"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/gpdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/green"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/gwlb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/hbr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/hitsdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/hologram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/imm"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/imp"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ims"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/iot"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kvcachestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/live"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/market"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/maxcompute"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mhub"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/milvus"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mnsopen"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mscopensubscription"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mse"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/natgateway"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nlb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oceanbasepro"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oos"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/openapiexplorer"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/opensearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/paiflow"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardbx"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatelink"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/pvtz"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/quickbipublic"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/quotas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rdsai"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcesharing"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rkvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rocketmq"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ros"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sae"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/scdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/schedulerx2"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sddp"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/selectdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/servicecatalog"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/servicemesh"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sgw"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/starrocks"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sts"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/swasopen"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/tag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ververica"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vod"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpcipam"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpcpeer"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/wafopenapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/yundunbastionhost"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/yundundbaudit"
)

// A managed resource and a data source may share a name; two managed resources may
// not, whichever server implements them.
const (
	nsResource          = "resource"
	nsDataSource        = "data source"
	nsFunction          = "function"
	nsEphemeralResource = "ephemeral resource"
	nsListResource      = "list resource"
	nsAction            = "action"
)

var servicePackages = []conns.ServicePackage{
	actiontrail.ServicePackage(),
	adb.ServicePackage(),
	adcp.ServicePackage(),
	agentloop.ServicePackage(),
	aiworkspace.ServicePackage(),
	alb.ServicePackage(),
	alicloudfunction.ServicePackage(),
	alidns.ServicePackage(),
	aligreen.ServicePackage(),
	alikafka.ServicePackage(),
	amqpopen.ServicePackage(),
	antiddospublic.ServicePackage(),
	apig.ServicePackage(),
	arms.ServicePackage(),
	bpstudio.ServicePackage(),
	brainindustrial.ServicePackage(),
	bssopenapi.ServicePackage(),
	cas.ServicePackage(),
	cassandra.ServicePackage(),
	cbn.ServicePackage(),
	cddc.ServicePackage(),
	cdn.ServicePackage(),
	chatbot.ServicePackage(),
	clickhouse.ServicePackage(),
	cloudapi.ServicePackage(),
	cloudauth.ServicePackage(),
	cloudcontrol.ServicePackage(),
	cloudfirewall.ServicePackage(),
	cloudfw.ServicePackage(),
	natgateway.ServicePackage(),
	cloudphone.ServicePackage(),
	cloudsso.ServicePackage(),
	cms.ServicePackage(),
	computenest.ServicePackage(),
	config.ServicePackage(),
	cr.ServicePackage(),
	cs.ServicePackage(),
	das.ServicePackage(),
	datahub.ServicePackage(),
	dataworkspublic.ServicePackage(),
	dbfs.ServicePackage(),
	dbs.ServicePackage(),
	dcdn.ServicePackage(),
	ddosbgp.ServicePackage(),
	ddoscoo.ServicePackage(),
	dds.ServicePackage(),
	devopsrdc.ServicePackage(),
	dfs.ServicePackage(),
	databasegateway.ServicePackage(),
	dm.ServicePackage(),
	dms.ServicePackage(),
	dmsenterprise.ServicePackage(),
	drds.ServicePackage(),
	dts.ServicePackage(),
	dysmsapi.ServicePackage(),
	eais.ServicePackage(),
	eas.ServicePackage(),
	ebs.ServicePackage(),
	ecd.ServicePackage(),
	eci.ServicePackage(),
	ecs.ServicePackage(),
	edas.ServicePackage(),
	edsaic.ServicePackage(),
	edsuser.ServicePackage(),
	eflo.ServicePackage(),
	ehpc.ServicePackage(),
	eipanycast.ServicePackage(),
	elasticsearch.ServicePackage(),
	emr.ServicePackage(),
	ens.ServicePackage(),
	esa.ServicePackage(),
	ess.ServicePackage(),
	eventbridge.ServicePackage(),
	expressconnectrouter.ServicePackage(),
	fc.ServicePackage(),
	fcopen.ServicePackage(),
	fnf.ServicePackage(),
	ga.ServicePackage(),
	gdb.ServicePackage(),
	governance.ServicePackage(),
	gpdb.ServicePackage(),
	green.ServicePackage(),
	gwlb.ServicePackage(),
	hbase.ServicePackage(),
	hbr.ServicePackage(),
	hitsdb.ServicePackage(),
	hologram.ServicePackage(),
	imm.ServicePackage(),
	imp.ServicePackage(),
	ims.ServicePackage(),
	iot.ServicePackage(),
	kms.ServicePackage(),
	kvcachestore.ServicePackage(),
	live.ServicePackage(),
	market.ServicePackage(),
	maxcompute.ServicePackage(),
	mhub.ServicePackage(),
	milvus.ServicePackage(),
	mns.ServicePackage(),
	mnsopen.ServicePackage(),
	mscopensubscription.ServicePackage(),
	mse.ServicePackage(),
	nas.ServicePackage(),
	nlb.ServicePackage(),
	oceanbasepro.ServicePackage(),
	ons.ServicePackage(),
	oos.ServicePackage(),
	openapiexplorer.ServicePackage(),
	opensearch.ServicePackage(),
	oss.ServicePackage(),
	ots.ServicePackage(),
	paiflow.ServicePackage(),
	polardb.ServicePackage(),
	polardbx.ServicePackage(),
	privatelink.ServicePackage(),
	pvtz.ServicePackage(),
	quickbipublic.ServicePackage(),
	quotas.ServicePackage(),
	ram.ServicePackage(),
	rds.ServicePackage(),
	rdsai.ServicePackage(),
	resourcemanager.ServicePackage(),
	resourcesharing.ServicePackage(),
	rkvstore.ServicePackage(),
	rocketmq.ServicePackage(),
	ros.ServicePackage(),
	sae.ServicePackage(),
	sas.ServicePackage(),
	scdn.ServicePackage(),
	schedulerx2.ServicePackage(),
	sddp.ServicePackage(),
	selectdb.ServicePackage(),
	servicecatalog.ServicePackage(),
	servicemesh.ServicePackage(),
	sgw.ServicePackage(),
	slb.ServicePackage(),
	sls.ServicePackage(),
	smartag.ServicePackage(),
	starrocks.ServicePackage(),
	sts.ServicePackage(),
	swasopen.ServicePackage(),
	tag.ServicePackage(),
	ververica.ServicePackage(),
	vod.ServicePackage(),
	vpc.ServicePackage(),
	vpcipam.ServicePackage(),
	vpcpeer.ServicePackage(),
	vs.ServicePackage(),
	wafopenapi.ServicePackage(),
	yundunbastionhost.ServicePackage(),
	yundundbaudit.ServicePackage(),
}

var validateOnce sync.Once

// Validates the set once on first use, and panics on a violation: failing at
// startup beats a nil dereference on a later RPC.
func ServicePackages() []conns.ServicePackage {
	validateOnce.Do(func() {
		if err := validate(servicePackages); err != nil {
			panic("provider registration error: " + err.Error())
		}
	})
	return servicePackages
}

// Reports every violation at once. A declaration colliding with a provider.go map
// literal is assembleProvider's to find, not visible here.
func validate(sps []conns.ServicePackage) error {
	v := &validator{seen: map[string]map[string]string{}}
	for _, sp := range sps {
		for _, d := range sp.SDKResources {
			v.claim(nsResource, d.TypeName, sp.Name, "SDKResources", d.Factory != nil)
		}
		for _, d := range sp.Resources {
			v.claim(nsResource, d.TypeName, sp.Name, "Resources", d.Factory != nil)
		}
		for _, d := range sp.SDKDataSources {
			v.claim(nsDataSource, d.TypeName, sp.Name, "SDKDataSources", d.Factory != nil)
		}
		for _, d := range sp.DataSources {
			v.claim(nsDataSource, d.TypeName, sp.Name, "DataSources", d.Factory != nil)
		}
		for _, d := range sp.Functions {
			v.claim(nsFunction, d.TypeName, sp.Name, "Functions", d.Factory != nil)
		}
		for _, d := range sp.EphemeralResources {
			v.claim(nsEphemeralResource, d.TypeName, sp.Name, "EphemeralResources", d.Factory != nil)
		}
		for _, d := range sp.ListResources {
			v.claim(nsListResource, d.TypeName, sp.Name, "ListResources", d.Factory != nil)
		}
		for _, d := range sp.Actions {
			v.claim(nsAction, d.TypeName, sp.Name, "Actions", d.Factory != nil)
		}
	}
	return v.err()
}

type validator struct {
	seen     map[string]map[string]string // namespace → type name → "Service.Category"
	problems []string
}

func (v *validator) claim(namespace, typeName, service, category string, hasFactory bool) {
	origin := service + "." + category
	if typeName == "" {
		v.problems = append(v.problems, fmt.Sprintf("%s declared in %s has an empty TypeName", namespace, origin))
		return
	}
	if !hasFactory {
		v.problems = append(v.problems, fmt.Sprintf("%s %q declared in %s has a nil Factory", namespace, typeName, origin))
	}
	names := v.seen[namespace]
	if names == nil {
		names = map[string]string{}
		v.seen[namespace] = names
	}
	if prev, dup := names[typeName]; dup {
		v.problems = append(v.problems, fmt.Sprintf("%s %q is declared twice: %s and %s", namespace, typeName, prev, origin))
		return
	}
	names[typeName] = origin
}

func (v *validator) err() error {
	if len(v.problems) == 0 {
		return nil
	}
	return fmt.Errorf("%d invalid service package declaration(s):\n  %s", len(v.problems), strings.Join(v.problems, "\n  "))
}

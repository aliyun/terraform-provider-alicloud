### 前置说明
+ CLI部分功能需要依赖网络，且在办公网络下才能运行。
+ extend，$引用资源资源等高级特性镇元2.0的白屏暂不支持，如果您使用了这些特性白屏将切换为只读模式。
+ 安装请参考：安装[CloudSpec IDL](https://aliyuque.antfin.com/cloudspec/model/cli-install)。

### 主要支持特性
CloudSpec CLI主要提供的主要特性如下：

| 特性 | 说明 |
| --- | --- |
| 版本管理 | 查询当前CLI的版本、是否有更新。 |
| 代码格式化 | 在您手写IDL后，支持运行命令直接进行IDL代码的格式化。 |
| 脚手架 | 帮助您快速初始化项目工程。 |
| 编译检查 | 帮您检查整个IDL工程是否有语法错误。 |
| 规范检查 | 运行[CloudSpec规范](http://cloudspec.aliyun-inc.com/)检查，其中break change的检查支持单独执行。 |
| 基于resource自动生成operation | 基于资源自动生成operation，该特性请参考：[自动化生成说明](https://aliyuque.antfin.com/cloudspec/model/gpceg4uzvr9ck1ht) |
| VSCode插件管理 | 帮您快速安装、更新CloudSpec支持的VSCode插件。 |
| YAML生成 | 基于您的IDL工程导出阿里云支持的YAML格式的资源、API、RAM Meta。 |
| 执行选择器 | 应用原生的选择器语法对构建的CloudSpec对象执行选择，并返回选择的结果。 |
| 基于operation自动生成resource | 基于存量的operation自动生成resource。[可参考详细说明](https://aliyuque.antfin.com/cloudspec/model/couka8le4gooqozc) |


### 能力
#### version 版本更新
查询当前cloudspec的版本及检查更新，示例：

```bash
cloudspec version
```

输出如下：

```bash
1.0.0
当前版本是最新版本
```

#### vscode 安装插件
安装vscode插件，安装需要code作为可执行文件加到PATH中，否则会下载插件到本地，需要手工进入VSCode安装。

示例如下：

```bash
(base) lazy@B-73KY41WF-1106 ~ % cloudspec vscode
1.85.1
0ee08df0cf4527e40edc9aa28f4b5bd38bbff2b2
arm64
正在安装Vscode插件：code --install-extension /Users/lazy/cloudspec-0.0.8.vsix
Installing extensions...
(node:53146) [DEP0005] DeprecationWarning: Buffer() is deprecated due to security and usability issues. Please use the Buffer.alloc(), Buffer.allocUnsafe(), or Buffer.from() methods instead.
(Use `Electron --trace-deprecation ...` to show where the warning was created)
Extension 'cloudspec-0.0.8.vsix' was successfully installed.
```

#### create 快速初始化
脚手架能力，快速初始化service，resource及operation。

##### service 初始化服务
初始化service，要求初始化的目录下没有main.cspec入口文件。

支持的参数：

| 参数 | 是否必填 | 说明 |
| --- | --- | --- |
| --namespace | 是 | 服务所属的namespace，例如ECS |
| --name | 是 | 服务名称，例如Ecs |
| --version | 是 | 服务的版本，日期版本，例如2024-01-01 |
| --service_type | 是 | 服务下的operation的网关类别，支持  + rpc 所有的operation都是RPC风格的API  + roa 所有的operation都是ROA风格的API |


示例：

```bash
(base) lazy@B-73KY41WF-1106 cli % cloudspec create service --namespace=ECS --name=Ecs --version=2014-05-26 --service_type=rpc
初始化成功
修改文件：main.cspec
写入文件成功：/private/tmp/cli/main.cspec
# 查看生成的问题
(base) lazy@B-73KY41WF-1106 cli % cat main.cspec
$version: 1
namespace: alicloud.ECS.Ecs.v20140526

@runtimeType("pop")
@apiStyle("rpc")
service Ecs {
  version: "2014-05-26"
}%  
```

##### operation 初始化operation
初始化operation。

支持的参数：

| 参数 | 是否必填 | 说明 |
| --- | --- | --- |
| --name | 是 | operation名称，例如CreateInstance |
| --backend_type | 是 | 后端服务类型，支持：  + http  + hsf  + dubbo  + http_hsf |


示例如下：

```bash
(base) lazy@B-73KY41WF-1106 cli % cloudspec create operation --name=CreateInstance --backend_type=http                       
生成成功
修改文件：main.cspec
写入文件成功：/private/tmp/cli/main.cspec
修改文件：operations/CreateInstance.cspec
写入文件成功：/private/tmp/cli/operations/CreateInstance.cspec
```

##### resource 初始化资源
初始化resource。

支持的参数：

| 参数 | 是否必填 | 说明 |
| --- | --- | --- |
| --name | 是 | 资源名称，例如Instance |
| --type | 是 | 后端服务类型，支持：  + normal 普通资源  + singleton 单例资源  + virtual 虚拟资源  + readonly 只读资源  + relation 关系资源 |


示例如下：

```bash
(base) lazy@B-73KY41WF-1106 cli % cloudspec create resource --name=Instance --type=normal 
生成成功
修改文件：main.cspec
写入文件成功：/private/tmp/cli/main.cspec
修改文件：resources/Instance.cspec
写入文件成功：/private/tmp/cli/resources/Instance.cspec
```


#### build 编译语法检查
编译当前目录，检查是否有语法错误。

成功示例如下：

```bash
(base) lazy@B-73KY41WF-1106 Config_pop_Config_2019-01-08 % pwd
/Users/lazy/code/Config_pop_Config_2019-01-08
(base) lazy@B-73KY41WF-1106 Config_pop_Config_2019-01-08 % cloudspec build
解析成功
解析时间：1621ms
```

失败示例如下：

```bash
(base) lazy@B-73KY41WF-1106 error-spec % cloudspec build
解析失败
Position -> File: /main.cspec:6:5, Line: 6, Column: 5, Start: 6:5, Stop: 6:5.  The component A$b/C is not defined
解析时间：6ms
解析错误信息：Position -> File: /main.cspec:6:5, Line: 6, Column: 5, Start: 6:5, Stop: 6:5.  The component A$b/C is not defined
```

#### select 执行选择器
对编译的cloudspec对象运行选择器，选择器语法请参考：[https://aliyuque.antfin.com/cloudspec/model/vpp7m0c2l4q1spc1](https://aliyuque.antfin.com/cloudspec/model/vpp7m0c2l4q1spc1) 

支持参数如下：

```bash
cloudspec select --selector=""
```

selector是选择器的内容，必填，示例如下：

```bash
(base) lazy@B-73KY41WF-1106 Config_pop_Config_2020-09-07 % pwd
/Users/lazy/code/Config_pop_Config_2020-09-07
(base) lazy@B-73KY41WF-1106 Config_pop_Config_2020-09-07 %cloudspec select --selector="operation > struct[id|name=UpdateAggregateRemediation_Input]"
选择执行成功
解析时间：1894ms
选择器初始化时间：72ms
选择器执行时间：79ms
选择结果：
alicloud.Config.Config.v20200907#UpdateAggregateRemediation_Input
```

#### check 检查CloudSpec规范
执行约束器检查，**不检查breakChange**，默认运行镇元使用的约束器，如果运行自定义的约束器，请通过--rule_url覆盖：

```bash
cloudspec check --rule_url=""
```

示例如下：

```bash
(base) lazy@B-73KY41WF-1106 test % cloudspec check
选择执行成功
解析时间：4ms
规则执行时间：695ms
运行规则数量：169
错误数量：71
警告数量：0
note数量：3
选择结果：
约束器名称：alicloud.test#M-AC-0002，规范链接：http://cloudspec.aliyun-inc.com/?q=M-AC-0002
事件等级：DANGER
事件信息：API风格 取值范围：RPC (Remote Procedure Call)：更多用于面向过程的设计
根组件ID：alicloud.test#opC
根组件类型：operation
位置信息：文件路径：main.cspec，行号：19，列号：10
约束器名称：alicloud.test#M-AC-0002，规范链接：http://cloudspec.aliyun-inc.com/?q=M-AC-0002
事件等级：DANGER
...
```

#### bc 检查BreakChange规范
执行break change检查，默认运行镇元使用的约束器，如果运行自定义的约束器，请通过--rule_url覆盖，另外必须通过--master_path指定master的地址。

:::warning
关于master地址：请使用绝对路径，另外本地开发时，请将master clone到另外的文件夹下。

:::

示例如下：

```bash
(base) lazy@B-73KY41WF-1106 idl_bc % cloudspec bc test -mp=/Users/lazy/code/cloudspec-cli/idl_bc/test_master
选择执行成功
解析时间：5ms
规则执行时间：628ms
运行规则数量：84
错误数量：4
警告数量：0
note数量：0
选择结果：
约束器名称：alicloud.test#R-BR-0005，规范链接：http://cloudspec.aliyun-inc.com/?q=R-BR-0005
事件等级：ERROR
事件信息：资源-主键-不允许主键形式发生变更
根组件ID：alicloud.test#Ra
根组件类型：resource
位置信息：文件路径：/Users/lazy/code/cloudspec-cli/idl_bc/test_master/main.cspec，行号：11，列号：9
约束器名称：alicloud.test#R-BR-0005，规范链接：http://cloudspec.aliyun-inc.com/?q=R-BR-0005
事件等级：ERROR
事件信息：资源-主键-不允许主键形式发生变更
根组件ID：alicloud.test#Rb
根组件类型：resource
位置信息：文件路径：/Users/lazy/code/cloudspec-cli/idl_bc/test_master/main.cspec，行号：15，列号：9
约束器名称：alicloud.test#R-BR-0067，规范链接：http://cloudspec.aliyun-inc.com/?q=R-BR-0067
事件等级：ERROR
事件信息：API:opC 出参删除了属性: $.pendingRemove
根组件ID：alicloud.test#opC
根组件类型：operation
位置信息：文件路径：/Users/lazy/code/cloudspec-cli/idl_bc/test_master/main.cspec，行号：19，列号：10
约束器名称：alicloud.test#R-AC-0012，规范链接：http://cloudspec.aliyun-inc.com/?q=R-AC-0012
事件等级：DANGER
事件信息：不允许删除返回参数。
根组件ID：alicloud.test#opC
根组件类型：operation
位置信息：文件路径：/Users/lazy/code/cloudspec-cli/idl_bc/test_master/main.cspec，行号：31，列号：2

```

针对示例：

+ test 路径是当前最新的IDL项目；
+ /Users/lazy/code/cloudspec-cli/idl_bc/test_master 是master的仓库地址。

#### auto 自动化操作
自动化生成操作，支持参数：

| 参数 | 说明 |
| --- | --- |
| --auto_type | 自动化生成的类型，默认值为autoGenerateOperations，支持类型如下：  + autoGenerateOperations 在资源上打上@[autoGenerateOperations](https://aliyuque.antfin.com/cloudspec/model/gpceg4uzvr9ck1ht#MK00z) annotate自动化生成对应的operation。  + autoGenerateResources 该命令用于发起api资源化自动推导异步任务，基于当前空间下的未资源化operation生成资源分组建议，如果分组建议成功生成，在operation上打上@autoGenerateResource annotate，显示该API所属的资源以及操作类型  + autoGenerateResourcesFinal 根据@autoGenerateResource annotate 指定的资源名称、operation类型，生成resource |


自动化生成operation示例：

```bash
(base) lazy@B-73KY41WF-1106 auto % cloudspec auto 
自动生成成功
时间：9ms
写入文件成功：/private/tmp/auto/operations/TagResources.cspec
写入文件成功：/private/tmp/auto/main.cspec
写入文件成功：/private/tmp/auto/operations/ListTagResources.cspec
写入文件成功：/private/tmp/auto/operations/UntagResources.cspec
写入文件成功：/private/tmp/auto/resources/Instance.cspec
(base) lazy@B-73KY41WF-1106 auto % /Users/lazy/code/cloudspec-cli/dist/cloudspec/cloudspec build
解析成功
解析时间：113ms
```



自动化生成resources示例：

```bash
linjun@alideMacBook-Pro prject_2022-02-02 % cloudspec auto --auto_type=autoGenerateResources 
推导任务进行中，当前进度为0% 
```

触发任务后，等待若干时间后，出现如下输出，则说明资源分组建议完成。

```bash
写入文件成功：/prject_2022-02-02/operations/UpdateInstance.cspec
写入文件成功：/prject_2022-02-02/operations/DeleteInstance.cspec
写入文件成功：/prject_2022-02-02/operations/CreateInstance.cspec
写入文件成功：/prject_2022-02-02/operations/GetInstance.cspec
写入文件成功：/prject_2022-02-02/operations/ListInstances.cspec
```

搜索代码关键字“@autoGenerateResource ”，可以找到各个未资源化的API的资源分组建议，如下图所示：

![](https://intranetproxy.alipay.com/skylark/lark/0/2024/png/123756397/1729479017986-29079d8b-477c-4935-bfcf-47834df1cc55.png)

在当前版本中，@autoGenerateResource注解中的isAssociate字段未生效，主要关注resourceName和type字段，开发者可以对这两个字段进行修改，手动变更自动生成提供的资源名称及操作类型。

在手动调整完毕后，执行autoGenerateResourcesFinal命令，可以获取到资源的最终生成结果。同时可以看到资源生成的摘要信息，包括参与推导的API总数、新生成的资源数、新生成的资源列表、被修改的资源数、被修改的资源列表、被修改的资源列表中加入API的详情、API数量及类型不足而无法生成的资源分组、不允许被附加的API列表、无法形成资源属性的API分组。

```bash
linjun@alideMacBook-Pro prject_2022-02-02 % cloudspec auto --auto_type=autoGenerateResourcesFinal
推导的API总数: 81                                                               
新生成的资源数: 10
新生成的资源列表:
["RamPolicyExportTask","ProjectBuild","RamPolicy","TaskPolicy","RegistryNamespace","ResourceTypeExample","RegistryModule","Resource","RamPolicyExportTaskVersion","BuildResult"]
被修改的资源数: 5
被修改的资源列表:
["Group","Task","SceneTestingTask","ResourceExportTask","Module"]
被修改的资源列表中加入API的详情:
{"SceneTestingTask":["CloneSceneTestingTask","ListRelationSceneTestingTasks"],"Group":["CloneGroup","DissociateGroupRelation"],"Task":["CloneTask","GetTaskParameter"],"ResourceExportTask":["CancelResourceExportTask","RemoveResourceExportTaskVersion","ExecuteResourceExportTask","ListResourceExportTaskVersions"],"Module":["CloneModule","GenerateModule","ListRelationModules","UploadModule"]}
API数量及类型不足无法生成资源的分组:
{"RelationTask":["ListRelationTasks"],"SessionChat":["ListSessionChats","SubmitSessionChat"],"JobResult":["SyncJobResults"],"Product":["ListProducts"],"RegistryModuleVersion":["ListRegistryModuleVersions","PublishRegistryModuleVersion"],"ResourceType":["ListResourceTypesAbilities"],"ParameterSetRelation":["ListParameterSetRelation"],"SharedAccount":["AddSharedAccounts","RemoveSharedAccounts"],"ResourceTypeAbility":["SyncResourceTypeAbility"],"ExplorerHistory":["CreateExplorerHistory","ListExplorerHistories"],"GroupRelation":["ListGroupRelation"],"ExportTask":["SyncExportTaskResults"],"ChatFeedback":["SubmitSessionChatFeedback"],"RabbitmqPublisherAttachment":["AttachRabbitmqPublisher","DetachRabbitmqPublisher"],"TerraformerVersion":["ListTerraformerVersions"],"TerraformVersion":["ListAvailableTerraformVersions"],"GroupRelationAttachment":["AssociateGroupRelation"],"TerraformProviderVersion":["ListTerraformProviderVersions"],"Session":["CreateSession","ListSessions"]}
不允许被附加的API列表:
{}
无法形成资源属性的API分组:
{}
```

#### yaml 生成yaml
将IDL构建为镇元支持的yaml，支持参数：

| 参数 | 说明 |
| --- | --- |
| --output_path | 指定yaml输出的路径，需要为全路径，默认为/tmp/cloudspec |


示例：

```bash
(base) lazy@B-73KY41WF-1106 auto % cloudspec yaml -o /tmp/test
生成成功
时间：104ms
生成的文件在路径：/tmp/test
写入文件成功：/tmp/test/resources/Instance.yaml
写入文件成功：/tmp/test/apis/TagResources.yaml
写入文件成功：/tmp/test/apis/UntagResources.yaml
写入文件成功：/tmp/test/apis/CreateInstance.yaml
写入文件成功：/tmp/test/apis/DeleteInstance.yaml
写入文件成功：/tmp/test/apis/GetInstance.yaml
写入文件成功：/tmp/test/apis/UpdateInstance.yaml
写入文件成功：/tmp/test/apis/ListInstances.yaml
写入文件成功：/tmp/test/apis/ListTagResources.yaml
写入文件成功：/tmp/test/ram/ram.yaml
写入文件成功：/tmp/test/openStruct.yaml
(base) lazy@B-73KY41WF-1106 auto % ls -al /tmp/test 
total 8
drwxr-xr-x   6 lazy  wheel  192  1 14 20:37 .
drwxrwxrwt  17 root  wheel  544  1 14 20:37 ..
drwxr-xr-x  10 lazy  wheel  320  1 14 20:37 apis
-rw-r--r--   1 lazy  wheel   26  1 14 20:37 openStruct.yaml
drwxr-xr-x   3 lazy  wheel   96  1 14 20:37 ram
drwxr-xr-x   3 lazy  wheel   96  1 14 20:37 resources
```

#### terraform 生成 Terraform 代码 
:::warning
该特性需要将 CLI 升级到至少 1.1.3 版本。

:::

在生成 Terraform 代码前，你需要将 Terraform provider 的源码 clone 到本地的一个路径，[https://github.com/aliyun/terraform-provider-alicloud](https://github.com/aliyun/terraform-provider-alicloud)。

支持的参数如下：

| 参数 | 说明 |
| --- | --- |
| -r | 指定待生成的资源名称，例如 Instance |
| -o | Terraform provider 源码的路径 |


示例：

```bash
lazy@B-73KY41WF-1106 IaCService_pop_IaCService_2021-08-06 % cloudspec terraform -o /Users/lazy/github/terraform-provider-alicloud -r Project
开始生成terraform代码，请稍等
Running generate terraform command with path: /Users/lazy/code-idl/IaCService_pop_IaCService_2021-08-06
Running validate command with path: /Users/lazy/code-idl/IaCService_pop_IaCService_2021-08-06
判断入口文件为：main.cspec
编译成功
File created at: /Users/lazy/github/terraform-provider-alicloud/alicloud/connectivity/client.go
File created at: /Users/lazy/github/terraform-provider-alicloud/alicloud/provider.go
File created at: /Users/lazy/github/terraform-provider-alicloud/alicloud/service_alicloud_iac_service_v2.go
File created at: /Users/lazy/github/terraform-provider-alicloud/alicloud/resource_alicloud_iac_service_project_test.go
File created at: /Users/lazy/github/terraform-provider-alicloud/alicloud/resource_alicloud_iac_service_project.go
```

从 Terraform 源码的路径可以看到：

```bash
lazy@B-73KY41WF-1106 terraform-provider-alicloud % git status
On branch init
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
	modified:   alicloud/connectivity/client.go
	modified:   alicloud/provider.go

Untracked files:
  (use "git add <file>..." to include in what will be committed)
	alicloud/resource_alicloud_iac_service_project.go
	alicloud/resource_alicloud_iac_service_project_test.go
	alicloud/service_alicloud_iac_service_v2.go

no changes added to commit (use "git add" and/or "git commit -a")
```

#### graph 生成 服务/资源/API 可视化图
:::warning
该特性需要将 CLI 升级到至少 1.1.5 版本。

:::

如果需要生成 XMind 格式图，请确保本地安装了 Xmind（阿里郎可直接安装）；如果需要生成 png 格式的图片，请确保本地安装了Graphviz，通过 brew安装命令：

```bash
brew install graphviz
```

支持的参数如下：

| 参数 | 说明 |
| --- | --- |
| -n | 必填，指定待生成的服务、资源或者 operation名称，例如 Instance |
| -t |  生成图分类，支持  + xmind 生成 xmind 格式图，默认值  + png，生成 png格式图片 |
| -o | 生成图片存储的位置，默认为 ~/Downloads，请指定全路径。 |


##### 生成服务全景图示例
生成IaCService服务的 Xmind 图：(-n 指定服务名称，IDL 中的 service 定义的名称)

```bash
lazy@B-73KY41WF-1106 IaCService_pop_IaCService_2021-08-06 % cloudspec graph -n IaCService
Running generate graph command with path: /Users/lazy/code-idl/IaCService_pop_IaCService_2021-08-06 and outputDir: /Users/lazy/Downloads
Running validate command with path: /Users/lazy/code-idl/IaCService_pop_IaCService_2021-08-06
判断入口文件为：main.cspec
编译成功
CloudSpec 解析成功，开始生成图像
开始生成 XMind 文件：/Users/lazy/Downloads/IaCService.xmind
生成 XMind 文件成功：/Users/lazy/Downloads/IaCService.xmind
Process completed successfully.
```

![](https://intranetproxy.alipay.com/skylark/lark/0/2024/png/309278/1730967301008-def09f31-bab8-477b-9fc8-455c889b0151.png)

##### 生成资源全景图示例
生成资源的 XMind 格式图：

```bash
lazy@B-73KY41WF-1106 private-link % cloudspec graph -n VpcEndpoint
Running generate graph command with path: /Users/lazy/code-idl/private-link and outputDir: /Users/lazy/Downloads
Running validate command with path: /Users/lazy/code-idl/private-link
判断入口文件为：main.cspec
编译成功
CloudSpec 解析成功，开始生成图像
开始生成 XMind 文件：/Users/lazy/Downloads/VpcEndpoint.xmind
生成 XMind 文件成功：/Users/lazy/Downloads/VpcEndpoint.xmind
```

![](https://intranetproxy.alipay.com/skylark/lark/0/2024/png/309278/1730884366801-d4b87cdc-8e7f-4acd-8b22-3278d65a47f3.png)

> 标记为红色节点的属性为错误节点，需要修正。
>

生成 png 格式的图片：

```bash
lazy@B-73KY41WF-1106 private-link % cloudspec graph -n VpcEndpoint -t png
Running generate graph command with path: /Users/lazy/code-idl/private-link and outputDir: /Users/lazy/Downloads
Running validate command with path: /Users/lazy/code-idl/private-link
判断入口文件为：main.cspec
编译成功
CloudSpec 解析成功，开始生成图像
生成 dot 文件成功：/Users/lazy/Downloads/VpcEndpoint.dot
开始生成 PNG 图像：/Users/lazy/Downloads/VpcEndpoint.png
生成 PNG 图像成功：/Users/lazy/Downloads/VpcEndpoint.png
```

![](https://intranetproxy.alipay.com/skylark/lark/0/2024/png/309278/1730884446729-78e0bd13-1669-46ca-a5ec-0a392c7a5b0f.png)

##### 生成 operation 全景图示例
生成 operation CreateModule示例：

```bash
lazy@B-73KY41WF-1106 IaCService_pop_IaCService_2021-08-06 % cloudspec graph -n CreateModule
Running generate graph command with path: /Users/lazy/code-idl/IaCService_pop_IaCService_2021-08-06 and outputDir: /Users/lazy/Downloads
Running validate command with path: /Users/lazy/code-idl/IaCService_pop_IaCService_2021-08-06
判断入口文件为：main.cspec
编译成功
CloudSpec 解析成功，开始生成图像
开始生成 XMind 文件：/Users/lazy/Downloads/CreateModule.xmind
生成 XMind 文件成功：/Users/lazy/Downloads/CreateModule.xmind
Process completed successfully.
```

![](https://intranetproxy.alipay.com/skylark/lark/0/2024/png/309278/1730967393840-8b9718ec-ac67-4bb8-8697-595d5daea365.png)



#### tag 生成 TAG 的 outer/inner API
:::warning
使用此功能，CLI 的版本至少为 1.1.6。安装及更新 CLI：[https://aliyuque.antfin.com/cloudspec/model/cli-install](https://aliyuque.antfin.com/cloudspec/model/cli-install)

:::

生成接入 TAG 服务的 inner、outer API，支持参数：

| 参数 | 说明 |
| --- | --- |
| -t |  生成的 API 类型，支持：  + inner 生成内部 API；  + outer 生成外部（对客户）API。**（默认值）** |


#### sync-doc 从 API 文档完善资源属性文档
:::warning
[该特性需要将 CLI 升级到至少 1.1.7 版本](https://aliyuque.antfin.com/cloudspec/model/cli-install)。

:::

指定资源，根据其关联的 OpenAPI 的字段，自动推断资源属性的文档，请确保 OpenAPI 的文档都已经正确编写，且完成翻译后使用。

支持的参数如下：

| 参数 | 说明 |
| --- | --- |
| -n | 指定待同步的资源名称，例如 Instance； |
| -l |  是否使用大模型推断，默认为 false，支持 true/false。 |


示例：

```bash
lazy@B-73KY41WF-1106 gwlb % cloudspec sync-doc -n Listener
开始推断资源文档，请稍等
Running fix resource command with path: /Users/lazy/code-idl/gwlb
Running validate command with path: /Users/lazy/code-idl/gwlb
判断入口文件为：main.cspec
编译成功
推断资源属性：$.ListenerId文档开始
推断资源属性：$.ListenerId文档成功
推断资源属性：$.All文档开始
推断资源属性：$.All文档成功
推断资源属性：$.Status文档开始
推断资源属性：$.Status文档成功
推断资源属性：$.NextToken文档开始
推断资源属性：$.NextToken文档成功
推断资源属性：$.MaxResults文档开始
推断资源属性：$.MaxResults文档成功
推断资源属性：$.DryRun文档开始
推断资源属性：$.DryRun文档成功
推断资源属性：$.ServerGroupId文档开始
推断资源属性：$.ServerGroupId文档成功
推断资源属性：$.LoadBalancerId文档开始
推断资源属性：$.LoadBalancerId文档成功
推断资源属性：$.ResourceType文档开始
推断资源属性：$.ResourceType文档成功
推断资源属性：$.ListenerDescription文档开始
推断资源属性：$.ListenerDescription文档成功
推断资源属性：$.Filter[*].Name文档开始
推断资源属性：$.Filter[*].Name文档成功
推断资源属性：$.Filter[*].Values[*]文档开始
推断资源属性：$.Filter[*].Values[*]文档成功
推断资源属性：$.Filter[*].Values文档开始
推断资源属性：$.Filter[*].Values文档成功
推断资源属性：$.Filter[*]文档开始
推断资源属性：$.Filter[*]文档成功
推断资源属性：$.Filter文档开始
推断资源属性：$.Filter文档成功
推断资源属性：$.Skip文档开始
推断资源属性：$.Skip文档成功
推断资源属性：$.RegionId文档开始
推断资源属性：$.RegionId文档成功
推断资源属性：$.ClientToken文档开始
推断资源属性：$.ClientToken文档成功
推断资源属性：$.Tags[*].TagKey文档开始
推断资源属性：$.Tags[*].TagKey文档成功
推断资源属性：$.Tags[*].TagValue文档开始
推断资源属性：$.Tags[*].TagValue文档成功
推断资源属性：$.Tags[*]文档开始
推断资源属性：$.Tags[*]文档成功
推断资源属性：$.Tags文档开始
推断资源属性：$.Tags文档成功
重新生成受影响的文件
重新生成受影响的文件：/Users/lazy/code-idl/gwlb/resources/Listener.cspec成功
处理完成
```

#### test 资源及operation测试
cloudspec支持命令 test用于资源、operation 的测试，其子命令如下：

```bash
lazy@B-73KY41WF-1106 esa % cloudspec test
Usage: cloudspec test [OPTIONS] COMMAND [ARGS]...

  测试自动生成及运行.

Options:
  --help  Show this message and exit.

Commands:
  coverage  展示测试用例覆盖率
  gen       自动生成测试用例
  json      生成测试用例json格式
  migrate   迁移镇元测试用例json到 IDL
  run       执行测试用例
```

##### gen 测试用例自动生成
```bash
lazy@B-73KY41WF-1106 esa % cloudspec test gen --help
Usage: cloudspec test gen [OPTIONS] [PATH]

  自动生成测试用例

Options:
  -c, --component TEXT  资源或者API名称  [required]
  --help                Show this message and exit.
```

您需要在 IDL 目录root 目录下，执行命令：cloudspec test gen -c XXX，XXX 是资源名称，示例如下：

```bash
lazy@B-73KY41WF-1106 test-gen % cloudspec test gen -c Test
开始生成测试用例，请稍等
Running generate test command with path: /private/tmp/test-gen and component: Test
Running validate command with path: /private/tmp/test-gen
判断入口文件为：main.cspec
编译成功
编译耗时：170ms
12月 06, 2024 8:09:50 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 开始生成测试用例：resource_Test_test，资源：Test，阶段：init
12月 06, 2024 8:09:50 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate init
信息: 生成 init 阶段属性：Name
WARNING: sun.reflect.Reflection.getCallerClass is not supported. This will impact performance.
12月 06, 2024 8:09:52 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate init
信息: 生成 init 阶段属性：Other
12月 06, 2024 8:09:53 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate init
信息: 生成 init 阶段属性：RefId
12月 06, 2024 8:09:53 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 开始生成测试用例：resource_Ref_test，资源：Ref，阶段：init
12月 06, 2024 8:09:53 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate init
信息: 生成 init 阶段属性：Name
12月 06, 2024 8:09:54 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 生成测试用例：resource_Ref_test，资源：Ref，阶段：init 完成
12月 06, 2024 8:09:54 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 开始生成测试用例：resource_Ref_test，资源：Ref，阶段：destroy
12月 06, 2024 8:09:54 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 生成测试用例：resource_Ref_test，资源：Ref，阶段：destroy 完成
12月 06, 2024 8:09:54 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate init
信息: 生成 init 阶段属性：Tags
12月 06, 2024 8:09:54 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 生成测试用例：resource_Test_test，资源：Test，阶段：init 完成
12月 06, 2024 8:09:54 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 开始生成测试用例：resource_Test_test，资源：Test，阶段：modifies
12月 06, 2024 8:09:54 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段(单步修改)属性：Name
12月 06, 2024 8:09:55 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段(单步修改)属性：Age
12月 06, 2024 8:09:56 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段(单步修改)属性：Interests
12月 06, 2024 8:09:59 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段(单步修改)属性：config
12月 06, 2024 8:10:00 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段（单步修改）属性：$.Tags
12月 06, 2024 8:10:00 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段（全量替换）属性：Name
12月 06, 2024 8:10:01 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段（全量替换）属性：Age
12月 06, 2024 8:10:02 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段（全量替换）属性：Interests
12月 06, 2024 8:10:05 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段（全量替换）属性：config
12月 06, 2024 8:10:07 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段（全量替换）属性：$.Tags
12月 06, 2024 8:10:07 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate modifies
信息: 生成 modifies 阶段, 单独清空属性：$.Tags
12月 06, 2024 8:10:07 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 生成测试用例：resource_Test_test，资源：Test，阶段：modifies 完成
12月 06, 2024 8:10:07 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 开始生成测试用例：resource_Test_test，资源：Test，阶段：destroy
12月 06, 2024 8:10:07 下午 com.aliyun.cloudspec.automatic.impl.testCase.ResourceTestCaseGenerate generate
信息: 生成测试用例：resource_Test_test，资源：Test，阶段：destroy 完成
测试用例生成成功，开始更新本地文件
重新生成受影响的文件
重新生成受影响的文件：/private/tmp/test-gen/tests/Test_test.cspec成功
重新生成受影响的文件：/private/tmp/test-gen/main.cspec成功
```

```yaml
$version: 1
namespace: a.b.c.v20211012

$test Ref resource_Ref_test {
  init: {
    Name: "referencename"
  }
}

@testConfig({
  main: true
  uses: [resource_Ref_test]
})
$test Test resource_Test_test {
  init: {
    Name: "exampleName"
    Other: {
      OtherA: "otherAValue"
      OtherB: 100
    }
    RefId: "{{$.resource_Ref_test.Id}}"
    Tags: [{
      TagKey: "key1"
      TagValue: "value1"
    }
    , {
      TagKey: "key2"
      TagValue: "value2"
    }
    , {
      TagKey: "key3"
      TagValue: "value3"
    }]
  }
  modifies: [{
    Name: "newName"
  }
  , {
    Age: 25
  }
  , {
    Interests: [{
      InterestName: "reading"
      InterestLevel: 1
    }
    , {
      InterestName: "traveling"
      InterestLevel: 2
    }]
  }
  , {
    config: {
      host: "localhost"
      port: 3306
      user: "root"
      password: "password"
      database: "testdb"
    }
  }
  , {
    Tags: [{
      TagKey: "key4"
      TagValue: "value4"
    }
    , {
      TagKey: "key5"
      TagValue: "value5"
    }]
  }
  , {
    Name: "uniqueName"
    Age: 28
    Interests: [{
      InterestName: "swimming"
      InterestLevel: 3
    }
    , {
      InterestName: "cooking"
      InterestLevel: 4
    }]
    config: {
      host: "127.0.0.1"
      port: 3307
      user: "admin"
      password: "securepassword"
      database: "productiondb"
    }
    Tags: [{
      TagKey: "key6"
      TagValue: "value6"
    }]
  }
  , {
    Tags: []
  }]
  destroy: {
  }
}
```

```yaml
// 警告：请在明确知晓如何执行测试的前提下修改已有的Test Case，否则可能导致已有的Test Case无法执行
// Create at: 2024/11/26/17:28
// Author: lazy
$version: 1
namespace: a.b.c.v20211012

import "tests/Test_test.cspec"

@runtimeType("pop")
service c {
  version: "2021-10-12"
  resources: [Ref, Test]
}

@resourceBaseInfo({
  classification: "normal"
})
resource Ref {
  identifyDefinition: {
    /// 主键
    Id: string
  }
  properties: {
    /// 引用的名称，最长 32 个字符，a-z组成
    @required
    Name: string
  }
  create: createRef
  get: getRef
}

@autoMapping
operation createRef {
  input: {
    @required
    Name: string
  }
  output: {
    Id: string
  }
  errors: []
}

@autoMapping
operation getRef {
  input: {
    Id: string
  }
  output: {
    Name: string
  }
  errors: []
}

struct TestProp {
  /// 名称
  @required
  Name: string
  /// 年龄
  Age: int64
  /// 其他配置
  @required
  Other: {
    /// 其他A
    OtherA: string
    /// 其他B
    OtherB: int64
  }
  /// 用户兴趣
  Interests: array<{
    /// 兴趣名称
    InterestName: string
    /// 兴趣等级
    InterestLevel: int64
  }
  >
  // 只能从 update 修改，判断的时候也必须包含
  /// 标签
  Tags: array<{
    /// 标签键
    TagKey: string
    /// 标签值
    TagValue: string
  }
  >
  /// 引用ID
  @required
  RefId: string
  /// MySQL 数据库配置
  config: map<any>
}

// 自动生成 test 的测试用例
@resourceBaseInfo({
  classification: "normal"
})
@references([{
  relatedResource: Ref
  localProperty: "$.RefId"
  remoteProperty: "$.Id"
}])
resource Test {
  identifyDefinition: {
    Id: string
  }
  properties: TestProp
  create: createTest
  update: updateTest
  delete: deleteTest
  get: getTest
  list: listTests
}

operation deleteTest {
  input: {
    Id: string
  }
  output: {
    Id: string
  }
  errors: []
}

@autoMapping
operation listTests {
  input: {
    Name: string
    Age: int64
  }
  output: {
    @nested
    datas: array<TestProp>
  }
  errors: []
}

@autoMapping
operation updateTest {
  input: {
    Id: string
    Name: string
    Age: int64
    Tags: array<{
      TagKey: string
      TagValue: string
    }
    >
    Interests: array<{
      /// 兴趣名称
      InterestName: string
      /// 兴趣等级
      InterestLevel: int64
    }
    >
    config: map<any>
  }
  output: {
    Id: string
  }
  errors: []
}

@autoMapping
operation createTest {
  input: {
    Name: string
    Age: int64
    Other: {
      OtherA: string
      OtherB: int64
    }
    @required
    RefId: string
    Interests: array<{
      /// 兴趣名称
      InterestName: string
      /// 兴趣等级
      InterestLevel: int64
    }
    >
    config: map<any>
  }
  output: {
    Id: string
  }
  errors: []
}

@autoMapping
operation getTest {
  input: {
    Id: string
  }
  output: TestProp
  errors: []
}
```

##### run 执行测试用例
```bash
lazy@B-73KY41WF-1106 test-gen % cloudspec test run --help
Usage: cloudspec test run [OPTIONS] [PATH]

  执行测试用例

Options:
  -n, --name TEXT       测试用例名称
  -c, --component TEXT  资源或者API名称
  --help                Show this message and exit.
```

您可以用 -n 指定单个用例执行，也可以指定 -c 同时运行一个资源下的全部用例，示例:

```bash
cloudspec test run -c Record
```

```yaml
$version: 1
namespace: alicloud.ESA.ESA.v20240910

@testConfig({
  execConfigUuid: "bbeceb14f********0bfd0f"
  runtime: {
    regionId: "cn-hangzhou"
    endpoint: "esa.[RegionId].aliyuncs.com"
    // 账号不允许覆盖
    accountOverrideSupport: false
  }
})
$test Record record_pre_smtcdn {
  init: {
    SiteId: 495685814175328
    RecordName: "music.smtcdn.com"
    RecordType: "CNAME"
    Data: {
      Value: "music.aliyun.com"
    }
    Ttl: 1
    Comment: "test"
    Proxied: true
    BizName: "web"
  }
  modifies: [{
    Comment: "test1"
  }, {
    Comment: "test2"
  }]
  destroy: {}
}

@testConfig({
  //execConfigUuid: "bac6d5428704465*******4e7a6d2"
  main: true
  runtime: {
    regionId: "cn-hangzhou"
    endpoint: "esa.[RegionId].aliyuncs.com"
  }
  uses: [record_pre_smtcdn]
})
$test Record record_test {
  init: {
    SiteId: 488793825298016
    RecordName: "music.changes.com.cn"
    RecordType: "CNAME"
    Data: {
      Value: "music.aliyun.com"
    }
    Ttl: "{{$.record_pre_smtcdn.Ttl}}"
    Comment: "test"
    Proxied: true
    BizName: "{{$.record_pre_smtcdn.BizName}}"
  }
  modifies: [{
    Comment: "test11"
  }, {
    Comment: "test22"
  }]
  destroy: {}
}
```

执行后，会在本地自动打开默认浏览器，展示测试的进度、明细及测试覆盖度。

![](https://intranetproxy.alipay.com/skylark/lark/0/2024/png/309278/1733487549632-254dee8e-4cfa-4146-9ca2-0704323867be.png)

执行成功时会展示测试覆盖度，点击每个步骤的数字即可查询：

![](https://intranetproxy.alipay.com/skylark/lark/0/2024/png/309278/1733487578483-6813cb80-4cfd-4b92-8f12-d8403ca87353.png)

##### coverage 计算静态测试覆盖度
```bash
lazy@B-73KY41WF-1106 esa % cloudspec test coverage --help
当前版本是最新版本
Usage: cloudspec test coverage [OPTIONS] [PATH]

  展示测试用例覆盖率

Options:
  -c, --component TEXT  资源或者API名称
  -n, --name TEXT       测试用例名称
  --help                Show this message and exit.
```

您可以指定 -n 计算单个测试用例的覆盖度，也可以 -c 指定某个资源汇总计算某个资源全部用例的测试覆盖度，示例如下：

```bash
lazy@B-73KY41WF-1106 esa % cloudspec test coverage -c Record
开始计算测试用例覆盖率，请稍等
Running generate test coverage command with path: /Users/lazy/code-idl/esa
Running validate command with path: /Users/lazy/code-idl/esa
判断入口文件为：main.cspec
编译成功
编译耗时：4076ms
包含测试用例：record_test
当前测试覆盖度为静态计算得出，仅能评估测试设计覆盖情况，需要确保测试用例均运行通过。
测试覆盖度为：0.173611
属性覆盖度：
┌────────────────────┬─────┬──────────────┬───────────────┬───────────────────────┬───────────────────┐
│     Name Path      │Score│Write Coverage│Update Coverage│Array Data Update Score│Enum Coverage Score│
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│     $.Comment      │1.00 │     1.00     │     1.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│       $.Data       │0.25 │     0.50     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│ $.RecordMatchType  │1.00 │     1.00     │      N/A      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│      $.SiteId      │1.00 │     1.00     │      N/A      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.Data.Value    │0.25 │     0.50     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│$.Data.MatchingType │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│     $.AuthConf     │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│ $.Data.Fingerprint │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│   $.RecordCname    │0.00 │     0.00     │      N/A      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│$.AuthConf.AccessKey│0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.Data.Usage    │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│  $.Data.Algorithm  │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.Data.Type     │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.CreateTime    │0.00 │     0.00     │      N/A      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│ $.AuthConf.Version │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│       $.Ttl        │0.25 │     0.50     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│   $.Data.KeyTag    │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.RecordName    │1.00 │     1.00     │      N/A      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.HostPolicy    │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│ $.Data.Certificate │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.Data.Port     │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.UpdateTime    │0.00 │     0.00     │      N/A      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.RecordType    │1.00 │     1.00     │      N/A      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│$.AuthConf.AuthType │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│   $.Data.Weight    │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│  $.Data.Priority   │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│  $.Data.Selector   │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│     $.SiteName     │0.00 │     0.00     │      N/A      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│     $.Data.Tag     │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│     $.BizName      │0.25 │     0.50     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│$.AuthConf.SecretKey│0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│     $.Proxied      │0.25 │     0.50     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│ $.AuthConf.Region  │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.SourceType    │0.00 │     0.00     │      N/A      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│    $.Data.Flag     │0.00 │     0.00     │     0.00      │          N/A          │        N/A        │
├────────────────────┼─────┼──────────────┼───────────────┼───────────────────────┼───────────────────┤
│ $.RecordSourceType │0.00 │     0.00     │      N/A      │          N/A          │        N/A        │
└────────────────────┴─────┴──────────────┴───────────────┴───────────────────────┴───────────────────┘
```

##### json 转换为镇元支持的 JSON 格式
```bash
lazy@B-73KY41WF-1106 esa % cloudspec test json --help
Usage: cloudspec test json [OPTIONS] [PATH]

  生成测试用例json格式

Options:
  -c, --component TEXT  资源或者API名称
  -n, --name TEXT       测试用例名称
  -o, --output TEXT     导出时输出路径
  --help                Show this message and exit.
```

示例：

```bash
lazy@B-73KY41WF-1106 esa % cloudspec test json -n record_test
开始生成测试用例json，请稍等
Running generate test yaml command with path: /Users/lazy/code-idl/esa
Running validate command with path: /Users/lazy/code-idl/esa
判断入口文件为：main.cspec
编译成功
编译耗时：3858ms
路径为空，使用默认路径：/Users/lazy/Downloads
生成测试用例record_testjson成功，生成文件路径：/Users/lazy/Downloads/record_test.json
```

```json
{
	"caseOperations":{
		"resourceAttributes":[
			{
				"majorOperation":false,
				"ref":false,
				"resourceAttributeInJson":"{\"RecordName\":\"music.changes.com.cn\",\"Comment\":\"test\",\"SiteId\":488793825298016,\"Proxied\":true,\"RecordType\":\"CNAME\",\"Data\":{\"Value\":\"music.alyun.com\"},\"BizName\":\"${{ref(resource,ESA::Record::0.0.0.0.pre::record_pre_smtcdn.BizName)}}\",\"Ttl\":\"${{ref(resource,ESA::Record::0.0.0.0.pre::record_pre_smtcdn.Ttl)}}\"}",
				"resourceAttributes":[
					{
						"name":"SiteId",
						"type":"number",
						"value":488793825298016
					},
					{
						"name":"RecordName",
						"type":"string",
						"value":"music.changes.com.cn"
					},
					{
						"name":"RecordType",
						"type":"string",
						"value":"CNAME"
					},
					{
						"items":[
							{
								"name":"Value",
								"type":"string",
								"value":"music.alyun.com"
							}
						],
						"name":"Data",
						"type":"object"
					},
					{
						"name":"Ttl",
						"type":"ref",
						"value":"${{ref(resource,ESA::Record::0.0.0.0.pre::record_pre_smtcdn.Ttl)}}"
					},
					{
						"name":"Comment",
						"type":"string",
						"value":"test"
					},
					{
						"name":"Proxied",
						"type":"boolean",
						"value":true
					},
					{
						"name":"BizName",
						"type":"ref",
						"value":"${{ref(resource,ESA::Record::0.0.0.0.pre::record_pre_smtcdn.BizName)}}"
					}
				],
				"stepAction":"Create",
				"stepIndex":0
			},
			{
				"majorOperation":false,
				"ref":false,
				"resourceAttributeInJson":"{\"Comment\":\"test11\"}",
				"resourceAttributes":[
					{
						"name":"Comment",
						"type":"string",
						"value":"test11"
					}
				],
				"stepAction":"Update",
				"stepIndex":1
			},
			{
				"majorOperation":false,
				"ref":false,
				"resourceAttributeInJson":"{\"Comment\":\"test22\"}",
				"resourceAttributes":[
					{
						"name":"Comment",
						"type":"string",
						"value":"test22"
					}
				],
				"stepAction":"Update",
				"stepIndex":2
			},
			{
				"majorOperation":false,
				"ref":false,
				"resourceAttributeInJson":"{}",
				"resourceAttributes":[],
				"stepAction":"List",
				"stepIndex":3
			}
		]
	},
	"preHookOperations":[
		{
			"clearResource":false,
			"operationSteps":[
				{
					"code":"Record",
					"forceCreate":false,
					"keepUpToDate":false,
					"name":"record_pre_smtcdn",
					"product":"ESA",
					"resourceAttributeInJson":"{\"RecordName\":\"music.smtcdn.com\",\"Comment\":\"test\",\"SiteId\":495685814175328,\"Proxied\":true,\"RecordType\":\"CNAME\",\"Data\":{\"Value\":\"music.alyun.com\"},\"BizName\":\"web\",\"Ttl\":1}",
					"stepIndex":0,
					"stepName":"record_pre_smtcdn",
					"stepType":"resource_type"
				}
			]
		}
	]
}
```

#### copy 快速拷贝资源、operation、数据结构
新拷贝的组件会独立防止，引用的数据除 Open Struct外都会重新命名。

:::warning
该特性需要将 CLI 升级到至少 1.1.9 版本。

:::

支持的参数如下：

| 参数 | 说明 |
| --- | --- |
| -s | 原组件的名称 |
| -t | 目标组件名称 |


示例：

```bash
lazy@B-73KY41WF-1106 merge2 % cloudspec copy -s Create03 -t Create03Other                                            
Running copy command with path: /private/tmp/merge2
Running validate command with path: /private/tmp/merge2
判断入口文件为：main.cspec
编译成功
编译耗时：2011ms
Copy run finish
重新生成受影响的文件：/private/tmp/merge2/main.cspec成功
重新生成受影响的文件：/private/tmp/merge2/operations/Create03Other.cspec成功
```

#### merge 合并指定分支
注意，merge 的本质是将当前分支强制 reset 到指定分支，然后将当前分支和待合并分支的 diff 以追加的方式重新写回，因此从 Git 的提交记录看会有一条 force 的记录。

:::warning
该特性需要将 CLI 升级到至少 1.1.9 版本。

:::

支持的参数如下：

| 参数 | 说明 |
| --- | --- |
| -b | 待合并的分支，不指定默认为 master |
| -o | 合并的选项，不指定默认为：ACCEPT_ADD，REBUILD_SERVICE，IGNORE_DELETE  需要多个值时请使用英文逗号隔开，支持的全部选项及含义：  + ACCEPT_ADD 允许当前分支的新增  + ACCEPT_DELETE 不出现在当前分支的 component 移除  + ACCEPT_MODIFY 如果当前分支和待合并同时修改了某个 component，以当前分支为准  + REBUILD_SERVICE 重新构建 main.cspec  + IGNORE_DELETE 忽略当前不在当前分支中，但在待合并分支的 component |


示例：

```bash
lazy@B-73KY41WF-1106 merge2 % cloudspec merge              
branch 未指定，默认合并 master 分支
Running validate command with path: /private/tmp/merge2
判断入口文件为：main.cspec
编译成功
编译耗时：1993ms
CloudSpec 解析成功
开始编译 master 分支
Running validate command with path: /tmp/d313bc1d-69a1-4087-b843-7a082305b07c/master
判断入口文件为：main.cspec
编译成功
编译耗时：1044ms
master 分支 CloudSpec 解析成功
开始解析公共祖先分支：339ba0c2cc38263962c45b84acd0476cdfa0d343
Running validate command with path: /tmp/d313bc1d-69a1-4087-b843-7a082305b07c/sameParentRevision
判断入口文件为：main.cspec
编译成功
编译耗时：836ms
公共祖先分支 CloudSpec 解析成功
开始比较 master 分支和公共祖先分支
重新生成受影响的文件：/private/tmp/merge2/operations/YgTest.cspec成功
重新生成受影响的文件：/private/tmp/merge2/models/Haha.cspec成功
重新生成受影响的文件：/private/tmp/merge2/operations/Create03Other.cspec成功
重新生成受影响的文件：/private/tmp/merge2/main.cspec成功
重新生成受影响的文件：/private/tmp/merge2/operations/YuangenExample.cspec成功
请手动执行清理工作，删除备份目录：/tmp/d313bc1d-69a1-4087-b843-7a082305b07c
^_^ 自动合并完成，请手工将改动 push 到远程仓库

```

#### re-balance  重新平衡 cloudspec 工程
执行重平衡操作后，数据结构会拆解放到 models 目录，和 resource 在一个文件的 operation 会被放置到 operations 目录。

:::warning
该特性需要将 CLI 升级到至少 1.1.9 版本。

:::

支持的参数如下：

| 参数 | 说明 |
| --- | --- |
| 无需参数 | |


示例：

```bash
lazy@B-73KY41WF-1106 re-balance % cloudspec re-balance
ReBalance start run
Running validate command with path: /private/tmp/re-balance
判断入口文件为：main.cspec
编译成功
编译耗时：129ms
ReBalance run finish
重新生成受影响的文件：/private/tmp/re-balance/main.cspec成功
重新生成受影响的文件：/private/tmp/re-balance/models/A.cspec成功
重新生成受影响的文件：/private/tmp/re-balance/models/B.cspec成功
```

#### mcp  生成 MCP Server
:::warning
该特性需要将 CLI 升级到至少 1.1.11 版本；

当前生成 MCP 能力需要先发布 Python v2的 SDK，否则会导致生成错误。

:::

##### init 初始化配置
你可以在指定目录下初始化配置，配置文件固定为aliyun-mcp-generate.yaml，支持参数如下：

```bash
Usage: cloudspec mcp init [OPTIONS] [PATH]

  初始化 MCP 相关的配置文件

Options:
  -c, --config-path TEXT  配置文件目录或者路径，默认是当前目录，如果是目录则会在目录下生成一个aliyun-mcp-
                          generate.yaml
  --help                  Show this message and exit.
```

不指定参数，会自动创建一个生成的模板，如：

```bash
(alibaba-cloud-mcp-server) lazy@B-73KY41WF-1106 test3 % cloudspec mcp init
MCP 配置文件写入成功，路径为：/private/tmp/test3/aliyun-mcp-generate.yaml
```

配置文件内容如下：

```yaml
mcpName: "alibaba-cloud-mcp-server"
mcpVersion: "1.0.0"
description: "Alibaba Cloud Mcp Server"
outputPath: "/tmp/alicloud-mcp-server"
apis:
- popCode: "Ecs"
  popVersion: "2014-05-26"
  apis:
  - "*Instance*"
- popCode: "Oss"
  popVersion: "2019-05-17"
  apis:
  - "*Bucket"
  - "*Buckets"
  - "*Object"
  - "*Objects"
  - "!SelectObject"
```

你可以手工编辑这个文件配置需要生成到 mcp server中的OpenAPI，也可以定点排除，支持的规则如下：

+ * 匹配全部 OpenAPI
+ *A 匹配以 A 结尾的 OpenAPI
+ A* 匹配以 A 开头的 OpenAPI
+ A 匹配名称为 A 的 OpenAPI
+ !*A 排除所有以 A 结尾的OpenAPI
+ !A 排除指定 OpenAPI A

> 请注意，排除规则最后执行。
>

##### generate 按配置文件生成 mcp server
生成 mcp server代码，当不指定具体配置文件时，默认探测当前目录下的aliyun-mcp-generate.yaml，支持参数：

```bash
(alibaba-cloud-mcp-server) lazy@B-73KY41WF-1106 test3 % cloudspec mcp generate --help 
Usage: cloudspec mcp generate [OPTIONS] [PATH]

  执行 MCP 生成操作

Options:
  -c, --config-path TEXT  配置文件目录或者路径，默认是当前目录，如果是目录，会探测aliyun-mcp-generate.yaml
  --help                  Show this message and exit.
```

你可以通过 -c 全路径指定一个非当前目录下的配置文件生成。

生成示例如下：

```bash
(alibaba-cloud-mcp-server) lazy@B-73KY41WF-1106 test3 % cloudspec mcp generate
读取配置文件路径：/private/tmp/test3/aliyun-mcp-generate.yaml
成功创建2个过滤器
MCP Server将生成到: /tmp/alicloud-mcp-server
WARNING: sun.reflect.Reflection.getCallerClass is not supported. This will impact performance.
生成MCP Server成功，输出文件数量：5
1: 写入文件成功：/tmp/alicloud-mcp-server/components.py
2: 写入文件成功：/tmp/alicloud-mcp-server/data.py
3: 写入文件成功：/tmp/alicloud-mcp-server/pyproject.toml
4: 写入文件成功：/tmp/alicloud-mcp-server/server.py
5: 写入文件成功：/tmp/alicloud-mcp-server/README.md
MCP Server生成完成，输出路径：/tmp/alicloud-mcp-server
```

###### 完成安装
如果你没有安装 uv，可以使用以下命名完成安装：

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

然后去generate 生成的代码目录执行：

```bash
cd /tmp/alicloud-mcp-server
uv sync 
```

然后激活.venv配置

```bash
lazy@B-73KY41WF-1106 alicloud-mcp-server % source .venv/bin/activate
```

命名行会加上前缀：

```bash
(alibaba-cloud-mcp-server) lazy@B-73KY41WF-1106 alicloud-mcp-server % 
```

###### 调试tools
启动 mcp 调试：

```bash
mcp dev server.py
```

可以看到启动一个本地的调试进程和网页：

```bash
(alibaba-cloud-mcp-server) lazy@B-73KY41WF-1106 alicloud-mcp-server % mcp dev server.py 
DEBUG:mcp.server.lowlevel.server:Initializing server 'alibaba-cloud-mcp-server'
DEBUG:mcp.server.lowlevel.server:Registering handler for ListToolsRequest
DEBUG:mcp.server.lowlevel.server:Registering handler for CallToolRequest
DEBUG:mcp.server.lowlevel.server:Registering handler for ListResourcesRequest
DEBUG:mcp.server.lowlevel.server:Registering handler for ReadResourceRequest
DEBUG:mcp.server.lowlevel.server:Registering handler for PromptListRequest
DEBUG:mcp.server.lowlevel.server:Registering handler for GetPromptRequest
DEBUG:mcp.server.lowlevel.server:Registering handler for ListResourceTemplatesRequest
DEBUG:asyncio:Using selector: KqueueSelector
DEBUG:tzlocal:/etc/localtime found
DEBUG:tzlocal:1 found:
 {'/etc/localtime is a symlink to': 'Asia/Shanghai'}
INFO:apscheduler.scheduler:Adding job tentatively -- it will be properly scheduled when the scheduler starts
INFO:apscheduler.scheduler:Added job "EcsRamRoleCredentialsProvider.__init__.<locals>.refresh_task" to job store "default"
INFO:apscheduler.scheduler:Scheduler started
DEBUG:apscheduler.scheduler:Looking for jobs to run
DEBUG:apscheduler.scheduler:Next wakeup is due at 2025-04-20 19:19:43.321129+08:00 (in 59.993282 seconds)
Starting MCP inspector...
⚙️ Proxy server listening on port 6277
🔍 MCP Inspector is up and running at http://127.0.0.1:6274 🚀

```

去浏览器打开这个页面：

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745147982216-e7af4324-570a-42d0-95f2-8cb0faf145ea.png)

点击 Connect，可从 Environment Variables 配置环境变量，支持：

| 环境变量 | 说明 |
| --- | --- |
| ALIBABA_CLOUD_ACCESS_KEY_ID | access key |
| ALIBABA_CLOUD_ACCESS_KEY_SECRET | access secret key |
| ALIBABA_CLOUD_REGION_ID | region配置 |


:::warning
建议使用阿里云 CLI ，如果本地安装了 CLI 且已经完成了profile 配置，会默认使用本地的 CLI 配置，无需设置以上环境变量。

:::

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745148067610-1d90d152-6441-4e6e-976d-17eeee46949e.png)

看到状态变成已连接，选中 tools 页面，点击 list tools：

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745148135992-9db25574-7b03-4623-9a98-fe960247400a.png)

到此即可开始调试生成的 tools，也可配置到支持 MCP 工具的大模型聊天工具中，配置的指令和调试页面左侧 run --with mcp mcp run server.py 一致。

###### 认证配置
MCP 中默认使用了阿里云 credentials 工具，会识别环境变量及 CLI 的 profile，如果你需要使用 AK/SK，需要配置环境变量：

```bash
ALIBABA_CLOUD_ACCESS_KEY_ID
ALIBABA_CLOUD_ACCESS_KEY_SECRET
# 非必须，使用 STS 时配置
ALIBABA_CLOUD_SECURITY_TOKEN
```

强烈建议使用 CLI 登录阿里云，然后默认使用 CLI 中配置的 profile，参考文档：[安装 CLI ](https://help.aliyun.com/zh/cli/installation-guide/?spm=a2c4g.11186623.help-menu-29991.d_2.683f76f1Gj15A0&scm=20140722.H_121988._.OR_help-T_cn~zh-V_1)| [配置 CLI](https://help.aliyun.com/zh/cli/configure-credentials?spm=a2c4g.11186623.help-menu-29991.d_3_0.14d033afshIjer&scm=20140722.H_121193._.OR_help-T_cn~zh-V_1)

> CloudSSO Mode CLI 刚支持，各语言库仍未更新，如果客户使用 CloudSSO Mode，暂时需要优先设置环境变量。
>

###### Cherry studio调试
建议使用[https://github.com/CherryHQ/cherry-studio/releases](https://github.com/CherryHQ/cherry-studio/releases)，下载后安装，配置百炼的 token：

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745241437381-6388c738-1136-4733-8dcc-80ecea040a21.png)

然后选择 mcp 配置：

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745241467758-c1ac6e05-5558-4ce0-a1bb-ab8d3ba2552d.png)

添加服务器，选择 stdio，命令配置如下：

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745241508633-7e8d7ce6-edf1-4950-8137-07e6ca37429e.png)

命令：

```bash
/tmp/alibabacloud-actiontrail-mcp-server/.venv/bin/mcp
```

> /tmp/alibabacloud-actiontrail-mcp-server 替换为上序配置中你定义的生成地址。
>

参数：

```bash
run
/tmp/alibabacloud-actiontrail-mcp-server/server.py
```

> /tmp/alibabacloud-actiontrail-mcp-server/ 替换为上序你配置的输出地址。
>

开始调试：

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745241589235-e0ddcb4f-dde6-4170-81c4-1f36acca02f7.png)

在聊天窗口直接配置 MCP 服务器为刚才添加的MCP server。

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745241630209-b994320f-77a8-41c2-9e8b-f599e918c420.png)

可以看到，大模型已经按 tools 的参数格式成功调用，并格式化展示了返回的数据。

###### CLINE
安装：[https://cline.bot/](https://cline.bot/)

配置百炼token：

打开Cline配置，输入百炼URL（[https://dashscope.aliyuncs.com/compatible-mode/v1](https://dashscope.aliyuncs.com/compatible-mode/v1)）、模型id（qwen-max）、百炼API key

配置如下：

```bash
{
  "mcpServers": {
    "hello-mcp-server": {
      "autoApprove": [],
      "disabled": false,
      "timeout": 600,
      "command": "/private/tmp/resource-center-mcp-server/.venv/bin/mcp",
      "args": [
        "run",
        "/private/tmp/resource-center-mcp-server/server.py"
      ],
      "transportType": "stdio",
      "env": {
        "HOME": "/Users/lazy"
      }
    }
  }
}
```

> 配置可以参考本地生成的 mcp-install.json；上序文件路径需要替换为你生成的 server 路径。
>
> 请注意：当前 env 的 HOME 必须配置，是你的电脑家目录。
>

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745310121592-3d6a1aad-8829-45aa-a89d-8ab86064528b.png)

变成绿色就可使用：

![](https://intranetproxy.alipay.com/skylark/lark/0/2025/png/309278/1745310169965-4e0c8728-ab05-4a03-948d-0c6fb743de63.png)




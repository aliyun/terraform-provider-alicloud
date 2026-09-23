# SLS Logtail Pipeline Config 示例

这个目录包含了 `alicloud_sls_logtail_pipeline_config` 资源的使用示例，包括基础配置和高级配置示例。

## 前提条件

1. **阿里云账号凭证**

```bash
export ALICLOUD_ACCESS_KEY="your_access_key"
export ALICLOUD_SECRET_KEY="your_secret_key"
export ALICLOUD_REGION="cn-shanghai"
```

2. **自动创建资源**

示例会自动创建 SLS Project 和 Logstore，使用随机后缀避免命名冲突。

## 使用方法

### 1. 初始化 Terraform

```bash
cd examples/sls_logtail_pipeline_config
terraform init
```

### 2. 查看执行计划

```bash
terraform plan
```

### 3. 应用配置

```bash
terraform apply
```

### 4. 查看输出

```bash
terraform output
```

### 5. 销毁资源

```bash
terraform destroy
```

## 配置示例

### 基础配置

创建一个简单的 Pipeline Config：
- **输入**：采集 `/home/*.log` 文件
- **处理器**：使用正则解析提取 key1、key2
- **输出**：发送到 SLS Logstore

```hcl
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "example" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = "terraform logtail pipeline config example"
}

resource "alicloud_log_store" "example" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = "example-store"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "example" {
  project     = alicloud_log_project.example.project_name
  config_name = "${var.name}-${random_integer.default.result}"

  # inputs: 输入配置
  inputs = [
    {
      Type                     = "input_file"
      FilePaths                = "[\\\"/home/*.log\\\"]"
      EnableContainerDiscovery = false
      MaxDirSearchDepth        = 0
      FileEncoding             = "utf8"
    }
  ]

  # processors: 处理器配置
  processors = [
    {
      Type      = "processor_parse_regex_native"
      SourceKey = "content"
      Regex     = ".*"
      Keys      = "[\\\"key1\\\",\\\"key2\\\"]"
    }
  ]

  # flushers: 输出配置
  flushers = [
    {
      Type          = "flusher_sls"
      Logstore      = alicloud_log_store.example.logstore_name
      TelemetryType = "logs"
      Region        = "cn-shanghai"
      Endpoint      = "cn-shanghai-intranet.log.aliyuncs.com"
    }
  ]
}

output "pipeline_config_id" {
  description = "The ID of the Pipeline Config (format: project:config_name)"
  value       = alicloud_sls_logtail_pipeline_config.example.id
}

output "pipeline_config_name" {
  description = "The name of the Pipeline Config"
  value       = alicloud_sls_logtail_pipeline_config.example.config_name
}

output "project" {
  description = "The SLS Project name"
  value       = alicloud_sls_logtail_pipeline_config.example.project
}

output "inputs_count" {
  description = "Number of input configurations"
  value       = length(alicloud_sls_logtail_pipeline_config.example.inputs)
}

output "processors_count" {
  description = "Number of processor configurations"
  value       = length(alicloud_sls_logtail_pipeline_config.example.processors)
}

output "flushers_count" {
  description = "Number of flusher configurations"
  value       = length(alicloud_sls_logtail_pipeline_config.example.flushers)
}
```

### 高级配置

包含更复杂的配置：
- 多路径输入（nginx 和 app 日志）
- Apache 日志格式解析
- 提取 IP、时间、HTTP 方法等字段

```hcl
resource "alicloud_log_store" "advanced" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = "nginx-logs"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "advanced" {
  project     = alicloud_log_project.example.project_name
  config_name = "${var.name}-advanced-${random_integer.default.result}"

  # 多路径输入
  inputs = [
    {
      Type      = "input_file"
      FilePaths = "[\\\"/var/log/nginx/*.log\\\",\\\"/var/log/app/*.log\\\"]"
      EnableContainerDiscovery = false
      MaxDirSearchDepth        = 0
      FileEncoding             = "utf8"
    }
  ]

  # 多个处理器（Apache 日志解析）
  processors = [
    {
      Type      = "processor_parse_regex_native"
      SourceKey = "content"
      Regex     = "([\\\\d\\\\.]+) \\\\S+ \\\\S+ \\\\[(\\\\S+) \\\\S+\\\\] \\\\\"(\\\\S+) (\\\\S+) (\\\\S+)\\\\\" ([\\\\d]+) ([\\\\d]+)"
      Keys      = "[\\\"ip\\\",\\\"time\\\",\\\"method\\\",\\\"path\\\",\\\"protocol\\\",\\\"status\\\",\\\"size\\\"]"
    }
  ]

  # 输出到多个目标
  flushers = [
    {
      Type          = "flusher_sls"
      Logstore      = alicloud_log_store.advanced.logstore_name
      TelemetryType = "logs"
      Region        = "cn-shanghai"
      Endpoint      = "cn-shanghai-intranet.log.aliyuncs.com"
    }
  ]

  log_sample = "Sample log entry for testing"
}

output "advanced_pipeline_config_id" {
  description = "The ID of the advanced Pipeline Config"
  value       = alicloud_sls_logtail_pipeline_config.advanced.id
}
```

## 字段类型说明

由于使用了 `TypeMap`，所有值都需要使用特定格式：

- **字符串**：直接写，如 `Type = "input_file"`
- **布尔值**：写 `true` 或 `false`（会自动转换）
- **数字**：写数字（会自动转换），如 `MaxDirSearchDepth = 0`
- **数组**：写 JSON 字符串，注意转义双引号，如 `FilePaths = "[\\\"/home/*.log\\\"]"`

## 常见问题

### Q1: 如何修改采集路径？

修改 `inputs` 中的 `FilePaths`：

```hcl
FilePaths = "[\\\"/var/log/app/*.log\\\",\\\"/var/log/nginx/*.log\\\"]"
```

### Q2: 如何添加更多处理器？

`processors` 是一个列表，但每个元素只能是一个 map。如果需要多个处理器，需要在数组中添加多个元素。

### Q3: 如何修改输出的 Logstore？

修改 `flushers` 中的 `Logstore` 字段：

```hcl
flushers = [
  {
    Type     = "flusher_sls"
    Logstore = "your-logstore-name"
    # ... 其他字段
  }
]
```

## 参考资料

- [SLS Logtail Pipeline Config API 文档](https://help.aliyun.com/document_detail/...)
- [Terraform Provider 文档](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs)

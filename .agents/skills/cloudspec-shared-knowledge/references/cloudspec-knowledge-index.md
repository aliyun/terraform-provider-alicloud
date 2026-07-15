# CloudSpec 知识目录

## 📚 目录概览

本文档是 CloudSpec 相关知识的完整目录结构，用于系统化地组织和关联 CloudSpec 生态中的各种概念、技术、工具和最佳实践。

---

## 🏗️ 1. CloudSpec 基础语法

### 1.1 核心概念
- [CloudSpec 基础概念](./docs/corpora/common/basic-grammer/quick-start.md#cloudspec模型是什么)
- [基本概念介绍](./docs/corpora/common/basic-grammer/quick-start.md#基本概念介绍)
- [模型类型](./docs/corpora/common/basic-grammer/quick-start.md#模型类型)

### 1.2 基础类型系统
- [基础类型详解](./docs/corpora/common/basic-grammer/basic-types.md)
- [复合类型说明](./docs/corpora/common/basic-grammer/composite-types.md)

### 1.3 语法规范
- [CloudSpec IDL 语法指南](./docs/corpora/common/basic-grammer/cloudspec-idl-grammer.md)
- [命名空间管理](./docs/corpora/common/basic-grammer/namespace.md)
- [保留字列表](./docs/corpora/common/basic-grammer/reserved-words.md)

### 1.4 模型组合
- [模型组合](./docs/corpora/common/basic-grammer/model-composition.md)

### 1.5 快速入门
- [CloudSpec IDL 生成指南](./docs/corpora/common/cloudspec-idl-generation-guide.md)

---

## 🏷️ 2. 注解系统（Annotations）

### 2.1 Service 注解
- [Service 注解详解](./docs/corpora/common/annotates/service-types.md)

### 2.2 Operation 注解
- [Operation 注解详解](./docs/corpora/common/annotates/operation-annotate.md)
  - [HTTP 注解](./docs/corpora/common/annotates/operation-annotate.md#a61-http-注解)
  - [分页注解](./docs/corpora/common/annotates/operation-annotate.md#a62-paginated-注解)
  - [操作信息注解](./docs/corpora/common/annotates/operation-annotate.md#a63-operationinfo-注解)
  - [后端配置注解](./docs/corpora/common/annotates/operation-annotate.md#a64-backendconfigurationhttp-注解)
  - [错误映射注解](./docs/corpora/common/annotates/operation-annotate.md#a69-errormapping-注解)
  - [异步配置注解](./docs/corpora/common/annotates/operation-annotate.md#a634-async-注解)
- [Operation 注解详解 - 设计指南](./docs/corpora/common/how-to-design-operation.md#operation-注解详解)

### 2.3 Resource 注解
- [Resource 注解详解](./docs/corpora/common/annotates/resource-annotate.md)

### 2.4 其他注解
- [赋值约束注解](./docs/corpora/common/annotates/value-constraint-annotate.md)
- [文档注解](./docs/corpora/common/annotates/document-annotate.md)
- [测试注解](./docs/corpora/common/annotates/test-annotate.md)
- [赋值约束注解](./docs/corpora/common/annotates/assignment-constraint-annotate.md)

---

## 🧪 3. 资源测试（Resource Testing）

### 3.1 测试指南
- [资源测试快速开始](./docs/corpora/resource-test/test-guide/resource-test-quick-start.md)

### 3.2 测试常见问题 FAQ
- [网络错误处理](./docs/corpora/resource-test/resource-test-faq/1-network-error.md)
- [重试策略错误](./docs/corpora/resource-test/resource-test-faq/2-retryPolicies-error.md)
- [资源不存在条件错误](./docs/corpora/resource-test/resource-test-faq/3-resourceNotExistCondition-error.md)
- [异步错误处理](./docs/corpora/resource-test/resource-test-faq/4-async-error.md)
- [差异错误处理](./docs/corpora/resource-test/resource-test-faq/5-diff-error.md)
- [枚举错误处理](./docs/corpora/resource-test/resource-test-faq/6-enum-error.md)
- [依赖错误处理](./docs/corpora/resource-test/resource-test-faq/7-dependency-error.md)
- [ListTagResources 错误](./docs/corpora/resource-test/resource-test-faq/7-ListTagResources-error.md)

### 3.3 CLI 工具
- [Aliyun CLI cspec plugin 快速开始](./docs/corpora/cli-guide/cloudspec-cli-quick-start.md)

---

## 📋 4. API 设计

### 4.1 设计原则
- [Operation 设计原则](./docs/corpora/common/how-to-design-operation.md#operation-设计原则)

### 4.2 命名规范
- [Operation 设计最佳实践 - 命名规范](./docs/corpora/common/how-to-design-operation.md#operation-设计最佳实践)

### 4.3 参数设计
- [Operation 设计基础 - 核心要素](./docs/corpora/common/how-to-design-operation.md#operation-设计基础)
- [Operation 配置规范 - 参数设计](./docs/corpora/common/how-to-design-operation.md#operation-配置规范)

### 4.4 错误处理
- [Operation 配置规范 - 错误处理](./docs/corpora/common/how-to-design-operation.md#operation-配置规范)
- [常见问题与解决方案 - 错误处理](./docs/corpora/common/how-to-design-operation.md#常见问题与解决方案)

### 4.5 版本管理
- TODO: 添加 API 版本管理规范

### 4.6 安全设计
- [Operation 配置规范 - 安全配置](./docs/corpora/common/how-to-design-operation.md#operation-配置规范)

### 4.7 性能优化
- [Operation 设计最佳实践 - 性能优化](./docs/corpora/common/how-to-design-operation.md#operation-设计最佳实践)

### 4.8 文档规范
- TODO: 添加 API 文档规范

---

## 5. 企业级能力

### 5.1 Terraform 集成
- TODO: 添加 Terraform 集成文档

### 5.2 RMC（资源管理控制台）集成
- TODO: 添加 RMC 集成文档

### 5.3 RAM（访问控制）集成
- TODO: 添加 RAM 集成文档

### 5.4 其他企业级能力


---

## 🏛️ 6. 平台方知识

### 6.1 镇元平台
- [镇元平台使用场景与元数据关系]


---

## 📝 7. 常见缩写简写对应关系

### 7.1 技术术语
- TODO: 添加技术术语缩写表

### 7.2 云服务相关
- TODO: 添加云服务缩写表

### 7.3 企业级能力
- TODO

### 7.4 开发相关
- TODO

---

## 🔧 8. 如何定义 Operation

### 8.1 设计基础
- [Operation 设计基础](./docs/corpora/common/how-to-design-operation.md#operation-设计基础)
- [Operation 设计原则](./docs/corpora/common/how-to-design-operation.md#operation-设计原则)

### 8.2 配置规范
- [Operation 配置规范](./docs/corpora/common/how-to-design-operation.md#operation-配置规范)
- [Operation 注解详解](./docs/corpora/common/how-to-design-operation.md#operation-注解详解)

### 8.3 最佳实践
- [Operation 设计最佳实践](./docs/corpora/common/how-to-design-operation.md#operation-设计最佳实践)

### 8.4 设计检查清单
- [Operation 设计检查清单](./docs/corpora/common/how-to-design-operation.md#operation-设计检查清单)

### 8.5 常见问题
- [常见问题与解决方案](./docs/corpora/common/how-to-design-operation.md#常见问题与解决方案)

---

## 🏗️ 9. 如何定义资源

### 9.1 设计指南
- [资源设计指南](./docs/corpora/common/how-to-design-resource.md)

### 9.2 设计检查清单
- TODO: 添加资源设计检查清单

### 9.3 最佳实践
- TODO: 添加资源设计最佳实践

### 9.4 常见问题
- TODO: 添加资源设计常见问题

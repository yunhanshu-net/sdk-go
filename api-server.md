# API-Server 设计蓝图

## 1. 系统概述

API-Server是云函数(yunhanshu)项目的核心组件，负责管理所有runner的元数据和API接口信息，处理前端请求并转发至runcher执行引擎，同时提供runner的生命周期管理功能。

## 2. 架构设计

### 2.1 整体架构

```
+----------------+     +-----------------+     +----------------+
|                |     |                 |     |                |
|  Web前端       +---->+   API-Server    +---->+    Runcher     |
|                |     |  (公有云部署)    |     |  (私有云部署)  |
|                |     |                 |     |                |
+----------------+     +-------+---------+     +--------+-------+
                               |                        |
                       +-------v---------+     +--------v-------+
                       |                 |     |                |
                       |   数据库存储    |     |     Runner     |
                       |                 |     |                |
                       +-----------------+     +----------------+
```

### 2.2 核心组件

1. **HTTP服务层**：处理REST API请求
2. **NATS消息层**：与runcher通信
3. **数据访问层**：管理数据库操作
4. **认证授权层**：管理用户权限
5. **Runner管理层**：处理runner生命周期

## 3. 数据模型设计

### 3.1 Runner表（优化版）
```sql
CREATE TABLE `runner` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `created_by` bigint(20) DEFAULT NULL,
  `updated_by` bigint(20) DEFAULT NULL,
  `en_name` varchar(64) NOT NULL COMMENT 'runner英文标识',
  `cn_name` varchar(128) NOT NULL COMMENT 'runner中文名称',
  `description` text COMMENT 'runner描述',
  `version` varchar(32) NOT NULL COMMENT '当前版本',
  `user_id` bigint(20) NOT NULL COMMENT '所属用户ID',
  `language` varchar(32) NOT NULL COMMENT 'SDK语言',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1-正常，2-禁用',
  `runcher_id` bigint(20) DEFAULT NULL COMMENT '部署的runcher ID',
  `is_public` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否公开',
  `logo` varchar(255) DEFAULT NULL COMMENT '项目logo',
  `tags` varchar(255) DEFAULT NULL COMMENT '标签，多个用;分隔',
  `max_memory` int(11) DEFAULT '256' COMMENT '最大内存限制(MB)',
  `max_cpu` int(11) DEFAULT '1' COMMENT '最大CPU核心数',
  `deploy_config` json DEFAULT NULL COMMENT '部署配置项',
  `health_status` tinyint(4) DEFAULT '0' COMMENT '健康状态：0-未知，1-健康，2-异常',
  `last_health_check` datetime DEFAULT NULL COMMENT '最后健康检查时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_en_name` (`user_id`, `en_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.2 Service_Tree表（优化版）
```sql
CREATE TABLE `service_tree` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `created_by` bigint(20) DEFAULT NULL,
  `updated_by` bigint(20) DEFAULT NULL,
  `runner_id` bigint(20) NOT NULL COMMENT '所属runner ID',
  `code` varchar(64) NOT NULL COMMENT '服务标识',
  `name` varchar(128) NOT NULL COMMENT '服务名称',
  `desc` text COMMENT '服务描述',
  `tags` varchar(255) DEFAULT NULL COMMENT '标签，多个用;分隔',
  `is_public` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否公开',
  `user_id` bigint(20) NOT NULL COMMENT '所属用户ID',
  `source_user_id` bigint(20) DEFAULT NULL COMMENT '源用户ID(fork情况)',
  `fork_id` bigint(20) DEFAULT NULL COMMENT '被fork的service_tree ID',
  `parent_id` bigint(20) DEFAULT NULL COMMENT '父服务ID',
  `children_count` int(11) NOT NULL DEFAULT '0' COMMENT '子服务数量',
  `full_path` varchar(512) NOT NULL COMMENT '全路径',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '层级',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `access_control` json DEFAULT NULL COMMENT '访问控制配置',
  `api_count` int(11) NOT NULL DEFAULT '0' COMMENT 'API数量',
  `dependencies` json DEFAULT NULL COMMENT '服务依赖关系',
  `version` varchar(32) DEFAULT NULL COMMENT '服务版本标识',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_runner_parent_code` (`runner_id`, `parent_id`, `code`),
  KEY `idx_runner_full_path` (`runner_id`, `full_path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.3 Runner_Func表（优化版）
```sql
CREATE TABLE `runner_func` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `created_by` bigint(20) DEFAULT NULL,
  `updated_by` bigint(20) DEFAULT NULL,
  `runner_id` bigint(20) NOT NULL COMMENT '所属runner ID',
  `tree_id` bigint(20) NOT NULL COMMENT '所属服务树ID',
  `code` varchar(64) NOT NULL COMMENT '接口标识',
  `name` varchar(128) NOT NULL COMMENT '接口名称',
  `desc` text COMMENT '接口描述',
  `tags` varchar(255) DEFAULT NULL COMMENT '标签，多个用;分隔',
  `request` json DEFAULT NULL COMMENT '请求参数',
  `response` json DEFAULT NULL COMMENT '响应参数',
  `callbacks` json DEFAULT NULL COMMENT '回调列表',
  `use_tables` json DEFAULT NULL COMMENT '使用的表',
  `is_public` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否公开',
  `user_id` bigint(20) NOT NULL COMMENT '所属用户ID',
  `source_user_id` bigint(20) DEFAULT NULL COMMENT '源用户ID(fork情况)',
  `fork_id` bigint(20) DEFAULT NULL COMMENT '被fork的func ID',
  `http_method` varchar(10) NOT NULL DEFAULT 'GET' COMMENT 'HTTP方法',
  `api_path` varchar(255) NOT NULL COMMENT 'API路径',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1-正常，2-禁用',
  `avg_execute_time` int(11) DEFAULT '0' COMMENT '平均执行时间(ms)',
  `max_execute_time` int(11) DEFAULT '0' COMMENT '最大执行时间(ms)',
  `rate_limit` int(11) DEFAULT NULL COMMENT 'API调用频率限制(次/分钟)',
  `cache_strategy` json DEFAULT NULL COMMENT '缓存策略',
  `timeout` int(11) DEFAULT '30000' COMMENT '执行超时时间(ms)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_runner_tree_code` (`runner_id`, `tree_id`, `code`),
  UNIQUE KEY `idx_runner_api_path` (`runner_id`, `api_path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.4 Runner_Version表
```sql
CREATE TABLE `runner_version` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `created_by` bigint(20) DEFAULT NULL,
  `updated_by` bigint(20) DEFAULT NULL,
  `runner_id` bigint(20) NOT NULL COMMENT '所属runner ID',
  `version` varchar(32) NOT NULL COMMENT '版本号',
  `change_log` text COMMENT '变更日志',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1-正常，2-回滚',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_runner_version` (`runner_id`, `version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.5 API_Metrics表（新增）
```sql
CREATE TABLE `api_metrics` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `func_id` bigint(20) NOT NULL COMMENT '关联的函数ID',
  `call_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '调用次数',
  `success_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '成功次数',
  `error_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '错误次数',
  `avg_response_time` int(11) NOT NULL DEFAULT '0' COMMENT '平均响应时间(ms)',
  `max_response_time` int(11) NOT NULL DEFAULT '0' COMMENT '最大响应时间(ms)',
  `p95_response_time` int(11) NOT NULL DEFAULT '0' COMMENT '95%响应时间(ms)',
  `last_called_at` datetime DEFAULT NULL COMMENT '最后调用时间',
  `stat_date` date NOT NULL COMMENT '统计日期',
  `memory_usage` int(11) DEFAULT '0' COMMENT '内存使用(MB)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_func_date` (`func_id`, `stat_date`),
  KEY `idx_stat_date` (`stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.6 Runner_Dependencies表（新增）
```sql
CREATE TABLE `runner_dependencies` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `created_by` bigint(20) DEFAULT NULL,
  `runner_id` bigint(20) NOT NULL COMMENT '依赖方runner ID',
  `depend_runner_id` bigint(20) NOT NULL COMMENT '被依赖runner ID',
  `depend_tree_id` bigint(20) DEFAULT NULL COMMENT '被依赖服务树ID',
  `depend_func_id` bigint(20) DEFAULT NULL COMMENT '被依赖函数ID',
  `depend_type` tinyint(4) NOT NULL COMMENT '依赖类型：1-Runner，2-ServiceTree，3-Func',
  `version_constraint` varchar(32) DEFAULT NULL COMMENT '版本约束',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1-正常，2-异常',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_runner_depend` (`runner_id`, `depend_type`, `depend_runner_id`, `depend_tree_id`, `depend_func_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.7 User_Permission表（新增）
```sql
CREATE TABLE `user_permission` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `created_by` bigint(20) DEFAULT NULL,
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `resource_type` tinyint(4) NOT NULL COMMENT '资源类型：1-Runner，2-ServiceTree，3-Func',
  `resource_id` bigint(20) NOT NULL COMMENT '资源ID',
  `permission` varchar(32) NOT NULL COMMENT '权限类型：read,write,execute,admin',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_resource` (`user_id`, `resource_type`, `resource_id`, `permission`),
  KEY `idx_resource` (`resource_type`, `resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.8 Runcher表（新增）
```sql
CREATE TABLE `runcher` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `name` varchar(64) NOT NULL COMMENT 'runcher名称',
  `description` text COMMENT 'runcher描述',
  `host` varchar(255) NOT NULL COMMENT 'runcher主机地址',
  `port` int(11) NOT NULL DEFAULT '4222' COMMENT 'NATS端口',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1-在线，2-离线',
  `last_heartbeat` datetime DEFAULT NULL COMMENT '最后心跳时间',
  `region` varchar(64) DEFAULT NULL COMMENT '部署区域',
  `max_runners` int(11) DEFAULT '100' COMMENT '最大runner数量',
  `current_runners` int(11) DEFAULT '0' COMMENT '当前runner数量',
  `cpu_usage` float DEFAULT '0' COMMENT 'CPU使用率',
  `memory_usage` float DEFAULT '0' COMMENT '内存使用率',
  `disk_usage` float DEFAULT '0' COMMENT '磁盘使用率',
  `tags` varchar(255) DEFAULT NULL COMMENT '标签，多个用;分隔',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_host_port` (`host`, `port`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.9 Operation_Log表（新增）
```sql
CREATE TABLE `operation_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `user_id` bigint(20) NOT NULL COMMENT '操作用户ID',
  `operation_type` varchar(32) NOT NULL COMMENT '操作类型',
  `resource_type` tinyint(4) NOT NULL COMMENT '资源类型：1-Runner，2-ServiceTree，3-Func',
  `resource_id` bigint(20) NOT NULL COMMENT '资源ID',
  `operation_desc` text COMMENT '操作描述',
  `operation_result` tinyint(4) NOT NULL DEFAULT '1' COMMENT '结果：1-成功，2-失败',
  `error_msg` text COMMENT '错误信息',
  `client_ip` varchar(64) DEFAULT NULL COMMENT '客户端IP',
  `operation_data` json DEFAULT NULL COMMENT '操作数据',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 4. 核心功能设计

### 4.1 Runner管理

#### 4.1.1 Runner创建流程
1. 接收前端创建请求，验证参数
2. 在数据库创建runner记录
3. 向runcher发送创建指令
4. 接收runcher创建结果
5. 更新runner状态

#### 4.1.2 Runner版本管理
1. 支持版本发布、查看历史版本
2. 支持一键回滚到指定版本
3. 提供版本之间差异比较功能

#### 4.1.3 Runner资源管理（新增）
1. 监控和控制runner资源使用
2. 支持动态调整资源上限
3. 提供资源使用统计报表

### 4.2 服务树管理

#### 4.2.1 服务树的CRUD操作
1. 创建/修改/删除服务树节点
2. 按层级查询服务树
3. 服务树节点排序

#### 4.2.2 服务树Fork功能
1. 选择目标服务树节点进行Fork
2. 支持整个服务树或单个节点Fork
3. 可追踪Fork关系和来源

#### 4.2.3 服务树依赖管理（新增）
1. 管理服务间依赖关系
2. 依赖冲突检测
3. 依赖版本约束

### 4.3 API管理

#### 4.3.1 API注册与同步
1. 从runner获取API信息
2. 同步到数据库
3. 提供API文档生成功能

#### 4.3.2 API调用流程
1. 接收前端请求
2. 查询API信息
3. 通过NATS向runcher发送请求
4. 接收runcher执行结果
5. 返回前端结果

#### 4.3.3 API性能监控（新增）
1. 记录API调用次数、响应时间
2. 性能异常检测
3. 性能优化建议

### 4.4 安全与权限

#### 4.4.1 认证授权
1. 用户身份认证
2. 基于角色的权限控制
3. API访问权限管理

#### 4.4.2 安全防护
1. 请求参数验证
2. 代码执行沙箱
3. API调用频率限制

#### 4.4.3 细粒度权限控制（新增）
1. 资源级别的权限控制
2. 操作类型权限（读、写、执行、管理）
3. 权限继承与覆盖机制

### 4.5 监控与统计

#### 4.5.1 系统监控
1. API调用次数、耗时统计
2. Runner状态监控
3. 系统资源使用监控

#### 4.5.2 告警机制
1. 定义告警规则
2. 多渠道告警通知
3. 告警历史查询

#### 4.5.3 数据分析（新增）
1. API调用趋势分析
2. 资源使用情况分析
3. 用户行为分析

## 5. API接口设计

### 5.1 Runner相关API
- `POST /api/runners` - 创建runner
- `GET /api/runners` - 获取runner列表
- `GET /api/runners/:id` - 获取runner详情
- `PUT /api/runners/:id` - 更新runner
- `DELETE /api/runners/:id` - 删除runner
- `POST /api/runners/:id/deploy` - 部署runner
- `POST /api/runners/:id/rollback` - 回滚runner
- `GET /api/runners/:id/metrics` - 获取runner资源使用情况（新增）
- `POST /api/runners/:id/scale` - 调整runner资源配置（新增）

### 5.2 服务树相关API
- `POST /api/trees` - 创建服务树节点
- `GET /api/trees` - 获取服务树列表
- `GET /api/trees/:id` - 获取服务树详情
- `PUT /api/trees/:id` - 更新服务树节点
- `DELETE /api/trees/:id` - 删除服务树节点
- `POST /api/trees/:id/fork` - Fork服务树节点
- `GET /api/trees/:id/dependencies` - 获取服务依赖关系（新增）
- `POST /api/trees/:id/dependencies` - 添加服务依赖（新增）

### 5.3 API相关接口
- `POST /api/funcs` - 创建/更新API
- `GET /api/funcs` - 获取API列表
- `GET /api/funcs/:id` - 获取API详情
- `DELETE /api/funcs/:id` - 删除API
- `POST /api/funcs/:id/fork` - Fork API
- `GET /api/funcs/:id/metrics` - 获取API调用统计（新增）
- `PUT /api/funcs/:id/rate-limit` - 设置API调用频率限制（新增）
- `PUT /api/funcs/:id/cache` - 设置API缓存策略（新增）

### 5.4 API调用接口
- `ANY /run/:user/:runner/*` - API统一调用入口

### 5.5 权限管理接口（新增）
- `POST /api/permissions` - 创建权限
- `GET /api/permissions` - 获取权限列表
- `DELETE /api/permissions/:id` - 删除权限
- `GET /api/resources/:type/:id/permissions` - 获取资源权限

## 6. 性能优化设计

### 6.1 缓存策略
1. API元数据缓存
2. 热点runner信息缓存
3. 用户权限缓存

### 6.2 数据库优化
1. 表分区策略
2. 索引优化
3. 读写分离

### 6.3 高并发处理
1. 异步处理机制
2. 限流策略
3. 降级策略

### 6.4 API结果缓存（新增）
1. 按API配置缓存结果
2. 缓存失效策略
3. 缓存命中率监控

## 7. 扩展性设计

### 7.1 插件系统
1. 定义插件接口
2. 支持自定义认证插件
3. 支持自定义监控插件

### 7.2 多语言SDK支持
1. 支持Go SDK
2. 预留Python、Java等SDK扩展接口

### 7.3 第三方集成
1. 支持AI代码生成集成
2. 支持CI/CD系统集成
3. 支持监控系统集成

### 7.4 触发器机制（新增）
1. 支持定时触发
2. 支持事件触发
3. 支持消息队列触发

## 8. 部署架构

### 8.1 单区域部署
```
[负载均衡] -> [API-Server集群] -> [NATS集群] -> [Runcher节点]
                    |
                    v
            [数据库主从集群]
```

### 8.2 多区域部署
```
                      [全局负载均衡]
                      /            \
    [区域A负载均衡]                  [区域B负载均衡]
           |                              |
  [区域A API-Server集群]          [区域B API-Server集群]
           |                              |
     [区域A NATS集群] <--(同步)--> [区域B NATS集群]
           |                              |
   [区域A Runcher节点]              [区域B Runcher节点]
```

### 8.3 混合云部署（新增）
```
                   [API-Server集群(公有云)]
                           |
                           v
                      [NATS网关]
                      /        \
[Runcher节点(公有云)] <--> [Runcher节点(私有云)] 
```

## 9. 安全性设计

### 9.1 数据安全
1. 敏感数据加密存储
2. 数据备份与恢复机制
3. 数据传输加密

### 9.2 应用安全
1. API访问鉴权
2. 代码沙箱执行
3. 资源隔离

### 9.3 运行时安全
1. 执行超时控制
2. 资源限制
3. 异常监控与处理

### 9.4 合规与审计（新增）
1. 操作日志记录
2. 权限变更审计
3. 代码安全扫描

## 10. 开发与运维计划

### 10.1 开发阶段规划
1. 基础框架搭建 - 2周
2. 核心功能开发 - 4周
3. 性能优化与测试 - 2周
4. 安全加固 - 1周

### 10.2 部署与运维
1. 环境准备与部署流程
2. 监控告警配置
3. 灰度发布策略
4. 应急预案

### 10.3 DevOps支持（新增）
1. CI/CD流程设计
2. 自动化测试
3. 蓝绿部署支持

## 11. 未来规划

### 11.1 功能迭代
1. AI代码生成集成
2. 可视化开发界面
3. 更多SDK语言支持

### 11.2 性能提升
1. 分布式调度优化
2. 冷热数据分离
3. 计算资源弹性伸缩

### 11.3 生态建设
1. 开放API市场
2. 社区贡献机制
3. 开发者工具套件

### 11.4 多云适配（新增）
1. 支持主流云厂商
2. 跨云部署方案
3. 云资源统一管理

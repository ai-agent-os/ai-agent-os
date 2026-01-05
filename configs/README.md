# 配置文件说明

## 配置文件列表

### 核心配置
- `global.yaml` - 全局共享配置（NATS、JWT、Control Service 等）
- `api-gateway.yaml` - API 网关配置（端口 9090）
- `app-server.yaml` - 应用服务配置（端口 9091）
- `app-storage.yaml` - 存储服务配置（端口 9092）
- `app-runtime.yaml` - 运行时服务配置（端口 9093）
- `agent-server.yaml` - 智能体服务配置（端口 9095）
- `hr-server.yaml` - HR服务配置（端口 9097，用户管理、组织架构、职责管理）
- `control-service.yaml` - 控制服务配置

### 其他配置
- `casbin_model.conf` - Casbin 权限模型配置

## 端口分配

- `9090` - API Gateway（网关）
- `9091` - App Server（应用服务）
- `9092` - App Storage（存储服务）
- `9093` - App Runtime（运行时服务）
- `9094` - Hub Server（Hub 服务）
- `9095` - Agent Server（智能体服务）
- `9096` - Control Service（控制服务）
- `9097` - HR Server（HR 服务）

## 使用说明

这些配置文件用于系统的各个组件，提供灵活的配置管理。每个服务都会读取对应的配置文件来初始化。

## 配置优先级

1. 命令行参数
2. 环境变量
3. 配置文件
4. 默认值
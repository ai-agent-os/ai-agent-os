# 配置文件说明

## 配置文件列表

### 核心配置
- `ai-agent-os.yaml` - AI Agent OS 全局配置
- `app-runtime.yaml` - 应用运行时服务配置
- `podman.yaml` - Podman 容器服务配置
- `function-server.yaml` - 函数服务器配置
- `container.yaml` - 容器相关配置

### 模板文件
- `main.go.template` - main.go 文件模板
- `readme.md.template` - README 文件模板

## 使用说明

这些配置文件用于系统的各个组件，提供灵活的配置管理。每个服务都会读取对应的配置文件来初始化。

## 配置优先级

1. 命令行参数
2. 环境变量
3. 配置文件
4. 默认值
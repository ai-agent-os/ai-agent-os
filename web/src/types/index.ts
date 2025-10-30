// API响应基础类型
export interface ApiResponse<T = any> {
  code: number
  data: T
  message?: string
}

// 用户相关类型
export interface UserInfo {
  id: number
  username: string
  email: string
  register_type: string
  avatar: string
  email_verified: boolean
  status: string
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  code?: string
}

// 应用相关类型
export interface App {
  id: number
  user: string
  code: string
  name: string
  nats_id: number
  host_id: number
  status: 'enabled' | 'disabled'
  version: string
  created_at: string
  updated_at: string
}

export interface CreateAppRequest {
  code: string
  name: string
}

// 服务目录相关类型
export interface ServiceTree {
  id: number
  name: string
  code: string
  parent_id: number
  type: 'package' | 'function'
  description: string
  tags: string
  app_id: number
  ref_id: number
  full_code_path: string
  created_at: string
  updated_at: string
  children?: ServiceTree[]
}

export interface CreateServiceTreeRequest {
  user: string
  app: string
  name: string
  code: string
  parent_id?: number
  description?: string
  tags?: string
}

// 函数相关类型
export interface Function {
  id: number
  request: any
  response: FieldConfig[]
  app_id: number
  tree_id: number
  method: string
  router: string
  has_config: boolean
  create_tables: string
  callbacks: string
  template_type: string
  created_at: string
  updated_at: string
}

// 字段配置类型
export interface FieldConfig {
  code: string
  name: string
  data: {
    type: string
    format?: string
    source?: string
    example?: string
    default_value?: string
  }
  desc?: string
  search?: string | null
  table_permission?: string | null
  widget: WidgetConfig
  callbacks?: any
  permission?: string | null
  validation?: string
}

// 组件配置类型
export interface WidgetConfig {
  type: string
  config: Record<string, any>
}

// 组件类型枚举
export enum WidgetType {
  INPUT = 'input',
  SELECT = 'select',
  TEXT_AREA = 'text_area',
  FILE_UPLOAD = 'file_upload',
  USER = 'user',
  DATETIME = 'datetime',
  NUMBER = 'number',
  SWITCH = 'switch',
  CHECKBOX = 'checkbox',
  RADIO = 'radio'
}

// 搜索类型
export interface SearchParams {
  eq?: string      // 精确匹配 eq=id:1
  like?: string    // 模糊匹配 like=title:xxx
  gte?: string     // 大于等于 gte=created_at:timestamp
  lte?: string     // 小于等于 lte=created_at:timestamp
  sort?: string    // 排序 sort=id:desc
  page?: number    // 页码
  page_size?: number // 页大小
}

// 表格响应类型
export interface TableResponse<T = any> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 路由配置
export interface RouteConfig {
  path: string
  name: string
  component: any
  meta?: {
    title?: string
    icon?: string
    requireAuth?: boolean
  }
}

// WebSocket消息类型
export interface WSMessage {
  type: string
  data: any
  timestamp: number
}
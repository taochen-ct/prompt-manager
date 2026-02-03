# Prompt Manager API 文档

## 概述

Prompt Manager 后端 API 基于 Go + Gin 框架开发，提供完整的用户、提示词(Prompt)和版本管理功能。

## 基础信息

- **Base URL**: `http://localhost:8080`
- **API 前缀**: `/api/v1`
- **Content-Type**: `application/json`
- **响应格式**: 统一 JSON 响应

### 统一响应结构

```json
{
  "code": 0,
  "data": {},
  "message": "success"
}
```

---

## User API

### 创建用户

**接口**: `POST /api/v1/user/create`

**请求参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| username | string | 是 | 用户名 (3-32字符) |
| nickname | string | 否 | 昵称 |
| department | string | 否 | 部门 |

**请求示例**:
```json
{
  "username": "admin",
  "nickname": "管理员",
  "department": "技术部"
}
```

**响应参数**:

| 字段 | 类型 | 描述 |
|------|------|------|
| id | string | 用户ID |
| username | string | 用户名 |
| nickname | string | 昵称 |
| department | string | 部门 |
| createdAt | string | 创建时间 |
| updatedAt | string | 更新时间 |

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "id": "1",
    "username": "admin",
    "nickname": "管理员",
    "department": "技术部",
    "createdAt": "2024-01-01 00:00:00",
    "updatedAt": "2024-01-01 00:00:00"
  },
  "message": "success"
}
```

---

### 获取用户

**接口**: `GET /api/v1/user/info/:id`

**路径参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 用户ID |

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "id": "1",
    "username": "admin",
    "nickname": "管理员",
    "department": "技术部",
    "createdAt": "2024-01-01 00:00:00",
    "updatedAt": "2024-01-01 00:00:00"
  },
  "message": "success"
}
```

---

### 更新用户

**接口**: `POST /api/v1/user/update/:id`

**路径参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 用户ID |

**请求参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| nickname | string | 否 | 昵称 |
| department | string | 否 | 部门 |

**请求示例**:
```json
{
  "nickname": "更新后的昵称",
  "department": "产品部"
}
```

---

### 删除用户

**接口**: `POST /api/v1/user/delete`

**请求参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 用户ID |

**请求示例**:
```json
{
  "id": "1"
}
```

---

## Prompt API

### 创建提示词

**接口**: `POST /api/v1/prompt/create`

**请求参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|----|------|
| name | string | 是  | 提示词名称 |
| createdBy | string | 是  | 创建者ID |
| username | string | 是  | 创建者用户名 |
| path | string | 是  | 路径 |

**请求示例**:
```json
{
  "name": "文案生成助手",
  "createdBy": "user123",
  "username": "管理员",
  "path": "/text-generate"
}
```

**响应参数**:

| 字段 | 类型 | 描述 |
|------|------|------|
| id | string | 提示词ID |
| name | string | 提示词名称 |
| path | string | 路径 |
| latestVersion | string | 最新版本号 |
| isPublish | boolean | 是否发布 |
| createBy | string | 创建者ID |
| username | string | 创建者用户名 |
| createAt | string | 创建时间 |
| updateAt | string | 更新时间 |

---

### 获取提示词

**接口**: `GET /api/v1/prompt/info/:id`

**路径参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 提示词ID |

---

### 获取提示词内容

**接口**: `GET /api/v1/prompt/content/*path`

**路径参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| path | string | 是 | 提示词路径 |

**查询参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| path | string | 是 | 提示词路径 |

**业务逻辑**:
- 如果提示词未发布 (`isPublish=false`)，返回错误 `"no published version"`
- 如果提示词已发布，根据 `latestVersion` (版本ID) 查询版本详情返回

**响应示例** (已发布):
```json
{
  "code": 0,
  "data": {
    "prompt": {
      "id": "xxx-xxx-xxx",
      "name": "文案生成助手",
      "path": "文案生成助手",
      "latestVersion": "version-xxx-xxx",
      "isPublish": true,
      "createBy": "user123",
      "username": "管理员",
      "createAt": "2024-01-01 00:00:00",
      "updateAt": "2024-01-01 00:00:00"
    },
    "version": {
      "id": "version-xxx-xxx",
      "promptId": "xxx-xxx-xxx",
      "version": "1.0.0",
      "content": "你是一个专业的文案生成助手...",
      "variables": "[\"topic\", \"tone\"]",
      "isPublish": true,
      "changeLog": "初始版本",
      "createdBy": "user123",
      "username": "管理员",
      "createdAt": "2024-01-01 00:00:00",
      "updatedAt": "2024-01-01 00:00:00"
    }
  },
  "message": "success"
}
```

**响应示例** (未发布):
```json
{
  "code": 400,
  "data": null,
  "message": "no published version"
}
```

---

### 更新提示词

**接口**: `POST /api/v1/prompt/update`

**请求参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 提示词ID |
| name | string | 是 | 提示词名称 |
| isPublish | boolean | 否 | 是否发布 |

**请求示例**:
```json
{
  "id": "xxx-xxx-xxx",
  "name": "更新后的名称",
  "isPublish": false
}
```

---

### 删除提示词

**接口**: `POST /api/v1/prompt/delete/:id`

**路径参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 提示词ID |

---

### 提示词列表

**接口**: `GET /api/v1/prompt/list`

**查询参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| offset | int | 否 | 偏移量 (默认0) |
| limit | int | 否 | 限制数量 (默认10) |

**响应示例**:
```json
{
  "code": 0,
  "data": [
    {
      "id": "xxx-xxx-xxx",
      "name": "文案生成助手",
      "path": "文案生成助手",
      "latestVersion": "version-xxx-xxx",
      "isPublish": true,
      "createBy": "user123",
      "username": "管理员",
      "createAt": "2024-01-01 00:00:00",
      "updateAt": "2024-01-01 00:00:00"
    }
  ],
  "message": "success"
}
```

---

## Prompt Version API

### 创建版本

**接口**: `POST /api/v1/version/create`

**请求参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| promptId | string | 是 | 所属提示词ID |
| version | string | 是 | 版本号 |
| content | string | 是 | 提示词内容 |
| variables | string | 否 | 变量定义 (JSON数组字符串) |
| changeLog | string | 否 | 更新日志 |
| createdBy | string | 是 | 创建者ID |
| username | string | 是 | 创建者用户名 |
| isPublish | boolean | 否 | 是否发布 (默认false) |

**请求示例**:
```json
{
  "promptId": "xxx-xxx-xxx",
  "version": "1.0.0",
  "content": "你是一个专业的文案生成助手，请根据用户提供的信息生成吸引人的文案。",
  "variables": "[\"topic\", \"tone\"]",
  "changeLog": "初始版本",
  "createdBy": "user123",
  "username": "管理员",
  "isPublish": true
}
```

**业务逻辑**:
- 当 `isPublish=true` 时，自动更新对应 Prompt 的 `latestVersion` 和 `isPublish` 字段

**响应参数**:

| 字段 | 类型 | 描述 |
|------|------|------|
| id | string | 版本ID |
| promptId | string | 所属提示词ID |
| version | string | 版本号 |
| content | string | 提示词内容 |
| variables | string | 变量定义 |
| isPublish | boolean | 是否发布 |
| changeLog | string | 更新日志 |
| createdBy | string | 创建者ID |
| username | string | 创建者用户名 |
| createdAt | string | 创建时间 |
| updatedAt | string | 更新时间 |

---

### 获取版本

**接口**: `GET /api/v1/version/info/:id`

**路径参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 版本ID |

---

### 获取提示词的所有版本

**接口**: `GET /api/v1/version/prompt/:promptId`

**路径参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| promptId | string | 是 | 提示词ID |

---

### 获取提示词的最新版本

**接口**: `GET /api/v1/version/prompt/:promptId/latest`

**路径参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| promptId | string | 是 | 提示词ID |

---

### 更新版本

**接口**: `POST /api/v1/version/update`

**请求参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 版本ID |
| version | string | 是 | 版本号 |
| content | string | 是 | 提示词内容 |
| variables | string | 否 | 变量定义 |
| changeLog | string | 否 | 更新日志 |
| isPublish | boolean | 否 | 是否发布 |

**请求示例**:
```json
{
  "id": "version-xxx",
  "version": "1.1.0",
  "content": "更新后的提示词内容...",
  "variables": "[\"topic\", \"tone\", \"audience\"]",
  "changeLog": "优化提示词",
  "isPublish": true
}
```

---

### 删除版本

**接口**: `POST /api/v1/version/delete/:id`

**路径参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 版本ID |

---

### 版本列表

**接口**: `GET /api/v1/version/list`

**查询参数**:

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| offset | int | 否 | 偏移量 (默认0) |
| limit | int | 否 | 限制数量 (默认10) |

---

## 错误码说明

| 错误码 | 描述 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 500 | 服务器内部错误 |

---

## 完整工作流示例

### 1. 创建用户
```bash
curl -X POST http://localhost:8080/api/v1/user/create \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "nickname": "管理员", "department": "技术部"}'
```

### 2. 创建提示词
```bash
curl -X POST http://localhost:8080/api/v1/prompt/create \
  -H "Content-Type: application/json" \
  -d '{"name": "文案生成助手", "createdBy": "admin", "username": "管理员"}'
```

### 3. 创建并发布版本
```bash
curl -X POST http://localhost:8080/api/v1/version/create \
  -H "Content-Type: application/json" \
  -d '{
    "promptId": "xxx-xxx-xxx",
    "version": "1.0.0",
    "content": "你是一个专业的文案生成助手...",
    "createdBy": "admin",
    "username": "管理员",
    "isPublish": true
  }'
```

### 4. 获取发布内容
```bash
curl -X GET "http://localhost:8080/api/v1/prompt/content?path=文案生成助手"
```

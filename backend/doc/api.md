# We-Dear API 文档

## 基础信息

- 基础URL: `http://localhost:8080/api`
- 所有请求和响应均使用 JSON 格式
- 认证方式: Bearer Token (除登录接口外，所有接口都需要在请求头中携带 token)

## 认证相关

### 登录

```http
POST /login
```

**请求参数:**

| 参数名   | 类型   | 必填 | 描述   |
|----------|--------|------|--------|
| username | string | 是   | 用户名 |
| password | string | 是   | 密码   |

**响应示例:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "1234567890",
    "username": "doctor1",
    "name": "张医生",
    "role": "doctor",
    "avatar": "/uploads/avatar.jpg"
  }
}
```

### 修改密码

```http
POST /change-password
```

**请求参数:**

| 参数名      | 类型   | 必填 | 描述                           |
|-------------|--------|------|--------------------------------|
| userId      | string | 否   | 用户ID (管理员修改他人密码时必填) |
| oldPassword | string | 否   | 原密码 (修改自己密码时必填)     |
| newPassword | string | 是   | 新密码                         |

## 医生管理

### 获取所有医生

```http
GET /doctors
```

**响应示例:**

```json
[
  {
    "id": "1234567890",
    "name": "张医生",
    "username": "doctor1",
    "departmentId": "dept1",
    "role": "doctor",
    "status": "active",
    "avatar": "/uploads/avatar.jpg"
  }
]
```

### 创建医生

```http
POST /doctors
```

**权限要求:** 管理员

**请求参数:**

| 参数名       | 类型   | 必填 | 描述     |
|--------------|--------|------|----------|
| name         | string | 是   | 姓名     |
| username     | string | 是   | 用户名   |
| password     | string | 是   | 密码     |
| departmentId | string | 是   | 科室ID   |
| avatar       | string | 否   | 头像URL  |

### 更新医生信息

```http
PUT /doctors/:id
```

### 删除医生

```http
DELETE /doctors/:id
```

**权限要求:** 管理员

## 患者管理

### 获取所有患者

```http
GET /patients
```

### 获取患者详情

```http
GET /patients/:id
```

### 创建患者

```http
POST /patients
```

**请求参数:**

| 参数名       | 类型   | 必填 | 描述     |
|--------------|--------|------|----------|
| name         | string | 是   | 姓名     |
| gender       | string | 是   | 性别     |
| birthDate    | string | 是   | 出生日期 |
| phone        | string | 是   | 电话     |
| address      | string | 否   | 地址     |
| doctorId     | string | 是   | 主治医生ID |

## 科室管理

### 获取所有科室

```http
GET /departments
```

### 创建科室

```http
POST /departments
```

**权限要求:** 管理员

**请求参数:**

| 参数名 | 类型   | 必填 | 描述     |
|--------|--------|------|----------|
| name   | string | 是   | 科室名称 |
| code   | string | 是   | 科室代码 |

### 更新科室

```http
PUT /departments/:id
```

### 删除科室

```http
DELETE /departments/:id
```

## 聊天相关

### 获取聊天列表

```http
GET /chat/list
```

### 获取聊天历史

```http
GET /chat/:patientId
```

### 医生发送消息

```http
POST /chat/:patientId/doctor
```

**请求参数:**

| 参数名    | 类型   | 必填 | 描述     |
|-----------|--------|------|----------|
| content   | string | 是   | 消息内容 |
| sender    | string | 是   | 发送者ID |

### 患者发送消息

```http
POST /chat/:patientId/patient
```

### 获取AI建议

```http
GET /chat/:patientId/suggestions
```

**查询参数:**

| 参数名    | 类型   | 描述     |
|-----------|--------|----------|
| messageId | string | 消息ID   |

## 随访记录

### 获取随访记录

```http
GET /patients/:id/followup
```

### 创建随访记录

```http
POST /followup
```

### 更新随访记录

```http
PUT /followup/:id
```

### 删除随访记录

```http
DELETE /followup/:id
```

## 医疗记录

### 获取医疗记录

```http
GET /patients/:id/medical
```

### 创建医疗记录

```http
POST /medical
```

### 更新医疗记录

```http
PUT /medical/:id
```

### 删除医疗记录

```http
DELETE /medical/:id
```

## 文件上传

### 上传文件

```http
POST /upload
```

**请求格式:** multipart/form-data

**请求参数:**

| 参数名 | 类型 | 必填 | 描述     |
|--------|------|------|----------|
| file   | file | 是   | 文件数据 |

**响应示例:**

```json
{
  "url": "/uploads/filename.jpg"
}
```

## 状态码说明

| 状态码 | 描述                |
|--------|-------------------|
| 200    | 请求成功           |
| 201    | 创建成功           |
| 400    | 请求参数错误       |
| 401    | 未授权/token无效   |
| 403    | 权限不足           |
| 404    | 资源不存在         |
| 500    | 服务器内部错误     |
# 慢病管理系统Demo

这是一个基于Vue 3 + Vite和Go Gin的慢病管理系统演示项目。

## 功能特点

- 医生可以查看患者列表
- 查看患者档案信息
- 与患者进行文字、图片、语音交流
- AI辅助回复建议
- 患者信息从CSV文件读取

## 项目结构

```
.
├── frontend/          # Vue 3 + Vite前端项目
├── backend/           # Go Gin后端项目
└── data/             # 数据文件目录
    └── patients.csv  # 患者数据
```

## 运行项目

### 前端

```bash
cd frontend
npm install
npm run dev
```

前端将在 http://localhost:5173 运行

### 后端

```bash
cd backend
go mod tidy
go run main.go
```

后端API将在 http://localhost:8080 运行

## API接口

- GET /api/patients - 获取所有患者列表
- GET /api/patients/:id - 获取指定患者信息
- POST /api/patients/:id/messages - 发送消息
- POST /api/upload - 上传文件

## 技术栈

### 前端
- Vue 3
- Vite
- TypeScript
- Element Plus
- Vue Router
- Pinia

### 后端
- Go
- Gin
- CORS 
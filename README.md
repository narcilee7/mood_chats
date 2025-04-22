# A mood ChatBot

## Tech Stack

Go + MongoDB + Vue.js + 科大讯飞AI Cloud

## Project Structure

1. text and voice chatbot (WebRTC + 讯飞API)
2. 内置情绪识别、用户画像系统: NLP分析 + 标签建模
3. 支持不同的“AI人格”
4. 消息流系统、上下文记忆系统
5. 自建算法分析语义+情感倾向
6. 情绪时序图谱
7. 用户标签管理（偏好、发言频率、情感倾向）
8. 实现上下文会话引擎：
   1. 多会话并发
   2. 会话上下文存储结构
   3. 自动清理机制
   4. 支持流式

### 服务端

1. Viper配置管理
2. Zap日志管理
3. Gin Prometheus监控
4. Gorm数据库 MongoDB
5. 科大讯飞AI接口

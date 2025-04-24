export type TextAreaType = 'index' | 'chat'

// 消息角色类型
export type Role = 'user' | 'assistant' | 'system'

// 消息类型
export interface Message {
  role: Role
  content: string
  timestamp: number
}

// 情绪分析结果类型
export interface Emotion {
  type: string      // 情绪类型：happy, sad, angry, etc.
  score: number     // 情绪强度
  keywords: string[] // 关键词
}

// 会话类型
export interface Session {
  id: string
  userId: string
  title: string
  messages: Message[]
  createdAt: number
  updatedAt: number
}

// 用户配置类型
export interface UserProfile {
  id: string
  emotionTags: Emotion[]
  sessionCount: number
  lastActive: number
  githubConfig: GithubConfig
}

// GitHub 配置类型
export interface GithubConfig {
  clientId: string
  clientSecret: string
  redirectUrl: string
  rawData: Record<string, any>
  provider: string
  email: string
  name: string
  firstName: string
  lastName: string
  nickName: string
  description: string
  userId: string
  avatarUrl: string
  location: string
  accessToken: string
  accessTokenSecret: string
  refreshToken: string
  expiresAt: Date
  idToken: string
} 
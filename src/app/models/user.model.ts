// src/app/models/user.model.ts
export interface User {
    id: number;           // 用户ID
    username: string;     // 用户名
    email: string;        // 电子邮箱
    created_at: Date;     // 创建时间
  }
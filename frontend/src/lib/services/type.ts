export interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}

// Auth Service - Token Management
export interface UserInfo {
  username: string;
  nickname: string;
  id?: string;
}

export interface Prompt {
  id: string;
  name: string;
  path: string;
  latestVersion: string;
  isPublish: boolean;
  createBy: string;
  username: string;
  createAt: string;
  updateAt: string;
  isFavorite?: boolean;
}

export interface Category {
  id: string;
  title: string;
  icon: string;
  count: number;
  url: string;
}

export interface PromptList {
  list: Prompt[];
  page: number;
  total: number;
  limit: number;
}
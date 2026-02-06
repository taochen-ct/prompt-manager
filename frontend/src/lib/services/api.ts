const API_PREFIX = '/api/v1';

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

export const auth = {
  getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('auth_token');
    }
    return null;
  },

  setToken(token: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token);
    }
  },

  clearToken(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  },

  isLoggedIn(): boolean {
    return !!this.getToken();
  },

  getUser(): UserInfo | null {
    if (typeof window !== 'undefined') {
      const userStr = localStorage.getItem('auth_user');
      if (userStr) {
        return JSON.parse(userStr)
      }
      return null;
    }
    return null;
  },

  setUser(user: UserInfo): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_user', JSON.stringify(user));
    }
  },

  clearUser(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_user');
    }
  },

  clearAll(): void {
    this.clearToken();
    this.clearUser();
  }
};

// Endpoints that don't require authorization
const PUBLIC_ENDPOINTS = ['/user/login', '/user/create', '/ping', '/prompt/content'];

class ApiService {
  private async request<T>(
      endpoint: string,
      options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_PREFIX}${endpoint}`;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options.headers as Record<string, string>,
    };

    // Add Authorization header if not a public endpoint
    const token = auth.getToken();
    if (token && !PUBLIC_ENDPOINTS.some(p => endpoint.startsWith(p))) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    const result: ApiResponse<T> = await response.json();

    if (result.code !== 0) {
      throw new Error(result.message || 'API Error');
    }

    return result.data;
  }

  // User APIs
  async login(data: { username: string; password: string }) {
    const result = await this.request<{ token: string; user: UserInfo; expireAt: string }>('/user/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    // Save token and user info
    auth.setToken(result.token);
    auth.setUser(result.user);
    return result;
  }

  async logout() {
    const user = auth.getUser();
    if (user) {
      try {
        await this.request('/user/logout', {
          method: 'POST',
          body: JSON.stringify({username: user.username}),
        });
      } catch (e) {
        // Ignore logout error, still clear local data
      }
    }
    auth.clearAll();
  }

  async createUser(data: {
    username: string;
    password: string;
    nickname?: string;
    department?: string;
  }) {
    return this.request('/user/create', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getUser(id: string) {
    return this.request(`/user/info/${id}`);
  }

  async updateUser(id: string, data: { nickname?: string; department?: string }) {
    return this.request(`/user/update/${id}`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async deleteUser(id: string) {
    return this.request('/user/delete', {
      method: 'POST',
      body: JSON.stringify({id}),
    });
  }

  // Prompt APIs
  async createPrompt(data: {
    name: string;
    createdBy: string;
    username: string;
    path?: string;
    category?: string;
  }): Promise<Prompt> {
    return this.request('/prompt/create', {
      method: 'POST',
      body: JSON.stringify({
        ...data,
        path: data.path || data.name,
      }),
    });
  }

  async getPrompt(id: string) {
    return this.request(`/prompt/info/${id}`);
  }

  async getPromptContent(path: string) {
    return this.request(`/prompt/content?path=${encodeURIComponent(path)}`);
  }

  async updatePrompt(data: {
    id: string;
    name: string;
    isPublish?: boolean;
    category?: string;
  }) {
    return this.request('/prompt/update', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async deletePrompt(id: string) {
    return this.request(`/prompt/delete/${id}`, {
      method: 'POST',
    });
  }

  async getPromptList(params: { offset?: number; limit?: number }): Promise<PromptList> {

    const query = new URLSearchParams();
    query.set('username', auth.getUser()?.username as string);
    if (params.offset !== undefined) query.set('offset', params.offset.toString());
    if (params.limit !== undefined) query.set('limit', params.limit.toString());

    const endpoint = `/prompt/list?${query.toString()}`;
    return this.request(endpoint);
  }

  // Version APIs
  async createVersion(data: {
    promptId: string;
    version: string;
    content: string;
    variables?: string | null;
    changeLog?: string | null;
    createdBy: string;
    username: string;
    isPublish?: boolean;
  }) {
    return this.request('/version/create', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getVersion(id: string) {
    return this.request(`/version/info/${id}`);
  }

  async getPromptVersions(promptId: string) {
    return this.request(`/version/prompt/${promptId}`);
  }

  async getLatestVersion(promptId: string) {
    return this.request(`/version/prompt/${promptId}/latest`);
  }

  async updateVersion(data: {
    id: string;
    version: string;
    content: string;
    variables?: string;
    changeLog?: string;
    isPublish?: boolean;
  }) {
    return this.request('/version/update', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async deleteVersion(id: string) {
    return this.request(`/version/delete/${id}`, {
      method: 'POST',
    });
  }

  async getVersionList(params: { offset?: number; limit?: number } = {}) {
    const query = new URLSearchParams();
    if (params.offset !== undefined) query.set('offset', params.offset.toString());
    if (params.limit !== undefined) query.set('limit', params.limit.toString());

    const endpoint = `/version/list${query.toString() ? `?${query.toString()}` : ''}`;
    return this.request(endpoint);
  }

  // Favorites APIs
  async addFavorite(data: { promptId: string }) {
    return this.request('/favorites/add', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async removeFavorite(data: { promptId: string }) {
    return this.request('/favorites/remove', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async checkFavorite(data: { promptId: string }) {
    return this.request('/favorites/check', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getFavoritesList(params: { offset?: number; limit?: number } = {}) {
    const query = new URLSearchParams();
    if (params.offset !== undefined) query.set('offset', params.offset.toString());
    if (params.limit !== undefined) query.set('limit', params.limit.toString());

    const endpoint = `/favorites/list${query.toString() ? `?${query.toString()}` : ''}`;
    return this.request(endpoint);
  }

  // Recently Used APIs
  async recordRecentlyUsed(data: { promptId: string }) {
    return this.request('/recently-used/record', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getRecentlyUsedList(params: { offset?: number; limit?: number } = {}) {
    const query = new URLSearchParams();
    if (params.offset !== undefined) query.set('offset', params.offset.toString());
    if (params.limit !== undefined) query.set('limit', params.limit.toString());

    const endpoint = `/recently-used/list${query.toString() ? `?${query.toString()}` : ''}`;
    return this.request(endpoint);
  }

  async removeRecentlyUsed(data: { promptId: string }) {
    return this.request('/recently-used/remove', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async cleanRecentlyUsed(data: { keepCount?: number } = {}) {
    return this.request('/recently-used/clean', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Category APIs
  async createCategory(data: {
    id: string;
    title: string;
    icon: string;
    url: string;
    createdBy: string;
    username: string;
  }) {
    return this.request('/category/create', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getCategory(id: string) {
    return this.request(`/category/info/${id}`);
  }

  async getCategoryList(): Promise<Category[]> {
    return this.request('/category/list');
  }

  async updateCategory(data: {
    id: string;
    title: string;
    icon: string;
    count?: number;
    url: string;
  }) {
    return this.request('/category/update', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async deleteCategory(id: string) {
    return this.request(`/category/delete/${id}`, {
      method: 'POST',
    });
  }

  // Category-based queries
  async getPromptsByCategory(category: string): Promise<any[]> {
    const allPrompts = await this.getPromptList({offset: 0, limit: 100});
    return allPrompts.list.filter(p =>
        p.path.toLowerCase().includes(category.toLowerCase()) ||
        p.name.toLowerCase().includes(category.toLowerCase())
    );
  }
}

export const api = new ApiService();

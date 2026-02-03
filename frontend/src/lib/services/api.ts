const API_BASE = 'http://localhost:8080/api/v1';

interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}

class ApiService {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE}${endpoint}`;

    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    const result: ApiResponse<T> = await response.json();

    if (result.code !== 0) {
      throw new Error(result.message || 'API Error');
    }

    return result.data;
  }

  // User APIs
  async createUser(data: {
    username: string;
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
      body: JSON.stringify({ id }),
    });
  }

  // Prompt APIs
  async createPrompt(data: {
    name: string;
    createdBy: string;
    username: string;
    path?: string;
  }) {
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

  async getPromptList(params: { offset?: number; limit?: number } = {}) {
    const query = new URLSearchParams();
    if (params.offset !== undefined) query.set('offset', params.offset.toString());
    if (params.limit !== undefined) query.set('limit', params.limit.toString());

    const endpoint = `/prompt/list${query.toString() ? `?${query.toString()}` : ''}`;
    return this.request(endpoint);
  }

  // Version APIs
  async createVersion(data: {
    promptId: string;
    version: string;
    content: string;
    variables?: string;
    changeLog?: string;
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

  // Favorites APIs (模拟数据)
  private favorites: string[] = [];

  async getFavorites(): Promise<any[]> {
    const allPrompts = await this.getPromptList({ offset: 0, limit: 100 }) as any[];
    return allPrompts.filter(p => this.favorites.includes(p.id)).map(p => ({
      ...p,
      isFavorite: true
    }));
  }

  async toggleFavorite(promptId: string): Promise<boolean> {
    const index = this.favorites.indexOf(promptId);
    if (index > -1) {
      this.favorites.splice(index, 1);
      return false;
    } else {
      this.favorites.push(promptId);
      return true;
    }
  }

  async isFavorite(promptId: string): Promise<boolean> {
    return this.favorites.includes(promptId);
  }

  // Recent APIs (模拟数据，存储最近访问的prompt)
  private recentPrompts: { id: string; timestamp: number }[] = [];

  async getRecentPrompts(limit: number = 10): Promise<any[]> {
    const allPrompts = await this.getPromptList({ offset: 0, limit: 100 }) as any[];
    const recentIds = this.recentPrompts
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, limit)
      .map(r => r.id);

    return recentIds
      .map(id => allPrompts.find(p => p.id === id))
      .filter(Boolean);
  }

  async addToRecent(promptId: string) {
    const existing = this.recentPrompts.findIndex(r => r.id === promptId);
    if (existing > -1) {
      this.recentPrompts.splice(existing, 1);
    }
    this.recentPrompts.unshift({ id: promptId, timestamp: Date.now() });
    // 只保留最近50条
    if (this.recentPrompts.length > 50) {
      this.recentPrompts = this.recentPrompts.slice(0, 50);
    }
  }

  // Category APIs (基于path过滤)
  async getPromptsByCategory(category: string): Promise<any[]> {
    const allPrompts = await this.getPromptList({ offset: 0, limit: 100 }) as any[];
    return allPrompts.filter(p =>
      p.path.toLowerCase().includes(category.toLowerCase()) ||
      p.name.toLowerCase().includes(category.toLowerCase())
    );
  }
}

export const api = new ApiService();

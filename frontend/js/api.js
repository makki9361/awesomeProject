const API_BASE_URL = 'http://localhost:8080/api';

class API {
    constructor() {
        this.getHeaders = () => {
            const headers = {
                'Content-Type': 'application/json'
            };
            const userId = localStorage.getItem('userId');
            if (userId) {
                headers['X-User-Id'] = userId;
            }
            return headers;
        };
    }

    async request(endpoint, options = {}) {
        try {
            const url = `${API_BASE_URL}${endpoint}`;
            console.log('API Request:', url, options);

            const response = await fetch(url, {
                ...options,
                headers: {
                    ...this.getHeaders(),
                    ...options.headers
                }
            });

            if (!response.ok) {
                let errorMessage;
                try {
                    const error = await response.json();
                    errorMessage = error.error || `HTTP ${response.status}`;
                } catch {
                    errorMessage = `HTTP ${response.status}: ${response.statusText}`;
                }
                throw new Error(errorMessage);
            }

            if (response.status === 204) {
                return null;
            }

            const data = await response.json();
            console.log('API Response:', data);
            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    async getUsers() {
        return this.request('/login');
    }

    async getCurrentUser() {
        return this.request('/users/me');
    }

    async getCategories(limit = 100, offset = 0) {
        return this.request(`/categories?limit=${limit}&offset=${offset}`);
    }

    async getCategory(id) {
        return this.request(`/categories/${id}`);
    }

    async createCategory(data) {
        return this.request('/categories', {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    async updateCategory(id, data) {
        return this.request(`/categories/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    async deleteCategory(id) {
        return this.request(`/categories/${id}`, {
            method: 'DELETE'
        });
    }

    async getRules(params = {}) {
        const queryParams = new URLSearchParams();
        if (params.category_id) queryParams.append('category_id', params.category_id);
        if (params.status) queryParams.append('status', params.status);
        if (params.search) queryParams.append('search', params.search);
        if (params.page) queryParams.append('page', params.page);
        if (params.page_size) queryParams.append('page_size', params.page_size);

        const query = queryParams.toString();
        return this.request(`/rules${query ? '?' + query : ''}`);
    }

    async getRule(id) {
        return this.request(`/rules/${id}`);
    }

    async createRule(data) {
        return this.request('/rules', {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    async updateRule(id, data) {
        return this.request(`/rules/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    async deleteRule(id) {
        return this.request(`/rules/${id}`, {
            method: 'DELETE'
        });
    }

    async publishRule(id, data) {
        return this.request(`/rules/${id}/publish`, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    async getUsersList(limit = 100, offset = 0) {
        return this.request(`/users?limit=${limit}&offset=${offset}`);
    }

    async getUser(id) {
        return this.request(`/users/${id}`);
    }

    async createUser(data) {
        return this.request('/users', {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    async updateUser(id, data) {
        return this.request(`/users/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    async deleteUser(id) {
        return this.request(`/users/${id}`, {
            method: 'DELETE'
        });
    }
}

const api = new API();
import { PUBLIC_API_URL } from '$env/static/public';
import { apiClient } from './client';
import type { 
  SuperUserDashboardStats, 
  Student, 
  Category, 
  StudentStatusType, 
  Invoice, 
  ActivityLog,
  User,
  Role,
  Notification
} from '$lib/types';

export const superUserApi = {
  getDashboard: (year?: number, range?: string) => {
    const params: Record<string, string> = {};
    if (year) params.year = String(year);
    if (range) params.range = range;
    const queryString = Object.keys(params).length ? '?' + new URLSearchParams(params).toString() : '';
    return apiClient.get<SuperUserDashboardStats>(`/api/v1/admin/dashboard${queryString}`);
  },
  
  getStudents: (search?: string) =>
    apiClient.get<Student[]>(
      `/api/v1/admin/students${search ? `?search=${search}` : ''}`
    ),
  getStudentsPaginated: (search?: string, page?: number, limit?: number) => {
    const params: Record<string, string> = search ? { search } : {};
    if (page) params.page = String(page);
    if (limit) params.limit = String(limit);
    const queryString = '?' + new URLSearchParams(params).toString();
    return apiClient.get<{ students: Student[]; pagination: { page: number; limit: number; total: number; pages: number } }>(`/api/v1/admin/students/paginated${queryString}`);
  },
  createStudent: (data: unknown) => apiClient.post('/api/v1/admin/students', data),
  updateStudent: (id: string, data: unknown) => apiClient.put(`/api/v1/admin/students/${id}`, data),
  deleteStudent: (id: string) => apiClient.delete(`/api/v1/admin/students/${id}`),
  toggleStudentStatus: (id: string) => apiClient.patch(`/api/v1/admin/students/${id}/toggle-status`),

  getCategories: () => apiClient.get<Category[]>('/api/v1/admin/categories'),
  createCategory: (data: unknown) => apiClient.post('/api/v1/admin/categories', data),
  updateCategory: (id: string, data: unknown) => apiClient.put(`/api/v1/admin/categories/${id}`, data),
  deleteCategory: (id: string) => apiClient.delete(`/api/v1/admin/categories/${id}`),

  getStatusTypes: () => apiClient.get<StudentStatusType[]>('/api/v1/admin/status-types'),
  createStatusType: (data: unknown) => apiClient.post('/api/v1/admin/status-types', data),
  updateStatusType: (id: string, data: unknown) => apiClient.put(`/api/v1/admin/status-types/${id}`, data),
  deleteStatusType: (id: string) => apiClient.delete(`/api/v1/admin/status-types/${id}`),

  generateInvoices: (hijriMonth: number, hijriYear: number) => {
    const d = new Date();
    return apiClient.post('/api/v1/admin/invoices/generate', {
      month: d.getMonth() + 1,
      year: d.getFullYear(),
      hijri_month: hijriMonth,
      hijri_year: hijriYear
    });
  },

  generateSemesterInvoices: (semester: number, hijriYear: number) =>
    apiClient.post('/api/v1/admin/invoices/generate-semester', { semester, hijri_year: hijriYear }),
  
  getInvoices: (filters?: Record<string, string>) => {
    const params = filters ? '?' + new URLSearchParams(filters).toString() : '';
    return apiClient.get<Invoice[]>(`/api/v1/admin/invoices${params}`);
  },
  getInvoicesPaginated: (filters?: Record<string, string>, page?: number, limit?: number) => {
    const params: Record<string, string> = filters ? { ...filters } : {};
    if (page) params.page = String(page);
    if (limit) params.limit = String(limit);
    const queryString = '?' + new URLSearchParams(params).toString();
    return apiClient.get<{ invoices: Invoice[]; pagination: { page: number; limit: number; total: number; pages: number } }>(`/api/v1/admin/invoices/paginated${queryString}`);
  },
  getStudentsWithInvoicesPaginated: (filters?: Record<string, string>, page?: number, limit?: number) => {
    const params: Record<string, string> = filters ? { ...filters } : {};
    if (page) params.page = String(page);
    if (limit) params.limit = String(limit);
    const queryString = '?' + new URLSearchParams(params).toString();
    return apiClient.get<{ students: Student[]; pagination: { page: number; limit: number; total: number; pages: number } }>(`/api/v1/admin/invoices/students-paginated${queryString}`);
  },

  getActivityLogs: () => apiClient.get<ActivityLog[]>('/api/v1/logs/'),
  getActivityLogsPaginated: (page?: number, limit?: number) => {
    const params: Record<string, string> = {};
    if (page) params.page = String(page);
    if (limit) params.limit = String(limit);
    const queryString = '?' + new URLSearchParams(params).toString();
    return apiClient.get<{ logs: ActivityLog[]; pagination: { page: number; limit: number; total: number; pages: number } }>(`/api/v1/logs/paginated${queryString}`);
  },
  
  getReports: (filters: Record<string, string>) => {
    const params = '?' + new URLSearchParams(filters).toString();
    return apiClient.get<unknown>(`/api/v1/reports/${params}`);
  },
  
  exportReports: (filters: Record<string, string>) => {
    const params = '?' + new URLSearchParams(filters).toString();
    return apiClient.download(`/api/v1/reports/export/csv${params}`, 'laporan_keuangan.csv');
  },

  getSettings: () => apiClient.get<unknown[]>('/api/v1/admin/settings'),
  updateSetting: (key: string, value: string) => apiClient.put('/api/v1/admin/settings', { key, value }),

  getNotifications: () => apiClient.get<Notification[]>('/api/v1/admin/notifications'),
  markNotificationRead: (id: string) => apiClient.patch(`/api/v1/admin/notifications/${id}/read`),

  // User management
  getAllUsers: () => apiClient.get<User[]>('/api/v1/admin/users'),
  getAllUsersPaginated: (page?: number, limit?: number) => {
    const params: Record<string, string> = {};
    if (page) params.page = String(page);
    if (limit) params.limit = String(limit);
    const queryString = '?' + new URLSearchParams(params).toString();
    return apiClient.get<{ users: User[]; pagination: { page: number; limit: number; total: number; pages: number } }>(`/api/v1/admin/users/paginated${queryString}`);
  },
  getAllRoles: () => apiClient.get<Role[]>('/api/v1/admin/users/roles'),
  updateUserRole: (userId: string, roleId: string) =>
    apiClient.patch(`/api/v1/admin/users/${userId}/role`, { role_id: roleId }),
  toggleUserActive: (userId: string) =>
    apiClient.patch(`/api/v1/admin/users/${userId}/toggle-active`)
};

import { apiClient } from './client';
import type { SuperUserDashboardStats, Student, Invoice, Category } from '$lib/types';

export const pengasuhApi = {
  getDashboard: (year?: number, range?: string) => {
    const params: Record<string, string> = {};
    if (year) params.year = String(year);
    if (range) params.range = range;
    const queryString = Object.keys(params).length ? '?' + new URLSearchParams(params).toString() : '';
    return apiClient.get<SuperUserDashboardStats>(`/api/v1/pengasuh/dashboard${queryString}`);
  },
  getStudents: (search?: string) =>
    apiClient.get<Student[]>(
      `/api/v1/pengasuh/students${search ? `?search=${search}` : ''}`
    ),
  getInvoices: (filters?: Record<string, string>) => {
    const params = filters ? '?' + new URLSearchParams(filters).toString() : '';
    return apiClient.get<Invoice[]>(`/api/v1/pengasuh/invoices${params}`);
  },
  getInvoicesPaginated: (filters: Record<string, string>, page: number, limit: number) => {
    const params = new URLSearchParams({
      ...filters,
      page: String(page),
      limit: String(limit)
    }).toString();
    return apiClient.get<{ invoices: Invoice[]; pagination: { page: number; limit: number; total: number; pages: number } }>(
      `/api/v1/pengasuh/invoices/paginated?${params}`
    );
  },
  getCategories: () =>
    apiClient.get<Category[]>(`/api/v1/pengasuh/categories`),
  getReports: (filters: Record<string, string>) => {
    const params = '?' + new URLSearchParams(filters).toString();
    return apiClient.get<unknown>(`/api/v1/reports/${params}`);
  },
  exportReports: (filters: Record<string, string>) => {
    const params = '?' + new URLSearchParams(filters).toString();
    return apiClient.download(`/api/v1/reports/export/csv${params}`, 'laporan_keuangan.csv');
  }
};

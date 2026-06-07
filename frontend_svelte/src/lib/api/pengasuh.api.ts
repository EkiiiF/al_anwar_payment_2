import { apiClient } from './client';
import type { SuperUserDashboardStats, Student, Invoice } from '$lib/types';

export const pengasuhApi = {
  getDashboard: () =>
    apiClient.get<SuperUserDashboardStats>('/api/v1/pengasuh/dashboard'),
  getStudents: (search?: string) =>
    apiClient.get<Student[]>(
      `/api/v1/pengasuh/students${search ? `?search=${search}` : ''}`
    ),
  getInvoices: (filters?: Record<string, string>) => {
    const params = filters ? '?' + new URLSearchParams(filters).toString() : '';
    return apiClient.get<Invoice[]>(`/api/v1/pengasuh/invoices${params}`);
  },
  getReports: (filters: Record<string, string>) => {
    const params = '?' + new URLSearchParams(filters).toString();
    return apiClient.get<unknown>(`/api/v1/reports/${params}`);
  },
  exportReports: (filters: Record<string, string>) => {
    const params = '?' + new URLSearchParams(filters).toString();
    return apiClient.download(`/api/v1/reports/export/csv${params}`, 'laporan_keuangan.csv');
  }
};

import { apiClient } from './client';
import type { 
  GuardianDashboardStats, 
  Invoice, 
  Payment, 
  Notification, 
  User 
} from '$lib/types';

export const guardianApi = {
  getDashboard: () => apiClient.get<GuardianDashboardStats>('/api/v1/guardian/dashboard'),
  getInvoices: () => apiClient.get<Invoice[]>('/api/v1/guardian/invoices'),
  getPaymentHistory: () => apiClient.get<Payment[]>('/api/v1/guardian/payments'),
  getNotifications: () => apiClient.get<Notification[]>('/api/v1/guardian/notifications'),
  markNotificationRead: (id: string) => apiClient.patch(`/api/v1/guardian/notifications/${id}/read`),
  createPayment: (invoiceIds: string[]) =>
    apiClient.post<{ token: string; redirect_url: string }>('/api/v1/payments/create', {
      invoice_ids: invoiceIds
    }),
  checkPaymentStatus: (orderId: string) => apiClient.get(`/api/v1/payments/status/${orderId}`),
  getProfile: () => apiClient.get<User>('/api/v1/auth/profile'),
  updateProfile: (data: unknown) => apiClient.put('/api/v1/auth/profile', data),
  changePassword: (data: unknown) => apiClient.put('/api/v1/auth/change-password', data)
};

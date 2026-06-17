import type { User } from './auth.types';
import type { Student } from './student.types';

export type InvoiceStatus = 'unpaid' | 'pending' | 'paid' | 'expired' | 'cancelled';

export interface Category {
  id: string;
  name: string;
  base_amount: number;
  description?: string;
}

export interface Invoice {
  id: string;
  student_id: string;
  student?: Student;
  category_id?: string;
  category?: Category;
  invoice_number: string;
  month: number;
  year: number;
  hijri_month: number;
  hijri_year: number;
  semester: number;
  academic_year_label: string;
  amount_due: number;
  status: InvoiceStatus;
  due_date: string;
  notified_at?: string;
  created_at: string;
  payments?: Payment[];
}

export interface Payment {
  id: string;
  invoice_id: string;
  invoice?: Invoice;
  external_id?: string;
  snap_token?: string;
  amount_paid: number;
  payment_method?: string;
  transaction_status?: string;
  payment_date?: string;
  evidence_url?: string;
  created_at: string;
}

export interface MonthlyPaymentStat {
  month: number;
  year: number;
  total: number;
}

export interface SemesterPaymentStat {
  semester_name: string;
  academic_year_label: string;
  total: number;
  invoice_count: number;
  paid_count: number;
  unpaid_count: number;
}

export interface HijriMonthInfo {
  hijri_month: number;
  hijri_month_name: string;
  hijri_year: number;
  semester: number;
  semester_name: string;
  academic_year_label: string;
  is_exam_month: boolean;
  is_registration: boolean;
}

export interface SuperUserDashboardStats {
  total_students: number;
  unpaid_invoices: number;
  paid_invoices: number;
  total_income_mo: number;
  monthly_payments: MonthlyPaymentStat[];
  semester_stats: SemesterPaymentStat[];
  current_hijri: HijriMonthInfo;
}

export interface GuardianStudentSummary {
  student_id: string;
  student_name: string;
  nis: string;
  category_name: string;
  status_name: string;
  is_active: boolean;
}

export interface GuardianDashboardStats {
  student: Student;
  total_unpaid: number;
  unpaid_count: number;
  total_invoices: number;
  recent_payments: Payment[];
}

export interface ActivityLog {
  id: string;
  user_id: string;
  user?: User;
  action: string;
  entity_name?: string;
  entity_id?: string;
  payload: string;
  ip_address?: string;
  user_agent?: string;
  created_at: string;
}

export interface Notification {
  id: string;
  user_id: string;
  type?: string;
  title: string;
  message: string;
  is_read: boolean;
  created_at: string;
}

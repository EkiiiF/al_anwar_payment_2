import type { User } from './auth.types';

export interface StudentStatusType {
  id: string;
  name: string;
  discount_percentage: number;
  is_active_billing: boolean;
  is_default: boolean;
  description: string;
}

export interface Category {
  id: string;
  name: string;
  base_amount: number;
  is_fixed: boolean;
  is_active: boolean;
  is_semester?: boolean;
  description: string;
}

export interface StudentName {
  first_name: string;
  middle_name?: string;
  last_name?: string;
}

export interface GuardianName {
  first_name: string;
  middle_name?: string;
  last_name?: string;
}

export interface Address {
  id: string;
  student_id: string;
  address_line: string;
  village: string;
  district: string;
  city: string;
  province: string;
  country: string;
  postal_code: string;
  is_primary: boolean;
}

export interface Student {
  id: string;
  name: StudentName;
  nik: string;
  birth_date: string;
  student_number: string;
  gender: string;
  muhadhoroh_level: number;  // 1-9 = Muhadhoroh, 0 = Lulus
  current_semester: number;  // 1 atau 2
  is_active: boolean;
  status_id: string;
  status?: StudentStatusType;
  guardians?: Guardian[];
  addresses?: Address[];
  invoices?: import('./invoice.types').Invoice[];
  billing_categories?: Category[];
  full_name?: string;
  created_at: string;
  updated_at: string;
}

export interface Guardian {
  id: string;
  student_id: string;
  student?: Student;
  user_id: string;
  user?: User;
  name: GuardianName;
  phone: string;
  email: string;
  relation: string;
  created_at: string;
  updated_at: string;
}

/**
 * Format label kelas muhadhoroh
 * @example getMuhadhorohLabel(3, 2) => "Muhadhoroh 3 Semester 2"
 * @example getMuhadhorohLabel(0, 0) => "Lulus"
 */
export function getMuhadhorohLabel(level: number, semester: number): string {
  if (level === 0) return 'Lulus';
  return `Muhadhoroh ${level} Semester ${semester}`;
}

/** Options untuk select kelas Muhadhoroh */
export const MUHADHOROH_OPTIONS = [
  { value: '1', label: 'Muhadhoroh 1' },
  { value: '2', label: 'Muhadhoroh 2' },
  { value: '3', label: 'Muhadhoroh 3' },
  { value: '4', label: 'Muhadhoroh 4' },
  { value: '5', label: 'Muhadhoroh 5' },
  { value: '6', label: 'Muhadhoroh 6' },
  { value: '7', label: 'Muhadhoroh 7' },
  { value: '8', label: 'Muhadhoroh 8' },
  { value: '9', label: 'Muhadhoroh 9' },
  { value: '0', label: 'Lulus' },
];

/** Options untuk select hubungan wali */
export const GUARDIAN_RELATION_OPTIONS = [
  // { value: 'Orang Tua', label: 'Orang Tua' },
  { value: 'Ayah', label: 'Ayah' },
  { value: 'Ibu', label: 'Ibu' },
  { value: 'Kakek', label: 'Kakek' },
  { value: 'Nenek', label: 'Nenek' },
  { value: 'Paman', label: 'Paman' },
  { value: 'Bibi', label: 'Bibi' },
  { value: 'Kakak', label: 'Kakak' },
  { value: 'Saudara', label: 'Saudara' },
  // { value: 'Wali', label: 'Wali' },
  // { value: 'Lainnya', label: 'Lainnya' },
];

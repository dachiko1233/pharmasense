export type RiskLevel = "CRITICAL" | "HIGH" | "MEDIUM" | "LOW";

export interface User {
  id: string;
  pharmacy_id: string;
  email: string;
  full_name: string;
  role: "admin" | "staff";
  is_active: boolean;
  last_login_at: string | null;
  created_at: string;
}

export interface Pharmacy {
  id: string;
  name: string;
  license_number: string;
  address: string;
  city: string;
  phone: string;
  email: string;
  language: string;
}

export interface Product {
  id: string;
  barcode: string;
  name: string;
  name_el: string;
  category: string;
  manufacturer: string;
  requires_prescription: boolean;
}

export interface InventoryBatchWithRisk {
  id: string;
  pharmacy_id: string;
  product_id: string;
  product_name: string;
  product_name_el: string;
  category: string;
  manufacturer: string;
  barcode: string;
  batch_number: string;
  expiry_date: string;
  initial_quantity: number;
  current_quantity: number;
  purchase_price: number;
  selling_price: number;
  supplier: string;
  received_date: string;
  risk_level: RiskLevel;
  days_until_expiry: number;
  avg_daily_sales: number;
  expected_sales: number;
  estimated_surplus: number;
  estimated_loss: number;
  suggested_discount_percent: number;
}

export interface DashboardData {
  critical_count: number;
  high_count: number;
  medium_count: number;
  low_count: number;
  total_estimated_loss: number;
  total_potential_savings: number;
  total_inventory_value: number;
  critical_alerts: InventoryBatchWithRisk[];
}

export interface TimelineMonth {
  month: string;
  critical: number;
  high: number;
  medium: number;
  low: number;
}

export interface TopRiskItem {
  product_name: string;
  batch_number: string;
  estimated_loss: number;
  risk_level: RiskLevel;
  category: string;
}

export interface AlertAction {
  id: string;
  batch_id: string;
  pharmacy_id: string;
  user_id: string;
  action_type: "DISCOUNT" | "TRANSFER" | "RETURN" | "DISMISS";
  discount_percent: number | null;
  notes: string;
  created_at: string;
  user_name: string;
  product_name: string;
  batch_number: string;
}

export interface SavingsReport {
  month: string;
  actions_count: number;
  estimated_saved: number;
}

export interface CategoryReport {
  category: string;
  total_batches: number;
  at_risk: number;
  estimated_loss: number;
}

export interface ListResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface InventoryFilters {
  categories: string[];
  suppliers: string[];
  risk_levels: string[];
}

import {
  AlertAction,
  CategoryReport,
  DashboardData,
  InventoryBatchWithRisk,
  InventoryFilters,
  ListResponse,
  LoginResponse,
  Pharmacy,
  Product,
  SavingsReport,
  TimelineMonth,
  TopRiskItem,
  User,
} from "@/types";

const BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001/api/v1";

function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("pharmasense_token");
}

async function request<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options?.headers as Record<string, string>),
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${BASE_URL}${path}`, { ...options, headers });

  if (res.status === 401) {
    if (typeof window !== "undefined") {
      localStorage.removeItem("pharmasense_token");
      localStorage.removeItem("pharmasense_user");
      window.location.href = "/login";
    }
    throw new Error("Unauthorized");
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error || "Request failed");
  }

  return res.json();
}

// Auth
export const authApi = {
  login: (email: string, password: string) =>
    request<LoginResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),
  me: () => request<User>("/auth/me"),
  logout: () => request("/auth/logout", { method: "POST" }),
};

// Pharmacy
export const pharmacyApi = {
  get: () => request<Pharmacy>("/pharmacy"),
  update: (data: Partial<Pharmacy>) =>
    request<Pharmacy>("/pharmacy", { method: "PATCH", body: JSON.stringify(data) }),
  getUsers: () => request<User[]>("/pharmacy/users"),
};

// Dashboard / Risk
export const riskApi = {
  dashboard: () => request<DashboardData>("/risk/dashboard"),
  timeline: () => request<TimelineMonth[]>("/risk/timeline"),
  topRisk: (limit = 10) =>
    request<TopRiskItem[]>(`/risk/top?limit=${limit}`),
  assessments: (params?: Record<string, string | number>) => {
    const q = new URLSearchParams(
      Object.entries(params || {}).map(([k, v]) => [k, String(v)])
    ).toString();
    return request<ListResponse<InventoryBatchWithRisk>>(
      `/risk/assessments${q ? `?${q}` : ""}`
    );
  },
  recalculate: () =>
    request("/risk/recalculate", { method: "POST" }),
};

// Inventory
export const inventoryApi = {
  list: (params?: Record<string, string | number>) => {
    const q = new URLSearchParams(
      Object.entries(params || {})
        .filter(([, v]) => v !== "" && v !== undefined)
        .map(([k, v]) => [k, String(v)])
    ).toString();
    return request<ListResponse<InventoryBatchWithRisk>>(
      `/inventory${q ? `?${q}` : ""}`
    );
  },
  get: (id: string) =>
    request<InventoryBatchWithRisk>(`/inventory/${id}`),
  filters: () => request<InventoryFilters>("/inventory/filters"),
  exportCSV: async () => {
    const token = getToken();
    const res = await fetch(`${BASE_URL}/inventory/export`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    if (!res.ok) throw new Error("Export failed");
    return res.blob();
  },
  templateURL: () => `${BASE_URL}/inventory/import/template`,
  importCSV: (file: File) => {
    const form = new FormData();
    form.append("file", file);
    const token = getToken();
    return fetch(`${BASE_URL}/inventory/import`, {
      method: "POST",
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: form,
    }).then((r) => r.json());
  },
  confirmImport: async (file: File): Promise<{ imported: number; errors: { row: number; error: string }[] }> => {
    const form = new FormData();
    form.append("file", file);
    const token = getToken();
    const res = await fetch(`${BASE_URL}/inventory/import/confirm`, {
      method: "POST",
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: form,
    });
    if (!res.ok) {
      const body = await res.json().catch(() => ({ error: res.statusText }));
      throw new Error(body.error || "Import failed");
    }
    return res.json();
  },
};

// Products
export const productsApi = {
  list: (params?: Record<string, string | number>) => {
    const q = new URLSearchParams(
      Object.entries(params || {})
        .filter(([, v]) => v !== "")
        .map(([k, v]) => [k, String(v)])
    ).toString();
    return request<ListResponse<Product>>(`/products${q ? `?${q}` : ""}`);
  },
};

// Alerts
export const alertsApi = {
  list: (riskLevel?: string) => {
    const q = riskLevel ? `?risk_level=${riskLevel}` : "";
    return request<{ data: InventoryBatchWithRisk[]; total: number }>(
      `/alerts${q}`
    );
  },
  takeAction: (
    batchId: string,
    action: {
      action_type: string;
      discount_percent?: number;
      notes?: string;
    }
  ) =>
    request(`/alerts/${batchId}/action`, {
      method: "POST",
      body: JSON.stringify(action),
    }),
  history: (limit = 50) =>
    request<AlertAction[]>(`/alerts/history?limit=${limit}`),
};

// Reports
export const reportsApi = {
  savings: () => request<SavingsReport[]>("/reports/savings"),
  categories: () => request<CategoryReport[]>("/reports/categories"),
  waste: () => request<SavingsReport[]>("/reports/waste"),
};

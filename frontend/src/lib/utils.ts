import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { RiskLevel } from "@/types";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat("el-CY", {
    style: "currency",
    currency: "EUR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
}

export function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

export function daysUntilExpiry(dateStr: string): number {
  const expiry = new Date(dateStr);
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  expiry.setHours(0, 0, 0, 0);
  return Math.floor((expiry.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
}

export const riskColors: Record<RiskLevel, { bg: string; text: string; border: string; light: string }> = {
  CRITICAL: {
    bg: "bg-red-600",
    text: "text-red-700",
    border: "border-red-200",
    light: "bg-red-50",
  },
  HIGH: {
    bg: "bg-orange-500",
    text: "text-orange-700",
    border: "border-orange-200",
    light: "bg-orange-50",
  },
  MEDIUM: {
    bg: "bg-yellow-500",
    text: "text-yellow-700",
    border: "border-yellow-200",
    light: "bg-yellow-50",
  },
  LOW: {
    bg: "bg-green-500",
    text: "text-green-700",
    border: "border-green-200",
    light: "bg-green-50",
  },
};

export function getRiskBadgeClass(level: RiskLevel): string {
  const map: Record<RiskLevel, string> = {
    CRITICAL: "bg-red-100 text-red-700 border-red-200",
    HIGH: "bg-orange-100 text-orange-700 border-orange-200",
    MEDIUM: "bg-yellow-100 text-yellow-700 border-yellow-200",
    LOW: "bg-green-100 text-green-700 border-green-200",
  };
  return map[level] || map.LOW;
}

export function getRiskRowClass(level: RiskLevel): string {
  const map: Record<RiskLevel, string> = {
    CRITICAL: "bg-red-50/50",
    HIGH: "bg-orange-50/50",
    MEDIUM: "bg-yellow-50/30",
    LOW: "",
  };
  return map[level] || "";
}

export function formatMonth(monthStr: string): string {
  const [year, month] = monthStr.split("-");
  return new Date(Number(year), Number(month) - 1).toLocaleDateString("en-GB", {
    month: "short",
    year: "numeric",
  });
}

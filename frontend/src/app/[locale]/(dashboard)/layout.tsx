"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useLocale } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import { Sidebar } from "@/components/layout/sidebar";
import { Header } from "@/components/layout/header";
import { isAuthenticated } from "@/lib/auth";
import { pharmacyApi } from "@/lib/api";

const pageTitles: Record<string, string> = {
  dashboard: "Dashboard",
  inventory: "Inventory",
  alerts: "Alerts",
  import: "Import",
  reports: "Reports",
  settings: "Settings",
};

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (!isAuthenticated()) {
      router.push(`/${locale}/login`);
    }
  }, [locale, router]);

  const { data: pharmacy } = useQuery({
    queryKey: ["pharmacy"],
    queryFn: pharmacyApi.get,
  });

  const currentPage = pathname.split("/").pop() || "dashboard";
  const breadcrumb = pageTitles[currentPage] || "";

  return (
    <div className="flex h-screen overflow-hidden bg-slate-50">
      <Sidebar />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Header
          pharmacyName={pharmacy?.name}
          breadcrumb={breadcrumb}
        />
        <main className="flex-1 overflow-y-auto p-6">
          {children}
        </main>
      </div>
    </div>
  );
}

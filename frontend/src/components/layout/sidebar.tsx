"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useTranslations, useLocale } from "next-intl";
import {
  LayoutDashboard,
  Package,
  Bell,
  Upload,
  BarChart3,
  Settings,
  Activity,
} from "lucide-react";
import { cn } from "@/lib/utils";

const navItems = [
  { key: "dashboard", icon: LayoutDashboard, href: "/dashboard" },
  { key: "inventory", icon: Package, href: "/inventory" },
  { key: "alerts", icon: Bell, href: "/alerts" },
  { key: "import", icon: Upload, href: "/import" },
  { key: "reports", icon: BarChart3, href: "/reports" },
  { key: "settings", icon: Settings, href: "/settings" },
];

export function Sidebar() {
  const t = useTranslations("nav");
  const locale = useLocale();
  const pathname = usePathname();

  return (
    <aside className="flex h-screen w-60 flex-col border-r border-slate-200 bg-white">
      {/* Logo */}
      <div className="flex h-16 items-center gap-2 border-b border-slate-200 px-6">
        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-600">
          <Activity className="h-4 w-4 text-white" />
        </div>
        <div>
          <span className="text-base font-bold text-slate-900">Pharma</span>
          <span className="text-base font-bold text-emerald-600">Sense</span>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 space-y-1 p-3">
        {navItems.map(({ key, icon: Icon, href }) => {
          const fullHref = `/${locale}${href}`;
          const isActive = pathname === fullHref || pathname.startsWith(fullHref + "/");
          return (
            <Link
              key={key}
              href={fullHref}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors",
                isActive
                  ? "bg-emerald-50 text-emerald-700"
                  : "text-slate-600 hover:bg-slate-50 hover:text-slate-900"
              )}
            >
              <Icon
                className={cn(
                  "h-4 w-4",
                  isActive ? "text-emerald-600" : "text-slate-400"
                )}
              />
              {t(key as keyof typeof t)}
            </Link>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="border-t border-slate-200 p-4">
        <p className="text-xs text-slate-400">PharmaSense v1.0</p>
        <p className="text-xs text-slate-400">Cyprus Pharmacy Platform</p>
      </div>
    </aside>
  );
}

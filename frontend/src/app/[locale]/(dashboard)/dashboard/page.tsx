"use client";

import { useQuery } from "@tanstack/react-query";
import { useTranslations, useLocale } from "next-intl";
import Link from "next/link";
import {
  AlertTriangle,
  TrendingDown,
  TrendingUp,
  Package,
  ArrowRight,
  Clock,
} from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { RiskBadge } from "@/components/layout/risk-badge";
import { riskApi } from "@/lib/api";
import { formatCurrency, formatDate, formatMonth } from "@/lib/utils";

function KPICard({
  title,
  value,
  subtitle,
  icon: Icon,
  color,
}: {
  title: string;
  value: string;
  subtitle?: string;
  icon: React.ElementType;
  color: string;
}) {
  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-slate-500">{title}</p>
            <p className="text-2xl font-bold text-slate-900 mt-1">{value}</p>
            {subtitle && <p className="text-xs text-slate-400 mt-1">{subtitle}</p>}
          </div>
          <div className={`flex h-12 w-12 items-center justify-center rounded-xl ${color}`}>
            <Icon className="h-6 w-6 text-white" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export default function DashboardPage() {
  const t = useTranslations("dashboard");
  const locale = useLocale();

  const { data: dashboard, isLoading: dashLoading } = useQuery({
    queryKey: ["dashboard"],
    queryFn: riskApi.dashboard,
  });

  const { data: timeline, isLoading: timelineLoading } = useQuery({
    queryKey: ["timeline"],
    queryFn: riskApi.timeline,
  });

  const { data: topRisk, isLoading: topLoading } = useQuery({
    queryKey: ["top-risk"],
    queryFn: () => riskApi.topRisk(10),
  });

  const timelineData = (timeline || []).map((m) => ({
    ...m,
    name: formatMonth(m.month),
  }));

  const topRiskData = (topRisk || []).map((item) => ({
    name:
      item.product_name.length > 22
        ? item.product_name.slice(0, 22) + "…"
        : item.product_name,
    loss: Math.round(item.estimated_loss),
    level: item.risk_level,
  }));

  return (
    <div className="space-y-6 max-w-7xl">
      {/* KPI Row */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {dashLoading ? (
          Array.from({ length: 4 }).map((_, i) => (
            <Card key={i}>
              <CardContent className="p-6">
                <Skeleton className="h-4 w-24 mb-2" />
                <Skeleton className="h-8 w-32" />
              </CardContent>
            </Card>
          ))
        ) : (
          <>
            <KPICard
              title={t("critical_items")}
              value={`${dashboard?.critical_count ?? 0} ${t("items")}`}
              subtitle={`+${dashboard?.high_count ?? 0} high risk`}
              icon={AlertTriangle}
              color="bg-red-600"
            />
            <KPICard
              title={t("estimated_loss")}
              value={formatCurrency(dashboard?.total_estimated_loss ?? 0)}
              subtitle="If no action taken"
              icon={TrendingDown}
              color="bg-red-500"
            />
            <KPICard
              title={t("potential_savings")}
              value={formatCurrency(dashboard?.total_potential_savings ?? 0)}
              subtitle="With recommended actions"
              icon={TrendingUp}
              color="bg-emerald-600"
            />
            <KPICard
              title={t("inventory_value")}
              value={formatCurrency(dashboard?.total_inventory_value ?? 0)}
              subtitle="Total stock at selling price"
              icon={Package}
              color="bg-slate-500"
            />
          </>
        )}
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Expiry Timeline */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>{t("expiry_timeline")}</CardTitle>
            <CardDescription>{t("expiry_timeline_desc")}</CardDescription>
          </CardHeader>
          <CardContent>
            {timelineLoading ? (
              <Skeleton className="h-64 w-full" />
            ) : (
              <ResponsiveContainer width="100%" height={260}>
                <BarChart data={timelineData} barSize={14}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
                  <XAxis
                    dataKey="name"
                    tick={{ fontSize: 11, fill: "#94a3b8" }}
                    tickLine={false}
                  />
                  <YAxis tick={{ fontSize: 11, fill: "#94a3b8" }} tickLine={false} axisLine={false} />
                  <Tooltip
                    contentStyle={{
                      border: "1px solid #e2e8f0",
                      borderRadius: "8px",
                      fontSize: "12px",
                    }}
                  />
                  <Legend wrapperStyle={{ fontSize: "12px" }} />
                  <Bar dataKey="critical" stackId="a" fill="#dc2626" name="Critical" radius={[0,0,0,0]} />
                  <Bar dataKey="high" stackId="a" fill="#f97316" name="High" />
                  <Bar dataKey="medium" stackId="a" fill="#eab308" name="Medium" />
                  <Bar dataKey="low" stackId="a" fill="#22c55e" name="Low" radius={[4,4,0,0]} />
                </BarChart>
              </ResponsiveContainer>
            )}
          </CardContent>
        </Card>

        {/* Top Risk Products */}
        <Card>
          <CardHeader>
            <CardTitle>{t("top_risk")}</CardTitle>
            <CardDescription>{t("top_risk_desc")}</CardDescription>
          </CardHeader>
          <CardContent>
            {topLoading ? (
              <div className="space-y-3">
                {Array.from({ length: 5 }).map((_, i) => (
                  <Skeleton key={i} className="h-8 w-full" />
                ))}
              </div>
            ) : (
              <div className="space-y-2">
                {topRiskData.slice(0, 8).map((item, i) => (
                  <div key={i} className="flex items-center gap-2">
                    <span className="w-5 text-xs text-slate-400 shrink-0">{i + 1}.</span>
                    <div className="flex-1 min-w-0">
                      <p className="text-xs font-medium text-slate-700 truncate">{item.name}</p>
                      <div className="mt-1 h-1.5 w-full rounded-full bg-slate-100">
                        <div
                          className={`h-1.5 rounded-full ${item.level === "CRITICAL" ? "bg-red-500" : "bg-orange-400"}`}
                          style={{
                            width: `${Math.min(100, (item.loss / (topRiskData[0]?.loss || 1)) * 100)}%`,
                          }}
                        />
                      </div>
                    </div>
                    <span className="text-xs font-semibold text-slate-700 shrink-0">
                      {formatCurrency(item.loss)}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Recent Critical Alerts */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>{t("recent_critical")}</CardTitle>
            <CardDescription>Items requiring immediate action</CardDescription>
          </div>
          <Button variant="outline" size="sm" asChild>
            <Link href={`/${locale}/alerts`} className="gap-2">
              {t("view_all")}
              <ArrowRight className="h-3 w-3" />
            </Link>
          </Button>
        </CardHeader>
        <CardContent>
          {dashLoading ? (
            <div className="space-y-3">
              {Array.from({ length: 3 }).map((_, i) => (
                <Skeleton key={i} className="h-16 w-full" />
              ))}
            </div>
          ) : !dashboard?.critical_alerts?.length ? (
            <div className="py-8 text-center text-slate-400">
              <Package className="h-8 w-8 mx-auto mb-2 opacity-40" />
              <p className="text-sm">{t("no_critical")}</p>
            </div>
          ) : (
            <div className="divide-y divide-slate-100">
              {dashboard.critical_alerts.map((alert) => (
                <div
                  key={alert.id}
                  className="flex items-center justify-between py-3 gap-4"
                >
                  <div className="flex items-center gap-3 min-w-0">
                    <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-red-50">
                      <AlertTriangle className="h-4 w-4 text-red-600" />
                    </div>
                    <div className="min-w-0">
                      <p className="text-sm font-medium text-slate-900 truncate">
                        {alert.product_name}
                      </p>
                      <p className="text-xs text-slate-500">
                        Batch: {alert.batch_number} · {alert.current_quantity} units
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4 shrink-0">
                    <div className="text-right hidden sm:block">
                      <div className="flex items-center gap-1 text-xs text-red-600">
                        <Clock className="h-3 w-3" />
                        {alert.days_until_expiry <= 0
                          ? t("expired")
                          : `${alert.days_until_expiry} ${t("days_left").replace("{{days}}", "")}`}
                      </div>
                      <p className="text-xs font-semibold text-slate-700 mt-0.5">
                        {formatCurrency(alert.estimated_loss)} loss
                      </p>
                    </div>
                    <RiskBadge level={alert.risk_level} />
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

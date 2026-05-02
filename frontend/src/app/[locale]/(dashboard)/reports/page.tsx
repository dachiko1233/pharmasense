"use client";

import dynamic from "next/dynamic";
import { useQuery } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { FileDown } from "lucide-react";
import { toast } from "sonner";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { reportsApi } from "@/lib/api";
import { formatCurrency, formatMonth } from "@/lib/utils";

const SavingsBarChart = dynamic(
  () => import("@/components/charts/SavingsBarChart").then((m) => m.SavingsBarChart),
  { ssr: false, loading: () => <Skeleton className="h-64 w-full" /> }
);

const ActionsLineChart = dynamic(
  () => import("@/components/charts/ActionsLineChart").then((m) => m.ActionsLineChart),
  { ssr: false, loading: () => <Skeleton className="h-48 w-full" /> }
);

export default function ReportsPage() {
  const t = useTranslations("reports");

  const { data: savings, isLoading: savLoading } = useQuery({
    queryKey: ["savings"],
    queryFn: reportsApi.savings,
  });

  const { data: categories, isLoading: catLoading } = useQuery({
    queryKey: ["categories"],
    queryFn: reportsApi.categories,
  });

  const savingsData = (savings || []).map((s) => ({
    name: formatMonth(s.month),
    saved: Math.round(s.estimated_saved),
    actions: s.actions_count,
  }));

  const handleExportPDF = () => {
    toast.info("PDF export coming soon!", { description: "This feature will be available in the next release." });
  };

  return (
    <div className="space-y-6 max-w-5xl">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-slate-500 text-sm">{t("subtitle")}</p>
        </div>
        <Button variant="outline" onClick={handleExportPDF} className="gap-2">
          <FileDown className="h-4 w-4" />
          {t("export_pdf")}
        </Button>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Savings Chart */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>{t("savings_chart")}</CardTitle>
            <CardDescription>Estimated money saved by taking action on expiring products</CardDescription>
          </CardHeader>
          <CardContent>
            {savLoading ? (
              <Skeleton className="h-64 w-full" />
            ) : savingsData.length === 0 ? (
              <div className="h-64 flex items-center justify-center text-slate-400">
                <div className="text-center">
                  <p className="text-sm">No actions taken yet</p>
                  <p className="text-xs mt-1">Take action on alerts to see savings here</p>
                </div>
              </div>
            ) : (
              <SavingsBarChart data={savingsData} />
            )}
          </CardContent>
        </Card>

        {/* Actions trend */}
        <Card>
          <CardHeader>
            <CardTitle>{t("waste_chart")}</CardTitle>
            <CardDescription>Number of actions taken per month</CardDescription>
          </CardHeader>
          <CardContent>
            {savLoading ? (
              <Skeleton className="h-48 w-full" />
            ) : savingsData.length === 0 ? (
              <div className="h-48 flex items-center justify-center text-slate-400 text-sm">
                No data yet
              </div>
            ) : (
              <ActionsLineChart data={savingsData} />
            )}
          </CardContent>
        </Card>

        {/* Category breakdown */}
        <Card>
          <CardHeader>
            <CardTitle>{t("category_table")}</CardTitle>
            <CardDescription>Categories with highest expiry risk</CardDescription>
          </CardHeader>
          <CardContent>
            {catLoading ? (
              <div className="space-y-3">
                {Array.from({ length: 5 }).map((_, i) => (
                  <Skeleton key={i} className="h-8 w-full" />
                ))}
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="text-left border-b border-slate-200">
                      <th className="pb-2 font-medium text-slate-600">{t("category")}</th>
                      <th className="pb-2 font-medium text-slate-600 text-right">{t("total_batches")}</th>
                      <th className="pb-2 font-medium text-slate-600 text-right">{t("at_risk")}</th>
                      <th className="pb-2 font-medium text-slate-600 text-right">{t("estimated_loss")}</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {(categories || []).map((cat) => (
                      <tr key={cat.category}>
                        <td className="py-3 font-medium text-slate-700">{cat.category}</td>
                        <td className="py-3 text-right text-slate-600">{cat.total_batches}</td>
                        <td className="py-3 text-right">
                          <span
                            className={`font-medium ${
                              cat.at_risk > 0 ? "text-orange-600" : "text-slate-400"
                            }`}
                          >
                            {cat.at_risk}
                          </span>
                        </td>
                        <td className="py-3 text-right">
                          <span
                            className={`font-semibold ${
                              cat.estimated_loss > 0 ? "text-red-600" : "text-slate-400"
                            }`}
                          >
                            {cat.estimated_loss > 0 ? formatCurrency(cat.estimated_loss) : "—"}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

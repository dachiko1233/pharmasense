"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import {
  AlertTriangle,
  Clock,
  Package,
  ShoppingCart,
  TrendingDown,
  Percent,
  Truck,
  RotateCcw,
  X,
  History,
} from "lucide-react";
import { toast } from "sonner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import { RiskBadge } from "@/components/layout/risk-badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { alertsApi } from "@/lib/api";
import { formatCurrency, formatDate } from "@/lib/utils";
import { InventoryBatchWithRisk } from "@/types";

function AlertCard({
  batch,
  onAction,
}: {
  batch: InventoryBatchWithRisk;
  onAction: (batchId: string, action: string, discount?: number) => void;
}) {
  const t = useTranslations("alerts");
  const [discountOpen, setDiscountOpen] = useState(false);
  const [discountVal, setDiscountVal] = useState(
    batch.suggested_discount_percent.toString()
  );

  const suggest =
    batch.risk_level === "CRITICAL"
      ? t("suggest_discount", { pct: batch.suggested_discount_percent })
      : batch.risk_level === "HIGH"
      ? t("suggest_discount", { pct: batch.suggested_discount_percent })
      : t("suggest_transfer");

  return (
    <>
      <Card
        className={`border-l-4 ${
          batch.risk_level === "CRITICAL"
            ? "border-l-red-500"
            : batch.risk_level === "HIGH"
            ? "border-l-orange-500"
            : "border-l-yellow-500"
        }`}
      >
        <CardContent className="p-5">
          <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
            {/* Product info */}
            <div className="flex items-start gap-3 flex-1 min-w-0">
              <div
                className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-lg ${
                  batch.risk_level === "CRITICAL" ? "bg-red-50" : "bg-orange-50"
                }`}
              >
                <AlertTriangle
                  className={`h-5 w-5 ${
                    batch.risk_level === "CRITICAL" ? "text-red-600" : "text-orange-500"
                  }`}
                />
              </div>
              <div className="min-w-0">
                <div className="flex items-center gap-2 flex-wrap">
                  <h3 className="font-semibold text-slate-900 truncate">
                    {batch.product_name}
                  </h3>
                  <RiskBadge level={batch.risk_level} />
                </div>
                <p className="text-xs text-slate-500 mt-0.5">
                  {batch.category} · Batch: {batch.batch_number} · {batch.supplier}
                </p>
              </div>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 shrink-0">
              <div className="text-center">
                <div className="flex items-center justify-center gap-1 text-slate-500 mb-1">
                  <Clock className="h-3.5 w-3.5" />
                  <span className="text-xs">{t("days_until_expiry")}</span>
                </div>
                <p
                  className={`text-lg font-bold ${
                    batch.days_until_expiry <= 14
                      ? "text-red-600"
                      : batch.days_until_expiry <= 30
                      ? "text-orange-500"
                      : "text-slate-700"
                  }`}
                >
                  {batch.days_until_expiry <= 0 ? "Expired" : batch.days_until_expiry}
                </p>
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center gap-1 text-slate-500 mb-1">
                  <Package className="h-3.5 w-3.5" />
                  <span className="text-xs">{t("current_qty")}</span>
                </div>
                <p className="text-lg font-bold text-slate-700">{batch.current_quantity}</p>
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center gap-1 text-slate-500 mb-1">
                  <ShoppingCart className="h-3.5 w-3.5" />
                  <span className="text-xs">{t("expected_sales")}</span>
                </div>
                <p className="text-lg font-bold text-slate-700">{batch.expected_sales}</p>
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center gap-1 text-slate-500 mb-1">
                  <TrendingDown className="h-3.5 w-3.5" />
                  <span className="text-xs">{t("estimated_loss")}</span>
                </div>
                <p className="text-lg font-bold text-red-600">
                  {formatCurrency(batch.estimated_loss)}
                </p>
              </div>
            </div>
          </div>

          {/* Suggested action */}
          <div className="mt-4 rounded-lg bg-slate-50 border border-slate-200 px-4 py-3">
            <p className="text-xs font-medium text-slate-500 mb-1">{t("suggested_action")}</p>
            <p className="text-sm text-slate-700">{suggest}</p>
          </div>

          {/* Action buttons */}
          <div className="mt-4 flex flex-wrap gap-2">
            <Button
              size="sm"
              onClick={() => setDiscountOpen(true)}
              className="gap-1.5 bg-emerald-600 hover:bg-emerald-700"
            >
              <Percent className="h-3.5 w-3.5" />
              {t("apply_discount")}
            </Button>
            <Button
              size="sm"
              variant="outline"
              onClick={() => onAction(batch.id, "TRANSFER")}
              className="gap-1.5"
            >
              <Truck className="h-3.5 w-3.5" />
              {t("transfer")}
            </Button>
            <Button
              size="sm"
              variant="outline"
              onClick={() => onAction(batch.id, "RETURN")}
              className="gap-1.5"
            >
              <RotateCcw className="h-3.5 w-3.5" />
              {t("return_supplier")}
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => onAction(batch.id, "DISMISS")}
              className="gap-1.5 text-slate-500"
            >
              <X className="h-3.5 w-3.5" />
              {t("dismiss")}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Discount dialog */}
      <Dialog open={discountOpen} onOpenChange={setDiscountOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("apply_discount")}</DialogTitle>
            <DialogDescription>
              {batch.product_name} · {batch.current_quantity} units · Expires{" "}
              {formatDate(batch.expiry_date)}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-2">
            <div className="space-y-2">
              <Label>Discount Percentage</Label>
              <div className="flex items-center gap-2">
                <Input
                  type="number"
                  min={1}
                  max={100}
                  value={discountVal}
                  onChange={(e) => setDiscountVal(e.target.value)}
                  className="w-24"
                />
                <span className="text-slate-500">%</span>
              </div>
              <p className="text-xs text-slate-500">
                Suggested: {batch.suggested_discount_percent}% · New price:{" "}
                {formatCurrency(
                  batch.selling_price * (1 - Number(discountVal) / 100)
                )}
              </p>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDiscountOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={() => {
                onAction(batch.id, "DISCOUNT", Number(discountVal));
                setDiscountOpen(false);
              }}
            >
              Apply Discount
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

export default function AlertsPage() {
  const t = useTranslations("alerts");
  const queryClient = useQueryClient();

  const { data: allAlerts, isLoading: allLoading } = useQuery({
    queryKey: ["alerts", "all"],
    queryFn: () => alertsApi.list(),
  });

  const { data: criticalAlerts } = useQuery({
    queryKey: ["alerts", "CRITICAL"],
    queryFn: () => alertsApi.list("CRITICAL"),
  });

  const { data: highAlerts } = useQuery({
    queryKey: ["alerts", "HIGH"],
    queryFn: () => alertsApi.list("HIGH"),
  });

  const { data: history } = useQuery({
    queryKey: ["alert-history"],
    queryFn: () => alertsApi.history(20),
  });

  const actionMutation = useMutation({
    mutationFn: ({
      batchId,
      action,
      discount,
    }: {
      batchId: string;
      action: string;
      discount?: number;
    }) =>
      alertsApi.takeAction(batchId, {
        action_type: action,
        discount_percent: discount,
      }),
    onSuccess: () => {
      toast.success(t("action_recorded"));
      queryClient.invalidateQueries({ queryKey: ["alerts"] });
      queryClient.invalidateQueries({ queryKey: ["alert-history"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : "Action failed");
    },
  });

  const handleAction = (batchId: string, action: string, discount?: number) => {
    actionMutation.mutate({ batchId, action, discount });
  };

  const critCount = criticalAlerts?.total ?? 0;
  const highCount = highAlerts?.total ?? 0;
  const allCount = allAlerts?.total ?? 0;

  const renderAlerts = (items: InventoryBatchWithRisk[] | undefined, loading: boolean) => {
    if (loading) {
      return (
        <div className="space-y-4">
          {Array.from({ length: 3 }).map((_, i) => (
            <Card key={i}>
              <CardContent className="p-5">
                <Skeleton className="h-24 w-full" />
              </CardContent>
            </Card>
          ))}
        </div>
      );
    }
    if (!items?.length) {
      return (
        <div className="py-16 text-center">
          <AlertTriangle className="h-8 w-8 mx-auto mb-3 text-slate-300" />
          <p className="text-slate-500">{t("no_alerts")}</p>
        </div>
      );
    }
    return (
      <div className="space-y-3">
        {items.map((batch) => (
          <AlertCard key={batch.id} batch={batch} onAction={handleAction} />
        ))}
      </div>
    );
  };

  return (
    <div className="space-y-6 max-w-5xl">
      <Tabs defaultValue="critical">
        <TabsList>
          <TabsTrigger value="critical">
            {t("critical_tab")} ({critCount})
          </TabsTrigger>
          <TabsTrigger value="high">
            {t("high_tab")} ({highCount})
          </TabsTrigger>
          <TabsTrigger value="all">
            {t("all_tab")} ({allCount})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="critical" className="mt-4">
          {renderAlerts(criticalAlerts?.data, allLoading)}
        </TabsContent>
        <TabsContent value="high" className="mt-4">
          {renderAlerts(highAlerts?.data, allLoading)}
        </TabsContent>
        <TabsContent value="all" className="mt-4">
          {renderAlerts(allAlerts?.data, allLoading)}
        </TabsContent>
      </Tabs>

      {/* Action History */}
      {history && history.length > 0 && (
        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center gap-2">
              <History className="h-4 w-4 text-slate-500" />
              <CardTitle className="text-base">{t("action_history")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="divide-y divide-slate-100">
              {history.map((action) => (
                <div key={action.id} className="flex items-center justify-between py-3 gap-4">
                  <div className="min-w-0">
                    <p className="text-sm font-medium text-slate-900 truncate">
                      {action.product_name}
                    </p>
                    <p className="text-xs text-slate-500">
                      Batch {action.batch_number} · by {action.user_name}
                    </p>
                  </div>
                  <div className="text-right shrink-0">
                    <span
                      className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                        action.action_type === "DISCOUNT"
                          ? "bg-emerald-100 text-emerald-700"
                          : action.action_type === "TRANSFER"
                          ? "bg-blue-100 text-blue-700"
                          : action.action_type === "RETURN"
                          ? "bg-purple-100 text-purple-700"
                          : "bg-slate-100 text-slate-600"
                      }`}
                    >
                      {action.action_type}
                      {action.discount_percent ? ` ${action.discount_percent}%` : ""}
                    </span>
                    <p className="text-xs text-slate-400 mt-1">
                      {new Date(action.created_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

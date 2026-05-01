"use client";

import { useEffect, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Building2, Users, Shield } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { pharmacyApi } from "@/lib/api";
import { getStoredUser } from "@/lib/auth";
import { Pharmacy } from "@/types";

export default function SettingsPage() {
  const t = useTranslations("settings");
  const queryClient = useQueryClient();
  const currentUser = getStoredUser();

  const { data: pharmacy, isLoading } = useQuery({
    queryKey: ["pharmacy"],
    queryFn: pharmacyApi.get,
  });

  const { data: users } = useQuery({
    queryKey: ["users"],
    queryFn: pharmacyApi.getUsers,
    enabled: currentUser?.role === "admin",
  });

  const [form, setForm] = useState<Partial<Pharmacy>>({});

  useEffect(() => {
    if (pharmacy) {
      setForm({
        name: pharmacy.name,
        address: pharmacy.address,
        city: pharmacy.city,
        phone: pharmacy.phone,
        email: pharmacy.email,
        language: pharmacy.language,
      });
    }
  }, [pharmacy]);

  const updateMutation = useMutation({
    mutationFn: (data: Partial<Pharmacy>) => pharmacyApi.update(data),
    onSuccess: () => {
      toast.success(t("saved"));
      queryClient.invalidateQueries({ queryKey: ["pharmacy"] });
    },
    onError: () => {
      toast.error("Update failed");
    },
  });

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    if (currentUser?.role !== "admin") {
      toast.error("Admin access required");
      return;
    }
    updateMutation.mutate(form);
  };

  return (
    <div className="space-y-6 max-w-2xl">
      {/* Pharmacy Info */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <Building2 className="h-5 w-5 text-emerald-600" />
            <CardTitle>{t("pharmacy_info")}</CardTitle>
          </div>
          <CardDescription>
            License: {pharmacy?.license_number || "—"}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 4 }).map((_, i) => (
                <div key={i} className="space-y-2">
                  <Skeleton className="h-4 w-24" />
                  <Skeleton className="h-9 w-full" />
                </div>
              ))}
            </div>
          ) : (
            <form onSubmit={handleSave} className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="sm:col-span-2 space-y-2">
                  <Label>{t("name")}</Label>
                  <Input
                    value={form.name || ""}
                    onChange={(e) => setForm({ ...form, name: e.target.value })}
                    disabled={currentUser?.role !== "admin"}
                  />
                </div>
                <div className="sm:col-span-2 space-y-2">
                  <Label>{t("address")}</Label>
                  <Input
                    value={form.address || ""}
                    onChange={(e) => setForm({ ...form, address: e.target.value })}
                    disabled={currentUser?.role !== "admin"}
                  />
                </div>
                <div className="space-y-2">
                  <Label>{t("city")}</Label>
                  <Input
                    value={form.city || ""}
                    onChange={(e) => setForm({ ...form, city: e.target.value })}
                    disabled={currentUser?.role !== "admin"}
                  />
                </div>
                <div className="space-y-2">
                  <Label>{t("phone")}</Label>
                  <Input
                    value={form.phone || ""}
                    onChange={(e) => setForm({ ...form, phone: e.target.value })}
                    disabled={currentUser?.role !== "admin"}
                  />
                </div>
                <div className="space-y-2">
                  <Label>{t("email")}</Label>
                  <Input
                    type="email"
                    value={form.email || ""}
                    onChange={(e) => setForm({ ...form, email: e.target.value })}
                    disabled={currentUser?.role !== "admin"}
                  />
                </div>
                <div className="space-y-2">
                  <Label>{t("language")}</Label>
                  <Select
                    value={form.language || "en"}
                    onValueChange={(v) => setForm({ ...form, language: v })}
                    disabled={currentUser?.role !== "admin"}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="en">English</SelectItem>
                      <SelectItem value="el">Ελληνικά</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              {currentUser?.role === "admin" && (
                <div className="flex justify-end pt-2">
                  <Button type="submit" disabled={updateMutation.isPending}>
                    {updateMutation.isPending ? "Saving..." : t("save")}
                  </Button>
                </div>
              )}
            </form>
          )}
        </CardContent>
      </Card>

      {/* User Management */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <Users className="h-5 w-5 text-emerald-600" />
            <CardTitle>{t("users")}</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          {currentUser?.role !== "admin" ? (
            <div className="flex items-center gap-2 text-slate-500 text-sm">
              <Shield className="h-4 w-4" />
              <p>Admin access required to view users</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200">
                    <th className="pb-2 text-left font-medium text-slate-600">{t("user_name")}</th>
                    <th className="pb-2 text-left font-medium text-slate-600">{t("user_email")}</th>
                    <th className="pb-2 text-left font-medium text-slate-600">{t("user_role")}</th>
                    <th className="pb-2 text-left font-medium text-slate-600">Status</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {(users || []).map((user) => (
                    <tr key={user.id}>
                      <td className="py-3 font-medium text-slate-900">{user.full_name}</td>
                      <td className="py-3 text-slate-600">{user.email}</td>
                      <td className="py-3">
                        <span
                          className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                            user.role === "admin"
                              ? "bg-emerald-100 text-emerald-700"
                              : "bg-slate-100 text-slate-600"
                          }`}
                        >
                          {user.role === "admin" ? t("admin") : t("staff")}
                        </span>
                      </td>
                      <td className="py-3">
                        <span
                          className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs ${
                            user.is_active ? "bg-green-100 text-green-700" : "bg-slate-100 text-slate-500"
                          }`}
                        >
                          {user.is_active ? "Active" : "Inactive"}
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
  );
}

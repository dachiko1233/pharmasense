"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useLocale, useTranslations } from "next-intl";
import { Activity, Eye, EyeOff } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { authApi } from "@/lib/api";
import { saveAuth } from "@/lib/auth";

export default function LoginPage() {
  const t = useTranslations("auth");
  const locale = useLocale();
  const router = useRouter();
  const [email, setEmail] = useState("admin@pharmasense.cy");
  const [password, setPassword] = useState("Demo1234!");
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await authApi.login(email, password);
      saveAuth(res.token, res.user);
      router.push(`/${locale}/dashboard`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* Left panel */}
      <div className="hidden lg:flex lg:w-1/2 flex-col justify-between bg-emerald-600 p-12 text-white">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-white/20">
            <Activity className="h-5 w-5 text-white" />
          </div>
          <div>
            <span className="text-xl font-bold">Pharma</span>
            <span className="text-xl font-bold text-emerald-200">Sense</span>
          </div>
        </div>

        <div>
          <h2 className="text-4xl font-bold leading-tight mb-4">
            {t("tagline")}
          </h2>
          <p className="text-emerald-100 text-lg">
            Catch expiring products before they expire. Save thousands of euros every month.
          </p>

          <div className="mt-12 grid grid-cols-2 gap-4">
            {[
              { label: "Avg Monthly Savings", value: "€2,400" },
              { label: "Products Monitored", value: "500+" },
              { label: "Pharmacies in Cyprus", value: "12+" },
              { label: "Risk Calculations", value: "Daily" },
            ].map((stat) => (
              <div key={stat.label} className="bg-white/10 rounded-xl p-4">
                <div className="text-2xl font-bold">{stat.value}</div>
                <div className="text-sm text-emerald-200 mt-1">{stat.label}</div>
              </div>
            ))}
          </div>
        </div>

        <p className="text-emerald-200 text-sm">
          © 2024 PharmaSense. Built for Cyprus pharmacies.
        </p>
      </div>

      {/* Right panel */}
      <div className="flex flex-1 flex-col items-center justify-center p-8 bg-slate-50">
        <div className="w-full max-w-sm">
          {/* Mobile logo */}
          <div className="lg:hidden flex items-center gap-2 mb-8">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-600">
              <Activity className="h-4 w-4 text-white" />
            </div>
            <span className="text-lg font-bold text-slate-900">
              Pharma<span className="text-emerald-600">Sense</span>
            </span>
          </div>

          <div className="mb-8">
            <h1 className="text-2xl font-bold text-slate-900">{t("login")}</h1>
            <p className="text-slate-500 mt-1 text-sm">
              Sign in to your pharmacy dashboard
            </p>
          </div>

          <form onSubmit={handleLogin} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">{t("email")}</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="admin@pharmasense.cy"
                required
                autoComplete="email"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">{t("password")}</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  required
                  autoComplete="current-password"
                  className="pr-10"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                >
                  {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
            </div>

            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? t("signing_in") : t("login")}
            </Button>
          </form>

          {/* Demo hint */}
          <div className="mt-6 rounded-lg border border-emerald-200 bg-emerald-50 p-4">
            <p className="text-xs font-medium text-emerald-700 mb-1">Demo Credentials</p>
            <p className="text-xs text-emerald-600 font-mono">admin@pharmasense.cy</p>
            <p className="text-xs text-emerald-600 font-mono">Demo1234!</p>
          </div>
        </div>
      </div>
    </div>
  );
}

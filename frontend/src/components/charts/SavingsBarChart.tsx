"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { formatCurrency } from "@/lib/utils";

interface SavingsDataPoint {
  name: string;
  saved: number;
  actions: number;
}

export function SavingsBarChart({ data }: { data: SavingsDataPoint[] }) {
  return (
    <ResponsiveContainer width="100%" height={260}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
        <XAxis dataKey="name" tick={{ fontSize: 11, fill: "#94a3b8" }} tickLine={false} />
        <YAxis
          tick={{ fontSize: 11, fill: "#94a3b8" }}
          tickLine={false}
          axisLine={false}
          tickFormatter={(v) => `€${v}`}
        />
        <Tooltip
          formatter={(v: number) => formatCurrency(v)}
          contentStyle={{ border: "1px solid #e2e8f0", borderRadius: "8px", fontSize: "12px" }}
        />
        <Bar dataKey="saved" fill="#059669" name="Estimated Saved" radius={[4, 4, 0, 0]} />
      </BarChart>
    </ResponsiveContainer>
  );
}

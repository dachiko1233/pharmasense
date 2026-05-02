"use client";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

interface SavingsDataPoint {
  name: string;
  saved: number;
  actions: number;
}

export function ActionsLineChart({ data }: { data: SavingsDataPoint[] }) {
  return (
    <ResponsiveContainer width="100%" height={200}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
        <XAxis dataKey="name" tick={{ fontSize: 11, fill: "#94a3b8" }} tickLine={false} />
        <YAxis tick={{ fontSize: 11, fill: "#94a3b8" }} tickLine={false} axisLine={false} />
        <Tooltip
          contentStyle={{ border: "1px solid #e2e8f0", borderRadius: "8px", fontSize: "12px" }}
        />
        <Line dataKey="actions" stroke="#059669" strokeWidth={2} dot={false} name="Actions" />
      </LineChart>
    </ResponsiveContainer>
  );
}

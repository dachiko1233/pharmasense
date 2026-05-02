"use client";

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

interface TimelineDataPoint {
  name: string;
  critical: number;
  high: number;
  medium: number;
  low: number;
}

export function ExpiryTimelineChart({ data }: { data: TimelineDataPoint[] }) {
  return (
    <ResponsiveContainer width="100%" height={260}>
      <BarChart data={data} barSize={14}>
        <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
        <XAxis dataKey="name" tick={{ fontSize: 11, fill: "#94a3b8" }} tickLine={false} />
        <YAxis tick={{ fontSize: 11, fill: "#94a3b8" }} tickLine={false} axisLine={false} />
        <Tooltip
          contentStyle={{ border: "1px solid #e2e8f0", borderRadius: "8px", fontSize: "12px" }}
        />
        <Legend wrapperStyle={{ fontSize: "12px" }} />
        <Bar dataKey="critical" stackId="a" fill="#dc2626" name="Critical" radius={[0, 0, 0, 0]} />
        <Bar dataKey="high" stackId="a" fill="#f97316" name="High" />
        <Bar dataKey="medium" stackId="a" fill="#eab308" name="Medium" />
        <Bar dataKey="low" stackId="a" fill="#22c55e" name="Low" radius={[4, 4, 0, 0]} />
      </BarChart>
    </ResponsiveContainer>
  );
}

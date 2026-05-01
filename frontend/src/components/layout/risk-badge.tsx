import { useTranslations } from "next-intl";
import { Badge } from "@/components/ui/badge";
import { RiskLevel } from "@/types";

const variantMap: Record<RiskLevel, "critical" | "high" | "medium" | "low"> = {
  CRITICAL: "critical",
  HIGH: "high",
  MEDIUM: "medium",
  LOW: "low",
};

interface RiskBadgeProps {
  level: RiskLevel;
  className?: string;
}

export function RiskBadge({ level, className }: RiskBadgeProps) {
  const t = useTranslations("risk");
  return (
    <Badge variant={variantMap[level]} className={className}>
      {t(level)}
    </Badge>
  );
}

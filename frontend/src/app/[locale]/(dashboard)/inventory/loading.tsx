import { Skeleton } from "@/components/ui/skeleton";
import { Card } from "@/components/ui/card";

export default function InventoryLoading() {
  return (
    <div className="space-y-4 max-w-7xl">
      <div className="flex flex-col sm:flex-row gap-3">
        <Skeleton className="h-10 flex-1" />
        <Skeleton className="h-10 w-36" />
        <Skeleton className="h-10 w-40" />
        <Skeleton className="h-10 w-44" />
      </div>

      <div className="flex items-center justify-between">
        <Skeleton className="h-4 w-24" />
        <Skeleton className="h-7 w-44" />
      </div>

      <Card>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-200 bg-slate-50">
                {Array.from({ length: 7 }).map((_, i) => (
                  <th key={i} className="p-4">
                    <Skeleton className="h-4 w-full" />
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {Array.from({ length: 8 }).map((_, i) => (
                <tr key={i}>
                  {Array.from({ length: 7 }).map((_, j) => (
                    <td key={j} className="p-4">
                      <Skeleton className="h-4 w-full" />
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
}

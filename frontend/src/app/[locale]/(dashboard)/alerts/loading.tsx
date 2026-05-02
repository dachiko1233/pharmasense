import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent } from "@/components/ui/card";

export default function AlertsLoading() {
  return (
    <div className="space-y-6 max-w-5xl">
      <div className="flex gap-2">
        <Skeleton className="h-9 w-32" />
        <Skeleton className="h-9 w-28" />
        <Skeleton className="h-9 w-24" />
      </div>

      <div className="space-y-3 mt-4">
        {Array.from({ length: 3 }).map((_, i) => (
          <Card key={i} className="border-l-4 border-l-slate-200">
            <CardContent className="p-5">
              <div className="flex flex-col sm:flex-row gap-4">
                <div className="flex items-start gap-3 flex-1">
                  <Skeleton className="h-10 w-10 rounded-lg shrink-0" />
                  <div className="flex-1 space-y-2">
                    <Skeleton className="h-5 w-48" />
                    <Skeleton className="h-4 w-64" />
                  </div>
                </div>
                <div className="grid grid-cols-4 gap-3">
                  {Array.from({ length: 4 }).map((_, j) => (
                    <Skeleton key={j} className="h-12 w-16" />
                  ))}
                </div>
              </div>
              <Skeleton className="h-12 w-full mt-4 rounded-lg" />
              <div className="flex gap-2 mt-4">
                <Skeleton className="h-8 w-32" />
                <Skeleton className="h-8 w-24" />
                <Skeleton className="h-8 w-28" />
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}

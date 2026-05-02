import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardHeader } from "@/components/ui/card";

export default function ImportLoading() {
  return (
    <div className="space-y-6 max-w-3xl">
      <Card>
        <CardHeader>
          <Skeleton className="h-5 w-32 mb-1" />
          <Skeleton className="h-4 w-64" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-40 w-full rounded-lg" />
          <div className="flex gap-3 mt-4">
            <Skeleton className="h-10 flex-1" />
            <Skeleton className="h-10 w-36" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <Skeleton className="h-5 w-40" />
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="h-10 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

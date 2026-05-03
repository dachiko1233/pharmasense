'use client';

import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { toast } from 'sonner';
import { Search, SlidersHorizontal, Download, X } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Skeleton } from '@/components/ui/skeleton';
import { RiskBadge } from '@/components/layout/risk-badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { inventoryApi } from '@/lib/api';
import { formatCurrency, formatDate, getRiskRowClass } from '@/lib/utils';
import { InventoryBatchWithRisk } from '@/types';

const PAGE_SIZE = 50;

export default function InventoryPage() {
  const t = useTranslations('inventory');
  const tc = useTranslations('common');
  const [search, setSearch] = useState('');
  const [riskLevel, setRiskLevel] = useState('');
  const [category, setCategory] = useState('');
  const [supplier, setSupplier] = useState('');
  const [page, setPage] = useState(1);
  const [sortBy, setSortBy] = useState('ra.risk_level');
  const [sortDir, setSortDir] = useState('ASC');

  const params = {
    search,
    risk_level: riskLevel,
    category,
    supplier,
    page,
    limit: PAGE_SIZE,
    sort_by: sortBy,
    sort_dir: sortDir,
  };

  const { data, isLoading } = useQuery({
    queryKey: ['inventory', params],
    queryFn: () => inventoryApi.list(params),
    placeholderData: (prev) => prev,
  });

  const { data: filters } = useQuery({
    queryKey: ['inventory-filters'],
    queryFn: inventoryApi.filters,
  });

  const clearFilters = () => {
    setSearch('');
    setRiskLevel('');
    setCategory('');
    setSupplier('');
    setPage(1);
  };

  const hasFilters = search || riskLevel || category || supplier;
  const total = data?.total ?? 0;
  const batches: InventoryBatchWithRisk[] = data?.data ?? [];
  const totalPages = Math.ceil(total / PAGE_SIZE);

  const handleExport = async () => {
    try {
      const blob = await inventoryApi.exportCSV();
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'inventory.csv';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch {
      toast.error('Export failed');
    }
  };

  return (
    <div className="space-y-4 max-w-7xl">
      {/* Top bar */}
      <div className="flex flex-col sm:flex-row gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
          <Input
            placeholder={t('search')}
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              setPage(1);
            }}
            className="pl-9"
          />
        </div>

        <div className="flex gap-2 flex-wrap">
          <Select
            value={riskLevel || 'all'}
            onValueChange={(v) => {
              setRiskLevel(v === 'all' ? '' : v);
              setPage(1);
            }}
          >
            <SelectTrigger className="w-36">
              <SlidersHorizontal className="h-3.5 w-3.5 mr-1.5 text-slate-400" />
              <SelectValue placeholder={t('risk_level')} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t('all')}</SelectItem>
              {(
                filters?.risk_levels ?? ['CRITICAL', 'HIGH', 'MEDIUM', 'LOW']
              ).map((r) => (
                <SelectItem key={r} value={r}>
                  {r}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            value={category || 'all'}
            onValueChange={(v) => {
              setCategory(v === 'all' ? '' : v);
              setPage(1);
            }}
          >
            <SelectTrigger className="w-40">
              <SelectValue placeholder={t('category')} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t('all')}</SelectItem>
              {(filters?.categories ?? []).map((c) => (
                <SelectItem key={c} value={c}>
                  {c}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            value={supplier || 'all'}
            onValueChange={(v) => {
              setSupplier(v === 'all' ? '' : v);
              setPage(1);
            }}
          >
            <SelectTrigger className="w-44">
              <SelectValue placeholder={t('supplier')} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t('all')}</SelectItem>
              {(filters?.suppliers ?? []).map((s) => (
                <SelectItem key={s} value={s}>
                  {s}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          {hasFilters && (
            <Button
              variant="ghost"
              size="sm"
              onClick={clearFilters}
              className="gap-1"
            >
              <X className="h-3.5 w-3.5" />
              {t('clear_filters')}
            </Button>
          )}

          <Button variant="outline" size="sm" onClick={handleExport}>
            <Download className="h-4 w-4 mr-1.5" />
            {t('export')}
          </Button>
        </div>
      </div>

      {/* Stats */}
      <div className="flex items-center justify-between text-sm text-slate-500">
        <span>{isLoading ? 'Loading...' : `${total} batches`}</span>
        <div className="flex items-center gap-2">
          <span>{tc('of')}</span>
          <Select value={sortBy} onValueChange={setSortBy}>
            <SelectTrigger className="h-7 w-44 text-xs">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ra.risk_level">Risk Level</SelectItem>
              <SelectItem value="ib.expiry_date">Expiry Date</SelectItem>
              <SelectItem value="ra.estimated_loss">Estimated Loss</SelectItem>
              <SelectItem value="p.name">Product Name</SelectItem>
            </SelectContent>
          </Select>
          <Button
            variant="ghost"
            size="sm"
            className="h-7 px-2 text-xs"
            onClick={() => setSortDir((d) => (d === 'ASC' ? 'DESC' : 'ASC'))}
          >
            {sortDir === 'ASC' ? '↑' : '↓'}
          </Button>
        </div>
      </div>

      {/* Table */}
      <Card>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-200 bg-slate-50">
                <th className="text-left p-4 font-medium text-slate-600">
                  {t('product')}
                </th>
                <th className="text-left p-4 font-medium text-slate-600 hidden md:table-cell">
                  {t('batch')}
                </th>
                <th className="text-left p-4 font-medium text-slate-600">
                  {t('expiry')}
                </th>
                <th className="text-right p-4 font-medium text-slate-600">
                  {t('quantity')}
                </th>
                <th className="text-center p-4 font-medium text-slate-600">
                  {t('risk')}
                </th>
                <th className="text-right p-4 font-medium text-slate-600">
                  {t('loss')}
                </th>
                <th className="text-left p-4 font-medium text-slate-600 hidden lg:table-cell">
                  Supplier
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {isLoading
                ? Array.from({ length: 8 }).map((_, i) => (
                    <tr key={i}>
                      {Array.from({ length: 6 }).map((_, j) => (
                        <td key={j} className="p-4">
                          <Skeleton className="h-4 w-full" />
                        </td>
                      ))}
                    </tr>
                  ))
                : batches.map((batch) => (
                    <tr
                      key={batch.id}
                      className={`hover:bg-slate-50 transition-colors ${getRiskRowClass(batch.risk_level)}`}
                    >
                      <td className="p-4">
                        <div>
                          <p className="font-medium text-slate-900">
                            {batch.product_name}
                          </p>
                          <p className="text-xs text-slate-400">
                            {batch.category}
                          </p>
                        </div>
                      </td>
                      <td className="p-4 hidden md:table-cell">
                        <span className="font-mono text-xs text-slate-600">
                          {batch.batch_number}
                        </span>
                      </td>
                      <td className="p-4">
                        <div>
                          <p className="font-medium text-slate-700">
                            {formatDate(batch.expiry_date)}
                          </p>
                          <p
                            className={`text-xs ${
                              batch.days_until_expiry <= 0
                                ? 'text-red-600 font-semibold'
                                : batch.days_until_expiry <= 30
                                  ? 'text-red-500'
                                  : batch.days_until_expiry <= 90
                                    ? 'text-orange-500'
                                    : 'text-slate-400'
                            }`}
                          >
                            {batch.days_until_expiry <= 0
                              ? 'Expired'
                              : `${batch.days_until_expiry} ${t('days')}`}
                          </p>
                        </div>
                      </td>
                      <td className="p-4 text-right">
                        <span className="font-medium text-slate-700">
                          {batch.current_quantity}
                        </span>
                        <p className="text-xs text-slate-400">units</p>
                      </td>
                      <td className="p-4 text-center">
                        <RiskBadge level={batch.risk_level} />
                      </td>
                      <td className="p-4 text-right">
                        <span
                          className={`font-semibold ${
                            batch.estimated_loss > 0
                              ? 'text-red-600'
                              : 'text-slate-400'
                          }`}
                        >
                          {batch.estimated_loss > 0
                            ? formatCurrency(batch.estimated_loss)
                            : '—'}
                        </span>
                      </td>
                      <td className="p-4 hidden lg:table-cell text-slate-500 text-xs">
                        {batch.supplier}
                      </td>
                    </tr>
                  ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between border-t border-slate-200 px-4 py-3">
            <p className="text-sm text-slate-500">
              Showing {(page - 1) * PAGE_SIZE + 1}–
              {Math.min(page * PAGE_SIZE, total)} of {total}
            </p>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage((p) => p - 1)}
                disabled={page === 1}
              >
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage((p) => p + 1)}
                disabled={page >= totalPages}
              >
                Next
              </Button>
            </div>
          </div>
        )}

        {!isLoading && batches.length === 0 && (
          <div className="py-16 text-center">
            <Search className="h-8 w-8 mx-auto mb-3 text-slate-300" />
            <p className="text-slate-500 font-medium">{t('no_results')}</p>
            {hasFilters && (
              <Button
                variant="ghost"
                size="sm"
                onClick={clearFilters}
                className="mt-2"
              >
                {t('clear_filters')}
              </Button>
            )}
          </div>
        )}
      </Card>
    </div>
  );
}

"use client";

import { useState, useRef } from "react";
import { useTranslations } from "next-intl";
import { Upload, FileText, Download, AlertCircle, CheckCircle2 } from "lucide-react";
import { toast } from "sonner";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { inventoryApi } from "@/lib/api";

interface PreviewRow {
  row: number;
  product_name: string;
  barcode: string;
  batch_number: string;
  expiry_date: string;
  quantity: string;
  purchase_price: string;
  error?: string;
}

interface ImportResult {
  total_rows: number;
  preview: PreviewRow[];
  error_count: number;
  template_url: string;
}

interface ConfirmResult {
  imported: number;
  errors: { row: number; error: string }[];
}

export default function ImportPage() {
  const t = useTranslations("import");
  const fileRef = useRef<HTMLInputElement>(null);
  const [dragOver, setDragOver] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<ImportResult | null>(null);
  const [confirmResult, setConfirmResult] = useState<ConfirmResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [importing, setImporting] = useState(false);

  const processFile = async (f: File) => {
    if (!f.name.endsWith(".csv")) {
      toast.error("Please upload a CSV file");
      return;
    }
    setFile(f);
    setResult(null);
    setConfirmResult(null);
    setLoading(true);
    try {
      const data = await inventoryApi.importCSV(f);
      setResult(data);
      if (data.error_count === 0) {
        toast.success(`${t("preview")}: ${data.total_rows} rows ready`);
      } else {
        toast.warning(`${data.error_count} validation errors found`);
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setLoading(false);
    }
  };

  const handleConfirmImport = async () => {
    if (!file) return;
    setImporting(true);
    try {
      const data = await inventoryApi.confirmImport(file);
      setConfirmResult(data);
      setResult(null);
      if (data.errors.length === 0) {
        toast.success(`${t("success")}: ${data.imported} rows imported`);
      } else {
        toast.warning(`${data.imported} imported, ${data.errors.length} errors`);
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Import failed");
    } finally {
      setImporting(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    const f = e.dataTransfer.files[0];
    if (f) processFile(f);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0];
    if (f) processFile(f);
  };

  return (
    <div className="space-y-6 max-w-3xl">
      <div className="flex items-center justify-between">
        <p className="text-slate-500 text-sm">{t("subtitle")}</p>
        <Button variant="outline" size="sm" asChild className="gap-2">
          <a href={inventoryApi.templateURL()} download>
            <Download className="h-4 w-4" />
            {t("download_template")}
          </a>
        </Button>
      </div>

      {/* Drop zone */}
      <Card>
        <CardContent className="p-6">
          <div
            className={`rounded-xl border-2 border-dashed p-12 text-center transition-colors ${
              dragOver
                ? "border-emerald-500 bg-emerald-50"
                : "border-slate-200 hover:border-slate-300 hover:bg-slate-50"
            }`}
            onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
            onDragLeave={() => setDragOver(false)}
            onDrop={handleDrop}
          >
            <div className="flex justify-center mb-4">
              <div className={`flex h-14 w-14 items-center justify-center rounded-xl ${
                dragOver ? "bg-emerald-100" : "bg-slate-100"
              }`}>
                <Upload className={`h-6 w-6 ${dragOver ? "text-emerald-600" : "text-slate-400"}`} />
              </div>
            </div>
            <p className="text-slate-700 font-medium mb-1">{t("drag_drop")}</p>
            <p className="text-slate-400 text-sm mb-4">{t("or")}</p>
            <Button
              variant="outline"
              onClick={() => fileRef.current?.click()}
              disabled={loading}
            >
              <FileText className="h-4 w-4 mr-2" />
              {t("browse")}
            </Button>
            <input
              ref={fileRef}
              type="file"
              accept=".csv"
              className="hidden"
              onChange={handleFileChange}
            />
            {file && (
              <p className="mt-3 text-sm text-emerald-600 font-medium">
                ✓ {file.name}
              </p>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Required columns */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Required Columns</CardTitle>
          <CardDescription>Your CSV must include these columns (header row)</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            {["product_name", "barcode", "batch_number", "expiry_date", "quantity", "purchase_price"].map((col) => (
              <span key={col} className="rounded-md bg-slate-100 px-2.5 py-1 font-mono text-xs text-slate-700">
                {col}
              </span>
            ))}
          </div>
          <p className="text-xs text-slate-500 mt-3">
            Optional: <span className="font-mono">selling_price</span>,{" "}
            <span className="font-mono">supplier</span>. Date format:{" "}
            <span className="font-mono">YYYY-MM-DD</span>
          </p>
        </CardContent>
      </Card>

      {/* Import result */}
      {confirmResult && (
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center gap-3">
              {confirmResult.errors.length === 0 ? (
                <CheckCircle2 className="h-8 w-8 text-emerald-500 shrink-0" />
              ) : (
                <AlertCircle className="h-8 w-8 text-orange-500 shrink-0" />
              )}
              <div>
                <p className="font-semibold text-slate-900">
                  {confirmResult.imported} {t("rows")} imported successfully
                </p>
                {confirmResult.errors.length > 0 && (
                  <p className="text-sm text-orange-600 mt-0.5">
                    {confirmResult.errors.length} rows skipped due to errors
                  </p>
                )}
              </div>
            </div>
            {confirmResult.errors.length > 0 && (
              <ul className="mt-3 space-y-1">
                {confirmResult.errors.map((e, i) => (
                  <li key={i} className="text-xs text-red-600">
                    Row {e.row}: {e.error}
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
      )}

      {/* Preview */}
      {result && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-base">{t("preview")}</CardTitle>
              <div className="flex items-center gap-2">
                {result.error_count === 0 ? (
                  <span className="flex items-center gap-1 text-emerald-600 text-sm">
                    <CheckCircle2 className="h-4 w-4" />
                    {result.total_rows} rows valid
                  </span>
                ) : (
                  <span className="flex items-center gap-1 text-red-600 text-sm">
                    <AlertCircle className="h-4 w-4" />
                    {t("errors_found", { count: result.error_count })}
                  </span>
                )}
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <table className="w-full text-xs">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="p-2 text-left font-medium text-slate-600">{t("row")}</th>
                    <th className="p-2 text-left font-medium text-slate-600">Product</th>
                    <th className="p-2 text-left font-medium text-slate-600">Batch#</th>
                    <th className="p-2 text-left font-medium text-slate-600">Expiry</th>
                    <th className="p-2 text-right font-medium text-slate-600">Qty</th>
                    <th className="p-2 text-left font-medium text-slate-600">{t("error")}</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {result.preview.map((row) => (
                    <tr key={row.row} className={row.error ? "bg-red-50" : ""}>
                      <td className="p-2 text-slate-500">{row.row}</td>
                      <td className="p-2 font-medium text-slate-700">{row.product_name}</td>
                      <td className="p-2 font-mono text-slate-600">{row.batch_number}</td>
                      <td className="p-2 text-slate-600">{row.expiry_date}</td>
                      <td className="p-2 text-right text-slate-600">{row.quantity}</td>
                      <td className="p-2 text-red-600">{row.error || ""}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <div className="mt-4 flex items-center justify-between">
              <p className="text-xs text-slate-500">
                Showing first 20 of {result.total_rows} rows
              </p>
              <Button
                disabled={result.error_count > 0 || importing}
                onClick={handleConfirmImport}
              >
                {importing ? "Importing..." : `${t("import_btn")} (${result.total_rows} rows)`}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

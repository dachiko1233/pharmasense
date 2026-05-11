import { redirect } from "next/navigation";

export default async function LocaleIndex({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (locale === "en") {
    redirect("/dashboard");
  } else {
    redirect(`/${locale}/dashboard`);
  }
}

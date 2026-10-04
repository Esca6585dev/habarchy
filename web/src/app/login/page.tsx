"use client";
import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Field } from "@/components/common/field";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { LocaleSwitcher } from "@/components/shell/locale-switcher";
import { ThemeToggle } from "@/components/shell/theme-toggle";

const schema = z.object({
  email: z.string().email(),
  password: z.string().min(1),
  totp_code: z.string().optional(),
});
type Form = z.infer<typeof schema>;

function LoginForm() {
  const t = useTranslations("auth");
  const router = useRouter();
  const params = useSearchParams();
  const [needTotp, setNeedTotp] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const form = useForm<Form>({ resolver: zodResolver(schema), defaultValues: { email: "", password: "", totp_code: "" } });

  const onSubmit = form.handleSubmit(async (values) => {
    setError(null);
    const res = await fetch("/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ...values, totp_code: values.totp_code || undefined }),
    });
    const json = (await res.json()) as { error?: { message: string; details?: { totp_required?: boolean } } };
    if (!res.ok) {
      if (json.error?.details?.totp_required) {
        setNeedTotp(true);
        return;
      }
      setError(json.error?.message ?? t("invalid"));
      return;
    }
    const next = params.get("next");
    router.replace(next && next.startsWith("/") ? next : "/");
    router.refresh();
  });

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle className="text-xl">{t("title")}</CardTitle>
        <CardDescription>{t("subtitle")}</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="space-y-4" noValidate>
          <Field label={t("email")} htmlFor="email" error={form.formState.errors.email?.message}>
            <Input id="email" type="email" autoComplete="username" autoFocus {...form.register("email")} />
          </Field>
          <Field label={t("password")} htmlFor="password" error={form.formState.errors.password?.message}>
            <Input id="password" type="password" autoComplete="current-password" {...form.register("password")} />
          </Field>
          {needTotp ? (
            <Field label={t("totp")} htmlFor="totp" hint={t("totpHint")}>
              <Input id="totp" inputMode="numeric" pattern="[0-9]*" maxLength={6} autoFocus {...form.register("totp_code")} />
            </Field>
          ) : null}
          {error ? (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          ) : null}
          <Button type="submit" className="w-full" disabled={form.formState.isSubmitting}>
            {form.formState.isSubmitting ? t("signingIn") : t("signIn")}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}

export default function LoginPage() {
  return (
    <div className="flex min-h-dvh flex-col bg-muted/40">
      <div className="flex items-center justify-end gap-1 p-3">
        <LocaleSwitcher />
        <ThemeToggle />
      </div>
      <div className="flex flex-1 items-center justify-center px-4 pb-16">
        <Suspense>
          <LoginForm />
        </Suspense>
      </div>
    </div>
  );
}

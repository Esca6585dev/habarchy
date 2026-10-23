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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { LocaleSwitcher } from "@/components/shell/locale-switcher";
import { ThemeToggle } from "@/components/shell/theme-toggle";
import { BrandMark } from "@/components/brand-mark";

// Country dial codes for phone login (Turkmenistan first, then the region and
// a few common destinations). The backend normalizes the full +<dial><number>.
const COUNTRIES = [
  { code: "TM", dial: "993", name: "Türkmenistan" },
  { code: "RU", dial: "7", name: "Россия" },
  { code: "UZ", dial: "998", name: "Oʻzbekiston" },
  { code: "KZ", dial: "7", name: "Qazaqstan" },
  { code: "TR", dial: "90", name: "Türkiye" },
  { code: "AZ", dial: "994", name: "Azərbaycan" },
  { code: "IR", dial: "98", name: "Īrān" },
  { code: "AF", dial: "93", name: "Afġānistān" },
  { code: "AE", dial: "971", name: "UAE" },
  { code: "DE", dial: "49", name: "Deutschland" },
  { code: "GB", dial: "44", name: "United Kingdom" },
  { code: "US", dial: "1", name: "USA" },
];

function PhoneLogin() {
  const t = useTranslations("auth");
  const router = useRouter();
  const params = useSearchParams();
  const [country, setCountry] = useState(COUNTRIES[0].code);
  const [phone, setPhone] = useState("");
  const [code, setCode] = useState("");
  const [sent, setSent] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const dial = COUNTRIES.find((c) => c.code === country)?.dial ?? "993";
  const fullPhone = () => `+${dial}${phone.replace(/\D/g, "")}`;

  const call = async (path: string, extra: Record<string, unknown>) => {
    const res = await fetch(`/api/auth/phone/${path}`, {
      method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ phone: fullPhone(), ...extra }),
    });
    if (!res.ok) {
      const j = (await res.json().catch(() => null)) as { error?: { message?: string } } | null;
      throw new Error(j?.error?.message ?? t("invalid"));
    }
  };
  const start = async () => { setError(null); setBusy(true); try { await call("start", {}); setSent(true); } catch (e) { setError((e as Error).message); } finally { setBusy(false); } };
  const verify = async () => {
    setError(null); setBusy(true);
    try {
      await call("verify", { code });
      const next = params.get("next");
      router.replace(next && next.startsWith("/") ? next : "/");
      router.refresh();
    } catch (e) { setError((e as Error).message); } finally { setBusy(false); }
  };

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <div className="w-36 shrink-0">
          <Field label={t("country")} htmlFor="country">
            <Select value={country} onValueChange={setCountry}>
              <SelectTrigger id="country"><SelectValue /></SelectTrigger>
              <SelectContent>
                {COUNTRIES.map((c) => <SelectItem key={c.code} value={c.code}>{c.name} +{c.dial}</SelectItem>)}
              </SelectContent>
            </Select>
          </Field>
        </div>
        <div className="flex-1">
          <Field label={t("phoneNumber")} htmlFor="phone">
            <Input id="phone" inputMode="tel" autoComplete="tel-national" placeholder="65 65 65 85" value={phone} onChange={(e) => setPhone(e.target.value)} disabled={sent} />
          </Field>
        </div>
      </div>
      {sent ? (
        <Field label={t("codeLabel")} htmlFor="code" hint={t("codeSentHint")}>
          <Input id="code" inputMode="numeric" pattern="[0-9]*" maxLength={8} autoFocus value={code} onChange={(e) => setCode(e.target.value)} />
        </Field>
      ) : null}
      {error ? <Alert variant="destructive"><AlertDescription>{error}</AlertDescription></Alert> : null}
      {!sent ? (
        <Button type="button" className="w-full" disabled={busy || phone.replace(/\D/g, "").length < 6} onClick={start}>{busy ? t("signingIn") : t("sendCode")}</Button>
      ) : (
        <div className="flex gap-2">
          <Button type="button" variant="outline" disabled={busy} onClick={() => { setSent(false); setCode(""); }}>{t("back")}</Button>
          <Button type="button" className="flex-1" disabled={busy || code.length < 4} onClick={verify}>{busy ? t("signingIn") : t("signIn")}</Button>
        </div>
      )}
    </div>
  );
}

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
        <span className="mb-1 flex h-9 w-9 items-center justify-center rounded-md bg-primary text-primary-foreground">
          <BrandMark className="h-5 w-5" />
        </span>
        <CardTitle className="text-xl">{t("title")}</CardTitle>
        <CardDescription>{t("subtitle")}</CardDescription>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue="email">
          <TabsList className="mb-4 grid w-full grid-cols-2">
            <TabsTrigger value="email">{t("emailTab")}</TabsTrigger>
            <TabsTrigger value="phone">{t("phoneTab")}</TabsTrigger>
          </TabsList>
          <TabsContent value="email">
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
          </TabsContent>
          <TabsContent value="phone">
            <PhoneLogin />
          </TabsContent>
        </Tabs>
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

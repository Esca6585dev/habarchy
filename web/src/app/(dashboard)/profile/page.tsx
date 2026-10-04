"use client";
import { useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { api, ApiError, ensureOk, unwrap } from "@/lib/api/client";
import { qk, useMe } from "@/lib/hooks";
import { PageHeader } from "@/components/common/page-header";
import { Field } from "@/components/common/field";
import { CopyButton } from "@/components/common/copy-button";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Avatar } from "@/components/chat/avatar";
import type { User } from "@/lib/api/types";

export default function ProfilePage() {
  const t = useTranslations("profile");
  const tc = useTranslations("common");
  const { data: me } = useMe();
  const qc = useQueryClient();
  const router = useRouter();
  const u = me as User | undefined;
  const [profile, setProfile] = useState({ full_name: "", bio: "" });
  const [avatarId, setAvatarId] = useState<string | null | undefined>(undefined);
  const fileRef = useRef<HTMLInputElement>(null);
  const [pw, setPw] = useState({ current: "", next: "" });
  const [setup, setSetup] = useState<{ secret: string; otpauth_url: string } | null>(null);
  const [code, setCode] = useState("");
  const [disablePw, setDisablePw] = useState("");

  const changePw = useMutation({
    mutationFn: async () => { ensureOk(await api.PUT("/api/admin/me/password", { body: { current_password: pw.current, new_password: pw.next } })); },
    onSuccess: () => { toast.success(t("passwordChanged")); setPw({ current: "", next: "" }); },
    onError: (e: ApiError) => toast.error(e.message),
  });
  const saveProfile = useMutation({
    mutationFn: async () => unwrap<User>(await api.PUT("/api/admin/me/profile", { body: { full_name: profile.full_name || u?.full_name || "", bio: profile.bio, avatar_id: (avatarId === undefined ? u?.avatar_id : avatarId) ?? undefined } })),
    onSuccess: () => { toast.success(t("profileSaved")); void qc.invalidateQueries({ queryKey: qk.me }); },
    onError: (e: ApiError) => toast.error(e.message),
  });
  const uploadAvatar = async (file: File) => {
    const fd = new FormData();
    fd.append("file", file);
    const res = await fetch("/api/backend/api/admin/attachments", { method: "POST", body: fd });
    if (!res.ok) return void toast.error(tc("error"));
    setAvatarId(((await res.json()) as { data: { id: string } }).data.id);
  };
  const start = useMutation({
    mutationFn: async () => unwrap<{ secret: string; otpauth_url: string }>(await api.POST("/api/admin/me/totp/setup")),
    onSuccess: setSetup,
    onError: (e: ApiError) => toast.error(e.message),
  });
  const confirm = useMutation({
    mutationFn: async () => { ensureOk(await api.POST("/api/admin/me/totp/confirm", { body: { code } })); },
    onSuccess: () => { toast.success(t("twoFactorOn")); setSetup(null); setCode(""); void qc.invalidateQueries({ queryKey: qk.me }); },
    onError: (e: ApiError) => toast.error(e.message),
  });
  const disable = useMutation({
    mutationFn: async () => { ensureOk(await api.POST("/api/admin/me/totp/disable", { body: { password: disablePw } })); },
    onSuccess: () => { toast.success(t("twoFactorOff")); setDisablePw(""); void qc.invalidateQueries({ queryKey: qk.me }); },
    onError: (e: ApiError) => toast.error(e.message),
  });
  const logoutAll = useMutation({
    mutationFn: async () => { await api.POST("/api/admin/auth/logout-all"); await fetch("/api/auth/logout", { method: "POST" }); },
    onSuccess: () => { router.replace("/login"); router.refresh(); },
  });

  return (
    <>
      <PageHeader title={t("title")} description={me?.email} />
      <div className="grid gap-6 lg:grid-cols-2">
        <Card className="lg:col-span-2">
          <CardHeader><CardTitle>{t("profileTitle")}</CardTitle></CardHeader>
          <CardContent>
            <div className="flex flex-col gap-4 sm:flex-row sm:items-start">
              <div className="flex flex-col items-center gap-2">
                <Avatar name={profile.full_name || u?.full_name} avatarId={avatarId === undefined ? u?.avatar_id : avatarId} size={72} />
                <input ref={fileRef} type="file" accept="image/*" className="hidden" onChange={(e) => { const f = e.target.files?.[0]; if (f) void uploadAvatar(f); e.target.value = ""; }} />
                <Button size="sm" variant="outline" onClick={() => fileRef.current?.click()}>{t("uploadPhoto")}</Button>
              </div>
              <div className="flex-1 space-y-3">
                <Field label={t("name")} htmlFor="pf-name"><Input id="pf-name" value={profile.full_name || u?.full_name || ""} onChange={(e) => setProfile({ ...profile, full_name: e.target.value })} /></Field>
                <Field label={t("bio")} htmlFor="pf-bio"><Textarea id="pf-bio" rows={2} value={profile.bio || u?.bio || ""} onChange={(e) => setProfile({ ...profile, bio: e.target.value })} /></Field>
                <Button onClick={() => saveProfile.mutate()} disabled={saveProfile.isPending}>{tc("save")}</Button>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader><CardTitle>{t("changePassword")}</CardTitle></CardHeader>
          <CardContent>
            <form className="space-y-3" onSubmit={(e) => { e.preventDefault(); changePw.mutate(); }}>
              <Field label={t("currentPassword")} htmlFor="cur"><Input id="cur" type="password" autoComplete="current-password" value={pw.current} onChange={(e) => setPw({ ...pw, current: e.target.value })} /></Field>
              <Field label={t("newPassword")} htmlFor="new"><Input id="new" type="password" autoComplete="new-password" minLength={8} value={pw.next} onChange={(e) => setPw({ ...pw, next: e.target.value })} /></Field>
              <Button type="submit" disabled={changePw.isPending || pw.next.length < 8}>{tc("save")}</Button>
            </form>
            <div className="mt-6 border-t pt-4">
              <Button variant="outline" onClick={() => logoutAll.mutate()}>{t("sessions")}</Button>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">{t("twoFactor")} {me?.totp_enabled ? <Badge variant="success">{t("enabled")}</Badge> : <Badge variant="outline">{t("disabled")}</Badge>}</CardTitle>
            <CardDescription>TOTP (Google Authenticator, Aegis, 1Password…)</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {!me?.totp_enabled && !setup ? <Button onClick={() => start.mutate()} disabled={start.isPending}>{t("enable")}</Button> : null}
            {setup ? (
              <div className="space-y-3">
                <p className="text-sm text-muted-foreground">{t("scan")}</p>
                <div className="flex items-center gap-2"><code className="break-all rounded bg-muted px-2 py-1 font-mono text-xs">{setup.secret}</code><CopyButton value={setup.secret} /></div>
                <a href={setup.otpauth_url} className="block truncate text-xs text-primary underline">{setup.otpauth_url}</a>
                <form className="flex items-end gap-2" onSubmit={(e) => { e.preventDefault(); confirm.mutate(); }}>
                  <Field label={t("code")} htmlFor="totp-code"><Input id="totp-code" inputMode="numeric" maxLength={6} value={code} onChange={(e) => setCode(e.target.value)} /></Field>
                  <Button type="submit" disabled={code.length !== 6 || confirm.isPending}>{t("confirm")}</Button>
                </form>
              </div>
            ) : null}
            {me?.totp_enabled ? (
              <form className="flex items-end gap-2" onSubmit={(e) => { e.preventDefault(); disable.mutate(); }}>
                <Field label={t("currentPassword")} htmlFor="dis-pw"><Input id="dis-pw" type="password" value={disablePw} onChange={(e) => setDisablePw(e.target.value)} /></Field>
                <Button type="submit" variant="destructive" disabled={!disablePw || disable.isPending}>{t("disable")}</Button>
              </form>
            ) : null}
          </CardContent>
        </Card>
      </div>
    </>
  );
}

"use client";
import { useTranslations } from "next-intl";
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription,
  AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button, type ButtonProps } from "@/components/ui/button";

export function ConfirmButton({
  title, description, onConfirm, children, confirmLabel, ...btn
}: ButtonProps & { title: string; description?: string; onConfirm: () => void; confirmLabel?: string }) {
  const t = useTranslations("common");
  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button {...btn}>{children}</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription>{description ?? t("confirmDeleteText")}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{t("cancel")}</AlertDialogCancel>
          <AlertDialogAction onClick={onConfirm}>{confirmLabel ?? t("delete")}</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

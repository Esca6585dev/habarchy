#!/usr/bin/env python3
"""Virtual phone for the Habarchy SMS gateway.

Pretends to be the Android Gateway app so the server side can be tested
without building or installing the APK. It authenticates with a gateway
key, long-polls the outbox, "sends" each SMS (prints it instead of using a
real radio), reports the result back, and heartbeats so the provider shows
as online. It can also forward a fake inbound SMS.

Only the Python standard library is used, so no install step is needed:

    python3 scripts/virtual_phone.py \
        --url https://habarchy.example.tm \
        --key gw_xxxxxxxxxxxxxxxx

Common options:

    --once            handle one lease batch, then exit (good for CI)
    --fail            report every SMS as failed (test fallback/retry)
    --fail-rate 0.3   fail ~30% of messages at random
    --send-delay 2    seconds to "transmit" each SMS before reporting
    --delivered       also send a delivery receipt after "sent"
    --inbound "+99365...:Salam"   forward one fake received SMS and exit
    --wait 20         long-poll seconds per outbox call (server caps it)

The key is the android_sms provider's gateway_key. Get it from the admin
panel (provider -> Pairing), or from the habarchy://gateway?...&key=... QR.
"""

from __future__ import annotations

import argparse
import json
import random
import signal
import sys
import time
import urllib.error
import urllib.request
from typing import Any

API_PREFIX = "/api/gateway/v1"

_stop = False


def _on_signal(_signum: int, _frame: Any) -> None:
    global _stop
    _stop = True
    print("\n… togtadylýar (signal). Soňky işini tamamlaýar.", flush=True)


def _color(txt: str, code: str, enabled: bool) -> str:
    return f"\033[{code}m{txt}\033[0m" if enabled else txt


class Phone:
    """A minimal gateway client speaking the phone-facing API."""

    def __init__(self, base_url: str, key: str, *, timeout: int, color: bool) -> None:
        self.base = base_url.rstrip("/")
        self.key = key.strip()
        self.timeout = timeout
        self.color = color

    def _request(self, method: str, path: str, body: Any = None, *, read_timeout: int | None = None) -> Any:
        url = self.base + API_PREFIX + path
        data = None
        headers = {"X-Gateway-Key": self.key, "Accept": "application/json"}
        if body is not None:
            data = json.dumps(body).encode("utf-8")
            headers["Content-Type"] = "application/json"
        req = urllib.request.Request(url, data=data, headers=headers, method=method)
        try:
            with urllib.request.urlopen(req, timeout=read_timeout or self.timeout) as resp:
                raw = resp.read()
                if not raw:
                    return None
                return json.loads(raw)
        except urllib.error.HTTPError as exc:
            detail = exc.read().decode("utf-8", "replace")
            raise ApiError(exc.code, detail) from exc
        except urllib.error.URLError as exc:
            raise ConnError(str(exc.reason)) from exc

    # --- endpoint wrappers -------------------------------------------------
    def me(self) -> dict:
        return self._request("GET", "/me")

    def heartbeat(self, info: dict) -> dict:
        return self._request("POST", "/heartbeat", info)

    def lease(self, wait: int, limit: int) -> list:
        # Give the HTTP read a little more than the server's long-poll window
        # so the connection is not torn down mid-wait.
        return self._request("GET", f"/outbox?wait={wait}&limit={limit}", read_timeout=wait + 15) or []

    def report(self, outbox_id: str, status: str, *, error_code: str = "", error_message: str = "", parts: int = 1) -> dict:
        return self._request(
            "POST",
            f"/outbox/{outbox_id}/result",
            {"status": status, "error_code": error_code, "error_message": error_message, "parts": parts},
        )

    def inbound(self, frm: str, text: str) -> dict:
        return self._request("POST", "/inbound", {"from": frm, "text": text})

    # --- pretty logging ----------------------------------------------------
    def log(self, tag: str, msg: str, code: str = "0") -> None:
        ts = time.strftime("%H:%M:%S")
        print(f"{_color(ts, '90', self.color)} {_color(tag, code, self.color)} {msg}", flush=True)


class ApiError(Exception):
    def __init__(self, status: int, body: str) -> None:
        super().__init__(f"HTTP {status}: {body}")
        self.status = status
        self.body = body


class ConnError(Exception):
    pass


def parts_for(text: str) -> int:
    """Rough SMS part count (GSM-7: 160, or 70 when non-Latin chars appear)."""
    limit = 160 if text.isascii() else 70
    per = 153 if text.isascii() else 67  # concatenated-SMS payload per part
    if len(text) <= limit:
        return 1
    return (len(text) + per - 1) // per


def handle_batch(phone: Phone, items: list, args: argparse.Namespace) -> int:
    done = 0
    for it in items:
        oid = it["id"]
        to = it.get("to", "?")
        text = it.get("text", "")
        slot = it.get("sim_slot", -1)
        sim = "default" if slot == -1 else f"SIM{slot}"
        phone.log("SMS →", f"{_color(to, '1;36', phone.color)} [{sim}] {text!r}", "33")

        if args.send_delay:
            time.sleep(args.send_delay)

        fail = args.fail or (args.fail_rate > 0 and random.random() < args.fail_rate)
        n = parts_for(text)
        try:
            if fail:
                phone.report(oid, "failed", error_code="radio_off", error_message="virtual phone: forced failure", parts=n)
                phone.log("  ✗", f"failed hasabat berildi (outbox {oid})", "31")
            else:
                phone.report(oid, "sent", parts=n)
                phone.log("  ✓", f"sent hasabat berildi · {n} bölek (outbox {oid})", "32")
                if args.delivered:
                    time.sleep(0.5)
                    phone.report(oid, "delivered", parts=n)
                    phone.log("  ✓✓", "delivered (ýetirildi) hasabat berildi", "32")
            done += 1
        except ApiError as exc:
            phone.log("  !", f"hasabat ýalňyşlygy: {exc}", "31")
    return done


def main() -> int:
    ap = argparse.ArgumentParser(
        description="Habarchy SMS gateway üçin sanal telefon (test guraly).",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__,
    )
    ap.add_argument("--url", required=True, help="Serweriň baza URL-i, meselem https://habarchy.example.tm")
    ap.add_argument("--key", required=True, help="android_sms provaýderiň gateway_key-i (gw_…)")
    ap.add_argument("--wait", type=int, default=20, help="Her outbox soragyndaky long-poll sekundy (default 20)")
    ap.add_argument("--limit", type=int, default=5, help="Bir gezekde alynýan SMS sany (default 5)")
    ap.add_argument("--send-delay", type=float, default=1.0, help="Her SMS-i 'ibermek' üçin sekunt (default 1)")
    ap.add_argument("--heartbeat", type=int, default=30, help="Heartbeat aralygy, sekunt (default 30; 0 = öçük)")
    ap.add_argument("--once", action="store_true", help="Bir gezek alyp işle, soň çyk")
    ap.add_argument("--fail", action="store_true", help="Ähli SMS-i failed hasabat ber")
    ap.add_argument("--fail-rate", type=float, default=0.0, help="Tötänleýin näsazlyk paýy, 0..1")
    ap.add_argument("--delivered", action="store_true", help="sent-den soň delivery receipt-de iber")
    ap.add_argument("--inbound", metavar="FROM:TEXT", help="Bir sany galp gelen SMS ugrat we çyk (meselem '+99365123456:Salam')")
    ap.add_argument("--no-color", action="store_true", help="Reňkli çykyşy öçür")
    args = ap.parse_args()

    phone = Phone(args.url, args.key, timeout=30, color=not args.no_color and sys.stdout.isatty())

    signal.signal(signal.SIGINT, _on_signal)
    signal.signal(signal.SIGTERM, _on_signal)

    # 1) Verify pairing.
    try:
        info = phone.me()
    except ConnError as exc:
        phone.log("BAGLANYŞYK", f"serwere ýetip bolmady: {exc}", "31")
        phone.log("", "  → --url dogrumy? Serwer işleýärmi? localhost bolsa telefon ýetip bilmez.", "90")
        return 2
    except ApiError as exc:
        phone.log("AWTH", f"jübütleme şowsuz: {exc}", "31")
        if exc.status == 401:
            phone.log("", "  → gateway_key nädogry ýa provaýder/proýekt öçük.", "90")
        return 2

    phone.log("JÜBÜTLENDI", f"provaýder={info.get('provider')!r} proýekt={info.get('project_id')}", "1;32")
    phone.log(
        "SAZLAMA",
        f"inbound_enabled={info.get('inbound_enabled')} sim_slot={info.get('sim_slot')} nobatda={info.get('pending')}",
        "36",
    )

    # One-shot inbound mode.
    if args.inbound:
        if ":" not in args.inbound:
            phone.log("INBOUND", "format: --inbound '+99365...:habaryň teksti'", "31")
            return 2
        frm, text = args.inbound.split(":", 1)
        if not info.get("inbound_enabled", False):
            phone.log("INBOUND", "serwerde inbound_enabled=false — ilki provaýderde açmaly.", "31")
            return 2
        try:
            res = phone.inbound(frm.strip(), text.strip())
            phone.log("INBOUND ←", f"ugradyldy id={res.get('id')} received_at={res.get('received_at')}", "32")
            return 0
        except ApiError as exc:
            phone.log("INBOUND", f"şowsuz: {exc}", "31")
            return 1

    phone.log("DIŇLEÝÄR", f"outbox long-poll {args.wait}s · limit {args.limit} · Ctrl+C bilen duruz", "1;34")

    last_hb = 0.0
    total = 0
    while not _stop:
        now = time.monotonic()
        if args.heartbeat and now - last_hb >= args.heartbeat:
            try:
                hb = phone.heartbeat({"model": "virtual-phone", "app": "virtual_phone.py", "battery": 100})
                last_hb = now
                if hb and hb.get("pending"):
                    phone.log("♥", f"heartbeat · nobatda {hb['pending']}", "90")
            except (ApiError, ConnError) as exc:
                phone.log("♥", f"heartbeat şowsuz: {exc}", "31")

        try:
            items = phone.lease(args.wait, args.limit)
        except ConnError as exc:
            phone.log("OUTBOX", f"baglanyşyk üzüldi, 3s garaşýar: {exc}", "31")
            time.sleep(3)
            continue
        except ApiError as exc:
            phone.log("OUTBOX", f"ýalňyşlyk: {exc}", "31")
            time.sleep(3)
            continue

        if items:
            phone.log("OUTBOX", f"{len(items)} SMS alyndy", "1;33")
            total += handle_batch(phone, items, args)

        if args.once:
            break

    phone.log("TAMAM", f"jemi {total} SMS işlenildi", "1;32")
    return 0


if __name__ == "__main__":
    sys.exit(main())

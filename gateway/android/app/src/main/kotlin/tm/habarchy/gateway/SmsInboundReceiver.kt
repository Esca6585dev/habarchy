package tm.habarchy.gateway

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.provider.Telephony
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale
import java.util.TimeZone

/** Forwards SMS the phone receives to the project webhook (sms.inbound), when enabled. */
class SmsInboundReceiver : BroadcastReceiver() {
    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action != Telephony.Sms.Intents.SMS_RECEIVED_ACTION) return
        val messages = Telephony.Sms.Intents.getMessagesFromIntent(intent) ?: return
        if (messages.isEmpty()) return
        val from = messages[0].displayOriginatingAddress ?: return
        val text = messages.joinToString("") { it.messageBody ?: "" }
        if (!InboundControl.effective(context)) return // switched off: never forward
        val fmt = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'", Locale.US).apply { timeZone = TimeZone.getTimeZone("UTC") }
        val ts = messages[0].timestampMillis
        Reporter.inbound(context, from, text, fmt.format(Date(ts)), ts)
        val store = GatewayStore(context)
        if (ts > store.lastInboundSyncAt) store.lastInboundSyncAt = ts
    }
}

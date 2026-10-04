package tm.habarchy.gateway

import android.app.PendingIntent
import android.content.Context
import android.content.Intent
import android.os.Build
import android.telephony.SmsManager
import android.telephony.SubscriptionManager

/** Sends one SMS (multipart when needed) and asks for SENT / DELIVERED broadcasts. */
class SmsSender(private val context: Context) {

    private fun manager(simSlot: Int): SmsManager {
        if (simSlot >= 0 && Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP_MR1) {
            try {
                val sm = context.getSystemService(Context.TELEPHONY_SUBSCRIPTION_SERVICE) as SubscriptionManager
                val info = sm.activeSubscriptionInfoList?.firstOrNull { it.simSlotIndex == simSlot }
                if (info != null) {
                    return if (Build.VERSION.SDK_INT >= 31) {
                        context.getSystemService(SmsManager::class.java).createForSubscriptionId(info.subscriptionId)
                    } else {
                        @Suppress("DEPRECATION")
                        SmsManager.getSmsManagerForSubscriptionId(info.subscriptionId)
                    }
                }
            } catch (_: SecurityException) {
                // READ_PHONE_STATE denied: fall through to the default SIM.
            }
        }
        return if (Build.VERSION.SDK_INT >= 31) context.getSystemService(SmsManager::class.java)
        else @Suppress("DEPRECATION") SmsManager.getDefault()
    }

    /** Returns the number of parts. Results arrive in [SmsResultReceiver]. */
    fun send(to: String, text: String, simSlot: Int, outboxId: String): Int {
        require(to.isNotBlank()) { "empty recipient" }
        val sms = manager(simSlot)
        val parts = sms.divideMessage(text)
        val sentIntents = ArrayList<PendingIntent>()
        val deliveredIntents = ArrayList<PendingIntent>()
        for (i in parts.indices) {
            sentIntents.add(pending(SmsResultReceiver.ACTION_SENT, outboxId, i, parts.size))
            deliveredIntents.add(pending(SmsResultReceiver.ACTION_DELIVERED, outboxId, i, parts.size))
        }
        SmsResultReceiver.expect(outboxId, parts.size)
        if (parts.size == 1) {
            sms.sendTextMessage(to, null, parts[0], sentIntents[0], deliveredIntents[0])
        } else {
            sms.sendMultipartTextMessage(to, null, parts, sentIntents, deliveredIntents)
        }
        return parts.size
    }

    private fun pending(action: String, outboxId: String, part: Int, total: Int): PendingIntent {
        val intent = Intent(action).setPackage(context.packageName)
            .putExtra(SmsResultReceiver.EXTRA_ID, outboxId)
            .putExtra(SmsResultReceiver.EXTRA_PART, part)
            .putExtra(SmsResultReceiver.EXTRA_TOTAL, total)
        val requestCode = (outboxId.hashCode() * 31 + part * 2 + if (action == SmsResultReceiver.ACTION_SENT) 0 else 1)
        val flags = PendingIntent.FLAG_UPDATE_CURRENT or (if (Build.VERSION.SDK_INT >= 23) PendingIntent.FLAG_IMMUTABLE else 0)
        return PendingIntent.getBroadcast(context, requestCode, intent, flags)
    }
}

package tm.habarchy.gateway

import android.app.Activity
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.telephony.SmsManager
import java.util.concurrent.ConcurrentHashMap

/**
 * Receives the SENT / DELIVERED results of the SMS parts and reports once
 * per outbox entry to the server (via [Reporter]).
 */
class SmsResultReceiver : BroadcastReceiver() {

    private class Progress(val total: Int) {
        var sentParts = 0
        var deliveredParts = 0
        var failedCode: String? = null
        var sentReported = false
        var deliveredReported = false
    }

    override fun onReceive(context: Context, intent: Intent) {
        val id = intent.getStringExtra(EXTRA_ID) ?: return
        val total = intent.getIntExtra(EXTRA_TOTAL, 1)
        val progress = inFlight.getOrPut(id) { Progress(total) }
        val store = GatewayStore(context)
        synchronized(progress) {
            when (intent.action) {
                ACTION_SENT -> {
                    val code = resultCode
                    if (code == Activity.RESULT_OK) {
                        progress.sentParts++
                    } else if (progress.failedCode == null) {
                        progress.failedCode = codeName(code)
                    }
                    val finished = progress.sentParts + (if (progress.failedCode != null) 1 else 0) >= progress.total || progress.failedCode != null
                    if (finished && !progress.sentReported) {
                        progress.sentReported = true
                        if (progress.failedCode == null) {
                            store.inc("sent")
                            store.log("✓ sent $id (${progress.total} part${if (progress.total > 1) "s" else ""})")
                            Reporter.report(context, id, "sent", null, null, progress.total)
                        } else {
                            store.inc("failed")
                            store.log("✕ failed $id: ${progress.failedCode}")
                            Reporter.report(context, id, "failed", progress.failedCode, "SmsManager: ${progress.failedCode}", progress.total)
                            inFlight.remove(id)
                        }
                    }
                }
                ACTION_DELIVERED -> {
                    if (resultCode == Activity.RESULT_OK) progress.deliveredParts++
                    if (progress.deliveredParts >= progress.total && !progress.deliveredReported) {
                        progress.deliveredReported = true
                        store.inc("delivered")
                        store.log("✓✓ delivered $id")
                        Reporter.report(context, id, "delivered", null, null, progress.total)
                        inFlight.remove(id)
                    }
                }
            }
        }
    }

    companion object {
        const val ACTION_SENT = "tm.habarchy.gateway.SMS_SENT"
        const val ACTION_DELIVERED = "tm.habarchy.gateway.SMS_DELIVERED"
        const val EXTRA_ID = "outbox_id"
        const val EXTRA_PART = "part"
        const val EXTRA_TOTAL = "total"

        private val inFlight = ConcurrentHashMap<String, Progress>()

        fun expect(id: String, parts: Int) {
            inFlight[id] = Progress(parts)
        }

        /** Server-side error codes (see androidgw.Retryable in the backend). */
        fun codeName(code: Int): String = when (code) {
            SmsManager.RESULT_ERROR_GENERIC_FAILURE -> "generic_failure"
            SmsManager.RESULT_ERROR_RADIO_OFF -> "radio_off"
            SmsManager.RESULT_ERROR_NULL_PDU -> "null_pdu"
            SmsManager.RESULT_ERROR_NO_SERVICE -> "no_service"
            SmsManager.RESULT_ERROR_LIMIT_EXCEEDED -> "limit_exceeded"
            SmsManager.RESULT_ERROR_SHORT_CODE_NOT_ALLOWED -> "short_code_not_allowed"
            SmsManager.RESULT_ERROR_SHORT_CODE_NEVER_ALLOWED -> "short_code_never_allowed"
            SmsManager.RESULT_NETWORK_REJECT -> "network_reject"
            SmsManager.RESULT_INVALID_ARGUMENTS -> "invalid_number"
            else -> "error_$code"
        }
    }
}

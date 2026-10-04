package tm.habarchy.gateway

import android.content.Context
import java.util.concurrent.Executors

/** Posts results to the server from a single background thread, retrying briefly. */
object Reporter {
    private val executor = Executors.newSingleThreadExecutor()

    fun report(context: Context, outboxId: String, status: String, errorCode: String?, errorMessage: String?, parts: Int) {
        val app = context.applicationContext
        executor.execute {
            val store = GatewayStore(app)
            val api = Api(store.url, store.key, GatewayService.appVersion(app))
            var delay = 2_000L
            repeat(5) { attempt ->
                try {
                    api.report(outboxId, status, errorCode, errorMessage, parts)
                    return@execute
                } catch (e: Exception) {
                    store.lastError = "report: ${e.message}"
                    if (attempt == 4) store.log("! could not report $status for $outboxId: ${e.message}")
                    Thread.sleep(delay)
                    delay = (delay * 2).coerceAtMost(30_000L)
                }
            }
        }
    }

    fun inbound(context: Context, from: String, text: String, receivedAtIso: String, timestampMillis: Long = 0L) {
        val app = context.applicationContext
        executor.execute {
            val store = GatewayStore(app)
            if (!InboundControl.effective(app, store) || store.url.isBlank()) return@execute
            if (!store.markForwarded("$from|${timestampMillis.takeIf { it > 0 } ?: receivedAtIso}|${text.hashCode()}")) return@execute
            try {
                Api(store.url, store.key, GatewayService.appVersion(app)).inbound(from, text, receivedAtIso)
                store.inc("inbound")
                store.log("↓ inbound from $from forwarded")
            } catch (e: Exception) {
                store.log("! inbound forward failed: ${e.message}")
            }
        }
    }
}

package tm.habarchy.gateway

import android.content.Context
import android.content.SharedPreferences
import org.json.JSONArray
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

/** Settings, counters and a small log ring buffer in SharedPreferences. */
class GatewayStore(context: Context) {
    private val p: SharedPreferences = context.applicationContext.getSharedPreferences("habarchy_gateway", Context.MODE_PRIVATE)

    var url: String
        get() = p.getString("url", "") ?: ""
        set(v) = p.edit().putString("url", v).apply()
    var key: String
        get() = p.getString("key", "") ?: ""
        set(v) = p.edit().putString("key", v).apply()
    var enabled: Boolean
        get() = p.getBoolean("enabled", false)
        set(v) = p.edit().putBoolean("enabled", v).apply()
    var forwardInbound: Boolean
        get() = p.getBoolean("forward_inbound", false)
        set(v) = p.edit().putBoolean("forward_inbound", v).apply()
    /** inbound_enabled of the provider as the server last reported it. */
    var serverInboundEnabled: Boolean
        get() = p.getBoolean("server_inbound_enabled", true)
        set(v) = p.edit().putBoolean("server_inbound_enabled", v).apply()
    /** Inbox messages older than this are never forwarded (set when forwarding is switched on). */
    var lastInboundSyncAt: Long
        get() = p.getLong("last_inbound_sync_at", 0)
        set(v) = p.edit().putLong("last_inbound_sync_at", v).apply()
    var inboundActive: Boolean
        get() = p.getBoolean("inbound_active", false)
        set(v) = p.edit().putBoolean("inbound_active", v).apply()

    /** Remembers forwarded SMS keys so the receiver and the inbox sync never send one twice. */
    @Synchronized
    fun markForwarded(key: String): Boolean {
        val arr = JSONArray(p.getString("forwarded", "[]") ?: "[]")
        for (i in 0 until arr.length()) if (arr.getString(i) == key) return false
        arr.put(key)
        while (arr.length() > 200) arr.remove(0)
        p.edit().putString("forwarded", arr.toString()).apply()
        return true
    }

    var lastPollAt: Long
        get() = p.getLong("last_poll_at", 0)
        set(v) = p.edit().putLong("last_poll_at", v).apply()
    var lastHeartbeatAt: Long
        get() = p.getLong("last_heartbeat_at", 0)
        set(v) = p.edit().putLong("last_heartbeat_at", v).apply()
    var lastError: String
        get() = p.getString("last_error", "") ?: ""
        set(v) = p.edit().putString("last_error", v).apply()
    var connected: Boolean
        get() = p.getBoolean("connected", false)
        set(v) = p.edit().putBoolean("connected", v).apply()
    var providerName: String
        get() = p.getString("provider_name", "") ?: ""
        set(v) = p.edit().putString("provider_name", v).apply()
    var pending: Long
        get() = p.getLong("pending", 0)
        set(v) = p.edit().putLong("pending", v).apply()

    fun inc(counter: String) {
        p.edit().putLong(counter, p.getLong(counter, 0) + 1).apply()
    }

    fun statusMap(running: Boolean): Map<String, Any?> = mapOf(
        "running" to running,
        "enabled" to enabled,
        "connected" to connected,
        "providerName" to providerName,
        "pending" to pending,
        "sent" to p.getLong("sent", 0),
        "failed" to p.getLong("failed", 0),
        "delivered" to p.getLong("delivered", 0),
        "inbound" to p.getLong("inbound", 0),
        "forwardInbound" to forwardInbound,
        "serverInboundEnabled" to serverInboundEnabled,
        "inboundActive" to inboundActive,
        "lastPollAt" to lastPollAt,
        "lastHeartbeatAt" to lastHeartbeatAt,
        "lastError" to lastError,
    )

    @Synchronized
    fun log(line: String) {
        val arr = JSONArray(p.getString("log", "[]") ?: "[]")
        val stamp = SimpleDateFormat("HH:mm:ss", Locale.US).format(Date())
        arr.put("$stamp  $line")
        while (arr.length() > MAX_LOG) arr.remove(0)
        p.edit().putString("log", arr.toString()).apply()
    }

    fun log(): List<String> {
        val arr = JSONArray(p.getString("log", "[]") ?: "[]")
        return (0 until arr.length()).map { arr.getString(it) }.reversed()
    }

    fun clearLog() = p.edit().putString("log", "[]").apply()

    companion object {
        private const val MAX_LOG = 200
    }
}

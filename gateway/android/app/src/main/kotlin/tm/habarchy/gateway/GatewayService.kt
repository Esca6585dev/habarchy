package tm.habarchy.gateway

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.app.Service
import android.content.Context
import android.content.Intent
import android.content.IntentFilter
import android.content.pm.ServiceInfo
import android.net.ConnectivityManager
import android.net.NetworkCapabilities
import android.os.BatteryManager
import android.os.Build
import android.os.IBinder
import android.os.PowerManager
import android.telephony.TelephonyManager
import androidx.core.app.NotificationCompat
import org.json.JSONObject
import java.util.concurrent.atomic.AtomicBoolean

/**
 * Foreground service: long-polls the Habarchy outbox and sends the SMS it
 * receives through the phone's SIM. Survives the UI being closed; restarts
 * on boot when enabled.
 */
class GatewayService : Service() {
    private val running = AtomicBoolean(false)
    private var worker: Thread? = null
    private var wakeLock: PowerManager.WakeLock? = null

    override fun onBind(intent: Intent?): IBinder? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        if (intent?.action == ACTION_STOP) {
            stopSelf()
            return START_NOT_STICKY
        }
        startForegroundCompat(buildNotification(getString(R.string.notif_connecting)))
        if (running.compareAndSet(false, true)) {
            isRunning = true
            acquireWakeLock()
            worker = Thread({ loop() }, "habarchy-gateway").also { it.start() }
        }
        return START_STICKY
    }

    override fun onDestroy() {
        running.set(false)
        isRunning = false
        worker?.interrupt()
        wakeLock?.let { if (it.isHeld) it.release() }
        GatewayStore(this).apply { connected = false }
        super.onDestroy()
    }

    private fun loop() {
        val store = GatewayStore(this)
        val sender = SmsSender(this)
        var backoff = 2_000L
        store.log("service started")
        while (running.get()) {
            if (store.url.isBlank() || store.key.isBlank()) {
                store.lastError = "not configured"
                updateNotification(getString(R.string.notif_not_configured))
                sleepQuiet(5_000)
                continue
            }
            val api = Api(store.url, store.key, appVersion(this))
            try {
                if (System.currentTimeMillis() - store.lastHeartbeatAt > HEARTBEAT_MS) {
                    val data = api.heartbeat(deviceInfo())
                    store.lastHeartbeatAt = System.currentTimeMillis()
                    store.pending = data.optLong("pending", 0)
                    if (store.providerName.isBlank()) {
                        store.providerName = api.me().optString("provider", "")
                    }
                }
                val items = api.lease(POLL_WAIT_SEC, 5)
                store.lastPollAt = System.currentTimeMillis()
                if (!store.connected) {
                    store.connected = true
                    store.lastError = ""
                    store.log("connected to ${store.url}" + if (store.providerName.isNotBlank()) " as \"${store.providerName}\"" else "")
                }
                updateNotification(getString(R.string.notif_online, store.providerName.ifBlank { "Habarçy" }))
                backoff = 2_000L
                for (i in 0 until items.length()) {
                    val item = items.getJSONObject(i)
                    val id = item.getString("id")
                    val to = item.getString("to")
                    val text = item.getString("text")
                    val simSlot = item.optInt("sim_slot", -1)
                    try {
                        val parts = sender.send(to, text, simSlot, id)
                        store.log("→ sending $id to ${mask(to)} ($parts part${if (parts > 1) "s" else ""})")
                    } catch (e: Exception) {
                        store.inc("failed")
                        store.log("✕ cannot send $id: ${e.message}")
                        Reporter.report(this, id, "failed", "send_exception", e.message ?: e.javaClass.simpleName, 0)
                    }
                }
            } catch (e: InterruptedException) {
                break
            } catch (e: Api.HttpError) {
                store.connected = false
                store.lastError = e.message ?: "http error"
                store.log("! ${e.message}")
                updateNotification(if (e.code == 401) getString(R.string.notif_unauthorized) else getString(R.string.notif_offline))
                sleepQuiet(if (e.code == 401) 30_000 else backoff)
                backoff = (backoff * 2).coerceAtMost(60_000L)
            } catch (e: Exception) {
                store.connected = false
                store.lastError = e.message ?: e.javaClass.simpleName
                store.log("! ${e.javaClass.simpleName}: ${e.message}")
                updateNotification(getString(R.string.notif_offline))
                sleepQuiet(backoff)
                backoff = (backoff * 2).coerceAtMost(60_000L)
            }
        }
        store.log("service stopped")
    }

    private fun deviceInfo(): JSONObject {
        val bm = getSystemService(Context.BATTERY_SERVICE) as BatteryManager
        val battery = bm.getIntProperty(BatteryManager.BATTERY_PROPERTY_CAPACITY)
        val cm = getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager
        val caps = cm.getNetworkCapabilities(cm.activeNetwork)
        val network = when {
            caps == null -> "none"
            caps.hasTransport(NetworkCapabilities.TRANSPORT_WIFI) -> "wifi"
            caps.hasTransport(NetworkCapabilities.TRANSPORT_CELLULAR) -> "cellular"
            else -> "other"
        }
        val tm = getSystemService(Context.TELEPHONY_SERVICE) as TelephonyManager
        val operator = try { tm.networkOperatorName ?: "" } catch (_: Exception) { "" }
        return JSONObject()
            .put("battery", battery)
            .put("network", network)
            .put("operator", operator)
            .put("model", "${Build.MANUFACTURER} ${Build.MODEL}")
            .put("android", Build.VERSION.RELEASE)
            .put("app_version", appVersion(this))
    }

    private fun sleepQuiet(ms: Long) {
        try { Thread.sleep(ms) } catch (_: InterruptedException) { running.set(false) }
    }

    private fun acquireWakeLock() {
        val pm = getSystemService(Context.POWER_SERVICE) as PowerManager
        wakeLock = pm.newWakeLock(PowerManager.PARTIAL_WAKE_LOCK, "habarchy:gateway").also { it.acquire() }
    }

    private fun startForegroundCompat(n: Notification) {
        if (Build.VERSION.SDK_INT >= 34) {
            startForeground(NOTIF_ID, n, ServiceInfo.FOREGROUND_SERVICE_TYPE_DATA_SYNC)
        } else {
            startForeground(NOTIF_ID, n)
        }
    }

    private fun buildNotification(text: String): Notification {
        val nm = getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
        if (Build.VERSION.SDK_INT >= 26 && nm.getNotificationChannel(CHANNEL_ID) == null) {
            nm.createNotificationChannel(NotificationChannel(CHANNEL_ID, getString(R.string.notif_channel), NotificationManager.IMPORTANCE_LOW))
        }
        val open = PendingIntent.getActivity(
            this, 0, Intent(this, MainActivity::class.java),
            PendingIntent.FLAG_UPDATE_CURRENT or (if (Build.VERSION.SDK_INT >= 23) PendingIntent.FLAG_IMMUTABLE else 0),
        )
        val stop = PendingIntent.getService(
            this, 1, Intent(this, GatewayService::class.java).setAction(ACTION_STOP),
            PendingIntent.FLAG_UPDATE_CURRENT or (if (Build.VERSION.SDK_INT >= 23) PendingIntent.FLAG_IMMUTABLE else 0),
        )
        return NotificationCompat.Builder(this, CHANNEL_ID)
            .setSmallIcon(android.R.drawable.stat_notify_chat)
            .setContentTitle(getString(R.string.app_name))
            .setContentText(text)
            .setContentIntent(open)
            .addAction(0, getString(R.string.notif_stop), stop)
            .setOngoing(true)
            .setOnlyAlertOnce(true)
            .build()
    }

    private fun updateNotification(text: String) {
        val nm = getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
        nm.notify(NOTIF_ID, buildNotification(text))
    }

    private fun mask(to: String): String = if (to.length > 6) to.take(4) + "***" + to.takeLast(2) else "***"

    companion object {
        const val CHANNEL_ID = "habarchy_gateway"
        const val ACTION_STOP = "tm.habarchy.gateway.STOP"
        private const val NOTIF_ID = 7701
        private const val POLL_WAIT_SEC = 20
        private const val HEARTBEAT_MS = 60_000L

        @Volatile
        var isRunning: Boolean = false
            private set

        fun start(context: Context) {
            val intent = Intent(context, GatewayService::class.java)
            if (Build.VERSION.SDK_INT >= 26) context.startForegroundService(intent) else context.startService(intent)
        }

        fun stop(context: Context) {
            context.stopService(Intent(context, GatewayService::class.java))
        }

        fun appVersion(context: Context): String = try {
            context.packageManager.getPackageInfo(context.packageName, 0).versionName ?: "dev"
        } catch (_: Exception) {
            "dev"
        }
    }
}

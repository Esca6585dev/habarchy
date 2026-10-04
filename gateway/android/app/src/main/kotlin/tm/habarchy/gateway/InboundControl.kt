package tm.habarchy.gateway

import android.Manifest
import android.content.ComponentName
import android.content.Context
import android.content.pm.PackageManager
import android.net.Uri
import android.provider.Telephony
import androidx.core.content.ContextCompat
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale
import java.util.TimeZone

/**
 * Turns inbound SMS forwarding on or off at the OS level. The receiver
 * component is enabled only while the user switch AND the server's
 * inbound_enabled are on and RECEIVE_SMS is granted; otherwise Android never
 * wakes the app for incoming SMS. When on, the inbox is also read (READ_SMS)
 * for SMS that arrived while the app was not running.
 */
object InboundControl {
    fun hasReceivePermission(context: Context): Boolean =
        ContextCompat.checkSelfPermission(context, Manifest.permission.RECEIVE_SMS) == PackageManager.PERMISSION_GRANTED

    fun hasReadPermission(context: Context): Boolean =
        ContextCompat.checkSelfPermission(context, Manifest.permission.READ_SMS) == PackageManager.PERMISSION_GRANTED

    fun effective(context: Context, store: GatewayStore = GatewayStore(context)): Boolean =
        store.forwardInbound && store.serverInboundEnabled && hasReceivePermission(context)

    /** Applies the current switches to the receiver component. */
    fun apply(context: Context) {
        val store = GatewayStore(context)
        val on = effective(context, store)
        val component = ComponentName(context, SmsInboundReceiver::class.java)
        val state = if (on) PackageManager.COMPONENT_ENABLED_STATE_ENABLED else PackageManager.COMPONENT_ENABLED_STATE_DISABLED
        if (context.packageManager.getComponentEnabledSetting(component) != state) {
            context.packageManager.setComponentEnabledSetting(component, state, PackageManager.DONT_KILL_APP)
            store.log(if (on) "inbound forwarding ON (receiver enabled)" else "inbound forwarding OFF (receiver disabled)")
        }
        store.inboundActive = on
    }

    /** Forwards inbox SMS newer than the last sync point (catch-up after the app was killed). */
    fun syncInbox(context: Context) {
        val store = GatewayStore(context)
        if (!effective(context, store) || !hasReadPermission(context)) return
        val since = store.lastInboundSyncAt
        if (since == 0L) {
            store.lastInboundSyncAt = System.currentTimeMillis()
            return
        }
        val fmt = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'", Locale.US).apply { timeZone = TimeZone.getTimeZone("UTC") }
        var newest = since
        try {
            context.contentResolver.query(
                Uri.parse("content://sms/inbox"),
                arrayOf(Telephony.Sms.ADDRESS, Telephony.Sms.BODY, Telephony.Sms.DATE),
                "${Telephony.Sms.DATE} > ?", arrayOf(since.toString()), "${Telephony.Sms.DATE} ASC",
            )?.use { c ->
                while (c.moveToNext()) {
                    val from = c.getString(0) ?: continue
                    val body = c.getString(1) ?: ""
                    val date = c.getLong(2)
                    Reporter.inbound(context, from, body, fmt.format(Date(date)), date)
                    if (date > newest) newest = date
                }
            }
        } catch (e: Exception) {
            store.log("! inbox read failed: ${e.message}")
        }
        store.lastInboundSyncAt = newest
    }
}

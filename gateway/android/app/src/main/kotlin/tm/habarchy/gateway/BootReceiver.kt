package tm.habarchy.gateway

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent

/** Restarts the gateway after a reboot when it was enabled. */
class BootReceiver : BroadcastReceiver() {
    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action == Intent.ACTION_BOOT_COMPLETED || intent.action == Intent.ACTION_MY_PACKAGE_REPLACED) {
            if (GatewayStore(context).enabled) GatewayService.start(context)
        }
    }
}

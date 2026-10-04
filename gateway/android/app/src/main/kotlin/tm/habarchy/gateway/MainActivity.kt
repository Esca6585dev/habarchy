package tm.habarchy.gateway

import android.Manifest
import android.content.Intent
import android.content.pm.PackageManager
import android.net.Uri
import android.os.Build
import android.os.PowerManager
import android.provider.Settings
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel

/** Bridges the Flutter UI to the native gateway service and settings. */
class MainActivity : FlutterActivity() {
    private var pendingPermissionResult: MethodChannel.Result? = null

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)
        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, CHANNEL).setMethodCallHandler { call, result ->
            val store = GatewayStore(this)
            when (call.method) {
                "getConfig" -> result.success(
                    mapOf("url" to store.url, "key" to store.key, "enabled" to store.enabled, "forwardInbound" to store.forwardInbound)
                )
                "saveConfig" -> {
                    store.url = (call.argument<String>("url") ?: "").trim().trimEnd('/')
                    store.key = (call.argument<String>("key") ?: "").trim()
                    store.forwardInbound = call.argument<Boolean>("forwardInbound") ?: false
                    result.success(true)
                }
                "start" -> {
                    if (!hasPermissions()) {
                        result.error("permission", "SEND_SMS permission missing", null)
                    } else {
                        store.enabled = true
                        GatewayService.start(this)
                        result.success(true)
                    }
                }
                "stop" -> {
                    store.enabled = false
                    GatewayService.stop(this)
                    result.success(true)
                }
                "getStatus" -> result.success(store.statusMap(GatewayService.isRunning))
                "getLog" -> result.success(store.log())
                "clearLog" -> { store.clearLog(); result.success(true) }
                "hasPermissions" -> result.success(hasPermissions())
                "requestPermissions" -> requestPermissions(result)
                "requestIgnoreBatteryOptimizations" -> { requestIgnoreBattery(); result.success(true) }
                "isIgnoringBatteryOptimizations" -> result.success(isIgnoringBattery())
                "testSms" -> {
                    val to = call.argument<String>("to") ?: ""
                    val text = call.argument<String>("text") ?: "Habarchy gateway test"
                    try {
                        SmsSender(this).send(to, text, -1, "test-" + System.currentTimeMillis())
                        result.success(true)
                    } catch (e: Exception) {
                        result.error("sms", e.message, null)
                    }
                }
                else -> result.notImplemented()
            }
        }
    }

    private fun requiredPermissions(): List<String> {
        val list = mutableListOf(Manifest.permission.SEND_SMS, Manifest.permission.READ_PHONE_STATE, Manifest.permission.RECEIVE_SMS)
        if (Build.VERSION.SDK_INT >= 33) list.add(Manifest.permission.POST_NOTIFICATIONS)
        return list
    }

    private fun hasPermissions(): Boolean =
        ContextCompat.checkSelfPermission(this, Manifest.permission.SEND_SMS) == PackageManager.PERMISSION_GRANTED

    private fun requestPermissions(result: MethodChannel.Result) {
        if (hasPermissions() && (Build.VERSION.SDK_INT < 33 ||
                    ContextCompat.checkSelfPermission(this, Manifest.permission.POST_NOTIFICATIONS) == PackageManager.PERMISSION_GRANTED)
        ) {
            result.success(true)
            return
        }
        pendingPermissionResult = result
        ActivityCompat.requestPermissions(this, requiredPermissions().toTypedArray(), REQ_PERMISSIONS)
    }

    override fun onRequestPermissionsResult(requestCode: Int, permissions: Array<out String>, grantResults: IntArray) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == REQ_PERMISSIONS) {
            pendingPermissionResult?.success(hasPermissions())
            pendingPermissionResult = null
        }
    }

    private fun isIgnoringBattery(): Boolean {
        val pm = getSystemService(POWER_SERVICE) as PowerManager
        return pm.isIgnoringBatteryOptimizations(packageName)
    }

    @Suppress("BatteryLife")
    private fun requestIgnoreBattery() {
        if (isIgnoringBattery()) return
        try {
            startActivity(Intent(Settings.ACTION_REQUEST_IGNORE_BATTERY_OPTIMIZATIONS, Uri.parse("package:$packageName")))
        } catch (_: Exception) {
            startActivity(Intent(Settings.ACTION_IGNORE_BATTERY_OPTIMIZATION_SETTINGS))
        }
    }

    companion object {
        const val CHANNEL = "habarchy/gateway"
        private const val REQ_PERMISSIONS = 4101
    }
}

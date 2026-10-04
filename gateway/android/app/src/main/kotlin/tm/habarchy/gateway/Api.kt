package tm.habarchy.gateway

import org.json.JSONArray
import org.json.JSONObject
import java.io.IOException
import java.net.HttpURLConnection
import java.net.URL

/** Minimal JSON client for /api/gateway/v1 (no third-party dependencies). */
class Api(private val baseUrl: String, private val key: String, private val appVersion: String) {

    class HttpError(val code: Int, message: String) : IOException("HTTP $code: $message")

    private fun open(path: String, method: String, readTimeoutMs: Int): HttpURLConnection {
        val conn = URL(baseUrl + path).openConnection() as HttpURLConnection
        conn.requestMethod = method
        conn.connectTimeout = 15_000
        conn.readTimeout = readTimeoutMs
        conn.setRequestProperty("X-Gateway-Key", key)
        conn.setRequestProperty("Accept", "application/json")
        conn.setRequestProperty("User-Agent", "HabarchyGateway/$appVersion (Android)")
        return conn
    }

    private fun exec(conn: HttpURLConnection, body: JSONObject?): JSONObject {
        if (body != null) {
            conn.doOutput = true
            conn.setRequestProperty("Content-Type", "application/json")
            conn.outputStream.use { it.write(body.toString().toByteArray()) }
        }
        val code = conn.responseCode
        val stream = if (code in 200..299) conn.inputStream else conn.errorStream
        val text = stream?.bufferedReader()?.use { it.readText() } ?: ""
        conn.disconnect()
        if (code !in 200..299) {
            val msg = try { JSONObject(text).optJSONObject("error")?.optString("message") ?: text } catch (_: Exception) { text }
            throw HttpError(code, msg.take(200))
        }
        return if (text.isBlank()) JSONObject() else JSONObject(text)
    }

    fun me(): JSONObject = exec(open("/api/gateway/v1/me", "GET", 20_000), null).getJSONObject("data")

    /** Long-polls the outbox; returns leased items (possibly empty). */
    fun lease(waitSec: Int, limit: Int): JSONArray =
        exec(open("/api/gateway/v1/outbox?wait=$waitSec&limit=$limit", "GET", (waitSec + 20) * 1000), null).getJSONArray("data")

    fun report(outboxId: String, status: String, errorCode: String?, errorMessage: String?, parts: Int) {
        val body = JSONObject().put("status", status).put("parts", parts)
        if (errorCode != null) body.put("error_code", errorCode)
        if (errorMessage != null) body.put("error_message", errorMessage.take(500))
        exec(open("/api/gateway/v1/outbox/$outboxId/result", "POST", 20_000), body)
    }

    fun heartbeat(info: JSONObject): JSONObject = exec(open("/api/gateway/v1/heartbeat", "POST", 20_000), info).getJSONObject("data")

    fun inbound(from: String, text: String, receivedAtIso: String) {
        exec(open("/api/gateway/v1/inbound", "POST", 20_000), JSONObject().put("from", from).put("text", text).put("received_at", receivedAtIso))
    }
}

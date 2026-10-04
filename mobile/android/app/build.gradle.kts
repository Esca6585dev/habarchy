import java.util.Properties
import java.io.FileInputStream

plugins {
    id("com.android.application")
    // The Flutter Gradle Plugin must be applied after the Android and Kotlin Gradle plugins.
    id("dev.flutter.flutter-gradle-plugin")
    // Uncomment after adding android/app/google-services.json (see README → Push):
    // id("com.google.gms.google-services")
}

// Release signing: android/key.properties (git-ignored) with
//   storeFile=../keystore/habarchy-release.jks
//   storePassword=...
//   keyAlias=habarchy
//   keyPassword=...
val keystoreProperties = Properties()
val keystorePropertiesFile = rootProject.file("key.properties")
if (keystorePropertiesFile.exists()) {
    keystoreProperties.load(FileInputStream(keystorePropertiesFile))
}

android {
    namespace = "tm.habarchy.admin"
    compileSdk = flutter.compileSdkVersion
    ndkVersion = flutter.ndkVersion

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    defaultConfig {
        applicationId = "tm.habarchy.admin"
        minSdk = maxOf(23, flutter.minSdkVersion) // flutter_secure_storage / firebase_messaging
        targetSdk = flutter.targetSdkVersion
        versionCode = flutter.versionCode
        versionName = flutter.versionName
        resValue("string", "app_name", "Habarçy")
        manifestPlaceholders["usesCleartextTraffic"] = "false"
    }

    flavorDimensions += "env"
    productFlavors {
        create("dev") {
            dimension = "env"
            applicationIdSuffix = ".dev"
            versionNameSuffix = "-dev"
            resValue("string", "app_name", "Habarçy Dev")
            // Emulator / LAN API over plain http during development.
            manifestPlaceholders["usesCleartextTraffic"] = "true"
        }
        create("prod") {
            dimension = "env"
        }
    }

    signingConfigs {
        if (keystorePropertiesFile.exists()) {
            create("release") {
                keyAlias = keystoreProperties["keyAlias"] as String
                keyPassword = keystoreProperties["keyPassword"] as String
                storeFile = file(keystoreProperties["storeFile"] as String)
                storePassword = keystoreProperties["storePassword"] as String
            }
        }
    }

    buildTypes {
        release {
            // Falls back to the debug key when key.properties is missing so
            // `flutter run --release` keeps working on a fresh checkout.
            signingConfig = if (keystorePropertiesFile.exists()) signingConfigs.getByName("release") else signingConfigs.getByName("debug")
            isMinifyEnabled = true
            isShrinkResources = true
            proguardFiles(getDefaultProguardFile("proguard-android-optimize.txt"), "proguard-rules.pro")
        }
    }
}

kotlin {
    compilerOptions {
        jvmTarget = org.jetbrains.kotlin.gradle.dsl.JvmTarget.JVM_17
    }
}

flutter {
    source = "../.."
}

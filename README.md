# SyncDrop

A privacy-first, self-hosted file backup engine designed to automatically sync media and files from your mobile device to personal hardware.

SyncDrop prioritizes high-speed, local-network (LAN) transfers via mDNS discovery when connected to your home Wi-Fi. It provides secure off-site sync using WireGuard and DuckDNS—keeping your private data strictly under your control with zero public storage ports exposed.

---

## Features

* **Zero Cloud Dependency:** Run your personal backup server on your own hardware (Windows, Linux, or Raspberry Pi).
* **Local mDNS Discovery:** Seamless zero-tap pairing when your phone joins your home Wi-Fi network.
* **LAN-Speed Sync:** Transfers media directly over local intranet bandwidth without consuming external data.
* **Encrypted Remote Backups:** Sync on the go over cellular networks via WireGuard tunnels and dynamic DNS (DuckDNS).
* **Zero Open Storage Ports:** No public HTTP ports exposed to the WAN; inbound traffic routes strictly through a WireGuard VPN tunnel.
* **Deduplication:** Checksum-based hashing (SHA-256) avoids redundant transfers of existing files.

---

## Architecture Overview

```text
+-------------------------------------------------------------+
|                        SyncDrop App                         |
+-------------------------------------------------------------+
               |                               |
       (Home Wi-Fi: mDNS)              (Remote / Cellular)
               |                               |
               v                               v
    +--------------------+            +------------------+
    | Direct Local LAN   |            | WireGuard VPN    |
    | (High-Speed Sync)  |            | (DuckDNS Target) |
    +--------------------+            +------------------+
               |                               |
               +---------------+---------------+
                               |
                               v
               +-------------------------------+
               |        SyncDrop Server        |
               |        (Go / Windows / Pi)    |
               +-------------------------------+
                               |
                               v
               +-------------------------------+
               |    Local Storage (SSD/Disk)   |
               +-------------------------------+
```

## Project Structure

This repository is organized as a monorepo containing both the server daemon and the client application:

```text
syncdrop/
├── server/             # Go backend, mDNS responder, file engine
├── client/             # Flutter mobile application (Android / iOS)
├── docs/               # Architecture diagrams, PRD, network guides
├── .github/workflows/  # CI/CD pipelines
├── LICENSE
└── README.md
```

## Getting Started

Prerequisites
- Go: `1.22+`
- Flutter SDK: `3.19+`
- Network: Local router with multicast/mDNS enabled (standard on most home routers).

### 1. Server Setup (Local Prototype)

#### 1.1. Clone the repository:
```bash
git clone [https://github.com/your-username/syncdrop.git](https://github.com/your-username/syncdrop.git)
cd syncdrop/server
```

#### 1.2. Configure environment variables:
```bash
cp .env.example .env
```
Edit `.env` to specify your target storage directory, port, and security tokens.

#### 1.3. Run the Go server:
```bash
go run cmd/server/main.go
```
The server starts listening on the configured port and broadcasts an mDNS service on your local network.

### 2. Client Setup (Flutter)

#### 2.1. Navigate to the client directory:
```base
cd ../client
```

#### 2.2. Fetch Flutter dependencies:
```base
flutter pub get
```

#### 2.3. Launch on a connected device:
```base
flutter run
```

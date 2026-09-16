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

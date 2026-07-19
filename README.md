<div align="center">

# GoDir

**Fast, beautiful local file sharing over your WiFi network.**

Share any folder on your local network with a single command. Browse, navigate, and download files from any device — phone, tablet, laptop — through a sleek dark-themed web interface.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-333?style=for-the-badge&logo=apple&logoColor=white)](https://github.com)

</div>

---

## Features

- **Zero dependencies** — single compiled binary, no runtime required
- **Auto IP detection** — finds your WiFi IP (192.x.x.x) automatically
- **Dark web UI** — responsive, animated, type-aware file icons
- **Fast** — pure Go HTTP server, handles thousands of concurrent connections
- **Secure by default** — only serves the directory you specify, no access to parent paths

## Quick Start

### Build from source

```bash
git clone https://github.com/n11kol11c/GoDir.git
cd GoDir
go build -o GoDir .
```

### Run

```bash
# Share the current directory
./GoDir

# Share a specific folder
./GoDir -dir ~/Documents

# Custom port
./GoDir -port 3000

# Share a folder on a custom port
./GoDir -dir /path/to/files -port 5000
```

### Install with `go install`

```bash
go install github.com/n11kol11c/GoDir@latest
```

## Usage

```
GoDir [flags]

Flags:
  -dir string    Directory to share (default ".")
  -port int      Port to serve on (default 8080)
```

When GoDir starts, it prints a styled banner with the URL to open:

```
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║ 				██████╗  ██████╗ ██████╗ ██╗██████╗                  ║ 
║ 				██╔════╝ ██╔═══██╗██╔══██╗██║██╔══██╗		             ║ 
║ 				██║  ███╗██║   ██║██║  ██║██║██████╔╝		             ║ 
║ 				██║   ██║██║   ██║██║  ██║██║██╔══██╗		             ║ 
║ 				╚██████╔╝╚██████╔╝██████╔╝██║██║  ██║		             ║ 
║ 				╚═════╝  ╚═════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝		             ║ 
║                                      						             ║ 
║               Share files on your local network              ║
╚══════════════════════════════════════════════════════════════╝

  📂 Serving:  /Users/you/Documents
  🌐 Network:  http://*****:8080
  💻 Local:    http://localhost:8080
  ⏱  Started:  2025-07-19 21:11:52
  🛡  OS:       *

  ✔ Server is live! Open the link above in any browser on your network.
```

Open the **Network URL** on any device connected to your WiFi to browse and download files.

## Web Interface

The browser UI provides:

| Feature | Description |
|---|---|
| **Breadcrumbs** | Click any parent folder in the path to navigate up |
| **File icons** | Emoji icons matched by file extension (images, videos, archives, code, etc.) |
| **File sizes** | Human-readable sizes (KB, MB, GB) |
| **Mod dates** | Last modified timestamp for each entry |
| **Download** | Click any file to download it directly |
| **Mobile** | Fully responsive layout for phones and tablets |
| **Animations** | Subtle fade-in rows and glow effects |

## Examples

```bash
# Share your Downloads folder for a quick transfer to your phone
./GoDir -dir ~/Downloads

# Run a temporary server on a shared project folder for your team
./GoDir -dir ./build -port 9090

# Share your home directory (use on trusted networks only)
./GoDir -dir ~
```

## Platform Notes

| OS | Default port | Notes |
|---|---|---|
| macOS | 8080 | Works out of the box. May prompt for network permissions on first run. |
| Linux | 8080 | Works out of the box. Use ports > 1024 to avoid root. |
| Windows | 8080 | Works out of the box. Firewall may prompt — allow access for your network. |

## Security Considerations

- GoDir only serves files under the directory you specify — it cannot escape to parent directories.
- Hidden files (dotfiles) are hidden from the web UI but still accessible via direct URL.
- Do not expose GoDir to the public internet. It has no authentication.
- Use it on **trusted local networks** only.

## License

MIT License. See [LICENSE](LICENSE) for details.

---

<div align="center">

Built with Go. No frameworks. No bloat. Just works.

</div>

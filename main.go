package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	version = "1.0.0"
	banner  = `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║ 				██████╗  ██████╗ ██████╗ ██╗██████╗            ║ 
║ 				██╔════╝ ██╔═══██╗██╔══██╗██║██╔══██╗		   ║ 
║ 				██║  ███╗██║   ██║██║  ██║██║██████╔╝		   ║ 
║ 				██║   ██║██║   ██║██║  ██║██║██╔══██╗		   ║ 
║ 				╚██████╔╝╚██████╔╝██████╔╝██║██║  ██║		   ║ 
║ 				╚═════╝  ╚═════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝		   ║ 
║                                      						   ║ 
║               Share files on your local network              ║
╚══════════════════════════════════════════════════════════════╝

	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorPurple = "\033[35m"
)

func main() {
	port := flag.Int("port", 8080, "Port to serve on")
	dir := flag.String("dir", ".", "Directory to share")
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("%s✖ Error resolving directory:%s %v\n", colorRed, colorReset, err)
	}

	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		log.Fatalf("%s✖ Directory does not exist:%s %s\n", colorRed, colorReset, absDir)
	}

	ip := getLocalIP()

	printBanner()

	fmt.Printf("  %s📂 Serving:%s  %s%s%s\n", colorCyan, colorReset, colorBold, absDir, colorReset)
	fmt.Printf("  %s🌐 Network:%s  %shttp://%s:%d%s\n", colorCyan, colorReset, colorGreen, ip, *port, colorReset)
	fmt.Printf("  %s💻 Local:%s    %shttp://localhost:%d%s\n", colorCyan, colorReset, colorGreen, *port, colorReset)
	fmt.Printf("  %s⏱  Started:%s  %s%s%s\n", colorCyan, colorReset, colorPurple, time.Now().Format("2006-01-02 15:04:05"), colorReset)
	fmt.Printf("  %s🛡  OS:%s       %s%s%s\n", colorCyan, colorReset, colorDim, runtime.GOOS, colorReset)
	fmt.Println()
	fmt.Printf("  %s%sPress Ctrl+C to stop serving%s\n", colorDim, colorYellow, colorReset)
	fmt.Println()

	handler := newFileServer(absDir)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Printf("\n\n  %s🛑 Shutting down GoDir...%s\n", colorRed, colorReset)
		server.Close()
	}()

	fmt.Printf("  %s✔ Server is live! Open the link above in any browser on your network.%s\n\n", colorGreen, colorReset)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("%s✖ Server error:%s %v\n", colorRed, colorReset, err)
	}

	fmt.Printf("  %s👋 GoDir stopped. Goodbye!%s\n\n", colorGreen, colorReset)
}

func printBanner() {
	lines := strings.Split(banner, "\n")
	cyan := "\033[36m"
	bold := "\033[1m"
	reset := "\033[0m"
	dim := "\033[2m"
	yellow := "\033[33m"

	fmt.Println()
	for i, line := range lines {
		if i < len(lines)-1 {
			fmt.Printf("  %s%s%s%s\n", cyan, bold, line, reset)
		} else {
			fmt.Printf("  %s%s%s%s\n", dim, yellow, line, reset)
		}
	}
	fmt.Println()
	fmt.Printf("  %s%s v%s %s— Fast, beautiful local file sharing%s\n\n", bold, cyan, version, dim, reset)
}

func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "localhost"
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if strings.Contains(iface.Name, "lo") || strings.Contains(iface.Name, "utun") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				octets := strings.Split(ip4.String(), ".")
				if len(octets) == 4 && octets[0] == "192" {
					return ip4.String()
				}
			}
		}
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				return ip4.String()
			}
		}
	}
	return "localhost"
}

type fileInfo struct {
	Name    string
	Size    string
	IsDir   bool
	ModTime string
	Icon    string
	Ext     string
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func getFileIcon(name string, isDir bool) string {
	if isDir {
		return "\U0001F4C1"
	}
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".ico":
		return "\U0001F5BC"
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm":
		return "\U0001F3AC"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma", ".m4a":
		return "\U0001F3B5"
	case ".pdf":
		return "\U0001F4C4"
	case ".doc", ".docx":
		return "\U0001F4DD"
	case ".xls", ".xlsx", ".csv":
		return "\U0001F4CA"
	case ".ppt", ".pptx":
		return "\U0001F4D1"
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".bz2":
		return "\U0001F4E6"
	case ".go", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".rs", ".rb", ".php", ".sh":
		return "\U0001F4BB"
	case ".txt", ".md", ".log", ".ini", ".cfg", ".conf":
		return "\U0001F4DC"
	case ".html", ".css", ".xml", ".json", ".yaml", ".yml":
		return "\U0001F310"
	case ".exe", ".msi", ".dmg", ".app":
		return "\u2699"
	case ".ttf", ".otf", ".woff", ".woff2":
		return "\U0001F524"
	default:
		return "\U0001F4C4"
	}
}

type breadcrumbPart struct {
	Name string
	Path string
}

type pageData struct {
	Title     string
	Path      string
	Bread     []breadcrumbPart
	Files     []fileInfo
	ParentDir string
	IPAddress string
	Version   string
	NumFiles  int
}

func newFileServer(rootDir string) http.Handler {
	mux := http.NewServeMux()

	funcMap := template.FuncMap{
		"fileCount": func(f []fileInfo) int { return len(f) },
	}

	tmpl := template.Must(template.New("dir").Funcs(funcMap).Parse(pageTemplate))

	fileServer := http.FileServer(http.Dir(rootDir))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cleanPath := filepath.Clean(r.URL.Path)
		if cleanPath == "/" || cleanPath == "." {
			cleanPath = "/"
		}

		fullPath := filepath.Join(rootDir, cleanPath)
		info, err := os.Stat(fullPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		entries, err := os.ReadDir(fullPath)
		if err != nil {
			http.Error(w, "Cannot read directory", http.StatusInternalServerError)
			return
		}

		files := make([]fileInfo, 0, len(entries))

		if cleanPath != "/" {
			files = append(files, fileInfo{
				Name:    "..",
				Size:    "-",
				IsDir:   true,
				ModTime: "-",
				Icon:    "\u2B06",
				Ext:     "",
			})
		}

		dirs := make([]fileInfo, 0)
		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			entryInfo, err := entry.Info()
			if err != nil {
				continue
			}
			f := fileInfo{
				Name:    name,
				Size:    formatSize(entryInfo.Size()),
				IsDir:   entry.IsDir(),
				ModTime: entryInfo.ModTime().Format("2006-01-02 15:04"),
				Icon:    getFileIcon(name, entry.IsDir()),
				Ext:     strings.ToLower(filepath.Ext(name)),
			}
			if entry.IsDir() {
				dirs = append(dirs, f)
			} else {
				files = append(files, f)
			}
			if len(dirs)+len(files)-1 > 5000 {
				break
			}
		}
		files = append(dirs, files...)

		trimmedPath := strings.Trim(cleanPath, "/")
		breadParts := make([]breadcrumbPart, 0)
		if trimmedPath != "" {
			parts := strings.Split(trimmedPath, "/")
			accumPath := ""
			for _, part := range parts {
				accumPath += "/" + part + "/"
				breadParts = append(breadParts, breadcrumbPart{
					Name: part,
					Path: accumPath,
				})
			}
		}

		ip := getLocalIP()

		data := pageData{
			Title:     "GoDir \u2014 " + cleanPath,
			Path:      cleanPath,
			Bread:     breadParts,
			Files:     files,
			ParentDir: filepath.Dir(strings.TrimSuffix(cleanPath, "/")),
			IPAddress: ip,
			Version:   version,
			NumFiles:  len(files),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
	})

	return logMiddleware(mux)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		duration := time.Since(start)

		status := "\033[32m"
		if wrapped.statusCode >= 400 {
			status = "\033[33m"
		}
		if wrapped.statusCode >= 500 {
			status = "\033[31m"
		}

		fmt.Printf("  %s%s%d%s %s %s%s%s\n",
			"\033[2m", status, wrapped.statusCode, "\033[0m",
			r.URL.Path,
			"\033[2m", duration.Round(time.Millisecond), "\033[0m")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

const pageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Inter:wght@400;500;600;700&display=swap');

        :root {
            --bg-primary: #0a0a0f;
            --bg-secondary: #12121a;
            --bg-tertiary: #1a1a2e;
            --bg-hover: #1e1e36;
            --border: #2a2a40;
            --accent: #6c5ce7;
            --accent-light: #a29bfe;
            --accent-glow: rgba(108, 92, 231, 0.3);
            --text-primary: #e2e8f0;
            --text-secondary: #94a3b8;
            --text-muted: #64748b;
            --green: #00b894;
            --cyan: #00cec9;
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            min-height: 100vh;
            line-height: 1.6;
        }

        .bg-grid {
            position: fixed; top: 0; left: 0; width: 100%; height: 100%;
            background-image:
                linear-gradient(rgba(108, 92, 231, 0.03) 1px, transparent 1px),
                linear-gradient(90deg, rgba(108, 92, 231, 0.03) 1px, transparent 1px);
            background-size: 50px 50px;
            pointer-events: none; z-index: 0;
        }

        .bg-glow {
            position: fixed; top: -50%; left: -50%; width: 200%; height: 200%;
            background: radial-gradient(circle at 50% 50%, rgba(108, 92, 231, 0.06) 0%, transparent 50%);
            pointer-events: none; z-index: 0;
            animation: glowPulse 8s ease-in-out infinite alternate;
        }

        @keyframes glowPulse { 0% { opacity: 0.5; } 100% { opacity: 1; } }

        .container {
            position: relative; z-index: 1;
            max-width: 1100px; margin: 0 auto; padding: 24px 20px;
        }

        header {
            display: flex; align-items: center; justify-content: space-between;
            margin-bottom: 32px; padding: 20px 28px;
            background: var(--bg-secondary); border: 1px solid var(--border);
            border-radius: 16px; backdrop-filter: blur(20px);
        }

        .logo { display: flex; align-items: center; gap: 14px; }

        .logo-icon {
            width: 44px; height: 44px;
            background: linear-gradient(135deg, var(--accent), var(--accent-light));
            border-radius: 12px; display: flex; align-items: center; justify-content: center;
            font-size: 22px; box-shadow: 0 0 20px var(--accent-glow);
        }

        .logo h1 {
            font-family: 'JetBrains Mono', monospace;
            font-size: 1.5rem; font-weight: 700;
            background: linear-gradient(135deg, var(--accent-light), var(--cyan));
            -webkit-background-clip: text; -webkit-text-fill-color: transparent;
            background-clip: text; letter-spacing: -0.5px;
        }

        .logo .version {
            font-size: 0.7rem; color: var(--text-muted);
            font-family: 'JetBrains Mono', monospace;
        }

        .header-info { display: flex; align-items: center; gap: 16px; }

        .status-badge {
            display: flex; align-items: center; gap: 8px;
            padding: 6px 14px; background: rgba(0, 184, 148, 0.1);
            border: 1px solid rgba(0, 184, 148, 0.2); border-radius: 100px;
            font-size: 0.75rem; font-weight: 500; color: var(--green);
        }

        .status-dot {
            width: 7px; height: 7px; background: var(--green);
            border-radius: 50%; animation: pulse 2s ease-in-out infinite;
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; box-shadow: 0 0 0 0 rgba(0, 184, 148, 0.4); }
            50% { opacity: 0.8; box-shadow: 0 0 0 6px rgba(0, 184, 148, 0); }
        }

        .breadcrumb {
            display: flex; align-items: center; gap: 4px; flex-wrap: wrap;
            padding: 16px 28px; margin-bottom: 24px;
            background: var(--bg-secondary); border: 1px solid var(--border);
            border-radius: 12px; font-family: 'JetBrains Mono', monospace; font-size: 0.82rem;
        }

        .breadcrumb a {
            color: var(--accent-light); text-decoration: none;
            padding: 3px 8px; border-radius: 6px; transition: all 0.2s;
        }
        .breadcrumb a:hover { background: var(--bg-tertiary); color: var(--cyan); }
        .breadcrumb .sep { color: var(--text-muted); font-size: 0.7rem; }
        .breadcrumb .current { color: var(--text-primary); font-weight: 600; }

        .file-table {
            background: var(--bg-secondary); border: 1px solid var(--border);
            border-radius: 16px; overflow: hidden; backdrop-filter: blur(20px);
        }

        .file-table table { width: 100%; border-collapse: collapse; }

        .file-table thead th {
            padding: 14px 24px; text-align: left;
            font-size: 0.72rem; font-weight: 600; text-transform: uppercase;
            letter-spacing: 1px; color: var(--text-muted);
            background: var(--bg-tertiary); border-bottom: 1px solid var(--border);
        }

        .file-table tbody tr {
            transition: all 0.15s ease; border-bottom: 1px solid rgba(42, 42, 64, 0.4);
        }
        .file-table tbody tr:last-child { border-bottom: none; }
        .file-table tbody tr:hover { background: var(--bg-hover); }
        .file-table td { padding: 12px 24px; font-size: 0.88rem; }

        .file-name { display: flex; align-items: center; gap: 12px; }
        .file-name a {
            color: var(--text-primary); text-decoration: none;
            font-weight: 500; transition: color 0.2s;
        }
        .file-name a:hover { color: var(--accent-light); }
        .file-name .icon { font-size: 1.2rem; width: 28px; text-align: center; flex-shrink: 0; }
        .file-size { font-family: 'JetBrains Mono', monospace; font-size: 0.8rem; color: var(--text-secondary); }
        .file-date { color: var(--text-muted); font-size: 0.82rem; }
        .file-ext {
            display: inline-block; padding: 2px 8px; background: var(--bg-tertiary);
            border-radius: 6px; font-family: 'JetBrains Mono', monospace;
            font-size: 0.7rem; color: var(--text-muted); border: 1px solid var(--border);
        }

        .empty-state { text-align: center; padding: 60px 20px; color: var(--text-muted); }
        .empty-state .empty-icon { font-size: 3rem; margin-bottom: 16px; opacity: 0.5; }
        .empty-state h3 { font-size: 1.1rem; color: var(--text-secondary); margin-bottom: 8px; }

        footer {
            margin-top: 32px; padding: 16px 28px;
            background: var(--bg-secondary); border: 1px solid var(--border);
            border-radius: 12px; display: flex; align-items: center;
            justify-content: space-between; font-size: 0.75rem; color: var(--text-muted);
        }
        footer a { color: var(--accent-light); text-decoration: none; }
        footer a:hover { color: var(--cyan); }

        .stats { display: flex; gap: 20px; align-items: center; }
        .stat { display: flex; align-items: center; gap: 6px; }

        @media (max-width: 768px) {
            header { flex-direction: column; gap: 16px; text-align: center; }
            .header-info { justify-content: center; }
            .file-table thead { display: none; }
            .file-table tbody tr { display: flex; flex-direction: column; padding: 16px 20px; gap: 6px; }
            .file-table td { padding: 0; }
            .file-table td:nth-child(3), .file-table td:nth-child(4) { display: none; }
            footer { flex-direction: column; gap: 12px; text-align: center; }
        }

        tr { animation: fadeIn 0.2s ease backwards; }
        tr:nth-child(1) { animation-delay: 0.02s; }
        tr:nth-child(2) { animation-delay: 0.04s; }
        tr:nth-child(3) { animation-delay: 0.06s; }
        tr:nth-child(4) { animation-delay: 0.08s; }
        tr:nth-child(5) { animation-delay: 0.10s; }
        tr:nth-child(6) { animation-delay: 0.12s; }
        tr:nth-child(7) { animation-delay: 0.14s; }
        tr:nth-child(8) { animation-delay: 0.16s; }
        tr:nth-child(9) { animation-delay: 0.18s; }
        tr:nth-child(10) { animation-delay: 0.20s; }

        @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

        .back-row td { border-bottom: 1px solid var(--border) !important; }
    </style>
</head>
<body>
    <div class="bg-grid"></div>
    <div class="bg-glow"></div>

    <div class="container">
        <header>
            <div class="logo">
                <div class="logo-icon">&#128193;</div>
                <div>
                    <h1>GoDir</h1>
                    <span class="version">v{{.Version}}</span>
                </div>
            </div>
            <div class="header-info">
                <div class="status-badge">
                    <span class="status-dot"></span>
                    Serving on {{.IPAddress}}
                </div>
            </div>
        </header>

        <nav class="breadcrumb">
            <a href="/">&#127968; root</a>
            {{range .Bread}}
                <span class="sep">&#8250;</span>
                <a href="{{.Path}}">{{.Name}}</a>
            {{end}}
        </nav>

        <div class="file-table">
            {{if .Files}}
            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Size</th>
                        <th>Modified</th>
                        <th>Type</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Files}}
                    <tr{{if eq .Name ".."}} class="back-row"{{end}}>
                        <td>
                            <div class="file-name">
                                <span class="icon">{{.Icon}}</span>
                                {{if eq .Name ".."}}
                                    <a href="{{$.ParentDir}}">..</a>
                                {{else if .IsDir}}
                                    <a href="{{$.Path}}{{.Name}}/">{{.Name}}/</a>
                                {{else}}
                                    <a href="{{$.Path}}{{.Name}}" download>{{.Name}}</a>
                                {{end}}
                            </div>
                        </td>
                        <td class="file-size">{{.Size}}</td>
                        <td class="file-date">{{.ModTime}}</td>
                        <td>
                            {{if .Ext}}
                                <span class="file-ext">{{.Ext}}</span>
                            {{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <div class="empty-state">
                <div class="empty-icon">&#128381;</div>
                <h3>This folder is empty</h3>
                <p>No files or folders to display</p>
            </div>
            {{end}}
        </div>

        <footer>
            <div class="stats">
                <div class="stat">
                    <span>&#128193;</span>
                    {{.NumFiles}} items
                </div>
                <div class="stat">
                    <span>&#127760;</span>
                    {{.IPAddress}}
                </div>
            </div>
            <div>
                Powered by <a href="#">GoDir</a> &mdash; Fast local file sharing
            </div>
        </footer>
    </div>
</body>
</html>`

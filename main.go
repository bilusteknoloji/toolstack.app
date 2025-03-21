package main

import (
	"bytes"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

const staticSiteFolder = "./site"

type templateDataIP struct {
	ClientIP          string
	TextSizeIPAddress string
}

func isDevelopment() bool {
	return os.Getenv("GOLANG_ENV") == "development"
}

type liveReloadServer struct {
	clients  map[*websocket.Conn]bool
	upgrader websocket.Upgrader
	mu       sync.Mutex
}

func (s *liveReloadServer) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[%s]: websocket upgrade: %s", color.HiBlackString("ERR"), color.RedString(err.Error()))
		return
	}
	defer conn.Close()

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			s.mu.Lock()
			delete(s.clients, conn)
			s.mu.Unlock()
			break
		}
	}
}

func (s *liveReloadServer) notifyClients() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
			log.Printf(
				"[%s]: websocket: %s, client: %v",
				color.HiBlackString("ERR"),
				color.RedString(err.Error()),
				client,
			)
			client.Close()
			delete(s.clients, client)
		}
	}
}

func newLiveReloadServer() *liveReloadServer {
	return &liveReloadServer{
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func getListenAddr() string {
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		return addr
	}

	return ":8000"
}

func watchFiles(dir string, liveReload *liveReloadServer) {
	if !isDevelopment() {
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				log.Printf("[%s]: %s", color.HiBlackString("CHANGE"), color.YellowString(event.Name))
				liveReload.notifyClients()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("[%s]: watch err: %s", color.HiBlackString("ERR"), color.RedString(err.Error()))
		}
	}
}

func getIP(r *http.Request) string {
	isDebug := strings.ToLower(r.URL.Query().Get("debug")) == "true"

	if isDebug {
		for k, v := range r.Header {
			log.Println("[DEBUG][HEADER]", k, v)
		}
	}
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		if isDebug {
			for i, possibleIP := range strings.Split(xff, ",") {
				trimmedIP := strings.TrimSpace(possibleIP)
				log.Println("[DEBUG][X-Forwarded-For]", trimmedIP)
				if i == 0 {
					log.Println("[DEBUG][X-Forwarded-For][FIRST]", trimmedIP)
					xff = trimmedIP
				}
			}
		}
		return xff
	}

	xrip := r.Header.Get("X-Real-IP")
	if xrip != "" {
		return xrip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "Unknown"
	}

	if ip == "::1" {
		return "127.0.0.1"
	}

	host := strings.Split(ip, "%")
	if len(host) > 0 {
		ip = host[0]
	}
	return ip
}

func parseTemplateAndInjectLiveReloadHandler(w http.ResponseWriter, r *http.Request) {
	filePath := staticSiteFolder + r.URL.Path
	if strings.HasSuffix(r.URL.Path, "/") {
		filePath = filePath + "index.html"
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		http.Error(
			w,
			"["+color.HiBlackString("ERR")+"]: template parse error: "+err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	var htmlContent string
	var buf bytes.Buffer

	switch {
	case strings.HasPrefix(r.URL.Path, "/ip/"):
		clientIP := getIP(r)
		acceptHeader := r.Header.Get("Accept")
		userAgent := r.Header.Get("User-Agent")

		if strings.Contains(acceptHeader, "text/plain") || strings.Contains(userAgent, "curl") ||
			strings.Contains(userAgent, "Wget") || strings.Contains(userAgent, "HTTPie") {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte(clientIP + "\n"))

			return
		}

		cssTextSizeIP := "4"
		if len(clientIP) > 18 {
			cssTextSizeIP = "2"
		}
		placeHolderIP := templateDataIP{
			ClientIP:          clientIP,
			TextSizeIPAddress: "text-" + cssTextSizeIP + "xl",
		}
		if err := tmpl.Execute(&buf, placeHolderIP); err != nil {
			http.Error(
				w,
				"["+color.HiBlackString("ERR")+"]: template execute error: "+err.Error(),
				http.StatusInternalServerError,
			)

			return
		}
		htmlContent = buf.String()
	default:
		htmlContentBytes, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "["+color.HiBlackString("ERR")+"]: file not found error: "+err.Error(), http.StatusNotFound)

			return
		}
		htmlContent = string(htmlContentBytes)
	}

	if isDevelopment() {
		log.Printf("[%s]: live reload injected to %s", color.HiBlackString("INJECT"), color.YellowString(filePath))

		htmlContent = strings.Replace(htmlContent, "</body>", `
			<script>
				const socket = new WebSocket("ws://" + window.location.host + "/ws");
				socket.onmessage = () => location.reload();
			</script>
			</body>
		`, 1)
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(htmlContent))
	logRequest(r, http.StatusOK, time.Since(time.Now()))
}

type responseLogger struct {
	http.ResponseWriter
	statusCode int
}

func (rl *responseLogger) WriteHeader(code int) {
	rl.statusCode = code
	rl.ResponseWriter.WriteHeader(code)
}

func logRequest(r *http.Request, status int, duration time.Duration) {
	statusColor := color.HiGreenString("%d", status)
	if status >= 400 {
		statusColor = color.RedString("%d", status)
	}
	log.Printf(
		"[%s]: %s - %s (%s)",
		color.HiBlackString(r.Method),
		color.WhiteString(r.RequestURI),
		statusColor,
		duration,
	)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rl := &responseLogger{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rl, r)

		logRequest(r, rl.statusCode, time.Since(start))
	})
}

func main() {
	addr := getListenAddr()

	var liveReloader *liveReloadServer

	if isDevelopment() {
		liveReloader = newLiveReloadServer()
	}

	fs := loggingMiddleware(http.FileServer(http.Dir("./site")))
	if liveReloader != nil {
		http.HandleFunc("/ws", liveReloader.webSocketHandler)

		go watchFiles("./site", liveReloader)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/healthz") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))

			return
		}

		if strings.HasSuffix(r.URL.Path, "/") || strings.HasSuffix(r.URL.Path, ".html") {
			parseTemplateAndInjectLiveReloadHandler(w, r)

			return
		}

		fs.ServeHTTP(w, r)
	})

	log.Printf("[%s]: Listening at: %s", color.HiBlackString("START"), color.WhiteString(addr))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

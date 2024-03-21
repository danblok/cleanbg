package httpapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/danblok/cleanbg/internal/types"
)

type HTTPServer struct {
	srv    *http.Server
	log    *slog.Logger
	client types.Cleaner
}

func NewHTTPServer(client types.Cleaner, addr string, log *slog.Logger) (*HTTPServer, error) {
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	s := &HTTPServer{srv: srv, log: log, client: client}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("POST /remove", s.handleRemove)

	styles := http.FileServer(http.Dir("./public/"))
	mux.Handle("/styles/", http.StripPrefix("/styles/", styles))
	s.srv.Handler = mux

	return s, nil
}

func (s *HTTPServer) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *HTTPServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("public/index.html")
	if err != nil {
		s.log.Error("http", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to read html page"})
		return
	}

	_, err = w.Write(file)
	if err != nil {
		s.log.Error("http", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to return html page"})
		return
	}
}

func (s *HTTPServer) handleRemove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		s.log.Error("http", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to parse multipart form"})
		return
	}
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		s.log.Error("http", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get file from form"})
		return
	}
	defer file.Close()
	s.log.Debug("http.remove", "size", fileHeader.Size)

	data, err := io.ReadAll(file)
	if err != nil {
		s.log.Error("http", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to read parsed file"})
		return
	}
	image, err := s.client.Clean(
		context.Background(),
		data,
	)
	if err != nil {
		s.log.Error("http", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to remove background"})
		return
	}

	w.Header().Set("Content-type", "application/octet-stream")
	_, err = w.Write(
		[]byte(
			fmt.Sprintf(
				`<img id="image-wo-bg" src="data:image%s;base64,%s"
class="border border-gray-200 rounded-md grid w-full place-items-center border-gray-200 h-[720px] sm:max-h-[720px] dark:border-gray-800 hover:bg-gray-100 dark:hover:bg-gray-50 object-scale-down"
/>`,
				"jpg",
				base64.StdEncoding.EncodeToString(image),
			),
		),
	)
	if err != nil {
		s.log.Error("http", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to return image"})
		return
	}

	s.log.Info("successfully removed background")
}

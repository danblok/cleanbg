package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/danblok/cleanbg/internal/grpcapi"
	"github.com/danblok/cleanbg/internal/log"
)

func main() {
	log := log.New("local")
	log.Warn("logger enabled")

	client, err := grpcapi.Connect("[::]:42069")
	if err != nil {
		log.Error("main", "error", err)
		return
	}

	http.HandleFunc("POST /remove", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Error("http", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to parse multipart form"})
			return
		}
		file, fileHeader, err := r.FormFile("image")
		if err != nil {
			log.Error("http", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get file from form"})
			return
		}
		defer file.Close()
		log.Debug("http.remove", "size", fileHeader.Size)

		data, err := io.ReadAll(file)
		if err != nil {
			log.Error("http", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to read parsed file"})
			return
		}
		image, err := client.Clean(
			context.Background(),
			data,
		)
		if err != nil {
			log.Error("http", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to remove background"})
			return
		}

		w.Header().Set("Content-type", "application/octet-stream")
		_, err = w.Write(
			[]byte(
				fmt.Sprintf(
					"<img id=\"image-wo-bg\" src=\"data:image/%s;base64,%s\" />",
					"jpeg",
					base64.StdEncoding.EncodeToString(image),
				),
			),
		)
		if err != nil {
			log.Error("http", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to return image"})
			return
		}

		log.Info("successfully removed background")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("public/index.html")
		if err != nil {
			log.Error("http", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to read html page"})
			return
		}

		_, err = w.Write(file)
		if err != nil {
			log.Error("http", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to return html page"})
			return
		}
	})

	log.Info("started http server")
	if err := http.ListenAndServe("[::]:8080", nil); err != nil {
		panic(err)
	}
}

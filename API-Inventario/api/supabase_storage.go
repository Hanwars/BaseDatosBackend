package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const supplyImageMaxBytes = 2 << 20 // 2 MiB; keep in sync with the Supabase bucket limit.

var supabaseStorageHTTPClient = &http.Client{Timeout: 60 * time.Second}

type supabaseStorageConfig struct {
	baseURL string
	apiKey  string
	bucket  string
}

func loadSupabaseStorageConfig() (supabaseStorageConfig, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("SUPABASE_URL")), "/")
	apiKey := strings.TrimSpace(os.Getenv("SUPABASE_SECRET_KEY"))
	bucket := strings.TrimSpace(os.Getenv("SUPABASE_STORAGE_BUCKET"))
	if bucket == "" {
		bucket = "emi-images"
	}
	if baseURL == "" {
		return supabaseStorageConfig{}, fmt.Errorf("falta configurar SUPABASE_URL")
	}
	parsed, err := url.ParseRequestURI(baseURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return supabaseStorageConfig{}, fmt.Errorf("SUPABASE_URL no es una URL HTTP/HTTPS válida")
	}
	if apiKey == "" {
		return supabaseStorageConfig{}, fmt.Errorf("falta configurar SUPABASE_SECRET_KEY en el backend")
	}
	if strings.ContainsAny(bucket, "/\\?#") {
		return supabaseStorageConfig{}, fmt.Errorf("SUPABASE_STORAGE_BUCKET contiene caracteres no permitidos")
	}
	return supabaseStorageConfig{baseURL: baseURL, apiKey: apiKey, bucket: bucket}, nil
}

func escapeStorageObjectPath(objectPath string) string {
	parts := strings.Split(objectPath, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
}

// supabaseStorageRequest performs a server-to-server request. The secret API key
// stays in Render's environment and is never returned to the browser.
func supabaseStorageRequest(
	ctx context.Context,
	method string,
	objectPath string,
	body io.Reader,
	contentType string,
) (*http.Response, error) {
	config, err := loadSupabaseStorageConfig()
	if err != nil {
		return nil, err
	}

	// sb_secret_* is an opaque API key, not a JWT. Send it only in
	// the apikey header. It maps to the service_role on Supabase, which
	// can access objects in a private bucket without a user JWT.
	requestURL := config.baseURL + "/storage/v1/object/" + url.PathEscape(config.bucket) + "/" + escapeStorageObjectPath(objectPath)
	request, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("no se pudo preparar la petición a Supabase Storage: %w", err)
	}
	request.Header.Set("apikey", config.apiKey)
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if method == http.MethodPost {
		request.Header.Set("x-upsert", "false")
		request.Header.Set("cache-control", "max-age=3600")
	}

	response, err := supabaseStorageHTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar con Supabase Storage: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		defer response.Body.Close()
		responseBody, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		return nil, fmt.Errorf("Supabase Storage respondió HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}
	return response, nil
}

// supabaseDownloadStorageObject retrieves an object from a private bucket.
// Private buckets cannot be downloaded using only an sb_secret_* API key:
// the Storage API requires either a user JWT or a time-limited signed URL.
// We create the signed URL server-side, then fetch it without exposing the
// Supabase secret key to the browser.
func supabaseDownloadStorageObject(
	ctx context.Context,
	objectPath string,
) (*http.Response, error) {
	config, err := loadSupabaseStorageConfig()
	if err != nil {
		return nil, err
	}

	signURL := config.baseURL + "/storage/v1/object/sign/" +
		url.PathEscape(config.bucket) + "/" + escapeStorageObjectPath(objectPath)
	signBody := strings.NewReader(`{"expiresIn":120}`)
	signRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, signURL, signBody)
	if err != nil {
		return nil, fmt.Errorf("no se pudo preparar la solicitud de URL firmada: %w", err)
	}
	signRequest.Header.Set("apikey", config.apiKey)
	signRequest.Header.Set("Content-Type", "application/json")

	signResponse, err := supabaseStorageHTTPClient.Do(signRequest)
	if err != nil {
		return nil, fmt.Errorf("no se pudo solicitar la URL firmada a Supabase: %w", err)
	}
	defer signResponse.Body.Close()

	if signResponse.StatusCode < http.StatusOK || signResponse.StatusCode >= http.StatusMultipleChoices {
		responseBody, _ := io.ReadAll(io.LimitReader(signResponse.Body, 2048))
		return nil, fmt.Errorf("Supabase no pudo firmar la descarga (HTTP %d): %s", signResponse.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var signed struct {
		SignedURL string `json:"signedURL"`
	}
	if err := json.NewDecoder(io.LimitReader(signResponse.Body, 16*1024)).Decode(&signed); err != nil {
		return nil, fmt.Errorf("Supabase devolvió una respuesta de URL firmada inválida: %w", err)
	}
	if signed.SignedURL == "" {
		return nil, fmt.Errorf("Supabase no devolvió signedURL")
	}

	storageBaseURL := config.baseURL + "/storage/v1"
	var downloadURL string
	if parsed, parseErr := url.Parse(signed.SignedURL); parseErr == nil && parsed.IsAbs() {
		base, _ := url.Parse(config.baseURL)
		if base == nil || parsed.Host != base.Host || parsed.Scheme != base.Scheme {
			return nil, fmt.Errorf("Supabase devolvió una URL firmada de un host inesperado")
		}
		downloadURL = signed.SignedURL
	} else if strings.HasPrefix(signed.SignedURL, "/") {
		downloadURL = storageBaseURL + signed.SignedURL
	} else {
		downloadURL = storageBaseURL + "/" + signed.SignedURL
	}

	downloadRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("no se pudo preparar la descarga firmada: %w", err)
	}
	// A signed URL already carries its short-lived authorization token.
	// Do not add the secret key or an Authorization header to this request.
	downloadResponse, err := supabaseStorageHTTPClient.Do(downloadRequest)
	if err != nil {
		return nil, fmt.Errorf("no se pudo descargar el objeto mediante la URL firmada: %w", err)
	}
	if downloadResponse.StatusCode < http.StatusOK || downloadResponse.StatusCode >= http.StatusMultipleChoices {
		defer downloadResponse.Body.Close()
		responseBody, _ := io.ReadAll(io.LimitReader(downloadResponse.Body, 2048))
		return nil, fmt.Errorf("Supabase no pudo entregar el objeto firmado (HTTP %d): %s", downloadResponse.StatusCode, strings.TrimSpace(string(responseBody)))
	}
	return downloadResponse, nil
}

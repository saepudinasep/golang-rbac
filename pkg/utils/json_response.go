package utils

import (
	"encoding/json"
	"net/http"
)

// ResponseFormat adalah standar struktur JSON REST API industri
type ResponseFormat struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty: jika nil, field tidak akan muncul di JSON
}

// -----------------------------------------------------------------------------
// HELPER UTAMA (BASE FUNCTION)
// -----------------------------------------------------------------------------

// WriteJSON adalah fungsi dasar untuk mengirimkan response JSON
func WriteJSON(w http.ResponseWriter, status int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := ResponseFormat{
		Success: success,
		Message: message,
		Data:    data,
	}

	_ = json.NewEncoder(w).Encode(res)
}

// ReadJSON adalah helper untuk membaca dan men-decode request body
func ReadJSON(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

// -----------------------------------------------------------------------------
// STANDARD INDUSTRY HTTP STATUS HELPERS
// -----------------------------------------------------------------------------

// OK - 200: Request sukses dan mengembalikan data
func OK(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusOK, true, message, data)
}

// Created - 201: Sukses membuat data baru (Resource Created)
func Created(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusCreated, true, message, data)
}

// BadRequest - 400: Error validasi input dari client, format JSON salah, dll.
func BadRequest(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Bad request, validasi gagal"
	}
	WriteJSON(w, http.StatusBadRequest, false, message, nil)
}

// Unauthorized - 401: Client belum login atau token JWT salah/missing
func Unauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Unauthorized, silakan login terlebih dahulu"
	}
	WriteJSON(w, http.StatusUnauthorized, false, message, nil)
}

// Forbidden - 403: Client sudah login, tapi tidak punya hak akses (Role tidak cocok)
func Forbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Forbidden, Anda tidak memiliki akses ke halaman ini"
	}
	WriteJSON(w, http.StatusForbidden, false, message, nil)
}

// NotFound - 404: Data atau Endpoint tidak ditemukan
func NotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Resource tidak ditemukan"
	}
	WriteJSON(w, http.StatusNotFound, false, message, nil)
}

// InternalServerError - 500: Terjadi error di database atau server crash
func InternalServerError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Terjadi kesalahan internal pada server"
	}
	WriteJSON(w, http.StatusInternalServerError, false, message, nil)
}

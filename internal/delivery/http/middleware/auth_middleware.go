package middleware

import (
	"context"
	"fmt"
	"golang-rbac/pkg/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClaimsContextKey contextKey = "userClaims"

// AuthMiddleware menyimpan jwtSecret sebagai dependency, bukan variable global.
// Sebelumnya JWT_SECRET didefinisikan ulang sebagai global var di file ini DAN
// di usecase secara terpisah — kalau salah satu lupa diubah, validasi token jadi
// gagal terus. Sekarang secret hanya didefinisikan sekali di config lalu di-inject.
type AuthMiddleware struct {
	jwtSecret []byte
}

func NewAuthMiddleware(jwtSecret []byte) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

// RoleMiddleware memvalidasi JWT lalu memastikan role user ada di allowedRoles.
func (a *AuthMiddleware) RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteJSON(w, http.StatusUnauthorized, false, "Token tidak ditemukan", nil)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse dan validasi token JWT
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("metode signing tidak terduga: %v", t.Header["alg"])
				}
				return a.jwtSecret, nil
			})

			if err != nil || !token.Valid {
				utils.WriteJSON(w, http.StatusUnauthorized, false, "Token tidak valid atau kedaluwarsa", nil)
				return
			}

			// Mengambil klaim data dari JWT
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.WriteJSON(w, http.StatusUnauthorized, false, "Gagal membaca data token", nil)
				return
			}

			userRole, ok := claims["role_name"].(string)
			if !ok {
				utils.WriteJSON(w, http.StatusForbidden, false, "Akses ditolak: Data role tidak ditemukan", nil)
				return
			}

			// Cek apakah role user ada di dalam daftar allowedRoles
			isAllowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				utils.WriteJSON(w, http.StatusForbidden, false, "Forbidden: Anda tidak memiliki akses untuk role ini", nil)
				return
			}

			// Simpan claims ke context supaya handler berikutnya bisa baca data user
			// (sebelumnya tidak ada cara bagi handler untuk tahu siapa user yang login).
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClaimsFromContext adalah helper agar handler bisa mengambil data JWT claims
// milik user yang sedang login.
func ClaimsFromContext(r *http.Request) (jwt.MapClaims, bool) {
	claims, ok := r.Context().Value(ClaimsContextKey).(jwt.MapClaims)
	return claims, ok
}

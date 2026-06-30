package http

import (
	"golang-rbac/internal/delivery/http/middleware"
	"golang-rbac/internal/domain"
	"golang-rbac/pkg/utils"
	"net/http"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewUserHandler(uc domain.UserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: uc}
}

// Register adalah endpoint baru: sebelumnya UserUsecase.Register sudah ada
// di business logic layer tapi tidak punya HTTP handler/route sama sekali.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, false, "Method tidak diizinkan", nil)
		return
	}

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := utils.ReadJSON(r, &input); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, false, "Format request tidak valid", nil)
		return
	}

	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.UserUsecase.Register(r.Context(), user); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	user.Password = "" // jangan kembalikan hash ke client
	utils.WriteJSON(w, http.StatusCreated, true, "Registrasi berhasil", user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, false, "Method tidak diizinkan", nil)
		return
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Menggunakan helper ReadJSON
	if err := utils.ReadJSON(r, &input); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, false, "Format request tidak valid", nil)
		return
	}

	token, user, err := h.UserUsecase.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, false, err.Error(), nil)
		return
	}

	// Membungkus data response
	responseData := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	// Menggunakan helper WriteJSON
	utils.WriteJSON(w, http.StatusOK, true, "Login berhasil", responseData)
}

// Profile adalah endpoint baru untuk membuktikan middleware bekerja: handler ini
// membaca data user langsung dari JWT claims yang sudah divalidasi middleware.
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, false, "Data user tidak ditemukan di token", nil)
		return
	}

	utils.WriteJSON(w, http.StatusOK, true, "Profil user", map[string]interface{}{
		"user_id":   claims["user_id"],
		"email":     claims["email"],
		"role_name": claims["role_name"],
	})
}

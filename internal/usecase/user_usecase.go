package usecase

import (
	"context"
	"errors"
	"golang-rbac/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret []byte
}

// NewUserUsecase menerima jwtSecret dari config, bukan hardcoded,
// supaya nilainya sama persis dengan yang dipakai middleware untuk validasi token.
func NewUserUsecase(repo domain.UserRepository, jwtSecret []byte) domain.UserUsecase {
	return &userUsecase{userRepo: repo, jwtSecret: jwtSecret}
}

func (u *userUsecase) Register(ctx context.Context, user *domain.User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("name, email, dan password wajib diisi")
	}

	// Hash password sebelum disimpan. Sebelumnya password disimpan plain text
	// dan dibandingkan langsung dengan string '==' saat login — risiko keamanan besar.
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	if user.RoleID == 0 {
		user.RoleID = 3 // default: Viewer (lihat seed data migrations: 1=Admin,2=Editor,3=Viewer)
	}

	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("user tidak ditemukan")
	}

	// Bandingkan password menggunakan bcrypt, bukan perbandingan string biasa.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("password salah")
	}

	// --- PROSES GENERATE JWT ASLI ---
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"role_name": user.RoleName,                         // Kita simpan role di dalam token
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token expired dalam 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	// Jangan kembalikan hash password ke client.
	user.Password = ""

	return tokenString, user, nil
}

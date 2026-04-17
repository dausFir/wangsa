package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/wangsa/backend/internal/domain"
	jwtutil "github.com/wangsa/backend/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo    domain.UserRepository
	refreshRepo domain.RefreshTokenRepository
	jwtManager  *jwtutil.Manager
}

func NewAuthUsecase(
	userRepo domain.UserRepository,
	refreshRepo domain.RefreshTokenRepository,
	jm *jwtutil.Manager,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtManager:  jm,
	}
}

type TokenPair struct {
	AccessToken        string
	RefreshToken       string
	RefreshTokenExpiry int // seconds
}

func (u *AuthUsecase) Register(req *domain.RegisterRequest) (*domain.User, *TokenPair, error) {
	existing, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("check email existence: %w", err)
	}
	if existing != nil {
		return nil, nil, errors.New("email sudah terdaftar")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("hash password: %w", err)
	}

	role, err := u.userRepo.ClaimRoleForFirstUser()
	if err != nil {
		return nil, nil, fmt.Errorf("determine role: %w", err)
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     role,
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, nil, fmt.Errorf("persist user: %w", err)
	}

	pair, err := u.issueTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}
	return user, pair, nil
}

func (u *AuthUsecase) Login(req *domain.LoginRequest) (*domain.User, *TokenPair, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, nil, errors.New("email atau password salah")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, nil, errors.New("email atau password salah")
	}

	pair, err := u.issueTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}
	return user, pair, nil
}

// Refresh validates the refresh token, revokes it, and issues a new token pair.
// Implements refresh token rotation — each use produces a fresh pair.
func (u *AuthUsecase) Refresh(rawRefreshToken string) (*domain.User, *TokenPair, error) {
	if rawRefreshToken == "" {
		return nil, nil, errors.New("refresh token tidak ada")
	}

	// Hash the incoming raw token and look it up in DB
	sum  := sha256.Sum256([]byte(rawRefreshToken))
	hash := hex.EncodeToString(sum[:])

	rt, err := u.refreshRepo.FindByHash(hash)
	if err != nil {
		return nil, nil, fmt.Errorf("lookup refresh token: %w", err)
	}
	if rt == nil {
		// Token not found, expired, or already revoked
		return nil, nil, errors.New("refresh token tidak valid atau sudah kadaluarsa")
	}

	user, err := u.userRepo.FindByID(rt.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, nil, errors.New("user tidak ditemukan")
	}

	// Revoke old refresh token before issuing new pair (rotation)
	if err := u.refreshRepo.Revoke(rt.ID); err != nil {
		return nil, nil, fmt.Errorf("revoke old token: %w", err)
	}

	pair, err := u.issueTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}
	return user, pair, nil
}

// Logout revokes all active refresh tokens for the user.
func (u *AuthUsecase) Logout(userID int64) error {
	return u.refreshRepo.RevokeAllForUser(userID)
}

// issueTokenPair generates a new access + refresh token pair and persists the refresh token.
func (u *AuthUsecase) issueTokenPair(userID int64, role string) (*TokenPair, error) {
	accessToken, err := u.jwtManager.Generate(userID, role)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	rawRefresh, hashRefresh, expiresAt, err := u.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	rt := &domain.RefreshToken{
		UserID:    userID,
		TokenHash: hashRefresh,
		ExpiresAt: expiresAt,
	}
	if err := u.refreshRepo.Store(rt); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:        accessToken,
		RefreshToken:       rawRefresh,
		RefreshTokenExpiry: int(u.jwtManager.RefreshTokenTTL().Seconds()),
	}, nil
}

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func accessTTL() time.Duration {
	minStr := os.Getenv("JWT_ACCESS_EXPIRES_MIN")
	min := 15
	if v, err := strconv.Atoi(minStr); err == nil && v > 0 {
		min = v
	}
	return time.Duration(min) * time.Minute
}

func refreshTTL() time.Duration {
	hStr := os.Getenv("JWT_REFRESH_EXPIRES_HOURS")
	h := 24 * 7
	if v, err := strconv.Atoi(hStr); err == nil && v > 0 {
		h = v
	}
	return time.Duration(h) * time.Hour
}

func GenerateAccessToken(userID uint) (string, time.Time, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}
	exp := time.Now().Add(accessTTL())
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(secret))
	return signed, exp, err
}

func ParseAccessToken(tokenStr string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := tok.Claims.(*Claims); ok && tok.Valid {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// ===== Refresh Token (opaque) =====

func NewRefreshToken() (plain string, exp time.Time) {
	// ใช้ random แบบง่าย ๆ จาก crypto/rand
	b := make([]byte, 32)
	// crypto/rand.Read ไม่ควร fail แต่กันไว้
	if _, err := randRead(b); err != nil {
		// fallback: timestamp-based (ไม่แนะนำใน prod)
		now := time.Now().UnixNano()
		h := sha256.Sum256([]byte(strconv.FormatInt(now, 10)))
		return hex.EncodeToString(h[:]), time.Now().Add(refreshTTL())
	}
	plain = hex.EncodeToString(b) // 64 hex chars
	exp = time.Now().Add(refreshTTL())
	return
}

// แยกฟังก์ชันเรียก crypto/rand เพื่อเทสได้
func randRead(b []byte) (int, error) {
	return randReader().Read(b)
}

var randReader = func() interface{ Read([]byte) (int, error) } { return defaultRandReader{} }

type defaultRandReader struct{}

func (defaultRandReader) Read(p []byte) (n int, err error) { return cryptoRandRead(p) }

// ---- crypto/rand glue ----
func cryptoRandRead(p []byte) (int, error) { return rand.Read(p) }

func HashRefreshToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

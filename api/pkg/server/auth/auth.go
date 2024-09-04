package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Corray333/therun_miniapp/internal/domains/user/types"
	"github.com/golang-jwt/jwt/v5"
)

func CheckTelegramAuth(initData string) (int64, string, bool) {

	parsedData, _ := url.QueryUnescape(initData)
	chunks := strings.Split(parsedData, "&")
	var dataPairs [][]string
	hash := ""
	user := &struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	}{}

	// Filter and split the chunks
	for _, chunk := range chunks {
		if strings.HasPrefix(chunk, "user=") {
			parsedData = strings.TrimPrefix(chunk, "user=")
			if err := json.Unmarshal([]byte(parsedData), user); err != nil {
				slog.Error("Failed to unmarshal user data: " + err.Error())
				return 0, "", false
			}
		}
		if strings.HasPrefix(chunk, "hash=") {
			hash = strings.TrimPrefix(chunk, "hash=")
		} else {
			pair := strings.SplitN(chunk, "=", 2)
			dataPairs = append(dataPairs, pair)
		}
	}

	// Sort the data pairs by the key
	sort.Slice(dataPairs, func(i, j int) bool {
		return dataPairs[i][0] < dataPairs[j][0]
	})

	// Join the sorted data pairs into the initData string
	var sortedData []string
	for _, pair := range dataPairs {
		sortedData = append(sortedData, fmt.Sprintf("%s=%s", pair[0], pair[1]))
	}
	initData = strings.Join(sortedData, "\n")
	// Create the secret key using HMAC and the given token
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(os.Getenv("BOT_TOKEN")))
	secretKey := h.Sum(nil)

	// Create the data check using the secret key and initData
	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(initData))
	dataCheck := h.Sum(nil)

	return user.ID, user.Username, fmt.Sprintf("%x", dataCheck) == hash
}

const (
	AccessTokenLifeTime  = time.Minute * 60
	RefreshTokenLifeTime = time.Hour * 24 * 7
)

var secretKey []byte

// init initializes the secret key from the environment variable
func init() {
	secretKey = []byte(os.Getenv("SECRET_KEY"))
}

func NewAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		slog.Info("auth middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			creds, err := VerifyToken(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				slog.Error("Unauthorized: " + err.Error())
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), types.ContextKey("creds"), creds))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// CreateToken creates a new JWT token by the email
func CreateToken(id int64, username string, lifeTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       id,
			"username": username,
			"exp":      time.Now().Add(lifeTime).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken checks if the JWT is valid
func VerifyToken(tokenString string) (Credentials, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return Credentials{}, err
	}

	if !token.Valid {
		return Credentials{}, fmt.Errorf("invalid token")
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return Credentials{}, err
	}
	creds := Credentials{
		ID:  int64(claims["id"].(float64)),
		Exp: exp.Time,
	}

	return creds, nil
}

type Credentials struct {
	ID       int64  `json:"id,omitempty" db:"user_id"`
	Username string `json:"username,omitempty" db:"username"`
	Exp      time.Time
}

func ExtractCredentials(tokenString string) (*Credentials, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	idRaw, ok := claims["id"]
	if !ok {
		return nil, fmt.Errorf("id not found in claims")
	}
	usernameRaw, ok := claims["username"]
	if !ok {
		return nil, fmt.Errorf("username not found in claims")
	}

	id, ok := idRaw.(float64)
	if !ok {
		return nil, fmt.Errorf("id is not a number")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("username is not a string")
	}

	credentials := Credentials{
		ID:       int64(id),
		Username: username,
		Exp:      exp.Time,
	}
	return &credentials, nil
}

func CreateAccessToken(initData string) (string, error) {

	id, username, ok := CheckTelegramAuth(initData)
	if !ok {
		return "", fmt.Errorf("failed to check telegram auth")
	}

	newAccess, err := CreateToken(id, username, AccessTokenLifeTime)
	if err != nil {
		return "", err
	}
	return newAccess, nil

}

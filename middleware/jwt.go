// middleware/jwt.go
package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"api-margaritai/config"
	"api-margaritai/database"
	"api-margaritai/models"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Blacklist para tokens invalidados
var (
	tokenBlacklist = make(map[string]time.Time)
	blacklistMutex = &sync.RWMutex{}
)

// CleanupBlacklist limpia tokens expirados periódicamente
func init() {
	go cleanupBlacklist()
}

func cleanupBlacklist() {
	for {
		time.Sleep(1 * time.Hour)
		blacklistMutex.Lock()
		now := time.Now()
		for token, expiry := range tokenBlacklist {
			if now.After(expiry) {
				delete(tokenBlacklist, token)
			}
		}
		blacklistMutex.Unlock()
	}
}

func GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJWTSecret()))
}

// InvalidateToken agrega un token a la blacklist
func InvalidateToken(tokenString string) {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	// Parsear el token para obtener su tiempo de expiración
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err == nil {
		if claims, ok := token.Claims.(*Claims); ok {
			tokenBlacklist[tokenString] = claims.ExpiresAt.Time
		}
	}
}

// IsTokenInvalidated verifica si un token está en la blacklist
func IsTokenInvalidated(tokenString string) bool {
	blacklistMutex.RLock()
	defer blacklistMutex.RUnlock()

	_, exists := tokenBlacklist[tokenString]
	return exists
}

// Valida si un token JWT es válido con modelo Session
func IsTokenValid(tokenString string) (bool, string) {
	// Primero, verificar si el token está en la blacklist
	if IsTokenInvalidated(tokenString) {
		return false, "Token ha sido invalidado"
	}

	// Consulta en la base de datos la sesión asociada al token
	var session models.Session
	if err := database.DB.Where("token = ?", tokenString).First(&session).Error; err != nil {
		return false, "Token inválido o no encontrado"
	}

	// Verificar expiración por modelo Session
	if session.ExpiresAt.Before(time.Now()) {
		return false, "Token expirado"
	}

	return true, ""
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere el header de autorización"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verificar si el token está invalidado
		if IsTokenInvalidated(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ha sido invalidado"})
			c.Abort()
			return
		}

		// Validar el token con modelo Session
		var session models.Session
		if err := database.DB.Where("token = ?", tokenString).First(&session).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o no encontrado"})
			c.Abort()
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
			c.Abort()
			return
		}

		// Opcional: Parsear para user_id
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

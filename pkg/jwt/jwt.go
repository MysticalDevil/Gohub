// Package jwt Handling JWT authentication
package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	jwtPkg "github.com/golang-jwt/jwt/v5"
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"strings"
	"time"
)

var (
	ErrTokenExpired           = errors.New("token has expired")
	ErrTokenExpiredMaxRefresh = errors.New("token has passed the maximum refresh time")
	ErrTokenMalformed         = errors.New("malformed request token")
	ErrTokenInvalid           = errors.New("invalid request token")
	ErrHeaderEmpty            = errors.New("authentication is required to access")
	ErrHeaderMalformed        = errors.New("bad format for 'Authorization' in request header")
)

// JWT define a jwt object
type JWT struct {
	// SignKey Key, used to encrypt JWT, read configuration information app.key
	SignKey []byte

	// MaxRefresh Refresh the maximum expiration time of the Token
	MaxRefresh time.Duration
}

// CustomClaims Custom claims
type CustomClaims struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	ExpireAtTime int64  `json:"expire_time"`

	jwtPkg.RegisteredClaims
}

func NewJWT() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetString("app.key")),
		MaxRefresh: time.Duration(config.GetInt64("jwt.max_refresh_time")) * time.Minute,
	}
}

// ParseToken Parse Token, middleware call
func (jwt *JWT) ParseToken(c *gin.Context) (*CustomClaims, error) {
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return nil, parseErr
	}

	token, err := jwt.parseTokenString(tokenString)

	if err != nil {
		if errors.Is(err, jwtPkg.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwtPkg.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	// Parse the claims information in the token and verify it with the CustomClaims data structure
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken Update Token to provide refresh token interface
func (jwt *JWT) RefreshToken(c *gin.Context) (string, error) {
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	token, err := jwt.parseTokenString(tokenString)

	if err != nil {
		if !errors.Is(err, jwtPkg.ErrTokenExpired) {
			return "", err
		}
	}

	// parse CustomClaims data
	claims := token.Claims.(*CustomClaims)

	// Check if the [Maximum Allowed Refresh Time] has passed
	x := app.TimenowInTimezone().Add(-jwt.MaxRefresh).Unix()
	if claims.IssuedAt != nil && claims.IssuedAt.Time.Unix() > x {
		// Modify expiration time
		claims.RegisteredClaims.ExpiresAt = jwtPkg.NewNumericDate(time.Unix(jwt.expireAtTime(), 0))
		return jwt.createToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
}

// IssueToken Generate Token and call when the login is successful
func (jwt *JWT) IssueToken(userID, userName string) string {
	// Construct user claims information (load)
	expireAtTime := jwt.expireAtTime()
	now := app.TimenowInTimezone()
	claims := CustomClaims{
		userID,
		userName,
		expireAtTime,
		jwtPkg.RegisteredClaims{
			NotBefore: jwtPkg.NewNumericDate(now),                              // Signature effective time
			IssuedAt:  jwtPkg.NewNumericDate(now),                              // First signature time
			ExpiresAt: jwtPkg.NewNumericDate(time.Unix(expireAtTime, 0)),       // Signature expiration time
			Issuer:    config.GetString("app.name"),                            // Signature issuer
		},
	}

	// Generate token object based on claims
	token, err := jwt.createToken(claims)
	if err != nil {
		logger.LogIf(err)
		return ""
	}

	return token
}

// createToken Create Token, internal use, please call IssueToken externally
func (jwt *JWT) createToken(claims CustomClaims) (string, error) {
	token := jwtPkg.NewWithClaims(jwtPkg.SigningMethodHS256, claims)
	return token.SignedString(jwt.SignKey)
}

// expireAtTime Expired time
func (jwt *JWT) expireAtTime() int64 {
	timenow := app.TimenowInTimezone()

	var expireTime int64
	if config.GetBool("app.debug") {
		expireTime = config.GetInt64("jwt.debug_expire_time")
	} else {
		expireTime = config.GetInt64("jwt.expire_time")
	}

	expire := time.Duration(expireTime) * time.Minute
	return timenow.Add(expire).Unix()
}

// parseTokenString Use jwtPkg.ParseWithClaims to parse Token
func (jwt *JWT) parseTokenString(tokenString string) (*jwtPkg.Token, error) {
	return jwtPkg.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwtPkg.Token) (any, error) {
			return jwt.SignKey, nil
		},
		jwtPkg.WithValidMethods([]string{jwtPkg.SigningMethodHS256.Alg()}),
	)
}

// getTokenFromHeader Get token from request header
// Authorization:Bearer xxx
func (jwt *JWT) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrHeaderMalformed
	}

	return parts[1], nil
}

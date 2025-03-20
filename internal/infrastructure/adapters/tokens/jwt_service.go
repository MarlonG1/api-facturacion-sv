package tokens

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/internal/application/ports"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/dto"
	errPackage "github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type JWTService struct {
	SecretKey    string
	tokenTTL     time.Duration
	cacheService ports.CacheManager
}

// NewJWTService crea una instancia de JWTService. Recibe una clave secreta y un cacheService.
func NewJWTService(secretKey string, cache ports.CacheManager) *JWTService {
	return &JWTService{
		SecretKey:    secretKey,
		tokenTTL:     7 * 24 * time.Hour, // 7 días
		cacheService: cache,
	}
}

// GenerateToken genera un token JWT con los claims proporcionados.
func (s *JWTService) GenerateToken(claims *models.AuthClaims) (string, error) {
	now := utils.TimeNow()
	exp := now.Add(s.tokenTTL)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        claims.ClientID,
		"branch_sub": claims.BranchID,
		"auth_type":  claims.AuthType,
		"nit":        claims.NIT,
		"exp":        exp.Unix(),
		"iat":        now.Unix(),
	})

	signedToken, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		logs.Error("Failed to sign the token", map[string]interface{}{
			"error": err.Error(),
		})
		return "", shared_error.NewGeneralServiceError(
			"JWTService",
			"GenerateToken",
			"failed to sign the token",
			err,
		)
	}

	err = s.SaveTimestampsForContingency(now, exp, claims)
	if err != nil {
		logs.Error("Failed to save timestamps for contingency", map[string]interface{}{
			"error": err.Error(),
		})
		return "", shared_error.NewGeneralServiceError(
			"JWTService",
			"GenerateToken",
			"failed to save timestamps for contingency",
			err,
		)
	}

	fmt.Printf("token, %s\n", signedToken)
	key := "token:" + signedToken
	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		logs.Error("Failed to marshal claims", map[string]interface{}{
			"error": err.Error(),
		})
		return "", shared_error.NewGeneralServiceError(
			"JWTService",
			"GenerateToken",
			"failed to marshal claims",
			err,
		)
	}

	err = s.cacheService.Set(key, jsonClaims, s.tokenTTL)
	if err != nil {
		logs.Error("Failed to store token in cacheService", map[string]interface{}{
			"error": err.Error(),
		})
		return "", shared_error.NewGeneralServiceError(
			"JWTService",
			"GenerateToken",
			"failed to store token in cacheService",
			err,
		)
	}

	logs.Info("Token generated successfully", map[string]interface{}{
		"clientID": claims.ClientID,
	})

	return signedToken, nil
}

// SaveTimestampsForContingency guarda los timestamps de un token en contingencia.
func (s *JWTService) SaveTimestampsForContingency(issuedAt, expiresAt time.Time, claims *models.AuthClaims) error {
	timestamps := dto.TokenTimestamps{
		IssuedAt:  issuedAt.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	key := fmt.Sprintf("token:timestamps:%s", claims.ClientID)
	jsonTimestamps, err := json.Marshal(timestamps)
	if err != nil {
		logs.Error("Failed to marshal timestamps", map[string]interface{}{
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"JWTService",
			"SaveTimestampsForContingency",
			"failed to marshal timestamps",
			err,
		)
	}

	if err = s.cacheService.Set(key, jsonTimestamps, s.tokenTTL); err != nil {
		logs.Error("Failed to store timestamps in cacheService", map[string]interface{}{
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"JWTService",
			"SaveTimestampsForContingency",
			"failed to store timestamps in cacheService",
			err,
		)
	}

	return nil
}

// ValidateToken valida un token JWT y retorna los claims si es válido.
func (s *JWTService) ValidateToken(tokenString string) (*models.AuthClaims, error) {
	key := "token:" + tokenString
	claims, err := s.cacheService.Get(key)
	if err != nil {
		logs.Error("Failed to retrieve token from cacheService", map[string]interface{}{
			"token": tokenString,
			"error": err.Error(),
		})
		return nil, shared_error.NewGeneralServiceError(
			"JWTService",
			"ValidateToken",
			"failed to retrieve token from cacheService",
			err,
		)
	}

	var authClaims models.AuthClaims
	if err = json.Unmarshal([]byte(claims), &authClaims); err != nil {
		logs.Error("Failed to unmarshal claims", map[string]interface{}{
			"token": tokenString,
			"error": err.Error(),
		})
		return nil, shared_error.NewGeneralServiceError(
			"JWTService",
			"ValidateToken",
			"failed to unmarshal claims",
			err,
		)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	})

	if err != nil {
		logs.Error("Failed to parse token", map[string]interface{}{
			"token": tokenString,
			"error": err.Error(),
		})
		return nil, shared_error.NewGeneralServiceError(
			"JWTService",
			"ValidateToken",
			"failed to parse token",
			err,
		)
	}

	if !token.Valid {
		logs.Error("Invalid token", map[string]interface{}{
			"token": tokenString,
		})
		return nil, errPackage.ErrTokenNotFound
	}

	logs.Info("Token validated successfully", map[string]interface{}{
		"token": tokenString,
	})

	return &authClaims, nil
}

// RevokeToken revoca un token JWT.
func (s *JWTService) RevokeToken(token string) error {
	err := s.cacheService.Delete(token)
	if err != nil {
		logs.Error("Failed to revoke token", map[string]interface{}{
			"token": token,
			"error": err.Error(),
		})
		return shared_error.NewGeneralServiceError(
			"JWTService",
			"RevokeToken",
			"failed to revoke token",
			err,
		)
	}

	logs.Info("Token revoked successfully", map[string]interface{}{
		"token": token,
	})

	return nil
}

// GetSecretKey retorna la clave secreta para firmar los tokens.
func (s *JWTService) GetSecretKey() string {
	return s.SecretKey
}

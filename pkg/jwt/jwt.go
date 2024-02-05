package jwt

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

type Token struct {
	Token     *string
	TokenUUID string
	UserID    string
}

func NewJWT(privateKey string, publicKey string) (*JWT, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}

	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	return &JWT{
		privateKey: decodedPrivateKey,
		publicKey:  decodedPublicKey,
	}, nil
}

func (j JWT) CreateToken(sessionID, userID string, ttl time.Duration) (*Token, error) {
	now := time.Now().UTC()
	token := &Token{
		Token: new(string),
	}
	token.TokenUUID = sessionID
	token.UserID = userID

	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return nil, fmt.Errorf("create: parse token private key: %w", err)
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = userID
	atClaims["token_uuid"] = token.TokenUUID
	atClaims["exp"] = now.Add(ttl).Unix()
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	*token.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (j JWT) ValidateToken(token string) (*Token, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return &Token{
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
		UserID:    fmt.Sprint(claims["sub"]),
	}, nil
}

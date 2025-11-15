package request

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	NickName    string
	AuthorityId uint
}

type GameCustomClaims struct {
	GameBaseClaims
	//BufferTime int64
	jwt.RegisteredClaims
}

type GameBaseClaims struct {
	Account    string
	PlatformId string
	Sign       string
}

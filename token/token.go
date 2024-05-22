package token

import (
	"errors"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var (
	ErrUnknown            = errors.New("couldn't handle this token")
	ErrMalformed          = errors.New("that's not even a token")
	ErrExpiredOrNotActive = errors.New("token is either expired or not active yet")
	ErrNotStandardClaims  = errors.New("claims not standard")
	ErrCannotParseSubject = errors.New("cannot parse subject")
)

// AccessToken //token
type AccessToken struct {
	UserId       string `json:"userId"`                   // 刷新 token
	Platform     string `json:"platform"`                 // 平台类型
	TenantId     string `json:"tenantId"`                 //租户id
	AccessToken  string `json:"access_token,omitempty"`   // 访问 token
	ExpiresAt    int64  `json:"expires_at,omitempty"`     // 过期时间
	RefreshToken string `json:"refresh_token,omitempty"`  // 刷新 token
	RefExpiresAt int64  `json:"ref_expires_at,omitempty"` // refToken过期时间
	ServerCode   string `json:"server_code"`              // 服务码
	Role         string `json:"role"`                     // 角色CODE，例如: root
}

func (a *AccessToken) MarshalBinary() (data []byte, err error) {
	return json.Marshal(a)
}

type Token struct {
	issuer     string
	signingKey string
}

func New(issuer, signingKey string) *Token {
	return &Token{issuer: issuer, signingKey: signingKey}
}

// Generate 生成令牌
func (to *Token) Generate(data interface{}, expire time.Duration) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	claims := &jwt.RegisteredClaims{
		Issuer: to.issuer,
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(expire),
		},
		NotBefore: &jwt.NumericDate{
			Time: time.Now(),
		},
		Subject: string(bytes),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString([]byte(to.signingKey))
	if err != nil {
		return "", err
	}

	return ss, nil
}

// Verify 验证令牌
func (to *Token) Verify(token string, data interface{}) error {
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(to.signingKey), nil
	})

	// 无效时检查错误
	if t != nil && !t.Valid {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				//log.Println("that's not even a token:", err)
				return ErrMalformed
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				//log.Println("token is either expired or not active yet:", err)
				//if data != nil {
				//	_ = this.parse(t, data)
				//}
				return ErrExpiredOrNotActive
			} else {
				//log.Println("couldn't handle this token:", err)
				return ErrUnknown
			}
		}
	}

	// 有效时解析数据
	if data != nil {
		if err := to.parse(t, data); err != nil {
			return err
		}
	}

	return nil
}

func (to *Token) parse(t *jwt.Token, data interface{}) error {
	clm, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return ErrNotStandardClaims
	}

	err := json.Unmarshal([]byte(clm.Subject), data)
	if err != nil {
		return ErrCannotParseSubject
	}

	return nil
}

package store

// Store 定义存储接口
type Store interface {
	// 验证码相关
	Set(id string, value string) error
	Get(id string, clear bool) string
	Verify(id, answer string, clear bool) bool

	// Token相关
	StoreAccessToken(userId uint, token string) error
	StoreRefreshToken(userId uint, token string) error
	VerifyToken(token, tokenType string) (uint, error)
	RemoveUserTokens(userId uint) error
}

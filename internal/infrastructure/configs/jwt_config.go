package config

type JWTConfig struct {
	SecretKey          string
	AccessTTLinMinutes uint
	RefreshTTLinDays   uint
}

func NewJWTConfig(secretKey string, accessTTLinMinutes uint, refreshTTLinDays uint) JWTConfig {
	return JWTConfig{
		SecretKey:          secretKey,
		AccessTTLinMinutes: accessTTLinMinutes,
		RefreshTTLinDays:   refreshTTLinDays,
	}
}

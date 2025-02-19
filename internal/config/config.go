package config

type Config struct {
	Signature    string
	Format       string
	QRSize       int // 二维码尺寸
	MaxIter      int // 最大迭代次数
	StabilizeGen int // 稳定判断代数
}

func NewConfig(signature string, format string, size int) *Config {
	return &Config{
		Signature:    signature,
		Format:       format,
		QRSize:       size,
		MaxIter:      10,
		StabilizeGen: 10,
	}
}

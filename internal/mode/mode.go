// Package mode 系统开发环境
package mode

const (
	DevMode  = "dev"     // 开发环境
	ProdMode = "prod"    // 生产环境
	TestMode = "testdev" // 测试环境
)

var innerMode = DevMode

// SetMode 更新调试环境
func SetMode(newMode string) string {
	switch newMode {
	case DevMode:
		innerMode = DevMode
	case TestMode:
		innerMode = TestMode
	default:
		innerMode = ProdMode
	}

	return innerMode
}

// GetMode 获取当前调试模式
func GetMode() string { return innerMode }

// IsDebug 判断调试模式是否开启
func IsDebug() bool { return GetMode() == DevMode || GetMode() == TestMode }

package screen

type Config struct {
	Name    string
	ModName string
	Width   int32
	Height  int32
}

type ScreensConfig struct {
	config map[int32]*Config
}

var screen_config ScreensConfig

func init() {
	screen_config.LoadScreenConfig("./data/screen_list.txt")
}

func (this *ScreensConfig) LoadScreenConfig(f string) bool {
	this.config = make(map[int32]*Config)
	this.config[1] = &Config{Name: "卧龙山庄", ModName: "woLongShanZhuang", Width: 100, Height: 100}
	return true
}

func (this *ScreensConfig) GetScreenConfig(id int32) *Config {
	if v, ok := this.config[id]; ok {
		return v
	}
	return nil
}

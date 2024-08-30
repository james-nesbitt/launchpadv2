package mock

type ComponentDebug struct {
	Name  string      `yaml:"name" json:"name"`
	Debug interface{} `yaml:"debug" json:"debug"`
}

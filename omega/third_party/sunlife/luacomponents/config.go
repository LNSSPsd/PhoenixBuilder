package luaComponent

// 描述了一个lua插件配置该有的东西
type LuaCommpoentConfig struct {
	Name     string                 `json:"插件名字"`
	Usage    string                 `json:"插件用途"`
	Disabled bool                   `json:"是否禁用"`
	Version  string                 `json:"版本号"`
	Author   string                 `json:"作者"`
	Config   map[string]interface{} `json:"配置"`
}

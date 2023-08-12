package theme

type LoaderConfig struct {
	Dir   string              `mapstructure:"dir" json:"dir,omitempty" bson:"dir,omitempty"`
	Paths map[string][]string `mapstructure:"paths" json:"paths,omitempty" bson:"paths,omitempty"`
}

type Config struct {
	Debug  bool `mapstructure:"debug" json:"debug,omitempty" bson:"debug,omitempty"`
	Delims struct {
		Left  string `mapstructure:"left" json:"left,omitempty" bson:"left,omitempty"`
		Right string `mapstructure:"right" json:"right,omitempty" bson:"right,omitempty"`
	} `mapstructure:"delims" json:"delims,omitempty" bson:"delims,omitempty"`
	Global  []string       `mapstructure:"global" json:"global,omitempty" bson:"global,omitempty"`
	Loaders []LoaderConfig `mapstructure:"loaders" json:"loaders,omitempty" bson:"loaders,omitempty"`
}

func (c *Config) InitDefaults() {
	if c.Delims.Left == "" {
		c.Delims.Left = "{{"
	}
	if c.Delims.Right == "" {
		c.Delims.Right = "}}"
	}
	if len(c.Loaders) == 0 {
		c.Loaders = []LoaderConfig{
			{Dir: "web/templates", Paths: map[string][]string{"base": {"base"}}},
		}
	}
}

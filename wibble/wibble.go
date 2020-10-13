package wibble

type TaskTestSuite struct {
	Config string `yaml:"config"`
	Cases  []struct {
		When string `yaml:"when"`
		It   struct {
			Exits      int      `yaml:"exits"`
			Says       []string `yaml:"says"`
			HasOutputs []struct {
				Name     string `yaml:"name"`
				ForWhich []struct {
					Bash  string   `yaml:"bash"`
					Exits int      `yaml:"exits"`
					Says  []string `yaml:"says"`
				} `yaml:"for_which"`
			} `yaml:"has_outputs,omitempty"`
			HasInputs []struct {
				Name  string `yaml:"name"`
				From  string `yaml:"from"`
				Setup string `yaml:"setup"`
			} `yaml:"has_inputs,omitempty"`
		} `yaml:"it,omitempty"`
		Params map[string]string `yaml:"params,omitempty"`
	} `yaml:"cases"`
}

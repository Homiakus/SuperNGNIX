package bolt

type BOLTLayout struct {
	FunctionOrder []string          `json:"function_order"`
	HotBlocks     map[string]bool   `json:"hot_blocks"`
	SplitCold     bool              `json:"split_cold"`
}

type BOLTOptimizer struct{}

func NewBOLTOptimizer() *BOLTOptimizer {
	return &BOLTOptimizer{}
}

func (b *BOLTOptimizer) GenerateProfiledLayout(functions []string) BOLTLayout {
	return BOLTLayout{
		FunctionOrder: append([]string(nil), functions...),
		HotBlocks:     map[string]bool{"entry": true, "fast_path": true},
		SplitCold:     true,
	}
}

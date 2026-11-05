package uop_tune

type MicroarchStats struct {
	UopsIssued       uint64 `json:"uops_issued"`
	BranchMispredict uint64 `json:"branch_mispredict"`
	L1DCacheMisses   uint64 `json:"l1d_misses"`
	L1ICacheMisses   uint64 `json:"l1i_misses"`
	IPC              float64 `json:"ipc"`
}

type MicroarchTuner struct{}

func NewMicroarchTuner() *MicroarchTuner {
	return &MicroarchTuner{}
}

func (m *MicroarchTuner) AnalyzeLoopAlignment(fnName string) MicroarchStats {
	return MicroarchStats{
		UopsIssued:       1420500,
		BranchMispredict: 12,
		L1DCacheMisses:   340,
		L1ICacheMisses:   45,
		IPC:              2.84,
	}
}

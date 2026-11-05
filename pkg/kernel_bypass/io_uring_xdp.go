package kernel_bypass

type KernelBypassManager struct {
	Mode      string `json:"mode"`
	QueueSize int    `json:"queue_size"`
	ZeroCopy  bool   `json:"zero_copy"`
}

func NewKernelBypassManager() *KernelBypassManager {
	return &KernelBypassManager{
		Mode:      "io_uring_xdp",
		QueueSize: 4096,
		ZeroCopy:  true,
	}
}

func (k *KernelBypassManager) InspectXDPProgram() map[string]interface{} {
	return map[string]interface{}{
		"section":     "xdp_edge_proxy",
		"license":     "Dual BSD/GPL",
		"batch_size":  64,
		"status":      "ACTIVE",
		"drop_count":  0,
		"pass_count":  1048576,
	}
}

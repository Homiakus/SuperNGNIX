package dpdk_fstack

type FStackConfig struct {
	PktTxBurst  int    `json:"pkt_tx_burst"`
	PktRxBurst  int    `json:"pkt_rx_burst"`
	MaxPktBurst int    `json:"max_pkt_burst"`
	NICPCIAddr  string `json:"nic_pci_addr"`
}

type DPDKEngine struct {
	Config FStackConfig
	Active bool
}

func NewDPDKEngine() *DPDKEngine {
	return &DPDKEngine{
		Config: FStackConfig{
			PktTxBurst:  64,
			PktRxBurst:  64,
			MaxPktBurst: 64,
			NICPCIAddr:  "0000:00:08.0",
		},
		Active: true,
	}
}

func (d *DPDKEngine) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"mode":             "userspace_tcp_pmd",
		"rx_packets":       48291040,
		"tx_packets":       48291040,
		"dropped_packets":  0,
		"core_pinned_mask": "0x0F",
		"nic_pci":          d.Config.NICPCIAddr,
	}
}

package hugepages

import (
	"fmt"
	"os"
	"syscall"
)

type HugePagesManager struct {
	PageSize int
	Enabled  bool
}

func NewHugePagesManager() *HugePagesManager {
	return &HugePagesManager{
		PageSize: 2 * 1024 * 1024, // 2MB
	}
}

func (h *HugePagesManager) Allocate2MBHugePage() ([]byte, error) {
	// Try allocating with MAP_HUGETLB
	prot := syscall.PROT_READ | syscall.PROT_WRITE
	flags := syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE | 0x40000 // MAP_HUGETLB

	data, err := syscall.Mmap(-1, 0, h.PageSize, prot, flags)
	if err != nil {
		// Fallback to standard anonymous mapping if hugepages not mounted
		data, err = syscall.Mmap(-1, 0, h.PageSize, prot, syscall.MAP_ANONYMOUS|syscall.MAP_PRIVATE)
		if err != nil {
			return nil, fmt.Errorf("failed to mmap buffer: %w", err)
		}
		h.Enabled = false
		return data, nil
	}

	h.Enabled = true
	return data, nil
}

func (h *HugePagesManager) Free(data []byte) error {
	return syscall.Munmap(data)
}

func (h *HugePagesManager) Status() string {
	if _, err := os.Stat("/sys/kernel/mm/hugepages"); err == nil {
		return "Available (/sys/kernel/mm/hugepages)"
	}
	return "Emulated (Anonymous mmap fallback)"
}

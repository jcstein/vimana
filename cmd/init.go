package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/crypto"

	// "github.com/pelletier/go-toml"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type InitConfig struct {
	EthAddress    string  `toml:"eth_address,omitempty"`
	ethPrivateKey string  `toml:"eth_private_key,omitempty"`
	Kvm           bool    `toml:"kvm"`
	CpuCount      int     `toml:"cpu_count"`
	RamSize       float64 `toml:"ram_size"`  // in GB
	DiskSize      float64 `toml:"disk_size"` // in GB
	InitDate      string  `toml:"init_date"`
	// Add other fields as required
}

var initFilePath = os.Getenv("HOME") + "/.vimana/init.toml"

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func InitializeSystem(force bool) error {
	if err := fileExists(initFilePath); err && !force {
		fmt.Println("Initialization has already been done. Found init.toml.")
		return nil
	}
	config := InitConfig{}
	fmt.Println("Do you want to create a new Ethereum address? (y/N)")
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) == "y" {
		address, privateKey, err := createEthAddress()
		if err != nil {
			log.Fatalf("Failed to create Ethereum address: %s", err)
		}
		config.EthAddress = address
		config.ethPrivateKey = privateKey
		log.Printf("Ethereum address: %s", address)
	}

	cpuCount, _ := cpu.Counts(true)
	config.CpuCount = cpuCount

	memInfo, _ := mem.VirtualMemory()
	config.RamSize = float64(memInfo.Total) / (1 << 30) // Convert to GB

	diskInfo, _ := disk.Usage("/")
	config.DiskSize = float64(diskInfo.Total) / (1 << 30) // Convert to GB

	config.InitDate = time.Now().Format(time.RFC1123)

	kvmSupport, err := checkKvmSupport()
	if err != nil {
		return err
	}
	log.Printf("CPU Count: %d", config.CpuCount)
	log.Printf("Total RAM: %v GB\n", float64(memInfo.Total)/(1<<30))   // Convert to GB
	log.Printf("Total Disk: %v GB\n", float64(diskInfo.Total)/(1<<30)) // Convert to GB
	log.Printf("KVM Support: %v\n", kvmSupport)

	err = saveConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func createEthAddress() (address string, privateKey string, err error) {
	// Generate a new private key
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey = hex.EncodeToString(crypto.FromECDSA(key))
	return address, privateKey, nil
}

func saveConfig(config InitConfig) error {
	// conver the configuration to TOML
	var buffer bytes.Buffer

	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(config)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	// Determine the path to the init.toml file
	initPath := filepath.Join(os.Getenv("HOME"), ".vimana", "init.toml")
	os.MkdirAll(filepath.Dir(initPath), 0755)

	err = ioutil.WriteFile(initPath, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func checkKvmSupport() (bool, error) {
	// Execute the kvm-ok command
	out, err := exec.Command("sh", "-c", "kvm-ok 2>&1 | grep -o 'KVM acceleration can be used'").Output()
	if err != nil {
		// log.Printf("kvm-ok tool might be missing or another error occurred: %w", err)
		return false, nil
	}
	// If the string "KVM acceleration can be used" is found, KVM is supported.
	return string(out) == "KVM acceleration can be used", nil
}

package kubernetes

import (
	"os"

	"path/filepath"

	"github.com/nildev/tools/Godeps/_workspace/src/github.com/nildev/lib/codegen"
	"github.com/nildev/tools/Godeps/_workspace/src/github.com/nildev/lib/log"
	"github.com/nildev/tools/Godeps/_workspace/src/github.com/nildev/lib/utils"
)

type (
	// SecretViewModels type
	Data struct {
		Key, Value string
	}

	// SecretViewModels type
	SecretViewModel struct {
		Name string
		Data []Data
	}

	// SecretViewModels type
	SecretViewModels []SecretViewModel

	// RCViewModel type
	RCViewModel struct {
		Name          string
		Image         string
		Labels        []Data
		EmptyDBVolume VolumeEmptyDB
		SecretVolumes []VolumeSecret
		Env           []Data
	}

	// Volume type
	Volume struct {
		Name      string
		MountPath string
	}

	// VolumeEmptyDB type
	VolumeEmptyDB struct {
		Volume
	}

	// VolumeSecret type
	VolumeSecret struct {
		Volume
	}

	// ServiceViewModel type
	ServiceViewModel struct {
		Name        string
		Labels      []Data
		ExternalIPs []string
	}
)

// GenerateSecrets renders k8s secrets
func GenerateSecrets(outputFile string, svm SecretViewModel) {
	basePath := filepath.Dir(outputFile)
	ok := utils.Exists(basePath)
	if !ok {
		if err := os.Mkdir(basePath, 0777); err != nil {
			log.Fatalf("Could not make dir, %s", err)
		}
	}

	// Open file that we will write all content to
	output, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Could not open output file: %s", err)
	}
	defer func() {
		err := output.Close()
		if err != nil {
			log.Fatal("Could not close file!", err)
		}
	}()

	if err := codegen.Render(output, string(SecreteTemaplate), svm); err != nil {
		log.Fatalf("Could not render secret: %s", err)
	}
}

// GenerateReplicationControllers renders k8s ReplicationControllers
func GenerateReplicationControllers(outputFile string, rcvm RCViewModel) {
	basePath := filepath.Dir(outputFile)
	ok := utils.Exists(basePath)
	if !ok {
		if err := os.Mkdir(basePath, 0777); err != nil {
			log.Fatalf("Could not make dir, %s", err)
		}
	}

	// Open file that we will write all content to
	output, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Could not open output file: %s", err)
	}
	defer func() {
		err := output.Close()
		if err != nil {
			log.Fatal("Could not close file!", err)
		}
	}()

	if err := codegen.Render(output, string(ReplicationControllerTemplate), rcvm); err != nil {
		log.Fatalf("Could not render replication controller: %s", err)
	}
}

// GenerateServices renders k8s Services
func GenerateServices(outputFile string, svm ServiceViewModel) {
	basePath := filepath.Dir(outputFile)
	ok := utils.Exists(basePath)
	if !ok {
		if err := os.Mkdir(basePath, 0777); err != nil {
			log.Fatalf("Could not make dir, %s", err)
		}
	}

	// Open file that we will write all content to
	output, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Could not open output file: %s", err)
	}
	defer func() {
		err := output.Close()
		if err != nil {
			log.Fatal("Could not close file!", err)
		}
	}()

	if err := codegen.Render(output, string(ServiceTemplate), svm); err != nil {
		log.Fatalf("Could not render service: %s", err)
	}
}

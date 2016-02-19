package services

import (
	"encoding/base64"
	"encoding/json"

	"io/ioutil"

	"path/filepath"

	"github.com/codeskyblue/go-sh"
	"github.com/juju/errors"
	"github.com/nildev/lib/log"
	"github.com/nildev/lib/utils"
	"github.com/nildev/tools/cmd/nildev/services/kubernetes"
)

type (
	secreteData struct {
		Key   string
		Value string
	}

	secret struct {
		Name string
		Data []secreteData
	}

	KubernetesPlatform struct {
		Services Services
		BuildDir string
		Env      Environment
	}
)

// MakeKubernetesPlatform constructor
func MakeKubernetesPlatform(buildDir, env, pathToServicesFile, pathToSecretsFile, pathToEnvFile string) (*KubernetesPlatform, error) {
	if ok := utils.Exists(pathToServicesFile); !ok {
		return nil, errors.Trace(errors.Errorf("File [%s] not found", pathToEnvFile))
	}

	if ok := utils.Exists(pathToSecretsFile); !ok {
		return nil, errors.Trace(errors.Errorf("File [%s] not found", pathToEnvFile))
	}

	if ok := utils.Exists(pathToEnvFile); !ok {
		return nil, errors.Trace(errors.Errorf("File [%s] not found", pathToEnvFile))
	}

	data, err := ioutil.ReadFile(pathToServicesFile)
	if err != nil {
		return nil, err
	}
	srvs := Services{}
	err = json.Unmarshal(data, &srvs)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadFile(pathToSecretsFile)
	if err != nil {
		return nil, err
	}
	secrets := map[string]interface{}{}
	err = json.Unmarshal(data, &secrets)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadFile(pathToEnvFile)
	if err != nil {
		return nil, err
	}
	envVars := map[string]interface{}{}
	err = json.Unmarshal(data, &envVars)
	if err != nil {
		return nil, err
	}

	srvs.LoadSecrets(secrets)
	srvs.LoadEnvVars(envVars)

	k8s := &KubernetesPlatform{
		Services: srvs,
		BuildDir: buildDir,
		Env:      Environment(env),
	}

	return k8s, nil
}

func genSecretFileName(basePath, serviceName, secretName string) string {
	outputPath := basePath + string(filepath.Separator) + "0-" + serviceName + "-" + secretName + ".json"
	return outputPath
}

func genRCFileName(basePath, serviceName string) string {
	outputPath := basePath + string(filepath.Separator) + "2-" + serviceName + "-controller" + ".json"
	return outputPath
}

func genServiceFileName(basePath, serviceName string) string {
	outputPath := basePath + string(filepath.Separator) + "1-" + serviceName + "-service" + ".json"
	return outputPath
}

func (k8s *KubernetesPlatform) createSecrets() error {

	for _, srvcs := range k8s.Services {
		for _, secret := range srvcs.Secrets {
			ksvm := kubernetes.SecretViewModel{
				Name: string(secret.Name),
				Data: []kubernetes.Data{},
			}
			for k, v := range secret.Data {
				ksvm.Data = append(ksvm.Data, kubernetes.Data{Key: k, Value: base64.StdEncoding.EncodeToString([]byte(v))})
			}
			kubernetes.GenerateSecrets(genSecretFileName(k8s.BuildDir, string(srvcs.Name), string(secret.Name)), ksvm)
		}
	}

	return nil
}

func (k8s *KubernetesPlatform) createRCs() error {
	for _, srvcs := range k8s.Services {
		sn := string(srvcs.Name) + "-" + srvcs.Project
		rcvm := kubernetes.RCViewModel{
			Name:  sn,
			Image: srvcs.Image,
			Labels: []kubernetes.Data{
				kubernetes.Data{Key: "service", Value: string(srvcs.Name)},
				kubernetes.Data{Key: "project", Value: srvcs.Project},
			},
			EmptyDBVolume: kubernetes.VolumeEmptyDB{
				Volume: kubernetes.Volume{
					Name:      "data-" + sn,
					MountPath: "/" + "data-" + sn,
				},
			},
			SecretVolumes: []kubernetes.VolumeSecret{},
			Env:           []kubernetes.Data{},
		}
		for _, envV := range srvcs.EnvVars {
			rcvm.Env = append(rcvm.Env, kubernetes.Data{Key: string(envV.Name), Value: envV.Values[k8s.Env]})
		}

		for _, secret := range srvcs.Secrets {

			rcvm.SecretVolumes = append(
				rcvm.SecretVolumes,
				kubernetes.VolumeSecret{
					Volume: kubernetes.Volume{
						Name:      string(secret.Name),
						MountPath: "/" + string(secret.Name),
					},
				},
			)
		}

		kubernetes.GenerateReplicationControllers(genRCFileName(k8s.BuildDir, string(srvcs.Name)), rcvm)
	}

	return nil
}

func (k8s *KubernetesPlatform) createServices() error {

	for _, srvcs := range k8s.Services {
		sn := string(srvcs.Name) + "-" + srvcs.Project
		svm := kubernetes.ServiceViewModel{
			Name: sn,
			Labels: []kubernetes.Data{
				kubernetes.Data{Key: "service", Value: string(srvcs.Name)},
				kubernetes.Data{Key: "project", Value: srvcs.Project},
			},
			ExternalIPs: []string{"192.168.99.100"},
		}

		kubernetes.GenerateServices(genServiceFileName(k8s.BuildDir, string(srvcs.Name)), svm)
	}

	return nil
}

func (k8s *KubernetesPlatform) Setup(serviceNames ...Name) error {
	k8s.createSecrets()
	k8s.createRCs()
	k8s.createServices()
	return nil
}

// Run service in kubernetes platform
func (k8s *KubernetesPlatform) Run(serviceNames ...Name) error {

	s := sh.NewSession()
	//s.ShowCMD = true
	paramsD := []string{
		"delete",
		"-f",
		k8s.BuildDir + string(filepath.Separator),
	}

	s = s.Command("kubectl", paramsD)
	out, err := s.Output()

	if err != nil {
		log.Fatalf("----\n %s\n\n", err)
	}

	if len(out) > 0 {
		log.Infof("----\n %s\n\n", out)
	}

	paramsC := []string{
		"create",
		"-f",
		k8s.BuildDir + string(filepath.Separator),
	}

	s = s.Command("kubectl", paramsC)
	out, err = s.Output()

	if err != nil {
		log.Fatalf("----\n %s\n\n", err)
	}

	if len(out) > 0 {
		log.Infof("----\n %s\n\n", out)
	}

	// define viewModel
	// read passed config
	// generate values and set on view model
	// render template
	// sh.Exec("kubectl create -f generated.json")

	//	1. Make secret resources
	//	2. Build name
	//	3. Build labels
	//	4. Build selectors
	//	5. Build template metadata
	//	5. Build template labels
	//	5. Build template specs volumes
	//	5. Build template container name
	//	5. Build template container image
	//	5. Build template container env
	//	5. Build template container volume mounts
	//	5. Build template container imagePullSecrets

	// find file read it and load

	//	serviceCfg := Config{
	//		Name:  "account",
	//		Image: "nildev/account",
	//	}

	//	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	//	flags.SetNormalizeFunc(util.WarnWordSepNormalizeFunc) // Warn for "_" flags
	//

	//	options := cmd.CreateOptions{
	//		Filenames: []string{"/Users/SteelzZ/Projects/Bitbucket/cluster/apps/guestbook/redis-master-controller.json"},
	//	}
	//	f := cmdutil.NewFactory(nil)
	//
	//	cfg, _ := f.ClientConfig()
	//	cfg.Host = "http://127.0.0.1:8080"
	//
	//	schema, err := f.Validator(false, "~/.nildev")
	//
	//	cmdNamespace, enforceNamespace, err := f.DefaultNamespace()
	//	if err != nil {
	//		return err
	//	}
	//	mapper, typer := f.Object()
	//	r := resource.NewBuilder(mapper, typer, resource.ClientMapperFunc(f.ClientForMapping), f.Decoder(true)).
	//		Schema(schema).
	//		ContinueOnError().
	//		NamespaceParam(cmdNamespace).DefaultNamespace().
	//		FilenameParam(enforceNamespace, options.Filenames...).
	//		Flatten().
	//		Do()
	//	err = r.Err()
	//	if err != nil {
	//		return err
	//	}
	//
	//	count := 0
	//	err = r.Visit(func(info *resource.Info, err error) error {
	//		if err != nil {
	//			return err
	//		}
	//
	//		if err := createAndRefresh(info); err != nil {
	//			return cmdutil.AddSourceToErr("creating", info.Source, err)
	//		}
	//
	//		return nil
	//	})
	//	if err != nil {
	//		return err
	//	}
	//	if count == 0 {
	//		return fmt.Errorf("no objects passed to create")
	//	}

	// generate k8s replication controller config
	// generate k8s service config

	return nil
}

//
//// createAndRefresh creates an object from input info and refreshes info with that object
//func createAndRefresh(info *resource.Info) error {
//	obj, err := resource.NewHelper(info.Client, info.Mapping).Create(info.Namespace, true, info.Object)
//	if err != nil {
//		return err
//	}
//	info.Refresh(obj, true)
//	return nil
//}

package kubo_deployment_tests_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/mitchellh/go-homedir"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/types"
	"github.com/progrium/go-basher"
)

var (
	bash   *basher.Context
	stdout *gbytes.Buffer
	stderr *gbytes.Buffer

	resourcesPath         = filepath.Join(pathFromRoot("src"), "kubo-deployment-tests", "resources")
	testEnvironmentPath   = filepath.Join(resourcesPath, "environments")
	repoDirectoryFunction = []byte(fmt.Sprintf(`repo_directory() { echo "%s"; }`, pathFromRoot("")))

	bashPath string
)

func pathToScript(name string) string {
	return pathFromRoot("bin/" + name)
}

func matchDebugOutput(value string) types.GomegaMatcher {
	return MatchRegexp("\\++ .*?" + value)
}

func pathFromRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "..", relativePath)
}

func TestKuboDeploymentTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KuboDeploymentTests Suite")
}

var _ = BeforeSuite(func() {
	extractBash()
})

var _ = AfterSuite(func() {
	os.Remove(bashPath)
})

var _ = BeforeEach(func() {
	bash, _ = basher.NewContext(bashPath, true)

	stdout = gbytes.NewBuffer()
	stderr = gbytes.NewBuffer()
	bash.Stdout = io.MultiWriter(GinkgoWriter, stdout)
	bash.Stderr = io.MultiWriter(GinkgoWriter, stderr)
	bash.SelfPath = "/bin/echo"

	bash.CopyEnv()
})

func extractBash() {
	bashDir, err := homedir.Expand("~/.basher")
	if err != nil {
		log.Fatal(err, "1")
	}

	bashPath = bashDir + "/bash"
	if _, err := os.Stat(bashPath); os.IsNotExist(err) {
		err = basher.RestoreAsset(bashDir, "bash")
		if err != nil {
			log.Fatal(err, "1")
		}
	}
}

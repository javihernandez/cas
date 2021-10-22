package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"testing"
)

var (
	testApiKey1 = os.Getenv("CNIL_GITHUB_TEST_API_KEY1")
	testApiKey2 = os.Getenv("CNIL_GITHUB_TEST_API_KEY2")
	revokedKey  = os.Getenv("CNIL_GITHUB_REVOKED_KEY")
	testHost    = os.Getenv("CNIL_GITHUB_TEST_HOST")
	testPort    = os.Getenv("CNIL_GITHUB_TEST_PORT")
	signerID1   = os.Getenv("CNIL_SIGNERID1")
	// signerID2       string = os.Getenv("CNIL_SIGNERID2")
	revokedSigner = os.Getenv("REVOKED_SIGNERID")
	revokedHash   = os.Getenv("REVOKED_HASH")
	untrustedHash = os.Getenv("UNTRUSTED_HASH")
	//TODO: Create a list of images and make this part of the test configuration rather than a secret
	imageToNotarize = os.Getenv("IMAGE_TO_NOTARIZE")

	attributeRE             = regexp.MustCompile(".*test_attr_value1.*")
	loginSuccessRe          = regexp.MustCompile("Login successful")
	statusTrustedRE         = regexp.MustCompile("Status:.*TRUSTED")
	statusUntrustedRE       = regexp.MustCompile("Status:.*UNTRUSTED")
	statusRevokedRE         = regexp.MustCompile("Status:.*REVOKED")
	notarizationsFoundRE    = regexp.MustCompile("notarizations found")
	noSignerIdRE            = regexp.MustCompile("no signer ID provided")
	apiKeyRevokedRE         = regexp.MustCompile("Apikey revoked")
	recursiveNotarizationRE = regexp.MustCompile(`notarized.*\d+.*items`)
	badApikeyRE             = regexp.MustCompile("api key not valid")
	logoutSuccessfulRE      = regexp.MustCompile("Logout successful")
	logoutNotLoggedInRE     = regexp.MustCompile("No logged-in user")
	resolvingDependenciesRE = regexp.MustCompile("Resolving dependencies")
	needToLoginRE           = regexp.MustCompile("you need to be logged in")
	countersRE              = regexp.MustCompile("counters")

	cmdStr   = GetCmdFqp()
	basePath = GetBasePath()
)

func GetBasePath() string {
	d, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to read current directory")
	}
	return d
}

func FindBinaryName() string {
	osName := runtime.GOOS
	binaryName := "cas"
	switch osName {
	case "windows":
		return binaryName + ".exe"
	default:
		return binaryName
	}
}

func GetCmdFqp() string {
	binary := FindBinaryName()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory, aborting tests")
		os.Exit(1)
	}
	workingDir := path.Dir(currentDir)
	return path.Join(workingDir, binary)
}

type CasTest struct {
	// Name of the test
	name string
	// CLI Args
	args []string
	// Array of regular expressions to validate stdout/stderr
	xOutput    []*regexp.Regexp
	fixture    string
	goldenFile string
}

func RetrieveGoldenFileBytes(filename string) []byte {
	g, err := ioutil.ReadFile(basePath + "/work/golden_files/" + filename)
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to read golden file %s", filename)
	}
	return g
}

func CompareBytes(expected []byte, actual []byte) {
	if !bytes.Equal(expected, actual) {
		err := ioutil.WriteFile("./work/actual_out", actual, 0644)
		if err != nil {
			log.Fatal("Failed to write output to disk")

		}
		log.Printf("Got:\n%s", actual)
		log.Printf("Expected:\n%s", string(expected))
		log.Fatal("output does not match golden file")
	}

}

type C struct{}

func ExecuteTests(testArray []CasTest, t *testing.T) {
	for _, tt := range testArray {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fixture != "" {
				c := C{}
				f := reflect.ValueOf(c).MethodByName(tt.fixture)
				f.Call(nil)
			}
			cmd := exec.Command(cmdStr, tt.args...)
			out, err := cmd.CombinedOutput()
			_ = cmd.Run()
			if err != nil {
				log.Println("Logging stdErr, tests may still pass")
				log.Println(fmt.Sprint(err) + ": " + string(out))
			}
			actual := string(out)
			for _, regex := range tt.xOutput {
				if !regex.MatchString(actual) {
					t.Log(fmt.Sprintf("Test %s failed", tt.name))
					xString := fmt.Sprintf("Expected %s but got %s", tt.xOutput, actual)
					log.Fatal(xString)
				}
			}
			if tt.goldenFile != "" {
				expected := RetrieveGoldenFileBytes(tt.goldenFile)
				CompareBytes(expected, out)
			}
		})
	}
}

func (c C) FixtureSetEnvApiKey() {
	_, f := os.LookupEnv("CAS_API_KEY")
	if f {
		err := os.Setenv("CAS_API_KEY", testApiKey2)
		if err != nil {
			log.Fatal("Failed to set environment variable for API Key")
		}
	} else {
		err := os.Setenv("CAS_API_KEY", testApiKey1)
		if err != nil {
			log.Fatal("Failed to set environment variable for API Key")
		}
	}
	log.Printf("API Key set to %s", os.Getenv("CAS_API_KEY"))

}

func (c C) FixturePullImage() {
	cmd := exec.Command("docker", "pull", imageToNotarize)
	out, err := cmd.CombinedOutput()
	_ = cmd.Run()
	if err != nil {
		log.Printf("Error pulling image %s", imageToNotarize)
	}
	log.Println(string(out))
}

func (c C) FixtureBuildBomImage(tagname string, filename string) {
	cmd := exec.Command("docker", "build", "-t", tagname, "../", "-f", filename)
	out, err := cmd.CombinedOutput()
	_ = cmd.Run()
	if err != nil {
		log.Printf("Error building bom image %s", filename)
	}
	log.Println(string(out))
}

/** Tests authenticated interaction with CNIL backend**/
func TestCNcloudContext(t *testing.T) {
	var tests = []CasTest{
		{"Login to private test instance", []string{"login", "--host", testHost, "--port", testPort, "--api-key", testApiKey1}, []*regexp.Regexp{loginSuccessRe}, "FixtureSetEnvApiKey", ""},
		{"Notarize a simple file", []string{"n", "integration_test.go"}, []*regexp.Regexp{statusTrustedRE}, "", ""},
		{"Authenticate a previously notarized file", []string{"a", "integration_test.go"}, []*regexp.Regexp{statusTrustedRE}, "", ""},
		{"Untrust a previously notarized file", []string{"untrust", "integration_test.go"}, []*regexp.Regexp{statusUntrustedRE}, "", ""},
		{"Authenticate a previously untrusted file", []string{"a", "integration_test.go"}, []*regexp.Regexp{statusUntrustedRE}, "", ""},
		{"Inspect a previously untrusted file using a date range", []string{"i", "integration_test.go", "--start", "2021/08/25-00:00:00", "--end", "2021/08/25-23:59:00", "--first", "10"}, []*regexp.Regexp{notarizationsFoundRE}, "", ""},
		{"Inspect a previously untrusted file (No Args)", []string{"inspect", "integration_test.go"}, []*regexp.Regexp{noSignerIdRE}, "", ""},
		{"Inspect a previously untrusted file --first flag", []string{"inspect", "integration_test.go", "--first", "1"}, []*regexp.Regexp{notarizationsFoundRE}, "", ""},
		{"Inspect a previously untrusted file --last flag", []string{"inspect", "integration_test.go", "--last", "1"}, []*regexp.Regexp{notarizationsFoundRE}, "", ""},
		{"Inspect a previously untrusted file hash flag", []string{"inspect", "--hash", untrustedHash, "--signerID", signerID1}, []*regexp.Regexp{notarizationsFoundRE}, "", ""},
		{"Notarize a file with specific attributes", []string{"n", "--attr", "test_attr1=test_attr_value1", "integration_test.go"}, []*regexp.Regexp{attributeRE}, "", ""},
		{"Authenticate with previous signerId", []string{"a", "integration_test.go", "--signerID", signerID1}, []*regexp.Regexp{statusTrustedRE}, "FixtureSetEnvApiKey", ""},
		{"Authenticate with a revoked signerId", []string{"a", "--hash", revokedHash, "--signerID", revokedSigner}, []*regexp.Regexp{apiKeyRevokedRE, statusRevokedRE}, "", ""},
		{"Notarize docker image", []string{"n", fmt.Sprintf("docker://%s", imageToNotarize)}, []*regexp.Regexp{statusTrustedRE}, "FixturePullImage", ""},
		{"Notarize git repo ", []string{"n", "git://../"}, []*regexp.Regexp{statusTrustedRE}, "", ""},
		{"Attempt to notarize using revoked API Key", []string{"n", "git://../", "--api-key", revokedKey}, []*regexp.Regexp{badApikeyRE}, "", ""},
		{"Attempt to notarize using an invalid API Key", []string{"n", "git://../", "--api-key", "lc._"}, []*regexp.Regexp{badApikeyRE}, "", ""},
		{"Logout", []string{"logout"}, []*regexp.Regexp{logoutSuccessfulRE}, "", ""},
		{"Logout without logging in", []string{"logout"}, []*regexp.Regexp{logoutNotLoggedInRE}, "", ""},
		{"Attempt to notarize after logging out", []string{"n", "git://../", "--api-key", "lc._"}, []*regexp.Regexp{needToLoginRE}, "", ""},
	}
	err := os.Setenv("CAS_SKIP_SIGNATURE_VERIFY", "true")
	if err != nil {
		log.Fatal("Failed setting environment variable CAS_SKIP_SIGNATURE_VERIFY")
	}
	ExecuteTests(tests, t)
}

func TestBomUnsupported(t *testing.T) {
	// Negative tests still can't use the above illustrated method
	pckManErrorRE := *regexp.MustCompile(`.*cannot identify package manager.*`)
	c := C{}
	c.FixturePullImage()
	cmd := exec.Command(cmdStr, "bom", fmt.Sprintf("docker://%s", imageToNotarize))
	out, err := cmd.CombinedOutput()
	_ = cmd.Run()
	log.Println(err)
	if err != nil {
		// We want to be here
		if !pckManErrorRE.MatchString(string(out)) {
			t.Fatal("Package manager error not found")
		}
	} else {
		t.Fatal("Unsupported docker image generated no errors")
	}
}

func TestBomSupportedDockerImage(t *testing.T) {
	err := exec.Command("docker", "pull", "codenotary/cas:bom-nodejs").Run()
	if err != nil {
		log.Fatal("failed pulling node-js image ")
	}
	var tests = []CasTest{
		{"Run cas bom on a supported docker image", []string{"bom", "docker://codenotary/cas:bom-nodejs"}, []*regexp.Regexp{regexp.MustCompile(`.*ca-certificates.*`), regexp.MustCompile(`.*openssl.*`)}, "", ""},
	}
	ExecuteTests(tests, t)

}

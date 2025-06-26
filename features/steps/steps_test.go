package steps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	projectRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("erro ao obter diretório de trabalho: %v", err)
	}
	for !dirExists(filepath.Join(projectRoot, "features")) {
		projectRoot = filepath.Dir(projectRoot)
		if projectRoot == "/" {
			t.Fatalf("diretório 'features' não encontrado")
		}
	}
	featuresPath := filepath.Join(projectRoot, "features")

	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{featuresPath},
	}
	status := godog.TestSuite{
		Name:                "godog",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()
	if status != 0 {
		t.Fail()
	}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

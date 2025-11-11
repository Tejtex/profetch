package pkgmgr

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Npm struct{}

func (n Npm) Detect(root string) (string, bool) {
    _, err := os.Stat(filepath.Join(root, "package.json"))
    return filepath.Join(root, "package.json"), err == nil

}

func (n Npm) CountDeps(path string) int {
	data, err := os.ReadFile(path)
    if err != nil {
        return 0
    }
    var pkg map[string]interface{}
    if err := json.Unmarshal(data, &pkg); err != nil {
        return 0
    }
    total := 0
    if deps, ok := pkg["dependencies"].(map[string]interface{}); ok {
        total += len(deps)
    }
    if deps, ok := pkg["devDependencies"].(map[string]interface{}); ok {
        total += len(deps)
    }
    return total
}
func (n Npm) GetName() string {
	return "npm"
}
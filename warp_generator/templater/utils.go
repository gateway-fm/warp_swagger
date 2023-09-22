package templater

import (
	"fmt"
	"os"
)

func GetTemplateInterfaces(iface ...any) (ifaces []any) {
	ifaces = append(ifaces, iface...)
	return ifaces
}
func CompleteFuncMap(names []string, a []any) map[string]any {
	m := make(map[string]any)
	for i := range names {
		m[names[i]] = a[i]
	}
	return m
}

func TempDir() error {
	if err := os.MkdirAll("./templates", 0777); err != nil {
		return fmt.Errorf("failed to mkdir:%w", err)
	}
	return nil
}
func CopyTemplatesToTemp(s string) error {

	pathCreate := "./templates" + "/" + s

	_, err := os.Create(pathCreate)
	if err != nil {
		return fmt.Errorf("failed to exec CopyTemplatesToTemp %w", err)
	}
	return nil
}

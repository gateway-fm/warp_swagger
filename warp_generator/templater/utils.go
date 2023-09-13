package templater

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

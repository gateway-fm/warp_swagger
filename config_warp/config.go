package config_warp //nolint:revive

type Warp struct {
	External *ExternalPkg
}
type ExternalPkg struct {
	Models      []*Models
	Output      string
	PackageName string
	ProtoName   string
	ProtoPath   string
	PackageURL  string
}

type Models struct {
	Name string
	Type string //   optional in the most cases
}

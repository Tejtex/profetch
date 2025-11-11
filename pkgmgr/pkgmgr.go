package pkgmgr

type PackageManager interface {
	Detect(root string) (string, bool);
	CountDeps(path string) int;
	GetName() string;
}
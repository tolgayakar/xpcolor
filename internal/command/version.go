package command

import "runtime/debug"

// version is set at build time via -ldflags.
var version string

func Version() string {
	if version != "" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(unset)"
	}
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	var revision, dirty string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			if s.Value == "true" {
				dirty = " (dirty)"
			}
		}
	}
	if revision == "" {
		return "(unset)"
	}
	return revision + dirty
}

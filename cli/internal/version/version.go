package version

// Version is the CLI version, injected at build time via -ldflags.
// Default to "dev" for unreleased local builds.
var Version = "dev"

package main

type Config struct {
	remoteBootstrapTarball string
	localBootstrapTarball string
	// Path on host where the chroot starts
	localDir string
	// Name of user created within chroot
	userName string
	repos    []Repo
}

type Repo struct {
	// Git remote that will be cloned
	remoteUrl string
}

// hard-coded
var defaultConfig = Config{
	// See https://archlinux.org/download/
	remoteBootstrapTarball: "http://mirrors.ocf.berkeley.edu/archlinux/iso/2022.11.01/archlinux-bootstrap-x86_64.tar.gz",
	localBootstrapTarball: "archlinux-bootstrap.tar.gz",
	userName:               "coder",
	// This is encoded in the tarball, double-check
	localDir: "root.x86_64",
	repos: []Repo{
		{
			remoteUrl: "git@github.com:christopherfujino/chris-monorepo",
		},
		{
			remoteUrl: "git@github.com:christopherfujino/dotfiles",
		},
	},
}

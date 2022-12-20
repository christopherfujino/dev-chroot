package main

type Config struct {
	RemoteBootstrapTarball string
	LocalBootstrapTarball  string
	//// Path on host where the chroot starts
	//LocalDir string
	// Name of user created within chroot
	UserName string
	UID      int
	Repos    []Repo
	// Contents of hash bang script run in chroot as root
	Provision string
}

type Repo struct {
	// Git remote that will be cloned
	remoteUrl string
}

// hard-coded
var defaultConfig = Config{
	// See https://archlinux.org/download/
	RemoteBootstrapTarball: "http://mirrors.ocf.berkeley.edu/archlinux/iso/2022.11.01/archlinux-bootstrap-x86_64.tar.gz",
	LocalBootstrapTarball:  "archlinux-bootstrap.tar.gz",
	UserName:               "coder",
	UID:                    1000,
	Repos: []Repo{
		{
			remoteUrl: "git@github.com:christopherfujino/chris-monorepo",
		},
		{
			remoteUrl: "git@github.com:christopherfujino/dotfiles",
		},
	},
	Provision: `
#!/usr/bin/env bash

# setup public keyring for pacman
pacman-key --init

# verifying the master keys
pacman-key --populate

pacman -Syu

# openssh needed to git clone via ssh
# unzip is needed by Flutter
# cmake is needed to build Neovim
# --needed means do not reinstall already present packages
pacman -S --needed \
	base-devel \
	man-db \
	vim \
	git \
	openssh \
	tmux \
	sudo \
	unzip \
	cmake \
	ripgrep \
	fzf

echo "Creating user {{.UserName}}..."
useradd --create-home --uid {{.UID}} --shell /bin/bash {{.UserName}}
passwd coder
echo "{{.UserName}} ALL=(ALL:ALL)" >> /etc/sudoers
`,
}

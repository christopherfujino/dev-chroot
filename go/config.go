package main

type Config struct {
	RemoteBootstrapTarball string
	LocalBootstrapTarball  string
	// Name of user created within chroot
	UserName string
	UID      int
	Repos    []Repo
	// Contents of a hash bang script run once in chroot as root
	Provision string
}

type Repo struct {
	// Git remote that will be cloned
	remoteUrl string
}

// hard-coded
var defaultConfig = Config{
	// See https://archlinux.org/download/
	RemoteBootstrapTarball: "https://mirrors.ocf.berkeley.edu/archlinux/iso/2023.04.01/archlinux-bootstrap-2023.04.01-x86_64.tar.gz",
	LocalBootstrapTarball:  "archlinux-bootstrap.tar.gz",
	UserName:               "coder",
	//UID:                    1000,
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

set -x

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
# Don't create home, we will git clone the home dir
useradd --uid {{.UID}} --shell /bin/bash {{.UserName}}
mkdir /home/{{.UserName}}
chown {{.UserName}}:{{.UserName}} /home/{{.UserName}}
sudo -u {{.UserName}} git clone --progress --verbose ssh://github@github.com/christopherfujino/chris-monorepo /home/{{.UserName}}
passwd coder
echo "{{.UserName}} ALL=(ALL:ALL) ALL" >> /etc/sudoers
`,
}

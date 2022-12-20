package main

type Config struct {
	remoteBootstrapTarball string
	localBootstrapTarball string
	// Path on host where the chroot starts
	localDir string
	// Name of user created within chroot
	userName string
	repos    []Repo
	// Contents of a /bin/sh script run in chroot as root
	Provision string
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
	repos: []Repo{
		{
			remoteUrl: "git@github.com:christopherfujino/chris-monorepo",
		},
		{
			remoteUrl: "git@github.com:christopherfujino/dotfiles",
		},
	},
	Provision: `
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

echo "Creating user coder..."
useradd -m -s /bin/bash coder
passwd coder
echo "coder ALL=(ALL:ALL)" >> /etc/sudoers
`,
}

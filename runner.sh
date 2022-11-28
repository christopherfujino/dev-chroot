#!/usr/bin/env bash

set -euo pipefail
YEAR='2022'
MONTH='11'
DATE='01'
ARCHITECTURE='x86_64'
# See https://archlinux.org/download/
REMOTE_BOOTSTRAP_TARBALL="http://mirrors.ocf.berkeley.edu/archlinux/iso/$YEAR.$MONTH.$DATE/archlinux-bootstrap-$ARCHITECTURE.tar.gz"
USERNAME='coder'

# This is encoded in the tarball, double-check
LOCAL_DIR="root.$ARCHITECTURE"

MONOREPO_REMOTE='git@github.com:christopherfujino/chris-monorepo'
DOTFILES='git@github.com:christopherfujino/dotfiles'

RUNNER='runner.sh'

function usage_exit {
  echo 'Usage: runner.sh [function]' >&2
  echo 'where [function] is one of:' >&2
  printf "\tbootstrap;\t\tdownload archlinux bootstrap tarball\n" >&2
  printf "\tattach;\t\t\tchroot to local arch linux dir\n" >&2
  printf "\tinitialize-root;\tinitialize chroot as root\n" >&2
  printf "\tinitialize-user;\tinitialize the user account\n" >&2
  exit 1
}

# run as user on host
function bootstrap {
  set -o xtrace

  if [ "$EUID" -eq 0 ]; then
    echo 'bootstrap function should not be called by root' >&2
    exit 1
  fi

  LOCAL_TARBALL='arch-tarball.tar.gz'

  if [ ! -f "$LOCAL_TARBALL" ]; then
    curl -l "$REMOTE_BOOTSTRAP_TARBALL" -o "$LOCAL_TARBALL"
  fi

  if [ ! -d "$LOCAL_DIR" ]; then
    # --numeric-owner since host might not use the same user id's as arch
    # Must have root permission as there are some UID 0 files
    sudo tar xzf "$LOCAL_TARBALL" --numeric-owner

    # Enable berkeley mirror
    # -E means extended regex
    # -i means update file in place
    sudo sed -E -i 's/^#(.*berkeley)/\1/' "$LOCAL_DIR/etc/pacman.d/mirrorlist"
    # disable CheckSpace setting
    sudo sed -E -i 's/^CheckSpace/#CheckSpace/' "$LOCAL_DIR/etc/pacman.conf"
  fi
}

# run as user on host
function attach {
  if [ "$EUID" -eq 0 ]; then
    echo 'attach function should not be called by root' >&2
    exit 1
  fi

  set -o xtrace

  # -ot is older than
  if [ ! -f "$LOCAL_DIR/$RUNNER" ] || [ "$LOCAL_DIR/$RUNNER" -ot "$RUNNER" ]; then
    sudo cp "$RUNNER" "$LOCAL_DIR/$RUNNER"
  fi

  if [ ! -f "$HOME/.ssh/id_rsa" ]; then
    echo "Expected $HOME/.ssh/id_rsa to exist to copy to chroot" >&2
    exit 1
  fi

  if [ ! -d "$LOCAL_DIR/.ssh" ]; then
    sudo cp -r "$HOME/.ssh/" "$LOCAL_DIR"
  fi

  # chroot
  sudo "$LOCAL_DIR/bin/arch-chroot" "$LOCAL_DIR/" "/$RUNNER" 'initialize-root'
}

# should be run as root in chroot
function initialize-root {
  if [ "$EUID" -ne 0 ]; then
    echo 'initialize-root function should be called by root' >&2
    exit 1
  fi

  set -o xtrace

  TOUCHFILE='/.initialized_root'
  if [ ! -f "$TOUCHFILE" ] || [ "$TOUCHFILE" -ot "/$RUNNER" ]; then
    echo "$TOUCHFILE cache miss"
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

    if ! id -u "$USERNAME" >/dev/null; then
      echo "creating user $USERNAME..."
      useradd -m -s /bin/bash "$USERNAME"

      passwd "$USERNAME"

      # Don't worry about security within chroot
      echo "$USERNAME ALL=(ALL:ALL) NOPASSWD: ALL" >> /etc/sudoers
    fi

    touch "$TOUCHFILE"
  fi

  # re-entrant call as user
  sudo -u "$USERNAME" "/$RUNNER" 'initialize-user'

  su --login "$USERNAME"
}

# run as root in chroot
function initialize-user {
  set -o xtrace

  if [ "$EUID" -eq 0 ]; then
    echo 'initialize-user function should not be called by root' >&2
    exit 1
  fi

  cd ~

  TOUCHFILE="$HOME/.initialized_user"

  if [ ! -f "$TOUCHFILE" ] || [ "$TOUCHFILE" -ot "/$RUNNER" ]; then
    if [ ! -d "$HOME/.ssh" ]; then
      sudo cp -r /.ssh "$HOME"
      sudo chown --recursive $(id -u):$(id -g) "$HOME/.ssh"
    fi

    if [ ! -d "$HOME/git/chris-monorepo" ]; then
      git clone "$MONOREPO_REMOTE" "$HOME/git/chris-monorepo"
    fi
    if [ ! -d "$HOME/git/dotfiles" ]; then
      git clone "$DOTFILES" "$HOME/git/dotfiles"
      # TODO this should just be repo init scripts
      ln -s -f "$HOME/dotfiles/.bashrc" "$HOME/.bashrc"
    fi

    touch "$TOUCHFILE"
  fi
}

if [ $# -ne 1 ]; then
  usage_exit
fi

case "$1" in
  'bootstrap')
    bootstrap
    ;;
  'attach')
    attach
    ;;
  'initialize-root')
    initialize-root
    ;;
  'initialize-user')
    initialize-user
    ;;
  *)
    usage_exit
    ;;
esac

exit 0

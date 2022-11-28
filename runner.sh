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

function usage_exit {
  echo 'Usage: runner.sh [function]' >&2
  echo 'where [function] is one of:' >&2
  printf "\tbootstrap;\t\tdownload archlinux bootstrap tarball\n" >&2
  printf "\tattach;\t\t\tchroot to local arch linux dir\n" >&2
  printf "\tinitialize-root;\tinitialize chroot as root\n" >&2
  printf "\tinitialize-user;\tinitialize the user account\n" >&2
  exit 1
}

function bootstrap {
  set -o xtrace

  LOCAL_TARBALL='arch-tarball.tar.gz'

  curl -l "$REMOTE_BOOTSTRAP_TARBALL" -o "$LOCAL_TARBALL"

  # --numeric-owner since host might not use the same user id's as arch
  sudo tar xzf "$LOCAL_TARBALL" --numeric-owner

  # Enable berkeley mirror
  # -E means extended regex
  # -i means update file in place
  sudo sed -E -i 's/^#(.*berkeley)/\1/' "$LOCAL_DIR/etc/pacman.d/mirrorlist"
  # disable CheckSpace setting
  sudo sed -E -i 's/^CheckSpace/#CheckSpace/' "$LOCAL_DIR/etc/pacman.conf"
}

function attach {
  set -o xtrace

  RUNNER='runner.sh'
  # -ot is older than
  if [ ! -f "$LOCAL_DIR/$RUNNER" ] || [ "$LOCAL_DIR/$RUNNER" -ot "$RUNNER" ]; then
    sudo cp "$RUNNER" "$LOCAL_DIR/$RUNNER"
  fi

  # chroot
  sudo "$LOCAL_DIR/bin/arch-chroot" "$LOCAL_DIR/" "/$RUNNER" 'initialize-root'
}

# should be run as root
function initialize-root {
  # TODO test userid is 0
  # TODO test if we already did this

  set -o xtrace

  # setup public keyring for pacman
  pacman-key --init

  # verifying the master keys
  pacman-key --populate

  pacman -Syu

  # openssh needed to git clone via ssh
  pacman -S \
    base-devel \
    man-db \
    vim \
    git \
    openssh \
    tmux \
    sudo

  if ! id -u "$USERNAME"; then
    echo "creating $USERNAME user..."
    useradd -m -s /bin/bash "$USERNAME"

    passwd "$USERNAME"
  else
    echo "User $USERNAME already exists"
  fi

  su --login "$USERNAME"
}

function initialize-user {
  set -o xtrace

  # TODO check not userid 0
  ssh-keygen
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
    echo "Oops $1 is not yet implemented" >&2
    exit 1
    initialize-user
    ;;
  *)
    usage_exit
    ;;
esac

exit 0

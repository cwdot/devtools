#!/usr/bin/env zsh

function install() {
  pushd .
  D=$1
  cd $D || exit
  go install
  echo "Installed $D"
  popd || exit

  ~/go/bin/$D version

  echo ""
  echo ""
}

install 1px

install hass

install gitter



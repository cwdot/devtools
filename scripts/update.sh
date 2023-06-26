#!/usr/bin/env zsh

function update() {
  pushd .
  D=$1
  cd ../$D || exit

  go get -u

  echo ""
  echo ""
}

update 1px

update hass

update gitter

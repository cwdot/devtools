#!/usr/bin/env zsh

function install() {
  pushd .
  D=$1
  cd ../$D || exit

  echo "[$D] installing..."
  echo "[$D] go mod tidy"
  go mod tidy

  echo "[$D] go install"
  go install
  echo "Installed $D"
  popd || exit

  echo "[$D] version"
  ~/go/bin/$D version

  if [[ "$(uname)" == "Darwin" ]]
  then
    output=$(brew --prefix)/share/zsh/site-functions/_$D
    ~/go/bin/"$D" completion zsh > $output
    echo "Created completion script for $D in $output"
  elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    output="${fpath[1]}/_$D"
    ~/go/bin/"$D" completion zsh > $output
    echo "Created completion script for $D in $output"
  fi

  echo ""
  echo ""
}

git pull --rebase

install 1px

install hass

install gitter

#!/usr/bin/env zsh

function install() {
  pushd .
  D=$1
  cd ../$D || exit

  go mod tidy

  go install
  echo "Installed $D"
  popd || exit

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

install 1px

install hass

install gitter



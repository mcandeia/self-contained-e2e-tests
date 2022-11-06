#!/bin/bash

brew install kind
brew install docker
brew install helm
brew install kubectl
brew install tilt

curl -L -o vcluster "https://github.com/loft-sh/vcluster/releases/latest/download/vcluster-darwin-arm64" && sudo install -c -m 0755 vcluster /usr/local/bin
go install github.com/google/ko@latest

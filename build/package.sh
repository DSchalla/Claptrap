#!/bin/bash
set -x
set -e
cd ../cmd/claptrap
go build
mv claptrap ../../build
cd ../../build
tar cvfz claptrap-plugin.tar.gz plugin.yaml claptrap
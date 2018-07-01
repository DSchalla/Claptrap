#!/bin/bash
set -x
set -e
cd ../cmd/claptrap
go build
mv claptrap ../../build
cd ../../build
cp -r "../static" .
tar cvfz claptrap-plugin.tar.gz plugin.yaml claptrap static/
rm -rf ./static

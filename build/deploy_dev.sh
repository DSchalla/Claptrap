#!/bin/bash
set -e
set -x
./package.sh;
rm ../../../mattermost/mattermost-server/plugins/com.dschalla.claptrap/claptrap;
cp claptrap ../../../mattermost/mattermost-server/plugins/com.dschalla.claptrap/claptrap;
rm -rf ../../../mattermost/mattermost-server/plugins/com.dschalla.claptrap/static;
cp -r ../static ../../../mattermost/mattermost-server/plugins/com.dschalla.claptrap/;
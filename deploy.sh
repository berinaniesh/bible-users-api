#!/bin/bash

go build .
scp bible-users-api berinaniesh.xyz:/home/berinaniesh/tmp/
ssh berinaniesh.xyz deploy-bible-users-api.sh

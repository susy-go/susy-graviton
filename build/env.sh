#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
sofdir="$workspace/src/github.com/susy-go"
if [ ! -L "$sofdir/susy-graviton" ]; then
    mkdir -p "$sofdir"
    cd "$sofdir"
    ln -s ../../../../../. susy-graviton
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$sofdir/susy-graviton"
PWD="$sofdir/susy-graviton"

# Launch the arguments with the configured environment.
exec "$@"

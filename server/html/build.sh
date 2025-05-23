#!/bin/bash

# Builds the template styles with Tailwind CSS's standalone CLI.
# https://tailwindcss.com/blog/standalone-cli

set -e
cd "$(dirname ${BASH_SOURCE[0]})"

tw="$(which tailwindcss)"
if [[ "$tw" ]]; then
    # Tailwind exists globally
    echo "Using Tailwind: $tw"
elif [[ ! -f ./tailwindcss ]]; then
    # TODO: accept that Windows still probably exists
    if [[ "$(uname -s)" == "Darwin" ]]; then
        os="macos"
    else
        os="linux"
    fi
    arch="$(uname -m)"
    if [[ "$arch" == "x86_64" ]]; then
        arch="x64"
    fi
    echo "Downloading tailwindcss-$os-$arch"
    curl -sLO "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$os-$arch"
    mv "tailwindcss-$os-$arch" tailwindcss

    chmod +x tailwindcss
    tw="./tailwindcss"
else
    tw="./tailwindcss"
fi

echo 'Starting build...'
"$tw" -i app.css -o app.min.css --minify "$@"

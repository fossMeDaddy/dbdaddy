#!/bin/sh

# RUNS STANDALONE

pkg_bin_dir="$HOME/.dbdaddy/bin"
pkg_bin="$pkg_bin_dir/dbdaddy"

pkg_version="$(curl --silent https://raw.githubusercontent.com/fossMeDaddy/dbdaddy/main/version)"

mkdir $HOME/.dbdaddy
mkdir $HOME/.dbdaddy/bin

# get os & arch here

# get os & arch specific release
# exc: aarch64 can appear in uname (map to arm64)
# exc: x86_64 can appear in uname (map to amd64)

curl https://github.com/fossMeDaddy/dbdaddy/releases/download/$pkg_version/dbdaddy -L -o $pkg_bin
chmod +x $pkg_bin

for file in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile" "$HOME/.zshenv"; do
    if [ -f "$file" ]; then
        if grep -q "alias dbdaddy=" "$file"; then
            echo "Alias 'dbdaddy' already exists in $file. Skipping."
        else
            echo "\nalias dbdaddy=$pkg_bin" >> "$file"
            echo "Appended alias to $file"
        fi
    fi
done

echo "Please restart terminal to run 'dbdaddy' globally."
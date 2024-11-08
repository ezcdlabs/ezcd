#!/bin/sh

# Set up Git aliases
git config --global alias.cm 'commit -m'
git config --global alias.cme 'commit --allow-empty -m'
git config --global alias.st 'status --short --branch'
git config --global alias.lg "log --graph --abbrev-commit --decorate --format=format:'%C(yellow)%h%C(reset)%C(auto)%d%C(reset) %C(normal)%s%C(reset) %C(dim blue)(%an - %ar - %ad)%C(reset)'"
git config --global alias.lgu "log --abbrev-commit --decorate --format=format:'%C(yellow)%h%C(reset)%C(auto)%d%C(reset) %C(normal)%s%C(reset) %C(dim blue)(%an - %ar - %ad)%C(reset)' -u"
git config --global alias.lgs "log --graph --abbrev-commit --decorate --stat --format=format:'%C(yellow)%h%C(reset)%C(auto)%d%C(reset) %C(normal)%s%C(reset) %C(dim blue)(%an - %ar - %ad)%C(reset)%n'"
git config --global alias.lgn "log --graph --abbrev-commit --decorate --numstat --format=format:'%C(yellow)%h%C(reset)%C(auto)%d%C(reset) %C(normal)%s%C(reset) %C(dim blue)(%an - %ar - %ad)%C(reset)%n'"
git config --global alias.lga "log --graph --abbrev-commit --decorate --format=format:'%C(yellow)%h%C(reset)%C(auto)%d%C(reset) %C(normal)%s%C(reset) %C(dim blue)(%an - %ar - %ad)%C(reset)' --all"

# Set git to use rebase
git config --global pull.rebase true
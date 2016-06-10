rm -fr ~/git

mkdir ~/git

# a clean repo with tags
cd ~/git
git init
git config user.email "john@doe.com"
git config user.name "John Doe"
touch tomate
git add -A
git commit -m "re v1"
git tag "notsemvertag"
git tag "v1.0.2"
git tag "v1.0.0"

git tag

# a dirty repo
rm -fr ~/git_dirty
mkdir ~/git_dirty
cd ~/git_dirty
git init
touch mew
git add -A

# a repo with untracked files
rm -fr ~/git_untracked
mkdir ~/git_untracked
cd ~/git_untracked
git init
touch mew2

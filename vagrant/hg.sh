cat <<EOT > ~/.hgrc
[ui]
username = John Doe <john@example.com>
EOT

rm -fr ~/hg
mkdir ~/hg

cd ~/hg
hg init
touch tomate
hg add
hg commit --m "toto"
hg tag notsemvertag
hg tag v1.0.2
hg tag v1.0.0
hg tags

# a dirty repo
rm -fr ~/hg_dirty
mkdir ~/hg_dirty
cd ~/hg_dirty
hg init
touch mew
hg add

# a repo with untracked files
rm -fr ~/hg_untracked
mkdir ~/hg_untracked
cd ~/hg_untracked
hg init
touch mew2

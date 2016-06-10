rm -fr ~/bzr

mkdir ~/bzr

# a clean repo with tags
cd ~/bzr
bzr whoami "Your Name <name@example.com>"
bzr init
touch tomate
bzr add *
bzr commit -m "re v1"
bzr tag "notsemvertag"
bzr tag "v1.0.2"
bzr tag "v1.0.0"

bzr tags --sort=time

# a dirty repo
rm -fr ~/bzr_dirty
mkdir ~/bzr_dirty
cd ~/bzr_dirty
bzr whoami "Your Name <name@example.com>"
bzr init
touch mew
bzr add *

# a repo with untracked files
rm -fr ~/bzr_untracked
mkdir ~/bzr_untracked
cd ~/bzr_untracked
bzr whoami "Your Name <name@example.com>"
bzr init
touch mew2

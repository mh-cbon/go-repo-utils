
echo ""
echo "################"
echo "bzr"

rm -fr ~/bzr

mkdir ~/bzr

# a clean repo with tags
cd ~/bzr
bzr whoami "Your Name <name@example.com>"
bzr init
touch tomate-notsemvertag
bzr add *
bzr commit -m "tomate notsemvertag"
bzr tag "notsemvertag"
touch tomate-1.0.2
bzr add *
bzr commit -m "tomate 1.0.2"
bzr tag "v1.0.2"
sleep 1 # need to ensure that at least one commit is not done within same second to test ordering
touch tomate-1.0.0
bzr add *
bzr commit -m "tomate 1.0.0"
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

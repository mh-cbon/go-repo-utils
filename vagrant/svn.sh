
echo ""
echo "################"
echo "svn"

rm -fr ~/svn
rm -fr ~/svn_work
rm -fr ~/svnrep

mkdir ~/svn
mkdir ~/svn_work
mkdir ~/svnrep
cd ~/svnrep
svnadmin create svn

svn import ~/svn file:///home/vagrant/svnrep/svn/trunk -m "Initial import of project1"
svn co file:///home/vagrant/svnrep/svn/trunk /home/vagrant/svn_work
svn mkdir file:///home/vagrant/svnrep/svn/tags/ -m "Add tag folder"

cd /home/vagrant/svn_work/
touch tomate-notsemvertag
svn add tomate-notsemvertag
svn commit -m "tomate notsemvertag"
svn copy file:///home/vagrant/svnrep/svn/trunk file:///home/vagrant/svnrep/svn/tags/notsemvertag -m "Release notsemvertag"
touch tomate-1.0.2
svn add tomate-1.0.2
svn commit -m "tomate 1.0.2"
svn copy file:///home/vagrant/svnrep/svn/trunk file:///home/vagrant/svnrep/svn/tags/v1.0.2 -m "Release v1.0.2"
sleep 1 # need to ensure that at least one commit is not done within same second to test ordering
touch tomate-1.0.0
svn add tomate-1.0.0
svn commit -m "tomate 1.0.0"
svn copy file:///home/vagrant/svnrep/svn/trunk file:///home/vagrant/svnrep/svn/tags/v1.0.0 -m "Release v1.0.0"

svn update
svn ls -v ^/tags


# a dirty repo
rm -fr ~/svn_dirty
rm -fr ~/svn_dirty_work
mkdir ~/svn_dirty
mkdir ~/svn_dirty_work
cd ~/svnrep
svnadmin create svn_dirty

svn import ~/svn_dirty file:///home/vagrant/svnrep/svn_dirty/trunk -m "Initial import of project1"
svn co file:///home/vagrant/svnrep/svn_dirty/trunk /home/vagrant/svn_dirty_work

cd /home/vagrant/svn_dirty_work/
touch tomate
svn add tomate

# a repo with untracked files
rm -fr ~/svn_untracked
rm -fr ~/svn_untracked_work
mkdir ~/svn_untracked
mkdir ~/svn_untracked_work
cd ~/svnrep
svnadmin create svn_untracked

svn import ~/svn_untracked file:///home/vagrant/svnrep/svn_untracked/trunk -m "Initial import of project1"
svn co file:///home/vagrant/svnrep/svn_untracked/trunk /home/vagrant/svn_untracked_work

cd /home/vagrant/svn_untracked_work/
touch tomate

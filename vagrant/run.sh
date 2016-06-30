
go build -o build/go-repo-utils main.go
vagrant ssh -c "sh /vagrant/vagrant/git.sh && sh /vagrant/vagrant/bzr.sh && sh /vagrant/vagrant/hg.sh && sh /vagrant/vagrant/svn.sh"
echo ""
echo "################"
echo "test"
vagrant ssh -c "cd ~/gow/src/github.com/mh-cbon/go-repo-utils && GO15VENDOREXPERIMENT=1 GOPATH=/home/vagrant/gow GOROOT=/home/vagrant/go ~/go/bin/go test"

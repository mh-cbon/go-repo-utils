DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/..
go build -o build/go-repo-utils main.go
vagrant init ubuntu/trusty64 || echo "has already init"
vagrant up
vagrant ssh -c "sh /vagrant/vagrant/setup.sh"
vagrant ssh -c "sh /vagrant/vagrant/git.sh"
vagrant ssh -c "sh /vagrant/vagrant/bzr.sh"
vagrant ssh -c "sh /vagrant/vagrant/hg.sh"
vagrant ssh -c "sh /vagrant/vagrant/svn.sh"
vagrant ssh -c "cd ~/gow/src/github.com/mh-cbon/go-repo-utils && GO15VENDOREXPERIMENT=1 GOPATH=/home/vagrant/gow GOROOT=/home/vagrant/go ~/go/bin/go test"

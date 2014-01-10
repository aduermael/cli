# This file describes the standard way to build Docker, using docker
#
# Usage:
#
# # Assemble the full dev environment. This is slow the first time.
# docker build -t docker .
#
# # Mount your source in an interactive container for quick testing:
# docker run -v `pwd`:/go/src/github.com/dotcloud/docker -privileged -i -t docker bash
#
# # Run the test suite:
# docker run -privileged docker hack/make.sh test
#
# # Publish a release:
# docker run -privileged \
#  -e AWS_S3_BUCKET=baz \
#  -e AWS_ACCESS_KEY=foo \
#  -e AWS_SECRET_KEY=bar \
#  -e GPG_PASSPHRASE=gloubiboulga \
#  docker hack/release.sh
#
# Note: Apparmor used to mess with privileged mode, but this is no longer
# the case. Therefore, you don't have to disable it anymore.
#

docker-version	0.6.1
FROM	stackbrew/ubuntu:12.04
MAINTAINER	Tianon Gravi <admwiggin@gmail.com> (@tianon)

# Add precise-backports to get s3cmd >= 1.1.0 (so we get ENV variable support in our .s3cfg)
RUN	echo 'deb http://archive.ubuntu.com/ubuntu precise-backports main universe' > /etc/apt/sources.list.d/backports.list

# Packaged dependencies
RUN	apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
	apt-utils \
	aufs-tools \
	build-essential \
	curl \
	dpkg-sig \
	git \
	iptables \
	libsqlite3-dev \
	lxc \
	mercurial \
	reprepro \
	ruby1.9.1 \
	ruby1.9.1-dev \
	s3cmd=1.1.0* \
	--no-install-recommends

# Get lvm2 source for compiling statically
RUN	git clone https://git.fedorahosted.org/git/lvm2.git /usr/local/lvm2 && cd /usr/local/lvm2 && git checkout -q v2_02_103
# see https://git.fedorahosted.org/cgit/lvm2.git/refs/tags for release tags
# note: we can't use "git clone -b" above because it requires at least git 1.7.10 to be able to use that on a tag instead of a branch and we only have 1.7.9.5

# Compile and install lvm2
RUN	cd /usr/local/lvm2 && ./configure --enable-static_link && make device-mapper && make install_device-mapper
# see https://git.fedorahosted.org/cgit/lvm2.git/tree/INSTALL

# Install Go
RUN	curl -s https://go.googlecode.com/files/go1.2.src.tar.gz | tar -v -C /usr/local -xz
ENV	PATH	/usr/local/go/bin:$PATH
ENV	GOPATH	/go:/go/src/github.com/dotcloud/docker/vendor
RUN	cd /usr/local/go/src && ./make.bash --no-clean 2>&1

# Compile Go for cross compilation
ENV	DOCKER_CROSSPLATFORMS	darwin/amd64 darwin/386
# TODO add linux/386 and linux/arm
RUN	cd /usr/local/go/src && bash -xc 'for platform in $DOCKER_CROSSPLATFORMS; do GOOS=${platform%/*} GOARCH=${platform##*/} ./make.bash --no-clean 2>&1; done'

# Grab Go's cover tool for dead-simple code coverage testing
RUN	go get code.google.com/p/go.tools/cmd/cover

# TODO replace FPM with some very minimal debhelper stuff
RUN	gem install --no-rdoc --no-ri fpm --version 1.0.2

# Setup s3cmd config
RUN	/bin/echo -e '[default]\naccess_key=$AWS_ACCESS_KEY\nsecret_key=$AWS_SECRET_KEY' > /.s3cfg

VOLUME	/var/lib/docker
WORKDIR	/go/src/github.com/dotcloud/docker

# Wrap all commands in the "docker-in-docker" script to allow nested containers
ENTRYPOINT	["hack/dind"]

# Upload docker source
ADD	.	/go/src/github.com/dotcloud/docker

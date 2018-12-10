# Copyright 2018 The bolaxy.org Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

APPNAME = apiServer
DATEVERSION=$(shell date -u +%Y%m%d)
COMMIT_SHA=$(shell git rev-parse --short HEAD)
COMMIT_BRANCH=${shell git rev-parse --abbrev-ref HEAD}
COMMIT_TAG=${shell git describe --abbrev=0 --tags}
LDFLAGS="-X main.VERSION=${DATEVERSION} -X main.GITCOMMIT=${COMMIT_SHA}"

.PHONY:all
all:build

.PHONY:build
build:generate
	-rm -rf build
	-mkdir -p build/{,log}
	-mkdir -p build/versionlog
	-mkdir -p build/static/
	go generate .
	go build -o build/${APPNAME} -ldflags $(LDFLAGS)
	cp -r config/ build/config
	rm -rf ./build/config/*.go ./build/config/voucher.key
	cp log.xml build/
	cp versionlog/log.json build/versionlog/
	cp -r scripts build/scripts
	cp -r static/dist build/static/
	cp -r static/lang build/static/
	cp glide.yaml build/
	cp path.yml build/

.PHONY:generate
generate:
	go generate .

.PHONY:rebuild
rebuild:
	go build -o build/${APPNAME} -ldflags $(LDFLAGS)

.PHONY:install
install:build
	mv ${APPNAME} ${GOPATH}/

.PHONY:clean
clean:
	-rm -rf build
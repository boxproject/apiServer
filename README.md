# apiServer

README: [ENGLISGH](./README.md) | [中文](./docs/readme_zh.md)

## Overview

BOX (Enterprise Token Safe Box) is an enterprise-level digital assets safe application that uses
the axiomatic techniques in blockchain, cryptography and communications security to protect
private keys and instructions. BOX, in principle, seeks to prevent the theft and tamper of private
keys and instructions. More Information to see [BOX Official Website](https://box.la/en) and [BOX White Paper](https://box.la/static/BOX_white_paper_en.pdf)。

## Table of Contents

+ [Features](#Features)
+ [Installation](#Installation)
+ [Quick Start](#QuickStart)
+ [Contributing](#Contributing)
+ [License](#License)

### Features

+ Automatic and secure transfer of assets according to approval flow along with bookkeeping function

+ Self-managed cryptographic keys and shared permissions

+ Multi-Layers of security for private key, transaction instruction and communication

+ Compliance with multiple token standards, one private key controls multiple chains

+ Free and Open Source

### Installation

+ Make sure you install the [prerequisites](./docs/requirements_en.md) first.
+ Generate certificate：

~~~sh
$ sh scripts/keys.sh
~~~

+ Modify the configuration file `config.toml.example` to improve proxy server information and your MySQL configuration information.
+ Rewrite the file name `config.toml.example` to `config.toml`.
+ Init your MySQL with the file `/db/box.sql`.
+ It would be best to modify the server mode from `debug` to `release`.
+ config wallet
  - Modify the basic path file `path.yml`.
  - Config the wallet MySQL whith `config.yml`.
  - Modify the relevant files in the `local` directory to configure the full node.

### QuickStart

+ Get the Source Code

~~~sh
$ git clone git@gitlab.2se.com:boxproject/apiServer.git
~~~

+ Build

~~~sh
$ make build
~~~

+ Rebuild

~~~sh
$ make rebuild
~~~

+ Start

~~~sh
$ ./apiServer start
~~~

### Contributing

Find more documentation for developers on [API.md](./docs/api.md)

### License

Licensed under the Apache License, Version 2.0, Copyright 2018. bolaxy.org authors.

```
 Copyright 2018. bolaxy.org authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
```


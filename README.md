# curlinator

URL download helper service

## Setup

### Prerequisites

Install Go for your platform: https://golang.org

I recommend gimme: https://github.com/travis-ci/gimme

```
mkdir -p ~/bin
curl -sL -o ~/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
chmod +x ~/bin/gimme
```

### Download

Clone this repository into the `src` directory of your
[`$GOPATH`](https://golang.org/doc/code.html#GOPATH)
(default: `~/go/src`):

```
git clone https://github.com/thomasheller/curlinator ~/go/src/github.com/thomasheller/curlinator
```

### Build

```
eval $(~/bin/gimme stable)
cd ~/go/src/github.com/thomasheller/curlinator
go get
go build
```

### Deploy as systemd service

- Put the `curlinator` binary in `/usr/local/bin`
- Put the `curlinator.service` file in `/lib/systemd/system`
- Add user and group `curlinator` and create the `/var/curlinator` directory (`useradd -U -m -d /var/curlinator curlinator`)

Enable service:

```
systemctl enable curlinator
service curlinator start
service curlinator status
```

## Usage

```
curl -d '{"url":"https://www.example.com"}' 'http://localhost:8000/add'

curl -d '{"url":"https://www.example.com"}' 'http://localhost:8000/delete'

curl 'http://localhost:8000/list'

curl 'http://localhost:8000/status'
```


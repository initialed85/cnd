# cnd

A simple commandline tool to help you automate things; `cnd` is short for `check-n-do` and provides some sugar
as long as your problem fits the following pattern:

- Check the state of something against the last time it was checked
- If changed, execute an action

## Concept

The user can invoke `cnd` as many times as they like, the given `check` will execute but the given `do` will only execute if the result of
the `check` differed from the last time `cnd` was ran.

## Usage

```shell
cnd \
  -job '(an optional job name; defaults to the current working directory name)' \
  -check '(a command that gathers our comparison state)' \
  -do '(a command that executes an action)'
```

## Example

### Redeploying a service when a Git repo changes

```shell
cnd -check 'git pull >/dev/null 2>&1 && git log -b master HEAD~1..HEAD' -do 'make down ; make up'
```

When executed and there are no changes, the output will look something like this:

```shell
2022/09/21 09:03:36 check-n-do; jobName="dinosaur"
2022/09/21 09:03:36 check; command="git pull >/dev/null 2>&1 && git log -b master HEAD~1..HEAD"
2022/09/21 09:03:38 check; output="commit e03e2d359daf575e242b15c1c8d70a013b85b199\nAuthor: Edward Beech <initialed85@gmail.com>\nDate:   Wed Sep 21 16:57:04 2022 +0800\n\n    Empty commit to see if my ghetto CD stuff is working"
2022/09/21 09:03:38 check; lastOutput="commit e03e2d359daf575e242b15c1c8d70a013b85b199\nAuthor: Edward Beech <initialed85@gmail.com>\nDate:   Wed Sep 21 16:57:04 2022 +0800\n\n    Empty commit to see if my ghetto CD stuff is working"
2022/09/21 09:03:38 check; changed=false
```

And if there are changes, the output will look something like this:

```shell
2022/09/21 09:04:36 check-n-do; jobName="dinosaur"
2022/09/21 09:04:36 check; command="git pull >/dev/null 2>&1 && git log -b master HEAD~1..HEAD"
2022/09/21 09:04:39 check; output="commit c8d7ccd6d0c820f14856e33f48fce6c3a1a634db\nAuthor: Edward Beech <initialed85@gmail.com>\nDate:   Wed Sep 21 17:04:29 2022 +0800\n\n    Another empty commit to see if my ghetto CD stuff is working"
2022/09/21 09:04:39 check; lastOutput="commit e03e2d359daf575e242b15c1c8d70a013b85b199\nAuthor: Edward Beech <initialed85@gmail.com>\nDate:   Wed Sep 21 16:57:04 2022 +0800\n\n    Empty commit to see if my ghetto CD stuff is working"
2022/09/21 09:04:39 check; changed=true
2022/09/21 09:04:39 do; command="make down ; make up"
2022/09/21 09:04:44 do; output="docker compose -f docker/docker-compose.yml down --remove-orphans --volumes || true\nContainer docker-frontend-1  Stopping\nContainer docker-frontend-1  Stopping\nContainer docker-backend-1  Stopping\nContainer docker-backend-1  Stopping\nContainer docker-backend-1  Stopped\nContainer docker-backend-1  Removing\nContainer docker-backend-1  Removed\nContainer docker-frontend-1  Stopped\nContainer docker-frontend-1  Removing\nContainer docker-frontend-1  Removed\nNetwork dinosaur-internal  Removing\nNetwork dinosaur-external  Removing\nNetwork dinosaur-internal  Removed\nNetwork dinosaur-external  Removed\ndocker build -t dinosaur-session -f docker/session/Dockerfile ./docker/session\nSending build context to Docker daemon  8.704kB\r\r\nStep 1/22 : FROM ubuntu:22.04\n ---> 2dc39ba059dc\nStep 2/22 : RUN apt-get update &&     apt-get install -y     curl git npm entr     golang-1.18 python3 default-jdk lua5.4 luarocks     procps file strace screen     net-tools inetutils-ping traceroute netcat tcpdump iproute2\n ---> Using cache\n ---> e945329d79f9\nStep 3/22 : RUN npm install -g typescript ts-node\n ---> Using cache\n ---> 46717d934463\nStep 4/22 : RUN curl https://sh.rustup.rs -sSf | sh -s -- -y\n ---> Using cache\n ---> 518759e3badb\nStep 5/22 : ENV PATH=${PATH}:/usr/lib/go-1.18/bin:/root/go/bin/:/root/.cargo/bin\n ---> Using cache\n ---> 0ab91dbee20b\nStep 6/22 : RUN go install github.com/sorenisanerd/gotty@latest\n ---> Using cache\n ---> ca817ab25f1b\nStep 7/22 : RUN mkdir -p /srv/cmd/\n ---> Using cache\n ---> 0c226f0dc3f9\nStep 8/22 : WORKDIR /srv/\n ---> Using cache\n ---> 9383bb62d031\nStep 9/22 : RUN luarocks install luasocket\n ---> Using cache\n ---> f5529435e81c\nStep 10/22 : RUN npm i --save-dev @types/node\n ---> Using cache\n ---> 7aa1a6c9f2ef\nStep 11/22 : RUN echo 'termcapinfo xterm* ti@:te@' >> /root/.screenrc\n ---> Using cache\n ---> 8b8e167d5e6d\nStep 12/22 : COPY docker-entrypoint.sh /docker-entrypoint.sh\n ---> Using cache\n ---> 91b0e7c1d04e\nStep 13/22 : COPY loop.sh /loop.sh\n ---> Using cache\n ---> 1baa4281d86d\nStep 14/22 : COPY build.sh /build.sh\n ---> Using cache\n ---> b779bc5c4f28\nStep 15/22 : COPY run.sh /run.sh\n ---> Using cache\n ---> d866e4bfe5f9\nStep 16/22 : COPY watch.sh /watch.sh\n ---> Using cache\n ---> eb766bb51055\nStep 17/22 : ENV GOTTY_PORT=${GOTTY_PORT:-8080}\n ---> Using cache\n ---> 63577e48f262\nStep 18/22 : ENV GOTTY_PATH=${GOTTY_PATH}\n ---> Using cache\n ---> 799685a63c38\nStep 19/22 : ENV BASE_FOLDER_PATH=${BASE_FOLDER_PATH:-/srv/}\n ---> Using cache\n ---> 64be509c0a1d\nStep 20/22 : ENV BUILD_CMD=${BUILD_CMD}\n ---> Using cache\n ---> 0b290e0759b2\nStep 21/22 : ENV RUN_CMD=${RUN_CMD}\n ---> Using cache\n ---> f5758eec0880\nStep 22/22 : ENTRYPOINT [\"/docker-entrypoint.sh\"]\n ---> Using cache\n ---> 08ee7a337db7\nSuccessfully built 08ee7a337db7\nSuccessfully tagged dinosaur-session:latest\ndocker compose -f docker/docker-compose.yml build --parallel\n#1 [docker-backend internal] load build definition from Dockerfile\n#1 transferring dockerfile: 32B 0.0s done\n#1 DONE 0.1s\n\n#2 [docker-frontend internal] load build definition from Dockerfile\n#2 transferring dockerfile: 32B 0.1s done\n#2 DONE 0.1s\n\n#3 [docker-frontend internal] load .dockerignore\n#3 transferring context: 2B done\n#3 DONE 0.0s\n\n#4 [docker-backend internal] load .dockerignore\n#4 transferring context: 2B 0.0s done\n#4 DONE 0.0s\n\n#5 [docker-backend internal] load metadata for docker.io/library/ubuntu:22.04\n#5 DONE 0.0s\n\n#6 [docker-backend internal] load metadata for docker.io/library/golang:1.18\n#6 DONE 1.9s\n\n#7 [docker-frontend internal] load metadata for docker.io/library/node:14.20.0\n#7 DONE 1.9s\n\n#8 [docker-backend stage-1  1/11] FROM docker.io/library/ubuntu:22.04\n#8 DONE 0.0s\n\n#9 [docker-backend build 1/9] FROM docker.io/library/golang:1.18@sha256:040bc0980888e2df09499c73156f3f869d843b7033d11b465fee7f1a358a494a\n#9 DONE 0.0s\n\n#10 [docker-backend internal] load build context\n#10 transferring context: 2.57kB 0.0s done\n#10 DONE 0.0s\n\n#11 [docker-backend stage-1  7/11] RUN mkdir -p /srv/pkg/sessions/languages\n#11 CACHED\n\n#12 [docker-backend stage-1  3/11] RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg\n#12 CACHED\n\n#13 [docker-backend stage-1  6/11] WORKDIR /srv/\n#13 CACHED\n\n#14 [docker-backend build 6/9] COPY cmd /srv/cmd\n#14 CACHED\n\n#15 [docker-backend stage-1  8/11] COPY pkg/sessions/languages /srv/pkg/sessions/languages\n#15 CACHED\n\n#16 [docker-backend stage-1 10/11] COPY docker/session /srv/docker/session\n#16 CACHED\n\n#17 [docker-backend stage-1  9/11] RUN mkdir -p /srv/docker/\n#17 CACHED\n\n#18 [docker-backend build 3/9] COPY go.mod /srv/go.mod\n#18 CACHED\n\n#19 [docker-backend build 8/9] COPY pkg /srv/pkg\n#19 CACHED\n\n#20 [docker-backend build 7/9] COPY internal /srv/internal\n#20 CACHED\n\n#21 [docker-backend stage-1  2/11] RUN apt-get update && apt-get install -y ca-certificates curl gnupg lsb-release golang-1.18\n#21 CACHED\n\n#22 [docker-backend build 5/9] RUN --mount=type=cache,target=/root/.cache/go-build go mod download\n#22 CACHED\n\n#23 [docker-backend stage-1  4/11] RUN echo   \"deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu   $(lsb_release -cs) stable\" | tee /etc/apt/sources.list.d/docker.list > /dev/null\n#23 CACHED\n\n#24 [docker-backend build 9/9] RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o main /srv/cmd/main.go\n#24 CACHED\n\n#25 [docker-backend build 2/9] WORKDIR /srv/\n#25 CACHED\n\n#26 [docker-backend stage-1  5/11] RUN apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin\n#26 CACHED\n\n#27 [docker-backend build 4/9] COPY go.sum /srv/go.sum\n#27 CACHED\n\n#28 [docker-backend stage-1 11/11] COPY --from=build /srv/main /srv/main\n#28 CACHED\n\n#29 [docker-frontend internal] load metadata for docker.io/library/nginx:1.23.1\n#29 DONE 1.9s\n\n#30 [docker-backend] exporting to image\n#30 exporting layers done\n#30 writing image sha256:a638b17a3f22988e74afc7a11f62b042103b882729fa578ed78ed01fca23813a done\n#30 naming to docker.io/library/docker-backend done\n#30 DONE 0.0s\n\n#31 [docker-frontend build 1/9] FROM docker.io/library/node:14.20.0@sha256:6adfb0c2a9db12a06893974bb140493a7482e2b3df59c058590594ceecd0c99b\n#31 DONE 0.0s\n\n#32 [docker-frontend stage-1 1/4] FROM docker.io/library/nginx:1.23.1@sha256:0b970013351304af46f322da1263516b188318682b2ab1091862497591189ff1\n#32 DONE 0.0s\n\n#33 [docker-frontend internal] load build context\n#33 transferring context: 1.20MB 0.0s done\n#33 DONE 0.0s\n\n#34 [docker-frontend build 6/9] COPY frontend/tsconfig.json /srv/tsconfig.json\n#34 CACHED\n\n#35 [docker-frontend build 8/9] COPY frontend/public /srv/public\n#35 CACHED\n\n#36 [docker-frontend build 4/9] COPY frontend/package-lock.json /srv/package-lock.json\n#36 CACHED\n\n#37 [docker-frontend stage-1 2/4] COPY docker/frontend/default.conf /etc/nginx/conf.d/default.conf\n#37 CACHED\n\n#38 [docker-frontend build 3/9] COPY frontend/package.json /srv/package.json\n#38 CACHED\n\n#39 [docker-frontend build 7/9] COPY frontend/src /srv/src\n#39 CACHED\n\n#40 [docker-frontend build 9/9] RUN npm run build\n#40 CACHED\n\n#41 [docker-frontend stage-1 3/4] COPY --from=build /srv/build/static /usr/share/nginx/html/static\n#41 CACHED\n\n#42 [docker-frontend build 2/9] WORKDIR /srv/\n#42 CACHED\n\n#43 [docker-frontend build 5/9] RUN npm ci\n#43 CACHED\n\n#44 [docker-frontend stage-1 4/4] COPY --from=build /srv/build/* /usr/share/nginx/html/\n#44 CACHED\n\n#30 [docker-frontend] exporting to image\n#30 exporting layers done\n#30 writing image sha256:6458f8593ff555f1a6ef183ea5d2f0e1402538bcb1226ea9934cbb5ddc1790f8 done\n#30 naming to docker.io/library/docker-frontend done\n#30 DONE 0.0s\n\nUse 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them\ndocker compose -f docker/docker-compose.yml up -d\nNetwork dinosaur-internal  Creating\nNetwork dinosaur-internal  Created\nNetwork dinosaur-external  Creating\nNetwork dinosaur-external  Created\nContainer docker-frontend-1  Creating\nContainer docker-backend-1  Creating\nContainer docker-frontend-1  Created\nContainer docker-backend-1  Created\nContainer docker-backend-1  Starting\nContainer docker-frontend-1  Starting\nContainer docker-frontend-1  Started\nContainer docker-backend-1  Started"
2022/09/21 09:04:44 check-n-do; done=true
```

If you wanted to run the above every minute as a `cron` job; simple edit your `crontab` with `crontab -e` and add the following:

```shell
  * *  *   *   *     /bin/bash -c "cd /home/edward/Projects/Home/dinosaur && cnd -check 'git pull >/dev/null 2>&1 && git log -b master HEAD~1..HEAD' -do 'make down ; make up'"
```

## Installation

```shell
go install github.com/initialed85/cnd@latest
```

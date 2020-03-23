# はじめてのdockerコマンド

本手順には次のdockerコマンドで実施した手順を記載します

```sh
$ docker --version
Docker version 19.03.2, build 6a30dfc
```

[DockerHub](https://hub.docker.com/)で公開されているコンテナイメージは自由に使うことができます。
ここでは、WebサーバであるnginxのうちalpineOSで起動するものをダウンロードしましょう。
このようなコンテナイメージが保存されている場所を、`レジストリ`と呼びます。
DockerHubが最も有名ですが、例えばGCPにはGoogle Container Registryというサービスがあり、自分たちのコンテナをプライベートに保存しておくことができます。
コンテナイメージをダウンロードする際は`docker pull`コマンドを次のように使います。

```sh
$ docker pull nginx:stable-alpine
stable-alpine: Pulling from library/nginx
9d48c3bd43c5: Pull complete
7a56a3a1208e: Pull complete
Digest: sha256:096c4b3464e2e465f20e9d704f1a0f8d27584df4d6758b6d00a14911cc9bb888
Status: Downloaded newer image for nginx:stable-alpine
docker.io/library/nginx:stable-alpine
```

ダウンロードしたイメージの一覧は、`docker images` または`docker image ls`のどちらでも実行できます。

```sh
$ docker images
REPOSITORY    TAG                 IMAGE ID            CREATED             SIZE
nginx         stable-alpine       ba8b51e283ea        35 minutes ago      14.3MB
...(省略)
```

イメージがどのようなものなのかを調べるコマンドに`docker inspect`があります。
出力はなかなか膨大ですが、意味のわかりそうなところを探してみてください。

```sh
$ docker inspect nginx:stable-alpine
...(省略)
```

さっそくこのコンテナを起動して、nginxサーバをコンテナとして起動してみましょう。
コンテナを起動するには`docker run`を使用します。runにはたくさんのオプションがありますが、ここではひとまず必要なものだけ指定します。

```sh
$ docker run -d -p 127.0.0.1:8080:80 nginx:stable-alpine
a034dc51f2039b20c698b6ffa14d6729fbe3ae514a2e26a188c13184503204c7
```

`-p`オプションは次のような意味を持っています。
`-p ホスト名:転送元ポート番号:転送先コンテナポート番号`
つまり、localhostである127.0.0.1の8080番ポートへのトラフィックをコンテナの80番に流そうとしています。

起動中のプロセスを確認するためには`docker ps`コマンドを実行します。

```sh
$ docker ps
CONTAINER ID        IMAGE                 COMMAND                  CREATED              STATUS              PORTS                    NAMES
a034dc51f203        nginx:stable-alpine   "nginx -g 'daemon of…"   About a minute ago   Up About a minute   127.0.0.1:8080->80/tcp   vigorous_kepler
```

実際にlocalhostにリクエストを送ってみましょう。

```sh
$ curl -i localhost:8080
HTTP/1.1 200 OK
Server: nginx/1.16.1
Date: Thu, 26 Sep 2019 12:29:39 GMT
```

このように、コンテナを使えば自分たちの起動したいプロセスを、dockerが動く環境であればどこでも動かすことができます。
コンテナで実行されるプロセスのログは、`docker logs <コンテナ名> | <コンテナID>`コマンドで見られます。
上の出力の例では、コンテナ名は`vigorous_kepler`で、コンテナIDは`a034dc51f203`です。

```sh
$ docker logs a0
172.17.0.1 - - [26/Sep/2019:12:29:39 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.54.0" "-"
```

このように、コンテナIDは全てを指定する必要はなく、起動中のコンテナで一意になる桁まで指定すれば正常にコマンドが実行できます。

このコンテナを削除するには、`docker rm -f <コンテナ名> | <コンテナID>`を実行します。

```sh
$ docker rm -f a0
a0
```

## コンテナのビルド

今度はDockerfileからコンテナをビルドしてみましょう。
コンテナのビルドは`docker build -t <イメージ名>:<タグ名> .`で実行できます。
このREADMEが配置されているディレクトリで次のように実行してください。

```sh
$ docker build -t first-build-container:v1.0 .
Sending build context to Docker daemon  10.24kB
Step 1/10 : FROM golang:latest as builder
 ---> ae999aee9560
Step 2/10 : WORKDIR /app
 ---> Running in 63ad5561e3df
Removing intermediate container 63ad5561e3df
 ---> 13276404f4b6
Step 3/10 : COPY go.mod main.go ./
 ---> f9a84303a792
Step 4/10 : RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
 ---> Running in 75c4d6d7f81c
Removing intermediate container 75c4d6d7f81c
 ---> 1095e7e7ce7a
Step 5/10 : FROM alpine:latest
 ---> b7b28af77ffe
Step 6/10 : RUN apk --no-cache add ca-certificates curl
 ---> Using cache
 ---> bda29d552247
Step 7/10 : WORKDIR /root/
 ---> Using cache
 ---> 0395899fefe0
Step 8/10 : COPY --from=builder /app/main .
 ---> Using cache
 ---> b4e92c590eef
Step 9/10 : EXPOSE 8080 80
 ---> Using cache
 ---> b8451cad7450
Step 10/10 : CMD ["./main"]
 ---> Using cache
 ---> ba8b51e283ea
Successfully built ba8b51e283ea
Successfully tagged first-build-container:v1.0
```

およそこのような出力が得られたはずです。
ビルドされたコンテナイメージの存在を確認してみましょう。

```sh
$ docker images
REPOSITORY                     TAG                 IMAGE ID            CREATED             SIZE
first-build-container          v1.0                ba8b51e283ea        57 minutes ago      14.3MB
```

最後にこのコンテナを実行して同じようにHTTPリクエストを送ってみてください。

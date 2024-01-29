# ClusterIP

`nginx-clusterip.yaml`を見てみましょう。
このマニフェストファイルには、`app=web`のラベルを持つPodと、そのラベルをselectorとするClusterIP Serviceを起動する定義が書いてあります。
kubectlコマンドを利用してデプロイしてください。

ClusterIPはクラスタの内部に公開するServiceなので、このままローカル環境にいてはトラフィックを送れません。
そのため、alpineコンテナを起動するPodをデバッグ用のPodとして用意します。
このようなケースでわざわざマニフェストファイルを書くのはめんどうなので、`kubectl run`コマンドで直接Podを作成しましょう。

```sh
$ kubectl run -it alpine --image=alpine -- ash
If you don't see a command prompt, try pressing enter.
/ #
```

プロンプトが表示されたら、`apk`パッケージのアップデート、`curl`コマンドのインストールを行いましょう。
インストールが終わったら、 `nginx-clusterip`に対してリクエストを送ってみてください。

```sh
/ # apk update
fetch http://dl-cdn.alpinelinux.org/alpine/v3.10/main/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/v3.10/community/x86_64/APKINDEX.tar.gz
v3.10.2-67-g8eb76d4bd1 [http://dl-cdn.alpinelinux.org/alpine/v3.10/main]
v3.10.2-66-ga328d90a4c [http://dl-cdn.alpinelinux.org/alpine/v3.10/community]
OK: 10336 distinct packages available

/ # apk add curl
(1/4) Installing ca-certificates (20190108-r0)
(2/4) Installing nghttp2-libs (1.39.2-r0)
(3/4) Installing libcurl (7.66.0-r0)
(4/4) Installing curl (7.66.0-r0)
Executing busybox-1.30.1-r2.trigger
Executing ca-certificates-20190108-r0.trigger
OK: 7 MiB in 18 packages

/ # curl -i nginx-clusterip
HTTP/1.1 200 OK
Server: nginx/1.16.1
```

さて、なぜ`nginx-clusterip`というホスト名が解決できたのでしょうか?

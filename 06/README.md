# 公開コンテナ名
`gcr.io/ca-tominaga-test/healthy-error-server`

# readinessProbeとlivenessProbe

本番環境でコンテナを起動し続けるためには、この2つの設定が必要不可欠です。
ここでは2つのProbeの動作を確認するためのコンテナを使って理解をしましょう。
まずはそのコンテナをビルドします。
このディレクトリにあるDockerfileを使ってコンテナをビルドして、Google Container Registry(GCR)にpushをします。
`YOUR_PROJECT_ID`の部分は、自分のGCPプロジェクトIDに変更して実行してください。

## コンテナのビルド

```sh
$ docker build -t gcr.io/YOUR_PROJECT_ID/probe .
Sending build context to Docker daemon  6.144kB
...
```

GCRにコンテナをpushするために、dockerコマンドに対してGCPに認証させる必要があります。

```sh
$ gcloud auth configure-docker
```

GCRにイメージを保存する場合、イメージ名の規則は `gcr.io/<project-id>/<イメージ名>`である必要があるため、先程のビルド時に指定をしました。
ビルドしたコンテナは、通常のdockerコマンドとまったく同じ作法でpushすることができます。

```sh
$ docker push gcr.io/YOUR_PRIOJECT_ID/probe
The push refers to repository [gcr.io/ca-container-book/probe]
c4fbb4c6eef7: Layer already exists
3e78f78901d8: Layer already exists
1bfeebd65323: Layer already exists
```

コンテナイメージがpushできたら、さっそくKubernetesにデプロイします。
`probe.yaml`を読んで、kubectlコマンドでデプロイしてください。

readinessProbeに成功している場合、kubectl get pods のREADYステータスが 1/1 となります。
このコンテナはルートパス`/`へのリクエストの結果が`/health`へのリクエストで成功するようになり、`/unhealth`へのリクエストで失敗するようになります。
つまり、この状態で`/unhealth`へリクエストを送ると、readinessProbeに失敗するようになります。

```sh
# 最初はリクエストに成功する
$ curl -i 34.84.223.243:30002
HTTP/1.1 200 OK

# /unhealthにリクエストを送る
$ curl -i 34.84.223.243:30002/unhealth
HTTP/1.1 200 OK

# リクエストに失敗するようになる
$ curl -i 34.84.223.243:30002
HTTP/1.1 503 Service Unavailable
```

この状態でREADYステータスを確認します。

```sh
$ kubectl get pods probe
NAME    READY   STATUS    RESTARTS   AGE
probe   0/1     Running   0          3m52s
```

しばらくしてからまた同じようにNodePortに対してトラフィックを送ってみましょう。

```sh
$ curl -i 34.84.223.243:30002
```

今度はレスポンスが返ってこないはずです。これは、readinessProbeがREADYではないPodにトラフィックを送信しないようにする振る舞いを示しています。

ところで、何らかの原因でPodが正常に処理を実行できないにもかかわらず、そのPodが起動しっぱなしでは困ります。
Podに障害が発生した際にプロセスを再起動させるものがlivenessProbeです。
probe PodとServiceを削除して、`probe.yaml`のコメントアウトを外してから再度デプロイしてください。
マニフェストファイルから作成したリソースを削除するには、`kubectl delete -f`を実行すると良いでしょう。

```sh
$ kubectl delete -f probe.yaml
pod "probe" deleted
service "probe" deleted
```

Podが起動し直したら、同じように `/unhealth`にリクエストを送ってみましょう。
リクエストを送ったら、`kubectl get pods probe -w`で、probe Podがどのような振る舞いをするかを観測します。

```sh
$ curl -i 34.84.223.243:30002/unhealth
HTTP/1.1 200 OK

$ kubectlo get pods probe -w
NAME    READY   STATUS    RESTARTS   AGE
probe   1/1     Running   0          61s
probe   0/1   Running   1     73s
probe   1/1   Running   1     84s
```

RESTARTSに注目してください。livenessProbeに失敗した結果、Podが再起動されたことが確認できます。

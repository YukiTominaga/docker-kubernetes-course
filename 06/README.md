# readinessProbeとlivenessProbe

本番環境でコンテナを起動し続けるためには、この2つの設定が必要不可欠です。
ここでは2つのProbeの動作を確認するためのコンテナを使って理解をしましょう。
`probe.yaml`を読んで、kubectlコマンドでデプロイしてください。

readinessProbeに成功している場合、kubectl get pods のREADYステータスが 1/1 となります。
このコンテナはルートパス`/`へのリクエストの結果が`/health`へのリクエストで成功するようになり、`/unhealth`へのリクエストで失敗するようになります。
つまり、この状態で`/unhealth`へリクエストを送ると、readinessProbeに失敗するようになります。

LoadBalancerの外部IPアドレスが `34.84.223.243` と仮定します。

```sh
# 最初はリクエストに成功する
$ curl -i 34.84.223.243
HTTP/1.1 200 OK

# /unhealthにリクエストを送る
$ curl -i 34.84.223.243/unhealth
HTTP/1.1 200 OK

# リクエストに失敗するようになる
$ curl -i 34.84.223.243
HTTP/1.1 503 Service Unavailable
```

この状態でREADYステータスを確認します。

```sh
$ kubectl get pods probe
NAME    READY   STATUS    RESTARTS   AGE
probe   0/1     Running   0          3m52s
```

しばらくしてからまた同じようにLoadBalancerに対してトラフィックを送ってみましょう。

```sh
$ curl -i 34.84.223.243
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
$ curl -i 34.84.223.243/unhealth
HTTP/1.1 200 OK

$ kubectlo get pods probe -w
NAME    READY   STATUS    RESTARTS   AGE
probe   1/1     Running   0          61s
probe   0/1   Running   1     73s
probe   1/1   Running   1     84s
```

RESTARTSに注目してください。livenessProbeに失敗した結果、Podが再起動されたことが確認できます。

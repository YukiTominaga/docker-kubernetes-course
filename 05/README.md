# LoadBalancer

実際にコンテナを外部に公開する際には、LoadBalancerを使います。
`nginx-loadbalancer.yaml`をkubectlでデプロイして、LoadBalancerを作成してください。
作成したらすぐに次のコマンドを実行してください。

```sh
$ kubectl get svc nginx-lb -w
NAME       TYPE           CLUSTER-IP   EXTERNAL-IP   PORT(S)        AGE
nginx-lb   LoadBalancer   10.28.0.78   <pending>     80:32085/TCP   10s
```

ここで、EXTERNAL-IPが"pending"となっているはずです。
これは、GCPのロードバランサが外部IPアドレスを確保している最中であることを意味しています。
ロードバランサが外部IPアドレスの確保に成功すると、"pending"部分が次のように書き換わります。

```sh
$ kubectl get svc nginx-lb -w
NAME       TYPE           CLUSTER-IP   EXTERNAL-IP   PORT(S)        AGE
nginx-lb   LoadBalancer   10.28.0.78   <pending>     80:32085/TCP   10s
nginx-lb   LoadBalancer   10.28.0.78   35.200.27.42   80:32085/TCP   54s
```

kubectl get コマンドの最後につけている`-w`は、一覧表示中のリソースに変更があった場合に自動的に新しい情報を出力してくれるオプションです。
この場合は、curlを使って次のようにロードバランサにトラフィックを送信できます。

```sh
$ curl -i 35.200.27.42
HTTP/1.1 200 OK
Server: nginx/1.16.1
```

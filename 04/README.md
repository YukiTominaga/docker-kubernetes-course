# NodePort

今度はNodePortです。
まずはNodePort Serviceをマニフェストファイル`nginx-nodeport.yaml`から作成しましょう。

NodePortにトラフィックを送信するためには、もちろんNodeのIPアドレスが必要です。
Kubernetesクラスタを構成するNodeのIPアドレスは`kubectl get nodes`で調べられます。

```sh
$ kubectl get nodes -o wide
NAME                                         STATUS   ROLES    AGE    VERSION          INTERNAL-IP   EXTERNAL-IP     OS-IMAGE
gke-gke-cluster-default-pool-fce83ce4-m9ng   Ready    <none>   6h5m   v1.12.8-gke.10   10.146.0.19   34.84.250.55    Container-Optimized OS from Google
gke-gke-cluster-default-pool-fce83ce4-w2vb   Ready    <none>   8h     v1.12.8-gke.10   10.146.0.17   34.84.223.243   Container-Optimized OS from Google
```

この出力の例では`34.84.250.55`または`34.84.223.243`です。
早速`34.84.223.243`の30001番にトラフィックを送ってみたいところですが、GCPではファイアウォールルールが許可されていないため、トラフィックを送信することができません。
まずはこのGKEクラスタが作成されているdefaultネットワークに対して、NodePortの範囲である`30000~32767`番のポートへのTCPを許可するルールを作ります。

```sh
$ gcloud compute firewall-rules create allow-nodeport \
    --allow=tcp:30000-32767 \
    --source-ranges=0.0.0.0/0 \
    --network=default
Creating firewall...⠧Created [https://www.googleapis.com/compute/v1/projects/ca-container-book/global/firewalls/allow-nodeport].
Creating firewall...done.
NAME            NETWORK  DIRECTION  PRIORITY  ALLOW            DENY  DISABLED
allow-nodeport  default  INGRESS    1000      tcp:30000-32767        False
```

こうすることでNodePortへのトラフィックが許可されます。

```sh
$ curl -i 34.84.223.243:30001
HTTP/1.1 200 OK
Server: nginx/1.16.1
```

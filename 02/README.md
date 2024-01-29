# はじめてのKubernetes

まずはPodを作成しましょう。
マニフェストファイルに記載されたリソースを作成する時は、`kubectl apply`コマンドを使用します。

```sh
$ kubectl apply -f basic-pod.yaml
pod/nginx created
pod/nginx-no-label created
```

起動中のリソースを一覧する際は、`kubectl get`コマンドを使います。
Podを一覧する時は`kubectl get pods`です。

```sh
$ kubectl get pods
NAME             READY   STATUS    RESTARTS   AGE
nginx            1/1     Running   0          2m44s
nginx-no-label   1/1     Running   0          2m17s
```

より詳しい情報を見る時は`-o wide`をつけます。

```sh
$ kubectl get pods -o wide
NAME             READY   STATUS    RESTARTS   AGE     IP          NODE                                         NOMINATED NODE
nginx            1/1     Running   0          3m24s   10.24.1.4   gke-gke-cluster-default-pool-fce83ce4-m9ng   <none>
nginx-no-label   1/1     Running   0          2m57s   10.24.1.5   gke-gke-cluster-default-pool-fce83ce4-m9ng   <none>
```

Podには全てラベルをつけることができます。特定のラベルがついているPodの一覧はよく使うので必ず覚えましょう。
`-l`オプションを使用することでラベルを指定することができます。

```sh
$ kubectl get pods -l app=web
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          7m30s
```

起動中のPodの中にあるコンテナにログインすることは非常によくありますので必ず覚えてください。
`kubectl exec`コマンドでログインができます。

```sh
$ kubectl exec -it nginx -- ash
/ # cat /etc/hosts
# Kubernetes-managed hosts file.
127.0.0.1	localhost
::1	localhost ip6-localhost ip6-loopback
fe00::0	ip6-localnet
fe00::0	ip6-mcastprefix
fe00::1	ip6-allnodes
fe00::2	ip6-allrouters
```

Nginxサーバのコンテナを起動するPodに直接トラフィックを送ってみましょう。
kubectlコマンドには、Podへのポートフォワーディングを行うコマンド`kubectl port-forward`があり、`kubectl port-forward <Pod名> ローカルポート番号:Podのポート番号`のように使います。

```sh
$ kubectl port-forward nginx 8081:80
Forwarding from 127.0.0.1:8081 -> 80
Forwarding from [::1]:8081 -> 80
```

起動中のPodの定義を少しだけ変更したい場合、わざわざマニフェストファイルを書き直して`kubectl apply`をするのは面倒です。
そんな時は`kubectl edit <Pod名>`で定義を直接変更することができます。

```sh
$ kubectl edit pods nginx
※ vimが起動して編集することができます
```

ここまでPodに関する操作のみ行ってきましたが、Kubernetesにはもっとたくさんのリソースがあります。
どんなリソースがあるのか見てみましょう。

```sh
$ kubectl api-resources
NAME                              SHORTNAMES   APIGROUP                       NAMESPACED   KIND
bindings                                                                      true         Binding
componentstatuses                 cs                                          false        ComponentStatus
configmaps                        cm                                          true         ConfigMap
endpoints                         ep                                          true         Endpoints
events                            ev                                          true         Event
...
```

それぞれのリソースにどのような定義が書けてどのような意味があるのか、という疑問がそのうち出てきます。
どれだけKubernetesを長く続けていても、全ての定義を覚えるのはなかなか難しいので`kubectl explain`コマンドで調べられることを覚えておきましょう。

```sh
$ kubectl explain Pod.spec
KIND:     Pod
VERSION:  v1

RESOURCE: spec <Object>

DESCRIPTION:
     Specification of the desired behavior of the pod. More info:
     https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status

     PodSpec is a description of a pod.
...
```

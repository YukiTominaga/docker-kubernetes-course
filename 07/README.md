# Deploymentとstern

Deploymentで複数の同等なPodを起動し、LoadBalancerに対してHTTPリクエストをして、複数のPodにトラフィックが分散されている様子を観測します。
まずは`nginx-deployment.yaml`を読んで内容を理解してデプロイください。(HorizontalPodAutoscalerは今は無視してください)

```sh
$ kubectl get pods -l app=nginx
NAME                            READY   STATUS              RESTARTS   AGE
nginx-deploy-54758469f5-gvctz   1/1     Running             0          40s
nginx-deploy-54758469f5-n54zk   1/1     Running             0          102s
nginx-deploy-54758469f5-xnwgj   0/1     ContainerCreating   0          40s
```

特定のPodのログを見る場合、`kubectl logs -f <Pod名>`で見られますが、これは単一のPodのログを見るのに適しています。
複数のPodのログを見るにはGCPのStackdriver Loggingを利用してもよいですが、ここでは[stern](https://github.com/wercker/stern)を使います。
手順に沿ってsternをインストールしてください(任意)。

sternは`stern <ログを見たいPodのprefix>`というように使います。今回は`nginx-deploy`を含むPodのログを見たいので次のようにします。

```sh
$ stern nginx-deploy
+ nginx-deploy-54758469f5-n54zk › nginx
+ nginx-deploy-54758469f5-xnwgj › nginx
+ nginx-deploy-54758469f5-gvctz › nginx
```

別のプロンプトを用意して、LoadBalancer `nginx-deploy-lb`の外部IP宛にトラフィックを送信してください。
すると、sternの出力にnginxへのアクセスログが吐き出されます。

```sh
nginx-deploy-54758469f5-xnwgj nginx 10.24.2.1 - - [26/Sep/2019:16:18:41 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.54.0" "-"
nginx-deploy-54758469f5-n54zk nginx 10.146.0.19 - - [26/Sep/2019:16:18:43 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.54.0" "-"
nginx-deploy-54758469f5-n54zk nginx 10.146.0.19 - - [26/Sep/2019:16:18:45 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.54.0" "-"
```

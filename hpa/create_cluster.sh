gcloud beta container clusters create hpa \
    --disk-size=20 \
    --num-nodes=1 \
    --region=asia-northeast1 \
    --node-locations=asia-northeast1-a,asia-northeast1-b,asia-northeast1-c \
    --enable-stackdriver-kubernetes \
    --cluster-version=latest \
    --enable-autoscaling \
    --max-nodes=100 \
    --min-nodes=1 \
    --preemptible

gcloud container node-pools create nginx-pool \
    --cluster=hpa \
    --region=asia-northeast1 \
    --node-locations=asia-northeast1-a,asia-northeast1-b,asia-northeast1-c \
    --enable-autoscaling \
    --max-nodes=100 \
    --min-nodes=3 \
    --disk-size=20 \
    --preemptible

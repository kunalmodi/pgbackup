A simple Docker image to backup a Postgres database, for usage in cron jobs.

A typical Kubernetes use case might look like this:
```
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  namespace: mynamespace
  name: cron-pg-backup
spec:
  schedule: "1 0 * * *"
  concurrencyPolicy: Forbid
  jobTemplate:
    spec:
      backoffLimit: 3
      template:
        spec:
          containers:
          - name: cron
            image: docker.pkg.github.com/kunalmodi/pgbackup/pgbackup:latest
            env:
            - name: GET_HOSTS_FROM
              value: dns
            - name: KEYS
              value: "backup/postgres.latest.dump,backup/postgres.{ds}.dump"
            - name: POSTGRES_URL
              valueFrom:
                configMapKeyRef:
                  name: postgres-config
                  key: url
            - name: S3_REGION
              valueFrom:
                secretKeyRef:
                  name: do-api
                  key: REGION
            - name: S3_ENDPOINT
              valueFrom:
                secretKeyRef:
                  name: do-api
                  key: ENDPOINT
            - name: S3_BUCKET
              valueFrom:
                secretKeyRef:
                  name: do-api
                  key: BUCKET
            - name: S3_KEY
              valueFrom:
                secretKeyRef:
                  name: do-api
                  key: SPACES_KEY
            - name: S3_SECRET
              valueFrom:
                secretKeyRef:
                  name: do-api
                  key: SPACES_SECRET
          restartPolicy: OnFailure
```

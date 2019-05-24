# Monitoring Database

The monitoring database is a MySQL instance on Cloud SQL. It stores the error logs for each prow job
failures; and the alert status for all error patterns.

## Database Connection With Cloud Shell

Navigate to SQL in Google Cloud Platform. Click on the knative-monitoring instance. In the Overview tab,
under Connect to this Instance section, click on Connect using Cloud Shell.

A `gcloud` connection command with root user will be populated in Cloud Shell. Change it to use your
database user and press enter to connect to the instance.

## Local Database Connection

Note: You need to authenticate as an IAM user with Cloud SQL Admin API role in knative-tests to use
Cloud SQL Proxy to connect to the database.

1. (One-time) Setup [Cloud SQL Proxy](https://cloud.google.com/sql/docs/mysql/quickstart-proxy-test).

1. Start the Cloud SQL Proxy on port 3307. This allows the database to be connected through localhost.

    ```level x 4 spaces
    cloud_sql_proxy --instances=knative-tests:us-central1:knative-monitoring=tcp:3307
    ```

1. Connect to the database with your database user name, password and the below configuration.

    ```level x 4 spaces
    Host: 127.0.0.1
    Port: 3307
    Database: monitoring (optional)
    ```

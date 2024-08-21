import os
import requests
import psycopg2
from kubernetes import client, config


def insert_success(conn, cursor, deployment_name):
    cursor.execute(f"""
        INSERT INTO HealthCheck
        (app_name, success_count, last_success, last_failure)
        VALUES ('{deployment_name}', 1, NOW()::varchar, 'Not Set')
        ON CONFLICT (app_name) DO UPDATE
        SET success_count = HealthCheck.success_count + 1,
        last_success = NOW()::varchar;""")
    conn.commit()


def insert_failure(conn, cursor, deployment_name):
    cursor.execute(f"""
        INSERT INTO HealthCheck
        (app_name, failure_count, last_success, last_failure)
        VALUES ('{deployment_name}', 1, 'Not Set', NOW()::varchar)
        ON CONFLICT (app_name) DO UPDATE
        SET failure_count = HealthCheck.failure_count + 1,
        last_failure = NOW()::varchar;""")
    conn.commit()


def get_service_port(namespace, deployment_name):
    config.load_incluster_config()
    core_api = client.CoreV1Api()

    services = core_api.list_namespaced_service(namespace)

    for service in services.items:
        try:
            if service.metadata.name == deployment_name:
                for port in service.spec.ports:
                    if port.protocol == "TCP":
                        return port.port
        except Exception:
            continue

    return None


def main():
    namespace = os.getenv("KUBERNETES_NAMESPACE")

    pq_host = os.getenv("PQ_HOST")
    pq_user = os.getenv("PQ_USER")
    pq_dbname = os.getenv("PQ_DBNAME")
    pq_password = os.getenv("PQ_PASSWORD")

    config.load_incluster_config()
    apps_api = client.AppsV1Api()

    deployments = apps_api.list_namespaced_deployment(namespace)

    conn = psycopg2.connect(
        host=pq_host,
        dbname=pq_dbname,
        user=pq_user,
        password=pq_password)
    cursor = conn.cursor()

    for deployment in deployments.items:
        try:
            if deployment.spec.selector.match_labels.get("monitor") != "true":
                continue
        except Exception:
            continue

        deployment_name = deployment.metadata.name

        try:
            try:
                host = f"{deployment_name}.default.svc.cluster.local"
                port = get_service_port(namespace, deployment_name)
                response = requests.get(
                    f"http://{host}:{port}/healthz")
            except Exception as e:
                print(f"Failed {deployment_name} with error: {str(e)}")
                insert_failure(conn, cursor, deployment_name)
                continue

            status_code = response.status_code

            if status_code == 200:
                insert_success(conn, cursor, deployment_name)
            else:
                insert_failure(conn, cursor, deployment_name)

        except Exception as e:
            print(f"Failed {deployment_name} with error: {str(e)}")

    cursor.close()
    conn.close()


if __name__ == "__main__":
    main()

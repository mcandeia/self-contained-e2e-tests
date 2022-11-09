# -*- mode: Python -*
load('ext://helm_resource', 'helm_resource', 'helm_repo')
load('ext://configmap', 'configmap_create')

k6_local_resource_name = 'k6'
def k6_local_resource(run_grafana=False, add_prometheus=True):
    if add_prometheus:
        helm_repo('stable', 'https://charts.helm.sh/stable')
        helm_repo('prometheus-community', 'https://prometheus-community.github.io/helm-charts')
        helm_resource('prometheus', 'prometheus-community/kube-prometheus-stack', flags=['--wait', '--set', 'prometheus.prometheusSpec.enableRemoteWriteReceiver=true'])
        k8s_yaml('./modules/configs/grafana-k6-dashboard.yaml')
        if run_grafana:
            local_resource('grafana-portforward', serve_cmd='kubectl port-forward svc/prometheus-grafana 3000:80', resource_deps=['prometheus'])
    local_resource(k6_local_resource_name, cmd='export IMG=ghcr.io/mcandeia/k6-operator:latest && rm -rf /tmp/.k6-operator >/dev/null && git clone https://github.com/grafana/k6-operator /tmp/.k6-operator && cd /tmp/.k6-operator && make deploy && cd - && rm -rf /tmp/.k6-operator')

def k6_test_resource(name='k6-test', namespace='', from_file=None, parallelism=1, prometheus_url='http://prometheus-kube-prometheus-prometheus.default.svc.cluster.local:9090/api/v1/write', prometheus_user=None, prometheus_password=None, add_dapr=True, resource_deps=[]):
    test_file_resource = 'k6-tests-config'
    configmap_create(test_file_resource, from_file=[from_file], namespace=namespace)
    k8s_resource(new_name=test_file_resource, objects=['%s:ConfigMap:%s' % (test_file_resource, namespace)])

    k8s_yaml(encode_yaml({
        'apiVersion': 'k6.io/v1alpha1',
        'kind': 'K6',
        'metadata': {
            'name': name,
            'namespace': namespace,
        },
        'spec': {
            'parallelism': parallelism,
            'script': {
                'configMap': {
                    'name': test_file_resource,
                    'file': os.path.basename(from_file)
                },
            },
            'arguments': '-o output-prometheus-remote' if prometheus_url != '' else '',
            'runner': {
                'image': 'ghcr.io/mcandeia/k6-custom:latest',
                'metadata': {
                    'annotations': {
                        'dapr.io/app-id': 'tester-app',
                        'dapr.io/enabled': 'true' if add_dapr else 'false'
                    },
                },
                'env': [{ 'name': 'K6_PROMETHEUS_REMOTE_URL', 'value': prometheus_url}, { 'name': 'K6_PROMETHEUS_USER', 'value': prometheus_user}, { 'name': 'K6_PROMETHEUS_PASSWORD', 'value': prometheus_password}] if prometheus_url != '' else []
            }
        }
    }))

    k8s_resource(new_name=name,
            objects=[name+':k6'], resource_deps=resource_deps + [test_file_resource, k6_local_resource_name])
    # local_resource("wait-for-tests-finished", cmd="kubectl wait --namespace k6-operator-system --for=condition=ready pod --all && kubectl wait --for=jsonpath='{.status.stage}'=finished k6 --all -n %s" % (namespace), resource_deps=['install_k6', k6_test_name])
    # collect results

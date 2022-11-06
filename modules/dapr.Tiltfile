# -*- mode: Python -*

load('ext://helm_resource', 'helm_resource', 'helm_repo')

dapr_api_version = 'dapr.io/v1alpha1'
def install(version='latest', namespace='dapr-system', wait=False):
    helm_repo('dapr', 'https://dapr.github.io/helm-charts/')
    helm_resource('dapr-install', 'dapr/dapr', namespace=namespace, flags=['--version='+version,'--wait', '--create-namespace'])

def component_create(name, type, namespace, version='v1', metadata={}):
    metadata_as_list = []
    for item in metadata.items():
        metadata_as_list.append({'name': str(item[0]), 'value': str(item[1])})
    return encode_yaml({
        'apiVersion': dapr_api_version,
        'kind': 'Component',
        'metadata': {
            'name': name,
            'namespace': namespace
        },
        'spec': {
            'type': type,
            'version': version,
            'metadata': metadata_as_list,
        }
    })

def with_annotations(deployment, app_name, app_port):
    deployment['spec']['template']['metadata']['annotations'] = { 'dapr.io/app-id': app_name, 'dapr.io/enabled': 'true', 'dapr.io/app-port': '8080'}
    return deployment
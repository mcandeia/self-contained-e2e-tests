# -*- mode: Python -*

load('ext://helm_resource', 'helm_resource', 'helm_repo')

dapr_api_version = 'dapr.io/v1alpha1'

def install(version='latest', namespace='dapr-system', wait=True):
    flags = ['--version='+version, '--create-namespace']
    if wait:
        flags.append('--wait')
    helm_repo('dapr', 'https://dapr.github.io/helm-charts/')
    helm_resource('dapr-install', 'dapr/dapr', namespace=namespace, flags=flags)

def component_create(name, type, namespace, version='v1', metadata={}, scopes=[]):
    k8s_yaml(component_yaml(name, type, namespace, version, metadata, scopes))
    k8s_resource(new_name=name,
             objects=[name+':component'], resource_deps=['dapr-install', 'dapr'])

def component_yaml(name, type, namespace, version='v1', metadata={}, scopes=[]):
    metadata_as_list = []
    for item in metadata.items():
        metadata_as_list.append({'name': str(item[0]), 'value': str(item[1])})
    component = {
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
    }
    if len(scopes) > 0:
        component['scopes'] = scopes
    return encode_yaml(component)

def with_annotations(deployment, app_name, app_port):
    deployment['spec']['template']['metadata']['annotations'] = { 'dapr.io/app-id': app_name, 'dapr.io/enabled': 'true', 'dapr.io/app-port': '8080'}
    return deployment
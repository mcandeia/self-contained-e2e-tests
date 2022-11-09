# -*- mode: Python -*
load('ext://deployment', 'deployment_yaml')
load('ext://helm_resource', 'helm_resource', 'helm_repo')
load('ext://ko', 'ko_build')

dapr_api_version = 'dapr.io/v1alpha1'
dapr_helm_resource_name = 'dapr'
dapr_helm_repo = 'dapr-repo'

def dapr_helm_resource(version='latest', namespace='dapr-system'):
    helm_repo(dapr_helm_repo, 'https://dapr.github.io/helm-charts/')
    helm_resource(dapr_helm_resource_name, '%s/dapr' % (dapr_helm_repo), namespace=namespace, flags=['--version='+version, '--create-namespace'])

def component_create(name, type, namespace, version='v1', metadata={}, scopes=[]):
    k8s_yaml(component_yaml(name, type, namespace, version, metadata, scopes))
    k8s_resource(new_name=name,
             objects=[name+':component'], resource_deps=[dapr_helm_resource_name])

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


def dapr_app_resource(app_name='dapr-app', app_path=None, namespace='default', app_port='8080', ko=False, resource_deps=[]):
    if not app_path:
        app_path = './apps/%s' % (app_name)
    build_app = ko_build if ko else docker_build
    build_app(app_name, app_path)
    deployment = decode_yaml(deployment_yaml(app_name, namespace=namespace, port=app_port))
    k8s_yaml(encode_yaml(with_annotations(deployment, app_name, app_port)))
    k8s_resource(app_name, resource_deps=resource_deps)

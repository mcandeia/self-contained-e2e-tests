# -*- mode: Python -*

load('ext://deployment', 'deployment_create', 'deployment_yaml')
load('ext://configmap', 'configmap_create')
load('ext://namespace', 'namespace_create')
load('ext://ko', 'ko_build')

load('./modules/k6.Tiltfile', install_k6='install', 'run_load_test')
load('./modules/dapr.Tiltfile', install_dapr='install', 'component_create',  'with_annotations')
update_settings(k8s_upsert_timeout_secs = 300)

dapr_namespace_config = 'dapr_namespace'
dapr_version_config = 'dapr_version'
tests_config = 'tests'
run_grafana_config = 'run_grafana'
add_prometheus_config = 'add_prometheus'
default_registry_config = 'default_registry'
infra_config = 'infra'
config.define_string(infra_config)
config.define_string(dapr_namespace_config)
config.define_string(dapr_version_config)
config.define_string(default_registry_config)
config.define_string_list(tests_config)
config.define_bool(run_grafana_config)
config.define_bool(add_prometheus_config)

cfg = config.parse()

default_registry(cfg[default_registry_config])
load_dynamic('./infra/%s/Tiltfile' % cfg[infra_config])

## Create VCluster.

test_namespace = cfg[dapr_namespace_config]
namespace_create(test_namespace)
install_k6(run_grafana=cfg[run_grafana_config], add_prometheus=cfg[add_prometheus_config])
install_dapr(version=cfg[dapr_version_config])

for test in cfg.get(tests_config, []):
    test_module = load_dynamic('./tests/' + test + '/Tiltfile')
    test_module['run'](test_namespace)
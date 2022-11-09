# -*- mode: Python -*

load('./dapr.Tiltfile', 'component_create')
load('ext://deployment', 'deployment_create')

redis_name = 'redis'
redis_port = '6379'

def redis_resource(name=redis_name, namespace=None):
    deployment_create(
        name,
        ports=redis_port,
        namespace=namespace,
        readiness_probe={'exec':{'command':['redis-cli','ping']}}
    )

    redis_svc_address = '%s.%s.svc.cluster.local:%s' % (name, namespace, str(redis_port))
    return redis_svc_address

def redis_statestore_resource(name='statestore', namespace='default', redis_url='', redis_password='', actor_statestore=False, scopes=[]):
    component_create(name=name, type='state.redis', version='v1', namespace=namespace, metadata={'redisHost':redis_url, 'redisPassword': redis_password,'actorStateStore': 'true' if actor_statestore else 'false' }, scopes= None if len(scopes) == 0 else scopes)


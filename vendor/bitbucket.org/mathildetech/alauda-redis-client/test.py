import os
import alauda_redis_client
os.environ["REDIS_TYPE_READER"] = 'normal'
os.environ["REDIS_HOST_READER"] = '127.0.0.1'
os.environ["REDIS_PORT_READER"] = '6379'
os.environ["REDIS_DB_NAME_READER"] = '0'

client = alauda_redis_client.RedisClientFactory.get_client_by_key('READER')
print client.incr('liaojian')
print client.incr('liaojian')


# import os
# import alauda_redis_client
# os.environ["REDIS_TYPE_READER"] = 'cluster'
# os.environ["REDIS_STARTUP_NODES_READER"] = '127.0.0.1:7000'
# os.environ["REDIS_READONLY_MODE_READER"] = 'true'
#
#
# client = alauda_redis_client.RedisClientFactory.get_client_by_key('READER')
# print client.incr('liaojian')
# print client.incr('liaojian')

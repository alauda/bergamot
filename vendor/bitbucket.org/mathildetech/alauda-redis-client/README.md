#README

##作用
封装了normal redis和cluster redis的client，只需要在环境中设置相应的环境变量，并引入这个module就能直接获取到对应的client

##安装
代码已经上传到我们的pypi上了。
###单独下载
    sudo pip install --trusted-host pypi.alauda.io --extra-index-url http://mathildetech:Mathilde1861@pypi.alauda.io/simple/ alauda_redis_client
###Dockerfile
以jakiro为例
    在文件 /jakiro/requirements-alauda.txt添加
    alauda_redis_client>=0.1.0

    在Dockerfile中添加
    RUN pip install --trusted-host pypi.alauda.io --extra-index-url \
    http://mathildetech:Mathilde1861@pypi.alauda.io/simple/ -r /jakiro/requirements-alauda.txt && \
    rm -rf /root/.cache/pip/


##使用方法
    #引用
	from alauda_redis_client import RedisClientFactory
	#获取client，"CLIENT1" 和 "CLIENT2"是一个KEY，用来读取不同的环境变量，因为我们允许读取多组redis配置，生成不同的client
	redis_client1 = RedisClientFactory.get_client_by_key("CLIENT1")
	redis_client2 = RedisClientFactory.get_client_by_key("CLIENT2")

###环境变量
下面环境变量的后缀，CLIENT1和CLIENT2都是可以自定义的**key**
####Normal REDIS
	REDIS_TYPE_CLIENT1=normal

	REDIS_HOST_CLIENT1=127.0.0.1

	REDIS_PORT_CLIENT1=6379

	REDIS_DB_NAME_CLIENT1=0

	#可选参数
	REDIS_DB_PASSWORD_CLIENT1=password

	#可选参数，连接池的最大连接数
	REDIS_MAX_CONNECTIONS_CLIENT1=32

	#可选参数，目前只支持BlockingConnectionPool
	REDIS_CONNECT_POOL_CLASS_CLIENT1=BlockingConnectionPool

	#可选参数，为你当前使用的key加上前缀，默认为空''
	REDIS_KEY_PREFIX_CLIENT1=docker_
####Cluster REDIS
	REDIS_TYPE_CLIENT2=cluster

	REDIS_STARTUP_NODES_CLIENT2=127.0.0.1:7000,127.0.0.1:7001

	#可选参数，默认是false
	REDIS_READONLY_MODE_CLIENT2=true

	#可选参数，连接池的最大连接数
	REDIS_MAX_CONNECTIONS_CLIENT2=32

	#可选参数，为你当前使用的key加上前缀，默认为空''
	REDIS_KEY_PREFIX_CLIENT2=docker_

##client 支持的方法
现在没有支持所有的redis client方法。只封装了jakiro用的。

	def incr(self, name, amount=1):

    def expire(self, name, time):

    def get(self, name):

    def set(self, name, value, ex=None, px=None, nx=False, xx=False):

    def hget(self, name, key):

    def hset(self, name, key, value):

    def hdel(self, name, *keys):

    def hmget(self, name, keys, *args):

    def hmset(self, name, mapping):

    def hgetall(self, name):

    def delete(self, *names):

    def setex(self, name, value, time):

    def blpop(self, names, timeout=0):

    def rpush(self, name, *values):

    def hscan_iter(self, name, match=None, count=None):

    def pipeline(self, transaction=False, shard_hint=None):

    def scan(self, cursor=0, match=None, count=None):

    def sadd(self, name, *values):

    def smembers(self, name):

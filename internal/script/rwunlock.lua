local key = KEYS[1]
local subkey = ARGV[1]
redis.call("HDEL", key, subkey)
local hlen = redis.call("HLEN", key)
if hlen <= 1 then
	redis.call("DEL", key)
end
return "OK"
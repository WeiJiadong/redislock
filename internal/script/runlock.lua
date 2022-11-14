local rkey = KEYS[1]
local subkey = "r"
redis.call("HDEL", rkey, subkey)
local hlen = redis.call("HLEN", rkey)
if hlen == 0 then
	redis.call("DEL", rkey)
end
return "ok"
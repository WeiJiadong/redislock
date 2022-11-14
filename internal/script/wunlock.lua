local wkey = KEYS[1]
local subkey = "w"
redis.call("HDEL", wkey, subkey)
return "ok"
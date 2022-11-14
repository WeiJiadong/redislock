local wkey = KEYS[1]
local subkey = "w"
local seconds = tonumber(ARGV[2])
local val = ARGV[1]
local wval = redis.call("HGET", wkey, subkey)
if wval ~= false then
	return "write lock conflict"
end
redis.call("HSET", wkey, subkey, val)
redis.call("EXPIRE", wkey, seconds)
return "ok"
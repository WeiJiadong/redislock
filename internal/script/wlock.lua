local wkey = KEYS[1]
local subkey = "w"
local seconds = tonumber(ARGV[2])
local wval = redis.call("HGET", wkey, subkey)
if wval ~= false and wval == ARGV[2] then
	return "FAIL"
end
redis.call("HSET", wkey, subkey, ARGV[1])
redis.call("EXPIRE", wkey, seconds)
return "OK"
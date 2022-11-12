local wkey = KEYS[1]
local subkey = "w"
local seconds = tonumber(ARGV[2])
local wval = redis.call("HGET", wkey, subkey)
if wval ~= false and wval == ARGV[2] then
	return "FAIL"
end
local rkey = KEYS[1]
subkey = "r"
redis.call("HSET", rkey, subkey, ARGV[1])
redis.call("EXPIRE", rkey, seconds)
return "OK"
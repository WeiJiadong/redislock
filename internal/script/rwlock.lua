local rwkey = "_rw"
local key = KEYS[1]
local rwval = ARGV[1]
local seconds = tonumber(ARGV[2])
local subkey = ARGV[3]
local rwnow = redis.call("HGET", key, rwkey)
if rwnow ~= false then
	if rwnow ~= rwval then
		return "FAIL"
	end
end
redis.call("HMSET", key, rwkey, rwval, subkey, "")
redis.call("EXPIRE", key, seconds)
return "OK"
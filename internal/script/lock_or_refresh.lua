local data = redis.call("get", KEYS[1])
if data == false then
	return redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
else
	if data == ARGV[1] then
		return redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
	else
		return "FAIL"
	end
end
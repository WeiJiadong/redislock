if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
else
	return "FAIL"
end
if redis.call("GET", KEYS[1]) == ARGV[1] then
	redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
	return "ok"
else
	return "lock not exist"
end
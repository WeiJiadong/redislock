if redis.call("GET", KEYS[1]) == ARGV[1] then
	redis.call("DEL", KEYS[1])
end
return "ok"
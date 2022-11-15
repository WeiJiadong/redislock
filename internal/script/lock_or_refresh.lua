local data = redis.call("GET", KEYS[1])
if data == false then
	redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
else
	if data == ARGV[1] then
		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
	else
		return "write lock conflict"
	end
end
return "ok"
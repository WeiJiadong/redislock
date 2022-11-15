local val = redis.call("GET", KEYS[1])
if val == false then
    redis.call("SET", KEYS[1], ARGV[2])
else 
    if val == ARGV[1] then
        redis.call("SET", KEYS[1], ARGV[2])
    else
        return "cas conflict"
    end
end
return "ok"
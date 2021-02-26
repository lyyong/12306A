/*
* @Author: 余添能
* @Date:   2021/2/20 12:28 下午
 */
package init_redis

import "github.com/go-redis/redis"

func CreateScriptBuyTicket() *redis.Script {
	script := redis.NewScript(`
		local key=tostring(KEYS[1])
		local min=tonumber(ARGV[1])
		local tickets = redis.call("ZRangeByScore", key, min,1000,"withscores")

		-- 表的大小为0，表示没有元素
		if #tickets==0 then 
		   return {0,0,0}
		end
		
		-- 奇数key 对应val是　value
    	-- 偶数key 对应val是　score

		local endStation
		local carriageAndSeatNo
		for i,v in pairs(tickets) do
			if i==1 then
				carriageAndSeatNo=v
			end
			if i==2 then
				endStation=v
				break;
			end
		end
		print(carriageAndSeatNo)
		local res=redis.call("ZREM",key,carriageAndSeatNo)
		return {endStation,carriageAndSeatNo,res}
	`)
	return script
}



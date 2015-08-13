module("woLongShanZhuang", package.seeall)

function OnInit(s)
	print("欢迎来到卧龙山庄.")
	-- s:Get_data()["lolo"] = "lolo"
	-- print(s:Get_data()["lolo"])
end

function OnHeartBeat( s)
	-- print(s:Get_thread():Get_thread_id().."  卧龙山庄心跳 "..os.time())
	
	-- print(s:Get_data()["lolo"])
	-- s:Get_data()["lolo"] = "lolo"..os.time()
end
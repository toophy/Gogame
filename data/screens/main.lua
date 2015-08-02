
require("data/screens/common")

function OnInitScreen( a )
	a:Add_screen("阿拉斯加", 1)
	a:Add_screen("上海", 1)
	a:Add_screen("厦门", 1)
	a:Add_screen("阿拉斯加2", 1)

	local btime = os.time()

	for i=1,10000000 do
	-- 	print(a:Get_randNum())
	-- print(a:Get_randNum()/98)
	-- local x = 98765432198765
	-- local y = 1000000
	-- local z = x + y
	a:Set_randNum(a:Get_randNum()/98)
	-- print(a:Get_randNum())
	end

	local etime = os.time()

	print(etime-btime)

	return 1
end

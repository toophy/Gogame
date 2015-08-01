
require("data/screens/common")

-- print("Lua : hello go game")
-- -- print("Lua call go sum 99+1 = "..mysum(99,1))

-- function OnInitScreen( a )
-- 	common.PrintTable(_G)

-- 	a:Add_screen("阿拉斯加", 1)
-- 	-- a:Add_screen("上海", 1)
-- 	-- a:Add_screen("厦门", 1)
-- 	-- a:Add_screen("阿拉斯加2", 1)


-- end

        function OnInitScreen( a )
        	person.name(a, "xx")
        	
        	print("\n\n")
        	common.PrintTable( a )
        	print("\n\n")

        	 print(a:name()) -- "Steeve"
        	a:name("Alicex")
        	print(a:name()) -- "Alice"
        end
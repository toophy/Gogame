module("main", package.seeall)

function OnInitScreen( a )
	
	a:Add_screen("阿拉斯加", 1)
	a:Add_screen("上海", 1)
	a:Add_screen("厦门", 1)
	a:Add_screen("阿拉斯加2", 1)

	print(common.GetCurrDir())

	return 1
end

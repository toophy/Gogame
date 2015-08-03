module("common", package.seeall)

-- 打印table
function PrintTable(root)
	local cache = {  [root] = "." }
	local function _dump(t,space,name)
		local temp = {}
		for k,v in pairs(t) do
			local key = tostring(k)
			if cache[v] then
				table.insert(temp,"* " .. key .. " {" .. cache[v].."}")
			elseif type(v) == "table" then
				local new_key = name .. "." .. key
				cache[v] = new_key
				table.insert(temp,"+ " .. key .. _dump(v,space .. (next(t,k) and "|" or " " ).. string.rep(" ",#key),new_key))
			else
				table.insert(temp,"- " .. key .. " [" .. tostring(v).."]")
			end
		end
		return table.concat(temp,"\n"..space)
	end
	print(_dump(root, "",""))
end

-- 获取当前目录
function GetCurrDir()
	obj = io.popen("cd")  --如果不在交互模式下，前面可以添加local 
	path = obj:read("*all"):sub(1,-2)    --path存放当前路径
	obj:close()   --关掉句柄

	return path
end


function filter(_, _, record)
	if record.log.level == nil then
		return -1, 0, 0
	end
	return 2, 0, record.log
end

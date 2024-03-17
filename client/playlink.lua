local base64 = import "b64"

response = ""
readingResp = ""
version = "0.0.1"
initalized = false

function split(inputstr, sep)
    if sep == nil then
        sep = "%s"
    end
    local t = {}
    for str in string.gmatch(inputstr, "([^" .. sep .. "]+)") do
        table.insert(t, str)
    end
    return t
end

local function process(packet, callback)
    if packet:match("^playlink|init_ok|%d%.%d%.%d") ~= nil then
        print("recv: it worked")
        print("playlink|connected|0.0.1")
        initalized = true
        return
    end

    if not initalized then
        error("You must run playlink.init() before using any other functions.")
    end

    if readingResp == "reading" and packet:find("endresp") then
        readingResp = ""
        callback(response)
        response = ""
    end

    if readingResp == "reading" then
        response = response .. base64.decode(packet)
        print("recv: currentMsg: " .. base64.decode(packet))
    end

    if readingResp == "init" and packet:find("startresp") then
        readingResp = "reading"
    end

    if packet:find("response") ~= nil then
        readingResp = "init"
    end
end

function get(url)
    if not initalized then
        error("You must run playlink.init() before using any other functions.")
    end
    print("playlink|request|" .. version .. "|GET|" .. url)
end

function init()
    print("playlink|init|" .. version)
end

return {
    process = process,
    get = get,
    init = init,
}

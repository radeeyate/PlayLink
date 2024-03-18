local base64 = import "b64"

response = ""
readingResp = ""
identifier = ""
ponged = false
status_code = 0
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

local function process(packet)
    if packet:match("^playlink|pong|%d%.%d%.%d") ~= nil then
        ponged = true
        return
    end

    if packet:match("^playlink|init_ok|%d%.%d%.%d") ~= nil then
        print("recv: it worked")
        print("playlink|connected|0.0.1")
        playlink.onConnection()
        initalized = true
        return
    end

    if not initalized then
        error("You must run playlink.init() before using any other functions.")
    end

    if readingResp == "reading" and packet:find("endresp") then
        readingResp = ""
        playlink.onResponse(
            {
                body = response,
                status_code = status_code,
                identifier = identifier,
            }
        )
        response = ""
        identifier = ""
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
        print("recv: resp: " .. packet)
        status_code = split(packet, "|")[3]
        identifier = base64.decode(split(packet, "|")[4])
    end
end

function request(method, url, identifier)
    if not initalized then
        error("You must run playlink.init() and be connected before using any other functions.")
    else
        if identifier == nil then
            identifier = "none"
        end
        print("playlink|request|" .. version .. "|" .. method .. "|" .. url .. "|" .. identifier)
    end
end

function init()
    print("playlink|init|" .. version)
end

playlink = {
    process = process,
    get = function(url, identifier) request("GET", url, identifier) end,
    init = init,
    onResponse = function(response) end,
    onConnection = function() end
}

return playlink

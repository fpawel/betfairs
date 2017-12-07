module Utils exposing (..)


websocketURL : String -> String -> String
websocketURL protocol host =
    if not (String.startsWith "http" protocol) then
        "wrong protocol `"
            ++ protocol
            ++ "wrong name of the protocol - expected a string that starts with `http`"
            |> Debug.crash
    else
        "ws"
            ++ (String.dropLeft 4 protocol)
            ++ "//"
            ++ host

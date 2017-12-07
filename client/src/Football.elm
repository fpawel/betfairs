module Football exposing (..)

import Json.Decode as D
import Html exposing (Html, text)
import Material.Table exposing (tr, td, th, tbody, table)


type alias Game =
    { id : Int
    , home : String
    , away : String
    , scoreHome : Int
    , scoreAway : Int
    , time : String
    , inplay : Bool
    }



-- VIEW


renderGames : List Game -> Html a
renderGames games =
    let
        headRow =
            List.map
                (\x -> th [] [ text x ])
                [ "Дома"
                , "Счёт"
                , "В гостях"
                , "Время"
                ]

        rows =
            List.map
                (\game ->
                    let
                        strScore =
                            if game.inplay then
                                toString game.scoreHome ++ " - " ++ toString game.scoreAway
                            else
                                ""
                    in
                        [ td [] [ text game.home ]
                        , td [] [ text strScore ]
                        , td [] [ text game.away ]
                        , td [] [ text game.time ]
                        ]
                )
                games
    in
        table
            []
            [ tbody [] (List.map (tr []) (headRow :: rows)) ]



-- DECODERS


parseGames : String -> Result String (List Game)
parseGames =
    D.decodeString <| D.list decoderGame


decoderGame : D.Decoder Game
decoderGame =
    D.map7 Game
        (D.field "id" D.int)
        (D.field "home" D.string)
        (D.field "away" D.string)
        (D.field "score_home" D.int)
        (D.field "score_away" D.int)
        (D.field "time" D.string)
        (D.field "in_play" D.bool)

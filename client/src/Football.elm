module Football exposing (..)

import Json.Decode as D
import Json.Decode.Pipeline exposing (decode, required)
import Html exposing (Html, text)
import Material.Table exposing (tr, td, th, tbody, table)


type alias Game =
    { id : Int
    , home : String
    , away : String
    , competition : String
    , country : String
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
                , "Страна"
                , "Чемпионат"
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
                        , td [] [ text game.country ]
                        , td [] [ text game.competition ]
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
    decode Game
        |> required "id" D.int
        |> required "home" D.string
        |> required "away" D.string
        |> required "competition" D.string
        |> required "country" D.string
        |> required "score_home" D.int
        |> required "score_away" D.int
        |> required "time" D.string
        |> required "in_play" D.bool

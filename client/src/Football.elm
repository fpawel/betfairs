module Football exposing (Model, Msg, init, update, view, subscriptions)

import String
import Json.Decode as D
import Json.Decode.Pipeline exposing (decode, required, hardcoded)
import Html exposing (Html, text)
import Material.Table exposing (tr, td, th, tbody, thead, table)
import Material.Spinner
import Utils exposing (..)
import WebSocket


-- MODEL


type alias Model =
    { protocol : String
    , host : String
    , games : List Game
    , sort : Sort
    }


type Msg
    = NewGames (List Game)
    | Reorder Sort


type Sort
    = SortOrder
    | SortCompetition
    | SortCountry


type alias Game =
    { order : Int
    , id : Int
    , home : String
    , away : String
    , competition : String
    , country : String
    , scoreHome : Int
    , scoreAway : Int
    , time : String
    , inplay : Bool
    }


init : String -> String -> ( Model, Cmd Msg )
init protocol host =
    ( { games = []
      , sort = SortOrder
      , protocol = protocol
      , host = host
      }
    , Cmd.none
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NewGames newGames ->
            ( { model | games = newGames }, Cmd.none )

        Reorder newSort ->
            ( { model | sort = newSort }, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    WebSocket.listen
        (websocketURL model.protocol model.host ++ "/football")
        (\str ->
            case parseGames str of
                Ok y ->
                    NewGames y

                Err err ->
                    Debug.crash err
        )



-- VIEW


view : Model -> Html Msg
view { games } =
    let
        hasCompetition =
            hasNotEmpty .competition games

        hasCountry =
            hasNotEmpty .country games

        headRow =
            [ th [] [ text "№" ]
            , th [] [ text "Дома" ]
            , th [] [ text "Счёт" ]
            , th [] [ text "В гостях" ]
            , th [] [ text "Время" ]
            ]
                ++ (if hasCountry then
                        [ th [] [ text "Страна" ] ]
                    else
                        []
                   )
                ++ (if hasCompetition then
                        [ th [] [ text "Чемпионат" ] ]
                    else
                        []
                   )

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
                        [ td [ Material.Table.numeric ] [ text <| toString (game.order + 1) ]
                        , td [] [ text game.home ]
                        , td [] [ text strScore ]
                        , td [] [ text game.away ]
                        , td [] [ text game.time ]
                        ]
                            ++ (if hasCountry then
                                    [ td [] [ text game.country ] ]
                                else
                                    []
                               )
                            ++ (if hasCompetition then
                                    [ td [] [ text game.competition ] ]
                                else
                                    []
                               )
                            |> tr []
                )
                games
    in
        if List.isEmpty games then
            Material.Spinner.spinner [ Material.Spinner.active True ]
        else
            table
                []
                [ thead [] [ tr [] headRow ]
                , tbody [] rows
                ]



-- DECODERS


parseGames : String -> Result String (List Game)
parseGames str =
    Result.map
        (List.indexedMap (\n x -> { x | order = n }))
        (D.decodeString (D.list decoderGame) str)


decoderGame : D.Decoder Game
decoderGame =
    decode Game
        |> hardcoded 0
        |> required "id" D.int
        |> required "home" D.string
        |> required "away" D.string
        |> required "competition" D.string
        |> required "country" D.string
        |> required "score_home" D.int
        |> required "score_away" D.int
        |> required "time" D.string
        |> required "in_play" D.bool



-- HELPERS


hasNotEmpty : (a -> String) -> List a -> Bool
hasNotEmpty f =
    List.filter (f >> String.isEmpty >> not)
        >> List.head
        >> Maybe.map (\_ -> True)
        >> Maybe.withDefault False

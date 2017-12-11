module Football exposing (Model, Msg, init, update, view, subs)

import String
import Json.Decode as D
import Json.Decode.Pipeline exposing (decode, required, hardcoded)
import Html exposing (Html, text, div)
import Html.Attributes
import Material.Table exposing (tr, td, th, tbody, thead, table)
import Material.Spinner as Spinner
import Utils exposing (..)
import WebSocket
import Material.Grid as Grid
import Material.Options as Options


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
    , mainPrices : List Float
    , totalMatched : Float
    , totalAvailable : Float
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


subs : (Msg -> msg) -> Model -> Sub msg
subs toMsg model =
    Sub.map toMsg (subscriptions model)


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
        noGames =
            List.isEmpty games

        spinner =
            if noGames then
                [ Spinner.spinner [ Spinner.active True ] ]
            else
                []

        gamesTable =
            if noGames then
                []
            else
                [ renderGamesTable games ]
    in
        Grid.grid []
            [ Grid.cell [ Grid.size Grid.All 3 ] spinner
            , Grid.cell [ Grid.size Grid.All 9 ] gamesTable
            ]


renderGamesTable : List Game -> Html Msg
renderGamesTable games =
    let
        hasCompetition =
            hasNotEmpty .competition games

        hasCountry =
            hasNotEmpty .country games
    in
        table
            []
            [ thead [] <| renderGamesHeaderRow
            , tbody [] <| List.map renderGameTableRow games
            ]


numToStr : number -> String
numToStr x =
    if x == 0 then
        ""
    else
        toString x


renderGameTableRow : Game -> Html Msg
renderGameTableRow game =
    let
        pricesSection =
            game.mainPrices
                |> List.map numToStr
                |> List.map
                    (text
                        >> List.singleton
                        >> td []
                    )
    in
        [ td [ Material.Table.numeric ] [ text <| toString (game.order + 1) ]
        , td [] [ text game.home ]
        , td []
            [ (if game.inplay then
                toString game.scoreHome ++ " - " ++ toString game.scoreAway
               else
                ""
              )
                |> text
            ]
        , td [] [ text game.away ]
        , td [] [ text game.time ]
        ]
            ++ pricesSection
            ++ [ td [] [ text <| numToStr <| game.totalMatched ]
               , td [] [ text <| numToStr <| game.totalAvailable ]
               , td [] [ text game.country ]
               , td [] [ text game.competition ]
               ]
            |> tr []


renderGamesHeaderRow : List (Html msg)
renderGamesHeaderRow =
    let
        colspan2 =
            [ Options.attribute <| Html.Attributes.colspan 2 ]
    in
        [ th [] [ text "№" ]
        , th [] [ text "Дома" ]
        , th [] [ text "Счёт" ]
        , th [] [ text "В гостях" ]
        , th [] [ text "Время" ]
        , th colspan2 [ text "П1" ]
        , th colspan2 [ text "Н" ]
        , th colspan2 [ text "П2" ]
        , th [] [ text "В паре" ]
        , th [] [ text "Не в паре" ]
        , th [] [ text "Страна" ]
        , th [] [ text "Чемпионат" ]
        ]
            |> tr []
            |> List.singleton



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
        |> required "main_prices" (D.list D.float)
        |> required "total_matched" D.float
        |> required "total_available" D.float



-- HELPERS


hasNotEmpty : (a -> String) -> List a -> Bool
hasNotEmpty f =
    List.filter (f >> String.isEmpty >> not)
        >> List.head
        >> Maybe.map (\_ -> True)
        >> Maybe.withDefault False

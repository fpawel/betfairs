module Football exposing (Model, Msg, init, update, view, subs)

import String
import Json.Decode as D
import Json.Decode.Pipeline exposing (decode, required, hardcoded, optional)
import Html exposing (Html, text, div)
import Material.Table exposing (tr, td, th, tbody, thead, table)
import Material.Spinner as Spinner
import Utils exposing (..)
import WebSocket
import Dict exposing (Dict)
import Set exposing (Set)


-- MODEL


type alias Model =
    { protocol : String
    , host : String
    , games : List Game
    , sort : Sort
    }


type Msg
    = NewGamesChanges GamesChanges
    | Reorder Sort


type Sort
    = SortOrder
    | SortCompetition
    | SortCountry


type alias Game =
    { id : Int
    , home : String
    , away : String
    , order : Int
    , competition : String
    , country : String
    , scoreHome : Int
    , scoreAway : Int
    , time : String
    , inplay : Bool
    , winBack : Float
    , winLay : Float
    , drawBack : Float
    , drawLay : Float
    , loseBack : Float
    , loseLay : Float
    , totalMatched : Float
    , totalAvailable : Float
    }


type alias GamesChanges =
    { new : List Game
    , out : List Int
    , upd : List GameChanges
    , reset : Bool
    }


type alias GameChanges =
    { id : Int
    , order : Maybe Int
    , competition : Maybe String
    , country : Maybe String
    , scoreHome : Maybe Int
    , scoreAway : Maybe Int
    , time : Maybe String
    , inplay : Maybe Bool
    , winBack : Maybe Float
    , winLay : Maybe Float
    , drawBack : Maybe Float
    , drawLay : Maybe Float
    , loseBack : Maybe Float
    , loseLay : Maybe Float
    , totalMatched : Maybe Float
    , totalAvailable : Maybe Float
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
        NewGamesChanges gamesChanges ->
            ( { model | games = updateGames model.games gamesChanges }, Cmd.none )

        Reorder newSort ->
            ( { model | sort = newSort }, Cmd.none )


updateGames : List Game -> GamesChanges -> List Game
updateGames games { new, out, upd, reset } =
    if reset then
        new
    else
        let
            outSet =
                out
                    |> Set.fromList

            newSet =
                Set.fromList (List.map .id new)

            updM =
                upd
                    |> List.map (\x -> ( x.id, x ))
                    |> Dict.fromList

            isJust =
                Maybe.map (\_ -> True)
                    >> Maybe.withDefault False

            play =
                games
                    |> List.filter
                        (\{ id } ->
                            (not <| Set.member id outSet)
                                && (not <| Set.member id newSet)
                        )
                    |> List.map
                        (\x ->
                            Dict.get x.id updM
                                |> Maybe.map (updateGame x)
                                |> Maybe.withDefault x
                        )
        in
            new
                ++ play
                |> List.sortBy .order


updateGame : Game -> GameChanges -> Game
updateGame x y =
    let
        comp fy fx =
            case fy y of
                Just t ->
                    t /= fx x

                _ ->
                    False
    in
        { x
            | order = Maybe.withDefault x.order y.order
            , competition = Maybe.withDefault x.competition y.competition
            , country = Maybe.withDefault x.country y.country
            , scoreHome = Maybe.withDefault x.scoreHome y.scoreHome
            , scoreAway = Maybe.withDefault x.scoreAway y.scoreAway
            , time = Maybe.withDefault x.time y.time
            , inplay = Maybe.withDefault x.inplay y.inplay
            , winBack = Maybe.withDefault x.winBack y.winBack
            , winLay = Maybe.withDefault x.winLay y.winLay
            , loseBack = Maybe.withDefault x.loseBack y.loseBack
            , loseLay = Maybe.withDefault x.loseLay y.loseLay
            , drawBack = Maybe.withDefault x.drawBack y.drawBack
            , drawLay = Maybe.withDefault x.drawLay y.drawLay
            , totalMatched = Maybe.withDefault x.totalMatched y.totalMatched
            , totalAvailable = Maybe.withDefault x.totalAvailable y.totalAvailable
        }



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
                    NewGamesChanges y

                Err err ->
                    Debug.crash err
        )



-- VIEW


view : Model -> Html Msg
view { games } =
    if List.isEmpty games then
        Spinner.spinner [ Spinner.active True ]
    else
        renderGamesTable games


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
    , td [] [ text <| numToStr <| game.winBack ]
    , td [] [ text <| numToStr <| game.winLay ]
    , td [] [ text <| numToStr <| game.drawBack ]
    , td [] [ text <| numToStr <| game.drawLay ]
    , td [] [ text <| numToStr <| game.loseBack ]
    , td [] [ text <| numToStr <| game.loseLay ]
    , td [] [ text <| numToStr <| game.totalMatched ]
    , td [] [ text <| numToStr <| game.totalAvailable ]
    , td [] [ text game.country ]
    , td [] [ text game.competition ]
    ]
        |> tr []


renderGamesHeaderRow : List (Html msg)
renderGamesHeaderRow =
    [ th [] [ text "№" ]
    , th [] [ text "Дома" ]
    , th [] [ text "Счёт" ]
    , th [] [ text "В гостях" ]
    , th [] [ text "Время" ]
    , th [] [ text "П1+" ]
    , th [] [ text "П1-" ]
    , th [] [ text "Н+" ]
    , th [] [ text "Н-" ]
    , th [] [ text "П2+" ]
    , th [] [ text "П2-" ]
    , th [] [ text "В паре" ]
    , th [] [ text "Не в паре" ]
    , th [] [ text "Страна" ]
    , th [] [ text "Чемпионат" ]
    ]
        |> tr []
        |> List.singleton



-- DECODERS


parseGames : String -> Result String GamesChanges
parseGames str =
    D.decodeString decoderGamesChanges str


decoderGame : D.Decoder Game
decoderGame =
    decode Game
        |> required "id" D.int
        |> required "home" D.string
        |> required "away" D.string
        |> required "order" D.int
        |> required "competition" D.string
        |> required "country" D.string
        |> required "score_home" D.int
        |> required "score_away" D.int
        |> required "time" D.string
        |> required "in_play" D.bool
        |> required "win_back" D.float
        |> required "win_lay" D.float
        |> required "draw_back" D.float
        |> required "draw_lay" D.float
        |> required "lose_back" D.float
        |> required "lose_lay" D.float
        |> required "total_matched" D.float
        |> required "total_available" D.float


decoderMaybe : String -> D.Decoder a -> D.Decoder (Maybe a -> b) -> D.Decoder b
decoderMaybe fieldStr d =
    optional fieldStr (D.maybe d) Nothing


decoderGamesChanges : D.Decoder GamesChanges
decoderGamesChanges =
    decode
        GamesChanges
        |> optional "new" (D.list decoderGame) []
        |> optional "out" (D.list D.int) []
        |> optional "upd" (D.list decoderGameCahnges) []
        |> optional "reset" D.bool False


decoderGameCahnges : D.Decoder GameChanges
decoderGameCahnges =
    decode GameChanges
        |> required "id" D.int
        |> decoderMaybe "order" D.int
        |> decoderMaybe "competition" D.string
        |> decoderMaybe "country" D.string
        |> decoderMaybe "score_home" D.int
        |> decoderMaybe "score_away" D.int
        |> decoderMaybe "time" D.string
        |> decoderMaybe "in_play" D.bool
        |> decoderMaybe "win_back" D.float
        |> decoderMaybe "win_lay" D.float
        |> decoderMaybe "draw_back" D.float
        |> decoderMaybe "draw_lay" D.float
        |> decoderMaybe "lose_back" D.float
        |> decoderMaybe "lose_lay" D.float
        |> decoderMaybe "total_matched" D.float
        |> decoderMaybe "total_available" D.float



-- HELPERS


hasNotEmpty : (a -> String) -> List a -> Bool
hasNotEmpty f =
    List.filter (f >> String.isEmpty >> not)
        >> List.head
        >> Maybe.map (\_ -> True)
        >> Maybe.withDefault False

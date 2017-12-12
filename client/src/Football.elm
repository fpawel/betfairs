module Football exposing (Model, Msg, init, update, view, subs)

import String
import Json.Decode as D
import Json.Decode.Pipeline exposing (decode, required, hardcoded, optional)
import Html exposing (..)
import Material.Table as Table
import Material.Spinner as Spinner
import Material.Options as Options
import Utils exposing (..)
import WebSocket
import Dict exposing (Dict)
import Set exposing (Set)
import Material
import Material.Button as Button
import Material.Grid as Grid
import Material.Icon as Icon
import Material.Tooltip as Tooltip
import Material.Dialog as Dialog


-- MODEL


type alias Model =
    { protocol : String
    , host : String
    , games : List Game
    , order : Table.Order
    , sortCol : SortColumn
    , mdl : Material.Model
    }


type Msg
    = NewGamesChanges GamesChanges
    | Reorder SortColumn
    | Mdl (Material.Msg Msg)


type SortColumn
    = SortOrder
    | SortCountry
    | SortCompetition
    | SortTotalMatched
    | SortTotalAvailable


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
      , order = Table.Ascending
      , protocol = protocol
      , host = host
      , sortCol = SortOrder
      , mdl = Material.model
      }
    , Cmd.none
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Mdl msg_ ->
            Material.update Mdl msg_ model

        NewGamesChanges gamesChanges ->
            let
                newGames =
                    model.games
                        |> updateGames gamesChanges
                        |> sortGames model.order model.sortCol
            in
                { model | games = newGames } ! []

        Reorder newSortCol ->
            let
                newOrder =
                    if model.order == Table.Ascending then
                        Table.Descending
                    else
                        Table.Ascending
            in
                { model
                    | order = newOrder
                    , sortCol = newSortCol
                    , games = sortGames newOrder newSortCol model.games
                }
                    ! []


sortGames : Table.Order -> SortColumn -> List Game -> List Game
sortGames order sortCol games =
    let
        sort =
            case ( order, sortCol ) of
                ( Table.Ascending, SortOrder ) ->
                    List.sortBy .order

                ( Table.Descending, SortOrder ) ->
                    List.sortWith (\x y -> reverse (.order x) (.order y))

                ( Table.Ascending, SortCountry ) ->
                    List.sortBy .country

                ( Table.Descending, SortCountry ) ->
                    List.sortWith (\x y -> reverse (.country x) (.country y))

                ( Table.Ascending, SortCompetition ) ->
                    List.sortBy .competition

                ( Table.Descending, SortCompetition ) ->
                    List.sortWith (\x y -> reverse (.competition x) (.competition y))

                ( Table.Ascending, SortTotalMatched ) ->
                    List.sortBy .totalMatched

                ( Table.Descending, SortTotalMatched ) ->
                    List.sortWith (\x y -> reverse (.totalMatched x) (.totalMatched y))

                ( Table.Ascending, SortTotalAvailable ) ->
                    List.sortBy .totalAvailable

                ( Table.Descending, SortTotalAvailable ) ->
                    List.sortWith (\x y -> reverse (.totalAvailable x) (.totalAvailable y))
    in
        sort games


updateGames : GamesChanges -> List Game -> List Game
updateGames { new, out, upd, reset } games =
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
view { games, order, sortCol, mdl } =
    if List.isEmpty games then
        Spinner.spinner [ Spinner.active True ]
    else
        Grid.grid []
            [ Grid.cell
                [ Grid.size Grid.All 11 ]
                [ renderGamesTable order sortCol games ]
            , Grid.cell
                [ Grid.size Grid.All 1 ]
                [ Button.render Mdl
                    [ 0 ]
                    mdl
                    [ Button.icon
                    , Button.fab
                    , Dialog.openOn "click"
                    ]
                    [ div []
                        [ Icon.view "settings" [ Tooltip.attach Mdl [ 0 ] ]
                        , Tooltip.render Mdl
                            [ 0 ]
                            mdl
                            [ Tooltip.left ]
                            [ text "Default tooltip" ]
                        ]
                    , Dialog.view
                        []
                        [ Dialog.title [] [ text "Greetings" ]
                        , Dialog.content []
                            [ p [] [ text "A strange game—the only winning move is not to play." ]
                            , p [] [ text "How about a nice game of chess?" ]
                            ]
                        , Dialog.actions []
                            [ Button.render Mdl
                                [ 0 ]
                                mdl
                                [ Dialog.closeOn "click" ]
                                [ text "Chess" ]
                            , Button.render Mdl
                                [ 0 ]
                                mdl
                                [ Dialog.closeOn "click"
                                ]
                                [ text "GTNW" ]
                            ]
                        ]
                    ]
                ]
            ]


renderGamesTable : Table.Order -> SortColumn -> List Game -> Html Msg
renderGamesTable order sortCol games =
    let
        hasCompetition =
            hasNotEmpty .competition games

        hasCountry =
            hasNotEmpty .country games
    in
        Table.table
            []
            [ Table.thead [] <| renderGamesHeaderRow order sortCol
            , Table.tbody [] <| List.map renderGameTableRow games
            ]


numToStr : number -> String
numToStr x =
    if x == 0 then
        ""
    else
        toString x


renderGameTableRow : Game -> Html Msg
renderGameTableRow game =
    [ Table.td [] [ text <| toString (game.order + 1) ]
    , Table.td [] [ text game.home ]
    , Table.td []
        [ (if game.inplay then
            toString game.scoreHome ++ " - " ++ toString game.scoreAway
           else
            ""
          )
            |> text
        ]
    , Table.td [] [ text game.away ]
    , Table.td [] [ text game.time ]
    , Table.td [] [ text <| numToStr <| game.winBack ]
    , Table.td [] [ text <| numToStr <| game.winLay ]
    , Table.td [] [ text <| numToStr <| game.drawBack ]
    , Table.td [] [ text <| numToStr <| game.drawLay ]
    , Table.td [] [ text <| numToStr <| game.loseBack ]
    , Table.td [] [ text <| numToStr <| game.loseLay ]
    , Table.td [] [ text <| numToStr <| game.totalMatched ]
    , Table.td [] [ text <| numToStr <| game.totalAvailable ]
    , Table.td [] [ text game.country ]
    , Table.td [] [ text game.competition ]
    ]
        |> Table.tr []


renderGamesHeaderRow : Table.Order -> SortColumn -> List (Html Msg)
renderGamesHeaderRow order sortCol =
    let
        sortTH sc str =
            Table.th
                [ Options.onClick (Reorder sc)
                , if sc == sortCol then
                    Table.sorted order
                  else
                    Options.nop
                ]
                [ text str ]
    in
        [ sortTH SortOrder "№"
        , Table.th [] [ text "Дома" ]
        , Table.th [] [ text "Счёт" ]
        , Table.th [] [ text "В гостях" ]
        , Table.th [] [ text "Время" ]
        , Table.th [] [ text "П1+" ]
        , Table.th [] [ text "П1-" ]
        , Table.th [] [ text "Н+" ]
        , Table.th [] [ text "Н-" ]
        , Table.th [] [ text "П2+" ]
        , Table.th [] [ text "П2-" ]
        , sortTH SortTotalMatched "В паре"
        , sortTH SortTotalAvailable "Не в паре"
        , sortTH SortCountry "Страна"
        , sortTH SortCompetition "Чемпионат"
        ]
            |> Table.tr []
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


reverse : comparable -> comparable -> Order
reverse x y =
    case compare x y of
        LT ->
            GT

        GT ->
            LT

        EQ ->
            EQ

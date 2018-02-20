module Football.Football exposing (Model, Msg, init, update, view, subs)

import String
import Html exposing (..)
import Html.Attributes exposing (..)
import Utils exposing (..)
import WebSocket
import Football.Data exposing (..)
import Football.CheckCompetition as CheckCompetition
import Ui.Header
import Table


-- MODEL


type alias Model =
    { protocol : String
    , host : String
    , games : List Game
    , tableState : Table.State
    , checkCompetition : CheckCompetition.Model
    }


type Msg
    = NewGamesChanges GamesChanges
    | ShowCheckCompetition
    | SetTableState Table.State
    | CheckCompetition CheckCompetition.Msg


init : String -> String -> ( Model, Cmd Msg )
init protocol host =
    ( { games = []
      , protocol = protocol
      , host = host
      , tableState = Table.initialSort "№"
      , checkCompetition = CheckCompetition.model
      }
    , Cmd.none
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        CheckCompetition msg_ ->
            let
                ( checkCompetition_, cmd_ ) =
                    CheckCompetition.update msg_ model.checkCompetition CheckCompetition
            in
                { model | checkCompetition = checkCompetition_ } ! [ cmd_ ]

        SetTableState newState ->
            { model | tableState = newState } ! []

        NewGamesChanges gamesChanges ->
            let
                newGames =
                    model.games
                        |> updateGames gamesChanges

                checkCompetition_ =
                    CheckCompetition.setGames model.checkCompetition newGames
            in
                { model | games = newGames, checkCompetition = checkCompetition_ } ! []

        ShowCheckCompetition ->
            { model | checkCompetition = CheckCompetition.setVisible model.checkCompetition True } ! []



-- SUBSCRIPTIONS


subs : (Msg -> msg) -> Model -> Sub msg
subs toMsg model =
    Sub.map toMsg (subscriptions model)


subscriptions : Model -> Sub Msg
subscriptions model =
    let
        listenGamesChanges =
            WebSocket.listen
                (websocketURL model.protocol model.host ++ "/football")
                (\str ->
                    case parseGames str of
                        Ok y ->
                            NewGamesChanges y

                        Err err ->
                            Debug.crash err
                )
    in
        [ listenGamesChanges
        , CheckCompetition.subs CheckCompetition model.checkCompetition
        ]
            |> Sub.batch



-- VIEW


view : Model -> Html Msg
view model =
    if List.isEmpty model.games then
        renderLoadingHeader
    else
        div []
            [ renderHeader
            , renderGamesTable model
            , CheckCompetition.view CheckCompetition model.checkCompetition
            ]


renderLoadingHeader : Html Msg      cxseaqDW2q1weds34r
renderLoadingHeader =
    Ui.Header.view
        [ renderTitle
        , Ui.Header.item
            { link = Nothing
            , target = "_blank"
            , action = Nothing
            , text = "Загрузка данных..."
            }
        ]


renderHeader : Html Msg
renderHeader =
    Ui.Header.view
        [ renderTitle
        , Ui.Header.spacer
        , Ui.Header.separator
        , Ui.Header.item
            { link = Nothing
            , target = "_blank"
            , action = Just ShowCheckCompetition
            , text = "Чемпионаты"
            }
        ]


renderTitle : Html a
renderTitle =
    Ui.Header.title
        { action = Nothing
        , target = "_self"
        , link = Nothing
        , text = "Футбол"
        }


renderGamesTable : Model -> Html Msg
renderGamesTable { games, checkCompetition, tableState } =
    let
        games_ =
            CheckCompetition.filterGames checkCompetition games
    in
        Table.view (configTable <| gamesHasInplay games_) tableState games_


configTable : Bool -> Table.Config Game Msg
configTable hasInplay =
    let
        dc =
            Table.defaultCustomizations
    in
        Table.customConfig
            { toId = .id >> toString
            , toMsg = SetTableState
            , columns =
                [ Table.intColumn "№" .order
                , Table.stringColumn "Дома" .home
                ]
                    ++ (if hasInplay then
                            [ columnScore ]
                        else
                            []
                       )
                    ++ [ Table.stringColumn "В гостях" .away
                       , columnTime
                       , Table.stringColumn "Чемпионат" .competition
                       , columnDollar "В паре" .totalMatched
                       , columnDollar "Не в паре" .totalAvailable
                       , columnPrice "П1+" .winBack
                       , columnPrice "П1-" .winLay
                       , columnPrice "Н+" .drawBack
                       , columnPrice "Н-" .drawLay
                       , columnPrice "П2+" .loseBack
                       , columnPrice "П2-" .loseLay
                       ]
            , customizations =
                { dc
                    | tableAttrs =
                        [ class "footbal-table"
                        ]
                }
            }


columnTime : Table.Column Game Msg
columnTime =
    Table.veryCustomColumn
        { name = "Время"
        , viewData =
            \x ->
                Table.HtmlDetails
                    [ style
                        [ ( "font-style", "italic" )
                        , ( "text-align", "center" )
                        ]
                    ]
                    [ text x.time ]
        , sorter = Table.increasingOrDecreasingBy .time
        }


columnScore : Table.Column Game Msg
columnScore =
    Table.stringColumn "Счёт"
        (\{ scoreHome, scoreAway, inplay } ->
            if inplay then
                toString scoreHome ++ " - " ++ toString scoreAway
            else
                ""
        )


columnDollar : String -> (a -> Float) -> Table.Column a msg
columnDollar name toValue =
    Table.veryCustomColumn
        { name = name
        , viewData =
            \x ->
                Table.HtmlDetails
                    [ style
                        [ ( "font-style", "italic" )
                        , ( "text-align", "right" )
                        ]
                    ]
                    (if toValue x > 0 then
                        [ text <| toString (toValue x) ++ " $" ]
                     else
                        []
                    )
        , sorter = Table.increasingOrDecreasingBy toValue
        }


columnPrice : String -> (a -> Float) -> Table.Column a msg
columnPrice name toValue =
    Table.veryCustomColumn
        { name = name
        , viewData =
            \x ->
                Table.HtmlDetails
                    [ style [ ( "text-align", "right" ) ] ]
                    (if toValue x > 0 then
                        [ text <| toString (toValue x) ]
                     else
                        []
                    )
        , sorter = Table.unsortable
        }


numToStr : number -> String
numToStr x =
    if x == 0 then
        ""
    else
        toString x



-- HELPERS


listZip : List a -> List b -> List ( a, b )
listZip xs ys =
    case ( xs, ys ) of
        ( x :: xBack, y :: yBack ) ->
            ( x, y ) :: listZip xBack yBack

        ( _, _ ) ->
            []


hasNotEmpty : (a -> String) -> List a -> Bool
hasNotEmpty f =
    List.filter (f >> String.isEmpty >> not)
        >> List.head
        >> Maybe.map (\_ -> True)
        >> Maybe.withDefault False

module Football.Football exposing (Model, Msg, init, update, view, subs)

import String
import Html exposing (..)
import Utils exposing (..)
import WebSocket
import Football.Data exposing (..)
import Ui.Modal
import Ui.Button
import Ui.Container
import Table


-- MODEL


type alias Model =
    { protocol : String
    , host : String
    , games : List Game
    , uiModal : Ui.Modal.Model
    , buttonSettings : Ui.Button.Model
    , tableState : Table.State
    }


type Msg
    = NewGamesChanges GamesChanges
    | UiModal Ui.Modal.Msg
    | SettingsDialog
    | SetTableState Table.State


init : String -> String -> ( Model, Cmd Msg )
init protocol host =
    ( { games = []
      , protocol = protocol
      , host = host
      , uiModal =
            { closable = True
            , backdrop = True
            , open = False
            }
      , buttonSettings = Ui.Button.model "Настройка" "primary" "medium"
      , tableState = Table.initialSort "№"
      }
    , Cmd.none
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SetTableState newState ->
            { model | tableState = newState } ! []

        UiModal msg_ ->
            { model | uiModal = Ui.Modal.update msg_ model.uiModal } ! []

        NewGamesChanges gamesChanges ->
            let
                newGames =
                    model.games
                        |> updateGames gamesChanges

                --|> sortGames model.order model.sortCol
            in
                { model | games = newGames } ! []

        SettingsDialog ->
            let
                uiModal =
                    model.uiModal
            in
                { model | uiModal = { uiModal | open = True } } ! []



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
view { games, uiModal, buttonSettings, tableState } =
    Ui.Container.view
        { direction = "row", align = "center", compact = False }
        []
        [ Table.view configTable tableState games
        , Ui.Modal.view
            (Ui.Modal.ViewModel
                [ text "Привет!" ]
                [ text "Пока!" ]
                UiModal
                "Настройки"
            )
            uiModal
        , Ui.Button.view SettingsDialog buttonSettings
        ]


configTable : Table.Config Game Msg
configTable =
    Table.config
        { toId = .id >> toString
        , toMsg = SetTableState
        , columns =
            [ Table.intColumn "№" .order
            , Table.stringColumn "Дома" .home
            , Table.stringColumn "Счёт"
                (\{ scoreHome, scoreAway, inplay } ->
                    if inplay then
                        toString scoreHome ++ " - " ++ toString scoreAway
                    else
                        ""
                )
            , Table.stringColumn "В гостях" .away
            , Table.stringColumn "Время" .time
            , Table.stringColumn "Страна" .country
            , Table.stringColumn "Чемпионат" .competition
            , dollarColumn "В паре" .totalMatched
            , dollarColumn "Не в паре" .totalAvailable
            , priceColumn "П1+" .winBack
            , priceColumn "П1-" .winLay
            , priceColumn "Н+" .drawBack
            , priceColumn "Н-" .drawLay
            , priceColumn "П2+" .loseBack
            , priceColumn "П2-" .loseLay
            ]
        }


dollarColumn : String -> (a -> Float) -> Table.Column a msg
dollarColumn name toValue =
    Table.customColumn
        { name = name
        , viewData =
            \x ->
                if toValue x > 0 then
                    toString (toValue x) ++ "$"
                else
                    ""
        , sorter = Table.increasingOrDecreasingBy toValue
        }


priceColumn : String -> (a -> Float) -> Table.Column a msg
priceColumn name toValue =
    Table.customColumn
        { name = name
        , viewData =
            \x ->
                if toValue x > 0 then
                    toString (toValue x)
                else
                    ""
        , sorter = Table.unsortable
        }


numToStr : number -> String
numToStr x =
    if x == 0 then
        ""
    else
        toString x



-- HELPERS


hasNotEmpty : (a -> String) -> List a -> Bool
hasNotEmpty f =
    List.filter (f >> String.isEmpty >> not)
        >> List.head
        >> Maybe.map (\_ -> True)
        >> Maybe.withDefault False

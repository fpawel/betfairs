module Football.Football exposing (Model, Msg, init, update, view, subs)

import String
import Html exposing (..)
import Html.Attributes exposing (style, attribute, class)
import Utils exposing (..)
import WebSocket
import Football.Data exposing (..)
import Ui.Modal
import Ui.IconButton
import Ui.Icons
import Ui.Container
import Ui.Checkbox
import Table
import Dict


-- MODEL


type alias Model =
    { protocol : String
    , host : String
    , games : List Game
    , uiModal : Ui.Modal.Model
    , tableState : Table.State
    , comps : Dict.Dict String CompModel
    }


type alias CompModel =
    { checkbox : Ui.Checkbox.Model
    , comp : String
    , value : Bool
    }


type Msg
    = NewGamesChanges GamesChanges
    | UiModal Ui.Modal.Msg
    | SettingsDialog
    | SetTableState Table.State
    | UiCheckbox String Ui.Checkbox.Msg
    | CheckboxChanged String Bool


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
      , tableState = Table.initialSort "№"
      , comps = Dict.empty
      }
    , Cmd.none
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UiCheckbox uid msg_ ->
            case Dict.get uid model.comps of
                Nothing ->
                    model ! []

                Just comp ->
                    let
                        ( newCheckbox, cmd_ ) =
                            Ui.Checkbox.update msg_ comp.checkbox

                        newComps =
                            Dict.insert uid { comp | checkbox = newCheckbox } model.comps
                    in
                        { model | comps = newComps } ! [ Cmd.map (UiCheckbox uid) cmd_ ]

        SetTableState newState ->
            { model | tableState = newState } ! []

        UiModal msg_ ->
            { model | uiModal = Ui.Modal.update msg_ model.uiModal } ! []

        NewGamesChanges gamesChanges ->
            let
                newGames =
                    model.games
                        |> updateGames gamesChanges

                newComps =
                    newGames
                        |> gamesCompetitions
                        |> List.map
                            (\comp ->
                                let
                                    v =
                                        Dict.toList model.comps
                                            |> List.map (Tuple.second)
                                            |> List.filter (.comp >> (==) comp)
                                            |> List.map (.checkbox >> .value)
                                            |> List.head
                                            |> Maybe.withDefault True

                                    cb =
                                        Ui.Checkbox.init ()
                                            |> Ui.Checkbox.setValue v
                                in
                                    ( cb.uid
                                    , { checkbox = cb
                                      , comp = comp
                                      , value = True
                                      }
                                    )
                            )
                        |> Dict.fromList

                --|> sortGames model.order model.sortCol
            in
                { model | games = newGames, comps = newComps } ! []

        SettingsDialog ->
            let
                uiModal =
                    model.uiModal
            in
                { model | uiModal = { uiModal | open = True } } ! []

        CheckboxChanged uid value ->
            let
                newComps =
                    Dict.map
                        (\k x ->
                            if k == uid then
                                { x | checkbox = Ui.Checkbox.setValue value x.checkbox }
                                --, checkbox = Ui.Checkbox.setValue value x.checkbox
                            else
                                x
                        )
                        model.comps
            in
                { model | comps = newComps } ! []



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

        cbs =
            model.comps
                |> Dict.toList
                |> List.map
                    (\( uid, x ) -> Ui.Checkbox.onChange (CheckboxChanged uid) x.checkbox)
    in
        [ listenGamesChanges ]
            ++ cbs
            |> Sub.batch



-- VIEW


view : Model -> Html Msg
view model =
    Ui.Container.rowEnd []
        [ div []
            [ rederGamesTable model
            , renderSettingsDialog model
            ]
        , div []
            [ Ui.IconButton.view SettingsDialog buttonSettingsModel
            , renderCompsCheckboxes model
            ]
        ]


renderCompsCheckboxes : Model -> Html Msg
renderCompsCheckboxes model =
    model.comps
        |> Dict.toList
        |> List.sortBy (Tuple.second >> .comp)
        |> List.map
            (\( uid, x ) ->
                li []
                    [ Ui.Checkbox.view x.checkbox |> Html.map (UiCheckbox uid)
                    , text x.comp
                    ]
            )
        |> ul []


renderSettingsDialog : Model -> Html Msg
renderSettingsDialog { uiModal } =
    Ui.Modal.view
        (Ui.Modal.ViewModel
            [ text "Привет!" ]
            [ text "Пока!" ]
            UiModal
            "Настройки"
        )
        uiModal


rederGamesTable : Model -> Html Msg
rederGamesTable { games, tableState } =
    Table.view (configTable <| gamesHasInplay games) tableState games


buttonSettingsModel : Ui.IconButton.Model Msg
buttonSettingsModel =
    { disabled = False
    , readonly = False
    , kind = "secondary"
    , size = "small"
    , glyph = Ui.Icons.plus []
    , side = "left"
    , text = "Настройки"
    }


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
                       , Table.stringColumn "Время" .time
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
            , customizations =
                { dc
                    | tableAttrs =
                        [ class "footbal-table"
                        ]
                }
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


dollarColumn : String -> (a -> Float) -> Table.Column a msg
dollarColumn name toValue =
    Table.veryCustomColumn
        { name = name
        , viewData =
            \x ->
                Table.HtmlDetails
                    [ class "dollar" ]
                    (if toValue x > 0 then
                        [ text <| toString (toValue x) ++ " $" ]
                     else
                        []
                    )
        , sorter = Table.increasingOrDecreasingBy toValue
        }


priceColumn : String -> (a -> Float) -> Table.Column a msg
priceColumn name toValue =
    Table.veryCustomColumn
        { name = name
        , viewData =
            \x ->
                Table.HtmlDetails
                    []
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

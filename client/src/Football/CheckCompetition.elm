module Football.CheckCompetition exposing (Model, Msg, model, setGames, update, view, subs, filter, setVisible)

import Html exposing (..)
import Html.Attributes exposing (style, attribute, class)
import Ui.Checkbox
import Dict
import Set
import Football.Data exposing (..)
import Table
import Ui.ButtonGroup
import Ui.Modal


type alias Model =
    { tableState : Table.State
    , items : List Item
    , uiModal : Ui.Modal.Model
    }


type alias Item =
    { checkbox : Ui.Checkbox.Model
    , name : String
    , size : Float
    , count : Int
    , country : String
    }


type Msg
    = UiCheckbox String Ui.Checkbox.Msg
    | CheckboxChanged String Bool
    | SetTableState Table.State
    | Check Bool
    | UiModal Ui.Modal.Msg


filter :
    Model
    -> List Game
    -> List Game
filter model =
    let
        st =
            model.items
                |> List.filter (.checkbox >> .value)
                |> List.map .name
                |> Set.fromList
    in
        List.filter (.competition >> \x -> Set.member x st)


update : Msg -> Model -> (Msg -> msg) -> ( Model, Cmd msg )
update msg model toMsg =
    case msg of
        UiModal msg_ ->
            { model | uiModal = Ui.Modal.update msg_ model.uiModal } ! []

        Check value ->
            { model
                | items = List.map (\x -> { x | checkbox = Ui.Checkbox.setValue value x.checkbox }) model.items
            }
                ! []

        SetTableState newState ->
            { model | tableState = newState } ! []

        UiCheckbox uid msg_ ->
            let
                m =
                    model.items
                        |> List.map (\x -> ( x.checkbox.uid, x ))
                        |> Dict.fromList
            in
                case Dict.get uid m of
                    Nothing ->
                        model ! []

                    Just x ->
                        let
                            ( checkbox_, cmd_ ) =
                                Ui.Checkbox.update msg_ x.checkbox

                            items =
                                Dict.insert uid { x | checkbox = checkbox_ } m
                                    |> Dict.toList
                                    |> List.map Tuple.second
                        in
                            { model | items = items } ! [ Cmd.map (UiCheckbox uid >> toMsg) cmd_ ]

        CheckboxChanged uid value ->
            let
                items =
                    List.map
                        (\x ->
                            if x.checkbox.uid == uid then
                                { x | checkbox = Ui.Checkbox.setValue value x.checkbox }
                            else
                                x
                        )
                        model.items
            in
                { model | items = items } ! []


model : Model
model =
    { items = []
    , tableState = Table.initialSort "Ставки"
    , uiModal =
        { closable = True
        , backdrop = True
        , open = False
        }
    }


setVisible : Model -> Bool -> Model
setVisible model open =
    let
        uiModal =
            model.uiModal
    in
        { model | uiModal = { uiModal | open = open } }


setGames : Model -> List Game -> Model
setGames model games =
    let
        m =
            model.items
                |> List.map (\x -> ( x.name, x ))
                |> Dict.fromList

        items =
            gamesCompetitions games
                |> List.map
                    (\( name, xs ) ->
                        let
                            size =
                                List.sum <| List.map (\x -> x.totalMatched + x.totalAvailable) xs

                            value =
                                Dict.get name m
                                    |> Maybe.map (.checkbox >> .value)
                                    |> Maybe.withDefault True

                            cb =
                                Ui.Checkbox.init ()

                            country =
                                case xs |> List.head |> Maybe.map .country of
                                    Just s ->
                                        s

                                    _ ->
                                        Debug.crash "empty list"
                        in
                            { checkbox = Ui.Checkbox.setValue value cb
                            , name = name
                            , size = size
                            , count = List.length xs
                            , country = country
                            }
                    )
    in
        { model | items = items, tableState = model.tableState }


renderCompetitions : Model -> Html Msg
renderCompetitions model =
    div
        [ style
            [ ( "height", "500px" )
            , ( "overflow", "auto" )
            ]
        ]
        [ Table.view configTable model.tableState model.items
        ]


renderTollButtons : Html Msg
renderTollButtons =
    Ui.ButtonGroup.model
        [ ( "Выбрать все", Check True )
        , ( "Снять выбор", Check False )
        ]
        |> Ui.ButtonGroup.view


view : (Msg -> msg) -> Model -> Html msg
view toMsg model =
    Ui.Modal.view
        (Ui.Modal.ViewModel
            [ renderCompetitions model ]
            [ renderTollButtons ]
            UiModal
            "Чемпионаты"
        )
        model.uiModal
        |> Html.map toMsg


configTable : Table.Config Item Msg
configTable =
    Table.customConfig
        { toId = .checkbox >> .uid
        , toMsg = SetTableState
        , columns =
            [ Table.veryCustomColumn
                { name = "Чемпионат"
                , sorter = Table.increasingOrDecreasingBy .name
                , viewData =
                    \x ->
                        Table.HtmlDetails []
                            [ table []
                                [ tr []
                                    [ td [] [ Ui.Checkbox.view x.checkbox |> Html.map (UiCheckbox x.checkbox.uid) ]
                                    , td
                                        [ attribute "valign" "middle"
                                        , style [ ( "text-align", "left" ) ]
                                        ]
                                        [ text <|
                                            if x.name == "" then
                                                "(без названия)"
                                            else
                                                x.name
                                        ]
                                    ]
                                ]
                            ]
                }
            , Table.stringColumn "Страна" .country
            , Table.veryCustomColumn
                { name = "Ставки"
                , sorter = Table.increasingOrDecreasingBy .size
                , viewData =
                    \x ->
                        Table.HtmlDetails
                            [ class "dollar" ]
                            (if x.size > 0 then
                                [ text <| toString x.size ++ " $" ]
                             else
                                []
                            )
                }
            , Table.veryCustomColumn
                { name = "Кол-во матчей"
                , sorter = Table.increasingOrDecreasingBy .count
                , viewData =
                    \x ->
                        Table.HtmlDetails
                            [ attribute "align" "right" ]
                            [ text <| toString x.count ]
                }
            ]
        , customizations =
            Table.defaultCustomizations
        }


subs : (Msg -> msg) -> Model -> Sub msg
subs toMsg =
    .items
        >> List.map
            (\x -> Ui.Checkbox.onChange (CheckboxChanged x.checkbox.uid) x.checkbox)
        >> Sub.batch
        >> Sub.map toMsg

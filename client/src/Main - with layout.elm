module Main exposing (..)

import Html exposing (..)
import Navigation exposing (Location)
import Football
import Material
import Material.Scheme as Scheme
import Material.Layout as Layout


--import Html.Attributes exposing (..)
--import Html.Events exposing (onClick)


main : Program Never Model Msg
main =
    Navigation.program LocationChanged
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }



-- MODEL


type alias Model =
    { location : Location
    , football : Football.Model
    , mdl : Material.Model
    }


init : Location -> ( Model, Cmd Msg )
init location =
    let
        ( football, footballCmd ) =
            Football.init location.protocol location.host
    in
        { location = location
        , football = football
        , mdl = Material.model
        }
            ! [ Cmd.map FootballMsg footballCmd, Layout.sub0 Mdl ]



-- UPDATE


type Msg
    = LocationChanged Location
    | FootballMsg Football.Msg
    | Mdl (Material.Msg Msg)
    | Nop


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Nop ->
            model ! []

        Mdl msg_ ->
            Material.update Mdl msg_ model

        LocationChanged newLocation ->
            ( { model | location = newLocation }, Cmd.none )

        FootballMsg msgFootball ->
            let
                ( newFootball, footballCmd ) =
                    Football.update msgFootball model.football
            in
                ( { model | football = newFootball }, Cmd.map FootballMsg footballCmd )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    [ Football.subs FootballMsg model.football
    , Layout.subs Mdl model.mdl
    ]
        |> Sub.batch



-- VIEW


view : Model -> Html Msg
view model =
    let
        football =
            Html.map FootballMsg <| Football.view model.football

        xx =
            div []
                [ football
                ]
                |> Scheme.top
    in
        Layout.render Mdl
            model.mdl
            []
            { header =
                [ Layout.row []
                    [ Layout.title []
                        [ h3 [] [ text "Футбол" ]
                        ]
                    ]
                ]
            , drawer = drawer
            , tabs = ( [ div [] [] ], [] )
            , main = [ football ]
            }
            |> Scheme.top


drawer : List (Html Msg)
drawer =
    [ Layout.title [] [ text "Example drawer" ]
    , Layout.navigation
        []
        [ Layout.link
            [ Layout.href "https://github.com/debois/elm-mdl" ]
            [ text "github" ]
        , Layout.link
            [ Layout.href "http://package.elm-lang.org/packages/debois/elm-mdl/latest/" ]
            [ text "elm-package" ]
        , Layout.link
            [ Layout.href "#cards"
              --, Options.onClick (Layout.toggleDrawer Mdl)
            ]
            [ text "Card component" ]
        ]
    ]



-- HELPERS

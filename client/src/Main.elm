module Main exposing (..)

import Html exposing (..)
import Navigation exposing (Location)
import Football
import Material.Scheme


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


type Msg
    = LocationChanged Location
    | FootballMsg Football.Msg


type alias Model =
    { location : Location
    , football : Football.Model
    }


init : Location -> ( Model, Cmd Msg )
init location =
    let
        ( football, footballCmd ) =
            Football.init location.protocol location.host
    in
        Model location football ! [ Cmd.map FootballMsg footballCmd ]



-- UPDATE


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
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
    Sub.map FootballMsg <| Football.subscriptions model.football



-- VIEW


view : Model -> Html Msg
view model =
    let
        football =
            Html.map FootballMsg <| Football.view model.football
    in
        div []
            [ football
            ]
            |> Material.Scheme.top



-- HELPERS

module Main exposing (..)

import Html exposing (..)
import Navigation exposing (Location)
import Football
import Utils exposing (..)
import WebSocket
import Debug


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
    , football : List Football.Game
    }


init : Location -> ( Model, Cmd Msg )
init location =
    ( Model location [], Cmd.none )



-- UPDATE


type Msg
    = LocationChanged Location
    | NewFootball (List Football.Game)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        LocationChanged newLocation ->
            ( { model | location = newLocation }, Cmd.none )

        NewFootball x ->
            ( { model | football = x }, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    WebSocket.listen
        (websocketURL model.location.protocol model.location.host ++ "/football")
        (\str ->
            case Football.parseGames str of
                Ok y ->
                    NewFootball y

                Err err ->
                    Debug.crash err
        )



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ Football.renderGames model.football ]

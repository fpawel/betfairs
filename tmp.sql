-- команды
CREATE TABLE IF NOT EXISTS teams (
  team_id SERIAL PRIMARY KEY,
  team_name TEXT UNIQUE
);

-- пример добавления команды 'Спартак'
INSERT INTO teams(team_name)
VALUES ('Спартак')
ON CONFLICT (team_name) DO NOTHING
RETURNING team_id;

-- чемпионаты
CREATE TABLE IF NOT EXISTS competitions (
  competition_id INT UNIQUE NOT NULL CHECK (competition_id > 0),
  competition_name TEXT,
  PRIMARY KEY (competition_id)
);

-- пример добавления чемпионата
INSERT INTO competitions(competition_id, competition_name)
VALUES (8596554, 'Israeli Liga Bet - North B')
ON CONFLICT (competition_id) DO NOTHING;


CREATE TABLE IF NOT EXISTS events (
  event_id INT UNIQUE NOT NULL CHECK (event_id > 0),
  open_date DATE,
  competition_id INT,
  home_id INT,
  away_id INT,
  country_code2 VARCHAR(2),
  FOREIGN KEY (competition_id)
  REFERENCES competitions (competition_id),
  FOREIGN KEY (home_id, away_id)
  REFERENCES teams (team_id),
  PRIMARY KEY (event_id, open_date)
);




CREATE TABLE IF NOT EXISTS markets (
  market_id INTEGER UNIQUE NOT NULL CHECK (market_id > 0),
  market_name TEXT,
  event_id INT NOT NULL CHECK (event_id > 0),
  FOREIGN KEY (event_id)
  REFERENCES events (event_id) ON DELETE CASCADE,
  PRIMARY KEY (market_id)
);
CREATE TABLE IF NOT EXISTS runners (
  runner_id INTEGER ,
  runner_name TEXT,
  market_id INT NOT NULL CHECK (market_id > 0),
  FOREIGN KEY (market_id)
  REFERENCES markets (market_id) ON DELETE CASCADE,
  PRIMARY KEY (market_id, runner_id)
);

DROP TABLE runners;
DROP TABLE markets;
DROP TABLE events;

TRUNCATE runners;
TRUNCATE runners, markets, events CASCADE ;

INSERT INTO events (event_id, competition, event_date)
VALUES
  (2, 'Чемпионат Англии', CURRENT_DATE),
  (3, 'Чемпионат России', CURRENT_DATE),
  (4, 'Чемпионат России', '2017-12-04'),
  (5, 'Чемпионат России', '2017-12-04'),
  (6, 'Чемпионат России', '2017-12-04')
ON CONFLICT (event_id) DO NOTHING ;
INSERT INTO markets (market_id, market_name, event_id)
VALUES
  (4, 'Ставки', 2),
  (5, 'Ставки', 3),
  (1, 'Ставки 1', 5),
  (2, 'Ставки 2', 5),
  (3, 'Ставки 3', 5)
ON CONFLICT (market_id) DO NOTHING ;

INSERT INTO runners (runner_id, runner_name, market_id)
VALUES
  (4, 'П', 4),
  (4, 'П', 5),
  (4, 'П', 2),
  (4, 'П', 3)
ON CONFLICT (market_id, runner_id) DO NOTHING ;

WITH yesterday_events AS (
    SELECT event_id FROM  events
    WHERE event_date < CURRENT_DATE
), yesterday_markets AS(
    SELECT market_id FROM markets
    WHERE event_id IN (yesterday_events)
)
DELETE FROM runners
WHERE market_id IN (yesterday_markets),
DELETE FROM markets WHERE event_id IN (yesterday_events)
;

DELETE FROM runners
WHERE market_id in
      (
        SELECT market_id from markets
        WHERE event_id in
              (
                SELECT event_id FROM  events
                WHERE event_date < CURRENT_DATE
              )
      );

DELETE FROM markets
WHERE event_id in
      (
        SELECT event_id FROM  events
        WHERE event_date < CURRENT_DATE
      );
DELETE FROM events WHERE event_date < CURRENT_DATE - '7 day'::interval;

UPDATE events SET event_date = '2017-11-12', competition = 'FAKE'
WHERE event_id = 28504185;
INSERT INTO events (event_id, competition, event_date)
VALUES (1, 'FAKE', '2017-11-12');


SELECT * FROM events;
SELECT * FROM markets;
SELECT * FROM runners;

INSERT INTO events (event_id, event_date, competition)
VALUES (1,CURRENT_DATE, 'Чемпионат Англии')
ON CONFLICT (event_id)
  DO UPDATE SET competition = 'Чемпионат испании';

DELETE  FROM events WHERE event_id = 1;
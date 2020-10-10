#!/bin/bash
 echo "ABCDEFG====================================="
 echo "CREATE DATABASE ac2;
use ac2;

CREATE TABLE users (
  userid VARCHAR(255) not null,
  pwhash BINARY(60) not null,
  attribute INT default 0,
  primary key (userid)
);

-- Change the initial user name and password into an argument.
INSERT INTO users (userid, pwhash, attribute) VALUES ('${AC2_APP_ADMINUSERNAME}', '${AC2_APP_ADMINPASSWORD_HASHED}', 1);

CREATE TABLE events (
  id BINARY(16) not null primary key,
  startdate DATETIME not null,
  track VARCHAR(64),
  weatherRandomness INT,
  P_hourOfDay INT,
  P_timeMultiplier INT,
  P_sessionDurationMinute INT,
  Q_hourOfDay INT,
  Q_timeMultiplier INT,
  Q_sessionDurationMinute INT,
  R_hourOfDay INT,
  R_timeMultiplier INT,
  R_sessionDurationMinute INT,
  pitWindowLengthSec INT,
  isRefuellingAllowedInRace BOOLEAN,
  mandatoryPitstopCount INT,
  isMandatoryPitstopRefuellingRequired BOOLEAN,
  isMandatoryPitstopTyreChangeRequired BOOLEAN,
  isMandatoryPitstopSwapDriverRequired BOOLEAN,
  tyreSetCount INT
);

INSERT into events (id, startdate, track, weatherRandomness, P_hourOfDay, P_timeMultiplier, P_sessionDurationMinute, Q_hourOfDay, Q_timeMultiplier, Q_sessionDurationMinute, R_hourOfDay, R_timeMultiplier, R_sessionDurationMinute, pitWindowLengthSec, isRefuellingAllowedInRace, mandatoryPitstopCount, isMandatoryPitstopRefuellingRequired, isMandatoryPitstopTyreChangeRequired, isMandatoryPitstopSwapDriverRequired, tyreSetCount) VALUES (0x0fba3964131745a9b8f4a9f762ba05f8, '2000-01-01 11:30:00', 'zandvoort_2019', 1, 12, 1, 10, 13, 1, 10, 14, 1, 10, 120, true, 1, false, true, true, 3);" | mysql -u root --password=${MYSQL_ROOT_PASSWORD}
 echo "ABCDEFG====================================="

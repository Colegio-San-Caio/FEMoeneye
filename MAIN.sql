-- MAIN.sql - KDP + eal7 FLOP QUADRUPLE 8038 - SQLite with triggers
-- run: sqlite3 eal7.db < MAIN.sql && sqlite3 eal7.db "SELECT * FROM matrix; SELECT * FROM flop_quadruple;"

DROP TABLE IF EXISTS kdp_profile;
DROP TABLE IF EXISTS eal7_events;

-- KDP Forum Profile (from main.y / lex/main.l / MAIN.r)
CREATE TABLE kdp_profile (
  nickname TEXT PRIMARY KEY,
  code TEXT NOT NULL,
  role TEXT DEFAULT 'Member',
  points INTEGER DEFAULT 1,
  boa TEXT DEFAULT 'application/vnd.x-mosfetq-dos-image',
  iana INTEGER DEFAULT 1460172,
  port INTEGER DEFAULT 8038
);

INSERT INTO kdp_profile (nickname, code, role, points) VALUES
('KDPMember:KL5UQ4f48673352', 'KL5UQ4f48673352', 'Member', 1);

-- eal7 events - like FEMoeneye main.ximg
CREATE TABLE eal7_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  ts DATETIME DEFAULT CURRENT_TIMESTAMP,
  page TEXT NOT NULL,
  event TEXT NOT NULL CHECK(event IN ('page_view','fem_run','ximg_view')),
  flops INTEGER NOT NULL DEFAULT 0,
  nickname TEXT REFERENCES kdp_profile(nickname),
  source TEXT DEFAULT 'sql'
);

-- TRIGGER: auto-calc FLOP QUADRUPLE = 4 + 16 + 64 = 84 on INSERT (like main.go did)
CREATE TRIGGER trg_flop_quadruple_insert
BEFORE INSERT ON eal7_events
FOR EACH ROW
BEGIN
  SELECT CASE NEW.event
    WHEN 'page_view' THEN 4
    WHEN 'fem_run' THEN 16
    WHEN 'ximg_view' THEN 64
    ELSE 4
  END INTO NEW.flops;
END;

-- Seed with your current persistent data (80 + 4 = 84 from eal7.sqlite.json)
INSERT INTO eal7_events (page, event, nickname) VALUES ('main.ximg', 'fem_run', 'KDPMember:KL5UQ4f48673352');
INSERT INTO eal7_events (page, event, nickname) VALUES ('eal7.x', 'ximg_view', 'KDPMember:KL5UQ4f48673352');
INSERT INTO eal7_events (page, event, nickname) VALUES ('eal7.html', 'page_view', 'KDPMember:KL5UQ4f48673352');

-- Views = same as /matrix endpoint in main.ximg Go machine code
CREATE VIEW matrix AS
SELECT page, event, COUNT(*) as count, SUM(flops) as total_flops
FROM eal7_events GROUP BY page, event;

CREATE VIEW flop_quadruple AS
SELECT COUNT(*) as total_events, SUM(flops) as total_flops_quadruple,
       SUM(CASE WHEN event='page_view' THEN 1 ELSE 0 END) as page_views,
       SUM(CASE WHEN event='fem_run' THEN 1 ELSE 0 END) as fem_runs,
       SUM(CASE WHEN event='ximg_view' THEN 1 ELSE 0 END) as ximg_views,
       '8038 flop quadruple' as port,
       'application/vnd.x-mosfetq-dos-image' as boa,
       1460172 as iana
FROM eal7_events;

-- R-compatible data.frame export
CREATE VIEW df_r AS
SELECT page, event, flops, SUM(flops) OVER (ORDER BY id) as total_flops FROM eal7_events;

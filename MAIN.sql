-- MAIN.sql - KDP + eal7 FLOP QUADRUPLE 8038 - SQLite FIXED
DROP TABLE IF EXISTS eal7_events;
DROP TABLE IF EXISTS kdp_profile;

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

CREATE TABLE eal7_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  ts DATETIME DEFAULT CURRENT_TIMESTAMP,
  page TEXT NOT NULL,
  event TEXT NOT NULL CHECK(event IN ('page_view','fem_run','ximg_view')),
  flops INTEGER NOT NULL DEFAULT 0,
  nickname TEXT REFERENCES kdp_profile(nickname),
  source TEXT DEFAULT 'sql'
);

-- FIXED TRIGGER: SQLite doesn't allow SELECT INTO NEW, use UPDATE AFTER INSERT
CREATE TRIGGER trg_flop_quadruple_insert
AFTER INSERT ON eal7_events
FOR EACH ROW
BEGIN
  UPDATE eal7_events SET flops = CASE NEW.event
    WHEN 'page_view' THEN 4
    WHEN 'fem_run' THEN 16
    WHEN 'ximg_view' THEN 64
    ELSE 4
  END WHERE id = NEW.id;
END;

INSERT INTO eal7_events (page, event, nickname) VALUES ('main.ximg', 'fem_run', 'KDPMember:KL5UQ4f48673352');
INSERT INTO eal7_events (page, event, nickname) VALUES ('eal7.x', 'ximg_view', 'KDPMember:KL5UQ4f48673352');
INSERT INTO eal7_events (page, event, nickname) VALUES ('eal7.html', 'page_view', 'KDPMember:KL5UQ4f48673352');

CREATE VIEW matrix AS SELECT page, event, COUNT(*) as count, SUM(flops) as total_flops FROM eal7_events GROUP BY page, event;
CREATE VIEW flop_quadruple AS SELECT COUNT(*) as total_events, SUM(flops) as total_flops_quadruple, 8038 as port, 'application/vnd.x-mosfetq-dos-image' as boa, 1460172 as iana, '4+16+64=84' as formula FROM eal7_events;

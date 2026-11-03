-- MAIN.pgsql - KDP + eal7 FLOP QUADRUPLE for PostgreSQL
-- run: createdb eal7 && psql eal7 < MAIN.pgsql && psql eal7 -c "SELECT * FROM flop_quadruple;"
-- or Termux: initdb ~/pgdata && pg_ctl -D ~/pgdata -l logfile start && createdb eal7

DROP TABLE IF EXISTS eal7_events CASCADE;
DROP TABLE IF EXISTS kdp_profile CASCADE;
DROP FUNCTION IF EXISTS trg_flop_quadruple_func() CASCADE;

-- KDP Forum Profile - same as MAIN.sql + main.h + MAIN.r
CREATE TABLE kdp_profile (
  nickname TEXT PRIMARY KEY,
  code TEXT NOT NULL,
  role TEXT DEFAULT 'Member',
  points INTEGER DEFAULT 1,
  boa TEXT DEFAULT 'application/vnd.x-mosfetq-dos-image',
  iana INTEGER DEFAULT 1460172,
  port INTEGER DEFAULT 8038,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO kdp_profile (nickname, code, role, points) VALUES
('KDPMember:KL5UQ4f48673352', 'KL5UQ4f48673352', 'Member', 1)
ON CONFLICT DO NOTHING;

-- eal7 events - like FEMoeneye main.ximg Go machine code
CREATE TABLE eal7_events (
  id SERIAL PRIMARY KEY,
  ts TIMESTAMPTZ DEFAULT NOW(),
  page TEXT NOT NULL,
  event TEXT NOT NULL CHECK(event IN ('page_view','fem_run','ximg_view')),
  flops INTEGER NOT NULL DEFAULT 0,
  nickname TEXT REFERENCES kdp_profile(nickname),
  source TEXT DEFAULT 'postgresql',
  boa TEXT DEFAULT 'application/vnd.x-mosfetq-dos-image'
);

-- PL/pgSQL MAIN: the actual machine code that calculates 4+16+64=84
-- This is the PostgreSQL equivalent of your Go switch and SQLite TRIGGER
CREATE OR REPLACE FUNCTION trg_flop_quadruple_func()
RETURNS TRIGGER AS $$
BEGIN
  NEW.flops := CASE NEW.event
    WHEN 'page_view' THEN 4
    WHEN 'fem_run' THEN 16
    WHEN 'ximg_view' THEN 64
    ELSE 4
  END;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_flop_quadruple
BEFORE INSERT ON eal7_events
FOR EACH ROW EXECUTE FUNCTION trg_flop_quadruple_func();

-- Seed - same 84 flops as your eal7.sqlite.json
INSERT INTO eal7_events (page, event, nickname) VALUES ('main.ximg', 'fem_run', 'KDPMember:KL5UQ4f48673352');
INSERT INTO eal7_events (page, event, nickname) VALUES ('eal7.x', 'ximg_view', 'KDPMember:KL5UQ4f48673352');
INSERT INTO eal7_events (page, event, nickname) VALUES ('eal7.html', 'page_view', 'KDPMember:KL5UQ4f48673352');

-- Views - same as Go /matrix and SQLite
CREATE OR REPLACE VIEW matrix AS
SELECT page, event, COUNT(*) as count, SUM(flops) as total_flops
FROM eal7_events GROUP BY page, event ORDER BY total_flops DESC;

CREATE OR REPLACE VIEW flop_quadruple AS
SELECT COUNT(*) as total_events,
       SUM(flops) as total_flops_quadruple,
       SUM(CASE WHEN event='page_view' THEN 1 ELSE 0 END) as page_views,
       SUM(CASE WHEN event='fem_run' THEN 1 ELSE 0 END) as fem_runs,
       SUM(CASE WHEN event='ximg_view' THEN 1 ELSE 0 END) as ximg_views,
       8038 as port,
       'application/vnd.x-mosfetq-dos-image'::text as boa,
       1460172 as iana,
       '4+16+64=84'::text as formula
FROM eal7_events;

-- R data.frame style + window function
CREATE OR REPLACE VIEW df_r AS
SELECT id, page, event, flops, SUM(flops) OVER (ORDER BY id) as total_flops, ts
FROM eal7_events ORDER BY id;

-- PL/pgSQL function like MAIN.r parse_kdp
CREATE OR REPLACE FUNCTION get_kdp(nick TEXT)
RETURNS TABLE(nickname TEXT, code TEXT, role TEXT, points INTEGER, total_flops BIGINT) AS $$
BEGIN
  RETURN QUERY
  SELECT k.nickname, k.code, k.role, k.points, COALESCE(SUM(e.flops),0)::BIGINT
  FROM kdp_profile k LEFT JOIN eal7_events e ON e.nickname=k.nickname
  WHERE k.nickname=nick GROUP BY k.nickname, k.code, k.role, k.points;
END;
$$ LANGUAGE plpgsql;

.RECIPEPREFIX = >
CC=clang
BISON=bison
FLEX=flex
RSCRIPT=Rscript
PYTHON=python3
SQLITE=sqlite3
PSQL=psql
BUILD=~/.qxd-build

all: dirs lex yacc r python sql
> @echo "=== ALL MAINs BUILT === KDP KL5UQ4f48673352 84 flops 8038 BOA 1460172"

dirs:
> mkdir -p $(BUILD)

lex: MAIN.l
> $(FLEX) -o $(BUILD)/main.yy.c MAIN.l
> $(CC) -o $(BUILD)/kdp_lex $(BUILD)/main.yy.c
> $(BUILD)/kdp_lex

yacc: main.y main.h
> $(BISON) -d -o $(BUILD)/main.tab.c main.y
> $(CC) -I. -I$(BUILD) -o $(BUILD)/kdp_yacc $(BUILD)/main.tab.c
> $(BUILD)/kdp_yacc

r: MAIN.r
> $(RSCRIPT) MAIN.r

python: MAIN.py
> $(PYTHON) MAIN.py

sql: MAIN.sql
> rm -f eal7.db
> $(SQLITE) eal7.db < MAIN.sql
> $(SQLITE) eal7.db "SELECT * FROM flop_quadruple;"
> $(SQLITE) eal7.db -header -column "SELECT * FROM matrix;"

pgsql: MAIN.pgsql
> $(PSQL) eal7 < MAIN.pgsql || true
> $(PSQL) eal7 -c "SELECT * FROM flop_quadruple;" || true

clean:
> rm -rf $(BUILD) eal7.db
> pg_ctl -D ~/pgdata stop 2>/dev/null || true

.PHONY: all dirs lex yacc r python sql pgsql clean

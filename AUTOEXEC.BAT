@echo off
REM FEAoeneye v5 - BOA Binary Observation Agency
REM IANA #1460172 - application/vnd.x-mosfetq-dos-image
REM Patreon: patreon.com/c/KAIOLAFKETELHUT
REM GitHub: Colegio-San-Caio/FEAoeneye - commit bdd4d2e
set FEA_ROOT=%~dp0
set PATH=%FEA_ROOT%;%PATH%
assoc .ximg=XIMGFile >nul 2>&1
ftype XIMGFile="%FEA_ROOT%curl.ximg" "%%1" >nul 2>&1
echo [AUTOEXEC] FEAoeneye T-QTS BOX init...
set PHIA_URL=https://jobs.ashbyhq.com/phia/d54a3c16-cf77-43a9-b3fd-1da3d7e9a6a0
set STABLE_SHA=22f483144eacf5f68e17fa0164fc1713b38a7b8172680f0f7435c47725c51540
echo [AUTOEXEC] Patreon: https://www.patreon.com/c/KAIOLAFKETELHUT
echo [AUTOEXEC] t=S/T complete - %STABLE_SHA%

REM --- IANA #1460172 - GNU MIME + Windows ---
echo [FEA] Installing ximg MIME for Windows
assoc .ximg=XIMGFile
ftype XIMGFile="%~dp0XIMGVIEW.EXE" "%%1"
echo application/vnd.x-mosfetq-dos-image = .ximg installed

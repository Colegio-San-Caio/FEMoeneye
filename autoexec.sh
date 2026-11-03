#!/data/data/com.termux/files/usr/bin/bash
echo "[AUTOEXEC.SH] FEAoeneye v5 - IANA #1460172"
echo "[AUTOEXEC.SH] curl.ximg = virus.exe = 7.3K"

# --- IANA #1460172 - application/vnd.x-mosfetq-dos-image - GNU MIME ---
echo "[FEA] Installing ximg MIME for GNU/Termux"
# Termux mime.types
mkdir -p $PREFIX/etc 2>/dev/null
echo "application/vnd.x-mosfetq-dos-image ximg XIMG" >> $PREFIX/etc/mime.types 2>/dev/null
echo "application/vnd.x-mosfetq-dos-image ximg XIMG" >> ~/.mime.types 2>/dev/null
sudo sh -c 'echo "application/vnd.x-mosfetq-dos-image ximg XIMG" >> /etc/mime.types' 2>/dev/null

# freedesktop shared-mime
mkdir -p ~/.local/share/mime/packages
cp x-mosfetq-dos-image.xml ~/.local/share/mime/packages/ 2>/dev/null
update-mime-database ~/.local/share/mime 2>/dev/null
echo "XIMG MIME installed: application/vnd.x-mosfetq-dos-image"

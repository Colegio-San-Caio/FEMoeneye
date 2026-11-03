# FEMoeneye - Direct FEM as CAM support in ERP (Apache)

Repo: https://github.com/Colegio-San-Caio/FEMoeneye
Twin: https://github.com/Colegio-San-Caio/FEAoeneye (BOA XIMG)
Parent: https://github.com/Colegio-San-Caio/FEAoeneye

## Purpose
Direct Finite Element Method (FEM) as CAM layer for ERP.

ERP (Apache OFBiz / Moqui / Odoo-style) -> FEMoeneye -> CAM Toolpath -> XIMG DOS image (application/vnd.x-mosfetq-dos-image)

IANA #1460172 compatible: output .ximg for BOA viewer.

## Architecture
ERP Order -> FEMoeneye Solver (direct stiffness) -> Stress mesh -> CAM support (G-code / XIMG)
- input: STEP/JSON work order from ERP
- FEM: 2D plane stress direct solver (numpy)
- CAM: contour extraction -> G-code + .ximg
- ERP: Apache HTTP callback to ERP API

## Apache Main
Built for Apache HTTPD / Tomcat ERP integration.
MIME: AddType application/vnd.x-mosfetq-dos-image .ximg

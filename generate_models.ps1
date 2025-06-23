Write-Host "Generando modelos con SQLBoiler..." -ForegroundColor Green
sqlboiler psql
if ($?) {
    Write-Host "Modelos generados correctamente." -ForegroundColor Green
} else {
    Write-Host "Error al generar los modelos." -ForegroundColor Red
}
Write-Host "Presiona cualquier tecla para continuar..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

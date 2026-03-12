param(
  [string]$OutputPath = ".\deployment\fluxstream.ico"
)

$ErrorActionPreference = 'Stop'

Add-Type -AssemblyName System.Drawing
Add-Type -AssemblyName System.Windows.Forms

$size = 256
$bitmap = New-Object System.Drawing.Bitmap $size, $size
$graphics = [System.Drawing.Graphics]::FromImage($bitmap)
$graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
$graphics.Clear([System.Drawing.Color]::Transparent)

$rect = New-Object System.Drawing.Rectangle 16, 16, 224, 224
$path = New-Object System.Drawing.Drawing2D.GraphicsPath
$radius = 52
$diameter = $radius * 2
$path.AddArc($rect.X, $rect.Y, $diameter, $diameter, 180, 90)
$path.AddArc($rect.Right - $diameter, $rect.Y, $diameter, $diameter, 270, 90)
$path.AddArc($rect.Right - $diameter, $rect.Bottom - $diameter, $diameter, $diameter, 0, 90)
$path.AddArc($rect.X, $rect.Bottom - $diameter, $diameter, $diameter, 90, 90)
$path.CloseFigure()

$gradient = New-Object System.Drawing.Drawing2D.LinearGradientBrush $rect, ([System.Drawing.Color]::FromArgb(255,37,99,235)), ([System.Drawing.Color]::FromArgb(255,14,165,233)), 45
$graphics.FillPath($gradient, $path)

$shadowRect = New-Object System.Drawing.Rectangle 32, 32, 192, 192
$shadowPath = New-Object System.Drawing.Drawing2D.GraphicsPath
$shadowPath.AddEllipse($shadowRect)
$shadowBrush = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::FromArgb(40,255,255,255))
$graphics.FillPath($shadowBrush, $shadowPath)

$bolt = New-Object System.Drawing.Drawing2D.GraphicsPath
$points = [System.Drawing.Point[]]@(
  (New-Object System.Drawing.Point 140, 44),
  (New-Object System.Drawing.Point 88, 138),
  (New-Object System.Drawing.Point 128, 138),
  (New-Object System.Drawing.Point 110, 214),
  (New-Object System.Drawing.Point 176, 112),
  (New-Object System.Drawing.Point 138, 112)
)
$bolt.AddPolygon($points)
$boltBrush = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::FromArgb(255,255,255,255))
$graphics.FillPath($boltBrush, $bolt)

$stroke = New-Object System.Drawing.Pen ([System.Drawing.Color]::FromArgb(55,255,255,255)), 4
$graphics.DrawPath($stroke, $path)

$ms = New-Object System.IO.MemoryStream
$bitmap.Save($ms, [System.Drawing.Imaging.ImageFormat]::Png)
$pngBytes = $ms.ToArray()
$ms.Dispose()
$graphics.Dispose()
$bitmap.Dispose()

$dir = Split-Path -Parent $OutputPath
if ($dir -and -not (Test-Path $dir)) {
  New-Item -ItemType Directory -Path $dir | Out-Null
}

$fs = [System.IO.File]::Open($OutputPath, [System.IO.FileMode]::Create)
$bw = New-Object System.IO.BinaryWriter($fs)

$bw.Write([UInt16]0)
$bw.Write([UInt16]1)
$bw.Write([UInt16]1)
$bw.Write([byte]0)
$bw.Write([byte]0)
$bw.Write([byte]0)
$bw.Write([byte]0)
$bw.Write([UInt16]1)
$bw.Write([UInt16]32)
$bw.Write([UInt32]$pngBytes.Length)
$bw.Write([UInt32]22)
$bw.Write($pngBytes)
$bw.Flush()
$bw.Dispose()
$fs.Dispose()

Write-Host "Icon generated:" (Resolve-Path $OutputPath).Path

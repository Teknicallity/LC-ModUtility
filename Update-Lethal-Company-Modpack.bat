# 2>NUL & @CLS & PUSHD "%~dp0" & "%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -nol -nop -ep bypass "[IO.File]::ReadAllText('%~f0')|iex" & POPD & EXIT /B

$directoryPath = ".\"

function Check-And-Delete-Bepinex-Related {
	$bepinexPath = Join-Path -Path $directoryPath -ChildPath "Bepinex"
	
	if (Test-Path -Path $bepinexPath -PathType Container){
		Remove-Item -Path $bepinexPath -Recurse -Force
		Write-Host "Bepinex folder deleted"
	} else{
		Write-Host "Bepinex directory does not exist. Creating directory"
	}
	
	$winhttpPath = Join-Path -Path $directoryPath -ChildPath "winhttp.dll"
	$doorstopPath = Join-Path -Path $directoryPath -ChildPath "doorstop_config.ini"
	
	if (Test-Path -Path $winhttpPath -PathType Leaf){
		Remove-Item -Path $winhttpPath -Force
	}
	if (Test-Path -Path $doorstopPath -PathType Leaf){
		Remove-Item -Path $doorstopPath -Force
	}
}

function Get-Download-Folder-Path {
	$userProfile = $env:USERPROFILE
	$downloadsFolder = Join-Path -Path $userProfile -ChildPath "Downloads"
	Write-Host "Downloads folder: $downloadsFolder"
	return $downloadsFolder
}

function Unzip-Bepinex($downloadsPath, $packName) {
	Write-Host "Unzipping: $packName"
	$zipInput = Join-Path -Path $downloadsPath -ChildPath $packName
	#$zipOutput = Join-Path -Path $directoryPath -ChildPath $directoryPath
	Expand-Archive -Force -Path $zipInput -DestinationPath $directoryPath #$zipOutput 
}

function Get-BepinexPack {
	param (
		[string]$downloadsFolder
	)
	$pattern = 'BepInExPack_v(\d+)'
	$matchingFiles = Get-ChildItem -Path $downloadsFolder | Where-Object { $_.Name -match $pattern }
	Find-Highest-Version $matchingFiles
}

function Find-Highest-Version($matchingFiles) {
	$maxVersion = 0
	$maxVersionFileName = ""

	# Iterate through the file names
	foreach ($fileName in $matchingFiles) {
		if ($fileName -match $pattern) {
			$version = [int]$matches[1]

			# Check if the current version is greater than the maximum version
			if ($version -gt $maxVersion) {
				$maxVersion = $version
				$maxVersionFileName = $fileName
			}
		}
	}
	$maxVersionFileName
}

#delete current bepinex folder and related files
#get the latest modpack version file Name
#unzip the file to the current directory
$err = 0

Check-And-Delete-Bepinex-Related
$downloadsPath = Get-Download-Folder-Path
$BepinexPackZip = Get-BepinexPack -downloadsFolder $downloadsPath
if ($BepinexPackZip -eq "") {
	Write-Host "ERROR: Cannot find BepinexPack zip file" -ForegroundColor red
	$err = 1
	
} else {
	try {
		Unzip-Bepinex $downloadsPath $BepinexPackZip
	} catch {
		Write-Host "ERROR: Cannot unzip" -ForegroundColor red
		$err = 1
	}
}

if ($err -eq 0) {
	Write-Host "Success!" -ForegroundColor Green
} else {
	Write-Host "Something went wrong" -ForegroundColor red
}


#Read-Host -Prompt "Press Enter to exit"
if ($Host.Name -eq "ConsoleHost")
{
    Write-Host "Press any key to exit..."
    $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyUp") > $null
}
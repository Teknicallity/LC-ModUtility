# 2>NUL & @CLS & PUSHD "%~dp0" & "%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -nol -nop -ep bypass "[IO.File]::ReadAllText('%~f0')|iex" & POPD & EXIT /B

param(
    [switch]$testrun,
    [Alias("v")]
    [switch]$verbose
)

$testrun=$false
$verbose=$false

$directoryPath = ".\"

function Remove-Bepinex-Files {
	$bepinexPath = Join-Path -Path $directoryPath -ChildPath "Bepinex"
	$winhttpPath = Join-Path -Path $directoryPath -ChildPath "winhttp.dll"
	$doorstopPath = Join-Path -Path $directoryPath -ChildPath "doorstop_config.ini"

	# If bepinex exists, delete
	if (Test-Path -Path $bepinexPath -PathType Container){
		if (!$testrun) {Remove-Item -Path $bepinexPath -Recurse -Force}
		Write-Host "Bepinex folder deleted$(if ($verbose) {": '$bepinexPath'"})"
	} else{
		Write-Host "Bepinex directory does not exist. Creating directory$(if ($verbose) {": '$bepinexPath'"})"
	}

	if (Test-Path -Path $winhttpPath -PathType Leaf){
		if (!$testrun) {Remove-Item -Path $winhttpPath -Force}
		if ($verbose) {Write-Host "Deleted '$winhttpPath'"}
	}
	if (Test-Path -Path $doorstopPath -PathType Leaf){
		if (!$testrun) {Remove-Item -Path $doorstopPath -Force}
		if ($verbose) {Write-Host "Deleted '$doorstopPath'"}
	}
}

function Get-Download-Folder-Path {
	$userProfile = $env:USERPROFILE
	$downloadsFolder = Join-Path -Path $userProfile -ChildPath "Downloads"
	Write-Host "Downloads folder: '$downloadsFolder'"
	return $downloadsFolder
}

function Unzip-Bepinex($downloadsPath, $packName) {
	$zipInput = Join-Path -Path $downloadsPath -ChildPath $packName
	Write-Host "Unzipping: " -ForegroundColor Yellow -NoNewLine
	Write-Host "$(if ($verbose) {"'$zipInput'"} else {"'$packName'"})"
	#$zipOutput = Join-Path -Path $directoryPath -ChildPath $directoryPath
	if (!$testrun){
		Expand-Archive -Force -Path $zipInput -DestinationPath $directoryPath #$zipOutput
	}
	if ($verbose) {Write-Host "Expanded archive"}
}

function Get-BepinexPack {
	param (
		[string]$downloadsFolder
	)
	$pattern = '^BepinExPack_v(\d+(\.\d+)*)\.zip$'
	$matchingFiles = Get-ChildItem -Path $downloadsFolder | Where-Object { $_.Name -match $pattern }
	if ($verbose) {Write-Host $matchingFiles}
	return Find-Highest-Version $matchingFiles
}

function Find-Highest-Version($matchingFiles) {
	$maxVersion = 0
	$maxVersionFileName = ""

	if ($verbose) {Write-Host "Filenames:"}
	# Iterate through the file names
	foreach ($fileName in $matchingFiles) {
		if ($verbose) {Write-Host "`t'$fileName'"}
		if ($fileName -match $pattern) {
			if ($verbose) {Write-Host "`t`tMatched"}
			# if ($verbose) {Write-Host "`t`t$($matches | Out-String)"}
			$version = [float]$matches[1]

			# Check if the current version is greater than the maximum version
			if ($verbose) {Write-Host "`t`t$version > $maxVersion " -NoNewLine}
			if ($version -gt $maxVersion) {
				if ($verbose) {Write-Host "is true"}
				$maxVersion = $version
				$maxVersionFileName = $fileName
			} elseif ($verbose) {Write-Host "is false"}
		} elseif ($verbose) {Write-Host "`t`tNo match"}
	}
	if ($verbose){
		Write-Host "Max Version: '$maxVersion'"
		Write-Host "Max Version Filename: '$maxVersionFileName'"
	}
	return $maxVersionFileName
}

#delete current bepinex folder and related files
#get the latest modpack version file Name
#unzip the file to the current directory
$err = 0

Remove-Bepinex-Files
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
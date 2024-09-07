# 2>NUL & @CLS & PUSHD "%~dp0" & "%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -nol -nop -ep bypass "[IO.File]::ReadAllText('%~f0')|iex" & POPD & EXIT /B

param(
    [switch]$testrun,
    [Alias("v")]
    [switch]$verbose
)

$testrun=$false
$verbose=$false

function Get-SteamInstallPath {
    try {
        $regKey = "HKCU:\Software\Valve\Steam"
        $steamPath = (Get-ItemProperty -Path $regKey).SteamPath
        return $steamPath
    } catch {
        Write-Host "Steam installation not found." -ForegroundColor Red
        return $null
    }
}

function Get-SteamLibraryFolders {
    param (
        [string]$SteamPath
    )

    $libraryFile = Join-Path -Path $SteamPath -ChildPath "steamapps\libraryfolders.vdf"
    if (Test-Path $libraryFile) {
        $vdfContent = Get-Content $libraryFile -Raw
        $libraryFolders = @()

        # Extract each path manually
        $steamLibraryMatches = [regex]::Matches($vdfContent, '"\d+"\s*{\s*"path"\s*"([^"]+)"')
        foreach ($match in $steamLibraryMatches) {
            $libraryFolders += $match.Groups[1].Value
        }
		if ($verbose) {Write-Host "Steam Libary Paths: $libraryFolders"}
        return $libraryFolders
    } else {
        Write-Host "Library folders file not found." -ForegroundColor Red
        return @()
    }
}

function Find-LethalCompany {
    param (
        [string[]]$LibraryPaths
    )

    $lethalCompanyAppId = "1966720"
    foreach ($path in $LibraryPaths) {
        $appManifest = Join-Path -Path $path -ChildPath "steamapps\appmanifest_$lethalCompanyAppId.acf"
        if (Test-Path $appManifest) {
            $gamePath = Join-Path -Path $path -ChildPath "steamapps\common\Lethal Company"
			if ($verbose) {Write-Host "Base Steam Path: $steamPath"}
            return $gamePath
        }
    }
    return $null
}

function Get-LethalCompanyFolder {
	$steamPath = Get-SteamInstallPath
    if (-not $steamPath) {
        return $null
    }

    $libraryPaths = Get-SteamLibraryFolders -SteamPath $steamPath
    if ($libraryPaths.Count -eq 0) {
        Write-Host "No Steam library folders found." -ForegroundColor Red
        return $null
    }

    $gamePath = Find-LethalCompany -LibraryPaths $libraryPaths
    if ($gamePath) {
        Write-Host "Lethal Company found at: " -NoNewline
		Write-Host "$gamePath" -ForegroundColor DarkGreen
		return $gamePath
    }
	Write-Host "Lethal Company not found in any Steam library." -ForegroundColor Red
	return $null
}

function Remove-Bepinex-Files {
	$bepinexPath = Join-Path -Path $lethalCompanyPath -ChildPath "Bepinex"
	$winhttpPath = Join-Path -Path $lethalCompanyPath -ChildPath "winhttp.dll"
	$doorstopPath = Join-Path -Path $lethalCompanyPath -ChildPath "doorstop_config.ini"

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

function Get-User-Download-Folder-Path {
	$userProfile = $env:USERPROFILE
	$downloadsFolder = Join-Path -Path $userProfile -ChildPath "Downloads"
	Write-Host "Downloads folder: '$downloadsFolder'"
	return $downloadsFolder
}

function Unarchive-Bepinex-Pack($downloadsPath, $packName) {
	$zipInput = Join-Path -Path $downloadsPath -ChildPath $packName
	Write-Host "Unzipping: " -NoNewLine
	Write-Host "$(if ($verbose) {"'$zipInput'"} else {"'$packName'"})" -ForegroundColor Yellow
	Write-Host "`tinto: $lethalCompanyPath"
	#$zipOutput = Join-Path -Path $lethalCompanyPath -ChildPath $lethalCompanyPath
	if (!$testrun){
		Expand-Archive -Force -Path $zipInput -DestinationPath $lethalCompanyPath #$zipOutput
	}
	if ($verbose) {Write-Host "Expanded archive"}
}

function Get-BepinexPack-Archive-Name {
	param (
		[string]$downloadsFolder
	)
	$pattern = '^BepinExPack_v(\d+(\.\d+)*)\.zip$'
	$matchingFiles = Get-ChildItem -Path $downloadsFolder | Where-Object { $_.Name -match $pattern }
	if ($verbose) {Write-Host $matchingFiles}
	return Find-Highest-Version-Pack $matchingFiles
}

function Find-Highest-Version-Pack($matchingFiles) {
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



$lethalCompanyPath = Get-LethalCompanyFolder
$err = 0

Remove-Bepinex-Files
$downloadsPath = Get-User-Download-Folder-Path
$BepinexPackZip = Get-BepinexPack-Archive-Name -downloadsFolder $downloadsPath
if ($BepinexPackZip -eq "") {
	Write-Host "ERROR: Cannot find BepinexPack zip file" -ForegroundColor red
	$err = 1

} else {
	try {
		Unarchive-Bepinex-Pack $downloadsPath $BepinexPackZip
	} catch {
		Write-Host "ERROR: Cannot unzip" -ForegroundColor red
		$err = 1
	}
}

if ($err -eq 0) {
	Write-Host "Success!" -ForegroundColor DarkGreen
} else {
	Write-Host "Something went wrong" -ForegroundColor red
}


#Read-Host -Prompt "Press Enter to exit"
if ($Host.Name -eq "ConsoleHost"){
    Write-Host "Press any key to exit..."
    $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyUp") > $null
}
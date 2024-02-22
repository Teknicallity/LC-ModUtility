
# LC-ModUtility

## Description:

- A cli utility to help the manual management of Lethal Company mods from ThunderStore

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [Limitations](#limitations)
4. [Planned](#planned)
5. [License](#license)



## Installation

- Drop 'LethalModUtility.exe' into your Lethal Company game directory

## Usage

1. Run 'LethalModUtility.exe' through terminal or by double-clicking
2. Choose one of the options by entering the corresponding number or letter:

    - `1`: Update mods - Read through '.\BepInEx\plugins.md' and download out of date mods to '\$User\Downloads\LC_New_Mods'.
    - `2`: Unzip pack from downloads - Read any files named '\$User\Downloads\BepinExPack_v#.zip', 
   choose the newest version and unzip into your Lethal Company game directory.
    - `3`: Creating new compressed modpack - Zip '.\BepinEx\', '.\winhttp.dll', and '.\doorstop_config.ini' into '.\BepinExPack_v#.zip'.
    - `4`: Download new mod - Takes a ThunderStore mod link, downloads the mod zip to '\$User\Downloads\LC_New_Mods',
   and adds the version with link to '.\BepInEx\plugins.md'.
    - `q`: Quit - Choose this option to exit the program.
   
3. Follow the prompts or instructions provided by the program for each selected option.

## Limitations

- Whenever downloading a mod from ThunderStore, the zip is only put it '\$User\Downloads\LC_New_Mods', not extracted and 
installed to the game.

## Planned

- Auto parse and installation of different mods using common ThunderStore mod directory formats.

## License

- This project is licensed under the terms of the GPL v3, found in the [LICENSE](LICENSE) file.


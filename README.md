
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

- Drop 'LC-ModUtility_#.#.#.exe' into your Lethal Company game directory

## Usage

1. Run 'LC-ModUtility_#.#.#.exe' through terminal or by double-clicking
2. Choose one of the options by entering the corresponding number or letter:

    - `1`: Update mods - Read through '.\BepInEx\plugins.md', download, and install mods.
    - `2`: Unzip pack from downloads - Read any files named '\$User\Downloads\BepinExPack_v#.zip', 
   choose the newest version and unzip into the .exe's local directory.
    - `3`: Creating new compressed modpack - Zip '.\BepinEx\', '.\winhttp.dll', and '.\doorstop_config.ini' into 
   '.\BepinExPack_v#.zip'.
    - `4`: Download new mod - Takes a ThunderStore mod link, downloads it, installs, it and adds the version with 
   link to '.\BepInEx\plugins.md'.
    - `q`: Quit - Choose this option to exit the program and write the modified modlist to '.\BepInEx\plugins.md'.
   
3. Follow the prompts or instructions provided by the program for each selected option.

## Limitations

- Error recovery is not ideal

## Planned

- Install mods completely fresh from a specified plugins.md file.

## License

- This project is licensed under the terms of the GPL v3, found in the [LICENSE](LICENSE) file.


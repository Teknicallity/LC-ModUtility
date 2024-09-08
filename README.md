
# LC-ModUtility

## Description:

- A cli utility to help the manual management of Lethal Company mods from ThunderStore. 
- Produces a simple zip file with all mods inside to give to your friends, along with the installer script (optional).
- Why? The Thunderstore app requires the Overwolf launcher, and both are filled with obtrusive ads.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Limitations](#limitations)
- [Planned](#planned)
- [License](#license)



## Installation

- Drop 'LC-ModUtility_#.#.#.exe' into your Lethal Company game directory.

## Usage

1. Run 'LC-ModUtility_#.#.#.exe' through terminal or by double-clicking.
2. Choose one of the options by entering the corresponding number or letter:

    - `1`: Update mods - Read through '.\BepInEx\plugins.md', download, and install mods.
    - `2`: Unzip pack from downloads - Read any files named '\$User\Downloads\BepinExPack_v#.zip', 
   choose the newest version and unzip into the .exe's local directory.
    - `3`: Creating new compressed modpack - Zip '.\BepinEx\', '.\winhttp.dll', and '.\doorstop_config.ini' into 
   '.\BepinExPack_v#.zip'.
    - `4`: Download new mod - Takes a ThunderStore mod link, downloads it, installs, it and adds the version with 
   link to '.\BepInEx\plugins.md'.
    - `5`: Re-download all mods as a "fresh start." Gives the option to keep the old configs.
    - `q`: Quit - Choose this option to exit the program and write the modified modlist to '.\BepInEx\plugins.md'.
   
3. Follow the prompts or instructions provided by the program for each selected option.

When finished modifying a pack, and after creating a new zip, you can send the zip file to whomever you want for them 
to unzip and install into the Lethal Company directory. For their convenience, and possibly to avoid headaches for you,
[an installation script](./Update-Lethal-Company-Modpack.bat) has been provided that will find their game directory and
install the modpack for them.

## Limitations

- Error recovery is not ideal

## Planned

- Git support for plugins.md and maybe config files
- Storing what files belong to what mod
- Easy mod removal
- Easier config management

## License

- This project is licensed under the terms of the GPL v3, found in the [LICENSE](LICENSE) file.


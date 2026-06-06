#ifndef MyAppVersion
  #define MyAppVersion "dev"
#endif

[Setup]
AppName=Capture Board Viewer
AppVersion={#MyAppVersion}
AppPublisher=GraftTeam
AppPublisherURL=https://github.com/PCloud63514/capture-board-viewer
AppSupportURL=https://github.com/PCloud63514/capture-board-viewer/issues
DefaultDirName={localappdata}\capture-board-viewer
OutputDir=dist
OutputBaseFilename=capture-board-viewer-setup-{#MyAppVersion}
PrivilegesRequired=lowest
DisableProgramGroupPage=yes
DisableDirPage=yes
Compression=lzma
SolidCompression=yes
WizardStyle=modern
UninstallDisplayName=Capture Board Viewer
UninstallDisplayIcon={app}\capture-board-selector.exe

[Languages]
Name: "korean"; MessagesFile: "compiler:Languages\Korean.isl"

[Files]
Source: "capture-board-selector.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{autodesktop}\Capture Board Viewer"; Filename: "{app}\capture-board-selector.exe"
Name: "{autoprograms}\Capture Board Viewer\Capture Board Viewer"; Filename: "{app}\capture-board-selector.exe"
Name: "{autoprograms}\Capture Board Viewer\제거"; Filename: "{uninstallexe}"

[Run]
Filename: "{app}\capture-board-selector.exe"; Description: "지금 실행하기"; Flags: nowait postinstall skipifsilent

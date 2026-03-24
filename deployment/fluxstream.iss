#define MyAppName "FluxStream"
#define MyAppVersion "2.0.0"
#define MyAppPublisher "FluxStream"
#define MyAppExeName "fluxstream.exe"
#ifndef SourceDir
  #define SourceDir "C:\xampp\htdocs\stream\dist\fluxstream-windows-amd64-service"
#endif
#ifndef OutputDir
  #define OutputDir "C:\xampp\htdocs\stream\dist"
#endif
#ifndef SetupIconPath
  #define SetupIconPath "C:\xampp\htdocs\stream\deployment\fluxstream.ico"
#endif

[Setup]
AppId={{8A7F7A82-8D71-4BA1-90A0-E34D2BFEA5C7}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\FluxStream
DefaultGroupName=FluxStream
PrivilegesRequired=admin
OutputDir={#OutputDir}
OutputBaseFilename=FluxStream-Setup
Compression=lzma
SolidCompression=yes
WizardStyle=modern
ArchitecturesInstallIn64BitMode=x64compatible
DisableProgramGroupPage=yes
SetupIconFile={#SetupIconPath}
UninstallDisplayIcon={app}\fluxstream.ico

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "german"; MessagesFile: "compiler:Languages\German.isl"
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"
Name: "french"; MessagesFile: "compiler:Languages\French.isl"
Name: "turkish"; MessagesFile: "compiler:Languages\Turkish.isl"

[CustomMessages]
english.DesktopIcon=Desktop shortcut
german.DesktopIcon=Desktop shortcut
spanish.DesktopIcon=Acceso directo en el escritorio
french.DesktopIcon=Raccourci sur le bureau
turkish.DesktopIcon=Masaustu kisayolu
english.LaunchApp=Launch FluxStream after setup
german.LaunchApp=FluxStream nach der Installation starten
spanish.LaunchApp=Iniciar FluxStream despues de la instalacion
french.LaunchApp=Lancer FluxStream apres l'installation
turkish.LaunchApp=Kurulumdan sonra FluxStream'i baslat

[Tasks]
Name: "desktopicon"; Description: "{cm:DesktopIcon}"; Flags: checkedonce

[Files]
Source: "{#SourceDir}\fluxstream.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#SourceDir}\fluxstream.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#SourceDir}\ffmpeg\*"; DestDir: "{app}\ffmpeg"; Flags: ignoreversion recursesubdirs createallsubdirs

[UninstallDelete]
Type: filesandordirs; Name: "{app}\data"
Type: filesandordirs; Name: "{app}\ffmpeg"
Type: files; Name: "{app}\fluxstream.exe"
Type: files; Name: "{app}\fluxstream.ico"
Type: dirifempty; Name: "{app}"

[Icons]
Name: "{group}\FluxStream"; Filename: "{app}\fluxstream.exe"; IconFilename: "{app}\fluxstream.ico"
Name: "{autodesktop}\FluxStream"; Filename: "{app}\fluxstream.exe"; Tasks: desktopicon; IconFilename: "{app}\fluxstream.ico"

[Run]
Filename: "{app}\fluxstream.exe"; Description: "{cm:LaunchApp}"; Flags: nowait postinstall skipifsilent; Check: ShouldLaunchDesktopMode

[UninstallRun]
Filename: "{app}\fluxstream.exe"; Parameters: "service stop"; Flags: runhidden waituntilterminated skipifdoesntexist; RunOnceId: "FluxStreamServiceStop"
Filename: "{app}\fluxstream.exe"; Parameters: "service uninstall"; Flags: runhidden waituntilterminated skipifdoesntexist; RunOnceId: "FluxStreamServiceUninstall"

[Code]
var
  ModePage: TInputOptionWizardPage;
  NetworkPage: TWizardPage;
  SSLPage: TWizardPage;
  UseCustomPortsCheck: TNewCheckBox;
  HttpPortLabel: TNewStaticText;
  HttpPortEdit: TNewEdit;
  RtmpPortLabel: TNewStaticText;
  RtmpPortEdit: TNewEdit;
  HttpsPortLabel: TNewStaticText;
  HttpsPortEdit: TNewEdit;
  RtmpsPortLabel: TNewStaticText;
  RtmpsPortEdit: TNewEdit;
  DomainLabel: TNewStaticText;
  DomainEdit: TNewEdit;
  PreloadSSLCheck: TNewCheckBox;
  EnableSSLNowCheck: TNewCheckBox;
  CertPathLabel: TNewStaticText;
  CertPathEdit: TNewEdit;
  CertBrowseButton: TNewButton;
  KeyPathLabel: TNewStaticText;
  KeyPathEdit: TNewEdit;
  KeyBrowseButton: TNewButton;

function InstallerLanguageCode(): String;
begin
  if ActiveLanguage = 'english' then
    Result := 'en'
  else if ActiveLanguage = 'german' then
    Result := 'de'
  else if ActiveLanguage = 'spanish' then
    Result := 'es'
  else if ActiveLanguage = 'french' then
    Result := 'fr'
  else
    Result := 'tr';
end;

function L(const TurkishText, EnglishText, GermanText, SpanishText, FrenchText: String): String;
begin
  if ActiveLanguage = 'english' then
    Result := EnglishText
  else if ActiveLanguage = 'german' then
    Result := GermanText
  else if ActiveLanguage = 'spanish' then
    Result := SpanishText
  else if ActiveLanguage = 'french' then
    Result := FrenchText
  else
    Result := TurkishText;
end;

function BoolToConfig(Value: Boolean): String;
begin
  if Value then
    Result := 'true'
  else
    Result := 'false';
end;

procedure AddConfigArg(var Params: String; const Key, Value: String);
begin
  if Params <> '' then
    Params := Params + ' ';
  Params := Params + '"' + Key + '=' + Value + '"';
end;

function ShouldInstallService(): Boolean;
begin
  Result := ModePage.Values[0];
end;

function ShouldLaunchDesktopMode(): Boolean;
begin
  Result := not ShouldInstallService();
end;

procedure UpdatePortPageState(Sender: TObject);
begin
  HttpPortEdit.Enabled := UseCustomPortsCheck.Checked;
  RtmpPortEdit.Enabled := UseCustomPortsCheck.Checked;
  HttpsPortEdit.Enabled := UseCustomPortsCheck.Checked;
  RtmpsPortEdit.Enabled := UseCustomPortsCheck.Checked;
  HttpPortLabel.Enabled := UseCustomPortsCheck.Checked;
  RtmpPortLabel.Enabled := UseCustomPortsCheck.Checked;
  HttpsPortLabel.Enabled := UseCustomPortsCheck.Checked;
  RtmpsPortLabel.Enabled := UseCustomPortsCheck.Checked;
end;

procedure UpdateSSLPageState(Sender: TObject);
begin
  CertPathEdit.Enabled := PreloadSSLCheck.Checked;
  CertBrowseButton.Enabled := PreloadSSLCheck.Checked;
  KeyPathEdit.Enabled := PreloadSSLCheck.Checked;
  KeyBrowseButton.Enabled := PreloadSSLCheck.Checked;
  CertPathLabel.Enabled := PreloadSSLCheck.Checked;
  KeyPathLabel.Enabled := PreloadSSLCheck.Checked;
  EnableSSLNowCheck.Enabled := PreloadSSLCheck.Checked;
end;

procedure BrowseForCert(Sender: TObject);
var
  SelectedFile: String;
begin
  SelectedFile := CertPathEdit.Text;
  if GetOpenFileName(L('CRT dosyasi secin','Select CRT file','CRT-Datei waehlen','Seleccione el archivo CRT','Selectionner le fichier CRT'), SelectedFile, '', 'Certificate Files|*.crt;*.pem|Tum Dosyalar|*.*', '') then
    CertPathEdit.Text := SelectedFile;
end;

procedure BrowseForKey(Sender: TObject);
var
  SelectedFile: String;
begin
  SelectedFile := KeyPathEdit.Text;
  if GetOpenFileName(L('KEY dosyasi secin','Select KEY file','KEY-Datei waehlen','Seleccione el archivo KEY','Selectionner le fichier KEY'), SelectedFile, '', 'Key Files|*.key;*.pem|Tum Dosyalar|*.*', '') then
    KeyPathEdit.Text := SelectedFile;
end;

function ValidatePortValue(const Caption, Value: String): Boolean;
var
  Port: Integer;
begin
  Port := StrToIntDef(Value, -1);
  Result := (Port > 0) and (Port <= 65535);
  if not Result then
    MsgBox(Caption + L(' icin 1-65535 araliginda gecerli bir port girin.',' must be a valid port between 1 and 65535.',' muss ein gueltiger Port zwischen 1 und 65535 sein.',' debe ser un puerto valido entre 1 y 65535.',' doit etre un port valide entre 1 et 65535.'), mbError, MB_OK);
end;

function ServiceExists(): Boolean;
var
  ResultCode: Integer;
begin
  Result := Exec(ExpandConstant('{sys}\sc.exe'), 'query FluxStream', '', SW_HIDE, ewWaitUntilTerminated, ResultCode) and (ResultCode = 0);
end;

function RunInstalledCommand(const Params, ErrorMessage: String): Boolean;
var
  ResultCode: Integer;
begin
  Result := Exec(ExpandConstant('{app}\fluxstream.exe'), Params, '', SW_HIDE, ewWaitUntilTerminated, ResultCode) and (ResultCode = 0);
  if not Result then
    RaiseException(ErrorMessage);
end;

procedure ApplyInstallerSettings();
var
  Params: String;
  CertTarget: String;
  KeyTarget: String;
begin
  Params := '';
  AddConfigArg(Params, 'language', InstallerLanguageCode());

  if UseCustomPortsCheck.Checked then
  begin
    AddConfigArg(Params, 'http_port', HttpPortEdit.Text);
    AddConfigArg(Params, 'embed_http_port', HttpPortEdit.Text);
    AddConfigArg(Params, 'rtmp_port', RtmpPortEdit.Text);
    AddConfigArg(Params, 'https_port', HttpsPortEdit.Text);
    AddConfigArg(Params, 'embed_https_port', HttpsPortEdit.Text);
    AddConfigArg(Params, 'rtmps_port', RtmpsPortEdit.Text);
  end;

  if Trim(DomainEdit.Text) <> '' then
    AddConfigArg(Params, 'embed_domain', Trim(DomainEdit.Text));

  if PreloadSSLCheck.Checked then
  begin
    ForceDirectories(ExpandConstant('{app}\data\certs'));
    CertTarget := ExpandConstant('{app}\data\certs\server.crt');
    KeyTarget := ExpandConstant('{app}\data\certs\server.key');
    if not CopyFile(CertPathEdit.Text, CertTarget, False) then
      RaiseException('SSL sertifika dosyasi kopyalanamadi.');
    if not CopyFile(KeyPathEdit.Text, KeyTarget, False) then
      RaiseException('SSL key dosyasi kopyalanamadi.');

    AddConfigArg(Params, 'ssl_cert_path', CertTarget);
    AddConfigArg(Params, 'ssl_key_path', KeyTarget);

    if EnableSSLNowCheck.Checked then
    begin
      AddConfigArg(Params, 'ssl_enabled', BoolToConfig(True));
      AddConfigArg(Params, 'rtmps_enabled', BoolToConfig(True));
      AddConfigArg(Params, 'embed_use_https', BoolToConfig(True));
    end;
  end;

  if Params <> '' then
    RunInstalledCommand('config set ' + Params, 'Kurulum secenekleri uygulanamadi.');
end;

procedure ConfigureInstallMode();
var
  ResultCode: Integer;
begin
  if ServiceExists() then
  begin
    Exec(ExpandConstant('{app}\fluxstream.exe'), 'service stop', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec(ExpandConstant('{app}\fluxstream.exe'), 'service uninstall', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
  end;

  if ShouldInstallService() then
  begin
    RunInstalledCommand('service install', 'FluxStream servisi kurulamadi.');
    RunInstalledCommand('service start', 'FluxStream servisi baslatilamadi.');
  end;
end;

procedure InitializeWizard();
var
  NoteLabel: TNewStaticText;
  CurrentTop: Integer;
begin
  ModePage := CreateInputOptionPage(
    wpSelectTasks,
    L('Kurulum Modu','Setup Mode','Installationsmodus','Modo de instalacion','Mode d''installation'),
    L('FluxStream''i nasil calistirmak istediginizi secin','Choose how you want to run FluxStream','Waehlen Sie, wie FluxStream ausgefuehrt werden soll','Elija como desea ejecutar FluxStream','Choisissez comment vous souhaitez executer FluxStream'),
    L('Bu secim istege baglidir. Kurulumdan sonra admin panelinden, Windows service araclarindan veya yeniden kurulumla degistirebilirsiniz.','This choice is optional. You can change it later from the admin panel, Windows service tools, or by reinstalling.','Diese Auswahl ist optional. Sie koennen sie spaeter im Admin-Panel, mit Windows-Diensten oder durch Neuinstallation aendern.','Esta opcion es opcional. Puede cambiarla mas tarde desde el panel, las herramientas de servicio de Windows o reinstalando.','Ce choix est facultatif. Vous pourrez le modifier plus tard depuis le panneau d''administration, les outils de service Windows ou via une reinstallation.'),
    True,
    False);
  ModePage.Add(L('Windows Service olarak kur ve otomatik baslat (onerilen)','Install as a Windows Service and start automatically (recommended)','Als Windows-Dienst installieren und automatisch starten (empfohlen)','Instalar como servicio de Windows e iniciar automaticamente (recomendado)','Installer comme service Windows et demarrer automatiquement (recommande)'));
  ModePage.Add(L('Servis kurma, uygulamayi elle/konsoldan baslat','Do not install a service, start the app manually / from console','Keinen Dienst installieren, Anwendung manuell / ueber die Konsole starten','No instalar servicio, iniciar la aplicacion manualmente / desde consola','Ne pas installer de service, lancer l''application manuellement / depuis la console'));
  ModePage.Values[0] := True;

  NetworkPage := CreateCustomPage(
    ModePage.ID,
    L('Ag Ayarlari (Opsiyonel)','Network Settings (Optional)','Netzwerkeinstellungen (optional)','Configuracion de red (opcional)','Parametres reseau (facultatif)'),
    L('Isterseniz varsayilan portlari kurulum sirasinda ozellestirebilirsiniz. Bu ayarlar kurulumdan sonra admin panelinden degistirilebilir.','You can customize the default ports during setup if needed. These settings can also be changed later from the admin panel.','Sie koennen die Standardports waehrend der Installation anpassen. Diese Einstellungen lassen sich spaeter auch im Admin-Panel aendern.','Puede personalizar los puertos predeterminados durante la instalacion. Tambien podra cambiarlos despues desde el panel.','Vous pouvez personnaliser les ports par defaut pendant l''installation. Ces parametres pourront aussi etre modifies ensuite depuis le panneau.'));

  NoteLabel := TNewStaticText.Create(NetworkPage);
  NoteLabel.Parent := NetworkPage.Surface;
  NoteLabel.Left := 0;
  NoteLabel.Top := 0;
  NoteLabel.Width := NetworkPage.SurfaceWidth;
  NoteLabel.Height := ScaleY(42);
  NoteLabel.WordWrap := True;
  NoteLabel.Caption := L('Bu sayfadaki ayarlar istege baglidir. Domain ve web portlari embed/player linklerinde kullanilir. RTMP ve RTMPS portlari yayin gonderimi icin ayridir. Tum ayarlar kurulumdan sonra panelden degistirilebilir.','These settings are optional. Domain and web ports are used in embed/player links. RTMP and RTMPS ports are dedicated to ingest. All settings can be changed later from the panel.','Diese Einstellungen sind optional. Domain und Web-Ports werden in Embed-/Player-Links verwendet. RTMP- und RTMPS-Ports sind fuer den Ingest vorgesehen. Alle Einstellungen koennen spaeter im Panel geaendert werden.','Estas opciones son opcionales. El dominio y los puertos web se usan en los enlaces de embed/player. Los puertos RTMP y RTMPS se usan para la ingesta. Todo puede cambiarse despues desde el panel.','Ces reglages sont facultatifs. Le domaine et les ports web sont utilises dans les liens embed/player. Les ports RTMP et RTMPS sont dedies a l''ingestion. Tous les reglages pourront etre modifies ensuite depuis le panneau.');

  UseCustomPortsCheck := TNewCheckBox.Create(NetworkPage);
  UseCustomPortsCheck.Parent := NetworkPage.Surface;
  UseCustomPortsCheck.Left := 0;
  UseCustomPortsCheck.Top := NoteLabel.Top + NoteLabel.Height + ScaleY(8);
  UseCustomPortsCheck.Width := NetworkPage.SurfaceWidth;
  UseCustomPortsCheck.Caption := L('Kurulum sirasinda ozel portlari uygula','Apply custom ports during setup','Benutzerdefinierte Ports waehrend der Installation anwenden','Aplicar puertos personalizados durante la instalacion','Appliquer des ports personnalises pendant l''installation');
  UseCustomPortsCheck.OnClick := @UpdatePortPageState;

  CurrentTop := UseCustomPortsCheck.Top + ScaleY(32);

  DomainLabel := TNewStaticText.Create(NetworkPage);
  DomainLabel.Parent := NetworkPage.Surface;
  DomainLabel.Left := 0;
  DomainLabel.Top := CurrentTop;
  DomainLabel.Caption := L('Public Domain / IP','Public Domain / IP','Oeffentliche Domain / IP','Dominio publico / IP','Domaine public / IP');

  DomainEdit := TNewEdit.Create(NetworkPage);
  DomainEdit.Parent := NetworkPage.Surface;
  DomainEdit.Left := ScaleX(180);
  DomainEdit.Top := CurrentTop - ScaleY(4);
  DomainEdit.Width := ScaleX(260);
  DomainEdit.Text := '';

  CurrentTop := CurrentTop + ScaleY(32);

  HttpPortLabel := TNewStaticText.Create(NetworkPage);
  HttpPortLabel.Parent := NetworkPage.Surface;
  HttpPortLabel.Left := 0;
  HttpPortLabel.Top := CurrentTop;
  HttpPortLabel.Caption := 'HTTP Port';

  HttpPortEdit := TNewEdit.Create(NetworkPage);
  HttpPortEdit.Parent := NetworkPage.Surface;
  HttpPortEdit.Left := ScaleX(180);
  HttpPortEdit.Top := CurrentTop - ScaleY(4);
  HttpPortEdit.Width := ScaleX(120);
  HttpPortEdit.Text := '8844';

  CurrentTop := CurrentTop + ScaleY(28);
  RtmpPortLabel := TNewStaticText.Create(NetworkPage);
  RtmpPortLabel.Parent := NetworkPage.Surface;
  RtmpPortLabel.Left := 0;
  RtmpPortLabel.Top := CurrentTop;
  RtmpPortLabel.Caption := 'RTMP Port';

  RtmpPortEdit := TNewEdit.Create(NetworkPage);
  RtmpPortEdit.Parent := NetworkPage.Surface;
  RtmpPortEdit.Left := ScaleX(180);
  RtmpPortEdit.Top := CurrentTop - ScaleY(4);
  RtmpPortEdit.Width := ScaleX(120);
  RtmpPortEdit.Text := '1935';

  CurrentTop := CurrentTop + ScaleY(28);
  HttpsPortLabel := TNewStaticText.Create(NetworkPage);
  HttpsPortLabel.Parent := NetworkPage.Surface;
  HttpsPortLabel.Left := 0;
  HttpsPortLabel.Top := CurrentTop;
  HttpsPortLabel.Caption := 'HTTPS Port';

  HttpsPortEdit := TNewEdit.Create(NetworkPage);
  HttpsPortEdit.Parent := NetworkPage.Surface;
  HttpsPortEdit.Left := ScaleX(180);
  HttpsPortEdit.Top := CurrentTop - ScaleY(4);
  HttpsPortEdit.Width := ScaleX(120);
  HttpsPortEdit.Text := '443';

  CurrentTop := CurrentTop + ScaleY(28);
  RtmpsPortLabel := TNewStaticText.Create(NetworkPage);
  RtmpsPortLabel.Parent := NetworkPage.Surface;
  RtmpsPortLabel.Left := 0;
  RtmpsPortLabel.Top := CurrentTop;
  RtmpsPortLabel.Caption := 'RTMPS Port';

  RtmpsPortEdit := TNewEdit.Create(NetworkPage);
  RtmpsPortEdit.Parent := NetworkPage.Surface;
  RtmpsPortEdit.Left := ScaleX(180);
  RtmpsPortEdit.Top := CurrentTop - ScaleY(4);
  RtmpsPortEdit.Width := ScaleX(120);
  RtmpsPortEdit.Text := '1936';

  SSLPage := CreateCustomPage(
    NetworkPage.ID,
    L('SSL On Yukleme (Opsiyonel)','SSL Preload (Optional)','SSL-Vorladen (optional)','Precarga SSL (opcional)','Prechargement SSL (facultatif)'),
    L('CRT ve KEY dosyalarini simdiden ekleyebilirsiniz. Istemezseniz kurulumdan sonra admin panelinden yukleyebilirsiniz.','You can add CRT and KEY files now. If you prefer, you can upload them later from the admin panel.','Sie koennen CRT- und KEY-Dateien jetzt hinzufuegen. Alternativ lassen sie sich spaeter ueber das Admin-Panel hochladen.','Puede anadir ahora los archivos CRT y KEY. Si lo prefiere, podra cargarlos despues desde el panel.','Vous pouvez ajouter les fichiers CRT et KEY maintenant. Si vous preferez, vous pourrez les televerser plus tard depuis le panneau.'));

  NoteLabel := TNewStaticText.Create(SSLPage);
  NoteLabel.Parent := SSLPage.Surface;
  NoteLabel.Left := 0;
  NoteLabel.Top := 0;
  NoteLabel.Width := SSLPage.SurfaceWidth;
  NoteLabel.Height := ScaleY(42);
  NoteLabel.WordWrap := True;
  NoteLabel.Caption := L('SSL adimi istege baglidir. Secilen sertifikalar data/certs klasorune kopyalanir. Kurulumdan sonra panelden degistirebilirsiniz.','The SSL step is optional. Selected certificates are copied into data/certs and can be changed later from the admin panel.','Der SSL-Schritt ist optional. Ausgewaehlte Zertifikate werden nach data/certs kopiert und koennen spaeter im Admin-Panel geaendert werden.','El paso SSL es opcional. Los certificados elegidos se copiaran a data/certs y podran cambiarse despues desde el panel.','L''etape SSL est facultative. Les certificats choisis seront copies dans data/certs et pourront etre modifies plus tard depuis le panneau.');

  PreloadSSLCheck := TNewCheckBox.Create(SSLPage);
  PreloadSSLCheck.Parent := SSLPage.Surface;
  PreloadSSLCheck.Left := 0;
  PreloadSSLCheck.Top := NoteLabel.Top + NoteLabel.Height + ScaleY(8);
  PreloadSSLCheck.Width := SSLPage.SurfaceWidth;
  PreloadSSLCheck.Caption := L('Kurulum sirasinda CRT ve KEY dosyalarini kopyala','Copy CRT and KEY files during setup','CRT- und KEY-Dateien waehrend der Installation kopieren','Copiar archivos CRT y KEY durante la instalacion','Copier les fichiers CRT et KEY pendant l''installation');
  PreloadSSLCheck.OnClick := @UpdateSSLPageState;

  EnableSSLNowCheck := TNewCheckBox.Create(SSLPage);
  EnableSSLNowCheck.Parent := SSLPage.Surface;
  EnableSSLNowCheck.Left := ScaleX(16);
  EnableSSLNowCheck.Top := PreloadSSLCheck.Top + ScaleY(24);
  EnableSSLNowCheck.Width := SSLPage.SurfaceWidth - ScaleX(16);
  EnableSSLNowCheck.Caption := L('Kurulumdan sonra HTTPS ve RTMPS''i aktif et','Enable HTTPS and RTMPS after setup','HTTPS und RTMPS nach der Installation aktivieren','Activar HTTPS y RTMPS despues de la instalacion','Activer HTTPS et RTMPS apres l''installation');
  EnableSSLNowCheck.Checked := True;

  CurrentTop := EnableSSLNowCheck.Top + ScaleY(32);

  CertPathLabel := TNewStaticText.Create(SSLPage);
  CertPathLabel.Parent := SSLPage.Surface;
  CertPathLabel.Left := 0;
  CertPathLabel.Top := CurrentTop;
  CertPathLabel.Caption := 'CRT';

  CertPathEdit := TNewEdit.Create(SSLPage);
  CertPathEdit.Parent := SSLPage.Surface;
  CertPathEdit.Left := 0;
  CertPathEdit.Top := CurrentTop + ScaleY(16);
  CertPathEdit.Width := SSLPage.SurfaceWidth - ScaleX(96);

  CertBrowseButton := TNewButton.Create(SSLPage);
  CertBrowseButton.Parent := SSLPage.Surface;
  CertBrowseButton.Left := CertPathEdit.Left + CertPathEdit.Width + ScaleX(8);
  CertBrowseButton.Top := CertPathEdit.Top - ScaleY(1);
  CertBrowseButton.Width := ScaleX(80);
  CertBrowseButton.Caption := L('Sec...','Browse...','Auswaehlen...','Buscar...','Parcourir...');
  CertBrowseButton.OnClick := @BrowseForCert;

  CurrentTop := CertPathEdit.Top + ScaleY(34);
  KeyPathLabel := TNewStaticText.Create(SSLPage);
  KeyPathLabel.Parent := SSLPage.Surface;
  KeyPathLabel.Left := 0;
  KeyPathLabel.Top := CurrentTop;
  KeyPathLabel.Caption := 'KEY';

  KeyPathEdit := TNewEdit.Create(SSLPage);
  KeyPathEdit.Parent := SSLPage.Surface;
  KeyPathEdit.Left := 0;
  KeyPathEdit.Top := CurrentTop + ScaleY(16);
  KeyPathEdit.Width := SSLPage.SurfaceWidth - ScaleX(96);

  KeyBrowseButton := TNewButton.Create(SSLPage);
  KeyBrowseButton.Parent := SSLPage.Surface;
  KeyBrowseButton.Left := KeyPathEdit.Left + KeyPathEdit.Width + ScaleX(8);
  KeyBrowseButton.Top := KeyPathEdit.Top - ScaleY(1);
  KeyBrowseButton.Width := ScaleX(80);
  KeyBrowseButton.Caption := L('Sec...','Browse...','Auswaehlen...','Buscar...','Parcourir...');
  KeyBrowseButton.OnClick := @BrowseForKey;

  UpdatePortPageState(nil);
  UpdateSSLPageState(nil);
end;

function NextButtonClick(CurPageID: Integer): Boolean;
begin
  Result := True;

  if CurPageID = NetworkPage.ID then
  begin
    if UseCustomPortsCheck.Checked then
      Result :=
        ValidatePortValue('HTTP Port', HttpPortEdit.Text) and
        ValidatePortValue('RTMP Port', RtmpPortEdit.Text) and
        ValidatePortValue('HTTPS Port', HttpsPortEdit.Text) and
        ValidatePortValue('RTMPS Port', RtmpsPortEdit.Text);
  end;

  if Result and (CurPageID = SSLPage.ID) and PreloadSSLCheck.Checked then
  begin
    if (Trim(CertPathEdit.Text) = '') or (not FileExists(CertPathEdit.Text)) then
    begin
      MsgBox('Gecerli bir CRT dosyasi secin veya SSL on yuklemeyi kapatin.', mbError, MB_OK);
      Result := False;
    end
    else if (Trim(KeyPathEdit.Text) = '') or (not FileExists(KeyPathEdit.Text)) then
    begin
      MsgBox('Gecerli bir KEY dosyasi secin veya SSL on yuklemeyi kapatin.', mbError, MB_OK);
      Result := False;
    end;
  end;
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    ApplyInstallerSettings();
    ConfigureInstallMode();
  end;
end;

function InitializeSetup(): Boolean;
begin
  Result := True;
end;
